FROM golang
RUN go install github.com/tsenart/vegeta@latest

CMD ['vegeta']