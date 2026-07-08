#Build
FROM golang:1.24 as build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o hackathon-app

#Application
FROM gcr.io/distroless/static-debian11

WORKDIR /app
COPY --from=build /app/hackathon-app .

EXPOSE 8000

CMD ["/app/hackathon-app"]