package qb

import "github.com/beego/beego/v2/client/orm/internal/models"

type core struct {
	dialect  Dialect
	registry *models.ModelCache
}
