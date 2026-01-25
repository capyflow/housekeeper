package repository

import (
	"context"
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
	mongoShareBoardCollection = "share_board"
	mongoNoteCollection       = "note"
)

type NoteRepo struct {
	ctx context.Context
	rdb *redis.Client
	mdb *mongo.Database
}

func NewNoteRepo(ctx context.Context, dsServer *ds.DatabaseServer) *NoteRepo {
	rdb, ok := dsServer.GetRedis("note")
	if !ok {
		panic("get redis failed")
	}
	mdb, ok := dsServer.GetMongo("housekeeper")
	if !ok {
		panic("get mongo failed")
	}
	nr := &NoteRepo{ctx: ctx, rdb: rdb, mdb: mdb.Client().Database("housekeeper")}
	nr.prepareIndex()
	return nr
}

// 准备索引
func (nr *NoteRepo) prepareIndex() {
	_, err := nr.mdb.Collection(mongoShareBoardCollection).Indexes().CreateOne(nr.ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		panic("create index failed, err: " + err.Error())
	}
}

// 创建共享看板
func (nr *NoteRepo) CreateShareBoard(ctx context.Context, shardBoard *model.ShareBoard) error {
	if len(shardBoard.Id) == 0 {
		shardBoard.Id = "sb_" + pkg.GenerateRandomString(15)
	}
	_, err := nr.mdb.Collection(mongoShareBoardCollection).InsertOne(ctx, shardBoard)
	if err != nil {
		logx.Errorf("NoteRepo|CreateShareBoard|InsertOne|Error|%v|%s", err, conv.ToJsonWithoutError(shardBoard))
		return err
	}
	return nil
}

// 更新共享看板
func (nr *NoteRepo) UpdateShareBoard(ctx context.Context, req *model.UpdateShareBoardReq) error {
	filter := bson.M{"id": req.Id}
	update := bson.M{
		"$set": bson.M{
			"board_name":  req.BoardName,
			"content":     req.Content,
			"device_name": req.DeviceName,
			"update_time": time.Now().UnixMilli(),
		},
	}
	result, err := nr.mdb.Collection(mongoShareBoardCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		logx.Errorf("NoteRepo|UpdateShareBoard|UpdateOne|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return err
	}
	if result.MatchedCount == 0 {
		logx.Errorf("NoteRepo|UpdateShareBoard|NotFound|%s", req.Id)
		return mongo.ErrNoDocuments
	}
	return nil
}

// 删除共享看板
func (nr *NoteRepo) DeleteShareBoard(ctx context.Context, id string) error {
	filter := bson.M{"id": id}
	result, err := nr.mdb.Collection(mongoShareBoardCollection).DeleteOne(ctx, filter)
	if err != nil {
		logx.Errorf("NoteRepo|DeleteShareBoard|DeleteOne|Error|%v|%s", err, id)
		return err
	}
	if result.DeletedCount == 0 {
		logx.Errorf("NoteRepo|DeleteShareBoard|NotFound|%s", id)
		return mongo.ErrNoDocuments
	}
	return nil
}

// 分页查询共享看板
func (nr *NoteRepo) ListShareBoard(ctx context.Context, req *model.ListShareBoardReq) ([]model.ShareBoard, int64, error) {
	// 构建查询条件
	filter := bson.M{}
	if len(req.Owner) > 0 {
		filter["owner"] = req.Owner
	}

	// 获取总数
	total, err := nr.mdb.Collection(mongoShareBoardCollection).CountDocuments(ctx, filter)
	if err != nil {
		logx.Errorf("NoteRepo|ListShareBoard|CountDocuments|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}

	// 分页查询
	skip := int64((req.Page - 1) * req.PageSize)
	limit := int64(req.PageSize)
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "create_time", Value: -1}})

	cursor, err := nr.mdb.Collection(mongoShareBoardCollection).Find(ctx, filter, opts)
	if err != nil {
		logx.Errorf("NoteRepo|ListShareBoard|Find|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	list := make([]model.ShareBoard, 0)
	if err = cursor.All(ctx, &list); err != nil {
		logx.Errorf("NoteRepo|ListShareBoard|All|Error|%v|%s", err, conv.ToJsonWithoutError(req))
		return nil, 0, err
	}

	return list, total, nil
}
