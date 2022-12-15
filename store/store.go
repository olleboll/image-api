package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	_ "bytes"
	_ "github.com/lib/pq"
	"github.com/olleboll/images/img"
)

type ImageStore interface {
	CreateImage(data []byte) (img.Image, error)
	UpdateImage(id int, data []byte) (img.Image, error)
	GetImages() ([]img.Image, error)
	GetImage(id int) (img.Image, error)
	GetImageData(id int, _data *[]byte) (error)
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
	err := db.Ping()
	if err != nil {
		fmt.Println("Error when connecting to db")
		fmt.Println(fmt.Errorf(err.Error()))
		return nil, err
	}
	fmt.Println("ConnectedToDb")

	// Check if things exist
	err = createTable(db)

	if err != nil {
		fmt.Println("Error when creating table")
		fmt.Println(fmt.Errorf(err.Error()))
		return nil, err
	}

	return &ImageDatabase{db: db}, nil
}

func (store *ImageDatabase) CreateImage(data []byte) (img.Image, error) {

	_image, err := img.GenerateImageData(data)

	if (err != nil) {
		// Return 422?
		return img.Image{}, err
	}

	queryString := `INSERT INTO images (data, metadata) VALUES ($1, $2) RETURNING id`
	metaJson, _ := json.Marshal(_image.Metadata)

	var id int
	err = store.db.QueryRow(queryString, _image.Data, metaJson).Scan(&id)

	if (err != nil) {
		// Databse error.. Return 500?
		return img.Image{}, err
	}
	
	return img.Image{
		Id: int(id),
		Metadata: _image.Metadata,
	}, nil
}

func (store *ImageDatabase) UpdateImage(id int, data []byte) (img.Image, error) {

	_image, err := img.GenerateImageData(data)

	if (err != nil) {
		// Return 422?
		return img.Image{}, err
	}

	queryString := `UPDATE images SET data = $1, metadata = $2 WHERE id = $3`
	metaJson, _ := json.Marshal(_image.Metadata)
	_, err = store.db.Exec(queryString, _image.Data, metaJson, id)

	if (err != nil) {
		// Databse error.. Return 500?
		return img.Image{}, err
	}
	
	return img.Image{
		Id: int(id),
		Metadata: _image.Metadata,
	}, nil
}
func (store *ImageDatabase) GetImages() ([]img.Image, error) {
	rows, _ := store.db.Query(`SELECT "id", "metadata" FROM "images"`)
	defer rows.Close()

	results := []img.Image{}

	for rows.Next() {
		var id int
		var metaString string
		err := rows.Scan(&id, &metaString)
		meta := img.Metadata{}
		json.Unmarshal([]byte(metaString), &meta)

		if err != nil {
			panic(err)
		}

		imageMeta := img.Image{
			Id: id,
			Metadata: meta,
		}

		results = append(results, imageMeta)
	}

	return results, nil
}

func (store *ImageDatabase) GetImage(id int) (img.Image, error) {
	rows, _ := store.db.Query(`SELECT "metadata" FROM "images" WHERE "id" = $1`, id)
	defer rows.Close()

	var metaString string
	rows.Next()
	err := rows.Scan(&metaString)
	meta := img.Metadata{}
	json.Unmarshal([]byte(metaString), &meta)

	if err != nil {
		panic(err)
	}

	image := img.Image{
		Id: id,
		Metadata: meta,
	}	

	return image, nil
}
func (store *ImageDatabase) GetImageData(id int, _data *[]byte) (error) {
	rows, _ := store.db.Query(`SELECT "data" FROM "images" WHERE "id" = $1`, id)
	defer rows.Close()

	rows.Next()
	err := rows.Scan(_data)
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