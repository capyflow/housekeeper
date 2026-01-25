package handler

import (
	"context"
	"errors"

	"github.com/capyflow/allspark-go/conv"
	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/service"
	"github.com/capyflow/housekeeper/pkg"
	vpkg "github.com/capyflow/vortexv3/pkg"
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
		return ctx.JsonResponse(vhttp.Codes.BadRequest, vpkg.VMsgResponse{
			"msg": "params invalid",
		})
	}

	// 调用service层登录方法
	resp, err := uh.userService.Login(ctx, &req)
	if err != nil {
		logx.Errorf("UserHandler|Login|Service|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		if errors.Is(err, pkg.ErrorEnum.ErrPasswordInvalid) {
			return ctx.JsonResponse(vhttp.Codes.BadRequest.WithSubCode(pkg.StateCodes.PasswordNotMatch), vpkg.VMsgResponse{
				"msg": "password not match",
			})
		} else if errors.Is(err, pkg.ErrorEnum.ErrUsernameInvalid) {
			return ctx.JsonResponse(vhttp.Codes.BadRequest.WithSubCode(pkg.StateCodes.UsernameNotMatch), vpkg.VMsgResponse{
				"msg": "username not match",
			})
		} else {
			return ctx.JsonResponse(vhttp.Codes.InternalError, vpkg.VMsgResponse{
				"msg": "login failed",
			})
		}
	}
	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}
