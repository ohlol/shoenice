package shoenice

import (
	"encoding/json"
	"fmt"
	"github.com/ohlol/graphite-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StatsInstance struct {
	Counters map[string]int `json:"counters"`
	Gauges map[string]int `json:"gauges"`
}

func NewStatsInstance() *StatsInstance {
	return &StatsInstance{
		Counters: make(map[string]int),
		Gauges: make(map[string]int),
	}
}

func (si *StatsInstance) Incr(stat string) {
	si.Counters[stat]++
}

func (si *StatsInstance) IncrN(stat string, val int) {
	si.Counters[stat] += val
}

func (si *StatsInstance) Gauge(stat string, val int) {
	si.Gauges[stat] = val
}

func (si *StatsInstance) fmtGraphite(prefix string) []graphite.Metric {
	var (
		key string
		metrics []graphite.Metric
	)

	now := time.Now().Unix()

	for stat, val := range si.Counters {
		key = strings.Join([]string{prefix, "counters", stat}, ".")
		metrics = append(metrics, graphite.Metric{Name: key, Value: strconv.Itoa(val), Timestamp: now})
	}

	for stat, val := range si.Gauges {
		key = strings.Join([]string{prefix, "gauges", stat}, ".")
		metrics = append(metrics, graphite.Metric{Name: key, Value: strconv.Itoa(val), Timestamp: now})
	}

	return metrics
}

func (si *StatsInstance) RunServer(listenaddr string, prefix string, sendInterval int, graphiteAddr string, graphitePort uint16) {
	statsHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		encoded, _ := json.Marshal(si)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(encoded))
	}

	go func() {
		http.HandleFunc("/stats", statsHandlerFunc)
		http.ListenAndServe(listenaddr, nil)
	}()

	go func() {
		graphite := graphite.Connect(graphite.GraphiteServer{Host: graphiteAddr, Port: graphitePort})
		for {
			graphite.Sendall(si.fmtGraphite(prefix))
			time.Sleep(time.Duration(sendInterval) * time.Second)
		}
	}()
}
