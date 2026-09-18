FROM golang:1.27.1-alpine3.24 AS build
WORKDIR /tmp/app
COPY . .
RUN apk add --no-cache git make
RUN CGO_ENABLED=0 make build

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /tmp/app/main /app
USER nonroot:nonroot
CMD [ "/app" ]
