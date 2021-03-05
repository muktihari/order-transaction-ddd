FROM golang:1.13.12-alpine AS builder
WORKDIR /builder

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o app ./main.go
RUN echo "nobody:x:65534:65534:nobody:/:" >> passwd

FROM scratch

COPY --from=builder /builder/passwd /etc/passwd
COPY --from=builder /builder/app ./app

USER nobody

ENTRYPOINT ["./app"]

LABEL maintainer="muktihaz@gmail.com"