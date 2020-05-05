FROM golang:alpine AS build

RUN apk update && apk upgrade

WORKDIR /hugocms

COPY ./admin ./admin

COPY ./adminapi ./adminapi

COPY ./article ./article

COPY ./config ./config

COPY ./hugo ./hugo

COPY ./internal ./internal

COPY ./plugin ./plugin

COPY ./pluginapi ./pluginapi

COPY ./protowrapper ./protowrapper

COPY ./session ./session

COPY ./signin ./signin

COPY ./user ./user

COPY ./*.go ./

COPY ./go.* ./

RUN go get

RUN go build -ldflags "-s"

FROM node AS webpack

WORKDIR /webpack

COPY ./assets/ui-webpack/package.json .

COPY ./assets/ui-webpack/package-lock.json .

RUN npm install

COPY ./assets/ui-webpack .

ARG mode

RUN npx webpack --mode $mode

FROM alpine

RUN apk update && apk upgrade

RUN apk add hugo

WORKDIR /hugocms

COPY --from=webpack /webpack/dist ./assets/

COPY ./html ./html

COPY --from=build /hugocms/hugocms .

ENTRYPOINT ["./hugocms"]
