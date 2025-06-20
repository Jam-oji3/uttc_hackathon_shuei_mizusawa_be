package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hackathon/model"
	"hackathon/repository"
	"hackathon/util"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type PostCreateUseCase struct {
	TxExecutor   repository.TransactionExecutor
	GeminiClient repository.GeminiClient
	PostRepo     repository.PostsRepository
	TrendRepo    repository.TrendsRepository
	SpoilerRepo  repository.SpoilersRepository
	DB           *sql.DB
}

func NewPostCreateUseCase(
	txExecutor repository.TransactionExecutor,
	geminiClient repository.GeminiClient,
	postRepo repository.PostsRepository,
	trendRepo repository.TrendsRepository,
	spoilerRepo repository.SpoilersRepository,
	db *sql.DB) *PostCreateUseCase {
	return &PostCreateUseCase{
		TxExecutor:   txExecutor,
		GeminiClient: geminiClient,
		PostRepo:     postRepo,
		TrendRepo:    trendRepo,
		SpoilerRepo:  spoilerRepo,
		DB:           db,
	}
}

func (uc *PostCreateUseCase) Execute(ctx context.Context, userId string, text string, replyTo *string, repostRef *string, mediaType *string, mediaUrl *string) (*model.Post, error) {
	id := util.GenerateULID()
	now := time.Now()

	post := model.Post{
		Id:        id,
		UserId:    userId,
		Text:      text,
		ReplyTo:   replyTo,
		RepostRef: repostRef,
		MediaType: mediaType,
		MediaURL:  mediaUrl,
		CreatedAt: now,
	}

	// バリデーション
	if err := validatePostData(&post); err != nil {
		return nil, err
	}

	generatedStr, err := uc.GeminiClient.GenerateContent(ctx, fmt.Sprintf(`
以下の投稿にアニメ、映画、ゲーム等、書籍等のネタバレが含まれる可能性があるか判定してください。
ネタバレが存在する可能性があるかないかを"SPOILER: true"のように明記してください。
作品名がわかる場合には"TITLE: ドラゴンボール"のように1つだけ作品名を書いてください。
作品名がわからない場合は"アニメ","映画"などのカテゴリー名か、"不明"をキーワードとしてください
また、トレンド集計に使うネタバレを含まないキーワードを"KEYWORDS: [A,B,C]"の形で最大3個出力してください。

【投稿】"%s"

【出力例】
SPOILER: true
TITLE: 推しの子 
KEYWORDS: [B小町,推しの子,第14話]
	`, text))

	if generatedStr == nil || strings.TrimSpace(*generatedStr) == "" {
		return nil, errors.New("gemini response is empty")
	}

	spoilerRe := regexp.MustCompile(`(?i)SPOILER:\s*(true|false)`)
	spoilerFlag := false
	if match := spoilerRe.FindStringSubmatch(*generatedStr); match != nil {
		spoilerFlag = strings.ToLower(match[1]) == "true"
	}

	var spoilerKeyword *string
	keywordRe := regexp.MustCompile(`(?i)TITLE:\s*(.+)`)
	if match := keywordRe.FindStringSubmatch(*generatedStr); match != nil {
		keyword := strings.TrimSpace(match[1])
		spoilerKeyword = &keyword
	}

	// KEYWORDS: の抽出
	keywordsRe := regexp.MustCompile(`(?i)KEYWORDS:\s*\[([^]]+)]`)
	var trendKeywords []string
	if match := keywordsRe.FindStringSubmatch(*generatedStr); match != nil {
		raw := match[1] // "B小町,推しの子,第14話"
		// カンマで区切って配列に
		for _, word := range strings.Split(raw, ",") {
			trimmed := strings.TrimSpace(word)
			if trimmed != "" {
				trendKeywords = append(trendKeywords, trimmed)
			}
		}
	}

	var spoiler model.Spoiler
	if spoilerFlag && spoilerKeyword != nil {
		spoiler = model.Spoiler{
			Id:      util.GenerateULID(),
			PostId:  post.Id,
			Keyword: *spoilerKeyword,
		}
	}

	unique := make(map[string]struct{})
	for _, noun := range trendKeywords {
		unique[noun] = struct{}{}
	}
	//時間単位で切り捨て
	hour := time.Now().Truncate(time.Hour)

	var trends []model.Trend
	for word := range unique {
		trends = append(trends, model.Trend{
			Id:    util.GenerateULID(), // IDはユニークなULID
			Word:  word,
			Count: 1,
			Hour:  hour,
		})
	}

	// トランザクション内で投稿作成
	_, err = uc.TxExecutor.DoInTx(ctx, uc.DB, func(ctx context.Context, tx *sql.Tx) (interface{}, error) {
		if err := uc.PostRepo.InsertPost(ctx, tx, post); err != nil {
			return nil, err
		}
		if err := uc.TrendRepo.InsertTrends(ctx, tx, trends); err != nil {
			return nil, err
		}
		if spoilerFlag && spoilerKeyword != nil {
			if err := uc.SpoilerRepo.InsertSpoiler(ctx, tx, &spoiler); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func validatePostData(post *model.Post) error {
	if post.Id == "" {
		return errors.New("post ID is required")
	}
	if post.UserId == "" {
		return errors.New("user ID is required")
	}
	if strings.TrimSpace(post.Text) == "" && post.MediaURL == nil {
		return errors.New("either text or media URL must be provided")
	}
	if len(post.Text) > 400 {
		return errors.New("text is too long (max 200 characters)")
	}

	if post.RepostRef != nil && post.ReplyTo != nil {
		return errors.New("repostRef and replyTo are mutually exclusive")
	}

	// MediaURLがあればURLとして妥当かチェック
	if post.MediaURL != nil {
		parsedURL, err := url.ParseRequestURI(*post.MediaURL)
		if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
			return errors.New("mediaURL must be a valid HTTP or HTTPS URL")
		}
	}

	// MediaTypeは空でもいいが、もしあれば簡単に制限
	if post.MediaType != nil {
		allowedMediaTypes := map[string]bool{
			"photo": true,
			"model": true,
		}
		if !allowedMediaTypes[*post.MediaType] {
			return errors.New("unsupported media type")
		}
	}

	return nil
}
