package report

import (

	//"hermesworld.com/grafana-reporter/grafana-reporter/grafanaaccess"
	//"hermesworld.com/grafana-reporter/grafana-reporter/schema"

	log "github.com/sirupsen/logrus"
	"github.com/nolte/go-grafana-exporter/pkg/config"
	"github.com/nolte/go-grafana-exporter/pkg/export"

	"github.com/spf13/cobra"
)

//var bulkExporter grafanaaccess.ExportBulk
//
var exportTimeFrom string
var exportTimeEnd string

var timeRange export.ExportTimeRange

var exportConfig config.BulkExportConfiguration
var outputDirectory string

var bulkExportConfig string

var ExportBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Export a bulk of Panels from given config File.",
	Long: `
Export a Set of Grafana Panels to a local Filesystem Structure, 
this can be helpful for sharing test results.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.Infof("Start Export a Bulk of Panels to %s", outputDirectory)

		//exporter := export.BulkExporter{}
		//exporter.New(exportConfig)
		//exporter.Export()

		err := exportConfig.Load(bulkExportConfig)
		if err != nil {
			log.Fatalf("Faild to load the Bulk Export Config: %v", err)
		}

		export.PathHelper = export.ReportPathPlainHelper{}

		result := export.ExportPanelsToLocalFileSystem(outputDirectory, exportConfig, timeRange)
		log.Debugf("Exported Groups %v", len(result.Results))

		//var exporter export.GrafanaExporter
		//err := exporter.New(viper.GetString("grafana.url"), viper.GetString("token"), &config, timeRange)
		//if err != nil {
		//	log.Fatalf("Fail to create Exporter", err)
		//}
		//exporter.Export(outputDirectory)

	},
}

func init() {
	log.Debug("init Export a Bulk of Panels to")
	ExportBulkCmd.PersistentFlags().StringVar(&bulkExportConfig, "exportConfig", "./export_description.yml", "Bulk Export Configuration")
	ExportBulkCmd.PersistentFlags().StringVar(&exportTimeFrom, "from", "1568458539", "Graph Start Time")
	ExportBulkCmd.PersistentFlags().StringVar(&exportTimeEnd, "end", "1568480139", "Graph End Time")
	ExportBulkCmd.PersistentFlags().StringVar(&outputDirectory, "export", "/tmp/grafana_export", "Export Target Folder")

	startTime, err := export.ConvertTimeToDate(exportTimeFrom)
	if err != nil {
		log.Panicf("Fail to Parse %s %v", exportTimeFrom, err)
	}

	endTime, err := export.ConvertTimeToDate(exportTimeEnd)
	if err != nil {
		log.Panicf("Fail to Parse %s %v", exportTimeEnd, err)
	}
	timeRange.New(startTime, endTime)
}
