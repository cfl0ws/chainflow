# stargazer_exporter

This Prometheus exporter is used to check the missed blocks of a given block chain address via the [Stargazer API](https://app.swaggerhub.com/apis-docs/Slamper/Stargazer/1.0.0#/Validators/GetValidatorMissesGrouped).

In order to run this you need to have the hash address.

## Build

```
go get github.com/thilinapiy/stargazer_exporter
cd $GOPATH/src/github.com/thilinapiy/stargazer_exporter
go build
```

## Install

**[Important]** Update the `stargazer_exporter.service` file's `--block-address` replacing `<change to your hash address>` with your hash address.

```
sudo chown prometheus:prometheus /usr/local/bin/stargazer_exporter
sudo chmod 755 /usr/local/bin/stargazer_exporter

sudo cp stargazer_exporter.service /etc/systemd/system/stargazer_exporter.service
sudo systemctl daemon-reload
sudo systemctl start stargazer_exporter.service
sudo systemctl enable stargazer_exporter.service
```

## Test

Run a curl on the end-point and it should give metrics output.

```
curl localhost:9119/metrics
```

## Integrate with Prometheus

Update `prometheus.yml` with following scrape configs and restart it.

```
scrape_configs:
  - job_name: 'stargazer_exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:9119'] 

```
