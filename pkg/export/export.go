package export

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	gapi "github.com/nolte/go-grafana-api"
	"github.com/nolte/grafana-exporter/pkg/config"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/spf13/viper"
)

var Client *gapi.Client

var PathHelper ReportPathHelper

type ReportPathPlainHelper struct {
}

func (h ReportPathPlainHelper) asReportPath(path string) string {
	return path
}

type ReportReplacePathHelper struct {
	OutputBase  string
	Replacement string
}

func (h ReportReplacePathHelper) asReportPath(path string) string {
	return strings.ReplaceAll(path, h.OutputBase, h.Replacement)
}

type ReportPathHelper interface {
	asReportPath(string) string
}

type BulkExportResult struct {
	Results map[string]BulkExportGroupResult
}

type ExportTimeRange struct {
	From time.Time
	End  time.Time
}

func (r *ExportTimeRange) New(from time.Time, end time.Time) {
	r.From = from
	r.End = end
}

func (r *ExportTimeRange) AsApiObject() gapi.TimeRange {
	apiObject := gapi.TimeRange{}
	apiObject.From = gapi.TimeToGrafanaString(r.From)
	apiObject.To = gapi.TimeToGrafanaString(r.End)
	return apiObject
}

type BulkExportGroupResult struct {
	Title       string
	Elements    []BulkExportElementResult
	Description string
}

func (r *BulkExportGroupResult) New(bulkConfigGroup config.ExportGroup) {
	r.Title = bulkConfigGroup.Title
	r.Description = bulkConfigGroup.Description
}

type BulkExportElementResult struct {
	Title                   string
	Description             string
	PanelID                 int64
	Path                    string
	ReportPath              string
	GrafanaLocatedDashboard *gapi.Dashboard
	DashboardVars           map[string][]string
	ExportSize              gapi.GrafanaPanelExportSize
	TimeRange               ExportTimeRange
}

func (e BulkExportElementResult) Panel() gapi.DashboardPanel {
	panel, _ := e.GrafanaLocatedDashboard.GetPanelFromDashboard(e.PanelID)
	return panel
}
func (e BulkExportElementResult) PanelTitle() string {
	currentPanel := e.Panel()
	return currentPanel.Title
}

func (e BulkExportElementResult) PanelDescription() string {
	currentPanel := e.Panel()
	return currentPanel.Description
}

func (e BulkExportElementResult) PanelLink() string {

	currentPanel := e.Panel()
	exportBase := viper.GetString("grafana.url")
	exportBase += e.GrafanaLocatedDashboard.FrontendURL(e.DashboardVars)
	timeRa, _ := e.TimeRange.AsApiObject().AsPartOfUrl()
	if !strings.HasSuffix(exportBase, "?") {
		exportBase += "&"
	}
	exportBase += timeRa
	exportBase += "&" + currentPanel.AsPartOfUrl()
	log.Debugf("Grafana Link: %s", exportBase)
	return exportBase
}

func ExportPanelsToLocalFileSystem(outputBase string, bulkExportConfig config.BulkExportConfiguration, timeRange ExportTimeRange) BulkExportResult {
	log.Debug("Plain Reporter")
	log.Debugf("New Bulk Exporter with %v Elements", len(bulkExportConfig.Exports))

	var result BulkExportResult

	results := make(map[string]BulkExportGroupResult)
	for groupName, group := range bulkExportConfig.Exports {
		outputBaseGroup := path.Join(outputBase, groupName)
		groupResult := BulkExportGroupResult{}
		//appending meta Data
		groupResult.New(group)

		log.Debugf("Export Group: \"%s\" with %v Panels to Path %s", groupName, len(group.Elements), "groupPath")
		var elements []BulkExportElementResult
		for _, element := range group.Elements {
			os.MkdirAll(outputBaseGroup, os.ModePerm)
			//activity := asExportActivity(element.Panel, outputBaseGroup, timeRange)
			log.Debugf("Export png from %s", element.Panel.DashboardID)

			timeRangeObj := timeRange.AsApiObject()

			outputName := path.Join(outputBaseGroup, getExportFileName(element.Panel))

			panelID := element.Panel.ID
			dashboardID := element.Panel.DashboardID

			exportPanelImageSize := element.ExportSize.AsApiObject()

			dashboardVars := element.Panel.DashboardVars
			log.Debugf("Export the Panel: %v from Dashboard: %v to: %s", dashboardID, panelID, outputName)
			err := Client.ExportPanelAsImage(
				dashboardID,
				element.Panel.OrgID,
				panelID,
				timeRangeObj,
				exportPanelImageSize,
				dashboardVars,
				"Europe/Berlin",
				outputName)

			if err != nil {
				log.Panic("faild to export ", err)
			}
			elementResult := BulkExportElementResult{}

			elementResult.Title = element.Title
			elementResult.Description = element.Description
			elementResult.TimeRange = timeRange
			elementResult.PanelID = panelID
			elementResult.ExportSize = exportPanelImageSize

			elementResult.DashboardVars = dashboardVars

			dashboard, err := Client.GetDashboard(dashboardID)
			if err != nil {
				log.Panic("faild to export ", err)
			}
			elementResult.GrafanaLocatedDashboard = dashboard

			// set the exported object Path informations
			elementResult.Path = outputName
			elementResult.ReportPath = PathHelper.asReportPath(outputName)

			//			elementResult.Title = element.Panel.Title
			//			elementResult.Description = element.Description
			//			elementResult.DashboardID = element.Panel.DashboardID
			//			elementResult.PanelID = element.Panel.ID
			//			elementResult.OrgID = element.Panel.OrgID
			//			elementResult.TimeRange = timeRange
			//			elementResult.ExportSize = activity.ExportSize
			//			elementResult.ExportResult = activity
			//
			elements = append(elements, elementResult)
		}
		groupResult.Elements = elements
		results[groupName] = groupResult
	}
	result.Results = results
	return result
}

//func asExportActivity(p config.ExportPanel, output string, timeRange ExportTimeRange) gapi.GrafanaPanelExport {
//	exportActivity := gapi.GrafanaPanelExport{}
//	exportActivity.Output = path.Join(output, getExportFileName(p))
//	exportActivity.Panel.ID = int(p.ID)
//	exportActivity.ExportRange = timeRange
//	exportActivity.Org.Id = p.OrgID
//	exportActivity.Tz = "Europe%2FBerlin"
//	exportActivity.DashboardVars = p.DashboardVars
//
//	// set dashboard identifier
//	exportActivity.Dashboard.UID = p.DashboardID
//	exportActivity.Dashboard.Title = "notused"
//
//	if p.Height == 0 {
//		exportActivity.ExportSize.Height = 500
//	} else {
//		exportActivity.ExportSize.Height = p.Height
//	}
//
//	if p.Width == 0 {
//		exportActivity.ExportSize.Width = 1000
//	} else {
//		exportActivity.ExportSize.Width = p.Width
//	}
//
//	return exportActivity
//}

func getExportFileName(p config.ExportPanel) string {
	var exportFileName string

	// define the output file name
	if p.ExportName != "" {
		exportFileName = p.ExportName
	} else {
		// generate uuid, used as exported filename
		u, err := uuid.NewV4()
		if err != nil {
			log.Panic(err)
		}
		exportFileName = u.String()
	}
	return fmt.Sprintf("%v.png", exportFileName)
}

func ConvertTimeToDate(dateString string) (time.Time, error) {
	i, err := strconv.ParseInt(dateString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(i, 0)
	return tm, nil
}
