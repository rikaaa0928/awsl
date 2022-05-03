FROM alpine

COPY ./awsl /usr/local/bin/awsl

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["awsl"]
