package metrics

import "regexp"

func init() {
	Register(collectBinlogMetrics)
}

var collectBinlogMetrics = map[string]MetricConfig{
	"binlog_log_size": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "log_size",
			Help:      "pika serve instance total binlog size in bytes",
			Type:      metricTypeGauge,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "log_size",
		},
	},

	"binlog_<3.1.0": {
		Parser: &versionMatchParser{
			verC: mustNewVersionConstraint(`<3.1.0`),
			Parser: &regexParser{
				name: "binlog_<3.1.0",
				reg:  regexp.MustCompile(`binlog_offset:(?P<binlog_offset_filenum>[^\s]*)\s*(?P<binlog_offset>[\d]*)`),
			},
		},
		MetricMeta: MetaDatas{
			{
				Name:      "binlog_offset_filenum",
				Help:      "pika serve instance binlog file num",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias},
				ValueName: "binlog_offset_filenum",
			},
			{
				Name:      "binlog_offset",
				Help:      "pika serve instance binlog offset",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias, "safety_purge", "expire_logs_days", "expire_logs_nums"},
				ValueName: "binlog_offset",
			},
		},
	},

	"binlog_>=3.1.0": {
		Parser: &versionMatchParser{
			verC: mustNewVersionConstraint(`>=3.1.0`),
			Parser: &regexParser{
				name: "binlog_>=3.1.0",
				reg:  regexp.MustCompile(`(?P<db>db[\d]+)\s*binlog_offset=(?P<binlog_offset_filenum>[^\s]*)\s*(?P<binlog_offset>[\d]*),*safety_purge=(?P<safety_purge>[^\s\n]*)`),
			},
		},
		MetricMeta: MetaDatas{
			{
				Name:      "binlog_offset_filenum_db",
				Help:      "pika serve instance binlog file num for each db",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias, "db"},
				ValueName: "binlog_offset_filenum",
			},
			{
				Name:      "binlog_offset_db",
				Help:      "pika serve instance binlog offset for each db",
				Type:      metricTypeGauge,
				Labels:    []string{LabelNameAddr, LabelNameAlias, "db", "safety_purge"},
				ValueName: "binlog_offset",
			},
		},
	},
}
