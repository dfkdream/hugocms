FROM golang:alpine AS build

RUN apk update && apk upgrade

WORKDIR /hugocms

COPY ./*.go ./

COPY ./go.* ./

RUN go get

RUN go build -ldflags "-s"

FROM node AS webpack

WORKDIR /webpack

COPY ./assets/ui-webpack .

RUN npm install

RUN npx webpack --mode production

FROM alpine

RUN apk update && apk upgrade

RUN apk add hugo

WORKDIR /bin

COPY --from=webpack /webpack/dist ./assets/

COPY ./html ./html

COPY --from=build /hugocms .

ENTRYPOINT ["./hugocms"]
