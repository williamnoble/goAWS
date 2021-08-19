```sh
GOOS=linux go build -o build/main main.go
zip -j build/main.zip build.main
aws lambda invoke --function-name hello-lambda --cli-binary-format raw-in-base64-out --payload '{"body": "william"}'  out.json
```

