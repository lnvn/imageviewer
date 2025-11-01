FROM golang:1.22 as build

WORKDIR /go/src/app
COPY . .

# RUN go mod download
# RUN go vet -v
# RUN go test -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /go/bin/app main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]