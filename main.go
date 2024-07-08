package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func main() {

	ctx := context.Background()

	config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load default configuration: %v", err)
	}
	client := s3.NewPresignClient(s3.NewFromConfig(config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
	}))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /upload-credential", func(w http.ResponseWriter, r *http.Request) {

		req := &GetUploadCredentialRequest{}

		// 要求されたファイルサイズが100MBより大きい場合はエラー
		if 100*1024*1024 < req.Size {
			// ...
		}
		// 要求された拡張子が、jpg/jpeg/jpe以外の場合はエラー
		if !slices.Contains([]string{"jpg", "jpeg", "jpe"}, req.Extension) {
			// ...
		}
		presign, err := client.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket:               aws.String("your-bucket-name"),
			Key:                  aws.String("your-object-name." + req.Extension),
			ContentType:          aws.String(mime.TypeByExtension("." + req.Extension)),
			ContentLength:        aws.Int64(req.Size),
			ServerSideEncryption: types.ServerSideEncryptionAes256,
		}, s3.WithPresignExpires(15*time.Minute))

		if err != nil {
			log.Fatalf("failed to generate a presigned URL: %v", err)
		}
		fmt.Printf("the presigned URL is: %s\n", presign.URL)
		fmt.Printf("the HTTP method is: %s\n", presign.Method)
		fmt.Printf("the HTTP header is: %v\n", presign.SignedHeader)

		if err := json.NewEncoder(w).Encode(&GetUploadCredentialResponse{
			Method:  presign.Method,
			URL:     presign.URL,
			Headers: presign.SignedHeader,
		}); err != nil {
			log.Printf("failed to encode: %v", err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type GetUploadCredentialRequest struct {
	Size      int64  `query:"size"`
	Extension string `query:"extension"`
}

type GetUploadCredentialResponse struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"header"`
}

// https://prod-omoide-app.s3.ap-northeast-1.amazonaws.com/images/46z7n2w8k3kieb5nv1ll/c9380qd7eapc70t9no60.jpg?

// X-Amz-Algorithm=AWS4-HMAC-SHA256&
// X-Amz-Credential=ASIA3NNHFLAFMHALHVPS%2F20240702%2Fap-northeast-1%2Fs3%2Faws4_request&
// X-Amz-Date=20240702T123416Z&
// X-Amz-Expires=21600&
// X-Amz-Security-Token=IQoJb3JpZ2luX2VjELz%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaDmFwLW5vcnRoZWFzdC0xIkgwRgIhAOvAG5XmHD1oNZmjB32oKbzfDtAaO7o3%2BiPiYIg8XziIAiEA1lA7GKL5aRlCdWcEUhtQ1U9UcQ9%2FIY3lIv2l%2Bjo1oQIqgQMIdhACGgw3ODQ3MTg5NzcwMzQiDIpn%2BOb24%2B9Jbvhh1SreAh%2BupjqcjpKMKfbLorl%2FBUq2nxK1ixtIa6F6n2LNzt7pjiATilgTaD0R5G%2F6V5DdjEG7hqmJH6rPMCnBQz9DHsaPnpFZa08U1fEFC8o81cgqyQsvGE2NAW7dpRSRREn2VXeJi00GeFbsMVqlPRYYiQIqLvmbjfsNWKp2b62Q3HpJsJ2UWE992gjr7ahkQT2KOHCkqShCtZ6QSKQxjDn0Duxbxh7%2FnwOj%2Bzr4lmil6j%2BUgIJfspQNjxA1eUJJkAyWYU8yTe4boFhKmYk740XeEfyTV5kXHR8Z0G9fJKm7cklJ%2B3csKHUn%2FLskypKn3I%2BLkedzFDflMIqwwqh3vfIQIu17krLvygzzj6dmDR7CiUf8uZZHsOD%2B%2FxulNiJcfOBZMbw0hSHi8Wp%2F%2FzEgd%2B0ypbr6qi8qKx3BOvUu3crMZ%2BDom%2FpEAyxJe89MC4vTCuhwzNlYU40THhG1vheiMz6xMMjnj7QGOp0B43po9NSmvyrFyyKHpANazKt7kf1dy1dVQm09rUPt7JOtjOF38jy%2FGY83mggk9HafeMGe2AHRenK2C3qY9%2F4eSfWk1ckQf9zhSHhVhwjQPHyD7AiEH6x0k%2FT1gzh%2BQq1IkzQ4Nj6CJTB7Ave9hcRW9laH6iS1qGsazMQe23du6HVvWZQ5CeRx%2Bc6SYNX8yaB4H5O93Jl2xNN3kJICdA%3D%3D&
// X-Amz-SignedHeaders=host&
// X-Amz-Signature=291f91ae6a079cbeaab970990e870470bc4dae7901498edd4805161daa5a2034
