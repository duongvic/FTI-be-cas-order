package wsgi

import (
	"casorder/middlewares"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"casorder/api"
	"casorder/db"
	"casorder/utils/logging"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("Preflight Request Received")
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Header("Location", c.Request.Header.Get("Origin"))
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Initialize() {

	flag.Parse()
	port := viper.GetString("server.port")
	ip := viper.GetString("server.ip")

	app := gin.Default() // create gin app
	app.Use(CORSMiddleware())
	app.Use(location.Default())
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	gin.DisableConsoleColor()

	// Logging to a file.
	logPath := fmt.Sprintf("%v%v", viper.GetString("logging.logFolder"), "api.log")
	f, _ := os.Create(logPath)
	gin.DefaultWriter = io.MultiWriter(f)
	app.Use(gin.Recovery())
	app.Use(middlewares.AuthMiddleware())
	app.Use(db.Inject(db.GetDB()))
	app.Use(logging.Inject(logging.GetLogger()))
	api.ApplyRoutes(app) // apply api router
	var serverAddr = fmt.Sprintf("%v:%v", ip, port)
	app.Run(serverAddr) // listen to given port
}
