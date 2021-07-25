package main

import (
	clean_files "github.com/lietu/clean-files"
	"github.com/lietu/clean-files/cleaner"
	"github.com/lietu/clean-files/service"
	flag "github.com/spf13/pflag"
)

var (
	monitor *bool   = flag.BoolP("monitor", "m", false, "Run in service mode and monitor files over time")
	dryRun  *bool   = flag.BoolP("dry-run", "d", false, "Only detect files to clean up, do not actually delete anything.")
	config  *string = flag.StringP("config", "c", "", "Override configuration search paths")
)

func main() {
	flag.Parse()

	clean_files.ConfigureLogger()

	if *monitor {
		service.Run(*config, *dryRun)
	} else {
		cleaner.Run(*config, *dryRun)
	}
}
