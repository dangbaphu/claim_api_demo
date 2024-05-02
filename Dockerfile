FROM golang:1.22

WORKDIR /app

COPY . .

RUN GO111MODULE=on


RUN GOPROXY=https://goproxy.cn go install -mod=mod github.com/githubnemo/CompileDaemon


ENTRYPOINT CompileDaemon --build="go build -o /main ./cmd/main" --command="/main"