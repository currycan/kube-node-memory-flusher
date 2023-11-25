GO=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
TAG=1.0.0
BIN=flusher
BIN_PATH=build/$(BIN)
IMAGE=currycan/kube-node-memory-flusher

build: $(wildcard ./cmd/*.go ./core/*.go ./version/*.go ./*.go)
	go build -ldflags "-s -w" -o $(BIN_PATH);

image: build
	docker build -t $(IMAGE):$(TAG) .

push: image
	docker push $(IMAGE):$(TAG)

clean:
	rm -f $(BIN_PATH)