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
		edit_icon := &schema.SchemaField{}
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
		collection.Schema.AddField(edit_icon)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update
		edit_icon := &schema.SchemaField{}
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
		collection.Schema.AddField(edit_icon)

		return dao.SaveCollection(collection)
	})
}
