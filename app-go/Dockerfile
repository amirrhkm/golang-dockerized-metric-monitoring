FROM alpine:latest

LABEL maintainer="amirrfov22@gmail.com"

RUN apk add --no-cache wine

WORKDIR /app

COPY app-go.exe .
COPY entrypoint.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/entrypoint.sh

EXPOSE 8008

ENTRYPOINT ["entrypoint.sh"]
CMD ["wine", "app-go.exe"]
