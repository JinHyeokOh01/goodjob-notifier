package main

import (
	"log"

	"github.com/JinHyeokOh01/goodjob-notifier/internal/crawler"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/handler"
)

func main() {
	// 1. 의존성 생성 (Notifier는 아직 안 만들었으니 나중에 추가)
	c := crawler.NewKHUCrawler("https://example.com")

	// 2. 핸들러 조립 (지금은 Notifier 자리에 nil을 넣거나 더미를 넣음)
	h := handler.NewAppHandler(c, nil)

	// 3. 실행
	if err := h.Run(); err != nil {
		log.Fatal(err)
	}
}
