package service

import (
	"context"
	"time"

	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/internal/repository"
	"github.com/capyflow/housekeeper/pkg"
)

type RecordService struct {
	ctx        context.Context
	recordRepo *repository.RecordRepo
}

func NewRecordService(ctx context.Context, recordRepo *repository.RecordRepo) *RecordService {
	return &RecordService{ctx: ctx, recordRepo: recordRepo}
}

// 创建打卡记录
func (rs *RecordService) CreateRecord(ctx context.Context, req *model.CreateRecordReq) (*model.Record, error) {
	recordID := "r_" + pkg.GenerateRandomString(15)
	info := &model.Record{
		Id:       recordID,
		Title:    req.Title,
		CreateTs: time.Now().UnixMilli(),
		RecordTs: req.RecordTs,
	}
	return info, rs.recordRepo.CreateRecord(ctx, info)
}

// 更新打卡记录（仅支持 title 和 record_ts）
func (rs *RecordService) UpdateRecord(ctx context.Context, req *model.UpdateRecordReq) error {
	return rs.recordRepo.UpdateRecord(ctx, req)
}

// 删除打卡记录
func (rs *RecordService) DeleteRecord(ctx context.Context, id string) error {
	return rs.recordRepo.DeleteRecord(ctx, id)
}

// 查询打卡记录详情
func (rs *RecordService) GetRecord(ctx context.Context, id string) (*model.Record, error) {
	return rs.recordRepo.GetRecord(ctx, id)
}

// 分页查询打卡记录
func (rs *RecordService) ListRecord(ctx context.Context, req *model.ListRecordReq) (*model.ListRecordResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	list, total, err := rs.recordRepo.ListRecord(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.ListRecordResp{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, nil
}
