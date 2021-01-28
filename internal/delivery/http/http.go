package http

import (
	"fmt"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

type httpDelivery struct {
	logger logging.Logger

	globalConfig *config.Config

	userUsecase  usecase.UserUsecase
	boardUsecase usecase.BoardUsecase
}

func NewHTTPDelivery(globalConfig *config.Config, userUsecase usecase.UserUsecase, boardRepo usecase.BoardUsecase) *httpDelivery {
	delivery := &httpDelivery{
		logger:       logging.NewLogger(),
		globalConfig: globalConfig,
		userUsecase:  userUsecase,
		boardUsecase: boardRepo,
	}
	return delivery
}

func (delivery *httpDelivery) Run() error {
	mux := http.NewServeMux()
	delivery.buildRoute(mux)

	delivery.logger.Informationalf("listen port on %v", delivery.globalConfig.ListenPort)
	return http.ListenAndServe(fmt.Sprintf(":%v", delivery.globalConfig.ListenPort), mux)
}

func (delivery *httpDelivery) buildRoute(mux *http.ServeMux) {
	mux.HandleFunc("/v1/token", delivery.routeToken)
	mux.HandleFunc("/v1/boards", delivery.routeBoards)
	mux.HandleFunc("/v1/boards/", delivery.routeBoards)
	mux.HandleFunc("/v1/classes/", delivery.routeClasses)
	mux.HandleFunc("/v1/users/", delivery.routeUsers)
}
