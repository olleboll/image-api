# Image Api

A simple api for storing and retrieving image data. The images gets stored in a postgres table
with the image data as bytes and some metadata as json. Upload images using base64 encoded string.
When retrieving image data they will be represented as base64 encoded strings aswell.

For simplicity uses Postgres for storage of both meta data and image data.

### Endpoints

​
`GET /v1/images`  
List metadata for stored images.

`GET /v1/images/<id>`  
Get metadata for image with id `<id>`.

`GET /v1/images/<id>/data`  
Get image data for image with id `<id>`. The data will be base64 encoded.  
Optional GET parameter:  
`?bbox=<x>,<y>,<w>,<h>` to get a cutout of the image.

`POST /v1/images`  
Upload new image. Request body should be image data.

`PUT /v1/images/<id>`  
Update image. Request body should be image data as base64 encoded string.

### Run locally

Copy env example values

`cp .env-example .env`

#### Run in docker (recommended)

Build the docker image

`docker build . -t image-api`

Start database and the api in docker

`docker-compose up`

#### Without docker

It should be possible to run a postgres server locally (or wherever) without docker just be sure to
set the .env variables correctly if you do.

Install dependencies

`go mod download`

Run the service

`go run cmd/main.go`

### How to use

The create endpoint expects base64 encoded string as the body. In postman this can be done by simple setting the
`Body` to raw and paste a base64 encoded string (You can use https://www.base64encode.org/ ). Make a post request to

`http://localhost:8000/v1/images`

the response should look something like:

```json
{
  "id": 1,
  "metadata": {
    "size": 944, // (KiB)
    "width": 2489,
    "height": 2489,
    "format": "jpeg",
    "created_at": "2022-12-19 11:30:26.70659 +0100 CET m=+5.618776001"
  }
}
```

You can then use the other endpoints:

GET `http://localhost:8000/v1/images`

```json
[
  {
    "id": 1,
    "metadata": {
      "size": 944,
      "width": 2489,
      "height": 2489,
      "format": "jpeg",
      "created_at": "2022-12-19 11:30:26.70659 +0100 CET m=+5.618776001"
    }
  }
]
```

GET `http://localhost:8000/v1/images/1`

```json
{
  "id": 1,
  "metadata": {
    "size": 944,
    "width": 2489,
    "height": 2489,
    "format": "jpeg",
    "created_at": "2022-12-19 11:30:26.70659 +0100 CET m=+5.618776001"
  }
}
```

GET `http://localhost:8000/v1/images/1/data`

`/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQU [...]`

GET `http://localhost:8000/v1/images/1/data?bbox=10,10,400,400`

`/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAUDBAQEAwUEBAQFBQU [...]`

PUT `http://localhost:8000/v1/images/1`

```json
{
  "id": 1,
  "metadata": {
    "size": 944,
    "width": 2489,
    "height": 2489,
    "format": "jpeg",
    "created_at": "2022-12-19 11:30:26.70659 +0100 CET m=+5.618776001"
  }
}
```
