FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
# COPY *.go ./
COPY . .
RUN go build -o /app/build/realworld-go-nolambda .

EXPOSE 8080

CMD [ "/app/build/realworld-go-nolambda" ]