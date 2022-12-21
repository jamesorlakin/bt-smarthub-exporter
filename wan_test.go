package wan

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getWanXml() []byte {
	file, err := os.ReadFile("smarthub_wan_conn.xml")
	if err != nil {
		panic("Could not read test file: " + err.Error())
	}
	return file
}

func TestXmlContents(t *testing.T) {
	file := getWanXml()
	connectionDetails, err := ParseWanXml(file)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
		return
	}

	assert.Equal(t, true, connectionDetails.IsConnected, "IsConnected")
	assert.Equal(t, 336463954090, connectionDetails.DownloadedBytes, "DownloadedBytes")
	assert.Equal(t, 34683717318, connectionDetails.UploadedBytes, "UploadedBytes")
	assert.Equal(t, 1000000000, connectionDetails.UploadRateBps, "UploadRateBps")
	assert.Equal(t, 1000000000, connectionDetails.DownloadRateBps, "DownloadRateBps")
}
