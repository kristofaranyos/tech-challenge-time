FROM golang:alpine AS build
WORKDIR /src/
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /bin/go_service ./cmd/main.go

FROM scratch
COPY --from=build /bin/go_service /bin/go_service
COPY --from=build /src/app.properties /
EXPOSE 8080
ENTRYPOINT ["/bin/go_service"]