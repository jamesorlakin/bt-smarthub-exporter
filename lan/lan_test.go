package lan

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	file, err := os.ReadFile("lan.js")
	if err != nil {
		panic("Could not read test file: " + err.Error())
	}

	devices, _ := ParseLanDevices(file)

	assert.Len(t, devices, 2)

	assert.Equal(t, devices[0].Mac, "E4:F0:42:81:F6:30", "Device 1 MAC")
	assert.Equal(t, devices[0].Hostname, "Chromecast", "Device 1 MAC")
	assert.Equal(t, devices[0].Ip, "192.168.0.71", "Device 1 IP")
	assert.Equal(t, devices[0].DownloadedBytes, 0, "Device 1 Downloaded")
	assert.Equal(t, devices[0].UploadedBytes, 0, "Device 1 Uploaded")

	assert.Equal(t, devices[1].Mac, "F8:32:E4:9E:A4:89", "Device 1 MAC")
	assert.Equal(t, devices[1].Hostname, "jl-quail-1", "Device 1 MAC")
	assert.Equal(t, devices[1].Ip, "192.168.0.107", "Device 1 IP")
	assert.Equal(t, devices[1].DownloadedBytes, 8765641, "Device 1 Downloaded")
	assert.Equal(t, devices[1].UploadedBytes, 1313609, "Device 1 Uploaded")
}

func Benchmark(b *testing.B) {
	file, err := os.ReadFile("lan.js")
	if err != nil {
		panic("Could not read test file: " + err.Error())
	}
	for i := 0; i < b.N; i++ {
		ParseLanDevices(file)
	}
}
