package s3

import (
    "context"
    "io"
    "mime"
    "path"
    "strings"
    "sync"

    "s3-smart-browser/internal/config"
    "s3-smart-browser/internal/types"

    "github.com/aws/aws-sdk-go-v2/aws"
    awsconfig "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Client struct {
	client *s3.Client
	bucket string
	cfg    *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	optFns := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.AWS.Region),
	}

	if cfg.AWS.AccessKeyID != "" && cfg.AWS.SecretAccessKey != "" {
		optFns = append(optFns, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWS.AccessKeyID,
				cfg.AWS.SecretAccessKey,
				"",
			)),
		)
	}

	if cfg.AWS.Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               cfg.AWS.Endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})
		optFns = append(optFns, awsconfig.WithEndpointResolverWithOptions(customResolver))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &Client{
		client: client,
		bucket: cfg.AWS.Bucket,
		cfg:    cfg,
	}, nil
}

func (c *Client) ListDirectory(ctx context.Context, directoryPath string) (*types.DirectoryListing, error) {
	prefix := c.normalizePath(directoryPath)
	if prefix != "" {
		prefix += "/"
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(c.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	result, err := c.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	listing := &types.DirectoryListing{
		Path: directoryPath,
	}

	// Обработка директорий (CommonPrefixes)
	for _, commonPrefix := range result.CommonPrefixes {
		if commonPrefix.Prefix != nil {
			dirName := strings.TrimPrefix(*commonPrefix.Prefix, prefix)
			dirName = strings.TrimSuffix(dirName, "/")

			listing.Files = append(listing.Files, types.FileInfo{
				Name:        dirName,
				Path:        *commonPrefix.Prefix,
				IsDirectory: true,
			})
		}
	}

	// Обработка файлов с использованием горутин
	fileChan := make(chan types.FileInfo, len(result.Contents))

	for _, obj := range result.Contents {
		if strings.HasSuffix(*obj.Key, "/") || *obj.Key == prefix {
			continue // Пропускаем директории и сам префикс
		}

		wg.Add(1)
        go func(obj s3types.Object) {
			defer wg.Done()

			fileName := strings.TrimPrefix(*obj.Key, prefix)
			fileInfo := types.FileInfo{
				Name:        fileName,
				Path:        *obj.Key,
				Size:        aws.ToInt64(obj.Size),
				IsDirectory: false,
			}

            // Определение MIME-типа по расширению файла (без дополнительных S3 вызовов)
            if ext := strings.ToLower(path.Ext(fileName)); ext != "" {
                if mt := mime.TypeByExtension(ext); mt != "" {
                    fileInfo.MimeType = mt
                } else {
                    fileInfo.MimeType = "application/octet-stream"
                }
            } else {
                fileInfo.MimeType = "application/octet-stream"
            }

			fileChan <- fileInfo
		}(obj)
	}

	wg.Wait()
	close(fileChan)

	for file := range fileChan {
		listing.Files = append(listing.Files, file)
	}

	return listing, nil
}

func (c *Client) GetFileURL(ctx context.Context, filePath string) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(filePath),
	}

	presignClient := s3.NewPresignClient(c.client)
	presignedReq, err := presignClient.PresignGetObject(ctx, input)
	if err != nil {
		return "", err
	}

	return presignedReq.URL, nil
}

func (c *Client) StreamFile(ctx context.Context, filePath string, w io.Writer) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(filePath),
	}

	result, err := c.client.GetObject(ctx, input)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	_, err = io.Copy(w, result.Body)
	return err
}

func (c *Client) normalizePath(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	if c.cfg.App.BaseDirectory != "" {
		path = strings.TrimPrefix(path, c.cfg.App.BaseDirectory)
		path = c.pathJoin(c.cfg.App.BaseDirectory, path)
	}

	return path
}

func (c *Client) pathJoin(elem ...string) string {
	return strings.Trim(strings.Join(elem, "/"), "/")
}

func (c *Client) GetFileInfo(ctx context.Context, filePath string) (*types.FileInfo, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(filePath),
	}

	result, err := c.client.HeadObject(ctx, input)
	if err != nil {
		return nil, err
	}

    info := &types.FileInfo{
		Name:        path.Base(filePath),
		Path:        filePath,
		Size:        aws.ToInt64(result.ContentLength),
		IsDirectory: false,
    }

    // Используем ContentType из S3, если он есть, иначе определяем по расширению
    if result.ContentType != nil && *result.ContentType != "" {
        info.MimeType = *result.ContentType
    } else if ext := strings.ToLower(path.Ext(filePath)); ext != "" {
        if mt := mime.TypeByExtension(ext); mt != "" {
            info.MimeType = mt
        } else {
            info.MimeType = "application/octet-stream"
        }
    } else {
        info.MimeType = "application/octet-stream"
    }

    return info, nil
}
