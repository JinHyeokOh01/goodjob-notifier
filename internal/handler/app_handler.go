package handler

import (
	"fmt"

	"github.com/JinHyeokOh01/goodjob-notifier/internal/models"
)

// Crawler 인터페이스 정의 (나중에 crawler 패키지가 이 규격을 맞춤)
type Crawler interface {
	GetLatestPost() (*models.Post, error)
}

// Notifier 인터페이스 정의
type Notifier interface {
	Send(post *models.Post) error
}

type AppHandler struct {
	crawler  Crawler
	notifier Notifier
}

func NewAppHandler(c Crawler, n Notifier) *AppHandler {
	return &AppHandler{crawler: c, notifier: n}
}

func (h *AppHandler) Run() error {
	fmt.Println("알림 시스템 가동...")
	// 여기에 대략적인 흐름만 작성
	return nil
}
