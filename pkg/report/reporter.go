package report

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/nolte/go-grafana-exporter/pkg/export"
)

var TemplateBase string
var TemplateSubdir string

type Report struct {
	ReportName     string
	ReportTemplate string
	ReportContent  export.BulkExportResult

	TimeRange   export.ExportTimeRange
	Stage       string
	Description string
}

func CreateReport(outputDir string, report *Report) {
	os.MkdirAll(outputDir, os.ModePerm)
	// Generate a Report from Given Templates
	var files []string

	err := filepath.Walk(path.Join(TemplateBase, "commons"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	err = filepath.Walk(path.Join(TemplateBase, TemplateSubdir), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	log.Info(files)

	t := template.Must(template.New("").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}).ParseFiles(files[:]...))

	reportFile := path.Join(outputDir, report.ReportName+".md")

	log.Debugf("Write the Report to %s", reportFile)
	f, err := os.Create(reportFile)
	defer f.Close()

	err = t.ExecuteTemplate(f, report.ReportTemplate, report)

	//log.Info(fmt.Println(buf.String()))
	//log.Info(output)
	log.Infof("Successfull exported to %s", outputDir)

}
