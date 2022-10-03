package app

import (
	"fmt"
	"net/http"

	subscriptionsGrpcApi "github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/api/grpc/v1/subscriptions"
	"github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/api/rest/v1/ping"
	subscriptionRestApi "github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/api/rest/v1/subscriptions"
	"github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/services"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/app"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/log"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/auth"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/utils"
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start() {
	app.LoadEnv()
	log.SetUpLogPath(utils.EnvVarDefault("APP_LOGS_PATH", "stdout"))
	creds := app.TLSCredentials()
	go func() {
		app.StartGRPC(setup, shutdown, app.HostGRPC(), createGrpcApi, &creds, log.Instance())
	}()
	app.StartHTTP(setup, shutdown, app.HostHTTP(), createRestApi(log.Instance()))
}

func setup() {
	services.Instance()
}

func shutdown() {
	err := services.Instance().Shutdown()
	log.Error("error during app shutdown", err.Error())
}

func createRestApi(logger *logrus.Logger) *gin.Engine {
	router := gin.Default()
	gin.SetMode(app.Mode())
	router.Use(app.Cors())
	router.Use(app.NewLoggerMiddleware(logger))
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	v1 := router.Group("/api/v1")

	v1.GET("/subscriptions/ping", ping.Ping)

	authorized := router.Group("/api/v1")
	authorized.Use(app.AuthReqired(authenicate))
	{
		authorized.GET("/subscriptions/debug/vars", app.RequiredOwnerRole(), expvar.Handler())
		authorized.GET("/subscriptions/safe-ping", app.RequiredOwnerRole(), ping.SafePing)

		authorized.POST("/subscriptions/event", app.RequiredOwnerRole(), subscriptionRestApi.AddEvent)
		authorized.POST("/subscriptions/event/email", app.RequiredOwnerRole(), subscriptionRestApi.AddSendEmailEvent)
	}
	return router
}

func createGrpcApi(s *grpc.Server) {
	subscriptionsGrpcApi.RegisterServiceServer(s)
}

func authenicate(token string) (*auth.VerificationResult, error) {
	return services.Instance().Auth().VerifyToken(token)
}
