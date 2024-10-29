package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("38xjn6fuphfjmu3")
		if err != nil {
			return err
		}

		// update
		edit_icon := &core.TextField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "e2ctp1go",
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
		collection, err := app.FindCollectionByNameOrId("38xjn6fuphfjmu3")
		if err != nil {
			return err
		}

		// update
		edit_icon := &core.TextField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "e2ctp1go",
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
