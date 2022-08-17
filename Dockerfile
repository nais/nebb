FROM golang:1.19 as build

WORKDIR /go/src/nebb
COPY . .

RUN go mod download
RUN make nebb

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build /go/src/nebb/bin/nebb /
USER nonroot
CMD ["/nebb"]