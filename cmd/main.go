package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/olleboll/images/img"
	"github.com/olleboll/images/store"
)

func main() {
	godotenv.Load()
	/*

		// Define the endpoints
		​
		`GET /v1/images`
		​
		List metadata for stored images.
		​
		`GET /v1/images/<id>`
		​
		Get metadata for image with id `<id>`.
		​
		`GET /v1/images/<id>/data`
		​
		Get image data for image with id `<id>`.
		​
		Optional GET parameter: `?bbox=<x>,<y>,<w>,<h>` to get a cutout of the image.
		​
		`POST /v1/images`
		​
		Upload new image. Request body should be image data.
		​
		`PUT /v1/images/<id>`
		​
		Update image. Request body should be image data.
	*/

	imageStore, err := store.Connect()

	if err != nil {
		log.Fatal("Could not connect to db")
		return
	}

	getImages := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var responseData []byte
		var images []img.Image
		images, err = imageStore.GetImages()

		if err != nil {
			returnError(w, 500)
			return
		}
		responseData, _ = json.Marshal(images)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	}

	getImage := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var responseData []byte
		var image img.Image

		imageId, err := strconv.Atoi(ps.ByName("imageId"))
		if err != nil {
			returnError(w, 400)
			return
		}
		image, err = imageStore.GetImage(imageId)

		if err != nil {
			returnError(w, 404)
			return
		}
		responseData, _ = json.Marshal(image)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	}

	getImageData := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var imageData []byte

		// Get query params for cropping
		imageId, err := strconv.Atoi(ps.ByName("imageId"))
		if err != nil {
			returnError(w, 400)
			return
		}

		// Get the image data
		// bbox=<x>,<y>,<w>,<h>

		queryValues := r.URL.Query()
		_bbox := queryValues.Get("bbox")
		if len(_bbox) > 0 {
			bbox := strings.Split(_bbox, ",")
			cutout := img.Cutout{}
			cutout.X, _ = strconv.Atoi(bbox[0])
			cutout.Y, _ = strconv.Atoi(bbox[1])
			cutout.W, _ = strconv.Atoi(bbox[2])
			cutout.H, _ = strconv.Atoi(bbox[3])
			err = imageStore.GetImageData(imageId, &imageData, &cutout)
		} else {
			err = imageStore.GetImageData(imageId, &imageData, nil)
		}

		if err != nil {
			returnError(w, 500)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/base64")
		w.Write(imageData)
	}

	createImage := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var responseData []byte
		imageData, err := io.ReadAll(r.Body)

		if err != nil {
			returnError(w, 400)
			return
		}

		meta, err := imageStore.CreateImage(&imageData)

		if err != nil {
			returnError(w, 400)
			return
		}
		responseData, _ = json.Marshal(meta)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	}

	updateImage := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var responseData []byte
		imageData, err := io.ReadAll(r.Body)

		if err != nil {
			returnError(w, 400)
			return
		}

		imageId, err := strconv.Atoi(ps.ByName("imageId"))
		if err != nil {
			returnError(w, 400)
			return
		}

		meta, _err := imageStore.UpdateImage(imageId, &imageData)

		if _err != nil {
			// Failed to save to db
			returnError(w, 500)
			return
		}

		responseData, _ = json.Marshal(meta)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	}

	router := httprouter.New()

	router.GET("/v1/images", getImages)
	router.POST("/v1/images", createImage)
	router.GET("/v1/images/:imageId", getImage)
	router.PUT("/v1/images/:imageId", updateImage)
	router.GET("/v1/images/:imageId/data", getImageData)

	log.Println("Listing for requests at http://localhost:8000/v1/images")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func returnError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte("Error"))
}
