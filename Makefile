.PHONY: run python-service go-service

python-service:
	@python ./python-service/greeter_server.py

go-service:
	@go run ./go-service/main.go