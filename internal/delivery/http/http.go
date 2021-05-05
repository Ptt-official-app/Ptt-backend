package http

import (
	"fmt"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

// TODO: explain what Devlivery do
type Delivery struct {
	logger  logging.Logger
	usecase usecase.Usecase
}

// TODO: explain what this method to
func NewHTTPDelivery(usecase usecase.Usecase, logger logging.Logger) *Delivery {
	delivery := &Delivery{
		logger:  logger,
		usecase: usecase,
	}
	return delivery
}

// TODO: explain what this method to
func (delivery *Delivery) Run(port int16) error {
	mux := http.NewServeMux()
	delivery.buildRoute(mux)

	delivery.logger.Informationalf("listen port on %v", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), mux)
}
