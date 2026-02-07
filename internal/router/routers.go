package router

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/capyflow/housekeeper/internal/handler"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

func PrepareRouter(notesHandler *handler.NotesHandler,
	userHandler *handler.UserHandler, staticPath string) *vhttp.VortexHttpRouterGroup {
	rootGroup := vhttp.NewRootGroup("")
	// API路由组（/v1前缀）
	apiGroup := rootGroup.AddGroup("/v1")
	prepareNotesRouter(apiGroup, notesHandler)
	prepareUserRouter(apiGroup, userHandler)
	prepareStaticRouter(rootGroup, staticPath)
	return rootGroup
}

func prepareStaticRouter(rootGroup *vhttp.VortexHttpRouterGroup, staticPath string) {
	// /webui 入口 - 返回 index.html
	rootGroup.AddStaticRouter("/webui", func(ctx *vhttp.Context) error {
		indexPath := filepath.Join(staticPath, "index.html")
		ctx.GinContext().File(indexPath)
		return nil
	})

	rootGroup.AddStaticRouter("/webui/*path", func(ctx *vhttp.Context) error {
		path := ctx.GinContext().Param("path")
		trimmedPath := strings.TrimPrefix(path, "/")
		if trimmedPath == "icon.png" {
			iconPath := filepath.Join(staticPath, "icon.png")
			ctx.GinContext().File(iconPath)
			return nil
		}
		if strings.HasPrefix(trimmedPath, "assets/") {
			filePath := strings.TrimPrefix(trimmedPath, "assets/")
			fullPath := filepath.Join(staticPath, "assets", filePath)
			ctx.GinContext().File(fullPath)
			return nil
		}
		indexPath := filepath.Join(staticPath, "index.html")
		ctx.GinContext().File(indexPath)
		return nil
	})

}

func prepareNotesRouter(rootGroup *vhttp.VortexHttpRouterGroup,
	notesHandler *handler.NotesHandler) {
	notesGroup := rootGroup.AddGroup("/notes")
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/create", notesHandler.HandlerCreateShareBoard)   // 创建分享
	notesGroup.AddRouter([]string{http.MethodPut}, "/share_board/update", notesHandler.HandlerUpdateShareBoard)    // 修改分享
	notesGroup.AddRouter([]string{http.MethodDelete}, "/share_board/delete", notesHandler.HandlerDeleteShareBoard) // 删除分享
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/list", notesHandler.HandlerListShareBoard)       // 获取分享

	notesGroup.AddRouter([]string{http.MethodPost}, "/note/create", notesHandler.HandlerCreateNote)   // 创建笔记
	notesGroup.AddRouter([]string{http.MethodPut}, "/note/update", notesHandler.HandlerUpdateNote)    // 修改笔记
	notesGroup.AddRouter([]string{http.MethodDelete}, "/note/delete", notesHandler.HandlerDeleteNote) // 删除笔记
	notesGroup.AddRouter([]string{http.MethodPost}, "/note/info", notesHandler.HandlerNoteInfo)       // 查询笔记详情
	notesGroup.AddRouter([]string{http.MethodPost}, "/note/list", notesHandler.HandlerListNote)       // 获取笔记

}

func prepareUserRouter(rootGroup *vhttp.VortexHttpRouterGroup, userHandler *handler.UserHandler) {
	userGroup := rootGroup.AddGroup("/user")
	userGroup.AddRouter([]string{http.MethodPost}, "/login", userHandler.Login, vhttp.WithSkipJwtVerify())
}
