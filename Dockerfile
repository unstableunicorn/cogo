FROM scratch

COPY bin/cogo /

ENTRYPOINT ["/cogo"]
