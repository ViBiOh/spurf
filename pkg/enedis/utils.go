package enedis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/logger"
)

func (a *App) appendSessionCookie(headers http.Header) {
	for _, cookie := range headers["Set-Cookie"] {
		if strings.HasPrefix(cookie, "JSESSIONID") {
			a.cookie = fmt.Sprintf("%s; %s", a.cookie, getCookieValue(cookie))
		}
	}
}

func safeWrite(w *strings.Builder, content string) {
	if _, err := w.WriteString(content); err != nil {
		logger.Error("%#v", errors.WithStack(err))
	}
}

func getCookieValue(cookie string) string {
	return strings.SplitN(cookie, ";", 2)[0]
}
