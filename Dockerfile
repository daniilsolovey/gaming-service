FROM alpine:edge

RUN apk update && apk add --no-cache\
    bash \
    curl \
    tzdata \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY gaming-service /bin/app
COPY config.yaml /etc/gaming-service.yaml

CMD ["/bin/app", "--config=/etc/gaming-service.yaml"]
