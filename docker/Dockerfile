FROM golang:1.11.2
ENV GO111MODULE=on
RUN chown www-data:www-data ./
CMD go build main.go && ./main