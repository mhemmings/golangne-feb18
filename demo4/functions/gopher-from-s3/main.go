package main

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Service *s3.S3
var s3Bucket = os.Getenv("S3_BUCKET")
var gopherImage image.Image

func main() {
	var err error
	s3Service = s3.New(session.New())

	// Get gopher image in memory on cold start
	gopherFile, err := os.Open("gopher.png")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	gopherImage, err = png.Decode(gopherFile)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer gopherFile.Close()

	lambda.Start(handler)
}

func handler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		key := record.S3.Object.Key
		// Using the key from the event, get the image from S3
		inputImage, err := downloadFromS3(s3Bucket, key)
		if err != nil {
			return err
		}

		b := inputImage.Bounds()
		offset := image.Pt(b.Max.X-170, b.Max.Y-110)
		outputImage := image.NewRGBA(b)
		draw.Draw(outputImage, b, inputImage, image.ZP, draw.Src)
		draw.Draw(outputImage, gopherImage.Bounds().Add(offset), gopherImage, image.ZP, draw.Over)

		var outputBytes []byte
		outputBuffer := bytes.NewBuffer(outputBytes)

		// Attempt to keep input type if it was png
		if strings.ToLower(path.Ext(key)) == ".png" {
			err = png.Encode(outputBuffer, outputImage)
		} else {
			err = jpeg.Encode(outputBuffer, outputImage, &jpeg.Options{jpeg.DefaultQuality})
		}
		if err != nil {
			return err
		}

		// Save file back, assuming input/output naming
		err = writeToS3(s3Bucket, strings.Replace(key, "input/", "output/", -1), outputBuffer.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFromS3(bucket string, key string) (image.Image, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	resp, err := s3Service.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if strings.ToLower(path.Ext(key)) == ".png" {
		return png.Decode(resp.Body)
	}

	// Just fallback to jpeg encoding
	return jpeg.Decode(resp.Body)
}

func writeToS3(bucket string, key string, file []byte) error {
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	}

	_, err := s3Service.PutObject(params)
	return err
}
