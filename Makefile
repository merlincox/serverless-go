build: deps
	make recompile

recompile:
	env GOOS=linux go build -o bin/lambda lambda/main.go

deps:
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-sdk-go/aws
	go get github.com/aws/aws-sdk-go/aws/session
	go get github.com/aws/aws-sdk-go/service/s3

xxsupervise:
	supervisor --no-restart-on exit -e go -i bin --exec make -- recompile

start-local:
	aws-sam-local local start-api

watch:
	make supervise & make start-local

clean:
	rm -rf bin/*
