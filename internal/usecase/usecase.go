package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	"github.com/Ptt-official-app/Ptt-backend/internal/mail"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/go-bbs"
)

type Usecase interface {
	GetUserByID(ctx context.Context, userID string) (bbs.UserRecord, error)
	GetUserFavorites(ctx context.Context, userID string) ([]interface{}, error)
	GetUserInformation(ctx context.Context, userID string) (map[string]interface{}, error)
	GetUserArticles(ctx context.Context, userID string) ([]interface{}, error)
	GetUserPreferences(ctx context.Context, userID string) (map[string]string, error)
	GetUserComments(ctx context.Context, userID string) ([]bbs.UserCommentRecord, error)
	GetUserDrafts(ctx context.Context, userID, draftID string) (bbs.UserDraft, error)
	UpdateUserDraft(ctx context.Context, userID, draftID string, text []byte) (bbs.UserDraft, error)
	DeleteUserDraft(ctx context.Context, userID, draftID string) error

	GetBoardByID(ctx context.Context, boardID string) (bbs.BoardRecord, error)
	GetBoards(ctx context.Context, userID string) []bbs.BoardRecord
	CreateBoard(ctx context.Context, boardID, title string) (bbs.BoardRecord, error)
	GetPopularBoards(ctx context.Context) ([]bbs.BoardRecord, error)
	GetBoardPostsLimitation(ctx context.Context, boardID string) (*BoardPostLimitation, error)
	GetClasses(ctx context.Context, userID, classID string) ([]bbs.BoardRecord, error)
	GetBoardArticles(ctx context.Context, boardID string, offset, length uint, cond *ArticleSearchCond) []bbs.ArticleRecord
	GetBoardArticle(ctx context.Context, boardID, filename string) ([]byte, error)
	GetBoardTreasures(ctx context.Context, boardID string, treasuresID []string) []interface{}
	CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error)
	GetRawArticle(boardID, filename string) (string, error)

	CreateAccessTokenWithUsername(username string) string
	GetUserIDFromToken(token string) (string, error)
	RecordLogin(userID, ip string)
	CheckPermission(token string, permissionID []Permission, userInfo map[string]string) error

	GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error)
	AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error)
	ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error)
	ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error
	UpdateUsefulness(ctx context.Context, userID, boardID, filename, appendType string) (repository.PushRecord, error)

	UpdateMail(mail mail.Mail) error
}

type SupportWebUsecase interface {
	GetArticleURL(boardID string, filename string) string
}

type loginRecord struct {
	at time.Time
	ip string
}

type usecase struct {
	logger       logging.Logger
	globalConfig *config.Config
	repo         repository.Repository
	mailProvider mail.Mail

	loginMu      sync.RWMutex
	loginRecords map[string]loginRecord
}

func NewUsecase(globalConfig *config.Config, repo repository.Repository) Usecase {
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
