package handler

import (
	"context"

	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/service"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

type CheckInHandler struct {
	ctx context.Context
	svc *service.CheckInService
}

func NewCheckInHandler(ctx context.Context, svc *service.CheckInService) *CheckInHandler {
	return &CheckInHandler{ctx: ctx, svc: svc}
}

// 创建打卡任务
func (ch *CheckInHandler) HandlerCreateCheckInTarget(ctx *vhttp.Context) error {
	var req model.CreateCheckInTargetReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("CheckInHandler|HandlerCreateCheckInTarget|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{"Message": "params is invalid"})
	}

	if len(req.Title) == 0 {
		logx.Errorf("CheckInHandler|HandlerCreateCheckInTarget|title is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]any{"Message": "title is required"})
	}

	session := ctx.GetSession()
	
	target, err := ch.svc.CreateCheckInTarget(ctx.Context(), &req, session.UID)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerCreateCheckInTarget|CreateCheckInTarget failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]any{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]any{"info": target})
}

// 更新打卡任务
func (ch *CheckInHandler) HandlerUpdateCheckInTarget(ctx *vhttp.Context) error {
	var req model.UpdateCheckInTargetReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("CheckInHandler|HandlerUpdateCheckInTarget|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}

	if len(req.Id) == 0 {
		logx.Errorf("CheckInHandler|HandlerUpdateCheckInTarget|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "id is required"})
	}

	err := ch.svc.UpdateCheckInTarget(ctx.Context(), &req)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerUpdateCheckInTarget|UpdateCheckInTarget failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]string{"Message": "success"})
}

// 删除打卡任务
func (ch *CheckInHandler) HandlerDeleteCheckInTarget(ctx *vhttp.Context) error {
	var req model.DeleteCheckInTargetReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("CheckInHandler|HandlerDeleteCheckInTarget|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}

	if len(req.Id) == 0 {
		logx.Errorf("CheckInHandler|HandlerDeleteCheckInTarget|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "id is required"})
	}

	err := ch.svc.DeleteCheckInTarget(ctx.Context(), req.Id)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerDeleteCheckInTarget|DeleteCheckInTarget failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]string{"Message": "success"})
}

// 查询打卡任务详情
func (ch *CheckInHandler) HandlerCheckInTargetInfo(ctx *vhttp.Context) error {
	var req model.CheckInTargetInfoReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("CheckInHandler|HandlerCheckInTargetInfo|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}

	if len(req.Id) == 0 {
		logx.Errorf("CheckInHandler|HandlerCheckInTargetInfo|id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "id is required"})
	}

	info, err := ch.svc.GetCheckInTarget(ctx.Context(), req.Id)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerCheckInTargetInfo|GetCheckInTarget failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, map[string]any{"info": info})
}

// 分页查询打卡任务
func (ch *CheckInHandler) HandlerListCheckInTargets(ctx *vhttp.Context) error {
	var req model.ListCheckInTargetReq
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("CheckInHandler|HandlerListCheckInTargets|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}

	session := ctx.GetSession()
	
	resp, err := ch.svc.ListCheckInTargets(ctx.Context(), &req, session.UID)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerListCheckInTargets|ListCheckInTargets failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}

// 打卡
func (ch *CheckInHandler) HandlerCheckIn(ctx *vhttp.Context) error {
	var req model.CheckInRequest
	if err := ctx.UnmarshalBody(&req); err != nil {
		logx.Errorf("CheckInHandler|HandlerCheckIn|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}

	if len(req.Id) == 0 {
		logx.Errorf("CheckInHandler|HandlerCheckIn|target_id is required")
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "target_id is required"})
	}

	session := ctx.GetSession()
	
	resp, err := ch.svc.CheckIn(ctx.Context(), &req, session.UID)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerCheckIn|CheckIn failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}

	if !resp.Success {
		return ctx.JsonResponse(vhttp.Codes.Success, resp)
	}

	return ctx.JsonResponse(vhttp.Codes.Success, resp)
}

// 获取打卡统计（热力图数据）
func (ch *CheckInHandler) HandlerGetCheckInStats(ctx *vhttp.Context) error {
	var req struct {
		Id     string `json:"id"`      // 任务 ID
		Days   int    `json:"days"`    // 统计天数
	}
	
	if err := ctx.UnmarshalBody(&req); err != nil || len(req.Id) == 0 {
		logx.Errorf("CheckInHandler|HandlerGetCheckInStats|UnmarshalBody failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.BadRequest, map[string]string{"Message": "params is invalid"})
	}
	
	if req.Days <= 0 {
		req.Days = 365 // 默认一年
	}
	
	stats, err := ch.svc.GetCheckInStats(ctx.Context(), req.Id, req.Days)
	if err != nil {
		logx.Errorf("CheckInHandler|HandlerGetCheckInStats|GetCheckInStats failed, err: %v", err)
		return ctx.JsonResponse(vhttp.Codes.InternalError, map[string]string{"Message": "internal server error"})
	}
	
	return ctx.JsonResponse(vhttp.Codes.Success, stats)
}
