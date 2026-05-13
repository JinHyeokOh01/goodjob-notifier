package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/crawler"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/notifier"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/repository"
	"github.com/JinHyeokOh01/goodjob-notifier/internal/models"
)

func main() {
	// 1. .env 파일에서 비밀번호와 설정값 불러오기
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	receiver := os.Getenv("EMAIL_RECEIVER")
	targetURL := "https://goodjob.khu.ac.kr/bbs/board.php?bo_table=s4_1&sca=%EA%B5%AD%EC%A0%9C&page=1"

	// 2. 우리가 만든 3가지 핵심 부품 생성 (인스턴스화)
	c := crawler.NewKHUCrawler(targetURL)
	repo := repository.NewFileRepository("last_post.txt")
	noti := notifier.NewEmailNotifier(sender, password, receiver)

	// 3. 크롤러 작동: 최신 게시글 목록 싹 다 가져오기
	fmt.Println("🔍 크롤링을 시작합니다...")
	posts, err := c.GetRecentPosts()
	if err != nil {
		log.Fatal(err)
	}

	// 4. 기억력 작동: 파일에서 마지막으로 본 글의 링크 가져오기
	lastLink, err := repo.GetLastPostLink()
	if err != nil {
		log.Fatal("저장소 읽기 오류:", err)
	}

	// [중요 로직] 최초 실행 처리
	if lastLink == "" {
		fmt.Println("💡 최초 실행입니다! 알림 봇 테스트를 위해 가장 최신 글 1개를 이메일로 전송합니다.")
		err = noti.SendAlert(posts[0])
		if err != nil {
			log.Fatal("이메일 발송 실패:", err)
		}
		repo.SaveLastPostLink(posts[0].Link) // 본 글 저장
		fmt.Println("✅ 테스트 이메일 발송 및 초기화 완료!")
		return
	}

	// 5. 새 글 판별 및 알림 발송 (여러 개의 새 글이 올라왔을 때를 대비)
	var newPosts []*models.Post
	for _, p := range posts {
		if p.Link == lastLink {
			break // 어제 본 글을 만나면 탐색 중지!
		}
		newPosts = append(newPosts, p) // 새 글이면 바구니에 담기
	}

	if len(newPosts) == 0 {
		fmt.Println("📭 새로운 공지가 없습니다. 봇을 종료합니다.")
	} else {
		fmt.Printf("🚨 총 %d개의 새로운 공지를 발견했습니다! 알림을 발송합니다.\n", len(newPosts))
		
		// 바구니에 담긴 새 글들에 대해 각각 이메일 쏘기
		for _, p := range newPosts {
			err = noti.SendAlert(p)
			if err != nil {
				fmt.Printf("이메일 발송 실패 (%s): %v\n", p.Title, err)
			} else {
				fmt.Printf("📧 이메일 발송 성공: %s\n", p.Title)
			}
		}

		// 파일 업데이트: 가장 최신 글(0번째 인덱스)의 링크로 덮어쓰기
		repo.SaveLastPostLink(posts[0].Link)
		fmt.Println("💾 저장소 업데이트 완료!")
	}
}