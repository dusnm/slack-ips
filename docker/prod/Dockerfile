FROM golang:1.26-trixie AS build

WORKDIR /src

COPY . .

ENV CGO_ENABLED=0
ENV GOPROXY=direct

RUN go mod download && \
    go mod verify && \
    go build \
    -ldflags='-s -w -extldflags "-static"' \
    -o /bin/slack-ips ./main.go

FROM debian:trixie AS runtime

COPY --from=build /bin/slack-ips /bin/slack-ips

RUN apt-get update && \
    apt-get install -y ca-certificates

CMD ["/bin/slack-ips"]