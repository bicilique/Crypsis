package http

import (
	"crypsis-backend/internal/delivery/middlewere"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	Router          *gin.Engine
	ClientHandler   *ClientHandler
	AdminHandler    *AdminHandler
	HydraAdminURL   string
	TokenMiddlewere middlewere.TokenMiddlewareConfig
}

func (c *RouterConfig) Setup() {
	c.setupClient()
	c.setupAdmin()
	c.setupPublic()
}

func (c *RouterConfig) setupPublic() {
	group := c.Router.Group("/api")
	group.POST("/admin/login", c.AdminHandler.Login)
}

func (c *RouterConfig) setupClient() {
	group := c.Router.Group("/api")
	// group.Use(middlewere.TokenMiddleware(c.getTokenMiddlewareConfig()))
	group.Use(middlewere.TokenMiddleware(c.TokenMiddlewere))

	group.POST("/files", c.ClientHandler.UploadFile)
	group.GET("/files/:id/download", c.ClientHandler.DownloadFile)
	group.PUT("/files/:id/update", c.ClientHandler.UpdateFile)
	group.DELETE("/files/:id/delete", c.ClientHandler.DeleteFile)

	group.GET("/files/list", c.ClientHandler.ListFiles)
	group.GET("/files/:id/metadata", c.ClientHandler.MetaDataFile)

	group.POST("/files/encrypt", c.ClientHandler.EncryptFile)
	group.POST("/files/decrypt", c.ClientHandler.DecryptFile)

	// Temporarily disabled due to potential security issues
	// group.POST("/files/:id/recover", c.ClientHandler.RecoverFile)
}

func (c *RouterConfig) setupAdmin() {
	group := c.Router.Group("/api")
	// group.Use(middlewere.TokenMiddleware(c.getTokenMiddlewareConfig()))
	// group.Use(middlewere.TokenMiddleware(c.TokenMiddlewere))
	group.Use(middlewere.AdminTokenMiddleware(c.TokenMiddlewere))

	// Admin Account Management
	group.GET("/admin/logout", c.AdminHandler.Logout)
	group.GET("/admin/refresh-token", c.AdminHandler.RefreshToken)
	group.GET("/admin/list", c.AdminHandler.ListAdmin)
	group.PATCH("/admin/username", c.AdminHandler.UpdateAdminUsername)
	group.PATCH("/admin/password", c.AdminHandler.UpdateAdminPassword)
	group.DELETE("/admin", c.AdminHandler.DeleteAdmin)
	group.POST("/admin/add", c.AdminHandler.AddAdmin)

	// Application / Client Management
	group.POST("/admin/apps", c.AdminHandler.AddApp)
	group.GET("/admin/apps", c.AdminHandler.ListApps)
	group.GET("/admin/apps/:id", c.AdminHandler.GetApp)
	group.DELETE("/admin/apps/:id", c.AdminHandler.DeleteApp)
	group.POST("/admin/apps/:id/recover", c.AdminHandler.RecoverApp)
	group.PUT("/admin/apps/:id/rotate-secret", c.AdminHandler.RotateSecret)

	// File Management
	group.GET("/admin/files", c.AdminHandler.ListFiles)
	group.GET("/admin/apps/:id/files", c.AdminHandler.ListFilesByAppId)
	group.GET("/admin/logs", c.AdminHandler.ListLogs)
	group.POST("/admin/files/re-key", c.AdminHandler.Rekey)
}
