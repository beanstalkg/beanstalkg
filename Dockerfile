FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD beanstalkg /
EXPOSE 11300
CMD ["/beanstalkg"]

# Note: MUST compile beanstalkg as a static binary prior to building
# this container.

# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o beanstalkg .
