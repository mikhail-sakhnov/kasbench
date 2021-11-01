FROM golang
RUN go install github.com/tsenart/vegeta@latest
COPY main.go /code/main.go
COPY go.mod /code/go.mod
COPY go.sum /code/go.sum
COPY docker_entrypoint.sh /code/docker_entrypoint.sh
ENV DURATION 60s
WORKDIR /code

RUN go mod tidy
RUN go build -o /bin/create_targets main.go

CMD [ "/code/docker_entrypoint.sh", "/kubeconfig", "/report"]