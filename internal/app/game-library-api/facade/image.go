package facade

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/client/s3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	// ImageTypeCover type for cover image
	ImageTypeCover = "cover"
	// ImageTypeScreenshot type for screenshot images
	ImageTypeScreenshot = "screenshot"

	// MaxImageSizeKB maximum upload image size in KB
	MaxImageSizeKB int64 = 1024
	// MaxCovers maximum number of cover files allowed
	MaxCovers = 1
	// MaxScreenshots maximum number of screenshot files allowed
	MaxScreenshots = 8
)

// allowedImageTypes defines allowed image extensions
var allowedImageTypes = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

// UploadGameImages handles the business logic for uploading game images
func (p *Provider) UploadGameImages(ctx context.Context, coverFiles, screenshotFiles []*multipart.FileHeader) ([]model.File, error) {
	// validate cover file
	if err := validateImage(coverFiles, ImageTypeCover); err != nil {
		return nil, err
	}

	// validate screenshot files
	if err := validateImage(screenshotFiles, ImageTypeScreenshot); err != nil {
		return nil, err
	}

	uploadedFiles := make([]model.File, 0, len(coverFiles)+len(screenshotFiles))
	var mu sync.Mutex

	eg, egCtx := errgroup.WithContext(ctx)

	// process cover images
	for _, fileHeader := range coverFiles {
		eg.Go(func() error {
			coverFile, pErr := p.processFile(egCtx, fileHeader)
			if pErr != nil {
				p.log.Error("failed to process cover image", zap.Error(pErr))
				return pErr
			}
			mu.Lock()
			uploadedFiles = append(uploadedFiles, model.File{
				FileName: fileHeader.Filename,
				FileID:   coverFile.FileID,
				FileURL:  coverFile.FileURL,
				Type:     ImageTypeCover,
			})
			mu.Unlock()
			return nil
		})
	}

	// process screenshot images
	for _, fileHeader := range screenshotFiles {
		eg.Go(func() error {
			screenshotFile, pErr := p.processFile(egCtx, fileHeader)
			if pErr != nil {
				p.log.Error("failed to process screenshot image", zap.Error(pErr))
				return pErr
			}
			mu.Lock()
			uploadedFiles = append(uploadedFiles, model.File{
				FileName: fileHeader.Filename,
				FileID:   screenshotFile.FileID,
				FileURL:  screenshotFile.FileURL,
				Type:     ImageTypeScreenshot,
			})
			mu.Unlock()
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	return uploadedFiles, nil
}

// validateImage validates image files against constraints
func validateImage(files []*multipart.FileHeader, imageType string) error {
	maxFiles := MaxCovers
	if imageType == ImageTypeScreenshot {
		maxFiles = MaxScreenshots
	}

	if len(files) == 0 {
		return errors.New("no files provided")
	}
	if len(files) > maxFiles {
		return fmt.Errorf("too many files, maximum is %d", maxFiles)
	}

	maxSizeBytes := MaxImageSizeKB * 1024
	for _, fileHeader := range files {
		// validate file extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedImageTypes[ext] {
			return fmt.Errorf("unsupported file type %s, use .png, .jpg, or .jpeg", ext)
		}

		// validate file size
		if fileHeader.Size > maxSizeBytes {
			return fmt.Errorf("file size exceeds maximum of %d KB", MaxImageSizeKB)
		}
	}
	return nil
}

// processFile processes a single file for upload
func (p *Provider) processFile(ctx context.Context, fileHeader *multipart.FileHeader) (s3.UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	file, err := fileHeader.Open()
	if err != nil {
		return s3.UploadResult{}, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// determine content type based on extension
	contentType := "image/png"
	if ext == ".jpg" || ext == ".jpeg" {
		contentType = "image/jpeg"
	}

	// upload to s3
	result, err := p.s3Client.Upload(ctx, file, contentType, map[string]string{
		"fileName": fileHeader.Filename,
	})
	if err != nil {
		return s3.UploadResult{}, fmt.Errorf("failed to upload to S3: %v", err)
	}

	return result, nil
}
