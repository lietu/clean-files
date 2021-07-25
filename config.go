package clean_files

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/kpango/glg"
	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
)

type Rule struct {
	Globs    []string
	Keep     time.Duration
}

type rawRule struct {
	Globs    []string `mapstructure:"globs"`
	KeepDays int      `mapstructure:"keep_days"`
	Keep     string   `mapstructure:"keep"`
}

type UserConfig struct {
	Rules map[string]Rule
}

func (r rawRule) validate(file, ruleName string) bool {
	if r.Keep != "" {
		if r.KeepDays != 0 {
			glg.Errorf("Error parsing %s rule %s: Can't define both keep and keep_days", file, ruleName)
			return false
		}
	}

	if r.KeepDays < 0 {
		glg.Errorf("Error parsing %s rule %s: keep_days must not be <0", file, ruleName)
		return false
	}

	if r.Keep != "" {
		_, err := time.ParseDuration(r.Keep)
		if err != nil {
			glg.Errorf("Error parsing %s rule %s: keep is not a valid duration: %s", file, ruleName, err)
			return false
		}
	}

	if r.Keep == "" && r.KeepDays == 0 {
		glg.Errorf("Error parsing %s rule %s: must define either keep or keep_days", file, ruleName)
		return false
	}

	return true
}

// Absolute defaults in case reading any config fails
func newUserConfig() UserConfig {
	return UserConfig{
		Rules: map[string]Rule{},
	}
}

func detectUserConfig() string {
	for _, path := range CurrentUserConfigs {
		// TODO: Parse env variables
		_, err := os.Stat(path)
		if err == nil {
			return path
		} else if !os.IsNotExist(err) {
			glg.Errorf("Error trying to stat() %s: %s", path, err)
		}
	}

	glg.Errorf("Failed to find default user config in any of the configured locations: %s", strings.Join(CurrentUserConfigs, ", "))
	return ""
}

func ReadUserConfig(configFile string) UserConfig {
	cfg := newUserConfig()

	if configFile == "" {
		configFile = detectUserConfig()
	}

	if configFile != "" {
		tomlContents, err := ioutil.ReadFile(configFile)
		if err != nil {
			glg.Errorf("Failed to read %s: %s", configFile, tomlContents)
			return cfg
		}

		tomlData := map[string]interface{}{}
		err = toml.Unmarshal(tomlContents, &tomlData)
		if err != nil {
			glg.Errorf("Failed to parse %s: %s", configFile, err)
			return cfg
		}

		globRules, ok := tomlData["globs"].(map[string]interface{})
		if ok {
			rules := map[string]Rule{}
			for name, rule := range globRules {
				r := rawRule{}
				err = mapstructure.Decode(rule, &r)
				if err != nil {
					glg.Errorf("Failed to decode rule %s in %s: %s", name, configFile, err)
				} else {
					if r.validate(configFile, "globs."+name) {
						keep := time.Duration(r.KeepDays) * time.Hour * 24
						if r.Keep != "" {
							keep, _ = time.ParseDuration(r.Keep)
						}

						rules[name] = Rule{
							Globs:    r.Globs,
							Keep:     keep,
						}
					}
				}
			}
			cfg.Rules = rules
		}
	} else {
		glg.Error("No configuration found.")
		return cfg
	}

	return cfg
}
