package repository

import (
	"context"

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
	mongoRecordCollection = "record"
)

type RecordRepo struct {
	ctx context.Context
	rdb *redis.Client
	mdb *mongo.Database
}

func NewRecordRepo(ctx context.Context, dsServer *ds.DatabaseServer) *RecordRepo {
	rdb, ok := dsServer.GetRedis("record")
	if !ok {
		panic("get redis failed")
	}
	mdb, ok := dsServer.GetMongo("housekeeper")
	if !ok {
		panic("get mongo failed")
	}
	rr := &RecordRepo{
		ctx: ctx,
		rdb: rdb,
		mdb: mdb.Client().Database("housekeeper"),
	}
	rr.prepareIndex()
	return rr
}

// 准备索引
func (rr *RecordRepo) prepareIndex() {
	_, err := rr.mdb.Collection(mongoRecordCollection).Indexes().CreateOne(rr.ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("create index failed, err: " + err.Error())
	}
}

// 创建打卡记录
func (rr *RecordRepo) CreateRecord(ctx context.Context, record *model.Record) error {
	if len(record.Id) == 0 {
		record.Id = "r_" + pkg.GenerateRandomString(15)
	}
	_, err := rr.mdb.Collection(mongoRecordCollection).InsertOne(ctx, record)
	if err != nil {
		logx.Errorf("RecordRepo|CreateRecord|InsertOne|Error|%v|%s", err, conv.ToJsonWithoutError(record))
		return err
	}
	return nil
}

// 更新打卡记录（仅支持 title 和 record_ts）
func (rr *RecordRepo) UpdateRecord(ctx context.Context, req *model.UpdateRecordReq) error {
	filter := bson.M{"id": req.Id}
	update := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"record_ts": req.RecordTs,
		},
	}
	result, err := rr.mdb.Collection(mongoRecordCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		logx.Errorf("RecordRepo|UpdateRecord|UpdateOne|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return err
	}
	if result.MatchedCount == 0 {
		logx.Errorf("RecordRepo|UpdateRecord|NotFound|%s", req.Id)
		return mongo.ErrNoDocuments
	}
	return nil
}

// 删除打卡记录
func (rr *RecordRepo) DeleteRecord(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	result, err := rr.mdb.Collection(mongoRecordCollection).DeleteOne(ctx, filter)
	if err != nil {
		logx.Errorf("RecordRepo|DeleteRecord|DeleteOne|Error|%v|%s", err, id)
		return err
	}
	if result.DeletedCount == 0 {
		logx.Errorf("RecordRepo|DeleteRecord|NotFound|%s", id)
		return mongo.ErrNoDocuments
	}
	return nil
}

// 查询打卡记录详情
func (rr *RecordRepo) GetRecord(ctx context.Context, id string) (*model.Record, error) {
	filter := bson.M{"id": id}
	var record model.Record
	if err := rr.mdb.Collection(mongoRecordCollection).FindOne(ctx, filter).Decode(&record); err != nil {
		logx.Errorf("RecordRepo|GetRecord|FindOne|Error|%v|%s", err, id)
		return nil, err
	}
	return &record, nil
}

// 分页查询打卡记录
func (rr *RecordRepo) ListRecord(ctx context.Context, req *model.ListRecordReq) ([]model.Record, int64, error) {
	filter := bson.M{}
	total, err := rr.mdb.Collection(mongoRecordCollection).CountDocuments(ctx, filter)
	if err != nil {
		logx.Errorf("RecordRepo|ListRecord|CountDocuments|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}

	skip := int64((req.Page - 1) * req.PageSize)
	limit := int64(req.PageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "create_ts", Value: -1}})

	cursor, err := rr.mdb.Collection(mongoRecordCollection).Find(ctx, filter, opts)
	if err != nil {
		logx.Errorf("RecordRepo|ListRecord|Find|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	list := make([]model.Record, 0)
	if err = cursor.All(ctx, &list); err != nil {
		logx.Errorf("RecordRepo|ListRecord|All|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}

	return list, total, nil
}
