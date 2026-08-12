package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
)

// TODO: explain what Devlivery do
type Delivery struct {
	logger  logging.Logger
	usecase Usecase
}

// TODO: explain what this method to
func NewHTTPDelivery(usecase Usecase) *Delivery {
	delivery := &Delivery{
		logger:  logging.NewLogger(),
		usecase: usecase,
	}
	return delivery
}

// TODO: explain what this method to
func (delivery *Delivery) Run(port uint16) error {
	mux := http.NewServeMux()
	delivery.buildRoute(mux)

	delivery.logger.Informationalf("listen port on %v", port)
	srv := http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
	}
	return srv.ListenAndServe()
}
