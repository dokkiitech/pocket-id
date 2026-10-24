package cookie

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddAccessTokenCookie(c *gin.Context, maxAgeInSeconds int, token string) {
	setCookie(c, AccessTokenCookieName, token, maxAgeInSeconds, "/")
}

func AddSessionIdCookie(c *gin.Context, maxAgeInSeconds int, sessionID string) {
	setCookie(c, SessionIdCookieName, sessionID, maxAgeInSeconds, "/")
}

func AddDeviceTokenCookie(c *gin.Context, deviceToken string) {
	setCookie(c, DeviceTokenCookieName, deviceToken, int(15*time.Minute.Seconds()), "/api/one-time-access-token")
}

// setCookie sets a secure, HTTP-only cookie with SameSite=Lax.
// SameSite=Lax prevents the cookie from being sent with cross-site POST requests, which protects against CSRF,
// while still allowing top-level GET navigations from relying parties (e.g. the OIDC authorization flow).
func setCookie(c *gin.Context, name string, value string, maxAgeInSeconds int, path string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAgeInSeconds, path, "", true, true)
}
