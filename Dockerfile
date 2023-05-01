## ===> Build
FROM golang:1.19.3-alpine3.16 AS build

WORKDIR /

# Copy the dependencies list
COPY ./src/go.mod ./
COPY ./src/go.sum ./
# Install dependencies
RUN go mod download
#Copy the backend code
COPY ./src/*.go ./
# Build the backend
RUN CGO_ENABLED=0 go build -o ./bot-exe


## ===> Deploy
FROM gcr.io/distroless/base-debian10

#create a directory for the app
WORKDIR /

# Copy the backend binary
COPY --from=build /bot-exe ./

# Execute the backend
CMD [ "./bot-exe" ]