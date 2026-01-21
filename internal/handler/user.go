package handler

import (
	"context"

	"github.com/capyflow/allspark-go/conv"
	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/service"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

type UserHandler struct {
	ctx         context.Context
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(ctx context.Context, userService *service.UserService) *UserHandler {
	return &UserHandler{
		ctx:         ctx,
		userService: userService,
	}
}

// 登录
func (uh *UserHandler) Login(ctx *vhttp.Context) error {
	var req model.LoginReq
	err := ctx.UnmarshalBody(&req)
	if err != nil {
		logx.Errorf("UserHandler|Login|UnmarshalBody|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{
			"msg": "params invalid",
		})
	}

	// 调用service层登录方法
	resp, err := uh.userService.Login(uh.ctx, &req)
	if err != nil {
		logx.Errorf("UserHandler|Login|Service|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]any{
			"msg": "login failed",
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}
