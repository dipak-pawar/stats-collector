package controller

import (
	"database/sql"
	"fmt"
	"github.com/dipak-pawar/stats-collector/db"
	"github.com/dipak-pawar/stats-collector/models"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type NodeMetrics struct {
	ds db.Store
}

func NewNodeMetrics(db db.Store) *NodeMetrics {
	return &NodeMetrics{db}
}

func Register(r *mux.Router, d *sql.DB) {

	nm := NewNodeMetrics(db.NewDBStore(d))

	m := r.PathPrefix("/v1/metrics").Subrouter()
	m.HandleFunc("/node/{nodename}", nm.CreateNodeMetricsHandler).Methods("POST")

	a := r.PathPrefix("/v1/analytics").Subrouter()
	a.HandleFunc("/nodes/average", nm.AverageNodeMetricsHandler).Methods("GET")
	a.HandleFunc("/nodes/average", nm.AverageNodeMetricsHandler).Methods("GET").Queries(
		"timeslice", "{timeslice:[0-9]+}")

}

func (c *NodeMetrics) CreateNodeMetricsHandler(w http.ResponseWriter, r *http.Request) {
	nm := new(models.NodeMetrics)
	if err := jsonapi.UnmarshalPayload(r.Body, nm); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	vars := mux.Vars(r)
	nodename := vars["nodename"]

	c.ds.SaveNodeMetrics(nodename, nm)

	fmt.Printf("========================\n %v", *nm)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(201)
	if err := jsonapi.MarshalPayload(w, nm); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (c *NodeMetrics) AverageNodeMetricsHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	key := params.Get("timeslice")
	if key == "" {
		key = "60"
	}

	timeslice, err := strconv.ParseFloat(key, 9)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	nodeMetrics, err := c.ds.FilterNodeMetrics(int64(timeslice))
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), 500)
		return
	}

	totalCpuUsage := 0.0
	totalMemoryUsage := 0.0
	fmt.Println(len(nodeMetrics), nodeMetrics)
	if len(nodeMetrics) > 1 {
		for _, m := range nodeMetrics {
			totalCpuUsage += m.CpuUsage
			totalMemoryUsage += m.MemoryUsage
		}

		fmt.Println(nodeMetrics[0].Timestamp)
		fmt.Println(nodeMetrics[1].Timestamp)

		timeDiff := float64(nodeMetrics[0].Timestamp - nodeMetrics[len(nodeMetrics)-1].Timestamp)
		if timeDiff < (1000000000 * timeslice) {
			timeslice = timeDiff / 1000000000
		}
	} else {
		if len(nodeMetrics) == 1 {
			totalCpuUsage += nodeMetrics[0].CpuUsage
			totalMemoryUsage += nodeMetrics[0].MemoryUsage
			timeslice = 1.0
		} else {
			timeslice = 0.0
		}
	}

	n := float64(len(nodeMetrics))
	av := new(models.AverageNodesMetrics)
	av.Timeslice = timeslice
	av.CpuUsed = avg(totalCpuUsage, n)
	av.MemoryUsed = avg(totalMemoryUsage, n)

	fmt.Println("----------------------Average------------------------")
	fmt.Println(av)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(200)
	if err := jsonapi.MarshalPayload(w, av); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

//func timeSlice()  {
//	timeDiff := float64(nodeMetrics[0].Timestamp - nodeMetrics[len(nodeMetrics)-1].Timestamp)
//	if timeDiff < (1000000000 * timeslice) {
//		timeslice = timeDiff / 1000000000
//	}
//	n := float64(len(nodeMetrics))
//}

func avg(total float64, len float64) float64 {
	if len > 0.0 {
		return total / len
	}
	return 0.0
}
