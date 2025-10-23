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
            fileInfo.MimeType = detectMimeTypeByExtension(fileName)

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
    } else {
        info.MimeType = detectMimeTypeByExtension(filePath)
    }

    return info, nil
}

// detectMimeTypeByExtension пытается определить MIME-тип по расширению имени файла
// Сначала использует стандартный mime.TypeByExtension, затем применяет расширенную карту
// наиболее распространённых типов. Возвращает "application/octet-stream" по умолчанию.
func detectMimeTypeByExtension(fileName string) string {
    ext := strings.ToLower(path.Ext(fileName))
    if ext == "" {
        return "application/octet-stream"
    }

    if mt := mime.TypeByExtension(ext); mt != "" {
        return mt
    }

    switch ext {
    // Images
    case ".jpg", ".jpeg", ".jpe":
        return "image/jpeg"
    case ".png":
        return "image/png"
    case ".gif":
        return "image/gif"
    case ".bmp":
        return "image/bmp"
    case ".webp":
        return "image/webp"
    case ".svg":
        return "image/svg+xml"
    case ".heic":
        return "image/heic"

    // Video
    case ".mp4", ".m4v":
        return "video/mp4"
    case ".mov":
        return "video/quicktime"
    case ".webm":
        return "video/webm"
    case ".avi":
        return "video/x-msvideo"
    case ".mkv":
        return "video/x-matroska"

    // Audio
    case ".mp3":
        return "audio/mpeg"
    case ".wav":
        return "audio/wav"
    case ".flac":
        return "audio/flac"
    case ".ogg", ".oga":
        return "audio/ogg"

    // Documents and archives
    case ".pdf":
        return "application/pdf"
    case ".zip":
        return "application/zip"
    case ".tar":
        return "application/x-tar"
    case ".gz", ".tgz":
        return "application/gzip"
    case ".7z":
        return "application/x-7z-compressed"
    case ".rar":
        return "application/vnd.rar"

    // Office
    case ".xlsx":
        return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
    case ".xls":
        return "application/vnd.ms-excel"
    case ".csv":
        return "text/csv"
    case ".doc":
        return "application/msword"
    case ".docx":
        return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    case ".ppt":
        return "application/vnd.ms-powerpoint"
    case ".pptx":
        return "application/vnd.openxmlformats-officedocument.presentationml.presentation"

    // Text and markup
    case ".txt":
        return "text/plain"
    case ".md":
        return "text/markdown"
    case ".json":
        return "application/json"
    case ".yaml", ".yml":
        return "application/x-yaml"
    case ".xml":
        return "application/xml"

    // Code (best-effort)
    case ".js":
        return "application/javascript"
    case ".ts":
        return "text/typescript"
    case ".go":
        return "text/x-go"
    case ".py":
        return "text/x-python"
    case ".java":
        return "text/x-java-source"
    case ".rb":
        return "text/x-ruby"
    case ".php":
        return "application/x-php"
    case ".cpp", ".cc", ".cxx":
        return "text/x-c++src"
    case ".c":
        return "text/x-csrc"
    case ".cs":
        return "text/x-csharp"
    case ".sh", ".bash":
        return "text/x-shellscript"
    }

    return "application/octet-stream"
}
