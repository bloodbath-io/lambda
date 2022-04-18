## Create function

```bash
aws lambda create-function \
--function-name lock-and-dispatch-event \
--zip-file fileb://main.zip \
--handler main \
--runtime go1.x \
--role "arn:aws:iam::987933226201:role/lambda-basic-execution"
```

## Test invoke

```bash
aws lambda invoke \
--function-name lock-and-dispatch-event \
--invocation-type "RequestResponse" \
response.txt
```

## Update function

```bash
aws lambda update-function-code \
--function-name lock-and-dispatch-event \
--zip-file fileb://main.zip
```

## Build it for AWS Lambda
It needs different build and parameters so it can be run on AWS

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
```

## Zip it

```bash
zip main.zip main
```

## Troubleshooting

### Lambda path error

```elixir
"errorMessage" => "fork/exec /var/task/main: exec format error", "errorType" => "PathError"}
```

This means your build doesn't work with AWS Lambda, please build it with the env written above