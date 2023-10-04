package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// IsDecryptor is a middleware that checks if
// the user has the 'Decryptor' role.
func IsDecryptor(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	p := profile.(map[string]interface{})
	roles := p["https://codec-server/roles"]
	r := roles.([]interface{})

	if slices.Contains(r, "Decryptor") {
		ctx.Next()
	} else {
		ctx.AbortWithStatus(403)
	}
}
