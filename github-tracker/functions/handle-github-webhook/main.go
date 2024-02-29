package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github-tracker/github-tracker/database"
	"github-tracker/github-tracker/models"
	"github-tracker/github-tracker/repository"
	"github-tracker/github-tracker/repository/entity"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const SignaturePrefix = "sha256="

var ErrInvalidWebhook = errors.New("invalid github webhook")

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	isValid, err := validateGitHubRequest(ctx, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error validating GitHub webhook",
		}, err
	}

	if !isValid {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusBadRequest,
			Body:            "invalid GitHub webhook",
		}, ErrInvalidWebhook
	}

	var webhook models.GitHubWebhook
	err = json.Unmarshal([]byte(request.Body), &webhook)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)

		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error while unmarshal GitHub webhook",
		}, nil
	}

	fmt.Println("webhook: ", webhook.HeadCommit)

	db, err := database.Connect(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error initilizing database",
		}, err
	}

	commitRepo := repository.NewCommit(db)

	err = insertGitHubWebhook(ctx, commitRepo, webhook, request.Body, time.Now())
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error creating commit in db",
		}, err
	}

	fmt.Println("commit created!")

	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            "ok",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func insertGitHubWebhook(ctx context.Context, repo repository.Commit, webhook models.GitHubWebhook, body string, createdTime time.Time) error {
	commit := entity.Commit{
		RepoName:       webhook.Repository.FullName,
		CommitID:       webhook.HeadCommit.ID,
		CommitMessage:  webhook.HeadCommit.Message,
		AuthorUsername: webhook.HeadCommit.Author.Username,
		AuthorEmail:    webhook.HeadCommit.Author.Email,
		Payload:        body,
		CreatedAt:      createdTime,
		UpdatedAt:      createdTime,
	}

	err := repo.Insert(ctx, &commit)

	return err
}

func validateGitHubRequest(ctx context.Context, request events.APIGatewayProxyRequest) (bool, error) {
	fullSignature, ok := request.Headers["x-hub-signature-256"]
	if !ok {
		return false, fmt.Errorf("missing x-hub-signature")
	}

	signature := strings.TrimPrefix(fullSignature, SignaturePrefix)

	secretARN := os.Getenv("GITHUB_SECRET")

	config, err := config.LoadDefaultConfig(ctx, config.WithRegion(models.REGION))
	if err != nil {
		return false, fmt.Errorf("error loading config: %s", err.Error())
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretARN),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return false, fmt.Errorf("error getting secret value: %s", err.Error())
	}

	if result.SecretString == nil {
		return false, fmt.Errorf("secret value is nil")
	}

	passwordJson := *result.SecretString

	var secret map[string]string
	err = json.Unmarshal([]byte(passwordJson), &secret)
	if err != nil {
		return false, fmt.Errorf("error unmarshaling secret value: %s", err.Error())
	}

	gitHubSecret := secret["secret"]

	newSha := calculateSignature(gitHubSecret, request.Body)

	return hmac.Equal([]byte(signature), []byte(newSha)), nil
}

func calculateSignature(secret string, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))

	return hex.EncodeToString(mac.Sum(nil))
}

func main() {
	lambda.Start(handleRequest)
}
