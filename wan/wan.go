// A collection of methods to parse the response of the Smart Hub XML
package wan

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// The raw response from /nonAuth/wan_conn.xml
// This XML is *super* weird. It's XML, but the value is multidimensional JSON array in attributes.
// Oh, and some of the inner values are then URL-encoded strings delimited by ';'. WTF.
// The interpretation of these is based on the JavaScript within /broadband.htm
type WanConnXmlResponse struct {
	// [0] connected/disconnected, [1,2] not sure
	Wan_conn_status_list XmlValueAttribute `xml:"wan_conn_status_list"`
	// 'volume' being traffic metrics. [0] not sure, [1] downloaded, [2] uploaded
	Wan_conn_volume_list XmlValueAttribute `xml:"wan_conn_volume_list"`
	// [0] up bits/sec, [1] down bits/sec, [2,3] not sure
	Wan_status_rate XmlValueAttribute `xml:"status_rate"`
	Sysuptime       XmlValueAttribute `xml:"sysuptime"`
}

// A workaround for Go not supporting getting the attribute during a chain
// Grabs the `value="[['stuff1;stuff2']]"` bit of an XML element.
type XmlValueAttribute struct {
	Value string `xml:"value,attr"`
}

type WanConnectionDetails struct {
	IsConnected     bool
	UptimeSeconds   int
	DownloadedBytes int
	UploadedBytes   int
	DownloadRateBps int
	UploadRateBps   int
}

func ParseWanXml(xmlContents []byte) (*WanConnectionDetails, error) {
	log.Tracef("Got XML: %v", string(xmlContents))
	body := WanConnXmlResponse{}
	err := xml.Unmarshal(xmlContents, &body)
	if err != nil {
		log.Errorf("Could not parse XML: %v", err)
		return nil, err
	}

	connectionDetails := &WanConnectionDetails{}

	uptime, err := strconv.Atoi(body.Sysuptime.Value)
	if err != nil {
		log.Errorf("Could not parse uptime int: %v", err)
		return nil, err
	}
	connectionDetails.UptimeSeconds = uptime

	connStatus, err := decodeNestedJsonArrayFirst(body.Wan_conn_status_list.Value)
	if err != nil {
		log.Errorf("Could not decode connection status element: %v", err)
		return nil, err
	}
	connectionDetails.IsConnected = connStatus[0] == "connected"

	connVolume, err := decodeNestedJsonIntArrayFirst(body.Wan_conn_volume_list.Value)
	if err != nil {
		log.Errorf("Could not decode connection volume element: %v", err)
		return nil, err
	}
	connectionDetails.DownloadedBytes = connVolume[1]
	connectionDetails.UploadedBytes = connVolume[2]

	connRate, err := decodeNestedJsonIntArrayFirst(body.Wan_status_rate.Value)
	if err != nil {
		log.Errorf("Could not decode connection rate element: %v", err)
		return nil, err
	}
	connectionDetails.UploadRateBps = connRate[0]
	connectionDetails.DownloadRateBps = connRate[1]

	return connectionDetails, nil
}

// decode the nested JSON array (which has single quotes?), taking only the first element and URL-decode it before splitting the first item via a semicolon.
// e.g. `[['123%3B456'], ...]` becomes {'123', '456'}
func decodeNestedJsonArrayFirst(value string) ([]string, error) {
	if value == "" {
		// This suggested the value wasn't from the XML
		return nil, errors.New("expected a JSON array value but was empty")
	}

	parsedValue := make([][]string, 0, 3)
	// The 'value' element arrays aren't strictly JSON, as they have single quotes which aren't in the spec.
	// This hackily replaces them...
	massagedJson := strings.ReplaceAll(value, `'`, `"`)
	err := json.Unmarshal([]byte(massagedJson), &parsedValue)
	if err != nil {
		return nil, err
	}

	firstElement := parsedValue[0][0]
	firstElementDecoded, err := url.PathUnescape(firstElement)
	if err != nil {
		return nil, err
	}
	return strings.Split(firstElementDecoded, ";"), nil
}

func decodeNestedJsonIntArrayFirst(value string) ([]int, error) {
	values, err := decodeNestedJsonArrayFirst(value)
	if err != nil {
		return nil, err
	}

	intValues := make([]int, 0, len(values))
	for _, v := range values {
		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		intValues = append(intValues, num)
	}

	return intValues, nil
}
