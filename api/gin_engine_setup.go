package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fabric8-services/fabric8-wit/api/authz"
	"github.com/fabric8-services/fabric8-wit/api/handler"
	"github.com/fabric8-services/fabric8-wit/auth"
	"github.com/fabric8-services/fabric8-wit/configuration"
	"github.com/fabric8-services/fabric8-wit/gormapplication"
	"github.com/fabric8-services/fabric8-wit/notification"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// NewGinEngine instanciates a new HTTP engine to server the requests
func NewGinEngine(appDB *gormapplication.GormDB, notificationChannel notification.Channel, config *configuration.ConfigurationData) *gin.Engine {
	httpEngine := gin.Default()
	httpEngine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	httpEngine.Use(CORSMiddleware())
	authMiddleware := authz.NewJWTAuthMiddleware(appDB)
	spacesResource := handler.NewSpacesResource(appDB, config, auth.NewAuthzResourceManager(config))
	workitemsResource := handler.NewWorkItemsResource(appDB, notificationChannel, config)
	httpEngine.GET("/api/spaces", spacesResource.List)
	httpEngine.GET("/api/spaces/:spaceID", spacesResource.GetByID)
	httpEngine.GET("/api/spaces/:spaceID/workitems", workitemsResource.List)
	httpEngine.GET("/api/workitems/:workitemID", workitemsResource.Show)
	// secured endpoints
	authGroup := httpEngine.Group("/")
	authGroup.Use(authMiddleware.MiddlewareFunc())
	// spaceAuthzService := authz.NewAuthzService(config, appDB)
	// authGroup.Use(authz.AuthzServiceManager(spaceAuthzService))
	authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
	authGroup.POST("/api/spaces/", spacesResource.Create)
	authGroup.POST("/api/spaces/:spaceID/workitems", workitemsResource.Create)
	authGroup.PATCH("/api/workitems/:workitemID", authz.NewWorkItemEditorAuthorizator(appDB, config), workitemsResource.Update)

	// If an /api/* route does not exist, redirect it to /legacyapi/* path
	// to be handled by goa
	httpEngine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Redirect(http.StatusTemporaryRedirect, strings.Replace(c.Request.URL.RequestURI(), "/api/", "/legacyapi/", 1))
		} else {
			c.String(http.StatusNotFound, "Not found!")
		}
	})

	return httpEngine
}
