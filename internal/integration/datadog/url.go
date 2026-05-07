package datadog

import "fmt"

// MonitorAppURL returns the web UI URL for a monitor id. site is the Datadog site suffix
// (for example datadoghq.com, datadoghq.eu) matching DD_SITE.
func MonitorAppURL(site string, monitorID int64) string {
	s := site
	if s == "" {
		s = defaultSite
	}
	return fmt.Sprintf("https://app.%s/monitors/%d", s, monitorID)
}
