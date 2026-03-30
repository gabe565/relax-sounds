package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"authRule": "verified = true",
			"oauth2": {
				"mappedFields": {
					"avatarURL": "avatar",
					"name": "name"
				}
			}
		}`), &collection); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"exceptDomains": null,
			"hidden": false,
			"id": "email3885137012",
			"name": "email",
			"onlyDomains": null,
			"presentable": false,
			"required": true,
			"system": true,
			"type": "email"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"authRule": "",
			"oauth2": {
				"mappedFields": {
					"avatarURL": "",
					"name": ""
				}
			}
		}`), &collection); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"exceptDomains": null,
			"hidden": false,
			"id": "email3885137012",
			"name": "email",
			"onlyDomains": null,
			"presentable": false,
			"required": false,
			"system": true,
			"type": "email"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
