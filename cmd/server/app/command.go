package app

import (
	"context"
	"fmt"
	"go-yandex-practicum-metrics/internal/configs"
	"go-yandex-practicum-metrics/internal/logger"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DefaultAddress         = ":8080"
	DefaultDatabaseDSN     = "postgres://user:password@localhost:5432/db"
	DefaultStoreInterval   = "300" // Значение по умолчанию - 300 секунд
	DefaultFileStoragePath = ""
	DefaultRestore         = "false"

	EnvAddress         = "ADDRESS"
	EnvDatabaseDSN     = "DATABASE_DSN"
	EnvStoreInterval   = "STORE_INTERVAL"
	EnvFileStoragePath = "FILE_STORAGE_PATH"
	EnvRestore         = "RESTORE"

	FlagServerAddress   = "address"
	FlagDatabaseDSN     = "database-dsn"
	FlagStoreInterval   = "store-interval"
	FlagFileStoragePath = "file-storage-path"
	FlagRestore         = "restore"

	FlagShortServerAddress   = "a"
	FlagShortDatabaseDSN     = "d"
	FlagShortStoreInterval   = "i"
	FlagShortFileStoragePath = "f"
	FlagShortRestore         = "r"

	DescriptionServerAddress   = "Server address"
	DescriptionDatabaseDSN     = "Database DSN"
	DescriptionStoreInterval   = "Interval in seconds for data store (0 - synchronous write)"
	DescriptionFileStoragePath = "Path to file storage"
	DescriptionRestore         = "Restore backup (true/false)"
)

func NewCommand() *cobra.Command {
	var config *configs.ServerConfig

	viper.AutomaticEnv()   // Автоматическая привязка переменных окружения
	viper.SetEnvPrefix("") // Убираем префикс для переменных окружения, если нужно

	var cmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			bindFlagsAndEnv(cmd)
			config = configs.NewServerConfig(
				viper.GetString(FlagServerAddress),
				viper.GetString(FlagDatabaseDSN),
				viper.GetString(FlagStoreInterval),
				viper.GetString(FlagFileStoragePath),
				viper.GetString(FlagRestore),
			)
			if err := validateConfig(config, cmd); err != nil {
				cmd.PrintErrf("%v\n", err)
				return nil
			}
			logger.Init(logger.DEBUG)
			container, err := NewContainer(config)
			if err != nil {
				cmd.PrintErrf("failed to create container: %v\n", err)
				return nil
			}
			server, err := NewServer(config, container)
			if err != nil {
				cmd.PrintErrf("failed to create server: %v\n", err)
				return nil
			}
			worker, err := NewWorker(config, container)
			if err != nil {
				cmd.PrintErrf("failed to create worker: %v\n", err)
				return nil
			}
			if config.GetDatabaseDSN() != "" {
				MigrationsUp("postgres", config.GetDatabaseDSN())
			}
			if err := server.Start(ctx); err != nil {
				cmd.PrintErrf("failed to start server: %v\n", err)
				return nil
			}
			if err := worker.Start(ctx); err != nil {
				cmd.PrintErrf("failed to start worker: %v\n", err)
				return nil
			}
			<-ctx.Done()
			return nil
		},
	}
	cmd.Flags().StringP(FlagServerAddress, FlagShortServerAddress, DefaultAddress, DescriptionServerAddress)
	cmd.Flags().StringP(FlagDatabaseDSN, FlagShortDatabaseDSN, DefaultDatabaseDSN, DescriptionDatabaseDSN)
	cmd.Flags().StringP(FlagFileStoragePath, FlagShortFileStoragePath, DefaultFileStoragePath, DescriptionFileStoragePath)
	cmd.Flags().StringP(FlagStoreInterval, FlagShortStoreInterval, DefaultStoreInterval, DescriptionStoreInterval)
	cmd.Flags().StringP(FlagRestore, FlagShortRestore, DefaultRestore, DescriptionRestore)
	return cmd
}

func bindFlagsAndEnv(cmd *cobra.Command) {
	viper.BindPFlag(FlagServerAddress, cmd.Flags().Lookup(FlagServerAddress))
	viper.BindPFlag(FlagDatabaseDSN, cmd.Flags().Lookup(FlagDatabaseDSN))
	viper.BindPFlag(FlagFileStoragePath, cmd.Flags().Lookup(FlagFileStoragePath))
	viper.BindPFlag(FlagStoreInterval, cmd.Flags().Lookup(FlagStoreInterval))
	viper.BindPFlag(FlagRestore, cmd.Flags().Lookup(FlagRestore))
	viper.BindEnv(FlagServerAddress, EnvAddress)
	viper.BindEnv(FlagDatabaseDSN, EnvDatabaseDSN)
	viper.BindEnv(FlagStoreInterval, EnvStoreInterval)
	viper.BindEnv(FlagFileStoragePath, EnvFileStoragePath)
	viper.BindEnv(FlagRestore, EnvRestore)
}

func validateConfig(config *configs.ServerConfig, cmd *cobra.Command) error {
	storeInterval, err := strconv.Atoi(config.GetStoreInterval())
	if err != nil || storeInterval < 0 {
		return fmt.Errorf("invalid STORE_INTERVAL: must be a non-negative integer")
	}
	if config.GetRestore() != "true" && config.GetRestore() != "false" {
		return fmt.Errorf("invalid RESTORE value: must be a boolean (true/false)")
	}
	return nil
}

func Run(cmd *cobra.Command) int {
	if err := cmd.Execute(); err != nil {
		return 1
	}
	return 0
}
