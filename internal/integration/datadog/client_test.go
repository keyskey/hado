package datadog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMonitors_decode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/monitor" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("DD-API-KEY") != "api" || r.Header.Get("DD-APPLICATION-KEY") != "app" {
			t.Fatalf("missing auth headers")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":42,"name":"test","type":"query alert","query":"","tags":["env:prod","hado:foo"]}]`))
	}))
	defer ts.Close()

	c, err := NewClient(Config{
		APIKey:     "api",
		AppKey:     "app",
		Site:       "datadoghq.com",
		APIBaseURL: ts.URL,
	}, ts.Client())
	if err != nil {
		t.Fatal(err)
	}

	mons, err := c.ListMonitors(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(mons) != 1 || mons[0].ID != 42 || mons[0].Name != "test" {
		t.Fatalf("got %+v", mons)
	}
}

func TestMonitorAppURL_defaultSite(t *testing.T) {
	got := MonitorAppURL("", 7)
	want := "https://app.datadoghq.com/monitors/7"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestMonitorAppURL_eu(t *testing.T) {
	got := MonitorAppURL("datadoghq.eu", 9)
	want := "https://app.datadoghq.eu/monitors/9"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
