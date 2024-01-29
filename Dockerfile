# Use the official Go image as the base image
FROM golang:1.21

RUN mkdir /app
RUN mkdir -p /app/config
RUN mkdir -p /app/data
RUN mkdir -p /app/assets

COPY ./assets /app/assets
COPY dist/zarldotdev /app
EXPOSE 8080

WORKDIR /app
ENTRYPOINT ["./zarldotdev"]
