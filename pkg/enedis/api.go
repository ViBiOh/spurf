package enedis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v2/pkg/errors"
	"github.com/ViBiOh/httputils/v2/pkg/request"
)

const (
	oneDay = 24 * time.Hour
)

func (a *app) login() error {
	if a.email == "" || a.password == "" {
		return errors.New("no credentials provided")
	}

	values := url.Values{}
	values.Add("IDToken1", a.email)
	values.Add("IDToken2", a.password)
	values.Add("SunQueryParamsString", "cmVhbG09cGFydGljdWxpZXJz")
	values.Add("encoded", "true")
	values.Add("gx_charset", "UTF-8")

	ctx := context.Background()
	_, _, headers, err := request.PostForm(ctx, loginURL, values, nil)
	if err != nil {
		return err
	}

	authCookies := strings.Builder{}
	for _, cookie := range headers["Set-Cookie"] {
		if !strings.Contains(cookie, "Domain=.enedis.fr") {
			continue
		}

		if strings.Contains(cookie, "Expires=Thu, 01-Jan-1970 00:00:10 GMT") {
			continue
		}

		if authCookies.Len() != 0 {
			safeWrite(&authCookies, ";")
		}
		safeWrite(&authCookies, getCookieValue(cookie))
	}

	a.cookie = authCookies.String()

	return nil
}

func (a *app) getData(ctx context.Context, startDate string, first bool) (*Consumption, error) {
	header := http.Header{}
	header.Set("Cookie", a.cookie)

	params := url.Values{}
	params.Add("p_p_id", "lincspartdisplaycdc_WAR_lincspartcdcportlet")
	params.Add("p_p_lifecycle", "2")
	params.Add("p_p_state", "normal")
	params.Add("p_p_mode", "view")
	params.Add("p_p_resource_id", "urlCdcHeure")
	params.Add("p_p_cacheability", "cacheLevelPage")
	params.Add("p_p_col_id", "column-1")
	params.Add("p_p_col_count", "2")

	startTime, err := time.ParseInLocation(frenchDateFormat, startDate, a.location)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	endDate := startTime.Add(oneDay).Format(frenchDateFormat)

	values := url.Values{}
	params.Add("_lincspartdisplaycdc_WAR_lincspartcdcportlet_dateDebut", startDate)
	params.Add("_lincspartdisplaycdc_WAR_lincspartcdcportlet_dateFin", endDate)

	body, status, headers, err := request.PostForm(ctx, fmt.Sprintf("%s%s", consumeURL, params.Encode()), values, header)
	if err != nil || status == http.StatusFound {
		if first {
			a.appendSessionCookie(headers)
			return a.getData(ctx, startDate, false)
		}

		if err == nil {
			return nil, errors.New("unable to authent to enedis on the second try")
		}

		return nil, err
	}

	payload, err := request.ReadBody(body)
	if err != nil {
		return nil, err
	}

	var response Consumption
	if err := json.Unmarshal(payload, &response); err != nil {
		return nil, errors.WithStack(err)
	}

	if response.Etat != nil && response.Etat.Valeur == "erreur" {
		return nil, errors.New("API error: %s", response.Etat.ErreurText)
	}

	for _, value := range response.Graphe.Data {
		value.Timestamp = startTime.Add(time.Duration(30*(value.Ordre-1)) * time.Minute).Unix()
	}

	return &response, nil
}
