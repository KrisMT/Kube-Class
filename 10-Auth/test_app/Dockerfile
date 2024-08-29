FROM golang:alpine AS build

WORKDIR /app

COPY go.mod main.go ./

RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/app

FROM alpine

COPY --from=build /bin/app /bin

EXPOSE 4000

ENTRYPOINT ["/bin/app"]
