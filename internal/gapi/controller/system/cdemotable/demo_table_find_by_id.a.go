package cdemotable

import (
	"ginp-api/internal/gapi/dto/comdto"
	"ginp-api/internal/gapi/service/system/sdemotable"

	"ginp-api/pkg/ginp"
)

type RequestDemoTableFindById struct {
	ID uint `json:"id"`
}

type RespondDemoTableFindById struct {
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/demo_table/findById",
		Handler:        ginp.BindParamsHandler(DemoTableFindById, RequestDemoTableFindById{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      false,
		NeedPermission: false,
		PermissionName: "DemoTable.findById",
		Swagger: &ginp.SwaggerInfo{
			Title:         "find demo_table by id",
			Description:   "",
			RequestParams: comdto.ReqFindById{},
		},
	})
}

func DemoTableFindById(c *ginp.ContextPlus, requestParams *RequestDemoTableFindById) {
	info, err := sdemotable.Model().FindOneById(requestParams.ID)
	if err != nil {
		c.FailData(err.Error())
		return
	}
	c.SuccessData(info)
}
