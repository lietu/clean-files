package cleaner

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kpango/glg"
	clean_files "github.com/lietu/clean-files"
)

func Run(configFile string, dryRun bool) {
	config := clean_files.ReadUserConfig(configFile)
	CleanByRules(config.Rules, dryRun)
}

func CleanByRules(rules map[string]clean_files.Rule, dryRun bool) {
	for name, rule := range rules {
		ruleFiles := findAllFiles(rule.Globs)
		keep := time.Duration(rule.KeepDays) * time.Hour * 24
		if keep == 0 {
			keep = rule.Keep
		}

		cleanFiles(name, ruleFiles, keep, dryRun)
	}
}

func findAllFiles(globMatches []string) []string {
	files := []string{}

	for _, globMatch := range globMatches {
		matches, err := filepath.Glob(globMatch)
		if err != nil {
			glg.Errorf("Failed to do a match for %s: %s", globMatch, err)
			continue
		}

		files = append(files, matches...)
	}

	return files
}

func cleanFiles(name string, filePaths []string, keep time.Duration, dryRun bool) {
	before := time.Now().Add(-keep)
	for _, filePath := range filePaths {
		stat, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				// Already deleted, no problem
				continue
			}

			glg.Errorf("Error checking %s age: %s", filePath, err)
			continue
		}

		if stat.ModTime().Before(before) {
			since := time.Since(stat.ModTime())
			glg.Debugf("%s for rule %s is %s old", filePath, name, since)

			if !dryRun {
				glg.Printf("Deleting %s for rule %s", filePath, name)
				err := os.Remove(filePath)
				if err != nil {
					if os.IsNotExist(err) {
						// Already deleted, no problem
						continue
					}

					glg.Errorf("Error deleting %s: %s", filePath, err)
					continue
				}
			}
		}
	}
}
