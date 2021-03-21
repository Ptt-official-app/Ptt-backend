package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

type httpDelivery struct {
	*mux.Router
	logger  logging.Logger
	usecase usecase.Usecase
}

func NewHTTPDelivery(usecase usecase.Usecase) *httpDelivery {
	delivery := &httpDelivery{
		Router:  mux.NewRouter().StrictSlash(true), // redirect /path/ to /path
		logger:  logging.NewLogger(),
		usecase: usecase,
	}
	delivery.buildRoute()
	return delivery
}

func (delivery *httpDelivery) Run(port int16) error {
	delivery.logger.Informationalf("listen port on %v", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), delivery)
}
