package model

import "time"

type Alias struct {
	IDAlias      int64   `gorm:"column:idAlias;primaryKey;autoIncrement"`
	IDSuggestion int64   `gorm:"column:idSuggestion;not null;index:fk_Alias_Suggestion1_idx;uniqueIndex:unique_index"`
	Alias        *string `gorm:"column:alias;uniqueIndex:unique_index"`
	Count        int64   `gorm:"column:count;not null;default:1"`
}
func (Alias) TableName() string { return "Alias" }

type Click struct {
	IDInteraction int64   `gorm:"column:idInteraction;primaryKey;not null;autoIncrement:false;index:fk_Click_Interaction1_idx"`
	IDSuggestion  int64   `gorm:"column:idSuggestion;not null;index:fk_Click_Suggestion1_idx"`
	RowValues     *string `gorm:"column:rowvalues"`
}
func (Click) TableName() string { return "Click" }

type DataType struct {
	IDDataType int64   `gorm:"column:idDataType;primaryKey;autoIncrement"`
	Type       *string `gorm:"column:type"`
}
func (DataType) TableName() string { return "DataType" }

type DatabaitCreateType struct {
	IDDatabaitCreateType int64   `gorm:"column:idDatabaitCreateType;primaryKey;autoIncrement"`
	Type                 *string `gorm:"column:type"`
}
func (DatabaitCreateType) TableName() string { return "DatabaitCreateType" }

type DatabaitNextAction struct {
	IDDatabaitNextAction int64   `gorm:"column:idDatabaitNextAction;primaryKey;autoIncrement"`
	Action               *string `gorm:"column:action"`
}
func (DatabaitNextAction) TableName() string { return "DatabaitNextAction" }

type DatabaitTemplateType struct {
	IDDatabaitTemplateType int64   `gorm:"column:idDatabaitTemplateType;primaryKey;autoIncrement"`
	Template               *string `gorm:"column:template"`
}
func (DatabaitTemplateType) TableName() string { return "DatabaitTemplateType" }

type DoubleClick struct {
	IDInteraction int64   `gorm:"column:idInteraction;primaryKey;not null;autoIncrement:false;index:fk_DoubleClick_Interaction1_idx"`
	IDSuggestion  int64   `gorm:"column:idSuggestion;not null;index:fk_DoubleClick_Suggestion1_idx"`
	RowValues     *string `gorm:"column:rowvalues"`
}
func (DoubleClick) TableName() string { return "DoubleClick" }

type EditSuggestion struct {
	IDEdit        int64 `gorm:"column:idEdit;not null;index:_index_edit_suggestion_idEdit_agsdh1872dg;uniqueIndex:idEdit"`
	IDSuggestion  int64 `gorm:"column:idSuggestion;not null;index:_index_edit_suggestion_idSuggestion_agsdh1872dg;uniqueIndex:idEdit"`
	IsPrevSuggest int64 `gorm:"column:isPrevSuggestion;not null"`
	IsNew         int64 `gorm:"column:isNew;not null"`
	IsChosen      int64 `gorm:"column:isChosen;not null"`
}
func (EditSuggestion) TableName() string { return "Edit_Suggestion" }

type EntryType struct {
	IDEntryType int64   `gorm:"column:idEntryType;primaryKey;autoIncrement"`
	Type        *string `gorm:"column:type"`
}
func (EntryType) TableName() string { return "EntryType" }

type Interaction struct {
	IDInteraction     int64     `gorm:"column:idInteraction;primaryKey;autoIncrement"`
	IDSession         int64     `gorm:"column:idSession;not null;index:fk_Interaction_Session1_idx"`
	IDInteractionType int64     `gorm:"column:idInteractionType;not null;index:fk_Interaction_InteractionType1_idx"`
	Timestamp         time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`
}
func (Interaction) TableName() string { return "Interaction" }

type DatabaitTweet struct {
	IDDatabaitTweet     int64      `gorm:"column:idDatabaitTweet;primaryKey;autoIncrement"`
	IDInteraction       int64      `gorm:"column:idInteraction;not null"`
	IDDatabait          int64      `gorm:"column:idDatabait;not null"`
	URL                 string     `gorm:"column:url;not null"`
	Likes               *int64     `gorm:"column:likes"`
	Retweets            *int64     `gorm:"column:retweets"`
	Created             time.Time  `gorm:"column:created;not null;default:CURRENT_TIMESTAMP"`
	NextActionTimestamp *time.Time `gorm:"column:nextActionTimestamp"`
	NextAction          *int64     `gorm:"column:nextAction"`

	InteractionRef Interaction        `gorm:"constraint:Name:_fk_idInteraction_DatabaitTweet_t615das;foreignKey:IDInteraction;references:IDInteraction"`
	NextActionRef  DatabaitNextAction `gorm:"constraint:Name:_fk_nextAction_DatabaitTweet_t615das;foreignKey:NextAction;references:IDDatabaitNextAction"`
}
func (DatabaitTweet) TableName() string { return "DatabaitTweet" }

type Edit struct {
	IDInteraction int64  `gorm:"column:IdInteraction;not null;index:idInteraction_index_adfhj126"`
	IDEdit        int64  `gorm:"column:idEdit;primaryKey;autoIncrement"`
	IDEntryType   int64  `gorm:"column:idEntryType;not null"`
	Mode          string `gorm:"column:mode;not null;default:normal"`
	IsCorrect     *int64 `gorm:"column:isCorrect;default:2"`

	EntryTypeRef   EntryType   `gorm:"constraint:Name:_fk_idEntryType_from_edit_asllhg1233;foreignKey:IDEntryType;references:IDEntryType"`
	InteractionRef Interaction `gorm:"constraint:Name:_fk_idInteraction_from_edit_asdlhg1235;foreignKey:IDInteraction;references:IDInteraction"`
}
func (Edit) TableName() string { return "Edit" }

type InteractionType struct {
	IDInteractionType int64   `gorm:"column:idInteractionType;primaryKey;autoIncrement"`
	Interaction       *string `gorm:"column:interaction"`
}
func (InteractionType) TableName() string { return "InteractionType" }

type Profile struct {
	IDProfile   int64     `gorm:"column:idProfile;primaryKey;autoIncrement"`
	IDRole      int64     `gorm:"column:idRole;not null;default:2;index:index_idRole_profileTable"`
	Username    *string   `gorm:"column:username;uniqueIndex:unique_username_profile"`
	Email       *string   `gorm:"column:email;uniqueIndex:unique_email_profile"`
	Password    *string   `gorm:"column:password"`
	PasswordRaw *string   `gorm:"column:passwordRaw"`
	DateCreated time.Time `gorm:"column:date_created;default:CURRENT_TIMESTAMP"`
	DateUpdated time.Time `gorm:"column:date_updated;default:CURRENT_TIMESTAMP"`
}
func (Profile) TableName() string { return "Profile" }

type RemoveUserData struct {
	IDRemoveUserData int64     `gorm:"column:id_removeuserdata;primaryKey;autoIncrement"`
	IDProfile        int64     `gorm:"column:id_profile;not null"`
	IDSession        int64     `gorm:"column:id_session;not null"`
	Timestamp        time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`
}
func (RemoveUserData) TableName() string { return "RemoveUserData" }

type Role struct {
	IDRole int64  `gorm:"column:idRole;primaryKey;autoIncrement"`
	Role   string `gorm:"column:role;not null"`
}
func (Role) TableName() string { return "Role" }

type SearchType struct {
	IDSearchType int64  `gorm:"column:idSearchType;primaryKey;autoIncrement"`
	Type         string `gorm:"column:type;not null"`
}
func (SearchType) TableName() string { return "SearchType" }

type SelectRange struct {
	IDInteraction int64   `gorm:"column:idInteraction;not null;index:fk_Click_Interaction2_idx"`
	IDSuggestion  int64   `gorm:"column:idSuggestion;not null;index:fk_Click_Suggestion2_idx"`
	RowValues     *string `gorm:"column:rowvalues"`
}
func (SelectRange) TableName() string { return "SelectRange" }

type Session struct {
	IDSession int64      `gorm:"column:idSession;primaryKey;autoIncrement"`
	IDProfile *int64     `gorm:"column:idProfile;index:fk_Session_Profile1_idx"`
	Start     time.Time  `gorm:"column:start;default:CURRENT_TIMESTAMP"`
	End       *time.Time `gorm:"column:end"`
}
func (Session) TableName() string { return "Session" }

type SuggestionType struct {
	IDSuggestionType int64   `gorm:"column:idSuggestionType;primaryKey;autoIncrement"`
	IDDataType       int64   `gorm:"column:idDataType;not null;index:fk_SuggestionType_DataType1_idx"`
	Name             *string `gorm:"column:name"`
	IsActive         int64   `gorm:"column:isActive;not null;default:1"`
	Regex            string  `gorm:"column:regex;not null;default:.*"`
	MakesRowUnique   *int64  `gorm:"column:makesRowUnique;default:0"`
	CanBeBlank       int64   `gorm:"column:canBeBlank;not null;default:0"`
	IsFreeEdit       int64   `gorm:"column:isFreeEdit;not null;default:1"`
	IsDate           int64   `gorm:"column:isDate;not null;default:0"`
	IsLink           int64   `gorm:"column:isLink;not null;default:0"`
	IsCurrency       int64   `gorm:"column:isCurrency;not null;default:0"`
	IsEditable       int64   `gorm:"column:isEditable;not null;default:1"`
	IsPrivate        int64   `gorm:"column:isPrivate;not null;default:0"`
	ColumnOrder      *int64  `gorm:"column:columnOrder"`
}
func (SuggestionType) TableName() string { return "SuggestionType" }

type CopyColumn struct {
	IDInteraction    int64 `gorm:"column:idInteraction;primaryKey"`
	IDSuggestionType int64 `gorm:"column:idSuggestionType;primaryKey"`

	InteractionRef    Interaction    `gorm:"constraint:Name:_fk_idInteraction_CopyColumn_alsdfh12356;foreignKey:IDInteraction;references:IDInteraction"`
	SuggestionTypeRef SuggestionType `gorm:"constraint:Name:_fk_idSuggestionType_CopyColumn_alsdfh12356;foreignKey:IDSuggestionType;references:IDSuggestionType"`
}
func (CopyColumn) TableName() string { return "CopyColumn" }

type Search struct {
	IDInteraction    int64   `gorm:"column:idInteraction;not null;index:fk_Search_Interaction1_idx"`
	IDSuggestionType int64   `gorm:"column:idSuggestionType;not null;index:_fk_idSuggestionType_12835gv"`
	IDSearchType     int64   `gorm:"column:idSearchType;not null;default:3;index:_index_idSearchType_182356"`
	IsPartial        int64   `gorm:"column:isPartial;not null;default:1"`
	IsMulti          int64   `gorm:"column:isMulti;not null;default:0"`
	IsFromURL        int64   `gorm:"column:isFromUrl;not null;default:0"`
	Value            *string `gorm:"column:value"`
	MatchedValues    []byte  `gorm:"column:matchedValues"`

	InteractionRef    Interaction    `gorm:"constraint:Name:fk_Search_Interaction1;foreignKey:IDInteraction;references:IDInteraction"`
	SuggestionTypeRef SuggestionType `gorm:"constraint:Name:search_ibfk_1;foreignKey:IDSuggestionType;references:IDSuggestionType"`
	SearchTypeRef     SearchType     `gorm:"constraint:Name:fk_id_search_type_1;foreignKey:IDSearchType;references:IDSearchType"`
}
func (Search) TableName() string { return "Search" }

type SearchMulti struct {
	IDInteraction    int64   `gorm:"column:idInteraction;primaryKey;index:fk_SearchMulti_Interaction1_idx"`
	IDSuggestionType int64   `gorm:"column:idSuggestionType;primaryKey"`
	IDSearchType     int64   `gorm:"column:idSearchType;not null;default:3"`
	Value            *string `gorm:"column:value"`

	SuggestionTypeRef SuggestionType `gorm:"constraint:Name:SearchMulti_ibfk_11231;foreignKey:IDSuggestionType;references:IDSuggestionType"`
	InteractionRef    Interaction    `gorm:"constraint:Name:fk_SearchMulti_Interaction112312;foreignKey:IDInteraction;references:IDInteraction"`
	SearchTypeRef     SearchType     `gorm:"constraint:Name:fk_id_Search_type_112312;foreignKey:IDSearchType;references:IDSearchType"`
}
func (SearchMulti) TableName() string { return "SearchMulti" }

type Sort struct {
	IDInteraction    int64 `gorm:"column:idInteraction;primaryKey;index:fk_Sort_Interaction1_idx"`
	IDSuggestionType int64 `gorm:"column:idSuggestionType;primaryKey;index:fk_Sort_SuggestionType1_idx"`
	IsAsc            int64 `gorm:"column:isAsc;not null;default:1"`
	IsTrigger        int64 `gorm:"column:isTrigger;not null;default:1"`
	IsMulti          int64 `gorm:"column:isMulti;not null;default:0"`

	InteractionRef    Interaction    `gorm:"constraint:Name:_fk_idInteraction_1827365;foreignKey:IDInteraction;references:IDInteraction"`
	SuggestionTypeRef SuggestionType `gorm:"constraint:Name:_fk_idSuggestionType_2827365;foreignKey:IDSuggestionType;references:IDSuggestionType"`
}
func (Sort) TableName() string { return "Sort" }

type SuggestionTypeValues struct {
	IDSuggestionType int64  `gorm:"column:idSuggestionType;not null;index:fk_idSuggestioType_1827635;uniqueIndex:PRIMARY_id_and_value;uniqueIndex:idSuggestionType"`
	Value            string `gorm:"column:value;uniqueIndex:PRIMARY_id_and_value;uniqueIndex:idSuggestionType"`
	Active           int64  `gorm:"column:active;not null;default:1"`
}
func (SuggestionTypeValues) TableName() string { return "SuggestionTypeValues" }

type UniqueId struct {
	IDUniqueID int64   `gorm:"column:idUniqueID;primaryKey;autoIncrement"`
	Active     int64   `gorm:"column:active;not null;default:1"`
	Notes      *string `gorm:"column:notes"`
}
func (UniqueId) TableName() string { return "UniqueId" }

type Comments struct {
	IDComment     int64  `gorm:"column:idComment;primaryKey;autoIncrement"`
	IDInteraction int64  `gorm:"column:idInteraction;not null"`
	IDUniqueID    int64  `gorm:"column:idUniqueID;not null"`
	Comment       string `gorm:"column:comment;not null"`
	VoteUp        int64  `gorm:"column:voteUp;not null;default:0"`
	VoteDown      int64  `gorm:"column:voteDown;not null;default:0"`

	InteractionRef Interaction `gorm:"constraint:Name:Comments___fk_idInteraction_ksjdfba87aidsb;foreignKey:IDInteraction;references:IDInteraction"`
	UniqueIDRef    UniqueId    `gorm:"constraint:Name:Comments___fk_iduniqid_oq83eyfgqwuyofhba;foreignKey:IDUniqueID;references:IDUniqueID"`
}
func (Comments) TableName() string { return "Comments" }

type CommentVote struct {
	IDCommentVote int64   `gorm:"column:idCommentVote;primaryKey;autoIncrement"`
	IDInteraction int64   `gorm:"column:idInteraction;not null"`
	IDComment     int64   `gorm:"column:idComment;not null"`
	Vote          string  `gorm:"column:vote;not null;check:vote IN ('voteUp','voteUp-deselect','voteDown','voteDown-deselect')"`
	Selected      *int64  `gorm:"column:selected"`

	CommentRef     Comments     `gorm:"constraint:Name:table_name___fk_comments_akhsdfashjld;foreignKey:IDComment;references:IDComment"`
	InteractionRef Interaction  `gorm:"constraint:Name:table_name___fk_interaction_akdhfa;foreignKey:IDInteraction;references:IDInteraction"`
}
func (CommentVote) TableName() string { return "CommentVote" }

type CommentsView struct {
	IDCommentsView int64  `gorm:"column:idCommentsView;primaryKey;autoIncrement"`
	IDUniqueID     *int64 `gorm:"column:idUniqueID"`
	IDInteraction  *int64 `gorm:"column:idInteraction"`

	UniqueIDRef    UniqueId    `gorm:"constraint:Name:CommentsView___fk_comments_alksdhfga1231;foreignKey:IDUniqueID;references:IDUniqueID"`
	InteractionRef Interaction `gorm:"constraint:Name:CommentsView___fk_interaction_aljhsdfg5123;foreignKey:IDInteraction;references:IDInteraction"`
}
func (CommentsView) TableName() string { return "CommentsView" }

type Databaits struct {
	IDDatabait             int64     `gorm:"column:idDatabait;primaryKey;autoIncrement"`
	IDInteraction          int64     `gorm:"column:idInteraction;not null"`
	IDUniqueID             *int64    `gorm:"column:idUniqueID"`
	IDDatabaitTemplateType int64     `gorm:"column:idDatabaitTemplateType;not null"`
	IDDatabaitCreateType   int64     `gorm:"column:idDatabaitCreateType;not null"`
	Databait               string    `gorm:"column:databait;not null"`
	Columns                *string   `gorm:"column:columns"`
	Vals                   *string   `gorm:"column:vals"`
	Notes                  string    `gorm:"column:notes;not null"`
	Created                time.Time `gorm:"column:created;not null;default:CURRENT_TIMESTAMP"`
	Closed                 time.Time `gorm:"column:closed;not null;default:'0000-00-00 00:00:00'"`
	NextAction             *int64    `gorm:"column:nextAction"`

	NextActionRef   DatabaitNextAction   `gorm:"constraint:Name:_fk_databaits_nextAction_b6345das;foreignKey:NextAction;references:IDDatabaitNextAction"`
	CreateTypeRef   DatabaitCreateType   `gorm:"constraint:Name:_fk_idDatabaitCreateType_b6345das;foreignKey:IDDatabaitCreateType;references:IDDatabaitCreateType"`
	TemplateTypeRef DatabaitTemplateType `gorm:"constraint:Name:_fk_idDatabaitTemplateType_b6345das;foreignKey:IDDatabaitTemplateType;references:IDDatabaitTemplateType"`
	InteractionRef  Interaction          `gorm:"constraint:Name:_fk_idInteraction_Databaits_a6a3344;foreignKey:IDInteraction;references:IDInteraction"`
	UniqueIDRef     UniqueId             `gorm:"constraint:Name:_fk_idUniqueID_Databaits_b111edss;foreignKey:IDUniqueID;references:IDUniqueID"`
}
func (Databaits) TableName() string { return "Databaits" }

type DatabaitVisit struct {
	IDInteraction int64   `gorm:"column:idInteraction;not null;uniqueIndex:_unique_id_interaction_databaitvisit"`
	IDDatabait    int64   `gorm:"column:idDatabait;not null"`
	Source        *string `gorm:"column:source"`

	DatabaitRef    Databaits   `gorm:"constraint:Name:_fk_idDatabait_DatabaitVisit_b123gda;foreignKey:IDDatabait;references:IDDatabait"`
	InteractionRef Interaction `gorm:"constraint:Name:_fk_idInteraction_DatabaitVisit_asdhjk16341;foreignKey:IDInteraction;references:IDInteraction"`
}
func (DatabaitVisit) TableName() string { return "DatabaitVisit" }

type EditDelRow struct {
	IDUniqueID int64  `gorm:"column:idUniqueID;not null"`
	IDEdit     int64  `gorm:"column:idEdit;not null"`
	Comment    string `gorm:"column:comment;not null"`

	EditRef     Edit     `gorm:"constraint:Name:_fk_edit_akjdhaas;foreignKey:IDEdit;references:IDEdit"`
	UniqueIDRef UniqueId `gorm:"constraint:Name:_fk_UniqueID_1027836asda;foreignKey:IDUniqueID;references:IDUniqueID"`
}
func (EditDelRow) TableName() string { return "Edit_DelRow" }

type HelpUs struct {
	IDHelpUs      int64     `gorm:"column:idHelpUs;primaryKey;autoIncrement"`
	IDInteraction int64     `gorm:"column:idInteraction;not null"`
	IDUniqueID    int64     `gorm:"column:idUniqueID;not null"`
	HelpUsType    string    `gorm:"column:helpUsType;not null"`
	Question      string    `gorm:"column:question;not null"`
	Answer        *string   `gorm:"column:answer"`
	Start         time.Time `gorm:"column:start;not null;default:CURRENT_TIMESTAMP"`
	Answered      time.Time `gorm:"column:answered;not null;default:'0000-00-00 00:00:00'"`
	ShowAnother   time.Time `gorm:"column:showAnother;not null;default:'0000-00-00 00:00:00'"`
	Closed        time.Time `gorm:"column:closed;not null;default:'0000-00-00 00:00:00'"`

	InteractionRef Interaction `gorm:"constraint:Name:HelpUs___fk_interaction_idInteraction;foreignKey:IDInteraction;references:IDInteraction"`
	UniqueIDRef    UniqueId    `gorm:"constraint:Name:HelpUs___fk_interaction_idUniqueID;foreignKey:IDUniqueID;references:IDUniqueID"`
}
func (HelpUs) TableName() string { return "HelpUs" }

type Suggestions struct {
	IDSuggestion     int64     `gorm:"column:idSuggestion;primaryKey;autoIncrement;index:idSuggestion_2;uniqueIndex:idSuggestion"`
	IDSuggestionType int64     `gorm:"column:idSuggestionType;not null;index:fk_Suggestion_SuggestionType1_idx"`
	IDUniqueID       int64     `gorm:"column:idUniqueID;not null;index:fk_Suggestion_UniqueID_idx"`
	IDProfile        int64     `gorm:"column:idProfile;not null;default:2"`
	Suggestion       string    `gorm:"column:suggestion;not null;default:''"`
	Active           int64     `gorm:"column:active;not null;default:1"`
	Confidence       *int64    `gorm:"column:confidence"`
	LastUpdated      time.Time `gorm:"column:last_updated;not null;default:CURRENT_TIMESTAMP"`

	SuggestionTypeRef SuggestionType `gorm:"constraint:Name:fk_idSuggestionType_123687;foreignKey:IDSuggestionType;references:IDSuggestionType"`
	UniqueIDRef       UniqueId       `gorm:"constraint:Name:fk_unique_row_id;foreignKey:IDUniqueID;references:IDUniqueID"`
}
func (Suggestions) TableName() string { return "Suggestions" }

type Copy struct {
	IDInteraction int64 `gorm:"column:idInteraction;primaryKey"`
	IDSuggestion  int64 `gorm:"column:idSuggestion;primaryKey"`

	InteractionRef Interaction `gorm:"constraint:Name:_fk_idInteraction_4447654;foreignKey:IDInteraction;references:IDInteraction"`
	SuggestionRef  Suggestions `gorm:"constraint:Name:_fk_idSuggestion_4417654;foreignKey:IDSuggestion;references:IDSuggestion"`
}
func (Copy) TableName() string { return "Copy" }

type EditNewRow struct {
	IDEdit       int64 `gorm:"column:idEdit;not null;uniqueIndex:idEdit_"`
	IDSuggestion int64 `gorm:"column:idSuggestion;not null;uniqueIndex:idEdit_"`
	IsCorrect    int64 `gorm:"column:isCorrect;not null;default:2"`

	EditRef       Edit        `gorm:"constraint:Name:_fk_idEdit_from_Edit_NewRow_asdkl123;foreignKey:IDEdit;references:IDEdit"`
	SuggestionRef Suggestions `gorm:"constraint:Name:_fk_idSuggestion_from_Edit_NewRow_asdkl123;foreignKey:IDSuggestion;references:IDSuggestion"`
}
func (EditNewRow) TableName() string { return "Edit_NewRow" }

type Paste struct {
	IDInteraction         int64   `gorm:"column:idInteraction;primaryKey;not null;autoIncrement:false"`
	PasteValue            string  `gorm:"column:pasteValue;not null"`
	CopyCellIDSuggestion  *int64  `gorm:"column:copyCellIdSuggestion"`
	CopyCellValue         *string `gorm:"column:copyCellValue"`
	PasteCellIDSuggestion int64   `gorm:"column:pasteCellIdSuggestion;not null"`
	PasteCellValue        string  `gorm:"column:pasteCellValue;not null"`

	CopySuggestionRef  Suggestions `gorm:"constraint:Name:_fk_copy_id_suggestion_197823;foreignKey:CopyCellIDSuggestion;references:IDSuggestion"`
	PasteSuggestionRef Suggestions `gorm:"constraint:Name:_fk_paste_id_suggestion_197823;foreignKey:PasteCellIDSuggestion;references:IDSuggestion"`
}
func (Paste) TableName() string { return "Paste" }

type SearchGoogle struct {
	IDInteraction int64  `gorm:"column:IdInteraction;not null;uniqueIndex:IdInteraction"`
	IDUniqueID    int64  `gorm:"column:idUniqueID;not null"`
	IDSuggestion  int64  `gorm:"column:idSuggestion;not null"`
	SearchValues  string `gorm:"column:searchValues;not null"`

	InteractionRef Interaction `gorm:"constraint:Name:_fk_idInteraction_SearchGoogle_a645das;foreignKey:IDInteraction;references:IDInteraction"`
	SuggestionRef  Suggestions `gorm:"constraint:Name:_fk_idSuggestion_SearchGoogle_a645das;foreignKey:IDSuggestion;references:IDSuggestion"`
	UniqueIDRef    UniqueId    `gorm:"constraint:Name:_fk_idUniqueID_SearchGoogle_a645das;foreignKey:IDUniqueID;references:IDUniqueID"`
}
func (SearchGoogle) TableName() string { return "SearchGoogle" }

type ViewChange struct {
	IDInteraction int64  `gorm:"column:idInteraction;primaryKey"`
	ViewName      string `gorm:"column:viewname;primaryKey"`
}
func (ViewChange) TableName() string { return "ViewChange" }

type Visit struct {
	IDVisit       int64   `gorm:"column:idVisit;primaryKey;autoIncrement;uniqueIndex:Visit_idVisit_uindex"`
	IDInteraction int64   `gorm:"column:idInteraction;not null"`
	Source        *string `gorm:"column:source"`
	SearchCol     *string `gorm:"column:searchCol"`
	SearchVal     *string `gorm:"column:searchVal"`

	InteractionRef Interaction `gorm:"constraint:Name:Visit___fk_interaction;foreignKey:IDInteraction;references:IDInteraction"`
}
func (Visit) TableName() string { return "Visit" }

type Sessions struct {
	SessionID string  `gorm:"column:session_id;primaryKey"`
	Expires   int64   `gorm:"column:expires;not null"`
	Data      *string `gorm:"column:data"`
}
func (Sessions) TableName() string { return "sessions" }
