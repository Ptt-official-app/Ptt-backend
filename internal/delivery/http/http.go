package http

import (
	"fmt"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

type httpDelivery struct {
	logger logging.Logger

	userUsecase  usecase.UserUsecase
	boardUsecase usecase.BoardUsecase
	tokenUsecase usecase.TokenUsecase
}

func NewHTTPDelivery(userUsecase usecase.UserUsecase, boardRepo usecase.BoardUsecase, tokenUsecase usecase.TokenUsecase) *httpDelivery {
	delivery := &httpDelivery{
		logger:       logging.NewLogger(),
		userUsecase:  userUsecase,
		boardUsecase: boardRepo,
		tokenUsecase: tokenUsecase,
	}
	return delivery
}

func (delivery *httpDelivery) Run(port int16) error {
	mux := http.NewServeMux()
	delivery.buildRoute(mux)

	delivery.logger.Informationalf("listen port on %v", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
