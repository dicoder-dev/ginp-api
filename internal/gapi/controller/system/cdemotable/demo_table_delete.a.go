package cdemotable

import (
	"ginp-api/internal/gapi/dto/comdto"
	"ginp-api/internal/gapi/service/system/sdemotable"

	"ginp-api/pkg/ginp"
)

type RequestDemoTableDelete struct {
	ID uint `json:"id"`
}

type RespondDemoTableDelete struct {
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/demo_table/delete",
		Handler:        ginp.BindParamsHandler(DemoTableDelete, RequestDemoTableDelete{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      true,
		NeedPermission: true,
		PermissionName: "DemoTable.delete",
		Swagger: &ginp.SwaggerInfo{
			Title:         "delete demo_table",
			Description:   "",
			RequestParams: comdto.ReqDelete{},
		},
	})
}

func DemoTableDelete(c *ginp.ContextPlus, params *RequestDemoTableDelete) {
	err := sdemotable.Model().DeleteById(params.ID)
	if err != nil {
		c.FailData("delete fail :" + err.Error())
		return
	}
	c.Success()
}
