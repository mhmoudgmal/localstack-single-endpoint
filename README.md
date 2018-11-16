[![Build Status](https://travis-ci.org/mhmoudgmal/localstack-single-endpoint.svg?branch=master)](https://travis-ci.org/mhmoudgmal/localstack-single-endpoint) [![Go Report Card](https://goreportcard.com/badge/github.com/mhmoudgmal/localstack-single-endpoint)](https://goreportcard.com/report/github.com/mhmoudgmal/localstack-single-endpoint) [![Coverage Status](https://coveralls.io/repos/github/mhmoudgmal/localstack-single-endpoint/badge.svg?branch=add-coverall)](https://coveralls.io/github/mhmoudgmal/localstack-single-endpoint?branch=add-coverall)

Localstack single endpoint
---

A simple proxy that tends to understand and analyze aws requests according to [aws docs](https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html)
to detremine which service is being invoked and forward the request to the corresponding Localstack endpoint.
