FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/wdxtools ./cmd/wdxtools

FROM scratch AS runtime
COPY --from=builder /out/wdxtools /wdxtools
