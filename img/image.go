package img

import (
	"bytes"
	"encoding/base64"
	"image"
	gif "image/gif"
	jpeg "image/jpeg"
	png "image/png"
	"strings"
	"time"
)

type Metadata struct {
	Size       int    `json:"size"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Format     string `json:"format"`
	Created_at string `json:"created_at"`
}

type Image struct {
	Id       int      `json:"id"`
	Data     string   `json:"data,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

type Cutout struct {
	X int
	Y int
	W int
	H int
}

func AnalyzeImageData(data *[]byte) (*Image, error) {
	// Create image.Image instace
	imageObj, format, err := decodeImageData(data)
	if err != nil {
		return nil, err
	}

	// Populate metadata
	bounds := imageObj.Bounds()
	imageMetadata := Metadata{
		Size:       len(*data) / 1024,
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
		Format:     format,
		Created_at: time.Now().String(),
	}

	return &Image{
		Data:     string(*data),
		Metadata: imageMetadata,
	}, nil
}

func GetImageCutout(data *[]byte, cutout Cutout) ([]byte, error) {
	// Get original image data
	originalImage, format, err := decodeImageData(data)

	if err != nil {
		return []byte{}, err
	}

	// Create subimage
	imageObj := originalImage.(interface {
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

func decodeImageData(data *[]byte) (image.Image, string, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(*data)))
	return image.Decode(reader)
}
