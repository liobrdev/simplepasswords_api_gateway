FROM golang:1.18-alpine@sha256:c5f96a1e7888a0c634f6c37692317a7d133048a4b9a1f878b84e0ce568831f54

RUN adduser app_user
USER app_user

WORKDIR /api_gateway
COPY . .
RUN go mod download
