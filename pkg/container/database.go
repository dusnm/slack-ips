package container

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func (c *Container) GetDB() *sql.DB {
	if c.db == nil {
		db, err := sql.Open("sqlite", c.GetConfig().DB.Path)
		if err != nil {
			c.logger.
				Fatal().
				Err(err).
				Msg("failure while establishing a connection to the database")
		}

		c.db = db
	}

	return c.db
}
