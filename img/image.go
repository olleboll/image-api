package img

import (
	"image"
	gif "image/gif"
	png "image/png"
	jpeg "image/jpeg"
	"encoding/base64"
	"bytes"
	"strings"
	"time"
)

type Metadata struct {
	Size       int  `json:"size"`
	Width      int  `json:"width"`
	Height     int  `json:"height"`
	Format     string `json:"format"`
	Created_at string `json:"created_at"`
}

type Image struct {
	Id       int  `json:"id"`
	Data     string  `json:"data,omitempty"`
	Metadata Metadata  `json:"metadata,omitempty"`
}

type Cutout struct {
	X int
	Y int
	W int
	H int
}

func GenerateImageData(data *[]byte) (*Image, error) {
	// Create image.Image instace
	imageObj, format, err := getImageData(data)
	if (err != nil) {
		// Return specific error for 422?
		return nil, err
	}

	// Populate metadata
	bounds := imageObj.Bounds()
	imageMetadata := Metadata{
		Width: bounds.Dx(),
		Height: bounds.Dy(),
		Format: format,
		Created_at: time.Now().String(),
	}
	imageMetadata.Size = imageMetadata.Width * imageMetadata.Height * 4 / 1024

	// Create Image struct
	_image := Image{
		Data: string(*data),
		Metadata: imageMetadata,
	}

	return &_image, nil
}

func GetImageCutout(data *[]byte, cutout Cutout) ([]byte, error) {
	// Get original image data
	_image, format, err := getImageData(data)

	// Create subimage
	imageObj := _image.(interface {
        SubImage(r image.Rectangle) image.Image
    }).SubImage(image.Rect(cutout.X, cutout.Y, cutout.W, cutout.H))

	// Depending on format, use the correct Encoder
	buf := new(bytes.Buffer)
	switch format {
		case "jpeg":
			err = jpeg.Encode(buf, imageObj, nil)
		case "png":
			err = png.Encode(buf, imageObj)
		case "gif":
			err = gif.Encode(buf, imageObj, nil)
	}

	return []byte(base64.StdEncoding.EncodeToString(buf.Bytes())), err
	
}

func getImageData(data *[]byte) (image.Image, string, error){
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(*data)))
	return image.Decode(reader)
}