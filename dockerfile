FROM golang:1.23.4
WORKDIR /app
COPY . /app
RUN cp .env.example .env
RUN go get .
EXPOSE 8080
CMD ["go","run","."]