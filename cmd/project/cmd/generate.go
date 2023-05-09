package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy/otus-highload/cmd/project/cmd/common"
	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg/logging"
)

const (
	flagFilePath  = "file_path"
	flagBatchSize = "batch_size"
)

type generateCmdFlags struct {
	FilePath  string `mapstructure:"file_path"`
	BatchSize uint32 `mapstructure:"batch_size"`
}

func defaultFlagsValues() generateCmdFlags {
	return generateCmdFlags{
		FilePath:  "./build/resources/csv/people.csv",
		BatchSize: 100000,
	}
}

// newGenerateCmd ...
func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
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

			file, err := os.Open(flagsValues.FilePath)
			if err != nil {
				return fmt.Errorf("opening file '%s': %w", flagsValues.FilePath, err)
			}
			defer file.Close()

			csvReader := csv.NewReader(file)

			password, err := model.User{Password: "password"}.EncryptPassword()
			if err != nil {
				return fmt.Errorf("encrypting password: %w", err)
			}

			objs := make([]model.User, 0, flagsValues.BatchSize)
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

				obj := model.User{
					FirstName:  name[1],
					SecondName: name[0],
					Age:        uint8(age),
					City:       record[2],
					Password:   password,
				}
				objs = append(objs, obj)

				if i >= flagsValues.BatchSize {
					if _, err = st.CopyUsers(ctx, objs); err != nil {
						return fmt.Errorf("copying users: %w", err)
					}

					fmt.Print(".")
					i = 0
					objs = make([]model.User, 0, flagsValues.BatchSize)
				}
			}
			if len(objs) > 0 {
				if _, err = st.CopyUsers(ctx, objs); err != nil {
					return fmt.Errorf("copying users: %w", err)
				}
			}

			return nil
		},
	}

	config := defaultFlagsValues()
	cmd.Flags().StringP(flagFilePath, "f", config.FilePath, "Path to file")
	cmd.Flags().Uint32P(flagBatchSize, "s", config.BatchSize, "Size of the batch")

	return cmd
}

func getFlags() (generateCmdFlags, error) {
	config := defaultFlagsValues()
	if err := viper.Unmarshal(&config); err != nil {
		return generateCmdFlags{}, fmt.Errorf("generateCmdFlags unmarshal: %w", err)
	}

	return config, nil
}
