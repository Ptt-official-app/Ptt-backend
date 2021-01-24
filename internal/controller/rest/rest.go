package rest

import (
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

var (
	logger = logging.NewLogger()

	globalConfig *config.Config
	userRepo     repository.UserRepository
	boardRepo    repository.BoardRepository
)

func NewRESTHandler(config *config.Config, user repository.UserRepository, board repository.BoardRepository) http.Handler {
	// FIXME:
	globalConfig = config
	userRepo = user
	boardRepo = board

	r := http.NewServeMux()
	buildRoute(r)
	return r
}
