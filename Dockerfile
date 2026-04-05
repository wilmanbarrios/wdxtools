FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/numcrn ./cmd/numcrn
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/diffh ./cmd/diffh

FROM scratch AS runtime
COPY --from=builder /out/numcrn /numcrn
COPY --from=builder /out/diffh /diffh
