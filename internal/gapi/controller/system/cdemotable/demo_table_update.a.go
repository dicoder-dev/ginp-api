package cdemotable

import (
	"ginp-api/internal/gapi/entity"
	"ginp-api/internal/gapi/service/system/sdemotable"
	"ginp-api/pkg/where"

	"ginp-api/pkg/ginp"
)

type RequestDemoTableUpdate struct {
	entity.DemoTable
}

type RespondDemoTableUpdate struct {
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/demo_table/update",
		Handler:        ginp.BindParamsHandler(DemoTableUpdate, RequestDemoTableUpdate{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      true,
		NeedPermission: true,
		PermissionName: "DemoTable.update",
		Swagger: &ginp.SwaggerInfo{
			Title:         "modify demo_table",
			Description:   "",
			RequestParams: RequestDemoTableUpdate{},
		},
	})
}

func DemoTableUpdate(c *ginp.ContextPlus, params *RequestDemoTableUpdate) {
	wheres := where.Format(where.OptEqual("id", params.DemoTable.ID))
	err := sdemotable.Model().Update(wheres, &params.DemoTable)
	if err != nil {
		c.FailData("修改失败" + err.Error())
		return
	}
	c.Success()
}
