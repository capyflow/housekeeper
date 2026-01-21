package router

import (
	"net/http"

	"github.com/capyflow/housekeeper/internal/handler"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

func PrepareRouter(notesHandler *handler.NotesHandler, userHandler *handler.UserHandler) *vhttp.VortexHttpRouterGroup {
	rootGroup := vhttp.NewRootGroup("/v1")
	prepareNotesRouter(rootGroup, notesHandler)
	prepareUserRouter(rootGroup, userHandler)
	return rootGroup
}

func prepareNotesRouter(rootGroup *vhttp.VortexHttpRouterGroup, notesHandler *handler.NotesHandler) {
	notesGroup := rootGroup.AddGroup("/notes")
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/create", notesHandler.HandlerCreateShareBoard)
	notesGroup.AddRouter([]string{http.MethodPut}, "/share_board/update", notesHandler.HandlerUpdateShareBoard)
	notesGroup.AddRouter([]string{http.MethodDelete}, "/share_board/delete", notesHandler.HandlerDeleteShareBoard)
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/list", notesHandler.HandlerListShareBoard)
}

func prepareUserRouter(rootGroup *vhttp.VortexHttpRouterGroup, userHandler *handler.UserHandler) {
	userGroup := rootGroup.AddGroup("/user")
	userGroup.AddRouter([]string{http.MethodPost}, "/login", userHandler.Login)
}
