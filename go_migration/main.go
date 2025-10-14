package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"drafty3/model" 
)

func main() {
	// Use SQLite DB named drafty3.db with FK support
	dsn := "drafty3.db?_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite: %v", err)
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&model.Alias{},
		&model.Click{},
		&model.DataType{},
		&model.DatabaitCreateType{},
		&model.DatabaitNextAction{},
		&model.DatabaitTemplateType{},
		&model.DoubleClick{},
		&model.EditSuggestion{},
		&model.EntryType{},
		&model.Interaction{},
		&model.DatabaitTweet{},
		&model.Edit{},
		&model.InteractionType{},
		&model.Profile{},
		&model.RemoveUserData{},
		&model.Role{},
		&model.SearchType{},
		&model.SelectRange{},
		&model.Session{},
		&model.SuggestionType{},
		&model.CopyColumn{},
		&model.Search{},
		&model.SearchMulti{},
		&model.Sort{},
		&model.SuggestionTypeValues{},
		&model.UniqueId{},
		&model.Comments{},
		&model.CommentVote{},
		&model.CommentsView{},
		&model.Databaits{},
		&model.DatabaitVisit{},
		&model.EditDelRow{},
		&model.HelpUs{},
		&model.Suggestions{},
		&model.Copy{},
		&model.EditNewRow{},
		&model.Paste{},
		&model.SearchGoogle{},
		&model.ViewChange{},
		&model.Visit{},
		&model.Sessions{},
	); err != nil {
		log.Fatalf("automigrate: %v", err)
	}

	log.Println("AutoMigrate complete. SQLite database created: drafty3.db")
}
