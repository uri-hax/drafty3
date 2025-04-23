package main

import (
    "encoding/csv"
    "os"
    "strconv"
    "time"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/core"
)

func suggestionsCSVImport(app *pocketbase.PocketBase) error {
    f, _ := os.Open("csv_data/Suggestions.csv")
    defer f.Close()

    rows, _ := csv.NewReader(f).ReadAll()
    headers := rows[0]
    coll, _ := app.FindCollectionByNameOrId("Suggestions")

    return app.RunInTransaction(func(tx core.App) error {
        for _, row := range rows[1:] {
            rec := core.NewRecord(coll)

            for i, val := range row {
                field := headers[i]

                switch field {
					case "idSuggestion":
						rec.Set("id", val)

					case "confidence":
						if val != "" {
							if n, err := strconv.Atoi(val); err == nil {
								rec.Set("confidence", n)
							}
						}

					case "last_updated":
						if t, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
							rec.Set("last_updated", t)
						}

					case "idSuggestionType", "idUniqueID", "idProfile", "suggestion", "active":
						if val != "" {
							rec.Set(field, val)
						}
                }

            }

            tx.SaveNoValidate(rec)
        }
        return nil
    })
}

func uniqueIdsCSVImport(app *pocketbase.PocketBase) error {
    f, _ := os.Open("csv_data/UniqueId.csv")
    defer f.Close()

    rows, _ := csv.NewReader(f).ReadAll()
    headers := rows[0]
    coll, _ := app.FindCollectionByNameOrId("UniqueId")

    return app.RunInTransaction(func(tx core.App) error {
        for _, row := range rows[1:] {
            rec := core.NewRecord(coll)

            for j, val := range row {
                field := headers[j]

                switch field {
                    case "idUniqueID":
                        rec.Set("id", val)

                    case "active":
                        if val != "" {
							rec.Set("active", val)
						}

                    case "notes":
                        rec.Set("notes", val)
                }
            }

            tx.SaveNoValidate(rec)
        }
        return nil
    })
}

func suggestionTypeCSVImport(app *pocketbase.PocketBase) error {
	f, _ := os.Open("csv_data/SuggestionType.csv")
	defer f.Close()

	rows, _ := csv.NewReader(f).ReadAll()
	headers := rows[0]
	coll, _ := app.FindCollectionByNameOrId("SuggestionType")

	return app.RunInTransaction(func(tx core.App) error {
		for _, row := range rows[1:] {
			rec := core.NewRecord(coll)

			for i, val := range row {
				field := headers[i]

				switch field {
                    case "idSuggestionType":
                        rec.Set("id", val)

                    case "idDataType", "name", "regex":
                        rec.Set(field, val)

                    case "isEditable", "isPrivate", "columnOrder":
                        if val != "" {
                            if n, err := strconv.Atoi(val); err == nil {
                                rec.Set(field, n)
                            }
                        }

                    case "isActive", "canBeBlank", "isFreeEdit", "isDate", "isLink", "isCurrency":
                        if val != "" {
                            rec.Set(field, val)
                        }

                    case "makesRowUnique":
                        rec.Set(field, val)
                }
			}

			tx.SaveNoValidate(rec)
		}
		return nil
	})
}