package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/olleboll/images/img"
	"github.com/olleboll/images/store"
)

func main() {

	// Dirty shortcut to make the database ready for connections

	/*
		make all err checks "one liners"

	*/

	godotenv.Load()

	if os.Getenv("SLOW_START") == "true" {
		time.Sleep(2 * time.Second)
	}

	imageStore, err := store.Connect()

	if err != nil {
		log.Fatal("Could not connect to db")
		return
	}

	getImagesMetadata := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var images []img.Image
		images, err = imageStore.GetImagesMetadata()

		if err != nil {
			returnError(w, 500)
			return
		}
		responseData, _ := json.Marshal(images)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}

	getImageMetadata := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var image img.Image

		imageId, err := strconv.Atoi(ps.ByName("imageId"))
		if err != nil {
			returnError(w, 400)
			return
		}
		image, err = imageStore.GetImageMetadata(imageId)

		if err != nil {
			returnError(w, 404)
			return
		}
		responseData, _ := json.Marshal(image)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}

	getImageData := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var imageData []byte

		imageId, err := strconv.Atoi(ps.ByName("imageId"))
		if err != nil {
			returnError(w, 400)
			return
		}

		// Get query params for cropping
		queryValues := r.URL.Query()
		_bbox := queryValues.Get("bbox")

		var cutout *img.Cutout

		if len(_bbox) > 0 {
			bbox := strings.Split(_bbox, ",")
			cutout.X, _ = strconv.Atoi(bbox[0])
			cutout.Y, _ = strconv.Atoi(bbox[1])
			cutout.W, _ = strconv.Atoi(bbox[2])
			cutout.H, _ = strconv.Atoi(bbox[3])
		}

		if err := imageStore.GetImageData(imageId, &imageData, cutout); err != nil {
			returnError(w, 404)
			return
		}
		w.Header().Set("Content-Type", "text/base64")
		w.WriteHeader(http.StatusOK)
		w.Write(imageData)
	}

	createImage := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		responseData, _ := json.Marshal(meta)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}

	updateImage := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

		responseData, _ := json.Marshal(meta)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	}

	router := httprouter.New()

	router.GET("/v1/images", getImagesMetadata)
	router.POST("/v1/images", createImage)
	router.GET("/v1/images/:imageId", getImageMetadata)
	router.PUT("/v1/images/:imageId", updateImage)
	router.GET("/v1/images/:imageId/data", getImageData)

	log.Println("Listing for requests at http://localhost:8000/v1/images")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func returnError(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte("Error"))
}
