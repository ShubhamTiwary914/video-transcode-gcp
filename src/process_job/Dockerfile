
#=== Stage-1: Build ===
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o processjob



#=== Stage-2: Run ===
FROM alpine:3.19

# install bash, ffmpeg, and dependencies
RUN apk add --no-cache bash ffmpeg

ENV LOGS_PATH="/app/logs"
ENV OUT_PATH="/app/logs/stdout"
ENV TMPFS_PATH="/tmp/buffer"
ENV MODE="prod"
ENV FFMPEG_LOG="0"
ENV HLS_BUCKETNAME=""
ENV FILE_ID=""
ENV INPUT_PATH=""


#setup dirs: tmpfs, logs
WORKDIR /app
RUN mkdir -p /tmp/buffer /app/logs/stdout

COPY --from=builder /app/processjob /app/
COPY --from=builder /app/transcoder.sh /app/

ENTRYPOINT [ "/app/processjob" ]
