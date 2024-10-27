/**
 * This file is part of tz.
 *
 * tz is free software: you can redistribute it and/or modify it under
 * the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your
 * option) any later version.
 *
 * tz is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
 * or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public
 * License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tz.  If not, see <https://www.gnu.org/licenses/>.
 **/
package main

import (
	platform "os"
)

var os OsWrapper = PlatformOsWrapper{}

type OsWrapper interface {
	Exit(code int)
	Getenv(key string) string
	LookupEnv(key string) (string, bool)
	ReadFile(name string) ([]byte, error)
	Setenv(key, value string) error
	Stderr() *platform.File
	Stdin() *platform.File
	Stdout() *platform.File
	Unsetenv(key string) error
	UserHomeDir() (string, error)
	WriteFile(name string, data []byte, perm platform.FileMode) error
}

type PlatformOsWrapper struct {}

func (_ PlatformOsWrapper) Exit(code int) {
	platform.Exit(code)
}

func (_ PlatformOsWrapper) Getenv(key string) string {
	return platform.Getenv(key)
}

func (_ PlatformOsWrapper) LookupEnv(key string) (string, bool) {
	return platform.LookupEnv(key)
}

func (_ PlatformOsWrapper) ReadFile(name string) ([]byte, error) {
	return platform.ReadFile(name)
}

func (_ PlatformOsWrapper) Setenv(key, value string) error {
	return platform.Setenv(key, value)
}

func (_ PlatformOsWrapper) Stderr() *platform.File {
	return platform.Stderr
}

func (_ PlatformOsWrapper) Stdin() *platform.File {
	return platform.Stdin
}

func (_ PlatformOsWrapper) Stdout() *platform.File {
	return platform.Stdout
}

func (_ PlatformOsWrapper) Unsetenv(key string) error {
	return platform.Unsetenv(key)
}

func (_ PlatformOsWrapper) UserHomeDir() (string, error) {
	return platform.UserHomeDir()
}

func (_ PlatformOsWrapper) WriteFile(name string, data []byte, perm platform.FileMode) error {
	return platform.WriteFile(name, data, perm)
}
