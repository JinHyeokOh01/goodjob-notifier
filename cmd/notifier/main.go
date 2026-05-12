package main

import (
	"fmt"
	"log"

	"github.com/JinHyeokOh01/goodjob-notifier/internal/crawler"
)

func main() {
	url := "https://goodjob.khu.ac.kr/bbs/board.php?bo_table=s4_1&sca=%EA%B5%AD%EC%A0%9C&page=1"

	c := crawler.NewKHUCrawler(url)

	posts, err := c.GetRecentPosts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🎉 크롤링 성공! 총 %d개의 글을 가져왔습니다.\n\n", len(posts))

	for i := 0; i < len(posts); i++ {
		fmt.Printf("[%d번 글]\n", i+1)
		fmt.Println("제목:", posts[i].Title)
		fmt.Println("등록일:", posts[i].Date)
		fmt.Println("링크:", posts[i].Link)
		fmt.Println("------------------------")
	}
}
