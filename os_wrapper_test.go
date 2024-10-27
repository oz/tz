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
	"errors"
	platform "os"
	"testing"
)

type TestingOsWrapper struct {
	EnvVars            map[string]string
	ExitCode           *int
	Files              map[string][]byte
	HomeDir            string
	PlatformStdinFile  *platform.File
	PlatformStderrFile *platform.File
	PlatformStdoutFile *platform.File
	TestStdinFile      *platform.File
	TestStderrFile     *platform.File
	TestStdoutFile     *platform.File
}

func NewTestingOsWrapper(t *testing.T) *TestingOsWrapper {
	savedOs := os
	savedOsArgs := platform.Args
	t.Cleanup(func () {
		os = savedOs
		platform.Args = savedOsArgs
	})

	stderr, err1 := platform.CreateTemp("", "stderr")
	if err1 != nil {
		t.Fatalf("Could not create stderr temp file for %v: %v", t.Name(), err1)
	}
	t.Cleanup(func () {
		stderr.Sync()
		file := stderr.Name()
		bytes, _ := platform.ReadFile(file)
		platform.Stderr.Write(bytes)
		platform.Remove(file)
	})

	stdout, err2 := platform.CreateTemp("", "stdout")
	if err2 != nil {
		t.Fatalf("Could not create stdout temp file for %v: %v", t.Name(), err2)
	}
	t.Cleanup(func () {
		stdout.Sync()
		file := stdout.Name()
		bytes, _ := platform.ReadFile(file)
		platform.Stdout.Write(bytes)
		platform.Remove(file)
	})

	osWrapper := &TestingOsWrapper{
		EnvVars:            make(map[string]string),
		Files:              make(map[string][]byte),
		TestStderrFile:     stderr,
		TestStdinFile:      platform.Stdin,
		TestStdoutFile:     stdout,
		PlatformStderrFile: platform.Stderr,
		PlatformStdinFile:  platform.Stdin,
		PlatformStdoutFile: platform.Stdout,
	}
	os = osWrapper
	return osWrapper
}

func (t *TestingOsWrapper) ConsumeStderr() string {
	t.TestStderrFile.Sync()
	file := t.TestStderrFile.Name()
	bytes, _ := platform.ReadFile(file)
	platform.Truncate(file, 0)
	return string(bytes)
}

func (t *TestingOsWrapper) ConsumeStdout() string {
	t.TestStdoutFile.Sync()
	file := t.TestStdoutFile.Name()
	bytes, _ := platform.ReadFile(file)
	platform.Truncate(file, 0)
	return string(bytes)
}

func (t *TestingOsWrapper) Exit(code int) {
	t.ExitCode = &code
}

func (t *TestingOsWrapper) Getenv(key string) string {
	return t.EnvVars[key]
}

func (t *TestingOsWrapper) LookupEnv(key string) (string, bool) {
	value, exists := t.EnvVars[key]
	return value, exists
}

func (t *TestingOsWrapper) ReadFile(name string) ([]byte, error) {
	if bytes, exists := t.Files[name]; exists {
		return bytes, nil
	} else {
		return platform.ReadFile(name)
	}
}

func (t *TestingOsWrapper) RedirectPlatformOutput() {
	platform.Stderr = t.TestStderrFile
	platform.Stdout = t.TestStdoutFile
}

func (t *TestingOsWrapper) RevertPlatformOutput() {
	platform.Stderr = t.PlatformStderrFile
	platform.Stdout = t.PlatformStdoutFile
}

func (t *TestingOsWrapper) Setargs(args []string) {
	platform.Args = append([]string{platform.Args[0]}, args...)
}

func (t *TestingOsWrapper) Setenv(key, value string) error {
	t.EnvVars[key] = value
	return nil
}

func (t *TestingOsWrapper) Stderr() *platform.File {
	return t.TestStderrFile
}

func (t *TestingOsWrapper) Stdin() *platform.File {
	return t.TestStdinFile
}

func (t *TestingOsWrapper) Stdout() *platform.File {
	return t.TestStdoutFile
}

func (t *TestingOsWrapper) Unsetenv(key string) error {
	delete(t.EnvVars, key)
	return nil
}

func (t *TestingOsWrapper) UserHomeDir() (string, error) {
	if len(t.HomeDir) < 1 {
		return t.HomeDir, errors.New("TestingOsWrapper.HomeDir is not yet set")
	} else {
		return t.HomeDir, nil
	}
}

func (t *TestingOsWrapper) WriteFile(name string, data []byte, perm platform.FileMode) error {
	if _, exists := t.Files[name]; exists {
		t.Files[name] = data
		return nil
	} else {
		return platform.WriteFile(name, data, perm)
	}
}
