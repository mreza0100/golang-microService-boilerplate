package main

import (
	"fmt"
	"microServiceBoilerplate/configs"
	"microServiceBoilerplate/services/hellgate/graph/generated"
	"microServiceBoilerplate/services/hellgate/graph/resolvers"
	"microServiceBoilerplate/services/hellgate/middlewares"
	"microServiceBoilerplate/utils"

	"microServiceBoilerplate/proto/connections"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/mreza0100/golog"
)

func initLogger() *golog.Core {
	return golog.New(golog.InitOpns{
		LogPath:   configs.LogPath,
		Name:      "🔥__Hellgate__🔥",
		WithTime:  true,
		DebugMode: utils.IsDevMode,
	})
}

func main() {
	lgr := initLogger()
	graphqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(*resolvers.New(resolvers.NewOpts{
		Lgr: lgr,

		SecurityConn: connections.SecuritySrvConn(lgr),
		RelationConn: connections.RelationSrvConn(lgr),
		FeedConn:     connections.FeedSrvConn(lgr),
		UserConn:     connections.UserSrvConn(lgr),
		PostConn:     connections.PostSrvConn(lgr),
	})))

	ginServer := gin.New()
	ginServer.Use(gin.Recovery())

	ginServer.GET("/", func(ctx *gin.Context) {
		playground.Handler("micro", "/")(ctx.Writer, ctx.Request)
	})

	ginServer.Use(middlewares.SetWritter())
	ginServer.Use(middlewares.SetRequest())

	ginServer.POST("/", func(ctx *gin.Context) {
		graphqlServer.ServeHTTP(ctx.Writer, ctx.Request)
	})

	lgr.RedLog("🔥🔥🔥 Hellgate is open now 🔥🔥🔥")
	ginServer.Run(fmt.Sprintf(":%v", configs.HellgateConfigs.Port))
}
