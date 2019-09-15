package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	gapi "github.com/nolte/go-grafana-api"
	"gopkg.in/yaml.v2"
)

type BulkExportConfiguration struct {
	Exports map[string]ExportGroup `yaml:"export"`
}

func (c *BulkExportConfiguration) Load(path string) error {
	log.Debugf("Load Export config from local file %s", path)
	// load the config file
	// var bulkConfig BulkExportConfiguration
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}
	log.Debugf("Unmarshal Config to local object structure %v", string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	log.Debugf("Unmarshal Config to local object structure %v", c)

	return nil
}

type ExportGroup struct {
	Title       string          `yaml:"title"`
	Elements    []ExportElement `yaml:"elements"`
	Description string          `yaml:"description"`
}

type ExportElement struct {
	Description string      `yaml:"description"`
	Title       string      `yaml:"title"`
	Panel       ExportPanel `yaml:"panel"`
	ExportSize  ExportSize  `yaml:"size"`
}

type ExportPanel struct {
	ID int64 `yaml:"id"`

	OrgID       int64  `yaml:"orgId"`
	DashboardID string `yaml:"dashboardId"`
	ExportName  string `yaml:"exportName"`

	DashboardVars map[string][]string `yaml:"vars"`
}

type ExportSize struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

func (e ExportSize) AsApiObject() gapi.GrafanaPanelExportSize {
	apiObj := gapi.GrafanaPanelExportSize{}
	apiObj.New(e.Width, e.Height)
	return apiObj
}
