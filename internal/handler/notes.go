package handler

import (
	"context"

	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/service"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

type NotesHandler struct {
	ctx      context.Context
	notesSvc *service.NotesService
}

func NewNotesHandler(ctx context.Context, notesSvc *service.NotesService) *NotesHandler {
	return &NotesHandler{ctx: ctx, notesSvc: notesSvc}
}

// 创建共享看板
func (nh *NotesHandler) HandlerCreateShareBoard(ctx *vhttp.Context) error {
	var req model.CreateShareBoardReq
	if err := ctx.UnmarshalBody(&req); nil != err {
		logx.Errorf("NotesHandler|HandlerCreateShareBoard|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{
			"Message": "params is invalid",
		})
	}

	session := ctx.GetSession()
	req.Owner = session.UID

	if len(req.BoardName) == 0 || len(req.Content) == 0 {
		logx.Errorf("NotesHandler|HandlerCreateShareBoard|params is invalid")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{
			"Message": "params is invalid",
		})
	}

	info, err := nh.notesSvc.CreateShareBoard(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("NotesHandler|HandlerCreateShareBoard|CreateShareBoard failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]any{
			"Message": "internal server error",
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]any{
		"info": info,
	})
}

// 更新共享看板
func (nh *NotesHandler) HandlerUpdateShareBoard(ctx *vhttp.Context) error {
	var req model.UpdateShareBoardReq
	if err := ctx.UnmarshalBody(&req); nil != err {
		logx.Errorf("NotesHandler|HandlerUpdateShareBoard|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{
			"Message": "params is invalid",
		})
	}

	if len(req.Id) == 0 {
		logx.Errorf("NotesHandler|HandlerUpdateShareBoard|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{
			"Message": "id is required",
		})
	}

	err := nh.notesSvc.UpdateShareBoard(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("NotesHandler|HandlerUpdateShareBoard|UpdateShareBoard failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{
			"Message": "internal server error",
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]string{
		"Message": "success",
	})
}

// 删除共享看板
func (nh *NotesHandler) HandlerDeleteShareBoard(ctx *vhttp.Context) error {
	var req model.DeleteShareBoardReq
	if err := ctx.UnmarshalBody(&req); nil != err {
		logx.Errorf("NotesHandler|HandlerDeleteShareBoard|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{
			"Message": "params is invalid",
		})
	}

	if len(req.Id) == 0 {
		logx.Errorf("NotesHandler|HandlerDeleteShareBoard|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{
			"Message": "id is required",
		})
	}

	err := nh.notesSvc.DeleteShareBoard(ctx.Context(), req.Id)
	if err != nil {
		logx.Errorf("NotesHandler|HandlerDeleteShareBoard|DeleteShareBoard failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{
			"Message": "internal server error",
		})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]string{
		"Message": "success",
	})
}

// 分页查询共享看板
func (nh *NotesHandler) HandlerListShareBoard(ctx *vhttp.Context) error {
	var req model.ListShareBoardReq
	if err := ctx.UnmarshalBody(&req); nil != err {
		logx.Errorf("NotesHandler|HandlerQueryShareBoard|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{
			"Message": "params is invalid",
		})
	}

	resp, err := nh.notesSvc.ListShareBoard(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("NotesHandler|HandlerListShareBoard|ListShareBoard failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{
			"Message": "internal server error",
		})
	}
	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}
