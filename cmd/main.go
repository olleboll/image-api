package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/olleboll/images/db"
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

	imageStore, err := db.Connect()

	if err != nil {
		log.Fatal("Could not connect to db")
		return
	}

	images := func(w http.ResponseWriter, req *http.Request) {

		var responseData []byte
		fmt.Println("GOT REQUEST")
		if req.Method == "POST" {
			fmt.Println("POSTING")
			imageData, err := io.ReadAll(req.Body)
			meta, err := imageStore.CreateImage(imageData)

			if err != nil {

			}
			responseData, _ = json.Marshal(meta)
			w.Header().Set("Content-Type", "application/json")

		} else if req.Method == "GET" {
			var images []img.Image
			images, err = imageStore.GetImages()

			if err != nil {
				// Return some error
				fmt.Println("ERROR")
			}
			responseData, _ = json.Marshal(images)
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(responseData)
	}

	http.HandleFunc("/v1/images", images)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func returnResponse(w http.ResponseWriter, status int, response string){

}