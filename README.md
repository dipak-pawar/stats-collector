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

To deploy it, we have make target `deploy`, which creates ns `metrics` and deploy `metrics-collector` and `postgres` svc using required secrets.

`make deploy`

After successful deployment, you can access db service using `HOST=$(minikube ip):31002`.

To connect DB you can run `PGPASSWORD=mysecretpassword psql -h $(minikube ip) -U postgres -d postgres -p 31002`

At the same time, we are exposing metrics-collector service on port `31001`. You can access deployed rest api service using `$(minikube ip):31002` 

# Supported API

* /v1/status
    - METHOD - GET
    - this is used for readyness and liveness probes
