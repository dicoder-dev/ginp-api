package cdemotable

import (
	"ginp-api/internal/gapi/dto/comdto"
	"ginp-api/internal/gapi/service/system/sdemotable"

	"ginp-api/pkg/ginp"
)

type RequestDemoTableSearch struct {
	comdto.ReqSearch
}

type RespondDemoTableSearch struct {
	List     interface{} `json:"list"`
	Total    uint        `json:"total"`
	PageNum  uint        `json:"page_num"`
	PageSize uint        `json:"page_size"`
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/demo_table/search",
		Handler:        ginp.BindParamsHandler(DemoTableSearch, RequestDemoTableSearch{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      true,
		NeedPermission: true,
		PermissionName: "DemoTable.search",
		Swagger: &ginp.SwaggerInfo{
			Title:         "search demo_table",
			Description:   "",
			RequestParams: comdto.ReqSearch{},
		},
	})
}

func DemoTableSearch(c *ginp.ContextPlus, params *RequestDemoTableSearch) {
	list, total, err := sdemotable.Model().FindList(params.Wheres, params.Extra)
	if err != nil {
		c.FailData(err.Error())
		return
	}

	resp := &RespondDemoTableSearch{
		List:     list,
		Total:    uint(total),
		PageNum:  uint(params.Extra.PageNum),
		PageSize: uint(params.Extra.PageSize),
	}
	c.SuccessData(resp)
}
