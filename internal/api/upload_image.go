package api

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/middleware"
	"github.com/OutOfStack/game-library/internal/web"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	// maxRequestBodySize upper bound on the request body size to prevent memory exhaustion
	maxRequestBodySize = 32 << 20 // 32 MB
	// maxFormMemory maximum amount of memory used to store multipart form data
	maxFormMemory = maxRequestBodySize / 2
)

// UploadGameImages uploads cover and screenshots images
func (p *Provider) UploadGameImages(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "uploadGameImages")
	defer span.End()

	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		p.respondError(w, "get claims from context", err)
		return
	}

	span.SetAttributes(attribute.String("user.id", claims.UserID()))

	// limit request body size to prevent memory exhaustion, then parse multipart form
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
	err = r.ParseMultipartForm(maxFormMemory) // #nosec G120 -- r.Body is bounded by http.MaxBytesReader above
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			web.RespondError(w, web.NewErrorFromMessage("request body too large", http.StatusRequestEntityTooLarge))
			return
		}
		web.RespondError(w, web.NewErrorFromMessage("failed to parse form", http.StatusBadRequest))
		return
	}

	coverFiles := r.MultipartForm.File["cover"]
	screenshotFiles := r.MultipartForm.File["screenshots"]

	uploadedFiles, err := p.gameFacade.UploadGameImages(ctx, coverFiles, screenshotFiles, claims.Name)
	if err != nil {
		p.respondError(w, "upload game images", err, zap.String("publisher", claims.Name))
		return
	}

	files := make([]UploadedFileInfo, len(uploadedFiles))
	for i, f := range uploadedFiles {
		files[i] = UploadedFileInfo{
			FileName: f.FileName,
			FileID:   f.FileID,
			FileURL:  f.FileURL,
			Type:     UploadedFileInfoType(f.Type),
		}
	}

	web.Respond(w, UploadImagesResponse{
		Files: files,
	}, http.StatusCreated)
}
