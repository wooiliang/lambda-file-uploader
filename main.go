package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response struct
type Response events.APIGatewayProxyResponse

func parseForm(key, mpheader string, body io.Reader) (string, string, []byte, error) {
	var buf []byte
	mediaType, params, err := mime.ParseMediaType(mpheader)
	if err != nil {
		return "", "", nil, err
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(body, params["boundary"])
		fmt.Printf("DEBUG:: boundary: %v\n", params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				fmt.Printf("DEBUG:: parse EOF\n")
				break
			}
			if err != nil {
				return "", "", nil, err
			}
			if p.FormName() == key {
				buf, err = ioutil.ReadAll(p)
				if err != nil {
					return "", "", nil, err
				}
				return p.FileName(), p.Header.Get("Content-Type"), buf, nil
			}
		}
	}
	return "", "", nil, errors.New("key not found")
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	body, _ := base64.StdEncoding.DecodeString(request.Body)
	_, contentType, file, err := parseForm("fileToUpload", request.Headers["content-type"], bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	return Response{
		StatusCode:      200,
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
		Body: base64.StdEncoding.EncodeToString(file),
	}, nil
}

func main() {
	lambda.Start(handler)
}
