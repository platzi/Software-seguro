package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github-tracker/github-tracker/database"
	"github-tracker/github-tracker/repository"
	"github-tracker/github-tracker/repository/entity"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Author  string   `json:"author"`
	Count   int      `json:"commits_count"`
	Commits []Commit `json:"commits"`
}

type Commit struct {
	ID      string `json:"id"`
	Repo    string `json:"repo"`
	Message string `json:"message"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	email := request.PathParameters["auth_email"]

	if email == "" {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusBadRequest,
			Body:            "missing auth_email param",
		}, fmt.Errorf("missing path parameter: auth_email")
	}

	db, err := database.Connect(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error initializing db",
		}, err
	}

	commitRepo := repository.NewCommit(db)

	commits, err := commitRepo.GetCommitsByAuthorEmail(ctx, email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error getting commits in the db",
		}, err
	}

	response := parseResponse(email, commits)

	jsonData, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			IsBase64Encoded: false,
			StatusCode:      http.StatusInternalServerError,
			Body:            "error marshaling response",
		}, err
	}

	return events.APIGatewayProxyResponse{
		IsBase64Encoded: false,
		StatusCode:      http.StatusOK,
		Body:            string(jsonData),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func parseResponse(email string, commits []entity.Commit) Response {
	responseCommits := []Commit{}

	for _, item := range commits {
		commit := Commit{
			ID:      item.CommitID,
			Message: item.CommitMessage,
			Repo:    item.RepoName,
		}

		responseCommits = append(responseCommits, commit)
	}

	return Response{
		Author:  email,
		Count:   len(commits),
		Commits: responseCommits,
	}
}

func main() {
	lambda.Start(handleRequest)
}
