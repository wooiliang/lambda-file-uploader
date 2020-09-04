package main

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response struct
type Response events.APIGatewayProxyResponse

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	file, _ := base64.StdEncoding.DecodeString(request.Body)
	return Response{
		StatusCode:      200,
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type": request.Headers["Content-Type"],
		},
		Body: base64.StdEncoding.EncodeToString(file),
	}, nil
}

func main() {
	lambda.Start(handler)
}
