package notifier

import (
	"fmt"
	"net/smtp"

	"github.com/JinHyeokOh01/goodjob-notifier/internal/models"
)

// EmailNotifier는 구글 SMTP를 이용해 이메일을 발송합니다.
type EmailNotifier struct {
	SMTPServer string
	SMTPPort   string
	Sender     string
	Password   string
	Receiver   string
}

// NewEmailNotifier 생성자 함수
func NewEmailNotifier(sender, password, receiver string) *EmailNotifier {
	return &EmailNotifier{
		SMTPServer: "smtp.gmail.com", // 구글 메일 서버 주소
		SMTPPort:   "587",            // 구글 메일 포트
		Sender:     sender,
		Password:   password,
		Receiver:   receiver,
	}
}

// SendAlert: 새 글 정보를 받아 예쁜 HTML 이메일로 전송합니다.
func (e *EmailNotifier) SendAlert(post *models.Post) error {
	// SMTP 인증 정보 세팅
	auth := smtp.PlainAuth("", e.Sender, e.Password, e.SMTPServer)

	// 이메일 제목 및 헤더 (한글 깨짐 방지를 위해 UTF-8 설정)
	subject := fmt.Sprintf("Subject: 🚨 [미래인재센터의 새 글이 업로드 되었습니다!] %s\n", post.Title)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	
	// 이메일 본문 (HTML로 예쁘게 꾸미기)
	body := fmt.Sprintf(`
		<h2>새 게시글 : 경희대학교 미래인재센터</h2>
		<ul>
			<li><b>제목:</b> %s</li>
			<li><b>등록일:</b> %s</li>
		</ul>
		<p><a href="%s" style="padding: 10px 15px; background-color: #004d99; color: white; text-decoration: none; border-radius: 5px;">공지사항 바로가기</a></p>
	`, post.Title, post.Date, post.Link)

	// 메세지 조립
	msg := []byte(subject + mime + body)
	addr := e.SMTPServer + ":" + e.SMTPPort
	
	// 이메일 쏘기!
	err := smtp.SendMail(addr, auth, e.Sender, []string{e.Receiver}, msg)
	if err != nil {
		return fmt.Errorf("이메일 발송 실패: %v", err)
	}

	return nil
}