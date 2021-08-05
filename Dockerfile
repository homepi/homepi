FROM golang:1.16-alpine AS builder

LABEL maintainer="Alireza Josheghani <josheghani.dev@gmail.com>"

RUN apk add build-base

# Creating work directory
RUN mkdir /code

# Adding project to work directory
ADD . /code

# Choosing work directory
WORKDIR /code

# build project
RUN go build -o homepi .

FROM alpine AS app

COPY --from=builder /code/homepi /usr/bin/homepi

EXPOSE 55283

#ENV HPI_ENVIRONMENT production

ENTRYPOINT ["/usr/bin/homepi"]
CMD ["server", "--host", "0.0.0.0", "--port", "55283"]
