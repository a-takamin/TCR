# TODO: 効率的にする
FROM golang:1.23.0
COPY . .
RUN go build
CMD ./tcr
