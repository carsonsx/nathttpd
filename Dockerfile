FROM centurylink/ca-certs

MAINTAINER carsonsx <carsonsx@qq.com>

COPY nathttpd /nathttpd

ENTRYPOINT ["/nathttpd"]
CMD ["--help"]
