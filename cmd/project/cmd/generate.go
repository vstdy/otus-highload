package cmd

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy/otus-highload/cmd/project/cmd/common"
	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg/logging"
)

const (
	flagUsersFilePath  = "users_file_path"
	flagUsersBatchSize = "users_batch_size"
	flagPostsFilePath  = "posts_file_path"
)

type generateCmdFlags struct {
	UsersFilePath  string `mapstructure:"users_file_path"`
	UsersBatchSize uint32 `mapstructure:"users_batch_size"`
	PostsFilePath  string `mapstructure:"posts_file_path"`
}

func defaultFlagsValues() generateCmdFlags {
	return generateCmdFlags{
		UsersFilePath:  "./build/resources/csv/people.csv",
		UsersBatchSize: 100000,
		PostsFilePath:  "./build/resources/txt/posts.txt",
	}
}

// newMigrateCmd creates a new migration command.
func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate operations",
	}

	cmd.AddCommand(generateUsers())
	cmd.AddCommand(generateFriends())
	cmd.AddCommand(generatePosts())

	return cmd
}

// generateUsers ...
func generateUsers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Generate users",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			ctx, logger := logging.GetCtxLogger(context.Background(), logging.WithLogLevel(config.LogLevel))

			st, err := config.BuildStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					logger.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			flagsValues, err := getFlags()
			if err != nil {
				return err
			}

			file, err := os.Open(flagsValues.UsersFilePath)
			if err != nil {
				return fmt.Errorf("opening file '%s': %w", flagsValues.UsersFilePath, err)
			}
			defer file.Close()

			csvReader := csv.NewReader(file)

			password, err := model.User{Password: "password"}.EncryptPassword()
			if err != nil {
				return fmt.Errorf("encrypting password: %w", err)
			}

			users := make([]model.User, 0, flagsValues.UsersBatchSize)
			for i := uint32(1); ; i++ {
				record, errN := csvReader.Read()
				if err = errN; err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("parsing file as CSV: %w", err)
				}

				name := strings.Split(record[0], " ")
				age, errN := strconv.ParseUint(record[1], 10, 8)
				if err = errN; err != nil {
					return fmt.Errorf("parsing age: %w", err)
				}

				user := model.User{
					FirstName:  name[1],
					SecondName: name[0],
					Age:        uint8(age),
					City:       record[2],
					Password:   password,
				}
				users = append(users, user)

				if i >= flagsValues.UsersBatchSize {
					if _, err = st.CopyUsers(ctx, users); err != nil {
						return fmt.Errorf("copying users: %w", err)
					}

					fmt.Print(".")
					i = 0
					users = make([]model.User, 0, flagsValues.UsersBatchSize)
				}
			}
			if len(users) > 0 {
				if _, err = st.CopyUsers(ctx, users); err != nil {
					return fmt.Errorf("copying users: %w", err)
				}
			}

			return nil
		},
	}

	config := defaultFlagsValues()
	cmd.Flags().StringP(flagUsersFilePath, "f", config.UsersFilePath, "Path to file")
	cmd.Flags().Uint32P(flagUsersBatchSize, "s", config.UsersBatchSize, "Size of the batch")

	return cmd
}

// generateFriends ...
func generateFriends() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "friends",
		Short: "Generate friends",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			ctx, logger := logging.GetCtxLogger(context.Background(), logging.WithLogLevel(config.LogLevel))

			st, err := config.BuildStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					logger.Error().Err(err).Msg("Shutting down the app")
				}
			}()
			//todo make configurable
			offset, userID := int64(1), int64(1)
			limit := int64(100)
			friendsNumber := int64(5)

			friends := make([]model.Friend, 0, limit*friendsNumber)

			for ; userID < offset+limit; userID++ {
				var friendID int64
				friendIDs := make(map[int64]struct{}, friendsNumber)

				for i := int64(0); i < friendsNumber; i++ {
					friendID = offset + rand.Int63n(limit)
					for _, ok := friendIDs[friendID]; friendID == userID || ok; {
						friendID = offset + rand.Int63n(limit)
						_, ok = friendIDs[friendID]
					}
					friendIDs[friendID] = struct{}{}

					friend := model.Friend{
						UserID:   userID,
						FriendID: friendID,
					}
					friends = append(friends, friend)
				}
			}

			if _, err = st.CopyFriends(ctx, friends); err != nil {
				return fmt.Errorf("copying friends: %w", err)
			}

			return nil
		},
	}

	return cmd
}

// generatePosts ...
func generatePosts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "posts",
		Short: "Generate posts",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			ctx, logger := logging.GetCtxLogger(context.Background(), logging.WithLogLevel(config.LogLevel))

			st, err := config.BuildStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					logger.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			flagsValues, err := getFlags()
			if err != nil {
				return err
			}
			//todo make configurable
			offset := int64(1)
			limit := int64(100)
			batchSize := 100

			file, err := os.Open(flagsValues.PostsFilePath)
			if err != nil {
				return fmt.Errorf("opening file '%s': %w", flagsValues.PostsFilePath, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)

			posts := make([]model.Post, 0, batchSize)
			for i := 1; scanner.Scan(); i++ {
				post := model.Post{
					Text:     scanner.Text(),
					AuthorID: offset + rand.Int63n(limit),
				}

				posts = append(posts, post)

				if i != batchSize {
					continue
				}
				if _, err = st.CopyPosts(ctx, posts); err != nil {
					return fmt.Errorf("copying posts: %w", err)
				}
				i = 0
				posts = make([]model.Post, 0, batchSize)
			}

			return nil
		},
	}

	config := defaultFlagsValues()
	cmd.Flags().StringP(flagPostsFilePath, "p", config.PostsFilePath, "Size of the batch")

	return cmd
}

func getFlags() (generateCmdFlags, error) {
	config := defaultFlagsValues()
	if err := viper.Unmarshal(&config); err != nil {
		return generateCmdFlags{}, fmt.Errorf("generateCmdFlags unmarshal: %w", err)
	}

	return config, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}
