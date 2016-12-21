# See https://blog.golang.org/docker for more on deploying Go servers with Docker

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD ./src /go/src

# Build the emporium command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install emporium

# Forward logs to docker log collector
RUN mkdir -p /var/log/emporium \
	&& ln -sf /dev/stdout /var/log/emporium/emporium.log

# Document that the service listens on port 3000
EXPOSE 3000

CMD export WAP_MERCHANT="${WAP_MERCHANT}" && WAP_AUTH="${WAP_AUTH}" && /go/bin/emporium
