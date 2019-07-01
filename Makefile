
GO111MODULE?=on
export GO111MODULE

.PHONY: test
test:
	go test -vet off $(shell go list ./...) -failfast

.PHONY: format
format:
	gofmt -s -l -w $(shell find  . -name '*.go' | grep -vEf .gofmt_exclude)

.PHONY: build
build:
	go generate && GO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./metrics main.go

.PHONY: clean
clean:
	rm -rf ./metrics
	go clean ./...

.PHONY: image
image:
	docker build -t dipakpawar231/stats-collector:0.1 .
	docker tag dipakpawar231/stats-collector:0.1 dipakpawar231/stats-collector:latest


.PHONY: push-image
push-image:
	docker push dipakpawar231/stats-collector:0.1
	docker push dipakpawar231/stats-collector:latest


.PHONY: deploy
deploy:
	kubectl create ns stats || true
	kubectl apply -f k8s/ -n stats

.PHONY: db
db:
	docker run -p 5430:5432 -e POSTGRESQL_ADMIN_PASSWORD=secret -d postgres:11