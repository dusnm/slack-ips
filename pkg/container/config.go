package container

import (
	"os"
	"path/filepath"

	"github.com/dusnm/slack-ips/pkg/config"
)

func (c *Container) GetConfig() *config.Config {
	if c.cfg == nil {
		cfgDir, err := os.UserConfigDir()
		if err != nil {
			c.logger.
				Fatal().
				Err(err).
				Msg("could not get user config dir")
		}

		// A hierarchy of config paths
		paths := []string{
			"./config.toml",
			filepath.Join(cfgDir, "slack-ips", "config.toml"),
		}

		filename := ""
		for _, path := range paths {
			// The first file to be found will be used
			// for loading configuration
			_, err = os.Stat(path)
			if err == nil {
				filename = path
				break
			}
		}

		if filename == "" {
			c.logger.
				Fatal().
				Msg("could not find config.toml in any of the search paths")
		}

		cfgFile, err := os.OpenFile(filename, os.O_RDONLY, 0o644)
		if err != nil {
			c.logger.
				Fatal().
				Err(err).
				Msg("Failed to open config.toml")
		}

		defer cfgFile.Close()

		cfg, err := config.New(cfgFile)
		if err != nil {
			c.logger.
				Fatal().
				Err(err).
				Msg("Failed to parse config.toml")
		}

		c.cfg = cfg
	}

	return c.cfg
}
