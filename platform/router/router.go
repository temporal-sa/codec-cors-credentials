package router

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/converter"

	"codec-cors-credentials/platform/authenticator"
	"codec-cors-credentials/platform/encryption"
	"codec-cors-credentials/platform/middleware"
	"codec-cors-credentials/web/app/callback"
	"codec-cors-credentials/web/app/home"
	"codec-cors-credentials/web/app/login"
	"codec-cors-credentials/web/app/logout"
	"codec-cors-credentials/web/app/user"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	// The cookie store must not be restricted to the same site, otherwise the
	// browser will not send the cookie with the JS fetch request from the
	// Temporal Web UI. When using same site none, the cookie must also be secure.
	store.Options(sessions.Options{
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
	router.Use(sessions.Sessions("auth-session", store))

	// Configure CORS, so that requests from the Temporal Web UI are allowed by
	// this codec server.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{os.Getenv("TEMPORAL_ORIGIN_URL")}
	corsConfig.AllowMethods = []string{"POST"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "X-Namespace"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", home.Handler)
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/user", middleware.IsAuthenticated, user.Handler)
	router.GET("/logout", logout.Handler)

	// Register the coded handler, for the /decode endpoint. Access to the /decode
	// endpoint is restricted to authenticated users with the "Decryptor" role.
	codecHandler := converter.NewPayloadCodecHTTPHandler(&encryption.Codec{}, converter.NewZlibCodec(converter.ZlibCodecOptions{AlwaysEncode: true}))
	router.POST("/decode", middleware.IsAuthenticated, middleware.IsDecryptor, gin.WrapH(codecHandler))

	return router
}
