package enedis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ViBiOh/httputils/v2/pkg/errors"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
)

func (a *app) appendSessionCookie(headers http.Header) {
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
