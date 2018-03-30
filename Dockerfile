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

WORKDIR /src/github.com/aditya87/redis_proxy/integration
RUN GOOS=linux GOARCH=amd64 go build .

FROM redis

RUN mkdir -p /app
COPY --from=0 /src/github.com/aditya87/redis_proxy/redis_proxy /app/redis_proxy
COPY --from=0 /src/github.com/aditya87/redis_proxy/integration/integration /app/integration
CMD /app/redis_proxy
