package shoenice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatsInstance struct {
	Counters map[string]int `json:"counters"`
	Gauges map[string]interface{} `json:"gauges"`
}

func NewStatsInstance() *StatsInstance {
	return &StatsInstance{
		Counters: make(map[string]int),
		Gauges: make(map[string]interface{}),
	}
}

func (si *StatsInstance) Incr(stat string, val int) {
	si.Counters[stat] += val
}

func (si *StatsInstance) Gauge(stat string, val interface{}) {
	si.Gauges[stat] = val
}

func (si *StatsInstance) RunServer(listenaddr string) {
	statsHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		encoded, _ := json.Marshal(si)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(encoded))
	}

	go func() {
		http.HandleFunc("/stats", statsHandlerFunc)
		http.ListenAndServe(listenaddr, nil)
	}()
}
