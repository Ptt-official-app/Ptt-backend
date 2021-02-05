package http

import (
	"fmt"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

type httpDelivery struct {
	logger  logging.Logger
	usecase usecase.Usecase
}

func NewHTTPDelivery(usecase usecase.Usecase) *httpDelivery {
	delivery := &httpDelivery{
		logger:  logging.NewLogger(),
		usecase: usecase,
	}
	return delivery
}

func (delivery *httpDelivery) Run(port int16) error {
	mux := http.NewServeMux()
	delivery.buildRoute(mux)

	delivery.logger.Informationalf("listen port on %v", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
