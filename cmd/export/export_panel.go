package export

import (
	log "github.com/sirupsen/logrus"
	gapi "github.com/nolte/go-grafana-api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dashboardID string
var dashboardOrgID int64
var panelID int64
var timeRange gapi.TimeRange
var exportPanelImageSize gapi.GrafanaPanelExportSize
var outputName string
var tz string

var ExportSinglePanelCmd = &cobra.Command{
	Use:   "panel",
	Short: "Export a single Panel.",
	Long: `
Export a Single Panel form a Grafana Dashboard by given Time Range,
as PNG, to the Local Filesystem
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		log.Infof("Start Export Panel %v to %s", panelID, outputName)

		token := viper.GetString("token")
		client, _ := gapi.New(token, viper.GetString("grafana.url"))

		//outputName :=  path.Join(outputName, export.getExportFileName(element.Panel))
		//err := client.ExportPanelAsImage(activity)
		err := client.ExportPanelAsImage(
			dashboardID,
			dashboardOrgID,
			panelID,
			timeRange,
			exportPanelImageSize,
			nil,
			tz,
			outputName)

		if err != nil {
			log.Panic("faild to export ", err)
		}

		log.Infof("Panel %v Successfull exported to %s", panelID, outputName)
	},
}

func init() {
	ExportSinglePanelCmd.PersistentFlags().StringVar(&dashboardID, "dashboardID", "V3TD6Z5Wk", "The Dashboard UID")
	ExportSinglePanelCmd.PersistentFlags().Int64Var(&dashboardOrgID, "dashboardOrgID", 1, "The Organisation ID where the Dashboard ar Placed")
	ExportSinglePanelCmd.PersistentFlags().Int64Var(&panelID, "dashboardPanelID", 2, "The ID From the Panel for the Export")

	ExportSinglePanelCmd.PersistentFlags().StringVar(&timeRange.From, "from", "1567685127437", "Graph Start Time")
	ExportSinglePanelCmd.PersistentFlags().StringVar(&timeRange.To, "end", "1567706727437", "Graph End Time")

	ExportSinglePanelCmd.PersistentFlags().IntVar(&exportPanelImageSize.Width, "width", 1000, "Exported PNG Width")
	ExportSinglePanelCmd.PersistentFlags().IntVar(&exportPanelImageSize.Height, "height", 500, "Exported PNG Height")
	ExportSinglePanelCmd.PersistentFlags().StringVar(&tz, "tz", "Europe/Berlin", "Grafana Export TimeZone")

	ExportSinglePanelCmd.PersistentFlags().StringVar(&outputName, "export", "output.png", "The Target Output Full file path, with *.png as file ending")
}
