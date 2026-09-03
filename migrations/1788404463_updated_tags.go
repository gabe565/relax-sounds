package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("38xjn6fuphfjmu3")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"help": "",
			"hidden": false,
			"id": "select1716930793",
			"maxSelect": 0,
			"name": "color",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"pink",
				"orange",
				"green",
				"cyan",
				"slate",
				"blue",
				"gray"
			]
		}`)); err != nil {
			return err
		}

		if err := app.Save(collection); err != nil {
			return err
		}

		colors := map[string]string{
			"Animal":      "orange",
			"Brain Wave":  "cyan",
			"City":        "slate",
			"Music":       "pink",
			"Nature":      "green",
			"Water":       "blue",
			"White Noise": "gray",
		}

		records, err := app.FindAllRecords(collection)
		if err != nil {
			return err
		}

		for _, record := range records {
			color, ok := colors[record.GetString("name")]
			if !ok || record.GetString("color") != "" {
				continue
			}
			record.Set("color", color)
			if err := app.Save(record); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("38xjn6fuphfjmu3")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("select1716930793")

		return app.Save(collection)
	})
}
