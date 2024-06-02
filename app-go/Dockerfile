FROM alpine:latest

LABEL maintainer="amirrfov22@gmail.com"

RUN apk add --no-cache wine

WORKDIR /app

COPY app-go.exe .

EXPOSE 8008

CMD ["wine", "app-go.exe"]
