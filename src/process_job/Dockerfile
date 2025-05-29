FROM golang:1.24-alpine

WORKDIR /app


ENV LOGS_PATH="/app/logs"
ENV OUT_PATH="/app/logs/stdout"
ENV TMPFS_PATH="/tmp/buffer"
ENV MODE="prod"
ENV HLS_BUCKETNAME=""
ENV FILE_ID=""
ENV INPUT_PATH=""


#setup dirs: tmpfs, logs
RUN mkdir /tmp/buffer
RUN mkdir /app/logs
RUN mkdir /app/logs/stdout


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o processjob

ENTRYPOINT [ "/app/processjob" ]
