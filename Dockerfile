FROM centurylink/ca-certs
ADD beanstalkg /
EXPOSE 11300
CMD ["/beanstalkg"]
