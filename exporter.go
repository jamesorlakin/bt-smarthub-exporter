package main

import (
	"net/http"
	"os"

	"github.com/jamesorlakin/smarthub/smarthub"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	metricUp = prometheus.NewDesc(
		"up",
		"Whether the scrape target is accessible",
		nil, nil,
	)
	metricIsConnected = prometheus.NewDesc(
		"smarthub_connected",
		"Whether the WAN link of the Smart Hub is connected",
		nil, nil,
	)
	metricDownloadedBytes = prometheus.NewDesc(
		"smarthub_downloaded_bytes",
		"How many bytes have been downloaded during the lifetime of the WAN connection",
		nil, nil,
	)
	metricUploadedBytes = prometheus.NewDesc(
		"smarthub_uploaded_bytes",
		"How many bytes have been uploaded during the lifetime of the WAN connection",
		nil, nil,
	)
	metricUploadRate = prometheus.NewDesc(
		"smarthub_upload_rate",
		"The WAN connection speed for upload (may be to the modem for FTTP, which itself has a different speed)",
		nil, nil,
	)
	metricDownloadRate = prometheus.NewDesc(
		"smarthub_download_rate",
		"The WAN connection speed for download (may be to the modem for FTTP, which itself has a different speed)",
		nil, nil,
	)
)

type Exporter struct {
	host string
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricUp
	ch <- metricIsConnected
	ch <- metricDownloadedBytes
	ch <- metricUploadedBytes
	ch <- metricUploadRate
	ch <- metricDownloadRate
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Debug("Scraping")
	defer log.Debug("Scrape finished")
	details, err := smarthub.ScrapeWanDetails(e.host)

	if err != nil {
		ch <- prometheus.MustNewConstMetric(metricUp, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(metricUp, prometheus.GaugeValue, 1)

		if details.IsConnected {
			ch <- prometheus.MustNewConstMetric(metricIsConnected, prometheus.GaugeValue, 1)
		} else {
			ch <- prometheus.MustNewConstMetric(metricIsConnected, prometheus.GaugeValue, 0)
		}
		ch <- prometheus.MustNewConstMetric(metricDownloadedBytes, prometheus.CounterValue, float64(details.DownloadedBytes))
		ch <- prometheus.MustNewConstMetric(metricUploadedBytes, prometheus.CounterValue, float64(details.UploadedBytes))
		ch <- prometheus.MustNewConstMetric(metricUploadRate, prometheus.CounterValue, float64(details.UploadRateBps))
		ch <- prometheus.MustNewConstMetric(metricDownloadRate, prometheus.CounterValue, float64(details.DownloadRateBps))
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Smart Hub exporter starting")

	host := os.Getenv("SMARTHUB_HOST")
	if host == "" {
		panic("SMARTHUB_HOST environment variable not set")
	}

	prometheus.MustRegister(&Exporter{
		host: host,
	})
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9101", nil)
}
