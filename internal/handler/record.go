package handler

import (
	"context"

	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/service"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

type RecordHandler struct {
	ctx       context.Context
	recordSvc *service.RecordService
}

func NewRecordHandler(ctx context.Context, recordSvc *service.RecordService) *RecordHandler {
	return &RecordHandler{ctx: ctx, recordSvc: recordSvc}
}

// 创建打卡记录
func (rh *RecordHandler) HandlerCreateRecord(ctx *vhttp.Context) error {
	var req model.CreateRecordReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("RecordHandler|HandlerCreateRecord|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{"Message": "params is invalid"})
	}
	if len(req.Title) == 0 {
		logx.Errorf("RecordHandler|HandlerCreateRecord|title is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{"Message": "title is required"})
	}

	info, err := rh.recordSvc.CreateRecord(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("RecordHandler|HandlerCreateRecord|CreateRecord failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]any{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]any{"info": info})
}

// 更新打卡记录
func (rh *RecordHandler) HandlerUpdateRecord(ctx *vhttp.Context) error {
	var req model.UpdateRecordReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("RecordHandler|HandlerUpdateRecord|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}
	if len(req.Id) == 0 {
		logx.Errorf("RecordHandler|HandlerUpdateRecord|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "id is required"})
	}

	err := rh.recordSvc.UpdateRecord(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("RecordHandler|HandlerUpdateRecord|UpdateRecord failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]string{"Message": "success"})
}

// 删除打卡记录
func (rh *RecordHandler) HandlerDeleteRecord(ctx *vhttp.Context) error {
	var req model.DeleteRecordReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("RecordHandler|HandlerDeleteRecord|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}
	if len(req.Id) == 0 {
		logx.Errorf("RecordHandler|HandlerDeleteRecord|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "id is required"})
	}

	err := rh.recordSvc.DeleteRecord(ctx.Context(), req.Id)
	if err != nil {
		logx.Errorf("RecordHandler|HandlerDeleteRecord|DeleteRecord failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]string{"Message": "success"})
}

// 查询打卡记录详情
func (rh *RecordHandler) HandlerRecordInfo(ctx *vhttp.Context) error {
	var req model.RecordInfoReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("RecordHandler|HandlerRecordInfo|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}
	if len(req.Id) == 0 {
		logx.Errorf("RecordHandler|HandlerRecordInfo|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "id is required"})
	}

	info, err := rh.recordSvc.GetRecord(ctx.Context(), req.Id)
	if err != nil {
		logx.Errorf("RecordHandler|HandlerRecordInfo|GetRecord failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]any{"info": info})
}

// 分页查询打卡记录
func (rh *RecordHandler) HandlerListRecord(ctx *vhttp.Context) error {
	var req model.ListRecordReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("RecordHandler|HandlerListRecord|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}

	resp, err := rh.recordSvc.ListRecord(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("RecordHandler|HandlerListRecord|ListRecord failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}
