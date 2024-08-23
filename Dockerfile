# syntax=docker/dockerfile:1

FROM golang:1.22-alpine@sha256:58e52f6ddd39098d23b1a34ec24024acd5fe182245960afe572d83c313aafbed as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simplepasswords_api_gateway .

FROM golang:1.22-alpine@sha256:58e52f6ddd39098d23b1a34ec24024acd5fe182245960afe572d83c313aafbed as runner
WORKDIR /app
RUN addgroup -g 1001 app_user
RUN adduser -S -u 1001 -G app_user app_user
COPY --from=build --chown=app_user:app_user --chmod=500 /app/simplepasswords_api_gateway .
