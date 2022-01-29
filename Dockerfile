
FROM ghcr.io/lwnmengjing/endpoint-discover:latest

MAINTAINER lwnmengjing <991154416@qq.com>

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]