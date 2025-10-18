package http

import (
	"context"
	"crypsis-backend/internal/delivery/middlewere"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/services"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	clientService services.FileInterface
}

func NewClientHandler(clietService services.FileInterface) *ClientHandler {
	return &ClientHandler{
		clientService: clietService,
	}
}

func (ch *ClientHandler) UploadFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	ctx := context.WithValue(c.Request.Context(), "request", c.Request)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	result, err := ch.clientService.UploadFile(ctx, clientID, header.Filename, file)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFailedToReadFile):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrFileIsEmpty):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrHashCalculationFailed):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to upload file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "File uploaded successfully", model.UploadFileResponse{
		FileName: header.Filename,
		FileID:   result,
	})
}

func (ch *ClientHandler) DownloadFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	ctx := context.WithValue(c.Request.Context(), "request", c.Request)
	fileID := c.Param("id")
	result, fileName, err := ch.clientService.DownloadFile(ctx, clientID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrFailedToReadFile):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrFileIsEmpty):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrHashCalculationFailed):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}

	mimeType := http.DetectContentType(result)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, mimeType, result)
}

func (ch *ClientHandler) EncryptFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	ctx := context.WithValue(c.Request.Context(), "request", c.Request)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()
	safeFilename := filepath.Base(header.Filename)
	ext := filepath.Ext(safeFilename)
	baseName := strings.TrimSuffix(safeFilename, ext)
	result, fileUID, err := ch.clientService.EncryptFile(ctx, clientID, safeFilename, file)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFailedToReadFile):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrFileIsEmpty):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrHashCalculationFailed):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to upload file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}

	newFilename := fmt.Sprintf("encrypted_%s_%s%s", baseName, fileUID, ext)
	c.Header("X-File-UID", fileUID)
	c.Header("Content-Disposition", "attachment; filename=\""+newFilename+"\"")
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", result)
}

func (ch *ClientHandler) DecryptFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	ctx := context.WithValue(c.Request.Context(), "request", c.Request)
	file, header, err := c.Request.FormFile("file")
	fileId := c.Request.FormValue("id")
	if err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid file upload", err.Error())
		return
	}
	if fileId == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid file ID", "File ID is required")
		return

	}
	defer file.Close()

	result, err := ch.clientService.DecryptFile(ctx, clientID, fileId, file)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrFailedToReadFile):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrFileIsEmpty):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		case errors.Is(err, model.ErrHashCalculationFailed):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to download file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return

	}
	mimeType := http.DetectContentType(result)
	c.Header("Content-Disposition", "attachment; filename="+header.Filename)
	c.Data(http.StatusOK, mimeType, result)
}

func (ch *ClientHandler) UpdateFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	ctx := context.WithValue(c.Request.Context(), "request", c.Request)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()
	fileID := c.Param("id")
	result, err := ch.clientService.UpdateFile(ctx, clientID, fileID, header.Filename, file)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to update file", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to update file", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to update file", err.Error())
		case errors.Is(err, model.ErrFailedToReadFile):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to update file", err.Error())
		case errors.Is(err, model.ErrFileIsEmpty):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to update file", err.Error())
		case errors.Is(err, model.ErrHashCalculationFailed):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to update file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "File updated successfully", model.UploadFileResponse{
		FileName: header.Filename,
		FileID:   result,
	})
}

func (ch *ClientHandler) DeleteFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	fileID := c.Param("id")
	err := ch.clientService.DeleteFile(c.Request.Context(), clientID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to delete file", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to delete file", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to delete file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "File deleted successfully", nil)
}

func (ch *ClientHandler) RecoverFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	fileID := c.Param("id")
	result, err := ch.clientService.RecoverFile(c.Request.Context(), clientID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to recover file", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to recover file", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to recover file", err.Error())
		case errors.Is(err, model.ErrFileAlreadyExists):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to recover file", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "File recovered successfully", result)
}

func (ch *ClientHandler) MetaDataFile(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	fmt.Println("Client ID:", clientID)

	fileID := c.Param("id")
	result, err := ch.clientService.GetFileMetadata(c.Request.Context(), clientID, fileID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())
		case errors.Is(err, model.ErrAppNotActive):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to upload file", err.Error())

		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to get file metadata", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to get file metadata", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to get file metadata", err.Error())
		case errors.Is(err, model.ErrFileAlreadyExists):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to get file metadata", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Fetch file metadata successfully", result)
}

func (ch *ClientHandler) ListFiles(c *gin.Context) {
	clientID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Validate sort parameters using centralized helper
	sortBy, order = helper.ValidateSortParams(sortBy, order, helper.AllowedFileSortFields)

	count, result, err := ch.clientService.ListFiles(c.Request.Context(), clientID, limit, offset, sortBy, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrFileNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to get file metadata", err.Error())
		case errors.Is(err, model.ErrUnauthorizedFileAccess):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Failed to get file metadata", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to get file metadata", err.Error())
		case errors.Is(err, model.ErrFileAlreadyExists):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to get file metadata", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Fetch files metadata successfull", &model.ListFilesResponse{
		Count: count,
		Files: *result,
	})
}
