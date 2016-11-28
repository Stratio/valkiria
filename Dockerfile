FROM golang+systemd # TODO: need image with systemd and golang
RUN echo a # TODO: update packages
EXPOSE 8101 8101
WORKDIR /go/src/github.com/stratio/valkiria
COPY . /go/src/github.com/stratio/valkiria
RUN make
CMD valkiria