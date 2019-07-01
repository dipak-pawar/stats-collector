package db

import (
	"database/sql"
	"testing"

	"fmt"
	"github.com/dipak-pawar/stats-collector/config"
	"github.com/dipak-pawar/stats-collector/models"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"time"
)

type StoreSuite struct {
	suite.Suite
	store  *dbStore
	db     *sql.DB
	dbName string
}

func (s *StoreSuite) SetupSuite() {
	var err error
	s.dbName = fmt.Sprintf("test_%d", time.Now().Unix())
	connConf := config.Postgres.String()
	DB := Connect(connConf)
	if _, err := DB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, s.dbName)); err != nil {
		s.T().Fatalf("Database creation failed: %s. %#v", s.dbName, err)
	}
	if err = DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	var localConfig config.Configuration
	if err := envconfig.Process("postgresql", &localConfig); err != nil {
		log.Fatal(err.Error())
	}
	localConfig.SetDatabaseName(s.dbName)
	localConnConf := localConfig.String()
	testDB := Connect(localConnConf)
	if err = SchemaMigrate(testDB); err != nil {
		s.T().Fatal("Migration failed.", err.Error())
	}

	s.db = testDB
	s.store = &dbStore{db: testDB}
}

func (s *StoreSuite) SetupTest() {
	_, err := s.db.Query("DELETE FROM node_metrics")
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = s.db.Query("DELETE FROM nodes")

	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {

	//if _, err := s.db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, s.dbName)); err != nil {
	//	s.T().Fatalf("Database drop failed: %s. %#v", s.dbName, err)
	//}
	s.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateNode() {
	s.store.CreateNode("test name")
	s.verifyNode("test name")
}

func (s *StoreSuite) TestGetNode() {
	_, err := s.db.Query(`INSERT INTO nodes (name) VALUES('node1')`)
	if err != nil {
		s.T().Fatal(err)
	}

	s.verifyNode("node1")

	node2, err := s.store.GetNode("node2")
	s.Nil(node2, "node with node2 name should not be found")
	s.Error(err, "not found")
}

func (s *StoreSuite) TestSaveNodeMetrics() {
	//fmt.Println(time.Now().Unix())
	nm := models.NodeMetrics{time.Now().UnixNano(), 10.00, 10.00}
	s.store.SaveNodeMetrics("n1", &nm)
	s.verifyNode("n1")
	nms := make([]models.NodeMetrics, 0)
	nms = append(nms, nm)
	s.verifyNodeMetrics(nms)

	nm1 := models.NodeMetrics{time.Now().UnixNano(), 12.00, 12.00}
	nms = append(nms, nm1)

	s.store.SaveNodeMetrics("n1", &nm1)
	s.verifyNode("n1")
	s.verifyNodeMetrics(nms)
}

func (s *StoreSuite) TestFilterNodeMetrics() {

	nm := models.NodeMetrics{time.Now().UnixNano(), 10.00, 10.00}
	nm1 := models.NodeMetrics{time.Now().Local().Add(time.Second * 10).UnixNano(), 12.00, 12.00}
	nm2 := models.NodeMetrics{time.Now().Local().Add(time.Second * 20).UnixNano(), 14.00, 14.00}
	s.store.SaveNodeMetrics("nf1", &nm)
	s.store.SaveNodeMetrics("nf1", &nm1)
	s.store.SaveNodeMetrics("nf1", &nm2)

	nodeMetrics, err := s.store.FilterNodeMetrics(20)

	assert.NoError(s.T(), err, "failed to filter nodemetrics as per given timeslice")
	assert.Len(s.T(), nodeMetrics, 2)
	assert.Equal(s.T(), nodeMetrics[0], &nm2)
	assert.Equal(s.T(), nodeMetrics[1], &nm1)
}

func (s *StoreSuite) verifyNode(name string) {
	res, err := s.db.Query(`SELECT COUNT(*) FROM nodes WHERE name= $1`, name)
	if err != nil {
		s.T().Fatal(err)
	}

	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}

	node, err := s.store.GetNode(name)
	if err != nil {
		s.T().Fatal(err)
	}

	if node == nil {
		s.T().Errorf("expected node != nil, got %v", *node)
	}

	if node.Name != name {
		s.T().Errorf("incorrect details, expected %v, got %v", name, node.Name)
	}
}

func (s *StoreSuite) verifyNodeMetrics(metrics []models.NodeMetrics) {
	nmResult, err := s.db.Query(`SELECT COUNT(*) FROM node_metrics`)
	if err != nil {
		s.T().Fatal(err)
	}

	var nmCount int
	for nmResult.Next() {
		err := nmResult.Scan(&nmCount)
		if err != nil {
			s.T().Error(err)
		}
	}

	if nmCount != len(metrics) {
		s.T().Errorf("incorrect count, wanted %d, got %d", len(metrics), nmCount)
	}

	// ToDo verify actual content
}
