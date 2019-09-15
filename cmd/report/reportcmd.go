package report

import (
	log "github.com/sirupsen/logrus"
	"github.com/nolte/grafana-exporter/pkg/export"
	"github.com/nolte/grafana-exporter/pkg/report"
	uuid "github.com/nu7hatch/gouuid"

	"github.com/spf13/cobra"
)

var bulkReport report.Report

var ExportBulkReportCmd = &cobra.Command{
	Use:   "report",
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

		var helper export.ReportReplacePathHelper
		helper.Replacement = "."
		helper.OutputBase = outputDirectory
		export.PathHelper = helper
		result := export.ExportPanelsToLocalFileSystem(outputDirectory, exportConfig, timeRange)
		bulkReport.ReportContent = result
		bulkReport.ReportTemplate = "REPORT"
		bulkReport.TimeRange = timeRange
		report.TemplateSubdir = "markdown"
		report.CreateReport(outputDirectory, &bulkReport)
		log.Debugf("Exported Groups %v", len(result.Results))

	},
}

func init() {
	ExportBulkCmd.AddCommand(ExportBulkReportCmd)
	u, err := uuid.NewV4()
	if err != nil {
		log.Panic(err)
	}
	exportFileName := u.String()

	ExportBulkReportCmd.PersistentFlags().StringVar(&report.TemplateBase, "reporttemplate", "./templates", "The Report GoTemplate Base Directory")
	ExportBulkReportCmd.PersistentFlags().StringVar(&bulkReport.ReportName, "run", exportFileName, "The Testrun name, used for archive the export to the local FileSystem")
	ExportBulkReportCmd.PersistentFlags().StringVar(&bulkReport.Stage, "stage", "", "Executed Stage")
	ExportBulkReportCmd.PersistentFlags().StringVar(&bulkReport.Description, "description", "", "Test Run Description")
	//
	//reportRun.TimeRange = timeRange
}
