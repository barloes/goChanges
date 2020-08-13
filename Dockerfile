FROM golang:1.14

WORKDIR /app/go-sample-app
RUN go get -u github.com/slotix/pageres-go-wrapper


# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download


COPY . .

# Build the Go app
RUN go build -o ./out/go-sample-app .

# Run the binary program produced by `go install`
CMD ["./out/go-sample-app"]