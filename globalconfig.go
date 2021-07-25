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

type GlobalConfig struct {
	CleanFrequency time.Duration
	Rules          map[string]Rule
}

type rawGlobalConfig struct {
	CleanFrequency string `mapstructure:"clean_frequency"`
}

// Absolute defaults in case reading any config fails
func newGlobalConfig() GlobalConfig {
	return GlobalConfig{
		CleanFrequency: time.Hour,
	}
}

func detectGlobalConfig() string {
	for _, path := range GlobalConfigs {
		if _, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
			glg.Errorf("Error trying to stat() %s: %s", path, err)
		}
	}

	glg.Errorf("Failed to find default config in any of the configured locations: %s", strings.Join(GlobalConfigs, ", "))
	return ""
}

func ReadGlobalConfig(configFile string) GlobalConfig {
	cfg := newGlobalConfig()

	if configFile == "" {
		configFile = detectGlobalConfig()
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

		_, ok := tomlData["main"]
		if ok {
			rgc := rawGlobalConfig{}
			err = mapstructure.Decode(tomlData["main"], &rgc)
			if err != nil {
				glg.Errorf("Failed to decode %s: %s", configFile, err)
			} else {
				cleanFrequency, err := time.ParseDuration(rgc.CleanFrequency)
				if err != nil {
					glg.Errorf("Failed to parse clean_frequency %s: %s", rgc.CleanFrequency, err)
				} else {
					cfg.CleanFrequency = cleanFrequency
				}
			}
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
