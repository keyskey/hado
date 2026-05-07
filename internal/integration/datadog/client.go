package datadog

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	dd "github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
)

const defaultSite = "datadoghq.com"

// Config holds API authentication and site. Keys are never read from disk by this package.
type Config struct {
	APIKey string
	AppKey string
	Site   string
	// APIBaseURL, if non-empty, is a full base URL (e.g. httptest server) used to override
	// the API host via Configuration.Scheme/Host (for tests).
	APIBaseURL string
}

// Client lists Datadog monitors via the official datadog-api-client-go.
type Client struct {
	cfg      Config
	site     string
	monitors *datadogV1.MonitorsApi
}

// NewClient returns a Client with the given config and HTTP client.
func NewClient(cfg Config, httpClient *http.Client) (*Client, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, fmt.Errorf("datadog: api key is empty")
	}
	if strings.TrimSpace(cfg.AppKey) == "" {
		return nil, fmt.Errorf("datadog: application key is empty")
	}
	site := strings.TrimSpace(cfg.Site)
	if site == "" {
		site = defaultSite
	}

	configuration := dd.NewConfiguration()
	if httpClient != nil {
		configuration.HTTPClient = httpClient
	}
	if base := strings.TrimSpace(cfg.APIBaseURL); base != "" {
		parsed, err := url.Parse(base)
		if err != nil {
			return nil, fmt.Errorf("datadog: parse api base url: %w", err)
		}
		if parsed.Host != "" {
			if parsed.Scheme != "" {
				configuration.Scheme = parsed.Scheme
			} else {
				configuration.Scheme = "http"
			}
			configuration.Host = parsed.Host
		}
	}

	apiClient := dd.NewAPIClient(configuration)
	return &Client{
		cfg: Config{
			APIKey:     strings.TrimSpace(cfg.APIKey),
			AppKey:     strings.TrimSpace(cfg.AppKey),
			Site:       site,
			APIBaseURL: strings.TrimSpace(cfg.APIBaseURL),
		},
		site:     site,
		monitors: datadogV1.NewMonitorsApi(apiClient),
	}, nil
}

// NewClientFromEnv builds a client using DD_API_KEY, DD_APP_KEY, and optional DD_SITE.
func NewClientFromEnv(httpClient *http.Client) (*Client, error) {
	return NewClient(Config{
		APIKey: os.Getenv("DD_API_KEY"),
		AppKey: os.Getenv("DD_APP_KEY"),
		Site:   os.Getenv("DD_SITE"),
	}, httpClient)
}

// ListMonitors returns all monitors; call filters the result. Pagination is handled by the generated client.
func (c *Client) ListMonitors(ctx context.Context) ([]Monitor, error) {
	ctx = c.apiContext(ctx)
	ch, cancel := c.monitors.ListMonitorsWithPagination(ctx)
	defer cancel()

	var v1mons []datadogV1.Monitor
	for result := range ch {
		if result.Error != nil {
			return nil, formatListMonitorsError(result.Error)
		}
		v1mons = append(v1mons, result.Item)
	}

	out := make([]Monitor, 0, len(v1mons))
	for _, m := range v1mons {
		if m.Id == nil {
			continue
		}
		name := ""
		if m.Name != nil {
			name = *m.Name
		}
		tags := m.Tags
		if tags == nil {
			tags = []string{}
		}
		out = append(out, Monitor{
			ID:   *m.Id,
			Name: name,
			Tags: tags,
		})
	}
	return out, nil
}

// Site returns the normalized Datadog site suffix (for example datadoghq.com).
func (c *Client) Site() string {
	return c.cfg.Site
}

func (c *Client) apiContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, dd.ContextAPIKeys, map[string]dd.APIKey{
		"apiKeyAuth": {Key: c.cfg.APIKey},
		"appKeyAuth": {Key: c.cfg.AppKey},
	})
	// Test hosts override URL host via Configuration; default servers still need DD_SITE.
	if strings.TrimSpace(c.cfg.APIBaseURL) == "" {
		ctx = context.WithValue(ctx, dd.ContextServerVariables, map[string]string{
			"site": c.site,
		})
	}
	return ctx
}

func formatListMonitorsError(err error) error {
	var openAPIErr dd.GenericOpenAPIError
	if errors.As(err, &openAPIErr) {
		msg := strings.TrimSpace(string(openAPIErr.Body()))
		if msg == "" {
			return fmt.Errorf("datadog: list monitors: %s", openAPIErr.Error())
		}
		return fmt.Errorf("datadog: list monitors: %s: %s", openAPIErr.Error(), msg)
	}
	return fmt.Errorf("datadog: list monitors: %w", err)
}
