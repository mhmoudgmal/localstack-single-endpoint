FROM golang:1.13.5-alpine3.10 as build
Maintainer Mahmoud Gamal <mhmoudgmal.89@gmail.com>

RUN apk update && apk add git gcc musl-dev

ENV APP_HOME /build-app

WORKDIR $APP_HOME

COPY go.mod $APP_HOME
COPY go.sum $APP_HOME
RUN go mod download

COPY . $APP_HOME
RUN go test -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o localstack-single-endpoint .

FROM alpine as app
Maintainer Mahmoud Gamal <mhmoudgmal.89@gmail.com>

COPY --from=build /build-app/localstack-single-endpoint /
COPY --from=build /build-app/services.yml /

ENV LOCALSTACK_HOST "localhost"

EXPOSE 9000

CMD ["sh", "-c", "/localstack-single-endpoint -localstackHost ${LOCALSTACK_HOST}"]
