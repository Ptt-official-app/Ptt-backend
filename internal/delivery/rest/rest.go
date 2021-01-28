package rest

import (
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func NewRESTHandler(globalConfig *config.Config, userUsecase usecase.UserUsecase, boardRepo repository.BoardRepository) http.Handler {
	mux := http.NewServeMux()
	rest := newRestHandler(globalConfig, userUsecase, boardRepo)
	rest.buildRoute(mux)
	return mux
}

type restHandler struct {
	logger logging.Logger

	globalConfig *config.Config

	boardRepo repository.BoardRepository

	userUsecase usecase.UserUsecase
}

func newRestHandler(globalConfig *config.Config, userUsecase usecase.UserUsecase, boardRepo repository.BoardRepository) *restHandler {
	return &restHandler{
		logger:       logging.NewLogger(),
		globalConfig: globalConfig,
		boardRepo:    boardRepo,
		userUsecase:  userUsecase,
	}
}

func (rest *restHandler) buildRoute(mux *http.ServeMux) {
	mux.HandleFunc("/v1/token", rest.routeToken)
	mux.HandleFunc("/v1/boards", rest.routeBoards)
	mux.HandleFunc("/v1/boards/", rest.routeBoards)
	mux.HandleFunc("/v1/classes/", rest.routeClasses)
	mux.HandleFunc("/v1/users/", rest.routeUsers)
}
