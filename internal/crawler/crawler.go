package crawler

import "github.com/JinHyeokOh01/goodjob-notifier/internal/models"

type KHUCrawler struct {
	TargetURL string
}

func NewKHUCrawler(url string) *KHUCrawler {
	return &KHUCrawler{TargetURL: url}
}

func (c *KHUCrawler) GetLatestPost() (*models.Post, error) {
	// 임시 데이터 반환 (나중에 실제 로직 구현)
	return &models.Post{ID: 1, Title: "Skeleton"}, nil
}
