package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/rs/zerolog"
)

var (
	schema = `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT,
			name TEXT,
			bank_account_number TEXT,
			city TEXT,
			ips_string TEXT,
			UNIQUE(username),
			UNIQUE(bank_account_number)
		);

		CREATE TABLE IF NOT EXISTS settings (
		   user_id TEXT PRIMARY KEY,
		   qr_fg_color VARCHAR(7),
		   qr_bg_color VARCHAR(7),
		   qr_shape VARCHAR(6) CHECK (qr_shape IN ('square', 'circle', 'liquid')), -- Poor man's enum
		   qr_logo BLOB,
		   qr_show_logo BOOLEAN NOT NULL DEFAULT 0 CHECK (qr_show_logo IN (0, 1)),
		   FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
    `
)

func InitDB(c *container.Container) {
	db := c.GetDB()
	logger := c.GetLogger().
		With().
		Str("component", "command:initdb").
		Logger()

	ensureDBPath(c, logger)

	_, err := db.Exec(schema)
	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failure while initializing the database")
	}

	logger.
		Debug().
		Msg("command completed successfully")
}

func ensureDBPath(c *container.Container, logger zerolog.Logger) {
	path := filepath.Dir(c.GetConfig().DB.Path)
	_, err := os.Stat(path)
	if err == nil {
		logger.
			Debug().
			Msg("database directory already exists")
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		logger.
			Fatal().
			Err(err).
			Msg("failure while checking if the database directory exists")
	}

	if err = os.MkdirAll(path, 0755); err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("failure while creating the database directory")
	}
}
