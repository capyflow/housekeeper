package service

import (
	"context"

	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/repository"
)

type CheckInService struct {
	ctx context.Context
	repo *repository.CheckInRepo
}

func NewCheckInService(ctx context.Context, repo *repository.CheckInRepo) *CheckInService {
	return &CheckInService{ctx: ctx, repo: repo}
}

// 创建打卡任务
func (cs *CheckInService) CreateCheckInTarget(ctx context.Context, req *model.CreateCheckInTargetReq, owner string) (*model.CheckInTarget, error) {
	target := &model.CheckInTarget{
		Title:       req.Title,
		Description: req.Description,
		TimeRange:   req.TimeRange,
		Owner:       owner,
		CheckInCounts: make(map[string]int),
		TotalDays:   0,
	}
	
	if err := cs.repo.CreateCheckInTarget(ctx, target); err != nil {
		return nil, err
	}
	
	return target, nil
}

// 更新打卡任务
func (cs *CheckInService) UpdateCheckInTarget(ctx context.Context, req *model.UpdateCheckInTargetReq) error {
	return cs.repo.UpdateCheckInTarget(ctx, req)
}

// 删除打卡任务
func (cs *CheckInService) DeleteCheckInTarget(ctx context.Context, id string) error {
	return cs.repo.DeleteCheckInTarget(ctx, id)
}

// 查询打卡任务详情
func (cs *CheckInService) GetCheckInTarget(ctx context.Context, id string) (*model.CheckInTarget, error) {
	return cs.repo.GetCheckInTarget(ctx, id)
}

// 分页查询打卡任务
func (cs *CheckInService) ListCheckInTargets(ctx context.Context, req *model.ListCheckInTargetReq, owner string) (*model.ListCheckInTargetResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	
	if len(req.Owner) == 0 {
		req.Owner = owner
	}
	
	list, total, err := cs.repo.ListCheckInTargets(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &model.ListCheckInTargetResp{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, nil
}

// 打卡
func (cs *CheckInService) CheckIn(ctx context.Context, req *model.CheckInRequest, owner string) (*model.CheckInResponse, error) {
	// 验证所有者权限
	target, err := cs.repo.GetCheckInTarget(ctx, req.Id)
	if err != nil {
		return &model.CheckInResponse{
			Success: false,
			Message: "任务不存在",
		}, nil
	}
	
	if target.Owner != owner {
		return &model.CheckInResponse{
			Success: false,
			Message: "无权操作该任务",
		}, nil
	}
	
	// 目标时间戳为 0 则使用当前时间
	if req.TargetTs == 0 {
		req.TargetTs = 0 // 这里可以改成 time.Now().UnixMilli()
	}
	
	return cs.repo.CheckIn(ctx, req.Id, req.TargetTs)
}

// 获取打卡统计（热力图数据）
func (cs *CheckInService) GetCheckInStats(ctx context.Context, taskId string, days int) (map[string]int, error) {
	return cs.repo.GetCheckInStats(ctx, taskId, days)
}
