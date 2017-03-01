FROM alpine:latest
COPY bin/ebscli /bin/ebscli
RUN chmod +x /bin/ebscli
ENTRYPOINT ["ebscli"]
CMD ["-h"]
