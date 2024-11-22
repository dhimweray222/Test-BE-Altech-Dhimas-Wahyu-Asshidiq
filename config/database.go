package config

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"test-backend-altech/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

var (
	host               = os.Getenv("DB_HOST")
	port               = os.Getenv("DB_PORT")
	username           = os.Getenv("DB_USERNAME")
	password           = os.Getenv("DB_PASSWORD")
	dbName             = os.Getenv("DB_NAME")
	minConns           = os.Getenv("DB_MIN_CONNS")
	maxConns           = os.Getenv("DB_MAX_CONNS")
	TimeOutDuration, _ = strconv.Atoi(os.Getenv("DB_CONNECTION_TIMEOUT"))
)

func NewPostgresDatabase() *pgxpool.Pool {
	logger := utils.NewLogger()

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Errorw("Failed to parse configuration", "dsn", dsn)
	}

	minConnsInt, err := strconv.Atoi(minConns)
	if err != nil {
		logger.Errorf("DB_MIN_CONNS expected to be integer", "minimum connections", minConns)
	}
	maxConnsInt, err := strconv.Atoi(maxConns)
	if err != nil {
		logger.Errorf("DB_MAX_CONNS expected to be integer", "maximum connections", maxConns)
	}

	poolConfig.MinConns = int32(minConnsInt)
	poolConfig.MaxConns = int32(maxConnsInt)
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Errorw("Failed to apply pool configuration", "dsn", dsn)
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := pool.Ping(c); err != nil {
		logger.Error(err)
	}

	logger.Infow("Database connected", "dsn", dsn)
	// Execute all SQL files in the "db" folder
	if err := executeSQLFiles(pool, "db", logger); err != nil {
		logger.Error("Failed to execute SQL files", "error", err)
		return nil
	}
	return pool
}

func executeSQLFiles(pool *pgxpool.Pool, folderPath string, logger *zap.SugaredLogger) error {
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".sql" {
			logger.Infow("Executing SQL file", "filePath", path)
			if err := executeSQLFile(pool, path, logger); err != nil {
				return fmt.Errorf("error executing SQL file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error processing SQL files in folder %s: %w", folderPath, err)
	}

	logger.Infow("All SQL files executed successfully", "folder", folderPath)
	return nil
}

func executeSQLFile(pool *pgxpool.Pool, filePath string, logger *zap.SugaredLogger) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Errorw("Failed to read SQL file", "filePath", filePath, "error", err)
		return err
	}

	_, err = pool.Exec(context.Background(), string(content))
	if err != nil {
		logger.Errorw("Failed to execute SQL script", "filePath", filePath, "error", err)
		return err
	}

	logger.Infow("SQL script executed successfully", "filePath", filePath)
	return nil
}
