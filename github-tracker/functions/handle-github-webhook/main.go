package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github-tracker/github-tracker/database"
	"github-tracker/github-tracker/models"
	"github-tracker/github-tracker/repository"
	"github-tracker/github-tracker/repository/entity"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var webhook models.GitHubWebhook
	err := json.Unmarshal([]byte(request.Body), &webhook)
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

func main() {
	lambda.Start(handleRequest)
}
