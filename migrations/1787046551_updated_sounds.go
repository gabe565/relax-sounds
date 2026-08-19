package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// fill missing old_id values
		var maximum int
		if err := app.RecordQuery(collection.Id).
			Select("COALESCE(MAX(old_id), 0)").
			Row(&maximum); err != nil {
			return err
		}

		var missing []*core.Record
		if err := app.RecordQuery(collection.Id).
			Where(dbx.Or(dbx.HashExp{"old_id": 0}, dbx.HashExp{"old_id": nil})).
			OrderBy("created").
			All(&missing); err != nil {
			return err
		}

		for _, record := range missing {
			maximum++
			record.Set("old_id", maximum)
			if err := app.SaveNoValidate(record); err != nil {
				return err
			}
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX `+"`"+`idx_6qzmvogfg7`+"`"+` ON `+"`"+`sounds`+"`"+` (`+"`"+`short_id`+"`"+`)"
			]
		}`), &collection); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"help": "",
			"hidden": false,
			"id": "f33ff9pv",
			"max": null,
			"min": 0,
			"name": "short_id",
			"onlyInt": true,
			"presentable": false,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("toxjtrsrx7pgzug")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": []
		}`), &collection); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"help": "",
			"hidden": false,
			"id": "f33ff9pv",
			"max": null,
			"min": null,
			"name": "old_id",
			"onlyInt": false,
			"presentable": false,
			"required": false,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
