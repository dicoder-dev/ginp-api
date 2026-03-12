package entity

import (
	"ginp-api/internal/gapi/typ"
	"ginp-api/internal/gen"
	"time"
)

const tableNameTest = "test"

type Test struct {
	ID uint `json:"id"`
	//... other fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" `
}

var _ typ.IEntity = (*Test)(nil) // U实体必须实现接口GenConfig

func (Test) GenConfig() *gen.EntityConfig {
	return &gen.EntityConfig{
		TableName: tableNameTest,
	}
}

func (Test) TableName() string {
	return tableNameTest
}
