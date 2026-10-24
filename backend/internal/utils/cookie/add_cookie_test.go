package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCookiesAreSameSiteLax(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		add  func(c *gin.Context)
	}{
		{name: "access token", add: func(c *gin.Context) { AddAccessTokenCookie(c, 60, "token") }},
		{name: "session id", add: func(c *gin.Context) { AddSessionIdCookie(c, 60, "session") }},
		{name: "device token", add: func(c *gin.Context) { AddDeviceTokenCookie(c, "device") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.add(c)

			cookies := w.Result().Cookies()
			require.Len(t, cookies, 1)
			assert.Equal(t, http.SameSiteLaxMode, cookies[0].SameSite)
			assert.True(t, cookies[0].Secure)
			assert.True(t, cookies[0].HttpOnly)
		})
	}
}
