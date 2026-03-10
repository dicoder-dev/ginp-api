package cdemotable

import (
	"ginp-api/internal/gapi/entity"
	"ginp-api/internal/gapi/service/system/sdemotable"

	"ginp-api/pkg/ginp"
)

type RequestDemoTableCreate struct {
	entity.DemoTable
}

type RespondDemoTableCreate struct {
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/demo_table/create",
		Handler:        ginp.BindParamsHandler(DemoTableCreate, RequestDemoTableCreate{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      false,
		NeedPermission: false,
		PermissionName: "DemoTable.create",
		Swagger: &ginp.SwaggerInfo{
			Title:         "create demo_table",
			Description:   "",
			RequestParams: RequestDemoTableCreate{},
		},
	})
}

func DemoTableCreate(c *ginp.ContextPlus, params *RequestDemoTableCreate) {
	info, err := sdemotable.Model().Create(&params.DemoTable)
	if err != nil {
		c.FailData(err.Error())
		return
	}
	c.SuccessData(info)
}
