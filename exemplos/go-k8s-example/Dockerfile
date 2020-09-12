FROM golang

ENV GO111MODULE=on

WORKDIR go-k8s-example

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main main.go

EXPOSE 8080

CMD [ "./main" ]