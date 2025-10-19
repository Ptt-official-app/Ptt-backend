package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"regexp"
	"strings"

	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	apipb "github.com/Ptt-official-app/Ptt-backend/internal/proto/api"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	webpttparser "github.com/Ptt-official-app/Ptt-backend/webptt_parser"
	"github.com/Ptt-official-app/go-bbs"
	"google.golang.org/grpc"

	"database/sql"

	_ "github.com/lib/pq"
)

type server struct {
	apipb.UnimplementedBoardServiceServer
	usecase usecase.Usecase
}

func boardd(usecase usecase.Usecase) {
	// flag.Parse()
	//
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3432))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	apipb.RegisterBoardServiceServer(s, &server{
		usecase: usecase,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Board(ctx context.Context, req *apipb.BoardRequest) (*apipb.BoardReply, error) {
	slog.Info("boardd::Board", "len(req.Ref)", len(req.Ref))
	initCacheBoards(s.usecase)
	boards := make([]*apipb.Board, len(req.Ref))
	for i, ref := range req.Ref {
		slog.Info("boardd::Board", "ref", ref.GetRef(), "bid", ref.GetBid(), "name", ref.GetName())
		var board bbs.BoardRecord
		var boardIndex uint32
		if ref.GetBid() > 0 {
			if int(ref.GetBid()) >= len(cachedBoard) {
				slog.Error("boardd::Board", "invalid bid", ref.GetBid(), "len(cachedBoard)", len(cachedBoard))
				return nil, fmt.Errorf("invalid bid: %d", ref.GetBid())
			}
			board = cachedBoard[ref.GetBid()]
		} else if ref.GetName() != "" {
			var ok bool
			boardIndex, ok = boardToBoardIndex[strings.ToLower(ref.GetName())]
			if !ok {
				slog.Error("boardd::Board board name not found", "invalid name", ref.GetName())
				return nil, fmt.Errorf("invalid name: %s", ref.GetName())
			}
			board = cachedBoard[boardIndex]
		} else {
			slog.Error("boardd::Board", "invalid ref", ref)
			return nil, fmt.Errorf("invalid ref: %v", ref)
		}
		slog.Info("boardd::Board", "board", board)
		// articles := s.usecase.GetBoardArticles(context.Background(), board.BoardID(), 0, ^uint(0), &usecase.ArticleSearchCond{})

		l := fmt.Sprintf("/%s/index.html", board.BoardID())
		b, err := webpttparser.GetPttPage(l)
		if err != nil {
			return nil, err
		}
		var page *webpttparser.BoardIndexPage
		page, err = webpttparser.ParseBoardIndexPage(b)
		if err != nil {
			slog.Error("boardd::List ParseBoardIndexPage error", "error", err)
			return nil, err
		}
		prevPageHref := page.PrevPage // 上頁
		reg := regexp.MustCompile(`/index(\d+).html`)
		matches := reg.FindStringSubmatch(prevPageHref)
		pageNum := 0
		if len(matches) == 2 {
			fmt.Sscanf(matches[1], "%d", &pageNum)
		} else {
			slog.Warn("boardd::Board unable to parse last page number", "lastPageHref", prevPageHref)
		}

		numPosts := pageNum*20 + len(page.Articles)

		slog.Info("boardd::Board", "boardName", board.BoardID(), "numPosts", numPosts)

		boards[i] = &apipb.Board{
			Bid:        boardToBoardIndex[strings.ToLower(board.BoardID())],
			Name:       board.BoardID(),
			Title:      board.Title(),
			NumUsers:   0,
			Bclass:     "",
			Attributes: 0,
			// NumPosts:   uint32(len(articles)),
			NumPosts: uint32(numPosts),
		}

		// boards[i] = &apipb.Board{
		// 	Bid:        ref.GetBid(),
		// 	Name:       "sysop",
		// 	Title:      "Mock Board Title",
		// 	NumUsers:   100,
		// 	Bclass:     "Mock Class",
		// 	Attributes: 0,
		// }
	}

	slog.Info("boardd::Board reply", "boards", boards)
	return &apipb.BoardReply{Boards: boards}, nil
}

func (s *server) List(ctx context.Context, req *apipb.ListRequest) (*apipb.ListReply, error) {
	slog.Info("boardd::List", "ref", req.Ref, "bid", req.Ref.GetBid(), "req", req)
	initCacheBoards(s.usecase)
	var boardName string
	if len(req.Ref.GetName()) > 0 {
		boardName = req.Ref.GetName()
	} else if req.Ref.GetBid() > 0 {
		if int(req.Ref.GetBid()) > len(cachedBoard) {
			slog.Error("boardd::List", "invalid bid", req.Ref.GetBid(), "len(cachedBoard)", len(cachedBoard))
			return nil, fmt.Errorf("invalid bid: %d", req.Ref.GetBid())
		}
		boardName = cachedBoard[req.Ref.GetBid()].BoardID()
	} else {
		slog.Error("boardd::List", "invalid ref", req.Ref)
		return nil, fmt.Errorf("invalid ref: %v", req.Ref)
	}
	slog.Info("boardd::List", "boardName", boardName)
	var offset = uint(req.GetOffset())
	var length = uint(req.GetLength())
	if length <= 0 {
		length = 20 // default length
	}

	var page *webpttparser.BoardIndexPage

	if offset == 0 && req.IncludeBottoms {
		// fetch bottoms only
		l := fmt.Sprintf("/%s/index.html", boardName)
		b, err := webpttparser.GetPttPage(l)
		if err != nil {
			return nil, err
		}
		slog.Info("boardd::List fetched page from webpttparser for bottoms", "boardName", boardName, "length", len(b))

		page, err = webpttparser.ParseBoardIndexPage(b)
		if err != nil {
			slog.Error("boardd::List ParseBoardIndexPage error", "error", err)
			return nil, err
		}
		slog.Info("boardd::List parsed page for bottoms", "boardName", boardName, "len(page.Bottoms)", len(page.Bottoms))
	} else {

		pageNum := offset/20 + 1
		slog.Info("boardd::List fetching page", "boardName", boardName, "pageNum", pageNum)
		l := fmt.Sprintf("/%s/index%d.html", boardName, pageNum)
		b, err := webpttparser.GetPttPage(l)
		if err != nil {
			return nil, err
		}
		slog.Info("boardd::List fetched page from webpttparser", "boardName", boardName, "length", len(b))

		page, err = webpttparser.ParseBoardIndexPage(b)
		if err != nil {
			slog.Error("boardd::List ParseBoardIndexPage error", "error", err)
			return nil, err
		}
		slog.Info("boardd::List parsed page", "boardName", boardName, "len(page.Articles)", len(page.Articles), "len(page.Bottoms)", len(page.Bottoms), "LastPage", page.LastPage)
	}
	// articles := s.usecase.GetBoardArticles(context.Background(), boardID, offset, length, &usecase.ArticleSearchCond{})
	// if length > uint(len(articles)) {
	// 	length = uint(len(articles)) // limit to the number of articles
	// }
	// if length > 20 {
	// 	length = 20 // max length
	// }
	posts := make([]*apipb.Post, len(page.Articles))
	for i, article := range page.Articles {
		// slog.Info("boardd::List", "article", article)
		posts[i] = &apipb.Post{
			Index:         uint32(i + 1),
			Filename:      article.FileName,
			RawDate:       article.Date,
			NumRecommends: int32(article.Recommend),
			Filemode:      markToMode(article.Mark),
			Owner:         article.Owner,
			Title:         article.Title,
			ModifiedNsec:  0,
		}
	}

	bottoms := make([]*apipb.Post, len(page.Bottoms))
	for i, article := range page.Bottoms {
		// slog.Info("boardd::List", "bottom article", article)
		bottoms[i] = &apipb.Post{
			Index:         uint32(i + 1 + len(page.Articles)),
			Filename:      article.FileName,
			RawDate:       article.Date,
			NumRecommends: int32(article.Recommend),
			Filemode:      markToMode(article.Mark),
			Owner:         article.Owner,
			Title:         article.Title,
			ModifiedNsec:  0,
		}
	}

	slog.Info("boardd::List", "len(posts)", len(posts))

	return &apipb.ListReply{
		Posts:   posts,
		Bottoms: bottoms,
	}, nil

}

// Non-mail file modes
const (
	FileLocal = 1 << iota
	FileMarked
	FileDigest
	FileBottom
	FileSolved
)

func markToMode(mark string) int32 {
	switch mark {
	case "!":
		return FileMarked | FileSolved
	case "M":
		return FileMarked
	case "S":
		return FileSolved
	}
	return 0
}

func (s *server) Content(ctx context.Context, req *apipb.ContentRequest) (*apipb.ContentReply, error) {
	slog.Info("boardd::Content", "boardref", req.BoardRef, "filename", req.Filename, "token", req.ConsistencyToken, "options", req.PartialOptions)
	initCacheBoards(s.usecase)
	var boardName = "" // eg. "Gossiping"
	if req.BoardRef.GetName() != "" {
		boardName = req.BoardRef.GetName()
	} else if req.BoardRef.GetBid() > 0 {
		if int(req.BoardRef.GetBid()) > len(cachedBoard) {
			slog.Error("boardd::Content", "invalid bid", req.BoardRef.GetBid(), "len(cachedBoard)", len(cachedBoard))
			return nil, fmt.Errorf("invalid bid: %d", req.BoardRef.GetBid())
		}
		boardName = cachedBoard[req.BoardRef.GetBid()].BoardID()
	} else {
		slog.Error("boardd::Content", "invalid boardref", req.BoardRef)
		return nil, fmt.Errorf("invalid boardref: %v", req.BoardRef)
	}
	slog.Info("boardd::Content", "boardName", boardName, "filename", req.Filename)

	// TODO: use connection pool
	globalConfig, err := config.NewDefaultConfig()
	if err != nil {
		slog.Error("failed to get config", "error", err)
	}
	db, err := sql.Open("postgres", globalConfig.BBSHome)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Check if the connection is alive
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping database", "error", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	defer db.Close()
	CreateArticleContentTable(db)

	content, err := ReadArticleContentFromDB(db, boardName, req.Filename)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("GetBoardArticleContent error", "error", err)
		return nil, err
	}
	if content != nil {
		slog.Info("boardd::Content", "found content in DB", "length", len(content))
		return &apipb.ContentReply{
			Content: &apipb.Content{
				Content: content,
			},
		}, nil
	}

	slog.Info("boardd::Content fetching content from webpttparser", "boardName", boardName, "filename", req.Filename)

	l := fmt.Sprintf("/%s/%s.html", boardName, req.Filename)
	b, err := webpttparser.GetPttPage(l)
	if err != nil {
		return nil, err
	}
	// store content to DB
	err = WriteArticleContentToDB(db, boardName, req.Filename, webpttparser.HandlePage(b))
	if err != nil {
		slog.Error("StoreBoardArticleContent error", "error", err)
		return nil, err
	}
	slog.Info("boardd::Content stored content to DB", "boardName", boardName, "filename", req.Filename)

	return &apipb.ContentReply{
		Content: &apipb.Content{
			Content: webpttparser.HandlePage(b),
		},
	}, nil

}

func (s *server) Hotboard(context.Context, *apipb.HotboardRequest) (*apipb.HotboardReply, error) {
	slog.Info("GetPopularBoards")
	records, err := s.usecase.GetPopularBoards(context.Background())
	if err != nil {
		slog.Error("GetPopularBoards error", "err", err)
		return nil, err
	}
	slog.Info("GetPopularBoards", "records", records)
	boards := make([]*apipb.Board, len(records))
	for i, record := range records {
		slog.Info("GetPopularBoards", "record", record)
		boards[i] = &apipb.Board{
			Bid:        0,
			Name:       record.BoardID(),
			Title:      record.Title(),
			NumUsers:   0,
			Bclass:     "",
			Attributes: 0,
		}
	}
	slog.Info("GetPopularBoards", "boards", boards)
	return &apipb.HotboardReply{
		Boards: boards,
	}, nil

	// return &apipb.HotboardReply{
	// 	Boards: []*apipb.Board{
	// 		{
	// 			Bid:        1,
	// 			Name:       "Gossiping",
	// 			Title:      "◎[老八] 共匪共諜就在本能寺",
	// 			NumUsers:   8714,
	// 			Bclass:     "綜合",
	// 			Attributes: 0,
	// 		},
	// 		{
	// 			Name:       "Stock",
	// 			Title:      "◎[股票] 漲停板",
	// 			NumUsers:   6031,
	// 			Bclass:     "學術",
	// 			Attributes: 0,
	// 		},
	// 		{
	// 			Name:       "C_Chat",
	// 			Title:      "◎[希洽] 發文時標題請防雷",
	// 			NumUsers:   3792,
	// 			Bclass:     "閒談",
	// 			Attributes: 0,
	// 		},
	// 		{
	// 			Name:       "Tech_Job",
	// 			Title:      "◎[科技] 這裡是科技板",
	// 			NumUsers:   378,
	// 			Bclass:     "工作",
	// 			Attributes: 0,
	// 		},
	// 		{
	// 			Name:       "TY_Research",
	// 			Title:      "◎強極渦劇場上映中",
	// 			NumUsers:   27,
	// 			Bclass:     "大氣",
	// 			Attributes: 0,
	// 		},
	// 	},
	// }, nil
}

// mock a bid table for BoardID (text) to Bid (int64) mapping

var isCacheBoardsInitialized = false
var cachedBoard = []bbs.BoardRecord{
	// mock data
	&bbs.UnimplementedBoardRecord{},
}

var boardToBoardIndex = map[string]uint32{
	"": 0,
}

func initCacheBoards(usecase usecase.Usecase) {
	if isCacheBoardsInitialized {
		return
	}
	boards := usecase.GetBoards(context.Background(), "")
	for _, board := range boards {
		boardIndex, ok := boardToBoardIndex[strings.ToLower(board.BoardID())]
		if !ok {
			boardIndex = uint32(len(boardToBoardIndex))
			cachedBoard = append(cachedBoard, board)
			boardToBoardIndex[strings.ToLower(board.BoardID())] = boardIndex
			slog.Info("initCacheBoards", "new boardIndex", boardIndex, "boardID", board.BoardID())
		}
	}
	isCacheBoardsInitialized = true
}
