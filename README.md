# 🌐 A time zone helper

tz helps you schedule things across time zones. It is an interactive TUI
program that displays time across a few time zones of your choosing.


# Usage

Simply run `tz` with no arguments to show the local time, as well as the
UTC time zone. It gets more interesting once you set the `TZ_LIST`
environment variable with a comma-separated list of [tz data][tzdata]
zone names.

<p align="center">
<img align="center" src="./docs/tz.png" />
</p>

The program will adjust to light and dark terminals themes.

[tzdata]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones


# Installing

The simplest thing is probably to grab a release, but no one will be
harmed if you build from source, as only linux/amd64 builds are provided
for now.

## Packages

If you're an Archlinux user, packages are also available:

  - [tz][tz-arch] follows releases and,
  - [tz-git][tz-arch-git] builds the `main` git branch.

[tz-arch]: https://aur.archlinux.org/packages/tz
[tz-arch-git]: https://aur.archlinux.org/packages/tz-git


# Configuration

## Zone Alias

tz is configured only through `TZ_LIST`, and that limits us to the tz
database names, but you can alias these names using a special value: the
tz name followed by `;` and your alias:

`TZ_LIST="Europe/Paris;EMEA office,US/Central;US office"`

## Name of local zone

tz, by default, displays your local timezone as "Local". If you rather like
to display a different name there instead, you can use the `TZ_LOCAL_NAME`
environment variable to override it.

`TZ_LOCAL_NAME="Cologe office`

## Display time format

In the header line, tz will display the current time. By default the time format
for the time is in 12H format. If you prefer to display the time in 24H format, 
you can do so by setting `TZ_24H`

`TZ_24H=1`

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


# License

The GPL3 license.

Copyright (c) 2021 Arnaud Berthomier
