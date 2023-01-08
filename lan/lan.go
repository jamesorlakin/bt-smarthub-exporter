package lan

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/robertkrimen/otto"
	log "github.com/sirupsen/logrus"
)

type Device struct {
	Mac             string
	Hostname        string
	Ip              string
	DownloadedBytes int
	UploadedBytes   int
}

// Evaluate the JavaScript LAN stats and extract device metrics from them
// Yes, the router generates JavaScript instead of JSON...
func ParseLanDevices(js []byte) ([]Device, error) {
	vm := otto.New()
	// Add a no-op to the context
	vm.Run(`function addCfg() {}`)
	// TODO: Halting?
	vm.Run(js)

	devicesRaw, err := vm.Get("known_device_list")
	if err != nil {
		log.Errorf("Could not get known_device_list from JS context: %v", err)
		return nil, err
	}
	devices := jsArrayToMacStringMap(&devicesRaw)
	log.Tracef("Got LAN devices map: %v", devices)

	rateRaw, err := vm.Get("rate")
	if err != nil {
		log.Errorf("Could not get rate from JS context: %v", err)
		return nil, err
	}
	rates := jsArrayToMacStringMap(&rateRaw)
	log.Tracef("Got LAN rate map: %v", rates)

	devicesSlice := make([]Device, 0, len(devices))
	for mac, v := range devices {
		device := mapToDevice(v, rates[mac])
		devicesSlice = append(devicesSlice, device)
	}

	return devicesSlice, nil
}

// Convert a Otto array value to a native Go map by `mac` field. Nil if error.
// The strings are URL encoded so this deals with that too.
func jsArrayToMacStringMap(jsValue *otto.Value) map[string]map[string]string {
	// should be of type []interface{}
	// which contains map[string]interface{} (and that interface is always string)
	value, _ := jsValue.Export()
	valueSlice, ok := value.([]interface{})
	if !ok {
		log.Errorf("%v", ok)
	}

	// Convert from []interface{} to map[string]map[string]string
	//                                   ^^ mac      ^^ JS object keys
	macValueMap := make(map[string]map[string]string)
	for _, item := range valueSlice {
		itemMapProcessed := make(map[string]string)
		// The last element might be literally null, for some reason.
		if item == nil {
			continue
		}
		itemMap := item.(map[string]interface{})

		mac := itemMap["mac"].(string)
		mac, _ = url.PathUnescape(mac)
		// The MAC in the rate object is lowercase compared to others in upper...
		mac = strings.ToUpper(mac)

		for k, v := range itemMap {
			decodedString, _ := url.PathUnescape(v.(string))
			itemMapProcessed[k] = decodedString
		}
		macValueMap[mac] = itemMapProcessed
	}

	log.Tracef("%v", macValueMap)
	return macValueMap
}

func mapToDevice(deviceObj, rateObj map[string]string) Device {
	// rateObj may be nil
	device := Device{
		Mac:      deviceObj["mac"],
		Hostname: deviceObj["hostname"],
		Ip:       deviceObj["ip"],
	}

	if rateObj != nil {
		device.UploadedBytes, _ = strconv.Atoi(rateObj["tx"])
		device.DownloadedBytes, _ = strconv.Atoi(rateObj["rx"])
	}

	return device
}
