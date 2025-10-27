###############################################################################
#                                goframe-build
###############################################################################
FROM --platform=linux/amd64 golang:1.24-alpine AS builder
ENV WORKDIR=/gowlive
WORKDIR $WORKDIR

RUN apk add --no-cache gcc musl-dev

COPY api/ $WORKDIR/api
COPY internal/ $WORKDIR/internal
COPY go.mod $WORKDIR/go.mod
COPY go.sum $WORKDIR/go.sum
COPY main.go $WORKDIR/main.go

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod tidy && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

###############################################################################
#                                vue-build
###############################################################################
FROM node:20-alpine AS vue-builder
WORKDIR /app
COPY web .
RUN rm -rf node_modules dist && \
    npm install && \
    npm run build:prod

###############################################################################
#                                gowlive-build
###############################################################################
FROM --platform=linux/amd64 alpine:3.19

LABEL maintainer="shichen437 <shichen437@126.com>"
LABEL version="0.0.6"
LABEL description="Gowlive"
LABEL license="Mit"
LABEL source="https://github.com/shichen437/Gowlive"

ENV WORKDIR=/gowlive
ENV TZ=Asia/Shanghai
ENV PROJECT_DATA=$WORKDIR/resources/data

USER root
RUN mkdir -p $PROJECT_DATA && \
    apk --no-cache add nginx ffmpeg libc6-compat tzdata coreutils && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

ADD web/nginx.conf /etc/nginx/nginx.conf
COPY web/.env.production /web/.env.production
COPY --from=vue-builder /app/dist /usr/share/nginx/html

COPY manifest/config/config.yaml.example $WORKDIR/config.yaml
COPY manifest/migrate $WORKDIR/manifest/migrate
COPY manifest/docker/entrypoint.sh /entrypoint.sh
COPY --from=builder /gowlive/main $WORKDIR/main

ENV PROJECT_SM4KEY="abcdefghijklmnopqrstuvwxyz123456"
VOLUME $PROJECT_DATA
EXPOSE 12580

RUN chmod +x /entrypoint.sh
WORKDIR $WORKDIR
ENTRYPOINT ["sh"]
CMD ["/entrypoint.sh"]