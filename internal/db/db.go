package db

import "github.com/PichuChen/go-bbs"

type DB interface {
	ReadBoardRecords() ([]bbs.BoardRecord, error)
	ReadUserRecords() ([]bbs.UserRecord, error)
	ReadBoardArticleFile(boardId string, filename string) ([]byte, error)
	ReadBoardArticleRecordsFile(boardId string) ([]bbs.ArticleRecord, error)
	ReadBoardTreasureRecordsFile(boardId string, treasureId []string) ([]bbs.ArticleRecord, error)
	ReadUserFavoriteRecords(userId string) ([]bbs.FavoriteRecord, error)
}
