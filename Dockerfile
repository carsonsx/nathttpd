FROM golang:1.7.1

MAINTAINER carsonsx <carsonsx@qq.com>

COPY . /go/src/github.com/carsonsx/nathttpd
WORKDIR /go/src/github.com/carsonsx/nathttpd

RUN CGO_ENABLED=0 go install -v -a -installsuffix nathttpd -ldflags "-s -w -X github.com/carsonsx/nathttpd/version.GITCOMMIT=`git rev-parse --short HEAD`"

ENTRYPOINT ["nathttpd"]
CMD ["--help"]
