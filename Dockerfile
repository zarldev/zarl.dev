# Use the official Go image as the base image
FROM golang:1.22
WORKDIR /app
COPY dist/zarldev .
EXPOSE 8080
CMD ["./zarldev"]
