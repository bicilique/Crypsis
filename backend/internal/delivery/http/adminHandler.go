package http

import (
	"crypsis-backend/internal/delivery/middlewere"
	"crypsis-backend/internal/helper"
	"crypsis-backend/internal/model"
	"crypsis-backend/internal/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type AdminHandler struct {
	appService   services.ApplicationInterface
	fileService  services.FileInterface
	adminService services.AdminInterface
	validator    *validator.Validate
}

func NewAdminHandler(appService services.ApplicationInterface, adminService services.AdminInterface, fileService services.FileInterface) *AdminHandler {
	return &AdminHandler{
		appService:   appService,
		adminService: adminService,
		fileService:  fileService,
		validator:    validator.New(),
	}
}

func (a *AdminHandler) Login(c *gin.Context) {
	var request model.LoginRequest
	var err error

	// Use ShouldBind to bind form data
	if err = c.BindJSON(&request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}

	// DOLOGIN
	result, err := a.adminService.Login(c, request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.AdminErrWrongCredentials):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Invalid username or password", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Login successful", map[string]interface{}{"access_token": result})

}

func (a *AdminHandler) Logout(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	accessToken, _ := middlewere.GetAccessTokenFromHeader(c)
	if accessToken == "" {
		model.JSONErrorResponse(c, http.StatusInternalServerError, "Access token not found", "")
		return
	}

	result, err := a.adminService.RevokeToken(c, adminID, accessToken)
	if err != nil {
		model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Logout successful", result)

}

func (a *AdminHandler) RefreshToken(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	accessToken, _ := middlewere.GetAccessTokenFromHeader(c)
	if accessToken == "" {
		model.JSONErrorResponse(c, http.StatusInternalServerError, "Access token not found", "")
		return
	}

	result, err := a.adminService.RefreshToken(c, adminID, accessToken)
	if err != nil {
		model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Refresh token successful", map[string]interface{}{"access_token": result})

}

func (a *AdminHandler) UpdateAdminUsername(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	var request model.UpdateAdminUsernameRequest
	var err error

	// Use ShouldBind to bind form data
	if err = c.BindJSON(&request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}

	err = a.adminService.UpdateUsername(c, adminID, request.Username)
	if err != nil {
		switch {
		case errors.Is(err, model.AdminErrWrongCredentials):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Invalid username or password", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Updated username successfully", "New username is "+request.Username)

}

func (a *AdminHandler) UpdateAdminPassword(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	var request model.UpdateAdminPasswordRequest
	var err error

	if err = c.BindJSON(&request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}

	err = a.adminService.UpdatePassword(c, adminID, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.AdminErrWrongCredentials):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Invalid username or password", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Updated password successfully", nil)

}

func (a *AdminHandler) DeleteAdmin(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	id := c.Query("id")
	if id == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid id", "Id is required")
		return
	} else if id == adminID {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid id", "You cannot delete yourself")
		return
	}

	result, err := a.adminService.DeleteAdmin(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, model.AdminErrNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Invalid admin id", err.Error())
		case errors.Is(err, model.AdminErrWrongCredentials):
			model.JSONErrorResponse(c, http.StatusUnauthorized, "Invalid admin id", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Apps fetched successfully", result)

}

func (a *AdminHandler) AddAdmin(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}
	var request model.LoginRequest
	var err error

	if err = c.BindJSON(&request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid form data", err.Error())
		return
	}

	if err = a.validator.Struct(request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", model.WrapValidationError(err))
		return
	}

	result, err := a.adminService.Register(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, model.AdminErrAlreadyExists):
			model.JSONErrorResponse(c, http.StatusConflict, "Failed to add admin", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to add admin", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Admin added successfully", result)

}

func (a *AdminHandler) ListAdmin(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
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
	sortBy, order = helper.ValidateSortParams(sortBy, order, helper.AllowedAdminSortFields)

	result, err := a.adminService.GetAdminList(c.Request.Context(), offset, limit, sortBy, order)
	if err != nil {
		switch {
		case errors.Is(err, model.AdminErrAlreadyExists):
			model.JSONErrorResponse(c, http.StatusConflict, "Failed to add admin", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to add admin", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Fetch admin list successfully", result)
}

func (a *AdminHandler) AddApp(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	var request model.AddAppRequest
	var err error

	if err = c.BindJSON(&request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err = a.validator.Struct(request); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid request body", model.WrapValidationError(err))
		return
	}

	result, err := a.appService.AddApp(c.Request.Context(), request.Name, request.Uri, request.RedirectUri)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppAlreadyExists):
			model.JSONErrorResponse(c, http.StatusConflict, "Failed to add app", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to add app", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return

	}
	model.JSONSuccessResponse(c, http.StatusOK, "App added successfully", result)

}

func (a *AdminHandler) GetApp(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}
	appID := c.Param("id")
	if appID == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid app id", "App id is required")
		return
	}

	result, err := a.appService.GetInfo(c.Request.Context(), appID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to get app", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to get app", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "App fetched successfully", result)

}

func (a *AdminHandler) ListApps(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
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
	sortBy, order = helper.ValidateSortParams(sortBy, order, helper.AllowedAppSortFields)

	count, result, err := a.appService.ListApps(c.Request.Context(), limit, offset, sortBy, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to list apps", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponseWithCount(c, http.StatusOK, "Apps fetched successfully", count, result)

}

func (a *AdminHandler) DeleteApp(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}
	appID := c.Param("id")
	if appID == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid app id", "App id is required")
		return
	}
	err := a.appService.DeleteApp(c.Request.Context(), appID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to delete app", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to delete app", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "App deleted successfully", nil)

}

func (a *AdminHandler) RecoverApp(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}
	appID := c.Param("id")
	if appID == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid app id", "App id is required")
		return
	}
	result, err := a.appService.RecoverApp(c.Request.Context(), appID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppAlreadyActive):
			model.JSONErrorResponse(c, http.StatusConflict, "Failed to recover app", err.Error())
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to recover app", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to recover app", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "App recovered successfully", result)

}

func (a *AdminHandler) RotateSecret(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}
	var err error
	appID := c.Param("id")
	if appID == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid app id", "App id is required")
		return
	}

	result, err := a.appService.RotateSecret(c.Request.Context(), appID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrAppNotFound):
			model.JSONErrorResponse(c, http.StatusNotFound, "Failed to rotate secret", err.Error())
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to rotate secret", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Secret rotated successfully", result)

}

func (a *AdminHandler) ListFiles(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
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

	count, result, err := a.fileService.ListFilesForAdmin(c.Request.Context(), adminID, "", limit, offset, sortBy, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to list files", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponseWithCount(c, http.StatusOK, "Files fetched successfully", count, result)

}

func (a *AdminHandler) ListFilesByAppId(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	appID := c.Param("id")
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

	if appID == "" {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Invalid app id", "App id is required")
		return
	}

	count, result, err := a.fileService.ListFilesForAdmin(c.Request.Context(), adminID, appID, limit, offset, sortBy, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to list files", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponseWithCount(c, http.StatusOK, "Files fetched successfully", count, result)

}

func (a *AdminHandler) ListLogs(c *gin.Context) {
	_, isAllowed := middlewere.GetUserIDFromToken(c)
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
	sortBy, order = helper.ValidateSortParams(sortBy, order, helper.AllowedLogSortFields)

	count, result, err := a.fileService.ListLogs(c.Request.Context(), limit, offset, sortBy, order)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to list logs", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponseWithCount(c, http.StatusOK, "Logs fetched successfully", count, result)
}

func (a *AdminHandler) Rekey(c *gin.Context) {
	adminID, isAllowed := middlewere.GetUserIDFromToken(c)
	if !isAllowed {
		return
	}

	var rekey model.RekeyRequest
	if err := c.ShouldBindJSON(&rekey); err != nil {
		model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to rekey", err.Error())
		return

	}
	keyUID := rekey.KeyUID
	resukt, err := a.fileService.ReKey(c.Request.Context(), adminID, keyUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidInput):
			model.JSONErrorResponse(c, http.StatusBadRequest, "Failed to rekey", err.Error())
		default:
			model.JSONErrorResponse(c, http.StatusInternalServerError, "Something went wrong", err.Error())
		}
		return
	}
	model.JSONSuccessResponse(c, http.StatusOK, "Rekeyed successfully", resukt)
}
