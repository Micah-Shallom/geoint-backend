package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/Micah-Shallom/geoint-backend/internal/config"
	"github.com/Micah-Shallom/geoint-backend/pkg/repository/storage"
	"github.com/Micah-Shallom/geoint-backend/utility"
	"github.com/minio/minio-go/v7"
)

func UploadImagery(logger *utility.Logger, analysisID string, imageType string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("Failed to open file: %v", err))
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	path := fmt.Sprintf("imagery/%s/%s_%s", analysisID, imageType, file.Filename)
	minioClient := storage.DB.Minio
	bucketName := config.Config.Minio.BucketName

	_, err = minioClient.PutObject(context.Background(), bucketName, path, src, file.Size, minio.PutObjectOptions{
		ContentType: "image/tiff",
	})
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("Failed to upload imagery to %s: %v", path, err))
		return "", fmt.Errorf("failed to upload imagery to %s: %v", path, err)
	}

	utility.LogAndPrint(logger, fmt.Sprintf("Imagery uploaded successfully to %s", path))

	url := fmt.Sprintf("https://%s/%s/%s", minioClient.EndpointURL().Host, bucketName, path)
	return url, nil
}

func UploadChangeMap(logger *utility.Logger, analysisID string, file io.Reader, fileSize int64) (string, error) {
	path := fmt.Sprintf("changemaps/%s/change_map.tif", analysisID)
	minioClient := storage.DB.Minio
	bucketName := config.Config.Minio.BucketName

	_, err := minioClient.PutObject(context.Background(), bucketName, path, file, fileSize, minio.PutObjectOptions{
		ContentType: "image/tiff",
	})
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("Failed to upload change map to %s: %v", path, err))
		return "", fmt.Errorf("failed to upload change map to %s: %v", path, err)
	}

	utility.LogAndPrint(logger, fmt.Sprintf("Change map uploaded successfully to %s", path))

	url := fmt.Sprintf("https://%s/%s/%s", minioClient.EndpointURL().Host, bucketName, path)
	return url, nil
}

func UploadSupportDocument(logger *utility.Logger, analysisID string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("Failed to open document: %v", err))
		return "", fmt.Errorf("failed to open document: %v", err)
	}
	defer src.Close()

	path := fmt.Sprintf("documents/%s/%s", analysisID, file.Filename)
	minioClient := storage.DB.Minio
	bucketName := config.Config.Minio.BucketName

	contentType := "application/octet-stream"
	if file.Header.Get("Content-Type") != "" {
		contentType = file.Header.Get("Content-Type")
	}

	_, err = minioClient.PutObject(context.Background(), bucketName, path, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("Failed to upload document to %s: %v", path, err))
		return "", fmt.Errorf("failed to upload document to %s: %v", path, err)
	}

	utility.LogAndPrint(logger, fmt.Sprintf("Document uploaded successfully to %s", path))

	url := fmt.Sprintf("https://%s/%s/%s", minioClient.EndpointURL().Host, bucketName, path)
	return url, nil
}

func DeleteAnalysisFiles(logger *utility.Logger, analysisID string) error {
	minioClient := storage.DB.Minio
	bucketName := config.Config.Minio.BucketName

	// Delete imagery
	imageryPrefix := fmt.Sprintf("imagery/%s/", analysisID)
	objectsCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    imageryPrefix,
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("Error listing objects: %v", object.Err))
			continue
		}
		err := minioClient.RemoveObject(context.Background(), bucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("Failed to delete %s: %v", object.Key, err))
		}
	}

	// Delete documents
	docsPrefix := fmt.Sprintf("documents/%s/", analysisID)
	objectsCh = minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    docsPrefix,
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("Error listing objects: %v", object.Err))
			continue
		}
		err := minioClient.RemoveObject(context.Background(), bucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("Failed to delete %s: %v", object.Key, err))
		}
	}

	// Delete changemaps
	changemapPrefix := fmt.Sprintf("changemaps/%s/", analysisID)
	objectsCh = minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    changemapPrefix,
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("Error listing objects: %v", object.Err))
			continue
		}
		err := minioClient.RemoveObject(context.Background(), bucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("Failed to delete %s: %v", object.Key, err))
		}
	}

	utility.LogAndPrint(logger, fmt.Sprintf("All files deleted for analysis %s", analysisID))
	return nil
}