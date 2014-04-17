package config

// The function getFilePath is to get real config path.
// It solves problem that if the configuration file had a password or other thing,
// and worried about other people see it in the file after sharing in the project.

// We can use like:
//     export ProjectName_config="a.config=/.../a.config;b=/.../b.config",
//     getFilePath("a.config") return "/.../a.config", and
// if export ProjectName_config="" or not, we use local configuration file: a.config.

// export test_config="test.config=$HOME/test.config"
// export test_config="test.config=$HOME/test.config;test2.config=$HOME/test2.config"
// export test_config=""

import (
	"fmt"
	"os"
	"strings"
)

const (
	pathHolder = "_config"
)

func getFilePath(fileName string) string {
	project_name := getProjectName()
	config_env := os.Getenv(project_name + pathHolder)

	if config_env == "" {
		return fileName
	} else {
		return parseENVPath(fileName, config_env)
	}
}

func getProjectName() string {
	project_path := os.Args[0]
	index := strings.LastIndex(project_path, "/")
	return project_path[index+1:]
}

func parseENVPath(fileName, envPath string) string {
	configPaths := strings.Split(envPath, ";")
	fmt.Println(configPaths)
	for _, v := range configPaths {
		kv := strings.Split(v, "=")
		if fileName == kv[0] {
			return kv[1]
		}
	}
	return fileName
}
