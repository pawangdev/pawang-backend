FROM golang:alpine AS build-image
RUN apk update && apk add build-base git gcc mercurial
ADD . /src
RUN cd /src && go build -o go-pawang

FROM alpine
WORKDIR /app
COPY --from=build-image /src/go-pawang /app/
ENTRYPOINT [ "./go-pawang" ]