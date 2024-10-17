# 🌐 A time zone helper

tz helps you schedule things across time zones. It's an interactive TUI
program that displays time across the time zones of your choosing.


# Usage

Run `tz` with no arguments to show the local time, as well as the UTC
time zone. It gets more useful when you pass some time zones to the
program, to list those below the local time zone.

For now, you need to select the time zones from the [tz_data][tzdata]
list. Yes, there are plans to make this friendlier for humans too.
You're welcome to file an issue about it. I enjoy reading those.

If you would rather not type the list everytime, you could set an alias
for your shell, or use the `TZ_LIST` environment variable with a
semi-colon separated list of [tz data][tzdata] zone names, or use a
configuration file (see *Configuration* below).

Check out `tz -h` for other flags.

<p align="center">
<img align="center" src="./docs/tz.png" />
</p>

The program will adjust to light and dark terminals themes.

[tzdata]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones


# Installing

I provide linux/amd64 builds for "official" releases, and you can build
from source for your favorite architecture. Kind souls have also
packaged the program for other OSes.

## Packages

### Brew

Brew has a tz package: `brew install tz`

### Archlinux

If you're an Archlinux user, packages are also available:

  - [tz][tz-arch] follows releases and,
  - [tz-git][tz-arch-git] builds the `main` git branch.

[tz-arch]: https://aur.archlinux.org/packages/tz
[tz-arch-git]: https://aur.archlinux.org/packages/tz-git

### Go
```
go install github.com/oz/tz@latest
```

# Configuration

The tz program uses standard time zones as described
[here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).

You can specify preferences through a configuration file, an environment variable,
or as command arguments.

They are applied in this order, overriding at each subsequent step:

1. Configuration file `~/.config/tz/conf.toml`
2. Environment variable `TZ_LIST`
3. Command line arguments `tz UTC`

The local time zone is always displayed first.

## Configuration File

Configs are read from `$HOME/.config/tz/conf.toml`.

Time zone `id`s should reference items from the standard tz database names.
Alias them by providing your own name with the `name` key.

Sample configuration: [example-conf.toml](./example-conf.toml)

## Environment Variable

This method only supports setting time zones. Keymaps must be configured through
the configuration file.

Specify time zones to display by setting the `TZ_LIST` environment variable. For
example, to display your local time, the time in California, and the time Paris,
set `TZ_LIST` to `US/Pacific;Europe/Paris`.

The `TZ_LIST` environment variable recognizes items from the standard tz
database names, but you can alias these, using a special value: use the
standard name followed by `,` and your alias. For example:

```bash
TZ_LIST="Europe/Paris,EMEA office;US/Central,US office"
```

If adding this to a shell configuration, remember to export it:

```bash
export TZ_LIST="Europe/Paris,EMEA office;US/Central,US office"
```


# Building

You need a recent-ish release of go with modules support:

```bash
git clone https://github.com/oz/tz
cd tz
go build
```


# Testing

```bash
go test -cover
```


# Debug

```bash
DEBUG=1 tz # Logs will write to debug.log
```


# Contributing

Please do file bugs, and feature requests.  I am accepting patches too,
those are the best, but please, open an issue first to discuss your
changes. 😄


# License

The GPL3 license.

Copyright (c) 2021-2022 Arnaud Berthomier
