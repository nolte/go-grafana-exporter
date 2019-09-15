package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	gapi "github.com/nolte/go-grafana-api"
	export "github.com/nolte/grafana-exporter/cmd/export"
	report "github.com/nolte/grafana-exporter/cmd/report"
	exports "github.com/nolte/grafana-exporter/pkg/export"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//var cfgFile string

var debugEnabled bool

var rootCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Grafana Panel/s to the Local FileSystem",
	Long: `
Simple Tool for exporting Infrmations from the Grafana and create
shareable Reports.
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Grafana Exporter Version 0.0.1 WiP")
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(report.ExportBulkCmd)
	rootCmd.AddCommand(export.ExportSinglePanelCmd)

	rootCmd.PersistentFlags().BoolVar(&debugEnabled, "debug", false, "Enable Debug log")

	rootCmd.PersistentFlags().StringP("token", "t", "", "Grafana Access Token")
	rootCmd.PersistentFlags().StringP("grafana_url", "u", "http://localhost:3000", "Grafana Url")

	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("grafana.url", rootCmd.PersistentFlags().Lookup("grafana_url"))

}

func initConfig() {

	viper.BindEnv("GRAFANA_TOKEN")

	viper.SetDefault("token", viper.GetString("GRAFANA_TOKEN"))
	viper.SetDefault("grafana.tz", viper.GetString("Europe%2FBerlin"))

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
		log.Debug("Enable Debug Log")
	}
	log.SetOutput(os.Stdout)

	token := viper.GetString("token")
	client, _ := gapi.New(token, viper.GetString("grafana.url"))
	exports.Client = client
}
