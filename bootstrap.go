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
	var staticPath = flag.String("s", "./web/dist", "static file path")
	flag.Parse()

	ctx := context.Background()

	if confPath == nil {
		panic("config file path is nil")
	}
	config, err := conf.LoadConfig(*confPath)
	if err != nil {
		panic("load config failed, err: " + err.Error())
	}

	hk := NewHousekeeper(ctx, config, *staticPath)
	hk.Start()
	system.GracefulShutdown(hk.Stop)
}

type Housekeeper struct {
	config       *conf.Config
	vortexEngine *vortex.VortexEngine
}

func NewHousekeeper(ctx context.Context, config *conf.Config, staticPath string) *Housekeeper {
	dsServer := ds.InitDatabaseServer(ctx, &ds.DsConfig{
		Redis: config.RdbConfig.Redis,
		Mongo: config.RdbConfig.Mongo,
	}, func(dbIdxs map[string]interface{}) {
		// redis 数据索引
		dbIdxs["note"] = 1
		dbIdxs["record"] = 2
		dbIdxs["checkin"] = 3

		// mongo 数据索引
		dbIdxs["housekeeper"] = "housekeeper"
	})

	// 初始化 dao 层
	NoteRepo := repository.NewNoteRepo(ctx, dsServer)
	RecordRepo := repository.NewRecordRepo(ctx, dsServer)
	CheckInRepo := repository.NewCheckInRepo(ctx, dsServer)

	// 初始化服务层
	NotesService := services.NewNotesService(ctx, NoteRepo)
	RecordService := services.NewRecordService(ctx, RecordRepo)
	CheckInService := services.NewCheckInService(ctx, CheckInRepo)
	UserService := services.NewUserService(ctx, config)

	// 初始化 handler 层
	NotesHandler := handler.NewNotesHandler(ctx, NotesService)
	RecordHandler := handler.NewRecordHandler(ctx, RecordService)
	CheckInHandler := handler.NewCheckInHandler(ctx, CheckInService)
	UserHandler := handler.NewUserHandler(ctx, UserService)

	routers := router.PrepareRouter(NotesHandler, RecordHandler, CheckInHandler, UserHandler, staticPath)

	e := vortex.NewVortexEngine(ctx,
		vortex.WithPort(int(config.Port)),
		vortex.WithStaticDir(staticPath+"/**/*"), // LoadHTMLGlob 需要 glob 模式
		vortex.WithHttpRouterRootGroup(routers),  // 展开多个路由组
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
