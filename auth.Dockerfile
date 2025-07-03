FROM golang:1.24-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o note-auth cmd/auth/main.go
RUN ls

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/note-auth . 
CMD [ "./note-auth" ]
EXPOSE 9091