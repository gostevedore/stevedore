# ARG image_from
# ARG golang_version

# FROM docker-hub.stevedore.test:5000/library/golang:${golang_version}-alpine AS build

# WORKDIR /usr/src/app
# COPY . ./
# RUN go build -o /app main.go

# FROM ${image_from}
# COPY --from=build /app /app

# CMD ["/app"]

ARG image_from

FROM ${image_from}

CMD ["echo", "Hey there! Welcome to app1"]
