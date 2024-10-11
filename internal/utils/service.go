package utils

// GetPrometheusLabels returns a map of labels used for prometheus scraping
func GetPrometheusLabels() map[string]string {
	return map[string]string{
		"prometheus.io/port":   "9091",
		"prometheus.io/scrape": "true",
	}
}
