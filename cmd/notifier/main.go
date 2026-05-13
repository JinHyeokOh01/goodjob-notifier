package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/crawler"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/models"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/notifier"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/repository"
)

// HandleRequest: 람다가 깨어날 때마다 딱 한 번 실행되는 메인 로직입니다.
func HandleRequest(ctx context.Context) (string, error) {
	fmt.Println("🤖 람다(Lambda) 알림 봇 가동 시작!")

	// 1. 환경변수 불러오기 (godotenv.Load()가 빠졌습니다!)
	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	receiver := os.Getenv("EMAIL_RECEIVER")
	tableName := os.Getenv("DYNAMO_TABLE_NAME")
	targetURL := "https://goodjob.khu.ac.kr/bbs/board.php?bo_table=s4_1&sca=%EA%B5%AD%EC%A0%9C&page=1"

	// 2. 부품 조립
	c := crawler.NewKHUCrawler(targetURL)
	repo := repository.NewDynamoRepository(tableName)
	noti := notifier.NewEmailNotifier(sender, password, receiver)

	// 3. 크롤링 및 DB 확인
	posts, err := c.GetRecentPosts()
	if err != nil {
		return "", fmt.Errorf("크롤링 오류: %v", err)
	}

	lastLink, err := repo.GetLastPostLink()
	if err != nil {
		return "", fmt.Errorf("DB 읽기 오류: %v", err)
	}

	// 4. 최초 실행 처리
	if lastLink == "" {
		fmt.Println("💡 DB가 비어있습니다. 최신 글을 저장합니다.")
		repo.SaveLastPostLink(posts[0].Link)
		noti.SendAlert(posts[0])
		return "최초 실행 완료", nil
	}

	// 5. 새 글 판별 및 알림
	var newPosts []*models.Post
	for _, p := range posts {
		if p.Link == lastLink {
			break
		}
		newPosts = append(newPosts, p)
	}

	if len(newPosts) == 0 {
		msg := "📭 새로운 공지가 없습니다."
		fmt.Println(msg)
		return msg, nil
	}

	fmt.Printf("🚨 %d개의 새로운 공지 발견! 알림을 발송합니다.\n", len(newPosts))
	for _, p := range newPosts {
		err = noti.SendAlert(p)
		if err != nil {
			fmt.Printf("이메일 발송 실패: %v\n", err)
		}
	}
	
	// DB 업데이트
	repo.SaveLastPostLink(posts[0].Link)
	
	resultMsg := "💾 알림 발송 및 DB 업데이트 완료!"
	fmt.Println(resultMsg)
	return resultMsg, nil
}

func main() {
	// 기존의 무한 대기 루프와 cron이 전부 사라지고, 
	// 람다가 이 프로그램을 실행할 수 있도록 핸들러를 넘겨줍니다.
	lambda.Start(HandleRequest)
}