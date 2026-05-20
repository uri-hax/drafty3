package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"drafty3/go_migration/data_model"
	"drafty3/go_migration/user_model"

	"github.com/gorilla/sessions"
	esession "github.com/labstack/echo-contrib/session"
)

// HealthCheckhandler returns a bollean
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}

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
	var suggestion data_model.Suggestions
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
	var suggestion data_model.Suggestions
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
	var alias data_model.Alias
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
	var alias data_model.Alias
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
	var click data_model.Click
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
	IDSuggestionType  int64   `json:"IDSuggestionType"`
	IDUniqueID        int64   `json:"IDUniqueID"`
	RowValues         *string `json:"RowValues"`
}

// CreateClick handles POST /api/clicks
func (h *ClickHandler) CreateClick(c echo.Context) error {
	log.Println("LOGGING:: CreateClick HIT")

	// read the cookie based session
	sessionID, err := getCookieSessionID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active session",
			"detail": err.Error(),
		})
	}

	// bind request JSON filled with info for Interaction and Click
	var payload createClickPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	log.Printf("payload: %+v", payload)

	// set up data models
	var interaction data_model.Interaction
	var click data_model.Click

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// find matching suggestions
		var suggestions []data_model.Suggestions
		if err := tx.
			Where("idSuggestionType = ? AND idUniqueID = ?", payload.IDSuggestionType, payload.IDUniqueID).
			Find(&suggestions).Error; err != nil {
			return err
		}

		// find the active suggestion
		var activeSuggestionID int64
		for _, s := range suggestions {
			if s.Active != nil && *s.Active == 1 {
				activeSuggestionID = s.IDSuggestion
				break
			}
		}

		// make sure we got an active suggestion
		if activeSuggestionID == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "no active suggestion found")
		}

		// create Interaction using IDSession from cookie
		interaction = data_model.Interaction{
			IDSession:         sessionID,
			IDInteractionType: payload.IDInteractionType,
		}
		if err := tx.Create(&interaction).Error; err != nil {
			return err
		}

		// create Click linked to this Interaction
		click = data_model.Click{
			IDInteraction: interaction.IDInteraction,
			IDSuggestion:  activeSuggestionID,
			RowValues:     payload.RowValues,
		}
		if err := tx.Create(&click).Error; err != nil {
			return err
		}

		return nil
	})

	// error handling for the transaction
	if err != nil {
		if httpErr, ok := err.(*echo.HTTPError); ok {
			return c.JSON(httpErr.Code, echo.Map{
				"error": httpErr.Message,
			})
		}

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
	var dt data_model.DataType
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
	var dt data_model.DataType
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
	var dct data_model.DatabaitCreateType
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
	var dct data_model.DatabaitCreateType
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
	var dna data_model.DatabaitNextAction
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
	var dna data_model.DatabaitNextAction
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
	var dtt data_model.DatabaitTemplateType
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
	var dtt data_model.DatabaitTemplateType
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
	var dc data_model.DoubleClick
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
	var dc data_model.DoubleClick
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
	var es data_model.EditSuggestion
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
	var es data_model.EditSuggestion
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
	var et data_model.EntryType
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
	var et data_model.EntryType
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
	var interaction data_model.Interaction
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
	var interaction data_model.Interaction
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
	var dt data_model.DatabaitTweet
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
	var dt data_model.DatabaitTweet
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
	var edit data_model.Edit
	if err := h.DB.First(&edit, "idEdit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Edit not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, edit)
}

// struct of what we expect from front end with info to make rows in Interaction, Edit, Suggestions, and EditSuggestion
type createEditPayload struct {
	IDInteractionType int64  `json:"IDInteractionType"`
	IDEntryType       int64  `json:"IDEntryType"`
	Mode              string `json:"Mode"`
	IsCorrect         int64  `json:"IsCorrect"`

	IDSuggestionType int64  `json:"IDSuggestionType"`
	IDUniqueID       int64  `json:"IDUniqueID"`
	Suggestion       string `json:"Suggestion"`
	Active           int64  `json:"Active"`
}

func (h *EditHandler) CreateEdit(c echo.Context) error {
	// read the cookie based session
	sessionID, err := getCookieSessionID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active session",
			"detail": err.Error(),
		})
	}

	// read the cookie based profile
	profileID, err := getCookieProfileID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active profile",
			"detail": err.Error(),
		})
	}

	// bind request JSON filled with info for Interaction, Edit, Suggestions, and EditSuggestion
	var payload createEditPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// set up data models
	var interaction data_model.Interaction
	var edit data_model.Edit
	var suggestion data_model.Suggestions
	var editSuggestion data_model.EditSuggestion

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// find matching suggestions for a cell
		var matchingSuggestions []data_model.Suggestions
		if err := tx.
			Where("idSuggestionType = ? AND idUniqueID = ?", payload.IDSuggestionType, payload.IDUniqueID).
			Find(&matchingSuggestions).Error; err != nil {
			return err
		}

		var isPrevSuggest int64 = 0
		var isNew int64 = 1

		// see if the suggestion in the payload matches any of the existing suggestions for that cell and change fields accordingly
		for _, s := range matchingSuggestions {
			if s.Suggestion == payload.Suggestion {
				isPrevSuggest = 1
				isNew = 0
				break
			}
		}

		var highestSuggestion data_model.Suggestions
		var nextConfidence int64 = 1

		// find the highest confidence suggestion for that cell
		err := tx.
			Where("idSuggestionType = ? AND idUniqueID = ?", payload.IDSuggestionType, payload.IDUniqueID).
			Order("confidence DESC").
			First(&highestSuggestion).Error

		// make sure we got a suggestion and handle error if not
		if err == nil {
			// set next confidence to be 1 higher than the highest confidence so far for that cell
			if highestSuggestion.Confidence != nil {
				nextConfidence = *highestSuggestion.Confidence + 1
			}

			// if the new suggestion is active then set the currently highest confidence suggestion to be inactive
			zero := int64(0)
			if err := tx.Model(&data_model.Suggestions{}).
				Where("idSuggestion = ?", highestSuggestion.IDSuggestion).
				Update("active", zero).Error; err != nil {
				return err
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		active := payload.Active

		// make sure chosen aligns with active
		var isChosen int64 = 0
		if active == 1 {
			isChosen = 1
		}

		confidence := nextConfidence

		// create Interaction using IDSession from cookie
		interaction = data_model.Interaction{
			IDSession:         sessionID,
			IDInteractionType: payload.IDInteractionType,
		}
		if err := tx.Create(&interaction).Error; err != nil {
			return err
		}

		// create Edit linked to this Interaction
		edit = data_model.Edit{
			IDInteraction: interaction.IDInteraction,
			IDEntryType:   payload.IDEntryType,
			Mode:          payload.Mode,
			IsCorrect:     payload.IsCorrect,
		}
		if err := tx.Create(&edit).Error; err != nil {
			return err
		}

		// create Suggestion linked to this Edit and the profile from the cookie
		suggestion = data_model.Suggestions{
			IDSuggestionType: payload.IDSuggestionType,
			IDUniqueID:       payload.IDUniqueID,
			IDProfile:        profileID,
			Suggestion:       payload.Suggestion,
			Active:           &active,
			Confidence:       &confidence,
		}
		if err := tx.Create(&suggestion).Error; err != nil {
			return err
		}

		// create EditSuggestion linking the Edit and Suggestion
		editSuggestion = data_model.EditSuggestion{
			IDEdit:        edit.IDEdit,
			IDSuggestion:  suggestion.IDSuggestion,
			IsPrevSuggest: isPrevSuggest,
			IsNew:         isNew,
			IsChosen:      isChosen,
		}
		if err := tx.Create(&editSuggestion).Error; err != nil {
			return err
		}

		return nil
	})

	// error handling for the transaction
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit flow",
			"detail": err.Error(),
		})
	}

	// return new rows in Edit, Suggestions, and EditSuggestion
	return c.JSON(http.StatusCreated, echo.Map{
		"edit":            edit,
		"suggestion":      suggestion,
		"edit_suggestion": editSuggestion,
	})
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
	var it data_model.InteractionType
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
	var it data_model.InteractionType
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

// GetProfile handles GET /api/users/profiles/:id
func (h *ProfileHandler) GetProfile(c echo.Context) error {
	// lookup row by idProfile
	id := c.Param("id")

	// try to find the row and error if can't
	var profile user_model.Profile
	if err := h.DB.First(&profile, "idProfile = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Profile not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, profile)
}

// CreateProfile handles POST /api/users/profiles
func (h *ProfileHandler) CreateProfile(c echo.Context) error {
	// bind request JSON to Profile struct
	var profile user_model.Profile
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
	var rud data_model.RemoveUserData
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
	var rud data_model.RemoveUserData
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
	var role data_model.Role
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
	var role data_model.Role
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
	var st data_model.SearchType
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
	var st data_model.SearchType
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
	var sr data_model.SelectRange
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
	var sr data_model.SelectRange
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

// GetSuggestionType handles GET /api/suggestiontypes/:name
func (h *SuggestionTypeHandler) GetSuggestionType(c echo.Context) error {
	// lookup row by name
	name := c.Param("name")

	// make sure name is provided
	if name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "name parameter is required",
		})
	}

	// try to find the row and error if can't
	var st data_model.SuggestionType
	if err := h.DB.
		Select("idSuggestionType").
		Where("LOWER(name) = LOWER(?)", name).
		First(&st).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "SuggestionType not found",
			"name":  name,
		})
	}

	// return the id of the row
	return c.JSON(http.StatusOK, echo.Map{
		"idSuggestionType": st.IDSuggestionType,
	})
}

// CreateSuggestionType handles POST /api/suggestiontypes
func (h *SuggestionTypeHandler) CreateSuggestionType(c echo.Context) error {
	// bind request JSON to SuggestionType struct
	var st data_model.SuggestionType
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
	var cc data_model.CopyColumn
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
	var cc data_model.CopyColumn
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
	var search data_model.Search
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
	sessionID, err := getCookieSessionID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active session",
			"detail": err.Error(),
		})
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
	interaction := data_model.Interaction{
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
	search := data_model.Search{
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
		searchMulti := data_model.SearchMulti{
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
	var sm data_model.SearchMulti
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
	var sm data_model.SearchMulti
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
	var sort data_model.Sort
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
	var sort data_model.Sort
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
	// lookup rows by idSuggestionType
	id := c.Param("id")

	// make sure id is provided
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "id parameter is required",
		})
	}

	// find all active rows for this suggestion type
	var stvs []data_model.SuggestionTypeValues
	if err := h.DB.
		Where("idSuggestionType = ? AND active = 1", id).
		Order("value ASC").
		Find(&stvs).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to fetch suggestion type values",
			"detail": err.Error(),
		})
	}

	// return matched rows
	return c.JSON(http.StatusOK, stvs)
}

// CreateSuggestionTypeValues handles POST /api/suggestiontypevalues
func (h *SuggestionTypeValuesHandler) CreateSuggestionTypeValues(c echo.Context) error {
	// bind request JSON to SuggestionTypeValues struct
	var stv data_model.SuggestionTypeValues
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
	var uid data_model.UniqueId
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
	var uid data_model.UniqueId
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
	var comment data_model.Comments
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
	var comment data_model.Comments
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
	var cv data_model.CommentVote
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
	var cv data_model.CommentVote
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
	var cv data_model.CommentsView
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
	var cv data_model.CommentsView
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
	var d data_model.Databaits
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
	var d data_model.Databaits
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
	var dv data_model.DatabaitVisit
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
	var dv data_model.DatabaitVisit
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
	var edr data_model.EditDelRow
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
	IDInteractionType int64  `json:"IDInteractionType"`
	IDEntryType       int64  `json:"IDEntryType"`
	Mode              string `json:"Mode"`
	IsCorrect         int64  `json:"IsCorrect"`
	IDUniqueID        int64  `json:"IDUniqueID"`
	Comment           string `json:"Comment"`
}

// CreateEditDelRow handles POST /api/editdelrows
func (h *EditDelRowHandler) CreateEditDelRow(c echo.Context) error {
	// read the cookie based session
	sessionID, err := getCookieSessionID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active session",
			"detail": err.Error(),
		})
	}

	// bind request JSON filled with info for Interaction and Edit and EditDelRow
	var payload createEditDelRowPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// set up data models
	var interaction data_model.Interaction
	var edit data_model.Edit
	var edr data_model.EditDelRow

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// create Interaction using IDSession from cookie
		interaction = data_model.Interaction{
			IDSession:         sessionID,
			IDInteractionType: payload.IDInteractionType,
		}
		if err := tx.Create(&interaction).Error; err != nil {
			return err
		}

		// create Edit linked to this Interaction
		edit = data_model.Edit{
			IDInteraction: interaction.IDInteraction,
			IDEntryType:   payload.IDEntryType,
			Mode:          payload.Mode,
			IsCorrect:     payload.IsCorrect,
		}
		if err := tx.Create(&edit).Error; err != nil {
			return err
		}

		// find suggestion rows for this unique id
		var suggestions []data_model.Suggestions
		if err := tx.
			Where("idUniqueID = ?", payload.IDUniqueID).
			Find(&suggestions).Error; err != nil {
			return err
		}

		// find the active one and set it to 0
		for _, s := range suggestions {
			if s.Active != nil && *s.Active == 1 {
				zero := int64(0)
				if err := tx.Model(&data_model.Suggestions{}).
					Where("idSuggestion = ?", s.IDSuggestion).
					Update("active", &zero).Error; err != nil {
					return err
				}
				break
			}
		}

		// deactivate the row itself in UniqueId
		zero := int64(0)
		if err := tx.Model(&data_model.UniqueId{}).
			Where("idUniqueID = ?", payload.IDUniqueID).
			Update("active", zero).Error; err != nil {
			return err
		}

		// create EditDelRow linked to Edit and UniqueId
		edr = data_model.EditDelRow{
			IDEdit:     edit.IDEdit,
			IDUniqueID: payload.IDUniqueID,
			Comment:    payload.Comment,
		}
		if err := tx.Create(&edr).Error; err != nil {
			return err
		}

		return nil
	})

	// error handling for the transaction
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit del row flow",
			"detail": err.Error(),
		})
	}

	// return created rows
	return c.JSON(http.StatusCreated, echo.Map{
		"interaction": interaction,
		"edit":        edit,
		"editdelrow":  edr,
	})
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
	var hu data_model.HelpUs
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
	var hu data_model.HelpUs
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
	var copy data_model.Copy
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
	var copy data_model.Copy
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
	var enr data_model.EditNewRow
	if err := h.DB.First(&enr, "idEdit = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "EditNewRow not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, enr)
}

// struct of what we expect from front end with info to make rows in Interaction, Edit, many Suggestions, and EditNewRow
type createEditNewRowCellPayload struct {
	IDSuggestionType int64  `json:"IDSuggestionType"`
	Suggestion       string `json:"Suggestion"`
	Active           int64  `json:"Active"`
	Confidence       int64  `json:"Confidence"`
}

type createEditNewRowPayload struct {
	IDInteractionType int64                         `json:"IDInteractionType"`
	IDEntryType       int64                         `json:"IDEntryType"`
	Mode              string                        `json:"Mode"`
	IsCorrect         int64                         `json:"IsCorrect"`
	Cells             []createEditNewRowCellPayload `json:"Cells"`
}

// CreateEditNewRow handles POST /api/editnewrows
func (h *EditNewRowHandler) CreateEditNewRow(c echo.Context) error {
	// read the cookie based session
	sessionID, err := getCookieSessionID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active session",
			"detail": err.Error(),
		})
	}

	// read the cookie based profile
	profileID, err := getCookieProfileID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error":  "failed to get active profile",
			"detail": err.Error(),
		})
	}

	// bind request JSON filled with info for Interaction, Edit, Suggestions, and EditNewRow
	var payload createEditNewRowPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid request body",
			"detail": err.Error(),
		})
	}

	// validate that at least one cell is provided
	if len(payload.Cells) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "at least one cell is required",
		})
	}

	// set up data models
	var interaction data_model.Interaction
	var edit data_model.Edit
	var enr data_model.EditNewRow
	var uid data_model.UniqueId
	createdSuggestions := make([]data_model.Suggestions, 0, len(payload.Cells))

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		// create the UniqueId row for this new row
		notes := ""
		uid = data_model.UniqueId{
			Active: 1,
			Notes:  &notes,
		}
		if err := tx.Create(&uid).Error; err != nil {
			return err
		}

		// create Interaction using IDSession from cookie
		interaction = data_model.Interaction{
			IDSession:         sessionID,
			IDInteractionType: payload.IDInteractionType,
		}
		if err := tx.Create(&interaction).Error; err != nil {
			return err
		}

		// create Edit linked to this Interaction
		edit = data_model.Edit{
			IDInteraction: interaction.IDInteraction,
			IDEntryType:   payload.IDEntryType,
			Mode:          payload.Mode,
			IsCorrect:     payload.IsCorrect,
		}
		if err := tx.Create(&edit).Error; err != nil {
			return err
		}

		// create one Suggestion per cell
		for _, cell := range payload.Cells {
			active := cell.Active
			confidence := cell.Confidence

			suggestion := data_model.Suggestions{
				IDSuggestionType: cell.IDSuggestionType,
				IDUniqueID:       uid.IDUniqueID,
				IDProfile:        profileID,
				Suggestion:       cell.Suggestion,
				Active:           &active,
				Confidence:       &confidence,
			}
			if err := tx.Create(&suggestion).Error; err != nil {
				return err
			}

			createdSuggestions = append(createdSuggestions, suggestion)
		}

		// create EditNewRow linked to this Edit and the first new Suggestion
		enr = data_model.EditNewRow{
			IDEdit:       edit.IDEdit,
			IDSuggestion: createdSuggestions[0].IDSuggestion,
			IsCorrect:    payload.IsCorrect,
		}
		if err := tx.Create(&enr).Error; err != nil {
			return err
		}

		return nil
	})

	// error handling for the transaction
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create edit new row flow",
			"detail": err.Error(),
		})
	}

	// return new rows in UniqueId, Edit, Suggestions, and EditNewRow
	return c.JSON(http.StatusCreated, echo.Map{
		"uniqueId":    uid,
		"idUniqueID":  uid.IDUniqueID,
		"edit":        edit,
		"suggestions": createdSuggestions,
		"editNewRow":  enr,
	})
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
	var paste data_model.Paste
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
	var paste data_model.Paste
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
	var sg data_model.SearchGoogle
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
	var sg data_model.SearchGoogle
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
	var vc data_model.ViewChange
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
	var vc data_model.ViewChange
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
	var visit data_model.Visit
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
	var visit data_model.Visit
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

// GetSessions handles GET /api/users/sessions/:id
func (h *SessionsHandler) GetSessions(c echo.Context) error {
	// lookup row by idSession
	id := c.Param("id")

	// try to find the row and error if can't
	var s user_model.Session
	if err := h.DB.First(&s, "idSession = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Session not found",
			"id":    id,
		})
	}

	// return the row
	return c.JSON(http.StatusOK, s)
}

// CreateSessions handles POST /api/users/sessions
func (h *SessionsHandler) CreateSessions(c echo.Context) error {
	// get expiration time and current time for session management
	const expiration = 20 * time.Minute
	now := time.Now()

	// read the cookie based session from middleware
	cookieSession, err := esession.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to read cookie session",
		})
	}

	// set cookie options - OLD
	//cookieSession.Options = &sessions.Options{
	//	Path:     "/",
	//	MaxAge:   int(expiration.Seconds()),
	//	HttpOnly: true,
	//}

	// set cookie options - NEW
	cookieSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(expiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	// first try to reuse existing session from cookie
	var profile user_model.Profile
	// use id_session from cookie
	if sessionID, ok := getInt64(cookieSession.Values["id_session"]); ok && sessionID != 0 {
		var existing user_model.Session
		// see if session exists
		if err := h.DB.First(&existing, "idSession = ?", sessionID).Error; err == nil {
			// see if it's not expired yet
			if now.Before(existing.End) {
				// get profile for this session and return both
				if err := h.DB.First(&profile, "idProfile = ?", existing.IDProfile).Error; err == nil {
					return c.JSON(http.StatusOK, echo.Map{
						"profile": profile,
						"session": existing,
					})
				}
			}
		}
	}

	// second try to reuse existing profile from cookie, or create one
	// use profile_id from cookie
	if profileID, ok := getInt64(cookieSession.Values["profile_id"]); ok && profileID != 0 {
		// see if profile exists
		if err := h.DB.First(&profile, "idProfile = ?", profileID).Error; err != nil {
			profile = user_model.Profile{}
			// if profile doesn't exist then make a new one (old cookie)
			if err := h.DB.Create(&profile).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error":  "failed to create profile",
					"detail": err.Error(),
				})
			}
		}
		// if no profile in cookie then make a new one (first time user)
	} else {
		profile = user_model.Profile{}
		if err := h.DB.Create(&profile).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error":  "failed to create profile",
				"detail": err.Error(),
			})
		}
	}

	// third try to create new session with calculated end time
	// create model instance with profile id and start and end time
	newSession := user_model.Session{
		IDProfile: profile.IDProfile,
		Start:     now,
		End:       now.Add(expiration),
	}

	// create the session in db and error if fail
	if err := h.DB.Create(&newSession).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":  "failed to create session",
			"detail": err.Error(),
		})
	}

	// fourth try to save identifiers in cookie
	cookieSession.Values["profile_id"] = profile.IDProfile
	cookieSession.Values["id_session"] = newSession.IDSession

	// save the cookie session and error if fail
	if err := cookieSession.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to save cookie session",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"profile": profile,
		"session": newSession,
	})
}

// function to extract numerical IDs and convert to int64
func getInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int64:
		return val, true
	case int:
		return int64(val), true
	case int32:
		return int64(val), true
	case float64:
		return int64(val), true
	default:
		return 0, false
	}
}

func getCookieSessionID(c echo.Context) (int64, error) {
	sess, err := esession.Get("session", c)
	if err != nil {
		return 0, err
	}

	sessionID, ok := getInt64(sess.Values["id_session"])
	if !ok || sessionID == 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "no active session cookie")
	}

	return sessionID, nil
}

func getCookieProfileID(c echo.Context) (int64, error) {
	sess, err := esession.Get("session", c)
	if err != nil {
		return 0, err
	}

	profileID, ok := getInt64(sess.Values["profile_id"])
	if !ok || profileID == 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "no active profile cookie")
	}

	return profileID, nil
}
