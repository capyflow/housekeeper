package repository

import (
	"context"
	"github.com/capyflow/allspark-go/ds"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionRepo struct {
	ctx context.Context
	rdb *redis.Client
	mdb *mongo.Database
}

// 创建收藏夹
func NewCollectionRepo(ctx context.Context, dsServer *ds.DatabaseServer) *CollectionRepo {
	rdb, ok := dsServer.GetRedis("collection")
	if !ok {
		panic("get redis failed")
	}
	mdb, ok := dsServer.GetMongo("housekeeper")
	if !ok {
		panic("get mongo failed")
	}
	return &CollectionRepo{
		ctx: ctx,
		rdb: rdb,
		mdb: mdb,
	}
}
