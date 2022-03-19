package storage

import (
	"context"
	"fmt"
	"log"
	"os"
)

func Bucket(client BucketClient) {
	ctx := context.Background()

	// create(ctx, client)
	// uploadObject(ctx, client)
	downloadObject(ctx, client)
	// deleteObject(ctx, client)
	listObjects(ctx, client)
}

func create(ctx context.Context, client BucketClient) {
	if err := client.Create(ctx, "aws-test"); err != nil {
		log.Fatalln(err)
	}
	log.Println("create: ok")
}

func uploadObject(ctx context.Context, client BucketClient) {
	file, err := os.Open("id.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	url, err := client.UploadObject(ctx, "aws-test", "id.txt", file)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("upload object:", url)
}

func downloadObject(ctx context.Context, client BucketClient) {
	file, err := os.Create("id.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err := client.DownloadObject(ctx, "aws-test", "id.txt", file); err != nil {
		log.Fatalln(err)
	}
	log.Println("download object: ok")
}

func listObjects(ctx context.Context, client BucketClient) {
	objects, err := client.ListObjects(ctx, "aws-test")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list objects:")
	for _, object := range objects {
		fmt.Printf("%+v\n", object)
	}
}
