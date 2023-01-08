// Methods to interface with a Smart Hub's HTTP endpoint
package smarthub

import (
	"io"
	"net/http"

	"github.com/jamesorlakin/smarthub/lan"
	"github.com/jamesorlakin/smarthub/wan"
	log "github.com/sirupsen/logrus"
)

func ScrapeWanDetails(host string) (*wan.WanConnectionDetails, error) {
	var err error
	defer func() {
		if err != nil {
			log.Errorf("Error requesting scrape from Smart Hub: %v", err)
		}
	}()

	// It seems this doesn't need any auth!
	// The only bit that seems to matter is the Referer header!
	url := "http://" + host + "/nonAuth/wan_conn.xml"
	log.Debugf("Building request to %v", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", "http://"+host+"/")
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	connDetails, err := wan.ParseWanXml(respBytes)
	if err != nil {
		return nil, err
	}

	return connDetails, nil
}

func ScrapeLanDetails(host string) ([]lan.Device, error) {
	var err error
	defer func() {
		if err != nil {
			log.Errorf("Error requesting scrape from Smart Hub: %v", err)
		}
	}()

	// It seems this doesn't need any auth!
	// The only bit that seems to matter is the Referer header!
	url := "http://" + host + "/cgi/cgi_basicMyDevice.js"
	log.Debugf("Building request to %v", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", "http://"+host+"/")
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lanDetails, err := lan.ParseLanDevices(respBytes)
	if err != nil {
		return nil, err
	}

	return lanDetails, nil
}
