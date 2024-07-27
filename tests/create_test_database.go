package tests

import (
	"ChaikaGoods/internal/config"
	"context"
	"fmt"
	"os/exec"
)

// createDatabase creates a new database with the given name for testing purposes.
func createDatabase(ctx context.Context, storageCfg config.StorageConfig) error {
	cmdStr := fmt.Sprintf("createdb -h %s -p %s -U %s %s",
		storageCfg.Host, storageCfg.Port, storageCfg.User, storageCfg.Database)

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", storageCfg.Password))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create database: %v, output: %s", err, string(output))
	}

	return nil
}

// restoreDatabase restores a database from a dump file with test data.
func restoreDatabase(ctx context.Context, storageCfg config.StorageConfig, dumpFile string) error {
	cmdStr := fmt.Sprintf("psql -h %s -p %s -U %s -d %s -f %s",
		storageCfg.Host, storageCfg.Port, storageCfg.User, storageCfg.Database, dumpFile)

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", storageCfg.Password))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restore database: %v, output: %s", err, string(output))
	}

	return nil
}

// deleteDatabase deletes a database after tests.
func deleteDatabase(ctx context.Context, storageCfg config.StorageConfig) error {
	cmdStr := fmt.Sprintf("dropdb -h %s -p %s -U %s %s",
		storageCfg.Host, storageCfg.Port, storageCfg.User, storageCfg.Database)

	cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", storageCfg.Password))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete database: %v, output: %s", err, string(output))
	}

	return nil
}
