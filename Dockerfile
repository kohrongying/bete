FROM golang:1.14-buster as build

ADD . /go/src/bete
WORKDIR /go/src/bete/cmd/bete
RUN go build -o /go/bin/bete

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/bete /bete
ENTRYPOINT ["/bete"]