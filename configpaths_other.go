// +build !windows

package clean_files

var CurrentUserConfigs = []string{
	"$HOME/.cleanup-files.toml",
	"/home/$USER/.cleanup-files.toml",
	"/Users/$USER/.cleanup-files.toml",
}

var GlobalConfigs = []string{
	"/etc/cleanup-files.toml",
	"/usr/local/etc/cleanup-files.toml",
}
