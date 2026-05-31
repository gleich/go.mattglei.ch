# syntax=docker/dockerfile:1
FROM golang:1.26.1 AS build

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/go.mattglei.ch .

FROM alpine:3.20.2

RUN apk update && apk add --no-cache ca-certificates=20260413-r0 tzdata=2026b-r0 git=2.45.4-r0

WORKDIR /src
COPY --from=build /bin/go.mattglei.ch /bin/go.mattglei.ch
COPY --from=build /src/config.toml .
COPY --from=build /src/styles.css .

CMD ["/bin/go.mattglei.ch"]