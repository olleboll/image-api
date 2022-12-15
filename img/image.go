package img

import (
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
	Metadata Metadata  `json:"metadata"`
}

func GenerateImageData(data []byte) (*Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(data)))
	imageObj, _, err := image.Decode(reader)
	if (err != nil) {
		// Return specific error for 422?
		return nil, err
	}
	imageRect := imageObj.Bounds()
	imageMetadata := Metadata{
		Size: cap(data),
		Width: imageRect.Max.X - imageRect.Min.X,
		Height: imageRect.Max.Y - imageRect.Min.Y,
		Format: "jpg!?",
		Created_at: time.Now().String(),
	}

	_image := Image{
		Data: string(data),
		Metadata: imageMetadata,
	}

	return &_image, nil
}