FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/numcrn ./cmd/numcrn

FROM builder AS tester
RUN go test ./... -v

FROM builder AS bencher
RUN go test ./... -bench=. -benchmem

FROM scratch AS runtime
COPY --from=builder /out/numcrn /numcrn
ENTRYPOINT ["/numcrn"]
