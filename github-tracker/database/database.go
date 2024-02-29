package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github-tracker/github-tracker/models"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	_ "github.com/lib/pq"
)

var (
	getDbOnce sync.Once
	db        *sql.DB
)

func connect(ctx context.Context) (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPasswordSecretId := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	config, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(models.REGION))
	if err != nil {
		return &sql.DB{}, err
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(dbPasswordSecretId),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return &sql.DB{}, err
	}

	if result.SecretString == nil {
		return &sql.DB{}, err
	}

	passwordJson := *result.SecretString

	var secret map[string]string
	err = json.Unmarshal([]byte(passwordJson), &secret)
	if err != nil {
		return &sql.DB{}, err
	}

	password := secret["password"]

	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=require search_path=public`,
		dbHost,
		dbPort,
		dbUser,
		password,
		dbName,
	)

	return sql.Open("postgres", dsn)
}

func Connect(ctx context.Context) (*sql.DB, error) {
	var err error

	getDbOnce.Do(func() {
		db, err = connect(ctx)
	})

	return db, err
}
