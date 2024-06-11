FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-litmos"]
COPY baton-litmos /