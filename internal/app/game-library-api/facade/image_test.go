package facade_test

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/OutOfStack/game-library/internal/client/s3"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"go.uber.org/mock/gomock"
)

const (
	coverFormDataParam       = "cover"
	screenshotsFormDataParam = "screenshots"

	contentLength = 1024 // 1 KB

	fileNameS3MDField = "fileName"

	jpg = ".jpg"
	png = ".png"
)

func (s *TestSuite) TestUploadGameImages_Success() {
	// create file headers
	coverFileName, scrFileName := td.String()+jpg, td.String()+png
	coverFile, err := s.createFileHeader(coverFormDataParam, coverFileName, td.Bytesn(contentLength))
	s.NoError(err)
	screenshotFile, err := s.createFileHeader(screenshotsFormDataParam, scrFileName, td.Bytesn(contentLength))
	s.NoError(err)

	coverFileID, coverFileURL := td.String(), "https://"+td.String()+".com/"+coverFileName
	scrFileID, scrFileURL := td.String(), "https://"+td.String()+".com/"+scrFileName

	s.s3ClientMock.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any(), map[string]string{fileNameS3MDField: coverFileName}).
		Return(s3.UploadResult{FileID: coverFileID, FileURL: coverFileURL}, nil).
		Times(1)

	s.s3ClientMock.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any(), map[string]string{fileNameS3MDField: scrFileName}).
		Return(s3.UploadResult{FileID: scrFileID, FileURL: scrFileURL}, nil).
		Times(1)

	res, err := s.provider.UploadGameImages(s.ctx, []*multipart.FileHeader{coverFile}, []*multipart.FileHeader{screenshotFile})

	s.NoError(err)
	s.Len(res, 2)

	fileNames := []string{res[0].FileName, res[1].FileName}
	s.Contains(fileNames, coverFileName)
	s.Contains(fileNames, scrFileName)
	fileIDs := []string{res[0].FileID, res[1].FileID}
	s.Contains(fileIDs, coverFileID)
	s.Contains(fileIDs, scrFileID)
	fileURLs := []string{res[0].FileURL, res[1].FileURL}
	s.Contains(fileURLs, coverFileURL)
	s.Contains(fileURLs, scrFileURL)
}

func (s *TestSuite) TestUploadGameImages_InvalidCoverFile() {
	// create invalid cover file (too large)
	coverFile, err := s.createFileHeader(coverFormDataParam, td.String()+jpg, td.Bytesn(2*1024*1024)) // 2MB, exceeds max size
	s.NoError(err)
	screenshotFile, err := s.createFileHeader(screenshotsFormDataParam, td.String()+png, td.Bytesn(200*1024)) // 200KB
	s.NoError(err)

	result, err := s.provider.UploadGameImages(s.ctx, []*multipart.FileHeader{coverFile}, []*multipart.FileHeader{screenshotFile})

	s.Error(err)
	s.Contains(err.Error(), "file size exceeds maximum")
	s.Empty(result)
}

func (s *TestSuite) TestUploadGameImages_InvalidScreenshotFile() {
	// create valid cover file but invalid screenshot file (unsupported type)
	coverFile, err := s.createFileHeader(coverFormDataParam, td.String()+jpg, td.Bytesn(contentLength))
	s.NoError(err)
	screenshotFile, err := s.createFileHeader(screenshotsFormDataParam, td.String()+".bmp", td.Bytesn(contentLength)) // unsupported type
	s.NoError(err)

	result, err := s.provider.UploadGameImages(s.ctx, []*multipart.FileHeader{coverFile}, []*multipart.FileHeader{screenshotFile})

	s.Error(err)
	s.Contains(err.Error(), "unsupported file type")
	s.Empty(result)
}

func (s *TestSuite) TestUploadGameImages_TooManyFiles() {
	// create too many cover files
	coverFile1, err := s.createFileHeader(coverFormDataParam, td.String()+jpg, td.Bytesn(contentLength))
	s.NoError(err)
	coverFile2, err := s.createFileHeader(coverFormDataParam, td.String()+jpg, td.Bytesn(contentLength))
	s.NoError(err)
	screenshotFile, err := s.createFileHeader(screenshotsFormDataParam, td.String()+png, td.Bytesn(contentLength))
	s.NoError(err)

	result, err := s.provider.UploadGameImages(s.ctx, []*multipart.FileHeader{coverFile1, coverFile2}, []*multipart.FileHeader{screenshotFile})

	s.Error(err)
	s.Contains(err.Error(), "too many files")
	s.Empty(result)
}

func (s *TestSuite) TestUploadGameImages_S3UploadError() {
	coverFile, err := s.createFileHeader(coverFormDataParam, td.String()+jpg, td.Bytesn(contentLength))
	s.NoError(err)
	screenshotFile, err := s.createFileHeader(screenshotsFormDataParam, td.String()+png, td.Bytesn(contentLength))
	s.NoError(err)

	s.s3ClientMock.EXPECT().
		Upload(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(s3.UploadResult{}, errors.New("s3 upload failed")).
		AnyTimes()

	result, err := s.provider.UploadGameImages(s.ctx, []*multipart.FileHeader{coverFile}, []*multipart.FileHeader{screenshotFile})

	s.Error(err)
	s.Contains(err.Error(), "failed to upload to S3")
	s.Empty(result)
}

func (s *TestSuite) TestUploadGameImages_NoFiles() {
	result, err := s.provider.UploadGameImages(s.ctx, []*multipart.FileHeader{}, []*multipart.FileHeader{})

	s.Error(err)
	s.Contains(err.Error(), "no files provided")
	s.Empty(result)
}

func (s *TestSuite) createFileHeader(fieldName, fileName string, content []byte) (*multipart.FileHeader, error) {
	// prepare a buffer and multipart writer
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// add a form file part
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile: %w", err)
	}

	if _, err = part.Write(content); err != nil {
		return nil, fmt.Errorf("writing content: %w", err)
	}

	// close the writer to flush the ending boundary
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("closing writer: %w", err)
	}

	// parse the multipart data back
	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(int64(len(content) + 1024))
	if err != nil {
		return nil, fmt.Errorf("ReadForm: %w", err)
	}
	defer func() {
		err = form.RemoveAll()
		s.NoError(err)
	}()

	files := form.File[fieldName]
	if len(files) == 0 {
		return nil, fmt.Errorf("no file for field %q", fieldName)
	}
	return files[0], nil
}
