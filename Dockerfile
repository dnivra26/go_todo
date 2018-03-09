FROM golang:1.8
ADD . /app/ 
WORKDIR /app 
RUN go build -o main . 
CMD ["/app/main"]