package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	dhttp "github.com/Ptt-official-app/Ptt-backend/internal/delivery/http"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"

	_ "github.com/PichuChen/postgresql-gobbs"
	"github.com/Ptt-official-app/go-bbs"
	_ "github.com/Ptt-official-app/go-bbs/pttbbs"
)

func main() {
	var logLevel = flag.Uint("logLevel", 4, `log level: 0: Emergency; 1: Alert; 2: Critical; 3: Error; 4: Warning; 5: Notice; 6: Info; 7: Debug`)
	flag.Usage = func() {
		_, _ = os.Stderr.WriteString("Usage: \n  Ptt-backend [ options ]\n\n")
		_, _ = os.Stderr.WriteString("Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if _, ok := os.LookupEnv("LOG_LEVEL"); !ok {
		if *logLevel > 7 {
			*logLevel = 7
		}
	_ = os.Setenv("LOG_LEVEL", strconv.Itoa(int(*logLevel)))
	}

	logger := logging.NewLogger()
	logger.Informationalf("server start")

	globalConfig, err := config.NewDefaultConfig()
	if err != nil {
		logger.Errorf("failed to get config: %v", err)
		return
	}

	db, err := bbs.Open("postgresql", globalConfig.BBSHome)
	if err != nil {
		logger.Errorf("open bbs db error: %v", err)
		return
	}

	repo, err := repository.NewRepository(db)
	if err != nil {
		logger.Errorf("failed to create user repository: %s\n", err)
		return
	}
	usecase := usecase.NewUsecase(globalConfig, repo)
	go boardd(usecase)
	httpDelivery := dhttp.NewHTTPDelivery(usecase)
	if err := httpDelivery.Run(globalConfig.ListenPort); err != nil {
		logger.Errorf("run http delivery error: %s\n", err)
	}

}
