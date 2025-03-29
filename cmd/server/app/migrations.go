package app

import (
	"go-yandex-practicum-metrics/internal/logger"
	"os/exec"
)

func MigrationsUp(dbDriver, dsn string) error {
	cmd := exec.Command("goose", "-dir", "migrations", dbDriver, dsn, "up")
	_, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
