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

# Run the emporium command by default when the container starts.
CMD /go/bin/emporium

# Document that the service listens on port 3000.
EXPOSE 3000
