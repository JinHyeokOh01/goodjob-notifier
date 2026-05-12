package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/JinHyeokOh01/goodjob-notifier/internal/models"
	"github.com/PuerkitoBio/goquery"
)

type KHUCrawler struct {
	TargetURL string
}

func NewKHUCrawler(url string) *KHUCrawler {
	return &KHUCrawler{TargetURL: url}
}

// GetRecentPosts 로 이름 변경 및 반환 타입을 슬라이스([]*models.Post)로 변경
func (c *KHUCrawler) GetRecentPosts() ([]*models.Post, error) {
	res, err := http.Get(c.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("네트워크 오류 발생: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("정상적인 응답이 아닙니다. 상태 코드: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("HTML 파싱 오류: %v", err)
	}

	var posts []*models.Post // 여러 글을 담을 빈 바구니(슬라이스) 준비

	// .div_tb_tr 중에서 안에 .div_td(실제 데이터)가 있는 행들만 싹 다 찾습니다. (헤더 자동 제외)
	doc.Find(".div_tb_tr").Has(".div_td").Each(func(i int, s *goquery.Selection) {
		titleNode := s.Find(".col_subject a")
		dateNode := s.Find(".col_date")

		link, exists := titleNode.Attr("href")
		if !exists || link == "" {
			return // 링크가 없으면 건너뜀 (Each 반복문 안에서의 continue 역할)
		}

		cleanTitle := strings.Join(strings.Fields(titleNode.Text()), " ")
		cleanDate := strings.TrimSpace(dateNode.Text())

		// 게시글 구조체 생성
		post := &models.Post{
			Title: cleanTitle,
			Date:  cleanDate,
			Link:  link,
		}

		// 바구니에 담기
		posts = append(posts, post)
	})

	if len(posts) == 0 {
		return nil, fmt.Errorf("게시글을 찾지 못했습니다. 선택자를 확인해주세요.")
	}

	return posts, nil
}
