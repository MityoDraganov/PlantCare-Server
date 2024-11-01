package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var bucket *storage.BucketHandle // Global bucket variable

// InitializeApp initializes the Google Cloud Storage client and sets the default bucket.
func InitializeApp(credentialsFile string) (*storage.Client, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("error initializing Storage client: %v", err)
	}

	// Set up the default bucket handle from the environment variable
	bucketName := os.Getenv("FIREBASE_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("FIREBASE_BUCKET_NAME environment variable is not set")
	}
	bucket = client.Bucket(bucketName)
	return client, nil
}
func UploadFile(file io.Reader, destinationPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	// Ensure bucket is initialized
	if bucket == nil {
		return "", fmt.Errorf("bucket is not initialized")
	}

	// Create a writer for the object in the bucket
	obj := bucket.Object(destinationPath)
	writer := obj.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to write file to bucket: %v", err)
	}

	// Close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Make the uploaded file public
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("failed to make file public: %v", err)
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", os.Getenv("FIREBASE_BUCKET_NAME"), destinationPath)
	return url, nil
}
