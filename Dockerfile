FROM scratch
ADD pkg/linux_amd64/ebscli /bin/ebscli
ENTRYPOINT ["ebscli"]
CMD ["-h"]
