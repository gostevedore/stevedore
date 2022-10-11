ARG image_from

FROM golang:alpine AS build

WORKDIR /usr/src/app
COPY . ./
RUN go build -o /app main.go

FROM ${image_from}
COPY --from=build /app /app

CMD ["/app"]