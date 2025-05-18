FROM golang:alpine AS base

LABEL org.opencontainers.image.authors="sharovik89@ya.ru"

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    APP_PATH="/home/go/src/github.com/sharovik/toggl-to-jira"

WORKDIR ${APP_PATH}

COPY . .

RUN apk add --no-cache bash && apk add --no-cache make && apk add build-base && apk add --no-cache git

RUN make build

FROM alpine:latest AS run
RUN apk --no-cache add ca-certificates

ENV APP_PATH="/home/go/src/github.com/sharovik/toggl-to-jira"

WORKDIR ${APP_PATH}

COPY --from=base ${APP_PATH}/bin ${APP_PATH}/bin
COPY --from=base ${APP_PATH}/.env ${APP_PATH}/.env
COPY --from=base ${APP_PATH}/ ${APP_PATH}/database.sqlite

# Command to run when starting the container
ENTRYPOINT ["./bin/toggl-to-jira"]
