package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_file := &core.FileField{}
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
		collection.Fields.Add(edit_file)

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_file := &core.FileField{}
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
		collection.Fields.Add(edit_file)

		return app.Save(collection)
	})
}
