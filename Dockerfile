FROM golang:onbuild

RUN go get gopkg.in/yaml.v2 &\
    go get gotest.tools/assert

RUN go test -v

RUN go build -o localstack-single-endpoint

EXPOSE 9000

CMD ["./localstack-single-endpoint"]
