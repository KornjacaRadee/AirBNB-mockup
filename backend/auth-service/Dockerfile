FROM golang:alpine as build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o server
COPY certificatedb.pem /app

FROM alpine
COPY --from=build_container /app/server /usr/bin
RUN mkdir blacklist
COPY ./blacklist/blacklist.txt ./blacklist/
RUN mkdir logging
COPY ./logging/log.log ./logging/
EXPOSE 8080
ENTRYPOINT ["server"]