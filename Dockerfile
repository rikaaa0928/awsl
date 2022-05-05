FROM golang:1.18.1 as build
COPY . /tmp/work
RUN cd /tmp/work && go env -w CGO_ENABLED="0" && go mod download && GOOS=linux GOARCH=amd64 go build -o awsl

FROM debian

COPY --from=build /tmp/work/awsl /usr/local/bin/awsl

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["awsl"]
