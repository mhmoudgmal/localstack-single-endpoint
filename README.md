Localstack single endpoint
---

A simple proxy that tends to understand and analyze aws requests according to [aws docs](https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html)
to detremine which service is being invoked and forward the request to the corresponding Localstack endpoint.
