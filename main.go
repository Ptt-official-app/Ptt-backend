package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strconv"

	webpttparser "github.com/Ptt-official-app/Ptt-backend/webptt_parser"

	"github.com/Ptt-official-app/Ptt-backend/internal/aids"
	"github.com/Ptt-official-app/Ptt-backend/internal/config"
	dhttp "github.com/Ptt-official-app/Ptt-backend/internal/delivery/http"
	"github.com/Ptt-official-app/Ptt-backend/internal/logging"
	apipb "github.com/Ptt-official-app/Ptt-backend/internal/proto/api"
	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"google.golang.org/grpc"

	_ "github.com/PichuChen/postgresql-gobbs"
	"github.com/Ptt-official-app/go-bbs"
	_ "github.com/Ptt-official-app/go-bbs/pttbbs"
)

func main() {
	var logLevel = flag.Uint("logLevel", 4, `log level: 0: Emergency; 1: Alert; 2: Critical; 3: Error; 4: Warning; 5: Notice; 6: Info; 7: Debug`)
	flag.Usage = func() {
		os.Stderr.WriteString("Usage: \n  Ptt-backend [ options ]\n\n")
		os.Stderr.WriteString("Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if _, ok := os.LookupEnv("LOG_LEVEL"); !ok {
		if *logLevel > 7 {
			*logLevel = 7
		}
		os.Setenv("LOG_LEVEL", strconv.Itoa(int(*logLevel)))
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

type server struct {
	apipb.UnimplementedBoardServiceServer
	usecase usecase.Usecase
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

var boardToBoardIndex = map[string]uint32{
	"": 0,
}
var cachedBoard = []bbs.BoardRecord{
	// mock data
	&bbs.UnimplementedBoardRecord{},
}

func (s *server) Board(ctx context.Context, req *apipb.BoardRequest) (*apipb.BoardReply, error) {
	slog.Info("boardd::Board", "len(req.Ref)", len(req.Ref))
	boards := make([]*apipb.Board, len(req.Ref))
	for i, ref := range req.Ref {
		slog.Info("boardd::Board", "ref", ref.GetRef(), "bid", ref.GetBid(), "name", ref.GetName())
		board, err := s.usecase.GetBoardByID(context.Background(), ref.GetName())
		if err != nil {
			slog.Error("GetBoardByID error", "err", err)
			return nil, err
		}

		boardIndex, ok := boardToBoardIndex[board.BoardID()]
		if !ok {
			// if not found, assign a new index
			boardIndex = uint32(len(boardToBoardIndex))
			cachedBoard = append(cachedBoard, board)
			boardToBoardIndex[board.BoardID()] = boardIndex
			slog.Info("boardd::Board", "new boardIndex", boardIndex, "boardID", board.BoardID())
		}
		articles := s.usecase.GetBoardArticles(context.Background(), board.BoardID(), 0, ^uint(0), &usecase.ArticleSearchCond{})

		boards[i] = &apipb.Board{
			Bid:        boardIndex,
			Name:       board.BoardID(),
			Title:      board.Title(),
			NumUsers:   0,
			Bclass:     "",
			Attributes: 0,
			NumPosts:   uint32(len(articles)),
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
	return &apipb.BoardReply{Boards: boards}, nil
}

func (s *server) List(ctx context.Context, req *apipb.ListRequest) (*apipb.ListReply, error) {
	slog.Info("boardd::List", "ref", req.Ref, "bid", req.Ref.GetBid(), "req", req)
	var boardID = ""
	if len(req.Ref.GetName()) > 0 {
		boardID = req.Ref.GetName()
	} else if req.Ref.GetBid() > 0 {
		if int(req.Ref.GetBid()) > len(cachedBoard) {
			slog.Error("boardd::List", "invalid bid", req.Ref.GetBid(), "len(cachedBoard)", len(cachedBoard))
			return nil, fmt.Errorf("invalid bid: %d", req.Ref.GetBid())
		}
		boardID = cachedBoard[req.Ref.GetBid()].BoardID()
	} else {
		slog.Error("boardd::List", "invalid ref", req.Ref)
		return nil, fmt.Errorf("invalid ref: %v", req.Ref)
	}
	slog.Info("boardd::List", "boardID", boardID)
	var offset = uint(req.GetOffset())
	var length = uint(req.GetLength())
	if length <= 0 {
		length = 20 // default length
	}

	articles := s.usecase.GetBoardArticles(context.Background(), boardID, offset, length, &usecase.ArticleSearchCond{})
	if length > uint(len(articles)) {
		length = uint(len(articles)) // limit to the number of articles
	}
	if length > 20 {
		length = 20 // max length
	}
	posts := make([]*apipb.Post, length)
	for i, article := range articles[:length] {
		// slog.Info("boardd::List", "article", article)
		posts[i] = &apipb.Post{
			Index:         uint32(i + 1),
			Filename:      aids.Aidu2Fn(aids.Aidc2Aidu(article.Filename())),
			RawDate:       article.Date(),
			NumRecommends: int32(article.Recommend()),
			Owner:         article.Owner(),
			Title:         article.Title(),
		}
	}

	slog.Info("boardd::List", "len(posts)", len(posts))

	return &apipb.ListReply{
		Posts:   posts,
		Bottoms: []*apipb.Post{},
	}, nil

}

func (s *server) Content(ctx context.Context, req *apipb.ContentRequest) (*apipb.ContentReply, error) {
	slog.Info("boardd::Content", "boardref", req.BoardRef, "filename", req.Filename, "token", req.ConsistencyToken, "options", req.PartialOptions)
	var boardID = ""
	if req.BoardRef.GetName() != "" {
		boardID = req.BoardRef.GetName()
	} else if req.BoardRef.GetBid() > 0 {
		if int(req.BoardRef.GetBid()) > len(cachedBoard) {
			slog.Error("boardd::Content", "invalid bid", req.BoardRef.GetBid(), "len(cachedBoard)", len(cachedBoard))
			return nil, fmt.Errorf("invalid bid: %d", req.BoardRef.GetBid())
		}
		boardID = cachedBoard[req.BoardRef.GetBid()].BoardID()
	} else {
		slog.Error("boardd::Content", "invalid boardref", req.BoardRef)
		return nil, fmt.Errorf("invalid boardref: %v", req.BoardRef)
	}
	slog.Info("boardd::Content", "boardID", boardID, "filename", req.Filename)

	l := fmt.Sprintf("/%s/%s.html", boardID, req.Filename)
	b, err := getPttPage(l)
	if err != nil {
		return nil, err
	}
	return &apipb.ContentReply{
		Content: &apipb.Content{
			Content: webpttparser.HandlePage(b),
		},
	}, nil

}

func getPttPage(url string) ([]byte, error) {
	return webpttparser.GetPttPage(url)
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
