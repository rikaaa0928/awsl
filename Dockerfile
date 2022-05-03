FROM alpine

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY awsl /usr/local/bin/awsl

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["awsl"]
