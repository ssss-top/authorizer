package cookie

import (
	"net/url"

	"github.com/authorizerdev/authorizer/server/constants"
	"github.com/authorizerdev/authorizer/server/envstore"
	"github.com/authorizerdev/authorizer/server/utils"
	"github.com/gin-gonic/gin"
)

// SetAdminCookie sets the admin cookie in the response
func SetAdminCookie(gc *gin.Context, token string) {
	secure := false
	httpOnly := true
	hostname := utils.GetHost(gc)
	host, _ := utils.GetHostParts(hostname)

	gc.SetCookie(envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminCookieName), token, 3600, "/", host, secure, httpOnly)
}

// GetAdminCookie gets the admin cookie from the request
func GetAdminCookie(gc *gin.Context) (string, error) {
	cookie, err := gc.Request.Cookie(envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminCookieName))
	if err != nil {
		return "", err
	}

	// cookie escapes special characters like $
	// hence we need to unescape before comparing
	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}

// DeleteAdminCookie sets the response cookie to empty
func DeleteAdminCookie(gc *gin.Context) {
	secure := false
	httpOnly := true
	hostname := utils.GetHost(gc)
	host, _ := utils.GetHostParts(hostname)

	gc.SetCookie(envstore.EnvStoreObj.GetStringStoreEnvVariable(constants.EnvKeyAdminCookieName), "", -1, "/", host, secure, httpOnly)
}
