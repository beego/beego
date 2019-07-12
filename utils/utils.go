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

// IsModuleEnabled checks whether go module is enabled according to
// https://github.com/golang/go/wiki/Modules#when-do-i-get-old-behavior-vs-new-module-based-behavior
func IsModuleEnabled(gopaths []string, cwd string) bool {
	gm := os.Getenv("GO111MODULE")
	switch gm {
	case "on":
		return true
	case "off":
		return false
	default:
		for _, gopath := range gopaths {
			if gopath == "" {
				continue
			}
			if InDir(cwd, filepath.Join(gopath, "src")) != "" {
				return false
			}
		}
		return true
	}
}

// InDir checks whether path is in the file tree rooted at dir.
// If so, InDir returns an equivalent path relative to dir.
// If not, InDir returns an empty string.
// InDir makes some effort to succeed even in the presence of symbolic links.
func InDir(path, dir string) string {
	if rel := inDirLex(path, dir); rel != "" {
		return rel
	}
	xpath, err := filepath.EvalSymlinks(path)
	if err != nil || xpath == path {
		xpath = ""
	} else {
		if rel := inDirLex(xpath, dir); rel != "" {
			return rel
		}
	}

	xdir, err := filepath.EvalSymlinks(dir)
	if err == nil && xdir != dir {
		if rel := inDirLex(path, xdir); rel != "" {
			return rel
		}
		if xpath != "" {
			if rel := inDirLex(xpath, xdir); rel != "" {
				return rel
			}
		}
	}
	return ""
}

// inDirLex is like inDir but only checks the lexical form of the file names.
// It does not consider symbolic links.
// return the suffix. Most uses of str.HasFilePathPrefix should probably
// be calling InDir instead.
func inDirLex(path, dir string) string {
	pv := strings.ToUpper(filepath.VolumeName(path))
	dv := strings.ToUpper(filepath.VolumeName(dir))
	path = path[len(pv):]
	dir = dir[len(dv):]
	switch {
	default:
		return ""
	case pv != dv:
		return ""
	case len(path) == len(dir):
		if path == dir {
			return "."
		}
		return ""
	case dir == "":
		return path
	case len(path) > len(dir):
		if dir[len(dir)-1] == filepath.Separator {
			if path[:len(dir)] == dir {
				return path[len(dir):]
			}
			return ""
		}
		if path[len(dir)] == filepath.Separator && path[:len(dir)] == dir {
			if len(path) == len(dir)+1 {
				return "."
			}
			return path[len(dir)+1:]
		}
		return ""
	}
}
