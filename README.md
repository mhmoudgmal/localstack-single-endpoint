[![Build Status](https://travis-ci.org/mhmoudgmal/localstack-single-endpoint.svg?branch=master)](https://travis-ci.org/mhmoudgmal/localstack-single-endpoint) [![Go Report Card](https://goreportcard.com/badge/github.com/mhmoudgmal/localstack-single-endpoint)](https://goreportcard.com/report/github.com/mhmoudgmal/localstack-single-endpoint) [![Coverage Status](https://coveralls.io/repos/github/mhmoudgmal/localstack-single-endpoint/badge.svg?branch=add-coverall)](https://coveralls.io/github/mhmoudgmal/localstack-single-endpoint?branch=add-coverall)

Localstack single endpoint
---

A simple proxy that tends to understand and analyze aws requests according to [aws docs](https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html)
to detremine which service is being invoked and forward the request to the corresponding Localstack endpoint.

How to use?
---
#### Docker

- Run the docker container `docker run -p 9000:9000 mhmoudgmal/localstack-single-endpoint`
- Then you can access localstack using a single endpoint port (9000).

For example invoking a lambda will be the same as for listing the tables in a dynamodb:

- `aws lambda invoke --endpoint-url http://localhost:9000 --function function-name --payload '{}'`
- `aws dynamodb list-tables --endpoint-url http://localhost:9000`
