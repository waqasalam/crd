.PHONY: all build container


build: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o pmd-crd controller/pmd-controller.go

container: build
	docker build -t pmd-crd:latest .
