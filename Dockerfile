FROM golang:1.17-bullseye

RUN mkdir -p /opt/wager_app
COPY . /opt/wager_app/

WORKDIR /opt/wager_app
RUN echo ">> download libraries"
RUN go mod download
RUN go mod tidy
RUN echo ">> run unit test"
RUN go test -v ./...
RUN echo ">> create app binary"
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api_server ./

ENTRYPOINT [ "/opt/wager_app/api_server" ]