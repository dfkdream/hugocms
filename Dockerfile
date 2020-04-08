FROM golang:alpine AS build

RUN apk update && apk upgrade

WORKDIR /hugocms

COPY . .

RUN find . -type f ! -name '*.go' -and ! -name 'go.*' -delete

RUN go get

RUN go build -ldflags "-s"

FROM node AS webpack

WORKDIR /webpack

COPY ./assets/ui-webpack/package.json .

COPY ./assets/ui-webpack/package-lock.json .

RUN npm install

COPY ./assets/ui-webpack .

RUN npx webpack --mode production

FROM alpine

RUN apk update && apk upgrade

RUN apk add hugo

WORKDIR /hugocms

COPY --from=webpack /webpack/dist ./assets/

COPY ./html ./html

COPY --from=build /hugocms/hugocms .

ENTRYPOINT ["./hugocms"]
