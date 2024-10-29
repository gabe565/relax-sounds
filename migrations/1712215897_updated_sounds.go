package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

var ConvertAfterStart bool

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_file := &core.FileField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "otvjmy0h",
			"name": "file",
			"type": "file",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"mimeTypes": [
					"audio/ogg",
					"audio/mpeg"
				],
				"thumbs": [],
				"maxSelect": 2,
				"maxSize": 5242880,
				"protected": false
			}
		}`), edit_file); err != nil {
			return err
		}
		collection.Fields.Add(edit_file)

		ConvertAfterStart = true
		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_file := &core.FileField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "otvjmy0h",
			"name": "file",
			"type": "file",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"mimeTypes": [
					"audio/ogg",
					"audio/mpeg"
				],
				"thumbs": [],
				"maxSelect": 1,
				"maxSize": 5242880,
				"protected": false
			}
		}`), edit_file); err != nil {
			return err
		}
		collection.Fields.Add(edit_file)

		return app.Save(collection)
	})
}
