# File cleanup utility

Utility and system service to clean up unused files.

## Installation

TODO: Make it easy to install this on Win, Lin, Mac?

## Configuration

Main configuration file in `/etc/cleanup-files.toml`, `/usr/local/etc/cleanup-files.toml`, `C:\ProgramData\cleanup-files\cleanup-files.toml` or `C:\Program Files\cleanup-files\cleanup-files.toml`, or defined via `-c path/to/cleanup-files.toml`.

```toml
[main]
clean_frequency = "60s" # supports ms, s, m, h
```

Main configuration and user configuration can also define cleanup rules. User configuration files are looked for in `/home/*/.cleanup-files.toml`, `/Users/*/.cleanup-files.toml` and, `C:\Users\*\.cleanup-files.toml` when in service mode.

```toml
[globs.downloads]
globs = ["/home/*/Downloads/**/*", "C:\\Users\\*\\Downloads\\**\\*"]
keep_days = 14  # Each day is exactly 24 hours, value of 0 is ignored
# or
keep = "1h" # supports ms, s, m, h
```

## Running

One-shot, only runs once, with only config for current user being parsed.

```
clean-files
# or
clean-files -c path/to/cleanup-files.toml
```

As a service, stays in the background and periodically checks for updates.

```
clean-files --monitor
```

## Development

Ensure you have [pre-commit](https://pre-commit.com/#install) installed and you run `pre-commit install` in the repo root.


# License

Short answer: This software is licensed with the BSD 3-clause -license.

Long answer: The license for this software is in [LICENSE.md](./LICENSE.md), the libraries used may have varying other licenses that you need to be separately aware of.


# Financial support

This project has been made possible thanks to [Cocreators](https://cocreators.ee) and [Lietu](https://lietu.net). You can help us continue our open source work by supporting us on [Buy me a coffee](https://www.buymeacoffee.com/cocreators).

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/cocreators)
