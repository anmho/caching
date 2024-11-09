TERRAFORM_DIR=./terraform
REGION := $(shell terraform -chdir=$(TERRAFORM_DIR) output -raw region)
CLUSTER_NAME := $(shell terraform -chdir=$(TERRAFORM_DIR) output -raw cluster_name)

.PHONY: build
build:
	@go build -o ./bin/api ./cmd/api

.PHONY: run
run: build
	@./bin/api

.PHONY: watch
watch:
	@air --build.cmd "make build" --build.bin "./bin/api"

.PHONY: test
test:
	@go test ./...

.PHONY: cluster
cluster:
	@terraform -chdir=./terraform apply

.PHONY: destroy
destroy:
	@terraform -chdir=./terraform destroy

.PHONY: kubectx
kubectx:
	aws eks --region $(REGION) update-kubeconfig --name $(CLUSTER_NAME)