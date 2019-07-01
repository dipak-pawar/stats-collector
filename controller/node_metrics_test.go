package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"fmt"
	"github.com/dipak-pawar/stats-collector/config"
	"github.com/dipak-pawar/stats-collector/controller"
	"github.com/dipak-pawar/stats-collector/db"
	"github.com/dipak-pawar/stats-collector/models"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"log"
	"time"
)

func TestNodeMetrics(t *testing.T) {
	dbName := fmt.Sprintf("int_test_%d", time.Now().Unix())
	connConf := config.Postgres.String()
	DB := db.Connect(connConf)
	if _, err := DB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)); err != nil {
		t.Fatalf("Database creation failed: %s. %#v", dbName, err)
	}
	if err := DB.Close(); err != nil {
		log.Fatal("Error closing the database connection:", err)
	}

	var localConfig config.Configuration

	if err := envconfig.Process("postgresql", &localConfig); err != nil {
		log.Fatal(err.Error())
	}
	localConfig.SetDatabaseName(dbName)
	localConnConf := localConfig.String()
	testDB := db.Connect(localConnConf)
	if err := db.SchemaMigrate(testDB); err != nil {
		t.Fatal("Migration failed.", err.Error())
	}

	router := mux.NewRouter()
	controller.Register(router, testDB)
	ts := httptest.NewServer(router)
	defer ts.Close()

	// requesting avg without data
	averageMetrics := findAverageMetrics(t, ts.URL)
	assert.Equal(t, averageMetrics, &models.AverageNodesMetrics{0.0, 0.0, 0.0})

	// requesting avg with single record
	timestamp := time.Now()
	timestamp20 := timestamp.Local().Add(time.Second * 20).UnixNano()
	timestamp70 := timestamp.Local().Add(time.Second * 70).UnixNano()

	// insert single record
	requestBytes := []byte(`{"data": {"type": "nodemetrics","attributes": {"timestamp": ` + fmt.Sprintf("%d", timestamp.UnixNano()) + `, "memory_usage": 10.00, "cpu_usage": 10.00}}}`)
	postNodeMetrics(t, ts.URL, requestBytes)

	averageMetrics = findAverageMetrics(t, ts.URL)
	assert.Equal(t, averageMetrics, &models.AverageNodesMetrics{1.0, 10.00, 10.0})

	// insert second record
	requestBytes = []byte(`{"data": {"type": "nodemetrics","attributes": {"timestamp": ` + fmt.Sprintf("%d", timestamp20) + `, "memory_usage": 20.00, "cpu_usage": 20.00}}}`)
	postNodeMetrics(t, ts.URL, requestBytes)

	averageMetrics = findAverageMetrics(t, ts.URL)
	assert.Equal(t, averageMetrics, &models.AverageNodesMetrics{20.0, 15.00, 15.0})

	// insert third record
	requestBytes = []byte(`{"data": {"type": "nodemetrics","attributes": {"timestamp": ` + fmt.Sprintf("%d", timestamp70) + `, "memory_usage": 30.00, "cpu_usage": 30.00}}}`)
	postNodeMetrics(t, ts.URL, requestBytes)

	averageMetrics = findAverageMetrics(t, ts.URL)
	assert.Equal(t, averageMetrics, &models.AverageNodesMetrics{50.00, 25.00, 25.00})

	averageMetrics = findAverageMetrics(t, ts.URL, "100")
	assert.Equal(t, averageMetrics, &models.AverageNodesMetrics{70.00, 20.00, 20.00})

}

func findAverageMetrics(t *testing.T, url string, timeslice ...string) *models.AverageNodesMetrics {
	subURL := url + "/v1/analytics/nodes/average"
	if len(timeslice) > 0 && timeslice[0] != "" {
		subURL += "?timeslice=" + timeslice[0]
	}
	req, _ := http.NewRequest("GET", subURL, nil)
	avClient := http.Client{}
	avResp, err := avClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer avResp.Body.Close()

	avRespPayload := new(models.AverageNodesMetrics)

	if err := jsonapi.UnmarshalPayload(avResp.Body, avRespPayload); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
	}
	return avRespPayload
}

func postNodeMetrics(t *testing.T, url string, requestBytes []byte) {

	reqPayload := new(models.NodeMetrics)
	if err := jsonapi.UnmarshalPayload(bytes.NewReader(requestBytes), reqPayload); err != nil {
		t.Fatal("Unable to unmarshal input:", err)
	}

	req, _ := http.NewRequest("POST", url+"/v1/metrics/node/n1", bytes.NewReader(requestBytes))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	respPayload := new(models.NodeMetrics)
	if err := jsonapi.UnmarshalPayload(resp.Body, respPayload); err != nil {
		t.Fatal("Unable to unmarshal body:", err)
		return
	}

	if !reflect.DeepEqual(reqPayload, respPayload) {
		t.Errorf("Data not matching. \nOriginal: %#v\nNew Data: %#v", reqPayload, respPayload)
	}
}
