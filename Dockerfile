# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Rutvik Bhatt <rutvik@syphontechnologies.com>"

# Set the Current Working Directory inside the container
WORKDIR /home/go/summa/summa-auth-api

# Copy go mod, sum files and supervisor conf file
COPY go.mod go.sum supervisord.conf ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/summa-auth-app .


######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /home/go/summa/summa-auth-api

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /home/go/summa/summa-auth-api/build/summa-auth-app .

COPY --from=builder /home/go/summa/summa-auth-api/supervisord.conf /etc/supervisord.conf
COPY --from=builder /home/go/summa/summa-auth-api/docs/swagger.json ./docs
COPY --from=builder /home/go/summa/summa-auth-api/docs/swagger.yaml ./docs
COPY --from=builder /home/go/summa/summa-auth-api/cert.pem .
COPY --from=builder /home/go/summa/summa-auth-api/key.pem .
COPY --from=builder /home/go/summa/summa-auth-api/config.yml .

#Install Supervisor
RUN apk update && apk add --no-cache supervisor

#make directory for summa-auth supervisor logs
RUN mkdir /var/log/summa-auth

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable via supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
#CMD ["./summa-auth-app"]