package repository

import (
	"os"
	"strings"
)

// FileRepository는 텍스트 파일을 이용해 데이터를 저장하고 불러옵니다.
type FileRepository struct {
	FilePath string
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{FilePath: path}
}

// GetLastPostLink: 파일에 적힌 마지막 게시글의 링크를 읽어옵니다.
func (r *FileRepository) GetLastPostLink() (string, error) {
	data, err := os.ReadFile(r.FilePath)
	if err != nil {
		// 파일이 아예 존재하지 않는 경우 (최초 실행 시) 에러 대신 빈 문자열 반환
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SaveLastPostLink: 새로운 게시글의 링크를 파일에 덮어씁니다.
func (r *FileRepository) SaveLastPostLink(link string) error {
	// 0644는 파일 권한(읽기/쓰기 가능)을 의미합니다.
	return os.WriteFile(r.FilePath, []byte(link), 0644)
}