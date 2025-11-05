FROM golang:1.25.4-alpine3.22 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code into the container and build the app
COPY . .
RUN go build -v -o /dist/geoint_be

# Deployment stage
FROM alpine:3.22
WORKDIR /usr/src/app
COPY --from=build /usr/src/app ./
COPY --from=build /dist/geoint_be /usr/local/bin/geoint_be

# Start the application
CMD geoint_be