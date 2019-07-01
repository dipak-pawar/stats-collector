# Stats Collector

# Prerequisites
- go (1.12.6)
- go-bindata
- minikube (> v0.34.1)

# Building
We are using go modules. In order to build this project, make sure that you have exported `GO111MODULE`. You can export it using `export GO111MODULE=on`.

Install `go-bindata` using `go get -u github.com/go-bindata/go-bindata/...`.

In order to simplify building a project, we have make target by name `build`.

`make build`

# Running Tests
Make sure that you have postgres database running which is exposed on 5430 on localhost. However we have make target by name `db`, which is starting postgres as required for tests

To run tests we have make target by name `test`.

### TL;DR:
```bash
make db
make test
```

# Deploying to K8s
In order to deploy your application in k8s environment, you need to have k8s cluster up and running upfront. I have used `minikube` to deploy this application using maninfests defined in `k8s/` directory.

You can start minikube by installing and runnning it using `minikube start`.

To deploy it, we have make target `deploy`, which creates ns `stats` and deploy `stats-collector` and `postgres` svc using required secrets.

`make deploy`

After successful deployment, you can access db service using `HOST=$(minikube ip):31002`.

To connect DB you can run `PGPASSWORD=mysecretpassword psql -h $(minikube ip) -U postgres -d postgres -p 31002`

At the same time, we are exposing stats-collector service on port `31001`. You can access deployed rest api service using `$(minikube ip):31002`

# Supported API

* /v1/status
    - METHOD - GET
    - this is used for readyness and liveness probes
* /v1/metrics/node/{nodename}
    - METHOD - POST
    - Payload - `{"data": {"type": "nodemetrics","attributes": {"timestamp": 1561878024549969344, "memory_usage": 10.0, "cpu_usage": 10.0}}}`
    - timestamp here is unix epoch in nanoseconds, you can get it using `date +%s%N`
* /v1/analytics/nodes/average?timeslice=30
    - METHOD - GET
    - this return the payload of the form `{"data":{"type":"","attributes":{"cpu_used":15,"memory_used":15,"timeslice":40}}}`

# Payload Schema
- API - /v1/metrics/node/{nodename}
    - Method - POST
    - Request Payload
        ```bash
        {"data": {"type": "","attributes": {"timestamp": int64, "memory_usage": float64, "cpu_usage": float64}}}
        e.g.
	   {"data": {"type": "","attributes": {"timestamp": 1561878024549969344, "memory_usage": 10.0, "cpu_usage": 10.0}}}
        ```
    - timestamp here is unix epoch in nanoseconds, you can get it using `date +%s%N`
    - Response Payload
        ```bash
        {"data":{"type":"","attributes":{"cpu_usage":float64,"memory_usage":float64,"timestamp":int64}}}
        e.g.
	    {"data":{"type":"","attributes":{"cpu_usage":10,"memory_usage":10,"timestamp":1561878024549969408}}}
        ```
- API - /v1/analytics/nodes/average?timeslice=30
    - Method - GET
    - Response Payload
        ```bash
        {"data":{"type":"","attributes":{"cpu_used": float64,"memory_used": float64,"timeslice": float64}}}
        e.g.
	   {"data":{"type":"","attributes":{"cpu_used":20.0,"memory_used":20.0,"timeslice":70.0}}}
        ```

# Using APIs using curl
- `export HOST=$(minikube ip):31001`
- `curl -X POST $HOST/v1/metrics/node/h1 -H Content-type:application/json -H Accept:application/json -d '{"data": {"type": "","attributes": {"timestamp": 1561878024549969344, "memory_usage": 10.0, "cpu_usage": 10.0}}}'`
- `curl "$HOST/v1/analytics/nodes/average"`

# Next Steps
    - Implement Process APIs
    - Improve Tests
    - Code Refactoring
