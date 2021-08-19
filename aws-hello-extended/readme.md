aws lambda invoke --function-name lambda-hello-extended --cli-binary-format raw-in-base64-out --payload '{"first_name":"William", "last_na
me":"Noble", "Age": 34}' out.json

	  https://github.com/sebastianlacuesta/localstack-sample-go-lambda