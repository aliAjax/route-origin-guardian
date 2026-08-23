FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/route-origin-guardian ./cmd/guardian
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/route-origin-guardian /route-origin-guardian
EXPOSE 8080 1790 11019
ENTRYPOINT ["/route-origin-guardian"]
