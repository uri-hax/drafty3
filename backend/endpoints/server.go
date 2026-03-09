package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	esession "github.com/labstack/echo-contrib/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"drafty3/endpoints/handler"
)

func registerRoutes(api *echo.Group, db *gorm.DB) {
	// create handlers with dataset db
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
	removeUserDataHandler := handler.NewRemoveUserDataHandler(db)
	roleHandler := handler.NewRoleHandler(db)
	searchTypeHandler := handler.NewSearchTypeHandler(db)
	selectRangeHandler := handler.NewSelectRangeHandler(db)
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

	// Suggestions
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

	// SuggestionType
	api.GET("/suggestiontypes/:name", suggestionTypeHandler.GetSuggestionType)
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
}

func registerUserRoutes(api *echo.Group, usersDB *gorm.DB) {
	profileHandler := handler.NewProfileHandler(usersDB)
	sessionsHandler := handler.NewSessionsHandler(usersDB)

	api.GET("/profiles/:id", profileHandler.GetProfile)
	api.POST("/profiles", profileHandler.CreateProfile)

	api.GET("/sessions/:id", sessionsHandler.GetSessions)
	api.POST("/sessions", sessionsHandler.CreateSessions)
}

func main() {
	e := echo.New()

	e.Use(esession.Middleware(sessions.NewCookieStore([]byte("temp_secret_key"))))

	dbRoot := os.Getenv("DB_ROOT")
	if dbRoot == "" {
		dbRoot = "../db"
	}

	csprofsPath := os.Getenv("DB_PATH_CSPROFS")
	if csprofsPath == "" {
		csprofsPath = filepath.Join(dbRoot, "drafty_new_gorm.db")
	}

	studentsPath := os.Getenv("DB_PATH_STUDENTS")
	if studentsPath == "" {
		studentsPath = filepath.Join(dbRoot, "students_gorm.db")
	}

	usersPath := os.Getenv("DB_PATH_USERS")
	if usersPath == "" {
		usersPath = filepath.Join(dbRoot, "users_gorm.db")
	}

	dbCsprofs, err := gorm.Open(sqlite.Open(csprofsPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect csprofs db:", err)
	}

	dbStudents, err := gorm.Open(sqlite.Open(studentsPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect students db:", err)
	}

	dbUsers, err := gorm.Open(sqlite.Open(usersPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect users db:", err)
	}

	api := e.Group("/api")

	registerRoutes(api.Group("/csprofs"), dbCsprofs)
	registerRoutes(api.Group("/students"), dbStudents)
	registerUserRoutes(api.Group("/users"), dbUsers)

	log.Println("Server running on http://localhost:8081")
	e.Logger.Fatal(e.Start(":8081"))
}