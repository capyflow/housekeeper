package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/capyflow/allspark-go/conv"
	"github.com/capyflow/allspark-go/ds"
	"github.com/capyflow/allspark-go/logx"
	"github.com/capyflow/housekeeper/internal/model"
	"github.com/capyflow/housekeeper/pkg"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoCheckInCollection = "checkin_target"
)

type CheckInRepo struct {
	ctx context.Context
	rdb *redis.Client
	mdb *mongo.Database
}

func NewCheckInRepo(ctx context.Context, dsServer *ds.DatabaseServer) *CheckInRepo {
	rdb, ok := dsServer.GetRedis("checkin")
	if !ok {
		panic("get redis failed")
	}
	mdb, ok := dsServer.GetMongo("housekeeper")
	if !ok {
		panic("get mongo failed")
	}
	cr := &CheckInRepo{
		ctx: ctx,
		rdb: rdb,
		mdb: mdb.Client().Database("housekeeper"),
	}
	cr.prepareIndex()
	return cr
}

// 准备索引
func (cr *CheckInRepo) prepareIndex() {
	_, err := cr.mdb.Collection(mongoCheckInCollection).Indexes().CreateOne(cr.ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("create index failed, err: " + err.Error())
	}
	
	// 为 owner 字段添加索引，支持按所有者查询
	_, err = cr.mdb.Collection(mongoCheckInCollection).Indexes().CreateOne(cr.ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "owner", Value: 1}, {Key: "create_time", Value: -1}},
	})
	if err != nil {
		logx.Errorf("create owner index failed, err: %v", err)
	}
}

// 创建打卡任务
func (cr *CheckInRepo) CreateCheckInTarget(ctx context.Context, target *model.CheckInTarget) error {
	if len(target.Id) == 0 {
		target.Id = "ct_" + pkg.GenerateRandomString(15)
	}
	target.CreateTime = time.Now().UnixMilli()
	target.UpdateTime = target.CreateTime
	
	// 初始化 check_in_counts
	if target.CheckInCounts == nil {
		target.CheckInCounts = make(map[string]int)
	}
	
	_, err := cr.mdb.Collection(mongoCheckInCollection).InsertOne(ctx, target)
	if err != nil {
		logx.Errorf("CheckInRepo|CreateCheckInTarget|InsertOne|Error|%v|%s", err, conv.ToJsonWithoutError(target))
		return err
	}
	return nil
}

// 更新打卡任务
func (cr *CheckInRepo) UpdateCheckInTarget(ctx context.Context, req *model.UpdateCheckInTargetReq) error {
	filter := bson.M{"id": req.Id}
	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"description": req.Description,
			"time_range":  req.TimeRange,
			"update_time": time.Now().UnixMilli(),
		},
	}
	result, err := cr.mdb.Collection(mongoCheckInCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		logx.Errorf("CheckInRepo|UpdateCheckInTarget|UpdateOne|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return err
	}
	if result.MatchedCount == 0 {
		logx.Errorf("CheckInRepo|UpdateCheckInTarget|NotFound|%s", req.Id)
		return mongo.ErrNoDocuments
	}
	return nil
}

// 删除打卡任务
func (cr *CheckInRepo) DeleteCheckInTarget(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	result, err := cr.mdb.Collection(mongoCheckInCollection).DeleteOne(ctx, filter)
	if err != nil {
		logx.Errorf("CheckInRepo|DeleteCheckInTarget|DeleteOne|Error|%v|%s", err, id)
		return err
	}
	if result.DeletedCount == 0 {
		logx.Errorf("CheckInRepo|DeleteCheckInTarget|NotFound|%s", id)
		return mongo.ErrNoDocuments
	}
	return nil
}

// 查询打卡任务详情
func (cr *CheckInRepo) GetCheckInTarget(ctx context.Context, id string) (*model.CheckInTarget, error) {
	filter := bson.M{"id": id}
	var target model.CheckInTarget
	if err := cr.mdb.Collection(mongoCheckInCollection).FindOne(ctx, filter).Decode(&target); err != nil {
		logx.Errorf("CheckInRepo|GetCheckInTarget|FindOne|Error|%v|%s", err, id)
		return nil, err
	}
	return &target, nil
}

// 分页查询打卡任务
func (cr *CheckInRepo) ListCheckInTargets(ctx context.Context, req *model.ListCheckInTargetReq) ([]model.CheckInTarget, int64, error) {
	filter := bson.M{}
	
	// 如果指定了 owner，则过滤
	if len(req.Owner) > 0 {
		filter["owner"] = req.Owner
	}
	
	total, err := cr.mdb.Collection(mongoCheckInCollection).CountDocuments(ctx, filter)
	if err != nil {
		logx.Errorf("CheckInRepo|ListCheckInTargets|CountDocuments|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}

	skip := int64((req.Page - 1) * req.PageSize)
	limit := int64(req.PageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := cr.mdb.Collection(mongoCheckInCollection).Find(ctx, filter, opts)
	if err != nil {
		logx.Errorf("CheckInRepo|ListCheckInTargets|Find|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	list := make([]model.CheckInTarget, 0)
	if err = cursor.All(ctx, &list); err != nil {
		logx.Errorf("CheckInRepo|ListCheckInTargets|All|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}

	return list, total, nil
}

// 打卡（记录打卡时间）
func (cr *CheckInRepo) CheckIn(ctx context.Context, taskId string, targetTs int64) (*model.CheckInResponse, error) {
	filter := bson.M{"id": taskId}
	
	// 获取任务信息
	var target model.CheckInTarget
	if err := cr.mdb.Collection(mongoCheckInCollection).FindOne(ctx, filter).Decode(&target); err != nil {
		logx.Errorf("CheckInRepo|CheckIn|FindOne|Error|%v|%s", err, taskId)
		return &model.CheckInResponse{
			Success: false,
			Message: "任务不存在",
		}, nil
	}
	
	// 检查时间段限制
	timeRange := target.TimeRange
	if timeRange.Enabled {
		nowHour := time.UnixMilli(targetTs).Hour()
		nowMinute := time.UnixMilli(targetTs).Minute()
		nowTotalMin := nowHour*60 + nowMinute
		
		startTotalMin := timeRange.Start
		endTotalMin := timeRange.End
		
		inRange := false
		if startTotalMin <= endTotalMin {
			// 同一天内
			inRange = nowTotalMin >= startTotalMin && nowTotalMin <= endTotalMin
		} else {
			// 跨天 (如 22:00-02:00)
			inRange = nowTotalMin >= startTotalMin || nowTotalMin <= endTotalMin
		}
		
		if !inRange {
			return &model.CheckInResponse{
				Success: false,
				Message: "当前时间不在允许的打卡时间段内",
				TargetId: taskId,
			}, nil
		}
	}
	
	// 归一化到当天开始
	dayStart := model.NormalizeDayStart(targetTs)
	dateKey := model.GetDateKey(dayStart)
	
	// 检查是否已经打过卡
	if count, exists := target.CheckInCounts[dateKey]; exists && count > 0 {
		return &model.CheckInResponse{
			Success:    false,
			Message:    "今日已打卡过",
			TargetId:   taskId,
			CheckInTime: targetTs,
		}, nil
	}
	
	// 先查询当前的 total_days
	currentTotalDays := len(target.CheckInCounts)
	
	// 更新打卡记录 - 简单直接的方式
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("check_in_counts.%s", dateKey): 1,
			"total_days":                                currentTotalDays + 1,
			"update_time":                               time.Now().UnixMilli(),
		},
	}
	
	result, err := cr.mdb.Collection(mongoCheckInCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		logx.Errorf("CheckInRepo|CheckIn|UpdateOne|Error|%v|%s", err, taskId)
		return &model.CheckInResponse{
			Success: false,
			Message: "打卡失败",
			TargetId: taskId,
		}, nil
	}
	
	if result.MatchedCount == 0 {
		return &model.CheckInResponse{
			Success: false,
			Message: "任务不存在",
			TargetId: taskId,
		}, nil
	}
	
	return &model.CheckInResponse{
		Success:     true,
		Message:     "打卡成功",
		TargetId:    taskId,
		CheckInTime: targetTs,
		TotalDays:   currentTotalDays + 1,
	}, nil
}

// 统计某段时间内的打卡情况
func (cr *CheckInRepo) GetCheckInStats(ctx context.Context, taskId string, days int) (map[string]int, error) {
	filter := bson.M{"id": taskId}
	var target model.CheckInTarget
	if err := cr.mdb.Collection(mongoCheckInCollection).FindOne(ctx, filter).Decode(&target); err != nil {
		return nil, err
	}
	
	return target.CheckInCounts, nil
}
