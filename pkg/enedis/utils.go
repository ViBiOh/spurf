package enedis

import (
	"net/http"
	"strings"
)

func (a *app) appendCookies(response *http.Response) {
	for _, cookie := range response.Cookies() {
		if strings.Contains(cookie.Domain, "enedis.fr") || cookie.Domain == "" {
			a.cookies = append(a.cookies, cookie)
		}
	}
}
