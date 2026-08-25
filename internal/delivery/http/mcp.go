package http

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxMCPArticleBytes = 1 << 20

var (
	mcpBoardIDPattern  = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)
	mcpFilenamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,128}$`)
)

type mcpBoardInput struct {
	BoardID string `json:"board_id" jsonschema:"The PTT board ID, such as Gossiping."`
}

type mcpArticleInput struct {
	BoardID  string `json:"board_id" jsonschema:"The PTT board ID, such as Gossiping."`
	Filename string `json:"filename" jsonschema:"The PTT article filename."`
}

type mcpBoardOutput struct {
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	Posts   uint32 `json:"posts"`
}

type mcpPostOutput struct {
	Filename   string `json:"filename"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	Date       string `json:"date"`
	Recommends int    `json:"recommends"`
}

type mcpBoardPageOutput struct {
	BoardID string          `json:"board_id"`
	Posts   []mcpPostOutput `json:"posts"`
}

type mcpArticleOutput struct {
	BoardID  string `json:"board_id"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func (delivery *Delivery) mcpHandler() http.Handler {
	server := mcp.NewServer(&mcp.Implementation{
		Name:        "pttapp",
		Title:       "PTT App",
		Description: "Read-only access to public PTT board and article data.",
		Version:     "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "fetch_board_info",
		Description: "Get public metadata and article count for a PTT board.",
	}, delivery.fetchMCPBoardInfo)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fetch_board_page",
		Description: "Get the most recent 20 public articles from a PTT board.",
	}, delivery.fetchMCPBoardPage)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "fetch_article",
		Description: "Get a public PTT article by board and filename.",
	}, delivery.fetchMCPArticle)

	return mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{
			Stateless:                    true,
			JSONResponse:                 true,
			MaxRequestBodyBytes:          maxMCPArticleBytes,
			PropagateRequestCancellation: true,
		},
	)
}

func validateMCPBoardID(boardID string) error {
	if !mcpBoardIDPattern.MatchString(boardID) {
		return fmt.Errorf("invalid board_id")
	}
	return nil
}

func validateMCPFilename(filename string) error {
	if !mcpFilenamePattern.MatchString(filename) {
		return fmt.Errorf("invalid filename")
	}
	return nil
}

func (delivery *Delivery) fetchMCPBoardInfo(ctx context.Context, _ *mcp.CallToolRequest, input mcpBoardInput) (*mcp.CallToolResult, mcpBoardOutput, error) {
	if err := validateMCPBoardID(input.BoardID); err != nil {
		return nil, mcpBoardOutput{}, err
	}
	board, err := delivery.usecase.GetBoardByName(ctx, input.BoardID)
	if err != nil {
		return nil, mcpBoardOutput{}, fmt.Errorf("board not found")
	}
	posts := delivery.usecase.GetBoardArticles(ctx, board.BoardID(), 0, ^uint(0), &usecase.ArticleSearchCond{})
	return nil, mcpBoardOutput{BoardID: board.BoardID(), Title: board.Title(), Posts: uint32(len(posts))}, nil
}

func (delivery *Delivery) fetchMCPBoardPage(ctx context.Context, _ *mcp.CallToolRequest, input mcpBoardInput) (*mcp.CallToolResult, mcpBoardPageOutput, error) {
	if err := validateMCPBoardID(input.BoardID); err != nil {
		return nil, mcpBoardPageOutput{}, err
	}
	board, err := delivery.usecase.GetBoardByName(ctx, input.BoardID)
	if err != nil {
		return nil, mcpBoardPageOutput{}, fmt.Errorf("board not found")
	}
	articles := delivery.usecase.GetBoardArticles(ctx, board.BoardID(), 0, 20, &usecase.ArticleSearchCond{})
	posts := make([]mcpPostOutput, 0, len(articles))
	for _, article := range articles {
		posts = append(posts, mcpPostFromRecord(article))
	}
	return nil, mcpBoardPageOutput{BoardID: board.BoardID(), Posts: posts}, nil
}

func (delivery *Delivery) fetchMCPArticle(ctx context.Context, _ *mcp.CallToolRequest, input mcpArticleInput) (*mcp.CallToolResult, mcpArticleOutput, error) {
	if err := validateMCPBoardID(input.BoardID); err != nil {
		return nil, mcpArticleOutput{}, err
	}
	if err := validateMCPFilename(input.Filename); err != nil {
		return nil, mcpArticleOutput{}, err
	}
	content, err := delivery.usecase.GetBoardArticle(ctx, input.BoardID, input.Filename)
	if err != nil {
		return nil, mcpArticleOutput{}, fmt.Errorf("article not found")
	}
	if len(content) > maxMCPArticleBytes {
		return nil, mcpArticleOutput{}, fmt.Errorf("article exceeds the maximum response size")
	}
	return nil, mcpArticleOutput{BoardID: input.BoardID, Filename: input.Filename, Content: string(content)}, nil
}

func mcpPostFromRecord(article bbs.ArticleRecord) mcpPostOutput {
	return mcpPostOutput{
		Filename:   article.Filename(),
		Title:      article.Title(),
		Author:     article.Owner(),
		Date:       article.Date(),
		Recommends: article.Recommend(),
	}
}
