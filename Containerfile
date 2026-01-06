FROM alpine AS build

RUN apk add --no-cache git make musl-dev go upx

ENV GOROOT=/usr/lib/go
ENV GOPATH=/go
ENV PATH=/go/bin:$PATH

RUN mkdir build

COPY *.go  go.* build/
RUN cd build && \
    go install && \
    CGO_ENABLED=0 GOOS=linux go build \
      -v  \
      -o glimport \
      -ldflags="-s -w" && \
    upx --ultra-brute -q glimport && upx -t glimport

#-----------------------------------------------------------------------------

FROM scratch

LABEL org.opencontainers.image.authors="Didier FABERT <didier.fabert@gmail.com>"
LABEL eu.tartarefr.glimport.version=1.0.0

COPY --from=build build/glimport /glimport
COPY LICENSE /LICENSE
COPY README.md /README.md

ENTRYPOINT [ "/glimport" ]
