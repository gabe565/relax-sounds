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
		edit_icon := &core.TextField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "49laf206",
			"name": "icon",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_icon)
		collection.Fields.Add(edit_icon)

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_icon := &core.TextField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "49laf206",
			"name": "icon",
			"type": "text",
			"required": true,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": "^fa-"
			}
		}`), edit_icon)
		collection.Fields.Add(edit_icon)

		return app.Save(collection)
	})
}
