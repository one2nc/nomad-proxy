FROM scratch

COPY nomad-proxy proxy
ENTRYPOINT ["./proxy"]
