package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/olleboll/images/img"
	"github.com/olleboll/images/store"
)

func Run(imageStore store.ImageStore) {

	getImagesMetadata := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var images []img.Image
		images, err := imageStore.GetImagesMetadata()

		if err != nil {
			returnError(w, 500, "Something went wrong")
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
			returnError(w, 400, "Invalid id param")
			return
		}
		image, err = imageStore.GetImageMetadata(imageId)

		if err != nil {
			returnError(w, 404, "Could not find image")
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
			returnError(w, 400, "Invalid id param")
			return
		}

		// Get query params for cropping
		queryValues := r.URL.Query()
		bbox := queryValues.Get("bbox")

		var cutout *img.Cutout

		if len(bbox) > 0 {
			bbox := strings.Split(bbox, ",")
			cutout = &img.Cutout{}
			cutout.X, _ = strconv.Atoi(bbox[0])
			cutout.Y, _ = strconv.Atoi(bbox[1])
			cutout.W, _ = strconv.Atoi(bbox[2])
			cutout.H, _ = strconv.Atoi(bbox[3])
		}

		if err := imageStore.GetImageData(imageId, &imageData, cutout); err != nil {
			returnError(w, 404, "Couold not find image")
			return
		}
		w.Header().Set("Content-Type", "text/base64")
		w.WriteHeader(http.StatusOK)
		w.Write(imageData)
	}

	createImage := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		imageData, err := io.ReadAll(r.Body)

		if err != nil {
			returnError(w, 400, "Invalid body.")
			return
		}

		meta, err := imageStore.CreateImage(&imageData)

		if err != nil {
			returnError(w, 400, "Failed to save to db. Make sure it is a base64 encoded image")
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
			returnError(w, 400, "Invalid body")
			return
		}

		imageId, err := strconv.Atoi(ps.ByName("imageId"))
		if err != nil {
			returnError(w, 400, "Invalid id param")
			return
		}

		meta, err := imageStore.UpdateImage(imageId, &imageData)

		if err != nil {
			returnError(w, 400, "Failed to save to db. Make sure it is a base64 encoded image")
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

func returnError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
