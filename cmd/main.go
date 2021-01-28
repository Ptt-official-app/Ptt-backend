package main

import (
	"fmt"
	"net/http"

	"github.com/PichuChen/go-bbs"
	_ "github.com/PichuChen/go-bbs/pttbbs"
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/delivery/rest"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
)

func main() {
	var logger = logging.NewLogger()

	logger.Informationalf("server start")

	globalConfig, err := config.NewDefaultConfig()
	if err != nil {
		logger.Errorf("failed to get config: %v", err)
		return
	}

	db, err := bbs.Open("pttbbs", globalConfig.BBSHome)
	if err != nil {
		logger.Errorf("open bbs db error: %v", err)
		return
	}

	boardRepo, err := repository.NewBoardRepository(db)
	if err != nil {
		logger.Errorf("failed to create board repository: %s\n", err)
		return
	}

	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		logger.Errorf("failed to create user repository: %s\n", err)
		return
	}

	userUsecase := usecase.NewUserUsecase(userRepo)

	r := rest.NewRESTHandler(globalConfig, userUsecase, boardRepo)

	logger.Informationalf("listen port on %v", globalConfig.ListenPort)
	err = http.ListenAndServe(fmt.Sprintf(":%v", globalConfig.ListenPort), r)
	if err != nil {
		logger.Errorf("listen serve error: %v", err)
	}
}
