package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"drafty3/go_migration/model"

	esession "github.com/labstack/echo-contrib/session"
	"github.com/gorilla/sessions"
)

// SUGGESTIONS HANDLER

// SuggestionsHandler holds DB connection
type SuggestionsHandler struct {
	DB *gorm.DB
}

// NewSuggestionsHandler returns a new SuggestionsHandler for the given DB
func NewSuggestionsHandler(db *gorm.DB) *SuggestionsHandler {
	return &SuggestionsHandler{DB: db}
}

// GetSuggestion handles GET /api/suggestions/:id
func (h *SuggestionsHandler) GetSuggestion(c echo.Context) error {
	// lookup row by idSuggestion
	id := c.Param("id")

	// try to find the row and error if can't
	var suggestion model.Suggestions
	if err := h.DB.First(&suggestion, "idSuggestion = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Suggestion not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, suggestion)
}

// CreateSuggestion handles POST /api/suggestions
func (h *SuggestionsHandler) CreateSuggestion(c echo.Context) error {
	// bind request JSON to Suggestions struct
	var suggestion model.Suggestions
	if err := c.Bind(&suggestion); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&suggestion).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create suggestion",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, suggestion)
}

// UpdateSuggestion handles PUT /api/suggestions/:id
func (h *SuggestionsHandler) UpdateSuggestion(c echo.Context) error {
	// get existing row from id parameter in URL
	id := c.Param("id")

	// try to find the row and error if can't
	var existing model.Suggestions
	if err := h.DB.First(&existing, "idSuggestion = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Suggestion not found",
			"id":    id,
		})
	}

	// bind request JSON to existing which mutates in place
	if err := c.Bind(&existing); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// save full struct back to DB
	if err := h.DB.Save(&existing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to update suggestion",
			"detail": err.Error(),
		})
	}

	// return updated row
	return c.JSON(http.StatusOK, existing)
}

// ALIAS HANDLER

// AliasHandler holds DB connection
type AliasHandler struct {
	DB *gorm.DB
}

// NewAliasHandler returns a new AliasHandler for the given DB
func NewAliasHandler(db *gorm.DB) *AliasHandler {
	return &AliasHandler{DB: db}
}

// GetAlias handles GET /api/alias/:id
func (h *AliasHandler) GetAlias(c echo.Context) error {
	// lookup row by idAlias
	id := c.Param("id")

	// try to find the row and error if can't
	var alias model.Alias
	if err := h.DB.First(&alias, "idAlias = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Alias not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, alias)
}

// CreateAlias handles POST /api/alias
func (h *AliasHandler) CreateAlias(c echo.Context) error {
	// bind request JSON to Alias struct
	var alias model.Alias
	if err := c.Bind(&alias); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&alias).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create alias",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, alias)
}

// CLICK HANDLER

// ClickHandler holds DB connection
type ClickHandler struct {
	DB *gorm.DB
}

// NewClickHandler returns a new ClickHandler for the given DB
func NewClickHandler(db *gorm.DB) *ClickHandler {
	return &ClickHandler{DB: db}
}

// GetClick handles GET /api/clicks/:id
func (h *ClickHandler) GetClick(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var click model.Click
	if err := h.DB.First(&click, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Click not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, click)
}

// struct of what we expect from front end with info to make rows in Interaction and Click
type createClickPayload struct {
	IDInteractionType int64   `json:"IDInteractionType"`
	IDSuggestion      int64   `json:"IDSuggestion"`
	RowValues         *string `json:"RowValues"`
}

// CreateClick handles POST /api/clicks
func (h *ClickHandler) CreateClick(c echo.Context) error {
	// read the cookie based session
	sess, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to read cookie session",
			"detail": err.Error(),
		})
	}

	// extract session_id from cookie from whatever form it's in
	var sessionID int64
	if raw, ok := sess.Values["session_id"]; ok {
		switch v := raw.(type) {
		case int64:
			sessionID = v
		case int:
			sessionID = int64(v)
		case float64:
			sessionID = int64(v)
		}
	}

	// bind request JSON filled with info for Interaction and Click
	var payload createClickPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// create Interaction using IDSession from cookie
	interaction := model.Interaction{
		IDSession:         sessionID,
		IDInteractionType: payload.IDInteractionType,
		// Timestamp is default
	}
	if err := h.DB.Create(&interaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction",
			"detail": err.Error(),
		})
	}

	// create Click linked to this Interaction
	click := model.Click{
		IDInteraction: interaction.IDInteraction,
		IDSuggestion:  payload.IDSuggestion,
		RowValues:     payload.RowValues,
	}
	if err := h.DB.Create(&click).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create click",
			"detail": err.Error(),
		})
	}

	// return new row in click
	return c.JSON(http.StatusCreated, click)
}

// DATATYPE HANDLER

// DataTypeHandler holds DB connection
type DataTypeHandler struct {
	DB *gorm.DB
}

// NewDataTypeHandler returns a new DataTypeHandler for the given DB
func NewDataTypeHandler(db *gorm.DB) *DataTypeHandler {
	return &DataTypeHandler{DB: db}
}

// GetDataType handles GET /api/datatypes/:id
func (h *DataTypeHandler) GetDataType(c echo.Context) error {
	// lookup row by idDataType
	id := c.Param("id")

	// try to find the row and error if can't
	var dt model.DataType
	if err := h.DB.First(&dt, "idDataType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DataType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dt)
}

// CreateDataType handles POST /api/datatypes
func (h *DataTypeHandler) CreateDataType(c echo.Context) error {
	// bind request JSON to DataType struct
	var dt model.DataType
	if err := c.Bind(&dt); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dt).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create data type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dt)
}

// DATABAITCREATETYPE HANDLER

// DatabaitCreateTypeHandler holds DB connection
type DatabaitCreateTypeHandler struct {
	DB *gorm.DB
}

// NewDatabaitCreateTypeHandler returns a new DatabaitCreateTypeHandler for the given DB
func NewDatabaitCreateTypeHandler(db *gorm.DB) *DatabaitCreateTypeHandler {
	return &DatabaitCreateTypeHandler{DB: db}
}

// GetDatabaitCreateType handles GET /api/databaitcreatetypes/:id
func (h *DatabaitCreateTypeHandler) GetDatabaitCreateType(c echo.Context) error {
	// lookup row by idDatabaitCreateType
	id := c.Param("id")

	// try to find the row and error if can't
	var dct model.DatabaitCreateType
	if err := h.DB.First(&dct, "idDatabaitCreateType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DatabaitCreateType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dct)
}

// CreateDatabaitCreateType handles POST /api/databaitcreatetypes
func (h *DatabaitCreateTypeHandler) CreateDatabaitCreateType(c echo.Context) error {
	// bind request JSON to DatabaitCreateType struct
	var dct model.DatabaitCreateType
	if err := c.Bind(&dct); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dct).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create databait create type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dct)
}

// DATABAITNEXTACTION HANDLER

// DatabaitNextActionHandler holds DB connection
type DatabaitNextActionHandler struct {
	DB *gorm.DB
}

// NewDatabaitNextActionHandler returns a new DatabaitNextActionHandler for the given DB
func NewDatabaitNextActionHandler(db *gorm.DB) *DatabaitNextActionHandler {
	return &DatabaitNextActionHandler{DB: db}
}

// GetDatabaitNextAction handles GET /api/databaitnextactions/:id
func (h *DatabaitNextActionHandler) GetDatabaitNextAction(c echo.Context) error {
	// lookup row by idDatabaitNextAction
	id := c.Param("id")

	// try to find the row and error if can't
	var dna model.DatabaitNextAction
	if err := h.DB.First(&dna, "idDatabaitNextAction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DatabaitNextAction not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dna)
}

// CreateDatabaitNextAction handles POST /api/databaitnextactions
func (h *DatabaitNextActionHandler) CreateDatabaitNextAction(c echo.Context) error {
	// bind request JSON to DatabaitNextAction struct
	var dna model.DatabaitNextAction
	if err := c.Bind(&dna); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dna).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create databait next action",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dna)
}

// DATABAITTEMPLATETYPE HANDLER

// DatabaitTemplateTypeHandler holds DB connection
type DatabaitTemplateTypeHandler struct {
	DB *gorm.DB
}

// NewDatabaitTemplateTypeHandler returns a new DatabaitTemplateTypeHandler for the given DB
func NewDatabaitTemplateTypeHandler(db *gorm.DB) *DatabaitTemplateTypeHandler {
	return &DatabaitTemplateTypeHandler{DB: db}
}

// GetDatabaitTemplateType handles GET /api/databaittemplatetypes/:id
func (h *DatabaitTemplateTypeHandler) GetDatabaitTemplateType(c echo.Context) error {
	// lookup row by idDatabaitTemplateType
	id := c.Param("id")

	// try to find the row and error if can't
	var dtt model.DatabaitTemplateType
	if err := h.DB.First(&dtt, "idDatabaitTemplateType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DatabaitTemplateType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dtt)
}

// CreateDatabaitTemplateType handles POST /api/databaittemplatetypes
func (h *DatabaitTemplateTypeHandler) CreateDatabaitTemplateType(c echo.Context) error {
	// bind request JSON to DatabaitTemplateType struct
	var dtt model.DatabaitTemplateType
	if err := c.Bind(&dtt); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dtt).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create databait template type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dtt)
}

// DOUBLECLICK HANDLER

// DoubleClickHandler holds DB connection
type DoubleClickHandler struct {
	DB *gorm.DB
}

// NewDoubleClickHandler returns a new DoubleClickHandler for the given DB
func NewDoubleClickHandler(db *gorm.DB) *DoubleClickHandler {
	return &DoubleClickHandler{DB: db}
}

// GetDoubleClick handles GET /api/doubleclicks/:id
func (h *DoubleClickHandler) GetDoubleClick(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var dc model.DoubleClick
	if err := h.DB.First(&dc, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DoubleClick not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dc)
}

// CreateDoubleClick handles POST /api/doubleclicks
func (h *DoubleClickHandler) CreateDoubleClick(c echo.Context) error {
	// bind request JSON to DoubleClick struct
	var dc model.DoubleClick
	if err := c.Bind(&dc); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create double click",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dc)
}

// EDITSUGGESTION HANDLER

// EditSuggestionHandler holds DB connection
type EditSuggestionHandler struct {
	DB *gorm.DB
}

// NewEditSuggestionHandler returns a new EditSuggestionHandler for the given DB
func NewEditSuggestionHandler(db *gorm.DB) *EditSuggestionHandler {
	return &EditSuggestionHandler{DB: db}
}

// GetEditSuggestion handles GET /api/editsuggestion/:id
func (h *EditSuggestionHandler) GetEditSuggestion(c echo.Context) error {
	// lookup row by idEdit
	id := c.Param("id")

	// try to find the row and error if can't
	var es model.EditSuggestion
	if err := h.DB.First(&es, "idEdit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "EditSuggestion not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, es)
}

// CreateEditSuggestion handles POST /api/editsuggestion
func (h *EditSuggestionHandler) CreateEditSuggestion(c echo.Context) error {
	// bind request JSON to EditSuggestion struct
	var es model.EditSuggestion
	if err := c.Bind(&es); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&es).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit suggestion",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, es)
}


// ENTRYTYPE HANDLER

// EntryTypeHandler holds DB connection
type EntryTypeHandler struct {
	DB *gorm.DB
}

// NewEntryTypeHandler returns a new EntryTypeHandler for the given DB
func NewEntryTypeHandler(db *gorm.DB) *EntryTypeHandler {
	return &EntryTypeHandler{DB: db}
}

// GetEntryType handles GET /api/entrytypes/:id
func (h *EntryTypeHandler) GetEntryType(c echo.Context) error {
	// lookup row by idEntryType
	id := c.Param("id")

	// try to find the row and error if can't
	var et model.EntryType
	if err := h.DB.First(&et, "idEntryType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "EntryType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, et)
}

// CreateEntryType handles POST /api/entrytypes
func (h *EntryTypeHandler) CreateEntryType(c echo.Context) error {
	// bind request JSON to EntryType struct
	var et model.EntryType
	if err := c.Bind(&et); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&et).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create entry type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, et)
}

// INTERACTION HANDLER

// InteractionHandler holds DB connection
type InteractionHandler struct {
	DB *gorm.DB
}

// NewInteractionHandler returns a new InteractionHandler for the given DB
func NewInteractionHandler(db *gorm.DB) *InteractionHandler {
	return &InteractionHandler{DB: db}
}

// GetInteraction handles GET /api/interactions/:id
func (h *InteractionHandler) GetInteraction(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var interaction model.Interaction
	if err := h.DB.First(&interaction, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Interaction not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, interaction)
}

// CreateInteraction handles POST /api/interactions
func (h *InteractionHandler) CreateInteraction(c echo.Context) error {
	// bind request JSON to Interaction struct
	var interaction model.Interaction
	if err := c.Bind(&interaction); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&interaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, interaction)
}

// DATABAITTWEET HANDLER

// DatabaitTweetHandler holds DB connection
type DatabaitTweetHandler struct {
	DB *gorm.DB
}

// NewDatabaitTweetHandler returns a new DatabaitTweetHandler for the given DB
func NewDatabaitTweetHandler(db *gorm.DB) *DatabaitTweetHandler {
	return &DatabaitTweetHandler{DB: db}
}

// GetDatabaitTweet handles GET /api/databaittweets/:id
func (h *DatabaitTweetHandler) GetDatabaitTweet(c echo.Context) error {
	// lookup row by idDatabaitTweet
	id := c.Param("id")

	// try to find the row and error if can't
	var dt model.DatabaitTweet
	if err := h.DB.First(&dt, "idDatabaitTweet = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DatabaitTweet not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dt)
}

// CreateDatabaitTweet handles POST /api/databaittweets
func (h *DatabaitTweetHandler) CreateDatabaitTweet(c echo.Context) error {
	// bind request JSON to DatabaitTweet struct
	var dt model.DatabaitTweet
	if err := c.Bind(&dt); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dt).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create databait tweet",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dt)
}

// EDIT HANDLER

// EditHandler holds DB connection
type EditHandler struct {
	DB *gorm.DB
}

// NewEditHandler returns a new EditHandler for the given DB
func NewEditHandler(db *gorm.DB) *EditHandler {
	return &EditHandler{DB: db}
}

// GetEdit handles GET /api/edits/:id
func (h *EditHandler) GetEdit(c echo.Context) error {
	// lookup row by idEdit
	id := c.Param("id")

	// try to find the row and error if can't
	var edit model.Edit
	if err := h.DB.First(&edit, "idEdit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Edit not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, edit)
}

// struct of what we expect from front end with info to make rows in Interaction and Click
type createEditPayload struct {
	IDInteractionType int64   `json:"IDInteractionType"`
	IDEntryType       int64   `json:"IDEntryType"`
	Mode              string  `json:"Mode"`           
	IsCorrect         int64  `json:"IsCorrect"`        
}

// CreateEdit handles POST /api/edits
func (h *EditHandler) CreateEdit(c echo.Context) error {
	// read the cookie based session
	sess, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to read cookie session",
			"detail": err.Error(),
		})
	}

	// extract session_id from cookie from whatever form it's in
	var sessionID int64
	if raw, ok := sess.Values["session_id"]; ok {
		switch v := raw.(type) {
		case int64:
			sessionID = v
		case int:
			sessionID = int64(v)
		case float64:
			sessionID = int64(v)
		}
	}

	// bind request JSON filled with info for Interaction and Edit
	var payload createEditPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// create Interaction using IDSession from cookie
	interaction := model.Interaction{
		IDSession:         sessionID,
		IDInteractionType: payload.IDInteractionType,
		// Timestamp is default
	}
	if err := h.DB.Create(&interaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction",
			"detail": err.Error(),
		})
	}

	// create Edit linked to this Interaction
	edit := model.Edit{
		IDInteraction: interaction.IDInteraction,
		IDEntryType:   payload.IDEntryType,
		Mode:          payload.Mode,
		IsCorrect:     payload.IsCorrect,
	}
	if err := h.DB.Create(&edit).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit",
			"detail": err.Error(),
		})
	}

	// return new row in edit
	return c.JSON(http.StatusCreated, edit)
}

// INTERACTIONTYPE HANDLER

// InteractionTypeHandler holds DB connection
type InteractionTypeHandler struct {
	DB *gorm.DB
}

// NewInteractionTypeHandler returns a new InteractionTypeHandler for the given DB
func NewInteractionTypeHandler(db *gorm.DB) *InteractionTypeHandler {
	return &InteractionTypeHandler{DB: db}
}

// GetInteractionType handles GET /api/interactiontypes/:id
func (h *InteractionTypeHandler) GetInteractionType(c echo.Context) error {
	// lookup row by idInteractionType
	id := c.Param("id")

	// try to find the row and error if can't
	var it model.InteractionType
	if err := h.DB.First(&it, "idInteractionType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "InteractionType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, it)
}

// CreateInteractionType handles POST /api/interactiontypes
func (h *InteractionTypeHandler) CreateInteractionType(c echo.Context) error {
	// bind request JSON to InteractionType struct
	var it model.InteractionType
	if err := c.Bind(&it); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&it).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, it)
}

// PROFILE HANDLER

// ProfileHandler holds DB connection
type ProfileHandler struct {
	DB *gorm.DB
}

// NewProfileHandler returns a new ProfileHandler for the given DB
func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{DB: db}
}

// GetProfile handles GET /api/profiles/:id
func (h *ProfileHandler) GetProfile(c echo.Context) error {
	// lookup row by idProfile
	id := c.Param("id")

	// try to find the row and error if can't
	var profile model.Profile
	if err := h.DB.First(&profile, "idProfile = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Profile not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, profile)
}

// CreateProfile handles POST /api/profiles
func (h *ProfileHandler) CreateProfile(c echo.Context) error {
	// bind request JSON to Profile struct
	var profile model.Profile
	if err := c.Bind(&profile); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&profile).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create profile",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, profile)
}

// REMOVEUSERDATA HANDLER

// RemoveUserDataHandler holds DB connection
type RemoveUserDataHandler struct {
	DB *gorm.DB
}

// NewRemoveUserDataHandler returns a new RemoveUserDataHandler for the given DB
func NewRemoveUserDataHandler(db *gorm.DB) *RemoveUserDataHandler {
	return &RemoveUserDataHandler{DB: db}
}

// GetRemoveUserData handles GET /api/removeuserdata/:id
func (h *RemoveUserDataHandler) GetRemoveUserData(c echo.Context) error {
	// lookup row by id_removeuserdata
	id := c.Param("id")

	// try to find the row and error if can't
	var rud model.RemoveUserData
	if err := h.DB.First(&rud, "id_removeuserdata = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "RemoveUserData not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, rud)
}

// CreateRemoveUserData handles POST /api/removeuserdata
func (h *RemoveUserDataHandler) CreateRemoveUserData(c echo.Context) error {
	// bind request JSON to RemoveUserData struct
	var rud model.RemoveUserData
	if err := c.Bind(&rud); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&rud).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create remove user data",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, rud)
}

// ROLE HANDLER

// RoleHandler holds DB connection
type RoleHandler struct {
	DB *gorm.DB
}

// NewRoleHandler returns a new RoleHandler for the given DB
func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

// GetRole handles GET /api/roles/:id
func (h *RoleHandler) GetRole(c echo.Context) error {
	// lookup row by idRole
	id := c.Param("id")

	// try to find the row and error if can't
	var role model.Role
	if err := h.DB.First(&role, "idRole = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Role not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, role)
}

// CreateRole handles POST /api/roles
func (h *RoleHandler) CreateRole(c echo.Context) error {
	// bind request JSON to Role struct
	var role model.Role
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&role).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create role",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, role)
}

// SEARCHTYPE HANDLER

// SearchTypeHandler holds DB connection
type SearchTypeHandler struct {
	DB *gorm.DB
}

// NewSearchTypeHandler returns a new SearchTypeHandler for the given DB
func NewSearchTypeHandler(db *gorm.DB) *SearchTypeHandler {
	return &SearchTypeHandler{DB: db}
}

// GetSearchType handles GET /api/searchtypes/:id
func (h *SearchTypeHandler) GetSearchType(c echo.Context) error {
	// lookup row by idSearchType
	id := c.Param("id")

	// try to find the row and error if can't
	var st model.SearchType
	if err := h.DB.First(&st, "idSearchType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SearchType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, st)
}

// CreateSearchType handles POST /api/searchtypes
func (h *SearchTypeHandler) CreateSearchType(c echo.Context) error {
	// bind request JSON to SearchType struct
	var st model.SearchType
	if err := c.Bind(&st); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&st).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create search type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, st)
}

// SELECTRANGE HANDLER

// SelectRangeHandler holds DB connection
type SelectRangeHandler struct {
	DB *gorm.DB
}

// NewSelectRangeHandler returns a new SelectRangeHandler for the given DB
func NewSelectRangeHandler(db *gorm.DB) *SelectRangeHandler {
	return &SelectRangeHandler{DB: db}
}

// GetSelectRange handles GET /api/selectranges/:id
func (h *SelectRangeHandler) GetSelectRange(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var sr model.SelectRange
	if err := h.DB.First(&sr, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SelectRange not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, sr)
}

// CreateSelectRange handles POST /api/selectranges
func (h *SelectRangeHandler) CreateSelectRange(c echo.Context) error {
	// bind request JSON to SelectRange struct
	var sr model.SelectRange
	if err := c.Bind(&sr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&sr).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create select range",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, sr)
}

// SUGGESTIONTYPE HANDLER

// SuggestionTypeHandler holds DB connection
type SuggestionTypeHandler struct {
	DB *gorm.DB
}

// NewSuggestionTypeHandler returns a new SuggestionTypeHandler for the given DB
func NewSuggestionTypeHandler(db *gorm.DB) *SuggestionTypeHandler {
	return &SuggestionTypeHandler{DB: db}
}

// GetSuggestionType handles GET /api/suggestiontypes/:id
func (h *SuggestionTypeHandler) GetSuggestionType(c echo.Context) error {
	// lookup row by idSuggestionType
	id := c.Param("id")

	// try to find the row and error if can't
	var st model.SuggestionType
	if err := h.DB.First(&st, "idSuggestionType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SuggestionType not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, st)
}

// CreateSuggestionType handles POST /api/suggestiontypes
func (h *SuggestionTypeHandler) CreateSuggestionType(c echo.Context) error {
	// bind request JSON to SuggestionType struct
	var st model.SuggestionType
	if err := c.Bind(&st); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&st).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create suggestion type",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, st)
}
 
// COPYCOLUMN HANDLER

// CopyColumnHandler holds DB connection
type CopyColumnHandler struct {
	DB *gorm.DB
}

// NewCopyColumnHandler returns a new CopyColumnHandler for the given DB
func NewCopyColumnHandler(db *gorm.DB) *CopyColumnHandler {
	return &CopyColumnHandler{DB: db}
}

// GetCopyColumn handles GET /api/copycolumns/:id
func (h *CopyColumnHandler) GetCopyColumn(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var cc model.CopyColumn
	if err := h.DB.First(&cc, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "CopyColumn not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, cc)
}

// CreateCopyColumn handles POST /api/copycolumns
func (h *CopyColumnHandler) CreateCopyColumn(c echo.Context) error {
	// bind request JSON to CopyColumn struct
	var cc model.CopyColumn
	if err := c.Bind(&cc); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&cc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create copy column",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, cc)
}

// SEARCH HANDLER

// SearchHandler holds DB connection
type SearchHandler struct {
	DB *gorm.DB
}

// NewSearchHandler returns a new SearchHandler for the given DB
func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{DB: db}
}

// GetSearch handles GET /api/searches/:id
func (h *SearchHandler) GetSearch(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var search model.Search
	if err := h.DB.First(&search, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Search not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, search)
}

// struct of what we expect
type createSearchPayload struct {
	IDInteractionType int64  `json:"IDInteractionType"`
	IDSuggestionType  int64  `json:"IDSuggestionType"`
	IDSearchType      int64  `json:"IDSearchType"`
	IsPartial         int64  `json:"IsPartial"`
	IsMulti           int64  `json:"IsMulti"`
	IsFromURL         int64  `json:"IsFromURL"`
	Value             string `json:"Value"`
	MatchedValues     string `json:"MatchedValues"`
}

// CreateSearch handles POST /api/searches
func (h *SearchHandler) CreateSearch(c echo.Context) error {
	// read the cookie based session
	sess, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to read cookie session",
			"detail": err.Error(),
		})
	}

	// extract session_id from cookie from whatever form it's in
	var sessionID int64
	if raw, ok := sess.Values["session_id"]; ok {
		switch v := raw.(type) {
		case int64:
			sessionID = v
		case int:
			sessionID = int64(v)
		case float64:
			sessionID = int64(v)
		}
	}

	// bind request JSON into payload
	var payload createSearchPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// create Interaction for this search
	interaction := model.Interaction{
		IDSession:         sessionID,
		IDInteractionType: payload.IDInteractionType,
		// Timestamp is defaulted in the db
	}
	if err := h.DB.Create(&interaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction",
			"detail": err.Error(),
		})
	}

	// create Search linked to this Interaction
	val := payload.Value
	search := model.Search{
		IDInteraction:    interaction.IDInteraction,
		IDSuggestionType: payload.IDSuggestionType,
		IDSearchType:     payload.IDSearchType,
		IsPartial:        payload.IsPartial,
		IsMulti:          payload.IsMulti,
		IsFromURL:        payload.IsFromURL,
		Value:            &val,
		MatchedValues:    []byte(payload.MatchedValues),
	}
	if err := h.DB.Create(&search).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create search",
			"detail": err.Error(),
		})
	}

	// if IsMulti is true then create SearchMulti row
	if payload.IsMulti != 0 {
		smVal := payload.Value
		searchMulti := model.SearchMulti{
			IDInteraction:    interaction.IDInteraction,
			IDSuggestionType: payload.IDSuggestionType,
			IDSearchType:     payload.IDSearchType,
			Value:            &smVal,
		}
		if err := h.DB.Create(&searchMulti).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error":  "failed to create search multi",
				"detail": err.Error(),
			})
		}
	}

	// return created Search row
	return c.JSON(http.StatusCreated, search)
}

// SEARCHMULTI HANDLER

// SearchMultiHandler holds DB connection
type SearchMultiHandler struct {
	DB *gorm.DB
}

// NewSearchMultiHandler returns a new SearchMultiHandler for the given DB
func NewSearchMultiHandler(db *gorm.DB) *SearchMultiHandler {
	return &SearchMultiHandler{DB: db}
}

// GetSearchMulti handles GET /api/searchmultis/:id
func (h *SearchMultiHandler) GetSearchMulti(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var sm model.SearchMulti
	if err := h.DB.First(&sm, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SearchMulti not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, sm)
}

// CreateSearchMulti handles POST /api/searchmultis
func (h *SearchMultiHandler) CreateSearchMulti(c echo.Context) error {
	// bind request JSON to SearchMulti struct
	var sm model.SearchMulti
	if err := c.Bind(&sm); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&sm).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create search multi",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, sm)
}

// SORT HANDLER

// SortHandler holds DB connection
type SortHandler struct {
	DB *gorm.DB
}

// NewSortHandler returns a new SortHandler for the given DB
func NewSortHandler(db *gorm.DB) *SortHandler {
	return &SortHandler{DB: db}
}

// GetSort handles GET /api/sorts/:id
func (h *SortHandler) GetSort(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var sort model.Sort
	if err := h.DB.First(&sort, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Sort not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, sort)
}

// CreateSort handles POST /api/sorts
func (h *SortHandler) CreateSort(c echo.Context) error {
	// bind request JSON to Sort struct
	var sort model.Sort
	if err := c.Bind(&sort); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&sort).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create sort",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, sort)
}

// SUGGESTIONTYPEVALUES HANDLER

// SuggestionTypeValuesHandler holds DB connection
type SuggestionTypeValuesHandler struct {
	DB *gorm.DB
}

// NewSuggestionTypeValuesHandler returns a new SuggestionTypeValuesHandler for the given DB
func NewSuggestionTypeValuesHandler(db *gorm.DB) *SuggestionTypeValuesHandler {
	return &SuggestionTypeValuesHandler{DB: db}
}

// GetSuggestionTypeValues handles GET /api/suggestiontypevalues/:id
func (h *SuggestionTypeValuesHandler) GetSuggestionTypeValues(c echo.Context) error {
	// lookup row by idSuggestionType
	id := c.Param("id")

	// try to find the row and error if can't
	var stv model.SuggestionTypeValues
	if err := h.DB.First(&stv, "idSuggestionType = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SuggestionTypeValues not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, stv)
}

// CreateSuggestionTypeValues handles POST /api/suggestiontypevalues
func (h *SuggestionTypeValuesHandler) CreateSuggestionTypeValues(c echo.Context) error {
	// bind request JSON to SuggestionTypeValues struct
	var stv model.SuggestionTypeValues
	if err := c.Bind(&stv); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&stv).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create suggestion type values",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, stv)
}

// UNIQUEID HANDLER

// UniqueIdHandler holds DB connection
type UniqueIdHandler struct {
	DB *gorm.DB
}

// NewUniqueIdHandler returns a new UniqueIdHandler for the given DB
func NewUniqueIdHandler(db *gorm.DB) *UniqueIdHandler {
	return &UniqueIdHandler{DB: db}
}

// GetUniqueId handles GET /api/uniqueids/:id
func (h *UniqueIdHandler) GetUniqueId(c echo.Context) error {
	// lookup row by idUniqueID
	id := c.Param("id")

	// try to find the row and error if can't
	var uid model.UniqueId
	if err := h.DB.First(&uid, "idUniqueID = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "UniqueId not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, uid)
}

// CreateUniqueId handles POST /api/uniqueids
func (h *UniqueIdHandler) CreateUniqueId(c echo.Context) error {
	// bind request JSON to UniqueId struct
	var uid model.UniqueId
	if err := c.Bind(&uid); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&uid).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create unique id",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, uid)
}

// COMMENTS HANDLER

// CommentsHandler holds DB connection
type CommentsHandler struct {
	DB *gorm.DB
}

// NewCommentsHandler returns a new CommentsHandler for the given DB
func NewCommentsHandler(db *gorm.DB) *CommentsHandler {
	return &CommentsHandler{DB: db}
}

// GetComments handles GET /api/comments/:id
func (h *CommentsHandler) GetComments(c echo.Context) error {
	// lookup row by idComment
	id := c.Param("id")

	// try to find the row and error if can't
	var comment model.Comments
	if err := h.DB.First(&comment, "idComment = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Comments not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, comment)
}

// CreateComments handles POST /api/comments
func (h *CommentsHandler) CreateComments(c echo.Context) error {
	// bind request JSON to Comments struct
	var comment model.Comments
	if err := c.Bind(&comment); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&comment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create comments",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, comment)
}

// COMMENTVOTE HANDLER

// CommentVoteHandler holds DB connection
type CommentVoteHandler struct {
	DB *gorm.DB
}

// NewCommentVoteHandler returns a new CommentVoteHandler for the given DB
func NewCommentVoteHandler(db *gorm.DB) *CommentVoteHandler {
	return &CommentVoteHandler{DB: db}
}

// GetCommentVote handles GET /api/commentvotes/:id
func (h *CommentVoteHandler) GetCommentVote(c echo.Context) error {
	// lookup row by idCommentVote
	id := c.Param("id")

	// try to find the row and error if can't
	var cv model.CommentVote
	if err := h.DB.First(&cv, "idCommentVote = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "CommentVote not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, cv)
}

// CreateCommentVote handles POST /api/commentvotes
func (h *CommentVoteHandler) CreateCommentVote(c echo.Context) error {
	// bind request JSON to CommentVote struct
	var cv model.CommentVote
	if err := c.Bind(&cv); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&cv).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create comment vote",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, cv)
}

// COMMENTSVIEW HANDLER

// CommentsViewHandler holds DB connection
type CommentsViewHandler struct {
	DB *gorm.DB
}

// NewCommentsViewHandler returns a new CommentsViewHandler for the given DB
func NewCommentsViewHandler(db *gorm.DB) *CommentsViewHandler {
	return &CommentsViewHandler{DB: db}
}

// GetCommentsView handles GET /api/commentsviews/:id
func (h *CommentsViewHandler) GetCommentsView(c echo.Context) error {
	// lookup row by idCommentsView
	id := c.Param("id")

	// try to find the row and error if can't
	var cv model.CommentsView
	if err := h.DB.First(&cv, "idCommentsView = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "CommentsView not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, cv)
}

// CreateCommentsView handles POST /api/commentsviews
func (h *CommentsViewHandler) CreateCommentsView(c echo.Context) error {
	// bind request JSON to CommentsView struct
	var cv model.CommentsView
	if err := c.Bind(&cv); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&cv).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create comments view",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, cv)
}

// DATABAITS HANDLER

// DatabaitsHandler holds DB connection
type DatabaitsHandler struct {
	DB *gorm.DB
}

// NewDatabaitsHandler returns a new DatabaitsHandler for the given DB
func NewDatabaitsHandler(db *gorm.DB) *DatabaitsHandler {
	return &DatabaitsHandler{DB: db}
}

// GetDatabaits handles GET /api/databaits/:id
func (h *DatabaitsHandler) GetDatabaits(c echo.Context) error {
	// lookup row by idDatabait
	id := c.Param("id")

	// try to find the row and error if can't
	var d model.Databaits
	if err := h.DB.First(&d, "idDatabait = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Databaits not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, d)
}

// CreateDatabaits handles POST /api/databaits
func (h *DatabaitsHandler) CreateDatabaits(c echo.Context) error {
	// bind request JSON to Databaits struct
	var d model.Databaits
	if err := c.Bind(&d); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&d).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create databaits",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, d)
}

// DATABAITVISIT HANDLER

// DatabaitVisitHandler holds DB connection
type DatabaitVisitHandler struct {
	DB *gorm.DB
}

// NewDatabaitVisitHandler returns a new DatabaitVisitHandler for the given DB
func NewDatabaitVisitHandler(db *gorm.DB) *DatabaitVisitHandler {
	return &DatabaitVisitHandler{DB: db}
}

// GetDatabaitVisit handles GET /api/databaitvisits/:id
func (h *DatabaitVisitHandler) GetDatabaitVisit(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var dv model.DatabaitVisit
	if err := h.DB.First(&dv, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "DatabaitVisit not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, dv)
}

// CreateDatabaitVisit handles POST /api/databaitvisits
func (h *DatabaitVisitHandler) CreateDatabaitVisit(c echo.Context) error {
	// bind request JSON to DatabaitVisit struct
	var dv model.DatabaitVisit
	if err := c.Bind(&dv); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&dv).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create databait visit",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, dv)
}

// EDITDELROW HANDLER

// EditDelRowHandler holds DB connection
type EditDelRowHandler struct {
	DB *gorm.DB
}

// NewEditDelRowHandler returns a new EditDelRowHandler for the given DB
func NewEditDelRowHandler(db *gorm.DB) *EditDelRowHandler {
	return &EditDelRowHandler{DB: db}
}

// GetEditDelRow handles GET /api/editdelrows/:id
func (h *EditDelRowHandler) GetEditDelRow(c echo.Context) error {
	// lookup row by idEdit
	id := c.Param("id")

	// try to find the row and error if can't
	var edr model.EditDelRow
	if err := h.DB.First(&edr, "idEdit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "EditDelRow not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, edr)
}

// struct of what we expect from front end with info to make rows in Interaction and Edit and EditDelRow
type createEditDelRowPayload struct {
	IDInteractionType int64   `json:"IDInteractionType"`
	IDEntryType       int64   `json:"IDEntryType"`
	Mode              string  `json:"Mode"`
	IsCorrect         int64   `json:"IsCorrect"`
	Comment           string  `json:"Comment"`
}

// CreateEditDelRow handles POST /api/editdelrows
func (h *EditDelRowHandler) CreateEditDelRow(c echo.Context) error {
	// read the cookie based session
	sess, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to read cookie session",
			"detail": err.Error(),
		})
	}

	// extract session_id from cookie from whatever form it's in
	var sessionID int64
	if raw, ok := sess.Values["session_id"]; ok {
		switch v := raw.(type) {
		case int64:
			sessionID = v
		case int:
			sessionID = int64(v)
		case float64:
			sessionID = int64(v)
		}
	}

	// bind request JSON filled with info for Interaction and Edit and EditDelRow
	var payload createEditDelRowPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// create Interaction using IDSession from cookie
	interaction := model.Interaction{
		IDSession:         sessionID,
		IDInteractionType: payload.IDInteractionType,
		// Timestamp is default
	}
	if err := h.DB.Create(&interaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction",
			"detail": err.Error(),
		})
	}

	// create Edit linked to this Interaction
	edit := model.Edit{
		IDInteraction: interaction.IDInteraction,
		IDEntryType:   payload.IDEntryType,
		Mode:          payload.Mode,
		IsCorrect:     payload.IsCorrect,
	}
	if err := h.DB.Create(&edit).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit",
			"detail": err.Error(),
		})
	}

	// create a UniqueId row
    newUnique := model.UniqueId{
        Active: 1, 
        Notes:  nil,
    }
    if err := h.DB.Create(&newUnique).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "error":  "failed to create unique id",
            "detail": err.Error(),
        })
    }

    // create EditDelRow linked to Edit and UniqueId
    edr := model.EditDelRow{
        IDEdit:     edit.IDEdit,
        IDUniqueID: newUnique.IDUniqueID,
        Comment:    payload.Comment,
    }
	if err := h.DB.Create(&edr).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit del row",
			"detail": err.Error(),
		})
	}

	// return new row in edit del row
	return c.JSON(http.StatusCreated, edr)
}

// HELPUS HANDLER

// HelpUsHandler holds DB connection
type HelpUsHandler struct {
	DB *gorm.DB
}

// NewHelpUsHandler returns a new HelpUsHandler for the given DB
func NewHelpUsHandler(db *gorm.DB) *HelpUsHandler {
	return &HelpUsHandler{DB: db}
}

// GetHelpUs handles GET /api/helpus/:id
func (h *HelpUsHandler) GetHelpUs(c echo.Context) error {
	// lookup row by idHelpUs
	id := c.Param("id")

	// try to find the row and error if can't
	var hu model.HelpUs
	if err := h.DB.First(&hu, "idHelpUs = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "HelpUs not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, hu)
}

// CreateHelpUs handles POST /api/helpus
func (h *HelpUsHandler) CreateHelpUs(c echo.Context) error {
	// bind request JSON to HelpUs struct
	var hu model.HelpUs
	if err := c.Bind(&hu); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&hu).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create help us",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, hu)
}

// COPY HANDLER

// CopyHandler holds DB connection
type CopyHandler struct {
	DB *gorm.DB
}

// NewCopyHandler returns a new CopyHandler for the given DB
func NewCopyHandler(db *gorm.DB) *CopyHandler {
	return &CopyHandler{DB: db}
}

// GetCopy handles GET /api/copies/:id
func (h *CopyHandler) GetCopy(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var copy model.Copy
	if err := h.DB.First(&copy, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Copy not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, copy)
}

// CreateCopy handles POST /api/copies
func (h *CopyHandler) CreateCopy(c echo.Context) error {
	// bind request JSON to Copy struct
	var copy model.Copy
	if err := c.Bind(&copy); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&copy).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create copy",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, copy)
}

// EDITNEWROW HANDLER

// EditNewRowHandler holds DB connection
type EditNewRowHandler struct {
	DB *gorm.DB
}

// NewEditNewRowHandler returns a new EditNewRowHandler for the given DB
func NewEditNewRowHandler(db *gorm.DB) *EditNewRowHandler {
	return &EditNewRowHandler{DB: db}
}

// GetEditNewRow handles GET /api/editnewrows/:id
func (h *EditNewRowHandler) GetEditNewRow(c echo.Context) error {
	// lookup row by idEdit
	id := c.Param("id")

	// try to find the row and error if can't
	var enr model.EditNewRow
	if err := h.DB.First(&enr, "idEdit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "EditNewRow not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, enr)
}

// struct of what we expect from front end with info to make rows in Interaction and Edit and EditNewRow
type createEditNewRowPayload struct {
	IDInteractionType int64   `json:"IDInteractionType"`
	IDEntryType       int64   `json:"IDEntryType"`   
	IDSuggestion      int64   `json:"IDSuggestion"`     
	Mode              string  `json:"Mode"`       
	IsCorrect         int64   `json:"IsCorrect"`         
}

// CreateEditNewRow handles POST /api/editnewrows
func (h *EditNewRowHandler) CreateEditNewRow(c echo.Context) error {
	// read the cookie based session
	sess, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to read cookie session",
			"detail": err.Error(),
		})
	}

	// extract session_id from cookie from whatever form it's in
	var sessionID int64
	if raw, ok := sess.Values["session_id"]; ok {
		switch v := raw.(type) {
		case int64:
			sessionID = v
		case int:
			sessionID = int64(v)
		case float64:
			sessionID = int64(v)
		}
	}

	// bind request JSON filled with info for Interaction and Edit and EditNewRow
	var payload createEditNewRowPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// create Interaction using IDSession from cookie
	interaction := model.Interaction{
		IDSession:         sessionID,
		IDInteractionType: payload.IDInteractionType,
		// Timestamp is default
	}
	if err := h.DB.Create(&interaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create interaction",
			"detail": err.Error(),
		})
	}

	// create Edit linked to this Interaction
	edit := model.Edit{
		IDInteraction: interaction.IDInteraction,
		IDEntryType:   payload.IDEntryType,
		Mode:          payload.Mode,
		IsCorrect:     payload.IsCorrect,
	}
	if err := h.DB.Create(&edit).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit",
			"detail": err.Error(),
		})
	}

	// create EditNewRow linked to this Edit
	enr := model.EditNewRow{
		IDEdit:       edit.IDEdit,
		IDSuggestion: payload.IDSuggestion,
		IsCorrect:    payload.IsCorrect,
	}
	if err := h.DB.Create(&enr).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit new row",
			"detail": err.Error(),
		})
	}

	// return new row in edit new row
	return c.JSON(http.StatusCreated, enr)
}

// PASTE HANDLER

// PasteHandler holds DB connection
type PasteHandler struct {
	DB *gorm.DB
}

// NewPasteHandler returns a new PasteHandler for the given DB
func NewPasteHandler(db *gorm.DB) *PasteHandler {
	return &PasteHandler{DB: db}
}

// GetPaste handles GET /api/pastes/:id
func (h *PasteHandler) GetPaste(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var paste model.Paste
	if err := h.DB.First(&paste, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Paste not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, paste)
}

// CreatePaste handles POST /api/pastes
func (h *PasteHandler) CreatePaste(c echo.Context) error {
	// bind request JSON to Paste struct
	var paste model.Paste
	if err := c.Bind(&paste); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&paste).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create paste",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, paste)
}

// SEARCHGOOGLE HANDLER

// SearchGoogleHandler holds DB connection
type SearchGoogleHandler struct {
	DB *gorm.DB
}

// NewSearchGoogleHandler returns a new SearchGoogleHandler for the given DB
func NewSearchGoogleHandler(db *gorm.DB) *SearchGoogleHandler {
	return &SearchGoogleHandler{DB: db}
}

// GetSearchGoogle handles GET /api/searchgoogles/:id
func (h *SearchGoogleHandler) GetSearchGoogle(c echo.Context) error {
	// lookup row by IdInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var sg model.SearchGoogle
	if err := h.DB.First(&sg, "IdInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SearchGoogle not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, sg)
}

// CreateSearchGoogle handles POST /api/searchgoogles
func (h *SearchGoogleHandler) CreateSearchGoogle(c echo.Context) error {
	// bind request JSON to SearchGoogle struct
	var sg model.SearchGoogle
	if err := c.Bind(&sg); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&sg).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create search google",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, sg)
}

// VIEWCHANGE HANDLER

// ViewChangeHandler holds DB connection
type ViewChangeHandler struct {
	DB *gorm.DB
}

// NewViewChangeHandler returns a new ViewChangeHandler for the given DB
func NewViewChangeHandler(db *gorm.DB) *ViewChangeHandler {
	return &ViewChangeHandler{DB: db}
}

// GetViewChange handles GET /api/viewchanges/:id
func (h *ViewChangeHandler) GetViewChange(c echo.Context) error {
	// lookup row by idInteraction
	id := c.Param("id")

	// try to find the row and error if can't
	var vc model.ViewChange
	if err := h.DB.First(&vc, "idInteraction = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "ViewChange not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, vc)
}

// CreateViewChange handles POST /api/viewchanges
func (h *ViewChangeHandler) CreateViewChange(c echo.Context) error {
	// bind request JSON to ViewChange struct
	var vc model.ViewChange
	if err := c.Bind(&vc); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&vc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create view change",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, vc)
}

// VISIT HANDLER

// VisitHandler holds DB connection
type VisitHandler struct {
	DB *gorm.DB
}

// NewVisitHandler returns a new VisitHandler for the given DB
func NewVisitHandler(db *gorm.DB) *VisitHandler {
	return &VisitHandler{DB: db}
}

// GetVisit handles GET /api/visits/:id
func (h *VisitHandler) GetVisit(c echo.Context) error {
	// lookup row by idVisit
	id := c.Param("id")

	// try to find the row and error if can't
	var visit model.Visit
	if err := h.DB.First(&visit, "idVisit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Visit not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, visit)
}

// CreateVisit handles POST /api/visits
func (h *VisitHandler) CreateVisit(c echo.Context) error {
	// bind request JSON to Visit struct
	var visit model.Visit
	if err := c.Bind(&visit); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// insert into DB
	if err := h.DB.Create(&visit).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create visit",
			"detail": err.Error(),
		})
	}

	// return created row
	return c.JSON(http.StatusCreated, visit)
}

// SESSIONS HANDLER

// SessionsHandler holds DB connection
type SessionsHandler struct {
	DB *gorm.DB
}

// NewSessionsHandler returns a new SessionsHandler for the given DB
func NewSessionsHandler(db *gorm.DB) *SessionsHandler {
	return &SessionsHandler{DB: db}
}

// GetSessions handles GET /api/sessionsstore/:id
func (h *SessionsHandler) GetSessions(c echo.Context) error {
	// lookup row by session_id
	id := c.Param("id")

	// try to find the row and error if can't
	var s model.Sessions
	if err := h.DB.First(&s, "session_id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Sessions not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, s)
}

// CreateSessions handles POST /api/sessionsstore
func (h *SessionsHandler) CreateSessions(c echo.Context) error {
	// get cookie session via middleware
	sess, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to read cookie session",
		})
	}

	// set cookie expiration to 20 minutes
	const expiration = 20 * time.Minute

	// set session options
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(expiration.Seconds()),
		HttpOnly: true,
	}

	// get current time
	now := time.Now().Unix()

	// try to reuse existing DB session
	if val, ok := sess.Values["session_id"].(int64); ok && val != 0 {
		var existing model.Sessions
		if err := h.DB.First(&existing, "session_id = ?", val).Error; err == nil {
			// check if expires time is in the future compared to now
			if existing.Expires > now {
				// if so return the existing one
				return c.JSON(http.StatusOK, existing)
			}
		}
	}

	// create a new DB session
	newSession := model.Sessions{
		Expires: time.Now().Add(expiration).Unix(),
		Data:    nil,
	}

	// try to add it to the DB
	if err := h.DB.Create(&newSession).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create DB session",
			"detail": err.Error(),
		})
	}

	// store new ID in cookie
	sess.Values["session_id"] = newSession.SessionID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to save cookie session",
		})
	}

	// return the new session
	return c.JSON(http.StatusOK, newSession)
}

