package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/olleboll/images/store"
	"github.com/olleboll/images/img"
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

	images := func(w http.ResponseWriter, req *http.Request) {

		// split the path
		path := strings.Split(req.URL.Path, "/")[3:]
		var imageId int
		if (len(path) > 0) {
			imageId, err = strconv.Atoi(path[0])
			if (err != nil) {
				// id params missing or not a valid id
			}
		}

		var responseData []byte

		switch {
			case len(path) == 0 && req.Method == "GET":

				var images []img.Image
				images, err = imageStore.GetImages()
	
				if err != nil {
					// Return some error
					fmt.Println("ERROR")
				}
				responseData, _ = json.Marshal(images)
				w.Header().Set("Content-Type", "application/json")
				
			case len(path) == 0 && req.Method == "POST":

				fmt.Println("POSTING")
				imageData, err := io.ReadAll(req.Body)
				meta, err := imageStore.CreateImage(imageData)
	
				if err != nil {
	
				}
				responseData, _ = json.Marshal(meta)
				w.Header().Set("Content-Type", "application/json")

			case len(path) == 1 && req.Method == "GET":
				var image img.Image
				image, err = imageStore.GetImage(imageId)
	
				if err != nil {
					// Return some error
					fmt.Println("ERROR")
				}
				responseData, _ = json.Marshal(image)
				w.Header().Set("Content-Type", "application/json")
				

			case len(path) == 1 && req.Method == "PUT":

				fmt.Println("POSTING")
				imageData, err := io.ReadAll(req.Body)

				if err != nil {
					// invalid body error
				}

				meta, err := imageStore.UpdateImage(imageId, imageData)

				if err != nil {
					// Failed to save to db
				}

				responseData, _ = json.Marshal(meta)
				w.Header().Set("Content-Type", "application/json")
				

			case len(path) == 2 && req.Method == "GET" && path[1] == "data":
				// Get query params for cropping

				var imageData []byte
				err = imageStore.GetImageData(imageId, &imageData)
	
				if err != nil {
					// Return some error
					fmt.Println("ERROR")
				}
				responseData = imageData
				w.Header().Set("Content-Type", "text/base64")
			default:
				http.NotFound(w, req)
				return
		}
		w.Write(responseData)
	}

	http.HandleFunc("/v1/images/", images)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func returnResponse(w http.ResponseWriter, status int, response string){

}