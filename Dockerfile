FROM golang:1.9.2

RUN mkdir -p /src/github.com/aditya87/redis_proxy
ENV GOPATH=/
ENV PATH=$PATH:/bin
ADD . /src/github.com/aditya87/redis_proxy
RUN go get github.com/tools/godep
RUN go get github.com/onsi/ginkgo/ginkgo

WORKDIR /src/github.com/aditya87/redis_proxy

RUN godep restore
RUN GOOS=linux GOARCH=amd64 go build .

FROM redis

ENV REDIS_PORT=7777
ENV REDIS_HOST=localhost
ENV REDIS_PASSWORD=
ENV CACHE_CAPACITY=5
ENV EXPIRATION_TIME=10
ENV PORT=3000

COPY --from=0 /src/github.com/aditya87/redis_proxy/redis_proxy .
CMD redis-server --port ${REDIS_PORT} --daemonize yes && ./redis_proxy
