package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	jose "github.com/go-jose/go-jose/v3"
)

func handler(event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayCustomAuthorizerResponse, error) {
	route := event.RouteArn
	path := event.RequestContext.HTTP.Path

	jsonData, err := json.Marshal(event)
	if err != nil {
		return denyRequest(fmt.Sprintf("error unmarshal event: %s", err.Error()))
	}

	fmt.Println(string(jsonData))

	if path == "/commit" {
		return allowRequest(route)
	}

	authToken, ok := event.Headers["authorization"]
	if !ok {
		return denyRequest("missing auth token")
	}

	tokenString := strings.TrimPrefix(authToken, "Bearer ")
	if tokenString == "" {
		return denyRequest("missing auth token")
	}

	return verifyToken(tokenString, route)
}

func verifyToken(authToken string, route string) (events.APIGatewayCustomAuthorizerResponse, error) {
	_, err := jose.ParseSigned(authToken)
	if err != nil {
		return denyRequest("invalid auth0 token")
	}

	return allowRequest(route)
}

func denyRequest(reason string) (events.APIGatewayCustomAuthorizerResponse, error) {
	fmt.Println(reason)

	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "anonymous",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   "Deny",
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{"*"},
				},
			},
		},
	}, fmt.Errorf(reason)
}

func allowRequest(route string) (events.APIGatewayCustomAuthorizerResponse, error) {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: "user",
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Effect:   "Allow",
					Action:   []string{"execute-api:Invoke"},
					Resource: []string{route},
				},
			},
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
