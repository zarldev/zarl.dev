# Use the official Go image as the base image
FROM golang:1.22

RUN mkdir /app
RUN mkdir -p /app/config
RUN mkdir -p /app/data
RUN mkdir -p /app/assets

COPY ./assets/css/app.css /app/assets/css/app.css
COPY ./assets/favicon /app/assets/favicon
COPY dist/zarldotdev /app
EXPOSE 8080

WORKDIR /app
ENTRYPOINT ["./zarldotdev"]
