package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(settings Settings, filePath string) {
	if settings.S3Bucket == "" || settings.S3Directory == "" {
		fmt.Println("S3 settings not set, skipping upload.")
		return
	}

	awsAccessKeyIDBytes, _ := base64.StdEncoding.DecodeString(settings.AwsAccessKeyIdEncoded)
	awsSecretAccessKeyBytes, _ := base64.StdEncoding.DecodeString(settings.AwsSecretAccessKeyEncoded)
	awsRegion := settings.AwsRegion

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(string(awsAccessKeyIDBytes), string(awsSecretAccessKeyBytes), ""),
	})
	if err != nil {
		fmt.Println("Failed to create session:", err)
		notifyS3UploadError(settings, err.Error())
		return
	}

	svc := s3.New(sess)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		notifyS3UploadError(settings, err.Error())
		return
	}
	defer file.Close()

	bucketName := settings.S3Bucket
	keyName := fmt.Sprintf("%s/%s", settings.S3Directory, file.Name())

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		notifyS3UploadError(settings, err.Error())
		return
	}

	fmt.Println("File successfully uploaded to", bucketName)
}

func notifyS3UploadError(settings Settings, errorMessage string) {
	notifyWebhook := settings.Webhook != ""
	if notifyWebhook {
		message := "Could not upload to S3 bucket. " + errorMessage + " Server Name: " + settings.ServerName
		SendWebhook(settings.Webhook, message)
	}
}
