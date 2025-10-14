package model

import "time"

// Alias
type Alias struct {
	IdAlias      int    `gorm:"column:idAlias;primaryKey;autoIncrement"`
	IdSuggestion int    `gorm:"column:idSuggestion;not null;index;index:unique_index,unique,priority:1"`
	Alias        string `gorm:"column:alias;size:500;index:unique_index,unique,priority:2"`
	Count        int    `gorm:"column:count;not null;default:1"`
}
func (Alias) TableName() string { return "Alias" }

// Click
type Click struct {
	IdInteraction int    `gorm:"column:idInteraction;primaryKey"`
	IdSuggestion  int    `gorm:"column:idSuggestion;not null;index"`
	Rowvalues     string `gorm:"column:rowvalues"`
}
func (Click) TableName() string { return "Click" }

// DataType
type DataType struct {
	IdDataType int    `gorm:"column:idDataType;primaryKey;autoIncrement"`
	Type       string `gorm:"column:type"`
}
func (DataType) TableName() string { return "DataType" }

// DatabaitCreateType
type DatabaitCreateType struct {
	IdDatabaitCreateType int    `gorm:"column:idDatabaitCreateType;primaryKey;autoIncrement"`
	Type                 string `gorm:"column:type"`
}
func (DatabaitCreateType) TableName() string { return "DatabaitCreateType" }

// DatabaitNextAction
type DatabaitNextAction struct {
	IdDatabaitNextAction int    `gorm:"column:idDatabaitNextAction;primaryKey;autoIncrement"`
	Action               string `gorm:"column:action"`
}
func (DatabaitNextAction) TableName() string { return "DatabaitNextAction" }

// DatabaitTemplateType
type DatabaitTemplateType struct {
	IdDatabaitTemplateType int    `gorm:"column:idDatabaitTemplateType;primaryKey;autoIncrement"`
	Template               string `gorm:"column:template"`
}
func (DatabaitTemplateType) TableName() string { return "DatabaitTemplateType" }

// DoubleClick
type DoubleClick struct {
	IdInteraction int    `gorm:"column:idInteraction;primaryKey"`
	IdSuggestion  int    `gorm:"column:idSuggestion;not null;index"`
	Rowvalues     string `gorm:"column:rowvalues"`
}
func (DoubleClick) TableName() string { return "DoubleClick" }

// Edit_Suggestion
type EditSuggestion struct {
	IdEdit           int  `gorm:"column:idEdit;not null;index;index:ux_edit_suggestion,unique,priority:1"`
	IdSuggestion     int  `gorm:"column:idSuggestion;not null;index;index:ux_edit_suggestion,unique,priority:2"`
	IsPrevSuggestion bool `gorm:"column:isPrevSuggestion;not null"`
	IsNew            bool `gorm:"column:isNew;not null"`
	IsChosen         bool `gorm:"column:isChosen;not null"`
}
func (EditSuggestion) TableName() string { return "Edit_Suggestion" }

// EntryType
type EntryType struct {
	IdEntryType int    `gorm:"column:idEntryType;primaryKey;autoIncrement"`
	Type        string `gorm:"column:type"`
}
func (EntryType) TableName() string { return "EntryType" }

// Interaction
type Interaction struct {
	IdInteraction     int       `gorm:"column:idInteraction;primaryKey;autoIncrement"`
	IdSession         int       `gorm:"column:idSession;not null;index"`
	IdInteractionType int       `gorm:"column:idInteractionType;not null;index"`
	Timestamp         time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`
}
func (Interaction) TableName() string { return "Interaction" }

// DatabaitTweet
type DatabaitTweet struct {
	IdDatabaitTweet     int       `gorm:"column:idDatabaitTweet;primaryKey;autoIncrement"`
	IdInteraction       int       `gorm:"column:idInteraction;not null;index"`
	IdDatabait          int       `gorm:"column:idDatabait;not null"`
	Url                 string    `gorm:"column:url;not null"`
	Likes               int       `gorm:"column:likes"`
	Retweets            int       `gorm:"column:retweets"`
	Created             time.Time `gorm:"column:created;not null;default:CURRENT_TIMESTAMP"`
	NextActionTimestamp time.Time `gorm:"column:nextActionTimestamp"`
	NextAction          int       `gorm:"column:nextAction"`
}
func (DatabaitTweet) TableName() string { return "DatabaitTweet" }

// Edit
type Edit struct {
	IdInteraction int    `gorm:"column:IdInteraction;not null;index"`
	IdEdit        int    `gorm:"column:idEdit;primaryKey;autoIncrement"`
	IdEntryType   int    `gorm:"column:idEntryType;not null"`
	Mode          string `gorm:"column:mode;not null;default:normal"`
	IsCorrect     int    `gorm:"column:isCorrect;default:2"`
}
func (Edit) TableName() string { return "Edit" }

// InteractionType
type InteractionType struct {
	IdInteractionType int    `gorm:"column:idInteractionType;primaryKey;autoIncrement"`
	Interaction       string `gorm:"column:interaction"`
}
func (InteractionType) TableName() string { return "InteractionType" }

// Profile
type Profile struct {
	IdProfile   int       `gorm:"column:idProfile;primaryKey;autoIncrement"`
	IdRole      int       `gorm:"column:idRole;not null;default:2;index"`
	Username    string    `gorm:"column:username;unique"`
	Email       string    `gorm:"column:email;unique"`
	Password    string    `gorm:"column:password"`
	PasswordRaw string    `gorm:"column:passwordRaw"`
	DateCreated time.Time `gorm:"column:date_created;default:CURRENT_TIMESTAMP"`
	DateUpdated time.Time `gorm:"column:date_updated;default:CURRENT_TIMESTAMP"`
}
func (Profile) TableName() string { return "Profile" }

// RemoveUserData
type RemoveUserData struct {
	IDRemoveuserdata int       `gorm:"column:id_removeuserdata;primaryKey;autoIncrement"`
	IDProfile        int       `gorm:"column:id_profile;not null"`
	IDSession        int       `gorm:"column:id_session;not null"`
	Timestamp        time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`
}
func (RemoveUserData) TableName() string { return "RemoveUserData" }

// Role
type Role struct {
	IdRole int    `gorm:"column:idRole;primaryKey;autoIncrement"`
	Role   string `gorm:"column:role;not null"`
}
func (Role) TableName() string { return "Role" }

// SearchType
type SearchType struct {
	IdSearchType int    `gorm:"column:idSearchType;primaryKey;autoIncrement"`
	Type         string `gorm:"column:type;not null"`
}
func (SearchType) TableName() string { return "SearchType" }

// SelectRange
type SelectRange struct {
	IdInteraction int    `gorm:"column:idInteraction;not null;index"`
	IdSuggestion  int    `gorm:"column:idSuggestion;not null;index"`
	Rowvalues     string `gorm:"column:rowvalues"`
}
func (SelectRange) TableName() string { return "SelectRange" }

// Session
type Session struct {
	IdSession int       `gorm:"column:idSession;primaryKey;autoIncrement"`
	IdProfile int       `gorm:"column:idProfile;index"`
	Start     time.Time `gorm:"column:start;default:CURRENT_TIMESTAMP"`
	End       time.Time `gorm:"column:end"`
}
func (Session) TableName() string { return "Session" }

// SuggestionType
type SuggestionType struct {
	IdSuggestionType int    `gorm:"column:idSuggestionType;primaryKey;autoIncrement"`
	IdDataType       int    `gorm:"column:idDataType;not null;index"`
	Name             string `gorm:"column:name"`
	IsActive         bool   `gorm:"column:isActive;not null;default:1"`
	Regex            string `gorm:"column:regex;not null;default:.*"`
	MakesRowUnique   bool   `gorm:"column:makesRowUnique;default:0"`
	CanBeBlank       bool   `gorm:"column:canBeBlank;not null;default:0"`
	IsFreeEdit       bool   `gorm:"column:isFreeEdit;not null;default:1"`
	IsDate           bool   `gorm:"column:isDate;not null;default:0"`
	IsLink           bool   `gorm:"column:isLink;not null;default:0"`
	IsCurrency       bool   `gorm:"column:isCurrency;not null;default:0"`
	IsEditable       int    `gorm:"column:isEditable;not null;default:1"`
	IsPrivate        int    `gorm:"column:isPrivate;not null;default:0"`
	ColumnOrder      int    `gorm:"column:columnOrder"`
}
func (SuggestionType) TableName() string { return "SuggestionType" }

// CopyColumn
type CopyColumn struct {
	IdInteraction    int `gorm:"column:idInteraction;primaryKey"`
	IdSuggestionType int `gorm:"column:idSuggestionType;primaryKey"`
}
func (CopyColumn) TableName() string { return "CopyColumn" }

// Search
type Search struct {
	IdInteraction    int    `gorm:"column:idInteraction;not null;index"`
	IdSuggestionType int    `gorm:"column:idSuggestionType;not null;index"`
	IdSearchType     int    `gorm:"column:idSearchType;not null;default:3;index"`
	IsPartial        bool   `gorm:"column:isPartial;not null;default:1"`
	IsMulti          int    `gorm:"column:isMulti;not null;default:0"`
	IsFromUrl        int    `gorm:"column:isFromUrl;not null;default:0"`
	Value            string `gorm:"column:value"`
	MatchedValues    []byte `gorm:"column:matchedValues"`
}
func (Search) TableName() string { return "Search" }

// SearchMulti
type SearchMulti struct {
	IdInteraction    int    `gorm:"column:idInteraction;primaryKey;index"`
	IdSuggestionType int    `gorm:"column:idSuggestionType;primaryKey;index"`
	IdSearchType     int    `gorm:"column:idSearchType;not null;default:3;index"`
	Value            string `gorm:"column:value"`
}
func (SearchMulti) TableName() string { return "SearchMulti" }

// Sort
type Sort struct {
	IdInteraction    int  `gorm:"column:idInteraction;primaryKey;index"`
	IdSuggestionType int  `gorm:"column:idSuggestionType;primaryKey;index"`
	IsAsc            bool `gorm:"column:isAsc;not null;default:1"`
	IsTrigger        bool `gorm:"column:isTrigger;not null;default:1"`
	IsMulti          bool `gorm:"column:isMulti;not null;default:0"`
}
func (Sort) TableName() string { return "Sort" }

// SuggestionTypeValues
type SuggestionTypeValues struct {
	IdSuggestionType int    `gorm:"column:idSuggestionType;not null;index;index:PRIMARY_id_and_value,unique,priority:1"`
	Value            string `gorm:"column:value;index:PRIMARY_id_and_value,unique,priority:2"`
	Active           bool   `gorm:"column:active;not null;default:1"`
}
func (SuggestionTypeValues) TableName() string { return "SuggestionTypeValues" }

// UniqueId
type UniqueId struct {
	IdUniqueID int    `gorm:"column:idUniqueID;primaryKey;autoIncrement"`
	Active     bool   `gorm:"column:active;not null;default:1"`
	Notes      string `gorm:"column:notes"`
}
func (UniqueId) TableName() string { return "UniqueId" }

// Comments
type Comments struct {
	IdComment     int    `gorm:"column:idComment;primaryKey;autoIncrement"`
	IdInteraction int    `gorm:"column:idInteraction;not null"`
	IdUniqueID    int    `gorm:"column:idUniqueID;not null"`
	Comment       string `gorm:"column:comment;not null"`
	VoteUp        int    `gorm:"column:voteUp;not null;default:0"`
	VoteDown      int    `gorm:"column:voteDown;not null;default:0"`
}
func (Comments) TableName() string { return "Comments" }

// CommentVote
type CommentVote struct {
	IdCommentVote int    `gorm:"column:idCommentVote;primaryKey;autoIncrement"`
	IdInteraction int    `gorm:"column:idInteraction;not null;index"`
	IdComment     int    `gorm:"column:idComment;not null"`
	Vote          string `gorm:"column:vote;not null;check:vote IN ('voteUp','voteUp-deselect','voteDown','voteDown-deselect')"`
	Selected      bool   `gorm:"column:selected"`
}
func (CommentVote) TableName() string { return "CommentVote" }

// CommentsView
type CommentsView struct {
	IdCommentsView int `gorm:"column:idCommentsView;primaryKey;autoIncrement"`
	IdUniqueID     int `gorm:"column:idUniqueID"`
	IdInteraction  int `gorm:"column:idInteraction"`
}
func (CommentsView) TableName() string { return "CommentsView" }

// Databaits
type Databaits struct {
	IdDatabait             int       `gorm:"column:idDatabait;primaryKey;autoIncrement"`
	IdInteraction          int       `gorm:"column:idInteraction;not null"`
	IdUniqueID             int       `gorm:"column:idUniqueID"`
	IdDatabaitTemplateType int       `gorm:"column:idDatabaitTemplateType;not null"`
	IdDatabaitCreateType   int       `gorm:"column:idDatabaitCreateType;not null"`
	Databait               string    `gorm:"column:databait;not null"`
	Columns                string    `gorm:"column:columns"`
	Vals                   string    `gorm:"column:vals"`
	Notes                  string    `gorm:"column:notes;not null"`
	Created                time.Time `gorm:"column:created;not null;default:CURRENT_TIMESTAMP"`
	Closed                 time.Time `gorm:"column:closed;not null"` 
	NextAction             int       `gorm:"column:nextAction"`
}
func (Databaits) TableName() string { return "Databaits" }

// DatabaitVisit
type DatabaitVisit struct {
	IdInteraction int    `gorm:"column:idInteraction;not null;unique"`
	IdDatabait    int    `gorm:"column:idDatabait;not null"`
	Source        string `gorm:"column:source"`
}
func (DatabaitVisit) TableName() string { return "DatabaitVisit" }

// Edit_DelRow
type EditDelRow struct {
	IdUniqueID int    `gorm:"column:idUniqueID;not null"`
	IdEdit     int    `gorm:"column:idEdit;not null"`
	Comment    string `gorm:"column:comment;not null"`
}
func (EditDelRow) TableName() string { return "Edit_DelRow" }

// HelpUs
type HelpUs struct {
	IdHelpUs      int       `gorm:"column:idHelpUs;primaryKey;autoIncrement"`
	IdInteraction int       `gorm:"column:idInteraction;not null"`
	IdUniqueID    int       `gorm:"column:idUniqueID;not null"`
	HelpUsType    string    `gorm:"column:helpUsType;not null"`
	Question      string    `gorm:"column:question;not null"`
	Answer        string    `gorm:"column:answer"`
	Start         time.Time `gorm:"column:start;not null;default:CURRENT_TIMESTAMP"`
	Answered      time.Time `gorm:"column:answered;not null"`    
	ShowAnother   time.Time `gorm:"column:showAnother;not null"` 
	Closed        time.Time `gorm:"column:closed;not null"`      
}
func (HelpUs) TableName() string { return "HelpUs" }

// Suggestions
type Suggestions struct {
	IdSuggestion     int       `gorm:"column:idSuggestion;primaryKey;autoIncrement"`
	IdSuggestionType int       `gorm:"column:idSuggestionType;not null;index"`
	IdUniqueID       int       `gorm:"column:idUniqueID;not null;index"`
	IdProfile        int       `gorm:"column:idProfile;not null;default:2"`
	Suggestion       string    `gorm:"column:suggestion;not null"`
	Active           bool      `gorm:"column:active;not null;default:1"`
	Confidence       int64     `gorm:"column:confidence"`
	LastUpdated      time.Time `gorm:"column:last_updated;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}
func (Suggestions) TableName() string { return "Suggestions" }

// Copy
type Copy struct {
	IdInteraction int `gorm:"column:idInteraction;primaryKey"`
	IdSuggestion  int `gorm:"column:idSuggestion;primaryKey"`
}
func (Copy) TableName() string { return "Copy" }

// Edit_NewRow
type EditNewRow struct {
	IdEdit       int `gorm:"column:idEdit;not null;index;index:ux_edit_newrow,unique,priority:1"`
	IdSuggestion int `gorm:"column:idSuggestion;not null;index;index:ux_edit_newrow,unique,priority:2"`
	IsCorrect    int `gorm:"column:isCorrect;not null;default:2"`
}
func (EditNewRow) TableName() string { return "Edit_NewRow" }

// Paste
type Paste struct {
	IdInteraction         int    `gorm:"column:idInteraction;primaryKey"`
	PasteValue            string `gorm:"column:pasteValue;not null"`
	CopyCellIdSuggestion  int    `gorm:"column:copyCellIdSuggestion"`
	CopyCellValue         string `gorm:"column:copyCellValue"`
	PasteCellIdSuggestion int    `gorm:"column:pasteCellIdSuggestion;not null"`
	PasteCellValue        string `gorm:"column:pasteCellValue;not null"`
}
func (Paste) TableName() string { return "Paste" }

// SearchGoogle
type SearchGoogle struct {
	IdInteraction int    `gorm:"column:IdInteraction;not null;unique"`
	IdUniqueID    int    `gorm:"column:idUniqueID;not null"`
	IdSuggestion  int    `gorm:"column:idSuggestion;not null"`
	SearchValues  string `gorm:"column:searchValues;not null"`
}
func (SearchGoogle) TableName() string { return "SearchGoogle" }

// ViewChange
type ViewChange struct {
	IdInteraction int    `gorm:"column:idInteraction;primaryKey"`
	Viewname      string `gorm:"column:viewname;primaryKey"`
}
func (ViewChange) TableName() string { return "ViewChange" }

// Visit
type Visit struct {
	IdVisit       int    `gorm:"column:idVisit;primaryKey;autoIncrement"`
	IdInteraction int    `gorm:"column:idInteraction;not null"`
	Source        string `gorm:"column:source"`
	SearchCol     string `gorm:"column:searchCol"`
	SearchVal     string `gorm:"column:searchVal"`
}
func (Visit) TableName() string { return "Visit" }

// sessions
type Sessions struct {
	SessionID string `gorm:"column:session_id;primaryKey"`
	Expires   uint   `gorm:"column:expires;not null"`
	Data      string `gorm:"column:data"`
}
func (Sessions) TableName() string { return "sessions" }
