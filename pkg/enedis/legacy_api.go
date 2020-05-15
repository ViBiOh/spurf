package enedis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/request"
)

const (
	loginURL         = "https://espace-client-connexion.enedis.fr/auth/UI/Login"
	legacyConsumeURL = "https://espace-client-particuliers.enedis.fr/group/espace-particuliers/suivi-de-consommation?"
	oneDay           = 24 * time.Hour
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
	resp, err := request.New().Post(loginURL).Form(ctx, values)
	if err != nil {
		return err
	}

	a.cookies = make([]*http.Cookie, 0)
	a.appendCookies(resp)

	return nil
}

func (a *app) getDataFromLegacy(ctx context.Context, startDate string, first bool) (Consumption, error) {
	startTime, err := time.ParseInLocation(isoDateFormat, startDate, a.location)
	if err != nil {
		return emptyConsumption, err
	}

	params := url.Values{}
	params.Add("p_p_id", "lincspartdisplaycdc_WAR_lincspartcdcportlet")
	params.Add("p_p_lifecycle", "2")
	params.Add("p_p_state", "normal")
	params.Add("p_p_mode", "view")
	params.Add("p_p_resource_id", "urlCdcHeure")
	params.Add("p_p_cacheability", "cacheLevelPage")
	params.Add("p_p_col_id", "column-1")
	params.Add("p_p_col_pos", "1")
	params.Add("p_p_col_count", "3")

	values := url.Values{}
	values.Add("_lincspartdisplaycdc_WAR_lincspartcdcportlet_dateDebut", startTime.Format(frenchDateFormat))
	values.Add("_lincspartdisplaycdc_WAR_lincspartcdcportlet_dateFin", startTime.AddDate(0, 0, 1).Format(frenchDateFormat))

	req, err := request.New().Post(fmt.Sprintf("%s%s", legacyConsumeURL, params.Encode())).ContentForm().Build(ctx, strings.NewReader(values.Encode()))
	if err != nil {
		return emptyConsumption, err
	}

	for _, cookie := range a.cookies {
		req.AddCookie(cookie)
	}

	resp, err := request.Do(req)
	if err != nil || (resp != nil && resp.StatusCode == http.StatusFound) {
		if first {
			a.appendCookies(resp)
			return a.getDataFromLegacy(ctx, startDate, false)
		}

		if err == nil {
			return emptyConsumption, errors.New("unable to authent to enedis on the second try")
		}

		return emptyConsumption, err
	}

	payload, err := request.ReadBodyResponse(resp)
	if err != nil {
		return emptyConsumption, err
	}

	return parseResponseFromLegacy(startTime, payload)
}

func parseResponseFromLegacy(startTime time.Time, payload []byte) (Consumption, error) {
	var consumption Consumption
	if err := json.Unmarshal(payload, &consumption); err != nil {
		return emptyConsumption, err
	}

	if consumption.Etat.Valeur == "erreur" {
		return emptyConsumption, fmt.Errorf("API error: %s", consumption.Etat.ErreurText)
	}

	if consumption.Etat.Valeur == "nonActive" {
		return emptyConsumption, errors.New("Non active data")
	}

	datas := make([]Value, len(consumption.Graphe.Data))
	for index, value := range consumption.Graphe.Data {
		value.Timestamp = startTime.Add(time.Duration(30*(value.Ordre-1)) * time.Minute).Unix()
		datas[index] = value
	}

	consumption.Graphe.Data = datas

	return consumption, nil
}
