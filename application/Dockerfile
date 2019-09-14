FROM golang:alpine
RUN mkdir /app 
ADD . /app/
WORKDIR /app
RUN apk add --no-cache git
RUN go get -u github.com/go-sql-driver/mysql
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser
EXPOSE 8080/tcp
CMD ["./main"]
