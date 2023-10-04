package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

// Handler for our logged-in user page.
func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	// Check if user has the Decryptor role, and set in the template data.
	p := profile.(map[string]interface{})
	roles := p["https://codec-server/roles"]
	r := roles.([]interface{})
	isDecryptor := slices.Contains(r, "Decryptor")

	data := gin.H{
		"profile":     profile,
		"isDecryptor": isDecryptor,
	}

	ctx.HTML(http.StatusOK, "user.html", data)
}
