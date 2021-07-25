package service

import (
	"time"

	clean_files "github.com/lietu/clean-files"
	"github.com/lietu/clean-files/cleaner"
)

func Run(configFile string, dryRun bool) {
	for {
		// TODO: Maybe use inotify or something to monitor config for changes
		// Reload configuration for every iteration
		config := clean_files.ReadGlobalConfig(configFile)
		cleaner.CleanByRules(config.Rules, dryRun)

		// TODO: Find all users and loop through them

		// Wait until next poll
		time.Sleep(config.CleanFrequency)
	}
}
