package metrics

import "regexp"

func init() {
	Register(collectStatsMetrics)
}

var collectStatsMetrics = map[string]MetricConfig{
	"total_connections_received": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "total_connections_received",
			Help:      "pika serve instance total count of received connections from clients",
			Type:      metricTypeCounter,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "total_connections_received",
		},
	},
	"instantaneous_ops_per_sec": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "instantaneous_ops_per_sec",
			Help:      "pika serve instance prcessed operations in per second",
			Type:      metricTypeGauge,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "instantaneous_ops_per_sec",
		},
	},
	"total_commands_processed": {
		Parser: &normalParser{},
		MetricMeta: &MetaData{
			Name:      "total_commands_processed",
			Help:      "pika serve instance total count of processed commands",
			Type:      metricTypeCounter,
			Labels:    []string{LabelNameAddr, LabelNameAlias},
			ValueName: "total_commands_processed",
		},
	},
	"is_bgsaving": {
		Parser: &regexParser{
			name: "is_bgsaving",
			reg:  regexp.MustCompile(`is_bgsaving:(?P<is_bgsaving>(No|Yes)),?(?P<bgsave_name>[^,\n]*)`),
		},
		MetricMeta: &MetaData{
			Name:      "is_bgsaving",
			Help:      "pika serve instance bg save info",
			Type:      metricTypeGauge,
			Labels:    []string{LabelNameAddr, LabelNameAlias, "bgsave_name"},
			ValueName: "is_bgsaving",
		},
	},
	"is_scaning_keyspace": {
		Parser: &regexParser{
			name: "is_scaning_keyspace",
			reg:  regexp.MustCompile(`is_scaning_keyspace:(?P<is_scaning_keyspace>(No|Yes))[\s\S]*#\s*Time:(?P<keyspace_time>[^\n]*)`),
		},
		MetricMeta: &MetaData{
			Name:      "is_scaning_keyspace",
			Help:      "pika serve instance scan keyspace info",
			Type:      metricTypeGauge,
			Labels:    []string{LabelNameAddr, LabelNameAlias, "keyspace_time"},
			ValueName: "is_scaning_keyspace",
		},
	},
	"is_compact": {
		Parser: &regexParser{
			name: "is_compact",
			reg:  regexp.MustCompile(`is_compact:(?P<is_compact>(No|Yes))[\s\S]*compact_cron:(?P<compact_cron>[^\n]*)[\s\S]*compact_interval:(?P<compact_interval>[^\n]*)`),
		},
		MetricMeta: &MetaData{
			Name:      "is_compact",
			Help:      "pika serve instance compact info",
			Type:      metricTypeGauge,
			Labels:    []string{LabelNameAddr, LabelNameAlias, "compact_cron", "compact_interval"},
			ValueName: "is_compact",
		},
	},
}
