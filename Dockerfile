FROM golang:1.8
RUN go get -u github.com/vimukthi-git/beanstalkg
WORKDIR /src/github.com/vimukthi-git/beanstalkg
COPY . /src/github.com/vimukthi-git/beanstalkg
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /beanstalkg .

FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /beanstalkg /
EXPOSE 11300
CMD ["/beanstalkg"]
