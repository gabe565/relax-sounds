package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "toxjtrsrx7pgzug",
			"created": "2023-04-03 22:02:57.302Z",
			"updated": "2023-04-03 22:02:57.302Z",
			"name": "sounds",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "f4zrmuol",
					"name": "name",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
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
				},
				{
					"system": false,
					"id": "l0o8fgxf",
					"name": "tags",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "38xjn6fuphfjmu3",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": []
					}
				},
				{
					"system": false,
					"id": "otvjmy0h",
					"name": "file",
					"type": "file",
					"required": true,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"maxSize": 5242880,
						"mimeTypes": [
							"audio/ogg"
						],
						"thumbs": []
					}
				},
				{
					"system": false,
					"id": "f33ff9pv",
					"name": "old_id",
					"type": "number",
					"required": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				}
			],
			"indexes": [],
			"listRule": "",
			"viewRule": "",
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
