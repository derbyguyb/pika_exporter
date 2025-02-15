package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/pourer/pika_exporter/discovery"
	"github.com/pourer/pika_exporter/exporter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	hostFile           = flag.String("pika.host-file", getEnv("PIKA_HOST_FILE", ""), "Path to file containing one or more pika nodes, separated by newline. NOTE: mutually exclusive with pika.addr and pika.host-http.")
	hostHttp           = flag.String("pika.host-http", getEnv("PIKA_HOST_HTTP", ""), "Path to http containing one or more pika nodes, {\"instances\": [{\"alias\": \"test1\", \"addr\": \"127.0.0.1:7000\"}]}. NOTE: mutually exclusive with pika.addr and pika.host-file.")
	addr               = flag.String("pika.addr", getEnv("PIKA_ADDR", ""), "Address of one or more pika nodes, separated by comma.")
	password           = flag.String("pika.password", getEnv("PIKA_PASSWORD", ""), "Password for one or more pika nodes, separated by comma.")
	alias              = flag.String("pika.alias", getEnv("PIKA_ALIAS", ""), "Pika instance alias for one or more pika nodes, separated by comma.")
	namespace          = flag.String("namespace", getEnv("PIKA_EXPORTER_NAMESPACE", "pika"), "Namespace for metrics.")
	metricsFile        = flag.String("metrics-file", getEnv("PIKA_EXPORTER_METRICS_FILE", ""), "Metrics definition file.")
	keySpaceStatsClock = flag.Int("keyspace-stats-clock", getEnvInt("PIKA_EXPORTER_KEYSPACE_STATS_CLOCK", -1), "Stats the number of keys at keyspace-stats-clock o'clock every day, in the range [0, 23].If < 0, not open this feature.")
	checkKeyPatterns   = flag.String("check.key-patterns", getEnv("PIKA_EXPORTER_CHECK_KEY_PARTTERNS", ""), "Comma separated list of key-patterns to export value and length/size, searched for with SCAN.")
	checkKeys          = flag.String("check.keys", getEnv("PIKA_EXPORTER_CHECK_KEYS", ""), "Comma separated list of keys to export value and length/size.")
	checkScanCount     = flag.Int("check.scan-count", getEnvInt("PIKA_EXPORTER_CHECK_SCAN_COUNT", 100), "When check keys and executing SCAN command, scan-count assigned to COUNT.")
	listenAddress      = flag.String("web.listen-address", getEnv("PIKA_EXPORTER_WEB_LISTEN_ADDRESS", ":9121"), "Address to listen on for web interface and telemetry.")
	metricPath         = flag.String("web.telemetry-path", getEnv("PIKA_EXPORTER_WEB_TELEMETRY_PATH", "/metrics"), "Path under which to expose metrics.")
	logLevel           = flag.String("log.level", getEnv("PIKA_EXPORTER_LOG_LEVEL", "info"), "Log level, valid options: panic fatal error warn warning info debug.")
	logFormat          = flag.String("log.format", getEnv("PIKA_EXPORTER_LOG_FORMAT", "json"), "Log format, valid options: txt and json.")
	showVersion        = flag.Bool("version", false, "Show version information and exit.")
)

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if envVal, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(envVal); err == nil {
			return v
		}
	}
	return defaultVal
}

func main() {
	flag.Parse()

	log.Println("Pika Metrics Exporter")
	if *showVersion {
		return
	}

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalln("parse log.level failed, err:", err)
	}
	log.SetLevel(level)
	switch *logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}

	var dis discovery.Discovery
	if *hostFile != "" {
		dis, err = discovery.NewFileDiscovery(*hostFile)
	} else if *hostHttp != "" {
		dis, err = discovery.NewHttpDiscovery(*hostHttp)
	} else {
		dis, err = discovery.NewCmdArgsDiscovery(*addr, *password, *alias)
	}
	if err != nil {
		log.Fatalln(" failed. err:", err)
	}

	e, err := exporter.NewPikaExporter(dis, *namespace, *checkKeyPatterns, *checkKeys, *checkScanCount, *keySpaceStatsClock)
	if err != nil {
		log.Fatalln("exporter init failed. err:", err)
	}
	defer e.Close()

	registry := prometheus.NewRegistry()
	registry.MustRegister(e)
	http.Handle(*metricPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
<head><title>Pika Exporter</title></head>
<body>
<h1>Pika Exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>`))
	})

	log.Printf("Providing metrics on %s%s", *listenAddress, *metricPath)
	for _, instance := range dis.GetInstances() {
		log.Println("Connecting to Pika:", instance.Addr, "Alias:", instance.Alias)
	}
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
