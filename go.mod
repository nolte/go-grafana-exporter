module github.com/nolte/grafana-exporter

go 1.13

require (
	github.com/nolte/go-grafana-api v0.0.1
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/nolte/go-grafana-api => /home/nolte/gostuff/src/github.com/nolte/go-grafana-api
