FROM golang:1.17-alpine AS build
ADD ./ /app
WORKDIR /app
RUN go build -o gac

FROM alpine
WORKDIR /app
COPY --from=build /app/gac .
CMD ["./gac"]