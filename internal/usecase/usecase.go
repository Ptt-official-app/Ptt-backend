package usecase

import (
	"sync"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/mail"
)

type loginRecord struct {
	at time.Time
	ip string
}

type usecase struct {
	logger       logging.Logger
	globalConfig *config.Config
	repo         Repository
	mailProvider mail.Mail

	loginMu      sync.RWMutex
	loginRecords map[string]loginRecord
}

func NewUsecase(globalConfig *config.Config, repo Repository) *usecase {
	mailProvider, _ := mail.NewMailProvider(globalConfig.MailDriver)
	return &usecase{
		logger:       logging.NewLogger(),
		globalConfig: globalConfig,
		repo:         repo,
		mailProvider: mailProvider,
		loginRecords: make(map[string]loginRecord),
	}
}

func (usecase *usecase) RecordLogin(userID, ip string) {
	usecase.recordLoginAt(userID, ip, time.Now())
}

func (usecase *usecase) recordLoginAt(userID, ip string, at time.Time) {
	usecase.loginMu.Lock()
	defer usecase.loginMu.Unlock()
	usecase.loginRecords[userID] = loginRecord{at: at, ip: ip}
}

func (usecase *usecase) getLoginRecord(userID string) (loginRecord, bool) {
	usecase.loginMu.RLock()
	defer usecase.loginMu.RUnlock()
	record, ok := usecase.loginRecords[userID]
	return record, ok
}
