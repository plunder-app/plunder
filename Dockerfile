# syntax=docker/dockerfile:experimental

# Build BOOTy as an init
FROM golang:1.14-alpine as dev
RUN apk add --no-cache git ca-certificates make
COPY . /go/src/github.com/plunder-app/plunder
WORKDIR /go/src/github.com/plunder-app/plunder
ENV GO111MODULE=on
RUN --mount=type=cache,sharing=locked,id=gomod,target=/go/pkg/mod/cache \
    --mount=type=cache,sharing=locked,id=goroot,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux make build

FROM scratch
COPY --from=dev /go/src/github.com/plunder-app/plunder/plunder /
ENTRYPOINT ["/plunder"]