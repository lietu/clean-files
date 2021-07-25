package clean_files

import "github.com/kpango/glg"

func ConfigureLogger() {
	glg.Get().EnableTimestamp()
	glg.Get().DisableJSON()
}
