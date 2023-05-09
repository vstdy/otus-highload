package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy/otus-highload/api/metrics"
	"github.com/vstdy/otus-highload/api/rest"
	"github.com/vstdy/otus-highload/cmd/project/cmd/common"
	"github.com/vstdy/otus-highload/pkg/logging"
)

const (
	flagConfigPath    = "config"
	flagLogLevel      = "log_level"
	flagTimeout       = "timeout"
	flagServerAddress = "server_address"
	envSecretKey      = "secret_key"
	flagDatabaseDSN   = "database_dsn"
)

// Execute prepares cobra.Command context and executes root cmd.
func Execute() error {
	return newRootCmd().ExecuteContext(common.NewBaseCmdCtx())
}

// newRootCmd creates a new root command.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupConfig(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			logger := logging.NewLogger(logging.WithLogLevel(config.LogLevel))

			// Build servers
			svc, err := config.BuildService()
			if err != nil {
				return fmt.Errorf("app initialization: service building: %w", err)
			}

			srv, err := rest.NewServer(svc, config.HTTPServer)
			if err != nil {
				return fmt.Errorf("app initialization: server building: %w", err)
			}

			metricsSrv, err := metrics.NewMetricsServer(config.HTTPMetricsServer)
			if err != nil {
				return fmt.Errorf("app initialization: metrics server building: %w", err)
			}

			// Run servers
			go func() {
				if err = srv.ListenAndServe(); err != http.ErrServerClosed {
					logger.Error().Err(err).Msg("HTTP server ListenAndServe")
				}
			}()

			go func() {
				if err = metricsSrv.ListenAndServe(); err != http.ErrServerClosed {
					logger.Error().Err(err).Msg("metrics HTTP server ListenAndServe")
				}
			}()

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt)
			<-stop

			// Stop servers
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()
			if err = srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}
			if err = metricsSrv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("metrics server shutdown failed: %w", err)
			}

			if err = svc.Close(); err != nil {
				return fmt.Errorf("service shutdown failed: %w", err)
			}
			logger.Info().Msg("server stopped")

			return nil
		},
	}

	config := common.BuildDefaultConfig()
	cmd.PersistentFlags().String(flagConfigPath, "./config.yml", "Config file path")
	cmd.PersistentFlags().StringP(flagLogLevel, "l", config.LogLevel.String(), "Logger level [debug,info,warn,error,fatal]")
	cmd.PersistentFlags().Duration(flagTimeout, config.Timeout, "Request timeout")
	cmd.PersistentFlags().StringP(flagDatabaseDSN, "d", config.PSQLStorage.DSN, "Database source name")
	cmd.Flags().StringP(flagServerAddress, "a", config.HTTPServer.ServerAddress, "Server address")

	cmd.AddCommand(newMigrateCmd())
	cmd.AddCommand(newGenerateCmd())

	return cmd
}

// setupConfig reads app config and stores it to cobra.Command context.
func setupConfig(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("flags binding: %w", err)
	}

	if err := viper.BindEnv(envSecretKey); err != nil {
		return fmt.Errorf("%s env binding: %w", envSecretKey, err)
	}

	configPath := viper.GetString(flagConfigPath)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("reading config file: %w", err)
		}
	}

	viper.AutomaticEnv()
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	config := common.BuildDefaultConfig()
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("config unmarshal: %w", err)
	}

	logLevel, err := zerolog.ParseLevel(viper.GetString(flagLogLevel))
	if err != nil {
		return fmt.Errorf("%s flag parsing: %w", flagLogLevel, err)
	}
	config.LogLevel = logLevel
	config.HTTPServer.LogLevel = logLevel

	common.SetConfigToCmdCtx(cmd, config)

	return nil
}
