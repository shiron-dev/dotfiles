.PHONY: help init plan apply destroy validate fmt lint check clean

init:
	terraform init

plan:
	terraform plan

apply:
	terraform apply

validate:
	terraform validate

fmt:
	terraform fmt -recursive

lint:
	tflint --init && tflint

check: fmt lint validate
