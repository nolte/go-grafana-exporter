package report

import (
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/nolte/go-grafana-exporter/pkg/export"
	"github.com/nolte/go-grafana-exporter/pkg/report"

	"github.com/spf13/cobra"
)

var internaleHugoPath string
var internaleHugoGroup string

var ExportBulkReportHugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Create a Hugo formatted Report",
	Long: `
Export a set of Panels from given FileSystem config, to a hugo blog entry. 
This can be used for archiving TestResults and create a PDF Exportable Report.
`,
	Run: func(cmd *cobra.Command, args []string) {

		err := exportConfig.Load(bulkExportConfig)
		if err != nil {
			log.Fatalf("Faild to load the Bulk Export Config: %v", err)
		}

		contentOutput := path.Join(outputDirectory, "content")
		contentOutput = path.Join(contentOutput, internaleHugoPath)
		contentOutput = path.Join(contentOutput, internaleHugoGroup)

		imageOutput := path.Join(outputDirectory, "static/img")
		imageOutputRun := path.Join(imageOutput, internaleHugoPath)
		imageOutputRun = path.Join(imageOutputRun, internaleHugoGroup)
		imageOutputRun = path.Join(imageOutputRun, bulkReport.ReportName)

		var helper export.ReportReplacePathHelper
		helper.Replacement = "/img"
		helper.OutputBase = imageOutput
		export.PathHelper = helper
		bulkReport.TimeRange = timeRange

		result := export.ExportPanelsToLocalFileSystem(imageOutputRun, exportConfig, timeRange)
		bulkReport.ReportContent = result
		bulkReport.ReportTemplate = "REPORT"
		report.TemplateSubdir = "hugo"

		report.CreateReport(contentOutput, &bulkReport)
		log.Debugf("Exported Groups %v", len(result.Results))

	},
}

func init() {
	ExportBulkReportCmd.AddCommand(ExportBulkReportHugoCmd)
	ExportBulkReportHugoCmd.PersistentFlags().StringVar(&internaleHugoPath, "internal", "/testreports/", "Internal Report Base Path")
	ExportBulkReportHugoCmd.PersistentFlags().StringVar(&internaleHugoGroup, "internalGroup", "loadtest", "Internal Hugo Test Group")
	//
	//reportRun.TimeRange = timeRange
}
