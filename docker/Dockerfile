FROM centurylink/ca-certs

MAINTAINER carsonsx <carsonsx@qq.com>

#ADD https://raw.githubusercontent.com/carsonsx/nathttpd/master/bin/nathttpd_linux_amd64 /nathttpd
COPY nathttpd_linux_amd64 /nathttpd

ENTRYPOINT ["/nathttpd"]
CMD ["--help"]
