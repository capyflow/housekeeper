package router

import (
	"net/http"
	"strings"

	"github.com/capyflow/housekeeper/internal/handler"
	"github.com/capyflow/housekeeper/internal/service"
	vhttp "github.com/capyflow/vortexv3/server/http"
)

var (
	// 中间件解析token
	jwtParseMiddleware = func(userService *service.UserService) vhttp.VortexHttpMiddleware {
		return func(ctx *vhttp.Context) error {
			tokenString := ctx.GetHeader("Authorization")
			if tokenString == "" {
				return ctx.JsonResponse(vhttp.Codes.Unauthorized, map[string]any{
					"msg": "token is required",
				})
			}
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			claims, err := userService.ParseToken(tokenString, userService.SkipExpireCheck)
			if err != nil {
				return ctx.JsonResponse(vhttp.Codes.Unauthorized, map[string]any{
					"msg": "token is invalid",
				})
			}
			ctx.Set("claims", claims)
			return nil
		}
	}
)

func PrepareRouter(notesHandler *handler.NotesHandler,
	userHandler *handler.UserHandler,
	userService *service.UserService) *vhttp.VortexHttpRouterGroup {
	rootGroup := vhttp.NewRootGroup("/v1")
	prepareNotesRouter(rootGroup, notesHandler, userService)
	prepareUserRouter(rootGroup, userHandler)
	return rootGroup
}

func prepareNotesRouter(rootGroup *vhttp.VortexHttpRouterGroup,
	notesHandler *handler.NotesHandler,
	userService *service.UserService) {
	notesGroup := rootGroup.AddGroup("/notes")
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/create",
		notesHandler.HandlerCreateShareBoard, jwtParseMiddleware(userService))
	notesGroup.AddRouter([]string{http.MethodPut}, "/share_board/update",
		notesHandler.HandlerUpdateShareBoard, jwtParseMiddleware(userService))
	notesGroup.AddRouter([]string{http.MethodDelete}, "/share_board/delete",
		notesHandler.HandlerDeleteShareBoard, jwtParseMiddleware(userService))
	notesGroup.AddRouter([]string{http.MethodPost}, "/share_board/list",
		notesHandler.HandlerListShareBoard, jwtParseMiddleware(userService))
}

func prepareUserRouter(rootGroup *vhttp.VortexHttpRouterGroup, userHandler *handler.UserHandler) {
	userGroup := rootGroup.AddGroup("/user")
	userGroup.AddRouter([]string{http.MethodPost}, "/login", userHandler.Login)
}
