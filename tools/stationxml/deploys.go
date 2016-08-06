package main

import (
	"path/filepath"
	"sort"

	"github.com/GeoNet/delta/meta"
)

func DataloggerDeploys(install string) (map[string]meta.DeployedDataloggerList, error) {

	dataloggerDeploys := make(map[string]meta.DeployedDataloggerList)

	// where the dataloggers were deployed
	var loggers meta.DeployedDataloggerList
	if err := meta.LoadList(filepath.Join(install, "dataloggers.csv"), &loggers); err != nil {
		return nil, err
	}
	for _, d := range loggers {
		if _, ok := dataloggerDeploys[d.Place]; ok {
			dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
		} else {
			dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
		}
	}

	var recs meta.InstalledRecorderList
	if err := meta.LoadList(filepath.Join(install, "recorders.csv"), &recs); err != nil {
		return nil, err
	}
	for _, r := range recs {
		d := meta.DeployedDatalogger{
			Install: meta.Install{
				Equipment: meta.Equipment{
					Make:   r.Make,
					Model:  r.DataloggerModel,
					Serial: r.Serial,
				},
				Span: meta.Span{
					Start: r.Start,
					End:   r.End,
				},
			},
			Place: r.StationCode,
			Role:  r.LocationCode,
		}
		if _, ok := dataloggerDeploys[d.Place]; ok {
			dataloggerDeploys[d.Place] = append(dataloggerDeploys[d.Place], d)
		} else {
			dataloggerDeploys[d.Place] = meta.DeployedDataloggerList{d}
		}
	}

	// sort each datalogger deployment
	for i, _ := range dataloggerDeploys {
		sort.Sort(dataloggerDeploys[i])
	}

	return dataloggerDeploys, nil
}
