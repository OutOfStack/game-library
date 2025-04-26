package api

import (
	"errors"
	"net/http"

	api "github.com/OutOfStack/game-library/internal/app/game-library-api/api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	// maxFormMemory maximum amount of memory used to store multipart form data
	maxFormMemory = 32 << 20 // 32 MB
)

// UploadGameImages godoc
// @Summary Upload game images
// @Description uploads cover and screenshots images
// @ID upload-game-images
// @Accept  multipart/form-data
// @Produce json
// @Param   cover 		formData file 	false "Cover image file (.png, .jpg, .jpeg), maximum 1MB"
// @Param   screenshots formData []file false "Screenshot image files (.png, .jpg, .jpeg), up to 8 files, maximum 1MB each" collectionFormat(multi)
// @Success 201 {object} api.UploadImagesResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/images [post]
func (p *Provider) UploadGameImages(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "uploadGameImages")
	defer span.End()

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.log.Error("get claims from context", zap.Error(err))
		web.Respond500(w)
		return
	}

	span.SetAttributes(attribute.String("user.id", claims.UserID()))

	// parse multipart form with a reasonable max memory
	err = r.ParseMultipartForm(maxFormMemory)
	if err != nil {
		web.RespondError(w, web.NewError(errors.New("failed to parse form"), http.StatusBadRequest))
		return
	}

	// process cover image if provided
	coverFiles := r.MultipartForm.File["cover"]
	screenshotFiles := r.MultipartForm.File["screenshots"]

	// Use the facade to handle the business logic
	uploadedFiles, err := p.gameFacade.UploadGameImages(ctx, coverFiles, screenshotFiles)
	if err != nil {
		p.log.Error("upload game images", zap.Error(err))
		web.RespondError(w, web.NewError(err, http.StatusBadRequest))
		return
	}

	files := make([]api.UploadedFileInfo, len(uploadedFiles))
	for i, f := range uploadedFiles {
		files[i] = api.UploadedFileInfo{
			FileName: f.FileName,
			FileID:   f.FileID,
			FileURL:  f.FileURL,
			Type:     f.Type,
		}
	}

	web.Respond(w, api.UploadImagesResponse{
		Files: files,
	}, http.StatusCreated)
}
