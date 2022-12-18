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
semi-colon separated list of [tz data][tzdata] zone names (see
*Configuration* below). Command-line arguments trump the environment
variable.

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
You should specify what time zones to display by setting the `TZ_LIST`
environment variable. The local time is always displayed first. To
display your local time, the time in California, and the time Paris, you
have to set `TZ_LIST` to `US/Pacific;Europe/Paris`

## Zone Alias

The `TZ_LIST` env. variable recognizes items from the standard tz
database names, but you can alias these, using a special value: use the
standard name followed by `,` and your alias. For example:

```
TZ_LIST="Europe/Paris,EMEA office;US/Central,US office"
```


# Building

You need a recent-ish release of go with modules support:

```
git clone https://github.com/oz/tz
cd tz
go build
```


# Testing

```
go test -cover
```


# Contributing

Please do file bugs, and feature requests.  I am accepting patches too,
those are the best, but please, open an issue first to discuss your
changes. 😄


# License

The GPL3 license.

Copyright (c) 2021-2022 Arnaud Berthomier
