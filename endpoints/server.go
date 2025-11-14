package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"drafty3/endpoints/handler"

	esession "github.com/labstack/echo-contrib/session" 
	"github.com/gorilla/sessions"                       
)

func main() {
	// read DB path from env or fall back to default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../db/drafty_gorm.db"
	}

	// open DB connection with gorm
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// create echo instance
	e := echo.New()

	// use session middleware
	e.Use(esession.Middleware(sessions.NewCookieStore([]byte("temp_secret_key"))))

	// create handlers with db
	suggestionsHandler := handler.NewSuggestionsHandler(db)
	aliasHandler := handler.NewAliasHandler(db)
	clickHandler := handler.NewClickHandler(db)
	dataTypeHandler := handler.NewDataTypeHandler(db)
	databaitCreateTypeHandler := handler.NewDatabaitCreateTypeHandler(db)
	databaitNextActionHandler := handler.NewDatabaitNextActionHandler(db)
	databaitTemplateTypeHandler := handler.NewDatabaitTemplateTypeHandler(db)
	doubleClickHandler := handler.NewDoubleClickHandler(db)
	editSuggestionHandler := handler.NewEditSuggestionHandler(db)
	entryTypeHandler := handler.NewEntryTypeHandler(db)
	interactionHandler := handler.NewInteractionHandler(db)
	databaitTweetHandler := handler.NewDatabaitTweetHandler(db)
	editHandler := handler.NewEditHandler(db)
	interactionTypeHandler := handler.NewInteractionTypeHandler(db)
	profileHandler := handler.NewProfileHandler(db)
	removeUserDataHandler := handler.NewRemoveUserDataHandler(db)
	roleHandler := handler.NewRoleHandler(db)
	searchTypeHandler := handler.NewSearchTypeHandler(db)
	selectRangeHandler := handler.NewSelectRangeHandler(db)
	sessionHandler := handler.NewSessionHandler(db)
	suggestionTypeHandler := handler.NewSuggestionTypeHandler(db)
	copyColumnHandler := handler.NewCopyColumnHandler(db)
	searchHandler := handler.NewSearchHandler(db)
	searchMultiHandler := handler.NewSearchMultiHandler(db)
	sortHandler := handler.NewSortHandler(db)
	suggestionTypeValuesHandler := handler.NewSuggestionTypeValuesHandler(db)
	uniqueIdHandler := handler.NewUniqueIdHandler(db)
	commentsHandler := handler.NewCommentsHandler(db)
	commentVoteHandler := handler.NewCommentVoteHandler(db)
	commentsViewHandler := handler.NewCommentsViewHandler(db)
	databaitsHandler := handler.NewDatabaitsHandler(db)
	databaitVisitHandler := handler.NewDatabaitVisitHandler(db)
	editDelRowHandler := handler.NewEditDelRowHandler(db)
	helpUsHandler := handler.NewHelpUsHandler(db)
	copyHandler := handler.NewCopyHandler(db)
	editNewRowHandler := handler.NewEditNewRowHandler(db)
	pasteHandler := handler.NewPasteHandler(db)
	searchGoogleHandler := handler.NewSearchGoogleHandler(db)
	viewChangeHandler := handler.NewViewChangeHandler(db)
	visitHandler := handler.NewVisitHandler(db)
	sessionsHandler := handler.NewSessionsHandler(db)

	// set up /api routes
	api := e.Group("/api")

	// Suggestions (at the top, with PUT)
	api.GET("/suggestions/:id", suggestionsHandler.GetSuggestion)
	api.POST("/suggestions", suggestionsHandler.CreateSuggestion)
	api.PUT("/suggestions/:id", suggestionsHandler.UpdateSuggestion)

	// Alias
	api.GET("/alias/:id", aliasHandler.GetAlias)
	api.POST("/alias", aliasHandler.CreateAlias)

	// Click
	api.GET("/clicks/:id", clickHandler.GetClick)
	api.POST("/clicks", clickHandler.CreateClick)

	// DataType
	api.GET("/datatypes/:id", dataTypeHandler.GetDataType)
	api.POST("/datatypes", dataTypeHandler.CreateDataType)

	// DatabaitCreateType
	api.GET("/databaitcreatetypes/:id", databaitCreateTypeHandler.GetDatabaitCreateType)
	api.POST("/databaitcreatetypes", databaitCreateTypeHandler.CreateDatabaitCreateType)

	// DatabaitNextAction
	api.GET("/databaitnextactions/:id", databaitNextActionHandler.GetDatabaitNextAction)
	api.POST("/databaitnextactions", databaitNextActionHandler.CreateDatabaitNextAction)

	// DatabaitTemplateType
	api.GET("/databaittemplatetypes/:id", databaitTemplateTypeHandler.GetDatabaitTemplateType)
	api.POST("/databaittemplatetypes", databaitTemplateTypeHandler.CreateDatabaitTemplateType)

	// DoubleClick
	api.GET("/doubleclicks/:id", doubleClickHandler.GetDoubleClick)
	api.POST("/doubleclicks", doubleClickHandler.CreateDoubleClick)

	// EditSuggestion
	api.GET("/editsuggestion/:id", editSuggestionHandler.GetEditSuggestion)
	api.POST("/editsuggestion", editSuggestionHandler.CreateEditSuggestion)

	// EntryType
	api.GET("/entrytypes/:id", entryTypeHandler.GetEntryType)
	api.POST("/entrytypes", entryTypeHandler.CreateEntryType)

	// Interaction
	api.GET("/interactions/:id", interactionHandler.GetInteraction)
	api.POST("/interactions", interactionHandler.CreateInteraction)

	// DatabaitTweet
	api.GET("/databaittweets/:id", databaitTweetHandler.GetDatabaitTweet)
	api.POST("/databaittweets", databaitTweetHandler.CreateDatabaitTweet)

	// Edit
	api.GET("/edits/:id", editHandler.GetEdit)
	api.POST("/edits", editHandler.CreateEdit)

	// InteractionType
	api.GET("/interactiontypes/:id", interactionTypeHandler.GetInteractionType)
	api.POST("/interactiontypes", interactionTypeHandler.CreateInteractionType)

	// Profile
	api.GET("/profiles/:id", profileHandler.GetProfile)
	api.POST("/profiles", profileHandler.CreateProfile)

	// RemoveUserData
	api.GET("/removeuserdata/:id", removeUserDataHandler.GetRemoveUserData)
	api.POST("/removeuserdata", removeUserDataHandler.CreateRemoveUserData)

	// Role
	api.GET("/roles/:id", roleHandler.GetRole)
	api.POST("/roles", roleHandler.CreateRole)

	// SearchType
	api.GET("/searchtypes/:id", searchTypeHandler.GetSearchType)
	api.POST("/searchtypes", searchTypeHandler.CreateSearchType)

	// SelectRange
	api.GET("/selectranges/:id", selectRangeHandler.GetSelectRange)
	api.POST("/selectranges", selectRangeHandler.CreateSelectRange)

	// Session
	api.GET("/sessions/:id", sessionHandler.GetSession)
	api.POST("/sessions", sessionHandler.CreateSession)

	// SuggestionType
	api.GET("/suggestiontypes/:id", suggestionTypeHandler.GetSuggestionType)
	api.POST("/suggestiontypes", suggestionTypeHandler.CreateSuggestionType)

	// CopyColumn
	api.GET("/copycolumns/:id", copyColumnHandler.GetCopyColumn)
	api.POST("/copycolumns", copyColumnHandler.CreateCopyColumn)

	// Search
	api.GET("/searches/:id", searchHandler.GetSearch)
	api.POST("/searches", searchHandler.CreateSearch)

	// SearchMulti
	api.GET("/searchmultis/:id", searchMultiHandler.GetSearchMulti)
	api.POST("/searchmultis", searchMultiHandler.CreateSearchMulti)

	// Sort
	api.GET("/sorts/:id", sortHandler.GetSort)
	api.POST("/sorts", sortHandler.CreateSort)

	// SuggestionTypeValues
	api.GET("/suggestiontypevalues/:id", suggestionTypeValuesHandler.GetSuggestionTypeValues)
	api.POST("/suggestiontypevalues", suggestionTypeValuesHandler.CreateSuggestionTypeValues)

	// UniqueId
	api.GET("/uniqueids/:id", uniqueIdHandler.GetUniqueId)
	api.POST("/uniqueids", uniqueIdHandler.CreateUniqueId)

	// Comments
	api.GET("/comments/:id", commentsHandler.GetComments)
	api.POST("/comments", commentsHandler.CreateComments)

	// CommentVote
	api.GET("/commentvotes/:id", commentVoteHandler.GetCommentVote)
	api.POST("/commentvotes", commentVoteHandler.CreateCommentVote)

	// CommentsView
	api.GET("/commentsviews/:id", commentsViewHandler.GetCommentsView)
	api.POST("/commentsviews", commentsViewHandler.CreateCommentsView)

	// Databaits
	api.GET("/databaits/:id", databaitsHandler.GetDatabaits)
	api.POST("/databaits", databaitsHandler.CreateDatabaits)

	// DatabaitVisit
	api.GET("/databaitvisits/:id", databaitVisitHandler.GetDatabaitVisit)
	api.POST("/databaitvisits", databaitVisitHandler.CreateDatabaitVisit)

	// EditDelRow
	api.GET("/editdelrows/:id", editDelRowHandler.GetEditDelRow)
	api.POST("/editdelrows", editDelRowHandler.CreateEditDelRow)

	// HelpUs
	api.GET("/helpus/:id", helpUsHandler.GetHelpUs)
	api.POST("/helpus", helpUsHandler.CreateHelpUs)

	// Copy
	api.GET("/copies/:id", copyHandler.GetCopy)
	api.POST("/copies", copyHandler.CreateCopy)

	// EditNewRow
	api.GET("/editnewrows/:id", editNewRowHandler.GetEditNewRow)
	api.POST("/editnewrows", editNewRowHandler.CreateEditNewRow)

	// Paste
	api.GET("/pastes/:id", pasteHandler.GetPaste)
	api.POST("/pastes", pasteHandler.CreatePaste)

	// SearchGoogle
	api.GET("/searchgoogles/:id", searchGoogleHandler.GetSearchGoogle)
	api.POST("/searchgoogles", searchGoogleHandler.CreateSearchGoogle)

	// ViewChange
	api.GET("/viewchanges/:id", viewChangeHandler.GetViewChange)
	api.POST("/viewchanges", viewChangeHandler.CreateViewChange)

	// Visit
	api.GET("/visits/:id", visitHandler.GetVisit)
	api.POST("/visits", visitHandler.CreateVisit)

	// Sessions store (sessions table)
	api.GET("/sessions/:id", sessionsHandler.GetSessions)
	api.POST("/sessions", sessionsHandler.CreateSessions)

	// start server
	log.Println("Server running on http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
