package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DB에 저장할 데이터 형태 (구조체)
type Record struct {
	ID   string `dynamodbav:"id"`
	Link string `dynamodbav:"link"`
}

// DynamoRepository는 AWS DynamoDB와 통신합니다.
type DynamoRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoRepository: AWS SDK를 이용해 통신 준비를 합니다.
func NewDynamoRepository(tableName string) *DynamoRepository {
	// 1. .env에 있는 AWS 자격 증명(열쇠)을 자동으로 읽어옵니다.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("AWS 설정 로드 실패: %v", err)
	}

	// 2. DynamoDB 전용 클라이언트 생성
	client := dynamodb.NewFromConfig(cfg)

	return &DynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

// GetLastPostLink: DynamoDB에서 마지막 게시글 링크를 읽어옵니다.
func (r *DynamoRepository) GetLastPostLink() (string, error) {
	// 우리가 찾을 데이터의 이름표(Key)는 "last_post" 입니다.
	key, err := attributevalue.MarshalMap(map[string]string{"id": "last_post"})
	if err != nil {
		return "", err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key,
	}

	// AWS에 데이터 요청!
	result, err := r.client.GetItem(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("DynamoDB 읽기 오류: %v", err)
	}

	// 테이블이 비어있는 경우 (최초 실행)
	if result.Item == nil {
		return "", nil
	}

	// 가져온 데이터를 Record 구조체로 예쁘게 변환
	var record Record
	err = attributevalue.UnmarshalMap(result.Item, &record)
	if err != nil {
		return "", fmt.Errorf("데이터 파싱 오류: %v", err)
	}

	return record.Link, nil
}

// SaveLastPostLink: 새로운 게시글의 링크를 DynamoDB에 저장(덮어쓰기)합니다.
func (r *DynamoRepository) SaveLastPostLink(link string) error {
	record := Record{
		ID:   "last_post", // 항상 같은 ID에 덮어써서 최신화합니다.
		Link: link,
	}

	item, err := attributevalue.MarshalMap(record)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	// AWS에 데이터 저장!
	_, err = r.client.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("DynamoDB 쓰기 오류: %v", err)
	}

	return nil
}