package utils

import (
	"github.com/astaxie/beego/pkg/infrastructure/utils"
)

// GetGOPATHs returns all paths in GOPATH variable.
func GetGOPATHs() []string {
	return utils.GetGOPATHs()
}
