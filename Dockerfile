FROM golang:latest

COPY ./ ./
ENV GOPATH=/
RUN go build cmd/balanceService.go
CMD ["./balanceService"]