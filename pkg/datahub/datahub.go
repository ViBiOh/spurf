package datahub

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

var (
	productionTokenURL = "https://gw.prd.api.enedis.fr"
	sandboxTokenURL    = "https://gw.hml.api.enedis.fr"
	isoDateLayout      = "2006-01-02"
)

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Reading defines a consumption value
type Reading struct {
	Value    string `json:"value"`
	Date     string `json:"date"`
	Interval string `json:"interval_length"`
}

// MasterReading defines a consumption aggregate
type MasterReading struct {
	Readings []Reading `json:"interval_reading"`
}

// Consumption defines a consumption response
type Consumption struct {
	Master MasterReading `json:"meter_reading"`
}

// App of package
type App interface {
	GetConsumption(context.Context, time.Time, time.Time) (Consumption, error)
	RefreshToken(context.Context) error
}

// Config of package
type Config struct {
	accessToken  *string
	refreshToken *string
	clientID     *string
	clientSecret *string
	redirectURI  *string
	usagePointID *string
	sandbox      *bool
}

type app struct {
	accessToken  string
	refreshToken string
	clientID     string
	clientSecret string
	redirectURI  string
	usagePointID string
	sandbox      bool

	mutex sync.RWMutex
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string, overrides ...flags.Override) Config {
	return Config{
		accessToken:  flags.New(prefix, "datahub").Name("AccessToken").Default("").Label("Access Token").ToString(fs),
		refreshToken: flags.New(prefix, "datahub").Name("RefreshToken").Default("").Label("Refresh Token").ToString(fs),
		clientID:     flags.New(prefix, "datahub").Name("ClientID").Default("").Label("Client ID").ToString(fs),
		clientSecret: flags.New(prefix, "datahub").Name("ClientSecret").Default("").Label("Client Secret").ToString(fs),

		redirectURI:  flags.New(prefix, "datahub").Name("RedirectUri").Default("https://api.vibioh.fr/dump/").Label("Redirect URI").ToString(fs),
		usagePointID: flags.New(prefix, "datahub").Name("UsagePointId").Default("").Label("Identifiant du point de livraison").ToString(fs),

		sandbox: flags.New(prefix, "datahub").Name("Sandbox").Default(false).Label("Sandbox mode").ToBool(fs),
	}
}

// New creates new App from Config
func New(config Config) App {
	return &app{
		accessToken:  strings.TrimSpace(*config.accessToken),
		refreshToken: strings.TrimSpace(*config.refreshToken),
		clientID:     strings.TrimSpace(*config.clientID),
		clientSecret: strings.TrimSpace(*config.clientSecret),
		redirectURI:  strings.TrimSpace(*config.redirectURI),
		usagePointID: strings.TrimSpace(*config.usagePointID),
		sandbox:      *config.sandbox,
	}
}

func (a *app) GetURL() string {
	if a.sandbox {
		return sandboxTokenURL
	}
	return productionTokenURL
}

func (a *app) prepareRequest(lastInsert time.Time, now time.Time) *request.Request {
	req := request.New()

	a.mutex.RLock()
	req.Header("Authorization", fmt.Sprintf("Bearer %s", a.accessToken))
	a.mutex.RUnlock()

	usedLast := lastInsert.Add(time.Hour * 24)
	usedNow := now

	if (usedNow.Sub(usedLast)) > (time.Hour * 24 * 6) {
		usedNow = usedLast.Add(time.Hour * 24 * 6)
	}

	if usedNow.After(lastInsert) {
		return nil
	}

	req.Get(fmt.Sprintf("%s/v4/metering_data/consumption_load_curve?usage_point_id=%s&start=%s&end=%s", a.GetURL(), url.QueryEscape(a.usagePointID), usedLast.Format(isoDateLayout), usedNow.Format(isoDateLayout)))

	return req
}

func (a *app) GetConsumption(ctx context.Context, lastInsert time.Time, now time.Time) (Consumption, error) {
	var payload Consumption

	req := a.prepareRequest(lastInsert, now)
	if req == nil {
		return payload, nil
	}

	resp, err := req.Send(ctx, nil)
	if err != nil && resp.StatusCode == http.StatusForbidden {
		if err := a.RefreshToken(ctx); err != nil {
			return payload, fmt.Errorf("unable to refresh token: %s", err)
		}

		resp, err = a.prepareRequest(lastInsert, now).Send(ctx, nil)
	}

	if err != nil {
		return payload, fmt.Errorf("unable to get data: %s", err)
	}

	body, err := request.ReadBodyResponse(resp)
	if err != nil {
		return payload, fmt.Errorf("unable to read data response: %s", err)
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return payload, fmt.Errorf("unable to unmarshal data: %s", err)
	}

	return payload, nil
}

func (a *app) RefreshToken(ctx context.Context) error {
	req := request.New()

	req.Post(fmt.Sprintf("%s/v1/oauth2/token?redirect_uri=%s", a.GetURL(), url.QueryEscape(a.redirectURI)))

	values := url.Values{}
	values.Set("client_id", a.clientID)
	values.Set("client_secret", a.clientSecret)
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", a.refreshToken)

	resp, err := req.Form(ctx, values)
	if err != nil {
		return fmt.Errorf("unable to refresh token: %s", err)
	}

	body, err := request.ReadBodyResponse(resp)
	if err != nil {
		return fmt.Errorf("unable to read token response: %s", err)
	}

	var payload token
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unable to unmarshal token: %s", err)
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.accessToken = payload.AccessToken
	a.refreshToken = payload.RefreshToken

	return nil
}
