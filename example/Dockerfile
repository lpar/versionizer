# First stage
# Builds a container with Go in, based on Debian Stretch, 
# and uses it to build the project.

FROM golang:1.13.8-stretch as builder
COPY . /hello-podman
WORKDIR /hello-podman
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o hello-podman

# Second stage, builds a minimal container based on Alpine,
# which is used to run the project.

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /hello-podman .
CMD ["./hello-podman"]

# Begin metadata
LABEL org.opencontainers.image.title="My groovy application"
LABEL org.opencontainers.image.description="This is an example metadata file for an application."
LABEL org.opencontainers.image.authors="user@example.com"
LABEL org.opencontainers.image.vendor="Yoyodyne Propulsion Systems"
LABEL org.opencontainers.image.licenses="GPL-3"
LABEL org.opencontainers.image.version="d5e990d-DEV"
LABEL org.opencontainers.image.revision="d5e990d289c4b546ce9f2856377ce1cd5287075c"
LABEL org.opencontainers.image.created="2020-03-13T16:11:32-05:00"
LABEL org.opencontainers.image.information="https://www.example.com/groovyapp"
LABEL org.opencontainers.image.documentation="https://www.example.com/docs"
LABEL org.opencontainers.image.sourceCode="https://github.com/yoyodyne/groovy"
# End metadata

# I like metadata
