// Copyright 2014 beego Author. All Rights Reserved.
// Copyright 2017 Faissal Elamraoui. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package env is used to parse environment.
package env

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var env sync.Map

func init() {
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		env.Store(splits[0], os.Getenv(splits[0]))
	}
	env.Store("GOBIN", GetGOBIN())   // set non-empty GOBIN when initialization
	env.Store("GOPATH", GetGOPATH()) // set non-empty GOPATH when initialization
}

// Get returns a value for a given key.
// If the key does not exist, the default value will be returned.
func Get(key string, defVal string) string {
	if val, ok := env.Load(key); ok {
		return val.(string)
	}
	return defVal
}

// MustGet returns a value by key.
// If the key does not exist, it will return an error.
func MustGet(key string) (string, error) {
	if val, ok := env.Load(key); ok {
		return val.(string), nil
	}
	return "", fmt.Errorf("no env variable with %s", key)
}

// Set sets a value in the ENV copy.
// This does not affect the child process environment.
func Set(key string, value string) {
	env.Store(key, value)
}

// MustSet sets a value in the ENV copy and the child process environment.
// It returns an error in case the set operation failed.
func MustSet(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		return err
	}
	env.Store(key, value)
	return nil
}

// GetAll returns all keys/values in the current child process environment.
func GetAll() map[string]string {
	envs := make(map[string]string, 32)
	env.Range(func(key, value interface{}) bool {
		switch key := key.(type) {
		case string:
			switch val := value.(type) {
			case string:
				envs[key] = val
			}
		}
		return true
	})
	return envs
}

// envFile returns the name of the Go environment configuration file.
// Copy from https://github.com/golang/go/blob/c4f2a9788a7be04daf931ac54382fbe2cb754938/src/cmd/go/internal/cfg/cfg.go#L150-L166
func envFile() (string, error) {
	if file := os.Getenv("GOENV"); file != "" {
		if file == "off" {
			return "", fmt.Errorf("GOENV=off")
		}
		return file, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("missing user-config dir")
	}
	return filepath.Join(dir, "go", "env"), nil
}

// GetRuntimeEnv returns the value of runtime environment variable,
// that is set by running following command: `go env -w key=value`.
func GetRuntimeEnv(key string) (string, error) {
	file, err := envFile()
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("missing runtime env file")
	}

	var runtimeEnv string
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	envStrings := strings.Split(string(data), "\n")
	for _, envItem := range envStrings {
		envItem = strings.TrimSuffix(envItem, "\r")
		envKeyValue := strings.Split(envItem, "=")
		if len(envKeyValue) == 2 && strings.TrimSpace(envKeyValue[0]) == key {
			runtimeEnv = strings.TrimSpace(envKeyValue[1])
		}
	}
	return runtimeEnv, nil
}

// GetGOBIN returns GOBIN environment variable as a string.
// It will NOT be an empty string.
func GetGOBIN() string {
	// The one set by user explicitly by `export GOBIN=/path` or `env GOBIN=/path command`
	gobin := strings.TrimSpace(Get("GOBIN", ""))
	if gobin == "" {
		var err error
		// The one set by user by running `go env -w GOBIN=/path`
		gobin, err = GetRuntimeEnv("GOBIN")
		if err != nil {
			// The default one that Golang uses
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		if gobin == "" {
			return filepath.Join(build.Default.GOPATH, "bin")
		}
		return gobin
	}
	return gobin
}

// GetGOPATH returns GOPATH environment variable as a string.
// It will NOT be an empty string.
func GetGOPATH() string {
	// The one set by user explicitly by `export GOPATH=/path` or `env GOPATH=/path command`
	gopath := strings.TrimSpace(Get("GOPATH", ""))
	if gopath == "" {
		var err error
		// The one set by user by running `go env -w GOPATH=/path`
		gopath, err = GetRuntimeEnv("GOPATH")
		if err != nil {
			// The default one that Golang uses
			return build.Default.GOPATH
		}
		if gopath == "" {
			return build.Default.GOPATH
		}
		return gopath
	}
	return gopath
}
