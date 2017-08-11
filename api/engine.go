package api

import (
	"github.com/fabric8-services/fabric8-wit/api/resource"
	"github.com/fabric8-services/fabric8-wit/configuration"
	"github.com/fabric8-services/fabric8-wit/gormapplication"
	"github.com/fabric8-services/fabric8-wit/space"
	"github.com/gin-gonic/gin"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go-adapter/gingonic"
)

// NewAPI2goEngine instantiates a new HTTP engine to serve requests
func NewAPI2goEngine(appDB *gormapplication.GormDB, config *configuration.ConfigurationData) *gin.Engine {
	httpEngine := gin.Default()
	api := api2go.NewAPIWithRouting(
		"",
		api2go.NewStaticResolver("/"),
		gingonic.New(httpEngine),
	)
	spacesResource := resource.NewSpacesResource(appDB, config)
	api.AddResource(space.Space{}, &spacesResource)
	httpEngine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return httpEngine
}
