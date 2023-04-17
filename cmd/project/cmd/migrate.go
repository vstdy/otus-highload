package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/vstdy/otus-highload/cmd/project/cmd/common"
	"github.com/vstdy/otus-highload/pkg/logging"
)

// newMigrateCmd creates a new migration command.
func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migration operations",
	}

	cmd.AddCommand(migrateUp())
	cmd.AddCommand(migrateDown())

	return cmd
}

// newMigrateCmd ...
func migrateUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate DB to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			_, logger := logging.GetCtxLogger(context.Background(), logging.WithLogLevel(config.LogLevel))

			st, err := config.BuildStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					logger.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			if err = st.MigrateUp(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

// newMigrateCmd ...
func migrateDown() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Roll back a single migration from the current version",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			_, logger := logging.GetCtxLogger(context.Background(), logging.WithLogLevel(config.LogLevel))

			st, err := config.BuildStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					logger.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			if err = st.MigrateDown(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
