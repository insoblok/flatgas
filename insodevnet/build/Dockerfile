FROM alpine:3.19

LABEL maintainer="Flatgas Core Team"

RUN apk add --no-cache bash curl jq

WORKDIR /app

COPY ./bin/ /usr/local/bin/

ENTRYPOINT ["inso"]
