package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_file := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "otvjmy0h",
			"name": "file",
			"type": "file",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"maxSize": 5242880,
				"mimeTypes": [
					"audio/ogg",
					"audio/mpeg"
				],
				"thumbs": [],
				"protected": false
			}
		}`), edit_file)
		collection.Schema.AddField(edit_file)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_file := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "otvjmy0h",
			"name": "file",
			"type": "file",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"maxSize": 5242880,
				"mimeTypes": [
					"audio/ogg"
				],
				"thumbs": [],
				"protected": false
			}
		}`), edit_file)
		collection.Schema.AddField(edit_file)

		return dao.SaveCollection(collection)
	})
}
