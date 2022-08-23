FROM gcr.io/distroless/static-debian11:nonroot
COPY ./bin/nebb /
USER nonroot
CMD ["/nebb"]