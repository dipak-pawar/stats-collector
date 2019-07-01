package db

import (
	"database/sql"
	"github.com/dipak-pawar/stats-collector/models"
	"github.com/pkg/errors"
	"log"
)

type Store interface {
	CreateNode(name string) error
	GetNode(name string) (*models.Node, error)
	SaveNodeMetrics(node string, metrics *models.NodeMetrics) error
	FilterNodeMetrics(timeslice int64) ([]*models.NodeMetrics, error)
}

type dbStore struct {
	db *sql.DB
}

func NewDBStore(store *sql.DB) Store {
	return &dbStore{store}
}

func (store *dbStore) CreateNode(name string) error {
	_, err := store.db.Query(`INSERT INTO nodes(name) VALUES ($1)`, name)
	return err
}

func (store *dbStore) GetNode(name string) (*models.Node, error) {
	row := store.db.QueryRow("SELECT id, name from nodes WHERE name = $1", name)
	node := &models.Node{}
	if err := row.Scan(&node.Id, &node.Name); err != nil {
		return nil, err
	}
	return node, nil
}

func (store *dbStore) SaveNodeMetrics(node string, metrics *models.NodeMetrics) error {
	n, err := store.GetNode(node)
	if err == sql.ErrNoRows {

		if err := store.CreateNode(node); err != nil {
			return err
		}
	}
	n, err = store.GetNode(node)
	if err != nil {
		return err
	}

	_, err = store.db.Query(`INSERT INTO node_metrics(cpu_usage, memory_usage, timestamp, node_id) VALUES ($1, $2, $3, $4)`, metrics.CpuUsage, metrics.MemoryUsage, metrics.Timestamp, n.Id)
	return err

}

func (store *dbStore) FilterNodeMetrics(timeslice int64) ([]*models.NodeMetrics, error) {
	row := store.db.QueryRow(`SELECT timestamp FROM node_metrics ORDER BY timestamp DESC LIMIT 1`)
	var epochTime int64
	if err := row.Scan(&epochTime); err != nil {
		return nil, err
	}

	t := epochTime - (timeslice * 1000000000)

	rows, err := store.db.Query(`SELECT cpu_usage, memory_usage, timestamp FROM node_metrics WHERE timestamp >= $1 ORDER BY timestamp DESC`, t)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to get metrics with timestamp greater than given timeslice seconds %d", timeslice)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(errors.Wrap(err, "failed to close database rows"))
		}
	}()

	var nodeMetrics []*models.NodeMetrics
	for rows.Next() {
		nm := models.NodeMetrics{}
		err := rows.Scan(&nm.CpuUsage, &nm.MemoryUsage, &nm.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		nodeMetrics = append(nodeMetrics, &nm)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}
	return nodeMetrics, nil
}
