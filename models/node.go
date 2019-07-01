package models

import (
	"github.com/satori/go.uuid"
)

type Node struct {
	Id   uuid.UUID
	Name string
}

type NodeMetrics struct {
	Timestamp   int64   `jsonapi:"attr,timestamp"`
	CpuUsage    float64 `jsonapi:"attr,cpu_usage"`
	MemoryUsage float64 `jsonapi:"attr,memory_usage"`
}

type AverageNodesMetrics struct {
	Timeslice  float64 `jsonapi:"attr,timeslice"`
	CpuUsed    float64 `jsonapi:"attr,cpu_used"`
	MemoryUsed float64 `jsonapi:"attr,memory_used"`
}
