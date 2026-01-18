package router

import (
	"net/http"

	"github.com/capyflow/housekeeper/internal/handler"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

func PrepareRouter(notesHandler *handler.NotesHandler) *vhttp.VortexHttpRouterGroup {
	rootGroup := vhttp.NewRootGroup("/v1")
	prepareNotesRouter(rootGroup, notesHandler)
	return rootGroup
}

func prepareNotesRouter(rootGroup *vhttp.VortexHttpRouterGroup, notesHandler *handler.NotesHandler) {
	notesGroup := rootGroup.AddGroup("/notes")
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/create", notesHandler.HandlerCreateShareBoard)
	notesGroup.AddRouter([]string{http.MethodPut}, "/share_board/update", notesHandler.HandlerUpdateShareBoard)
	notesGroup.AddRouter([]string{http.MethodDelete}, "/share_board/delete", notesHandler.HandlerDeleteShareBoard)
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/list", notesHandler.HandlerListShareBoard)
}
