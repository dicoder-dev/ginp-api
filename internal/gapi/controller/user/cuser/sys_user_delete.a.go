package cuser

import (
	"ginp-api/internal/gapi/dto/comdto"
	"ginp-api/internal/gapi/service/user/suser"

	"ginp-api/pkg/ginp"
)

type RequestSysUserDelete struct {
	ID uint `json:"id"`
}

type RespondSysUserDelete struct {
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/sys_user/delete",
		Handler:        ginp.BindParamsHandler(SysUserDelete, RequestSysUserDelete{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      true,
		NeedPermission: true,
		PermissionName: "SysUse.delete",
		Swagger: &ginp.SwaggerInfo{
			Title:         "delete user",
			Description:   "",
			RequestParams: comdto.ReqDelete{},
		},
	})
}

func SysUserDelete(c *ginp.ContextPlus, params *RequestSysUserDelete) {
	err := suser.Model().DeleteById(params.ID)
	if err != nil {
		c.Fail("删除失败" + err.Error())
		return
	}
	c.Success()
}
