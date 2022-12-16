### Image Api

Copy env example values

`cp .env-example .env`

For local execution the .env-example values will work
Build the docker image

`docker build . -t image-api`

Start database and the api in docker

`docker-compose up`

Use post-man (or whatever) to upload image data using
`POST http://localhost:8000/v1/images`
with a base64 encoded image as the body
