```
    __  __                     ____  _
   / / / /___  ________  ___  / __ \(_)
  / /_/ / __ \/ __  __ \/ _ \/ /_/ / /
 / __  / /_/ / / / / / /  __/ ____/ /
/_/ /_/\____/_/ /_/ /_/\___/_/   /_/
```
### This project is under construction.

## About HomePi
HomePI (SmartHome) is a web api application that works with raspberry pi and relay boards to control GPIO pins!
This project works with an android app that you can find here https://github.com/HomePi/AndroidKotlin

<a target="_blank" href="https://documenter.getpostman.com/view/471191/SVtbQR4D?version=latest">
  <img src="https://img.shields.io/badge/Postman-api%20documentation-orange?logo=postman&style=for-the-badge" alt="Postman API Documentation">
</a>

## Features
* 🚀 Custom accessories support
* 🎉 Built with Golang

### Environment variables
#### Required environment variables
```env
HPI_SQLITE3_PATH=path/to/homepi.db
HPI_ACCESS_TOKEN_SECRET=random-jwt-secret-access-token
HPI_REFRESH_TOKEN_SECRET=random-jwt-secret-refresh-token
```
#### Optional environment variables
```env
HPI_ACCESS_TOKEN_EXPIRE_TIME=240 # a duration that an access_token could be valid (default "240 minutes")
HPI_REFRESH_TOKEN_EXPIRE_TIME=1440 # a duration that an refresh_token could be valid (default "1440 minutes")
```

## Run the docker image
```bash
$ docker run --restart always --device /dev/ttyAMA0:/dev/ttyAMA0 --device /dev/mem:/dev/mem --volume ./db/data:/code/db/data --privileged -dp 55283:55283 homepi/homepi
```

## docker-compose example
```yaml
version: '3'

services:
  api:
    container_name: homepi_api
    image: homepi/api
    ports:
      - 55283:55283
    environment:
      HPI_SQLITE3_PATH: "/db/data/homepi.db"
      HPI_ACCESS_TOKEN_SECRET: "random-jwt-secret-access-token"
      HPI_REFRESH_TOKEN_SECRET: "random-jwt-secret-refresh-token"
      HPI_ACCESS_TOKEN_EXPIRE_TIME: "240"
      HPI_REFRESH_TOKEN_EXPIRE_TIME: "1440"
    restart: always
    devices:
      - /dev/ttyAMA0:/dev/ttyAMA0
      - /dev/mem:/dev/mem
    privileged: true
    volumes:
      - homepi-db:/db/data/
    networks:
      - homepi

volumes:
  homepi-db:

networks:
  homepi:
    driver: bridge
```

## Contributing
Thank you for considering contributing to the HomePi project!

## License
The HomePI is open-source software licensed under the MIT license.

