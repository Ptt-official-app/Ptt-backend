package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ptt-official-app/go-bbs"

	"github.com/Ptt-official-app/Ptt-backend/internal/repository"
)

// GetPopularArticles returns articles by descending comment_count
func (usecase *usecase) GetPopularArticles(ctx context.Context) ([]repository.PopularArticleRecord, error) {
	articles, err := usecase.repo.GetPopularArticles(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPopularArticles error: %w", err)
	}
	return articles, nil
}

func isANSIStyleToken(token string) bool {
	token = strings.TrimPrefix(token, "\x1b")
	if len(token) < 2 || token[0] != '[' || token[len(token)-1] != 'm' {
		return false
	}
	for i := 1; i < len(token)-1; i++ {
		if (token[i] < '0' || token[i] > '9') && token[i] != ';' {
			return false
		}
	}
	return true
}

func stripANSISequences(line string) string {
	var result strings.Builder
	result.Grow(len(line))

	for i := 0; i < len(line); {
		if line[i] == '\x1b' && i+1 < len(line) && line[i+1] == '[' {
			j := i + 2
			for j < len(line) && ((line[j] >= '0' && line[j] <= '9') || line[j] == ';') {
				j++
			}
			if j < len(line) && line[j] == 'm' {
				i = j + 1
				continue
			}
		}
		result.WriteByte(line[i])
		i++
	}

	return result.String()
}

func parseUsefulnessComment(line string) (appendType, userID string, ok bool) {
	fields := strings.Fields(stripANSISequences(line))
	first := 0
	for first < len(fields) && isANSIStyleToken(fields[first]) {
		first++
	}
	if len(fields)-first < 2 {
		return "", "", false
	}
	if fields[first] != "↑" && fields[first] != "↓" {
		return "", "", false
	}
	return fields[first], fields[first+1], true
}

func userUsefulnessScore(article, userID string) int {
	score := 0
	for _, line := range strings.Split(article, "\n") {
		appendType, commentUserID, ok := parseUsefulnessComment(line)
		if !ok || !strings.EqualFold(commentUserID, userID) {
			continue
		}

		switch appendType {
		case "↑":
			if score < 1 {
				score++
			}
		case "↓":
			if score > -1 {
				score--
			}
		}
	}
	return score
}

func (usecase *usecase) UpdateUsefulness(ctx context.Context, userID, boardID, filename, appendType string) (repository.PushRecord, error) {
	if appendType != "↑" && appendType != "↓" {
		return nil, fmt.Errorf("UpdateUsefulness error: unsupported usefulness type %q", appendType)
	}

	articleRecords, err := usecase.repo.GetBoardArticleRecords(ctx, boardID, 0, ^uint(0))
	if err != nil {
		return nil, fmt.Errorf("UpdateUsefulness error: %w", err)
	}

	var owner string
	found := false
	for _, record := range articleRecords {
		if record.Filename() == filename {
			owner = record.Owner()
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("UpdateUsefulness error: article %s not found", filename)
	}

	if strings.EqualFold(owner, userID) {
		return nil, fmt.Errorf("UpdateUsefulness error: owners cannot rate their own article")
	}

	article, err := usecase.GetBoardArticle(ctx, boardID, filename)
	if err != nil {
		return nil, fmt.Errorf("UpdateUsefulness error: %w", err)
	}

	currentScore := userUsefulnessScore(bbs.Big5ToUtf8(article), userID)
	if (appendType == "↑" && currentScore >= 1) || (appendType == "↓" && currentScore <= -1) {
		return nil, fmt.Errorf("UpdateUsefulness error: user %s already reached usefulness limit %d", userID, currentScore)
	}

	p, err := usecase.repo.AppendComment(ctx, userID, boardID, filename, appendType, "")
	if err != nil {
		return nil, fmt.Errorf("UpdateUsefulness error: %w", err)
	}
	return p, nil
}

// AppendComment append comment to specific article
func (usecase *usecase) AppendComment(ctx context.Context, userID, boardID, filename, appendType, text string) (repository.PushRecord, error) {
	result, err := usecase.repo.AppendComment(ctx, userID, boardID, filename, appendType, text)

	return result, err
}

// ForwardArticleToBoard returns forwarding to board results
func (usecase *usecase) ForwardArticleToBoard(ctx context.Context, userID, boardID, filename, boardName string) (repository.ForwardArticleToBoardRecord, error) {
	forwardArticle, err := usecase.repo.ForwardArticleToBoard(ctx, userID, boardID, filename, boardName)
	if err != nil {
		return nil, fmt.Errorf("ForwardArticleToBoard error: %w", err)
	}

	return forwardArticle, err
}

// ForwardArticleToEmail returns forwarding to email results
func (usecase *usecase) ForwardArticleToEmail(ctx context.Context, userID, boardID, filename, email string) error {
	articleRecords, err := usecase.repo.GetBoardArticleRecords(ctx, boardID, 0, ^uint(0))
	if err != nil {
		return fmt.Errorf("GetBoardArticleRecords error: %w", err)
	}
	var title string
	for _, article := range articleRecords {
		if article.Filename() == filename {
			title = article.Title()
			break
		}
	}
	if title == "" {
		return fmt.Errorf("cannot find article %s", filename)
	}
	buffer, err := usecase.repo.GetBoardArticle(ctx, boardID, filename)
	if err != nil {
		return fmt.Errorf("GetBoardArticle error: %w", err)
	}
	return usecase.mailProvider.Send(email, title, userID, buffer)
}

// CreateArticle create a new article on a board
func (usecase *usecase) CreateArticle(ctx context.Context, userID, boardID, title, article string) (bbs.ArticleRecord, error) {
	record, err := usecase.repo.CreateArticle(ctx, userID, boardID, title, article)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (usecase *usecase) GetArticleURL(boardID string, filename string) string {
	// TODO: generate article url by config file
	return fmt.Sprintf("https://pttapp.cc/bbs/%s/%s.html", boardID, filename)
}

func (usecase *usecase) GetRawArticle(boardID, filename string) (string, error) {
	raw, err := usecase.repo.GetRawArticle(boardID, filename)
	if err != nil {
		return "", err
	}

	return raw, nil
}
