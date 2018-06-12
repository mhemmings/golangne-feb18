package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

var count = 0

func main() {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			count++
		}
	}()

	lambda.Start(hello)
}

func hello() (string, error) {
	return fmt.Sprintf("Count is %d", count), nil
}
