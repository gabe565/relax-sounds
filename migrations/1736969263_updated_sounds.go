package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"hidden": false,
			"id": "otvjmy0h",
			"maxSelect": 3,
			"maxSize": 100000000,
			"mimeTypes": [
				"audio/ogg",
				"audio/mpeg",
				"audio/wav",
				"audio/flac"
			],
			"name": "file",
			"presentable": false,
			"protected": false,
			"required": true,
			"system": false,
			"thumbs": null,
			"type": "file"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"hidden": false,
			"id": "otvjmy0h",
			"maxSelect": 2,
			"maxSize": 5242880,
			"mimeTypes": [
				"audio/ogg",
				"audio/mpeg"
			],
			"name": "file",
			"presentable": false,
			"protected": false,
			"required": true,
			"system": false,
			"thumbs": null,
			"type": "file"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
