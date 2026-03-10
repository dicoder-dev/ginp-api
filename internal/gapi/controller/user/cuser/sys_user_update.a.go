package cuser

import (
	"ginp-api/internal/gapi/entity"
	"ginp-api/internal/gapi/service/user/suser"
	"ginp-api/pkg/where"

	"ginp-api/pkg/ginp"
)

type RequestSysUserUpdate struct {
	entity.User
}

type RespondSysUserUpdate struct {
}

func init() {
	ginp.RouterAppend(ginp.RouterItem{
		Path:           "/api/sys_user/update",
		Handler:        ginp.BindParamsHandler(SysUserUpdate, RequestSysUserUpdate{}),
		HttpType:       ginp.HttpPost,
		NeedLogin:      true,
		NeedPermission: true,
		PermissionName: "SysUse.update",
		Swagger: &ginp.SwaggerInfo{
			Title:         "modify user",
			Description:   "",
			RequestParams: RequestSysUserUpdate{},
		},
	})
}

func SysUserUpdate(c *ginp.ContextPlus, params *RequestSysUserUpdate) {
	wheres := where.Format(where.OptEqual("id", params.User.ID))
	err := suser.Model().Update(wheres, &params.User)
	if err != nil {
		c.Fail("修改失败" + err.Error())
		return
	}
	c.Success()
}
