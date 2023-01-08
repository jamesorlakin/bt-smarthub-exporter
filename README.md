# BT Smart Hub (2) Exporter

> *(It goes without saying but this definitely isn't anything official by BT!)*

This is a simple attempt at a Prometheus exporter for the BT Smart Hub 2 router.
Other variants of the Smart Hub router are completely untested but feel free to report if it works and open pull requests for fixes and/or additional metrics.

## Usage

Run the executable with the environment variable `SMARTHUB_HOST` set to the IP address of the router (typically `192.168.1.254`).
The exporter will then listen on port `9101` for the HTTP path `/metrics`.

### Docker/Kubernetes

The image tag is [`jamesorlakin/smarthub-exporter:v0.1.0`](https://hub.docker.com)

An example Deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: smarthub-exporter
  namespace: monitoring
  labels:
    app: smarthub-exporter
spec:
  selector:
    matchLabels:
      app: smarthub-exporter
  replicas: 1
  template:
    metadata:
      labels:
        app: smarthub-exporter
      annotations:
        prometheus.io/port: "9101"
        prometheus.io/scrape: "true"
        prometheus.io/path: "/metrics"
    spec:
      containers:
        - name: exporter
          image: jamesorlakin/smarthub-exporter:v0.1.0
          livenessProbe:
            periodSeconds: 10
            timeoutSeconds: 5
            httpGet:
              path: /metrics
              port: 9101
          ports:
            - containerPort: 9101
              name: metrics
              protocol: TCP
          env:
            - name: SMARTHUB_HOST
              value: "192.168.1.254"
          resources:
            limits:
              cpu: 50m
              memory: 50Mi
```

## Metrics

The merics are those provided by the browser-based web console. These are reverse engineered from the browser dev tools and not the BT app - the consequence of which means to grab LAN device data this exporter executes returned JavaScript from the Smart Hub using the embedded Otto library.
Yes, *executes* JavaScript. There'll be a better way - please feel free to tweak!

| Metrics                       | Description                                                                                                                                                                          |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| up                            | Whether the scrape target is accessible (set to 0 if parsing/connecting is unsuccessful)                                                                                             |
| smarthub_connected            | Whether the WAN link of the Smart Hub is connected (this was tested against an FTTP-based installation using Ethernet port 4 - ADSL should work but this could do with verification) |
| smarthub_uptime               | How many seconds the system has been up for                                                                                                                                          |
| smarthub_connection_uptime    | How many seconds the WAN connection has been up for                                                                                                                                  |
| smarthub_downloaded_bytes     | How many bytes have been received over the lifetime of the WAN connection (or is it reboot of the router?)                                                                           |
| smarthub_uploaded_bytes       | How many bytes have been sent over the lifetime of the WAN connection (or is it reboot of the router?)                                                                               |
| smarthub_upload_rate          | The synchronized connection speed for upload (for FTTP this is likely the Gigabit connection to the modem, not the subscribed speed)                                                 |
| smarthub_download_rate        | The synchronized connection speed for download (for FTTP this is likely the Gigabit connection to the modem, not the subscribed speed)                                               |
| smarthub_lan_downloaded_bytes | The number of bytes downloaded from the internet by each LAN device since the router restarted. Labels include `mac`, `hostname`, and `ip`.                                          |
| smarthub_lan_uploaded_bytes   | The number of bytes downloaded to the internet by each LAN device since the router restarted. Labels include `mac`, `hostname`, and `ip`.                                            |

## Grafana Dashboard

An attempt at a Grafana dashboard lives in [Grafana Dashboard.json](./Grafana%20Dashboard.json):

![Dashboard screenshot](./Grafana%20Dashboard.png)

## TODO

- [ ] Option to turn off LAN metrics (e.g. executing JS is undesireable)
- [ ] More testing
- [ ] ARM64 build
- [ ] More metrics(?)
