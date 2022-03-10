# This Dockerfile is based on the simple example given in the README
# at https://hub.docker.com/_/golang
FROM golang:1.17

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/ ./...

ENV TZ="Europe/London"
ENV BROKER_USERNAME="((username))"
ENV BROKER_PASSWORD="((password))"
ENV BROKER_API_URL="https://customer.elephantsql.com/api"
ENV BROKER_API_KEY="((api_key))"
ENV BROKER_API_DEFAULT_REGION="google-compute-engine::europe-west2"

CMD ["elephantsql-broker"]
