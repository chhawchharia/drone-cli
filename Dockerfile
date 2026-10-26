FROM amd64/alpine:3.20@sha256:0d2a3e86a16d902509fb9dc0aea896f7689506be44b0ba6f12bb72252aaa276f as alpine
RUN apk add -U --no-cache 'ca-certificates>=20261020-r0'

FROM amd64/alpine:3.20@sha256:0d2a3e86a16d902509fb9dc0aea896f7689506be44b0ba6f12bb72252aaa276f
ENV GODEBUG netdns=go

# Add security updates
RUN apk add -U --no-cache 'ca-certificates>=20261020-r0' zlib openssl busybox && \
    # Create a non-root user
    addgroup -S drone && \
    adduser -S -G drone drone

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY release/linux/amd64/drone /bin/

# Change ownership of the binary to the non-root user
RUN chmod +x /bin/drone && \
    chown drone:drone /bin/drone

# Switch to non-root user
USER drone

ENTRYPOINT ["/bin/drone"]
