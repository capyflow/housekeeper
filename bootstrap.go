package main

import (
	"context"
	"flag"

	"github.com/capyflow/allspark-go/ds"
	"github.com/capyflow/allspark-go/system"
	"github.com/capyflow/housekeeper/internal/conf"
	"github.com/capyflow/housekeeper/internal/handler"
	"github.com/capyflow/housekeeper/internal/repository"
	"github.com/capyflow/housekeeper/internal/router"
	services "github.com/capyflow/housekeeper/internal/service"
	"github.com/capyflow/housekeeper/locale"
	vortex "github.com/capyflow/vortexv3"
)

func main() {
	var confPath = flag.String("c", "./internal/conf/config.toml", "config file path")
	flag.Parse()

	ctx := context.Background()

	if confPath == nil {
		panic("config file path is nil")
	}
	config, err := conf.LoadConfig(*confPath)
	if err != nil {
		panic("load config failed, err: " + err.Error())
	}

	hk := NewHousekeeper(ctx, config)
	hk.Start()
	system.GracefulShutdown(hk.Stop)
}

type Housekeeper struct {
	config       *conf.Config
	vortexEngine *vortex.VortexEngine
}

func NewHousekeeper(ctx context.Context, config *conf.Config) *Housekeeper {
	dsServer := ds.InitDatabaseServer(ctx, &ds.DsConfig{
		Redis: config.RdbConfig.Redis,
		Mongo: config.RdbConfig.Mongo,
	}, func(dbIdxs map[string]interface{}) {
		// redis 数据索引
		dbIdxs["note"] = 1

		// mongo 数据索引
		dbIdxs["housekeeper"] = "housekeeper"
	})

	// 初始化dao层
	NoteRepo := repository.NewNoteRepo(ctx, dsServer)

	// 初始化服务层
	NotesService := services.NewNotesService(ctx, NoteRepo)
	UserService := services.NewUserService(ctx, config)

	// 初始化handler层
	NotesHandler := handler.NewNotesHandler(ctx, NotesService)
	UserHandler := handler.NewUserHandler(ctx, UserService)

	routers := router.PrepareRouter(NotesHandler, UserHandler)

	e := vortex.NewVortexEngine(ctx,
		vortex.WithPort(int(config.Port)),
		vortex.WithHttpRouterRootGroup(routers),
		vortex.WithJwtOption(config.Jwt),
		vortex.WithI18n(locale.I18nValue.GetMap()))

	return &Housekeeper{
		config:       config,
		vortexEngine: e,
	}
}

func (hk *Housekeeper) Start() {
	hk.vortexEngine.Start()
}

func (hk *Housekeeper) Stop(ctx context.Context) error {
	// TODO 实现程序退出
	return nil
}
