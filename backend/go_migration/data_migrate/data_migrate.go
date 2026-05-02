package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"drafty3/go_migration/data_model"
)

// function to create sqlite db for dataset dbs from gorm automigrate and log any errors
func main() {
	// where to make new sqlite db - can be changed as needed
	dsn := "../../db/drafty_new_gorm.db?_pragma=foreign_keys(1)"
	// open new sqlite db using gorm
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}

	// automigrate all models
	if err := db.AutoMigrate(
		&data_model.Alias{},
		&data_model.Click{},
		&data_model.DataType{},
		&data_model.DatabaitCreateType{},
		&data_model.DatabaitNextAction{},
		&data_model.DatabaitTemplateType{},
		&data_model.DoubleClick{},
		&data_model.EditSuggestion{},
		&data_model.EntryType{},
		&data_model.Interaction{},
		&data_model.DatabaitTweet{},
		&data_model.Edit{},
		&data_model.InteractionType{},
		&data_model.Profile{},
		&data_model.RemoveUserData{},
		&data_model.Role{},
		&data_model.SearchType{},
		&data_model.SelectRange{},
		&data_model.Session{},
		&data_model.SuggestionType{},
		&data_model.CopyColumn{},
		&data_model.Search{},
		&data_model.SearchMulti{},
		&data_model.Sort{},
		&data_model.SuggestionTypeValues{},
		&data_model.UniqueId{},
		&data_model.Comments{},
		&data_model.CommentVote{},
		&data_model.CommentsView{},
		&data_model.Databaits{},
		&data_model.DatabaitVisit{},
		&data_model.EditDelRow{},
		&data_model.HelpUs{},
		&data_model.Suggestions{},
		&data_model.Copy{},
		&data_model.EditNewRow{},
		&data_model.Paste{},
		&data_model.SearchGoogle{},
		&data_model.ViewChange{},
		&data_model.Visit{},
		&data_model.Sessions{},
	); err != nil {
		// log any errors
		log.Fatalf("automigrate: %v", err)
	}

	// log success
	log.Println("AutoMigrate complete. SQLite database created: drafty_new_gorm.db")
}
