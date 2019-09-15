# Grafana Bulk Exporter


## Functions

* Export a single Panel as PNG
* Export all Panels from a Dashboard  (**planed**)
* Create a Shareable Report
  * Simple Markdown Format
  * Hugo Structure


## Usage

Export a Grafana Access Token

```bash
export GRAFANA_TOKEN=eyJrIjoiZW1oTTA3eFoxNDRiRnNUTHlIcjJzcGhsOW5QRjJScVEiLCJuIjoiYXBpIiwiaWQiOjF9
```

For a command overview use the ``--help`` parameter like: 

```bash
$ ./grafana-exporter --help

Simple Tool for exporting Informations from the Grafana and create
shareable Reports.

Usage:
  export [flags]
  export [command]

Available Commands:
  bulk        Export a bulk of Panels from given config File.
  help        Help about any command
  panel       Export a single Panel.

Flags:
      --debug                Enable Debug log
  -u, --grafana_url string   Grafana Url (default "http://localhost:3000")
  -h, --help                 help for export
  -t, --token string         Grafana Access Token

Use "export [command] --help" for more information about a command.

```

### Bulk Exports

For Bulk Exports you will need a Configuration File, this config contains the Export Structure and the Exported Panels, used for ``--exportConfig``.

```yaml
aliases:
  - &anchorDashboardInfos
    dashboardId: V3TD6Z5Wk
    orgId: 1
  - &anchorPanelOverviewExportSize
    size:
      width: 300
      height: 160
  - &anchorPanelDashboardVars
    vars:
      someTest: [1, 3]

export:
  overview:
    title: "Overview"
    description: "First test From Yaml"
    elements:
      - title: "Any Custom Panel Title"
        description: "Any Custom Panel Description"
        panel:
          id: 4
          <<: *anchorDashboardInfos
          <<: *anchorPanelDashboardVars
          exportName: firsttest2
        <<: *anchorPanelOverviewExportSize
```

### Folder 


```bash
./grafana-exporter bulk report \
    --run test \
    --export /tmp/grafa_export/simpleReport \
    --exportConfig ./export_description.yml \
    --reporttemplate ./templates \
    --from 1568458539 \
    --end 1568480139 \
    --debug
```

#### Hugo Export
Generate the Report to a Existing Hugo Blog Structure

```bash
./grafana-exporter bulk report hugo \
    --run test \
    --export ./hugo-report-site \
    --exportConfig ./export_description.yml \
    --reporttemplate ./templates \
    --from 1568458539 \
    --end 1568480139 \
    --internal /testreports/ \
    --internalGroup loadtest \
    --debug
```

## Develop

### Build

For Quick building and creating a shareable archive, with templates and a example configuration, use the ``Makefile`` goal ``package``.

```bash
make package
```

you will be find a created archive at ``./bin/grafana-exporter.tar.gz``
