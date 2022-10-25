FROM golang:1.16-alpine
COPY . /app
WORKDIR /app
RUN go mod download
EXPOSE 9090
CMD [ "go","run","student_app.go" ]