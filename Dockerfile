FROM golang:1.19-alpine AS build
ADD ./ /app
WORKDIR /app
RUN apk add git
RUN go build -o gac

FROM alpine
ARG USER=nonroot
RUN adduser -D $USER
USER $USER
WORKDIR /app
COPY --from=build --chown=$USER:$USER /app/gac .
CMD ["./gac"]
