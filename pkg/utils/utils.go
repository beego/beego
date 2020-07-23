package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" && compareGoVersion(runtime.Version(), "go1.8") >= 0 {
		gopath = defaultGOPATH()
	}
	return filepath.SplitList(gopath)
}

func compareGoVersion(a, b string) int {
	reg := regexp.MustCompile("^\\d*")

	a = strings.TrimPrefix(a, "go")
	b = strings.TrimPrefix(b, "go")

	versionsA := strings.Split(a, ".")
	versionsB := strings.Split(b, ".")

	for i := 0; i < len(versionsA) && i < len(versionsB); i++ {
		versionA := versionsA[i]
		versionB := versionsB[i]

		vA, err := strconv.Atoi(versionA)
		if err != nil {
			str := reg.FindString(versionA)
			if str != "" {
				vA, _ = strconv.Atoi(str)
			} else {
				vA = -1
			}
		}

		vB, err := strconv.Atoi(versionB)
		if err != nil {
			str := reg.FindString(versionB)
			if str != "" {
				vB, _ = strconv.Atoi(str)
			} else {
				vB = -1
			}
		}

		if vA > vB {
			// vA = 12, vB = 8
			return 1
		} else if vA < vB {
			// vA = 6, vB = 8
			return -1
		} else if vA == -1 {
			// vA = rc1, vB = rc3
			return strings.Compare(versionA, versionB)
		}

		// vA = vB = 8
		continue
	}

	if len(versionsA) > len(versionsB) {
		return 1
	} else if len(versionsA) == len(versionsB) {
		return 0
	}

	return -1
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		return filepath.Join(home, "go")
	}
	return ""
}
