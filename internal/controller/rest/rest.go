package rest

import (
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

func NewRESTHandler(globalConfig *config.Config, userRepo repository.UserRepository, boardRepo repository.BoardRepository) http.Handler {
	mux := http.NewServeMux()
	rest := newRestHandler(globalConfig, userRepo, boardRepo)
	rest.buildRoute(mux)
	return mux
}

type restHandler struct {
	logger       logging.Logger
	globalConfig *config.Config
	userRepo     repository.UserRepository
	boardRepo    repository.BoardRepository
}

func newRestHandler(globalConfig *config.Config, userRepo repository.UserRepository, boardRepo repository.BoardRepository) *restHandler {
	return &restHandler{
		logger:       logging.NewLogger(),
		globalConfig: globalConfig,
		userRepo:     userRepo,
		boardRepo:    boardRepo,
	}
}

func (rest *restHandler) buildRoute(mux *http.ServeMux) {
	mux.HandleFunc("/v1/token", rest.routeToken)
	mux.HandleFunc("/v1/boards", rest.routeBoards)
	mux.HandleFunc("/v1/boards/", rest.routeBoards)
	mux.HandleFunc("/v1/classes/", rest.routeClasses)
	mux.HandleFunc("/v1/users/", rest.routeUsers)
}
