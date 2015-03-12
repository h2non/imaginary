# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Server port to listen
ENV PORT 80

# Copy the local package files to the container's workspace.
#ADD . /go/src/github.com/h2non/imgine

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/h2non/imgine

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/imgine

# Expose the server TCP port
EXPOSE 80
