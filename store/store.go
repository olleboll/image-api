package store

import (
	_ "bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/olleboll/images/img"
)

type ImageStore interface {
	CreateImage(data *[]byte) (img.Image, error)
	UpdateImage(id int, data *[]byte) (img.Image, error)
	GetImagesMetadata() ([]img.Image, error)
	GetImageMetadata(id int) (img.Image, error)
	GetImageData(id int, _data *[]byte, cutout *img.Cutout) error
}

type ImageDatabase struct {
	db *sql.DB
}

// Connec to Database.
// Shortcut => Create and migrate if it does not exist
func Connect() (ImageStore, error) {

	var RDS_HOST string = os.Getenv("RDS_HOST")
	var RDS_USER string = os.Getenv("RDS_USER")
	var RDS_PASSWORD string = os.Getenv("RDS_PASSWORD")
	var RDS_DATABASE string = os.Getenv("RDS_DATABASE")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", RDS_USER, RDS_PASSWORD, RDS_HOST, RDS_DATABASE)
	// Connect to database
	db, _ := sql.Open("postgres", connStr)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("ConnectedToDb")

	if err := createTable(db); err != nil {
		return nil, err
	}

	return &ImageDatabase{db: db}, nil
}

func (store *ImageDatabase) CreateImage(data *[]byte) (img.Image, error) {

	image, err := img.AnalyzeImageData(data)

	if err != nil {
		return img.Image{}, err
	}

	queryString := `INSERT INTO images (data, metadata) VALUES ($1, $2) RETURNING id`
	metaJson, _ := json.Marshal(image.Metadata)

	var id int
	if err := store.db.QueryRow(queryString, image.Data, metaJson).Scan(&id); err != nil {
		return img.Image{}, err
	}

	return img.Image{
		Id:       int(id),
		Metadata: image.Metadata,
	}, nil
}

func (store *ImageDatabase) UpdateImage(id int, data *[]byte) (img.Image, error) {

	image, err := img.AnalyzeImageData(data)

	if err != nil {
		return img.Image{}, err
	}

	queryString := `UPDATE images SET data = $1, metadata = $2 WHERE id = $3`
	metaJson, _ := json.Marshal(image.Metadata)
	_, err = store.db.Exec(queryString, image.Data, metaJson, id)

	if err != nil {
		return img.Image{}, err
	}

	return img.Image{
		Id:       int(id),
		Metadata: image.Metadata,
	}, nil
}
func (store *ImageDatabase) GetImagesMetadata() ([]img.Image, error) {
	rows, _ := store.db.Query(`SELECT "id", "metadata" FROM "images"`)
	defer rows.Close()

	results := []img.Image{}

	for rows.Next() {
		var id int
		var metaString string
		if err := rows.Scan(&id, &metaString); err != nil {
			panic(err)
		}
		meta := img.Metadata{}
		json.Unmarshal([]byte(metaString), &meta)

		imageMeta := img.Image{
			Id:       id,
			Metadata: meta,
		}

		results = append(results, imageMeta)
	}

	return results, nil
}

func (store *ImageDatabase) GetImageMetadata(id int) (img.Image, error) {
	rows, _ := store.db.Query(`SELECT "metadata" FROM "images" WHERE "id" = $1`, id)
	defer rows.Close()

	var metaString string
	rows.Next()
	if err := rows.Scan(&metaString); err != nil {
		panic(err)
	}
	meta := img.Metadata{}
	json.Unmarshal([]byte(metaString), &meta)

	image := img.Image{
		Id:       id,
		Metadata: meta,
	}

	return image, nil
}
func (store *ImageDatabase) GetImageData(id int, byteData *[]byte, cutout *img.Cutout) error {
	rows, _ := store.db.Query(`SELECT "data" FROM "images" WHERE "id" = $1`, id)
	defer rows.Close()

	rows.Next()
	err := rows.Scan(byteData)

	if cutout != nil {
		*byteData, err = img.GetImageCutout(byteData, *cutout)
	}

	return err
}

func createTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS images (
		id SERIAL PRIMARY KEY,
		data BYTEA,
		metadata JSON
	)`
	_, err := db.Exec(sql)
	return err
}
