package img

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"encoding/base64"
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

func GenerateImageData(data []byte) (*Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(data)))
	
	imageObj, format, err := image.Decode(reader)
	if (err != nil) {
		// Return specific error for 422?
		return nil, err
	}
	bounds := imageObj.Bounds()
	imageMetadata := Metadata{
		Width: bounds.Dx(),
		Height: bounds.Dy(),
		Format: format,
		Created_at: time.Now().String(),
	}
	imageMetadata.Size = imageMetadata.Width * imageMetadata.Height * 4 / 1024

	_image := Image{
		Data: string(data),
		Metadata: imageMetadata,
	}

	return &_image, nil
}

func GetImageCutout(img *Image) (*Image, error) {
	fmt.Println("Getting cutout")
	return &Image{}, nil
}