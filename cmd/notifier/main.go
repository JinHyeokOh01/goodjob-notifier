package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/crawler"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/notifier"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/repository"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/models"
)

func main() {
	// 1. 환경변수 및 부품 세팅
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	receiver := os.Getenv("EMAIL_RECEIVER")
	targetURL := "https://goodjob.khu.ac.kr/bbs/board.php?bo_table=s4_1&sca=%EA%B5%AD%EC%A0%9C&page=1"
	tableName := os.Getenv("DYNAMO_TABLE_NAME")

	c := crawler.NewKHUCrawler(targetURL)
	//repo := repository.NewFileRepository("last_post.txt")
	repo := repository.NewDynamoRepository(tableName)
	noti := notifier.NewEmailNotifier(sender, password, receiver)

	// 2. 프로그램 켜지자마자 1번 즉시 실행
	fmt.Println("🤖 알림 봇 가동 시작!")
	runBot(c, repo, noti)

	// 3. 스케줄러(Cron) 설정
	cr := cron.New()
	
	// 매일 오전 11시 0분, 오후 6시(18시) 0분에 실행
	_, err = cr.AddFunc("0 11,18 * * *", func() {
		fmt.Println("\n⏰ [정기 스캔] 정해진 시간(11:00 / 18:00)이 되어 웹사이트를 확인합니다...")
		runBot(c, repo, noti)
	})
	if err != nil {
		log.Fatal("스케줄러 등록 실패:", err)
	}

	cr.Start()
	fmt.Println("🚀 봇이 오전 11시/오후 6시 감시 모드에 들어갔습니다. (종료: Ctrl + C)")

	// 4. 메인 함수가 종료되지 않고 계속 실행되도록 막아두기
	// 운영체제(OS)에서 강제 종료 신호(Ctrl+C)가 올 때까지 대기합니다.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("\n🛑 알림 봇을 안전하게 종료합니다.")
}

// runBot: 실제 크롤링 및 알림 발송을 담당하는 핵심 로직 (기존 main 함수 내용)
//func runBot(c *crawler.KHUCrawler, repo *repository.FileRepository, noti *notifier.EmailNotifier) {
func runBot(c *crawler.KHUCrawler, repo *repository.DynamoRepository, noti *notifier.EmailNotifier) {
	posts, err := c.GetRecentPosts()
	if err != nil {
		fmt.Println("크롤링 오류:", err)
		return
	}

	lastLink, err := repo.GetLastPostLink()
	if err != nil {
		fmt.Println("저장소 읽기 오류:", err)
		return
	}

	if lastLink == "" {
		fmt.Println("💡 최초 실행이므로 최신 글을 저장합니다.")
		repo.SaveLastPostLink(posts[0].Link)
		noti.SendAlert(posts[0])
		return
	}

	var newPosts []interface{} // 새 글을 담을 바구니
	for _, p := range posts {
		if p.Link == lastLink {
			break
		}
		newPosts = append(newPosts, p)
	}

	if len(newPosts) == 0 {
		fmt.Println("📭 새로운 공지가 없습니다.")
	} else {
		fmt.Printf("🚨 %d개의 새로운 공지 발견! 알림을 발송합니다.\n", len(newPosts))
		for _, p := range newPosts {
			post := p.(*models.Post) // 타입 단언
			noti.SendAlert(post)
		}
		repo.SaveLastPostLink(posts[0].Link)
		fmt.Println("💾 마지막 본 글 저장 완료!")
	}
}