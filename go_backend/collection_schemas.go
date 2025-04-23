package main

import (
    "log"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/core"
    //"github.com/pocketbase/pocketbase/tools/types"
)

// create table Alias
// (
//     idAlias      int auto_increment primary key,
//     idSuggestion int           not null,
//     alias        varchar(500)  null,
//     count        int default 1 not null,
//     constraint unique_index unique (idSuggestion, alias(250))
// )
// charset = utf8;
// create index fk_Alias_Suggestion1_idx on Alias (idSuggestion);
func setupAliasCollection(app *pocketbase.PocketBase) error {
	// SQL: CREATE TABLE Alias
	if _, err := app.FindCollectionByNameOrId("Alias"); err == nil {
		log.Println("Collection 'Alias' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Alias")

	// SQL: idAlias int auto_increment primary key
	// (Handled by PocketBase's internal unique "id" field)

	// SQL: idSuggestion int NOT NULL, foreign key references Suggestions(idSuggestion)
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})

	// SQL: alias varchar(500) NULL
	collection.Fields.Add(&core.TextField{
		Name:     "alias",
		Required: false,
		Max:      500,
	})

	// SQL: count int DEFAULT 1 NOT NULL
	collection.Fields.Add(&core.NumberField{
		Name:     "count",
		Required: true,
	})

	// SQL: constraint unique_index unique (idSuggestion, alias(250))
	collection.AddIndex("unique_index", true, "idSuggestion, alias", "")

	// SQL: create index fk_Alias_Suggestion1_idx on Alias (idSuggestion)
	collection.AddIndex("fk_Alias_Suggestion1_idx", false, "idSuggestion", "")

	// SQL: charset = utf8 (handled globally by SQLite)
	if err := app.Save(collection); err != nil {
		return err
	}

	log.Println("Created 'Alias' collection with SQL mapping")
	return nil
}

// create table Click
// (
//     idInteraction int          not null primary key,
//     idSuggestion  int          not null,
//     rowvalues     varchar(500) null
// )
// charset = utf8;
// create index fk_Click_Interaction1_idx on Click (idInteraction);
// create index fk_Click_Suggestion1_idx on Click (idSuggestion);
func setupClickCollection(app *pocketbase.PocketBase) error {
	// SQL: CREATE TABLE Click
	if _, err := app.FindCollectionByNameOrId("Click"); err == nil {
		log.Println("Collection 'Click' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Click")

	// SQL: idInteraction int NOT NULL PRIMARY KEY
	// (Handled by PocketBase's internal unique "id" field)
	// Add relation field for idInteraction referencing Interaction
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})

	// SQL: idSuggestion int NOT NULL, foreign key references Suggestions(idSuggestion)
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})

	// SQL: rowvalues varchar(500) NULL
	collection.Fields.Add(&core.TextField{
		Name:     "rowvalues",
		Required: false,
		Max:      500,
	})

	// SQL: create index fk_Click_Interaction1_idx on Click (idInteraction)
	collection.AddIndex("fk_Click_Interaction1_idx", false, "idInteraction", "")

	// SQL: create index fk_Click_Suggestion1_idx on Click (idSuggestion)
	collection.AddIndex("fk_Click_Suggestion1_idx", false, "idSuggestion", "")

	// SQL: charset = utf8 (handled globally by SQLite)
	if err := app.Save(collection); err != nil {
		return err
	}

	log.Println("Created 'Click' collection with SQL mapping")
	return nil
}

// SQL: create table DataType
// (
//     idDataType int auto_increment primary key,
//     type       varchar(45) null
// )
// charset = utf8;
func setupDataTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DataType"); err == nil {
		log.Println("Collection 'DataType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DataType")

	// idDataType is handled internally.
	// SQL: type varchar(45) null
	collection.Fields.Add(&core.TextField{
		Name:     "type",
		Required: false,
		Max:      45,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DataType' collection with SQL mapping")
	return nil
}

// SQL: create table DatabaitCreateType
// (
//     idDatabaitCreateType int auto_increment primary key,
//     type                 varchar(50) null
// )
// charset = utf8;
func setupDatabaitCreateTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DatabaitCreateType"); err == nil {
		log.Println("Collection 'DatabaitCreateType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DatabaitCreateType")

	// idDatabaitCreateType is handled internally.
	// SQL: type varchar(50) null
	collection.Fields.Add(&core.TextField{
		Name:     "type",
		Required: false,
		Max:      50,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DatabaitCreateType' collection with SQL mapping")
	return nil
}

// SQL: create table DatabaitNextAction
// (
//     idDatabaitNextAction int auto_increment primary key,
//     action               varchar(50) null
// )
// charset = utf8;
func setupDatabaitNextActionCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DatabaitNextAction"); err == nil {
		log.Println("Collection 'DatabaitNextAction' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DatabaitNextAction")

	// idDatabaitNextAction is handled internally.
	// SQL: action varchar(50) null
	collection.Fields.Add(&core.TextField{
		Name:     "action",
		Required: false,
		Max:      50,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DatabaitNextAction' collection with SQL mapping")
	return nil
}

// SQL: create table DatabaitTemplateType
// (
//     idDatabaitTemplateType int auto_increment primary key,
//     template               varchar(50) null
// )
// charset = utf8;
func setupDatabaitTemplateTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DatabaitTemplateType"); err == nil {
		log.Println("Collection 'DatabaitTemplateType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DatabaitTemplateType")

	// idDatabaitTemplateType is handled internally.
	// SQL: template varchar(50) null
	collection.Fields.Add(&core.TextField{
		Name:     "template",
		Required: false,
		Max:      50,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DatabaitTemplateType' collection with SQL mapping")
	return nil
}

// SQL: create table DoubleClick
// (
//     idInteraction int not null primary key,
//     idSuggestion  int not null,
//     rowvalues     varchar(500) null
// )
// charset = utf8;
// create index fk_DoubleClick_Interaction1_idx on DoubleClick (idInteraction);
// create index fk_DoubleClick_Suggestion1_idx on DoubleClick (idSuggestion);
func setupDoubleClickCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DoubleClick"); err == nil {
		log.Println("Collection 'DoubleClick' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DoubleClick")

	// SQL: idInteraction int not null primary key
	// Reference Interaction table.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})

	// SQL: idSuggestion int not null, assumed to reference Suggestions table.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})

	// SQL: rowvalues varchar(500) null
	collection.Fields.Add(&core.TextField{
		Name:     "rowvalues",
		Required: false,
		Max:      500,
	})

	// Indexes
	collection.AddIndex("fk_DoubleClick_Interaction1_idx", false, "idInteraction", "")
	collection.AddIndex("fk_DoubleClick_Suggestion1_idx", false, "idSuggestion", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DoubleClick' collection with SQL mapping")
	return nil
}

// SQL: create table Edit_Suggestion
// (
//     idEdit           int not null,
//     idSuggestion     int not null,
//     isPrevSuggestion tinyint(1) not null,
//     isNew            tinyint(1) not null,
//     isChosen         tinyint(1) not null,
//     constraint idEdit unique (idEdit, idSuggestion)
// )
// charset = utf8;
// create index _index_edit_suggestion_idEdit_agsdh1872dg on Edit_Suggestion (idEdit);
// create index _index_edit_suggestion_idSuggestion_agsdh1872dg on Edit_Suggestion (idSuggestion);
func setupEditSuggestionCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Edit_Suggestion"); err == nil {
		log.Println("Collection 'Edit_Suggestion' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Edit_Suggestion")

	// SQL: idEdit int not null, reference Edit table.
	editCollection, err := app.FindCollectionByNameOrId("Edit")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idEdit",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  editCollection.Id,
	})

	// SQL: idSuggestion int not null, assumed to reference Suggestions table.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})

	// SQL: isPrevSuggestion tinyint(1) not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isPrevSuggestion",
		Required: true,
	})
	// SQL: isNew tinyint(1) not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isNew",
		Required: true,
	})
	// SQL: isChosen tinyint(1) not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isChosen",
		Required: true,
	})

	// Unique composite index on (idEdit, idSuggestion)
	collection.AddIndex("idEdit_1", true, "idEdit, idSuggestion", "")
	// Additional indexes:
	collection.AddIndex("_index_edit_suggestion_idEdit_agsdh1872dg_1", false, "idEdit", "")
	collection.AddIndex("_index_edit_suggestion_idSuggestion_agsdh1872dg_1", false, "idSuggestion", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Edit_Suggestion' collection with SQL mapping")
	return nil
}

// SQL: create table EntryType
// (
//     idEntryType int auto_increment primary key,
//     type        varchar(45) null
// )
// charset = utf8;
func setupEntryTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("EntryType"); err == nil {
		log.Println("Collection 'EntryType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("EntryType")

	// idEntryType handled internally.
	// SQL: type varchar(45) null
	collection.Fields.Add(&core.TextField{
		Name:     "type",
		Required: false,
		Max:      45,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'EntryType' collection with SQL mapping")
	return nil
}

// SQL: create table Interaction
// (
//     idInteraction     int auto_increment primary key,
//     idSession         int not null,
//     idInteractionType int not null,
//     timestamp         datetime default current_timestamp() null
// )
// charset = utf8;
// create index fk_Interaction_InteractionType1_idx on Interaction (idInteractionType);
// create index fk_Interaction_Session1_idx on Interaction (idSession);
func setupInteractionCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Interaction"); err == nil {
		log.Println("Collection 'Interaction' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Interaction")

	// SQL: idInteraction handled internally.
	// SQL: idSession int not null, reference Session table.
	sessionCollection, err := app.FindCollectionByNameOrId("Session")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSession",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  sessionCollection.Id,
	})
	// SQL: idInteractionType int not null, reference InteractionType table.
	interactionTypeCollection, err := app.FindCollectionByNameOrId("InteractionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteractionType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionTypeCollection.Id,
	})
	// SQL: timestamp datetime default current_timestamp() null
	collection.Fields.Add(&core.AutodateField{
		Name:        "timestamp",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})

	// Indexes:
	collection.AddIndex("fk_Interaction_InteractionType1_idx", false, "idInteractionType", "")
	collection.AddIndex("fk_Interaction_Session1_idx", false, "idSession", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Interaction' collection with SQL mapping")
	return nil
}

// SQL: create table DatabaitTweet
// (
//     idDatabaitTweet     int auto_increment primary key,
//     idInteraction       int not null,
//     idDatabait          int not null,
//     url                 varchar(2500) not null,
//     likes               int null,
//     retweets            int null,
//     created             timestamp default current_timestamp() not null,
//     nextActionTimestamp timestamp null,
//     nextAction          int null,
//     constraint _fk_idInteraction_DatabaitTweet_t615das foreign key (idInteraction) references Interaction (idInteraction),
//     constraint _fk_nextAction_DatabaitTweet_t615das foreign key (nextAction) references DatabaitNextAction (idDatabaitNextAction)
// )
// charset = utf8;
func setupDatabaitTweetCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DatabaitTweet"); err == nil {
		log.Println("Collection 'DatabaitTweet' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DatabaitTweet")

	// SQL: idDatabaitTweet handled internally.
	// SQL: idInteraction int not null, foreign key referencing Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})

	// SQL: idDatabait int not null, reference Databaits table.
	databaitsCollection, err := app.FindCollectionByNameOrId("Databaits")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idDatabait",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  databaitsCollection.Id,
	})

	// SQL: url varchar(2500) not null
	collection.Fields.Add(&core.URLField{
		Name:     "url",
		Required: true,
	})
	// SQL: likes int null
	collection.Fields.Add(&core.NumberField{
		Name:     "likes",
		Required: false,
	})
	// SQL: retweets int null
	collection.Fields.Add(&core.NumberField{
		Name:     "retweets",
		Required: false,
	})
	// SQL: created timestamp default current_timestamp() not null
	collection.Fields.Add(&core.AutodateField{
		Name:        "created",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// SQL: nextActionTimestamp timestamp null
	collection.Fields.Add(&core.DateField{
		Name:        "nextActionTimestamp",
		Presentable: true,
	})
	// SQL: nextAction int null, foreign key referencing DatabaitNextAction.
	databaitNextActionCollection, err := app.FindCollectionByNameOrId("DatabaitNextAction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "nextAction",
		Required:      false,
		CascadeDelete: false,
		CollectionId:  databaitNextActionCollection.Id,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DatabaitTweet' collection with SQL mapping")
	return nil
}

// SQL: create table Edit
// (
//     IdInteraction int not null,
//     idEdit        int auto_increment primary key,
//     idEntryType   int not null,
//     mode          varchar(25) default 'normal' not null,
//     isCorrect     tinyint(1) default 2 null,
//     constraint _fk_idEntryType_from_edit_asllhg1233 foreign key (idEntryType) references EntryType (idEntryType),
//     constraint _fk_idInteraction_from_edit_asdlhg1235 foreign key (IdInteraction) references Interaction (idInteraction)
// )
// charset = utf8;
// create index idInteraction_index_adfhj126 on Edit (IdInteraction);
func setupEditCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Edit"); err == nil {
		log.Println("Collection 'Edit' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Edit")

	// SQL: IdInteraction int not null, foreign key referencing Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "IdInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})

	// idEdit auto_increment primary key is handled internally.
	// SQL: idEntryType int not null, foreign key referencing EntryType.
	entryTypeCollection, err := app.FindCollectionByNameOrId("EntryType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idEntryType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  entryTypeCollection.Id,
	})

	// SQL: mode varchar(25) default 'normal' not null
	collection.Fields.Add(&core.TextField{
		Name:     "mode",
		Required: true,
		Max:      25,
	})
	// SQL: isCorrect tinyint(1) default 2 null
	collection.Fields.Add(&core.NumberField{
		Name:     "isCorrect",
		Required: false,
	})
	// Index on IdInteraction
	collection.AddIndex("idInteraction_index_adfhj126", false, "IdInteraction", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Edit' collection with SQL mapping")
	return nil
}

// SQL: create table InteractionType
// (
//     idInteractionType int auto_increment primary key,
//     interaction       varchar(45) null
// )
// charset = utf8;
func setupInteractionTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("InteractionType"); err == nil {
		log.Println("Collection 'InteractionType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("InteractionType")

	// idInteractionType handled internally.
	// SQL: interaction varchar(45) null
	collection.Fields.Add(&core.TextField{
		Name:     "interaction",
		Required: false,
		Max:      45,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'InteractionType' collection with SQL mapping")
	return nil
}

// SQL: create table Profile
// (
//     idProfile    int auto_increment primary key,
//     idRole       int default 2 not null,
//     username     varchar(45) null,
//     email        varchar(45) null,
//     password     varchar(500) null,
//     passwordRaw  varchar(100) null,
//     date_created datetime default current_timestamp() null,
//     date_updated datetime default current_timestamp() null,
//     constraint unique_email_profile unique (email),
//     constraint unique_username_profile unique (username)
// )
// charset = utf8;
// create index index_idRole_profileTable on Profile (idRole);
func setupProfileCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Profile"); err == nil {
		log.Println("Collection 'Profile' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Profile")
	// idProfile is handled internally.
	// SQL: idRole int default 2 not null, reference Role table.
	roleCollection, err := app.FindCollectionByNameOrId("Role")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idRole",
		Required:     true,
		CascadeDelete: false,
		CollectionId: roleCollection.Id,
	})
	// SQL: username varchar(45) null
	collection.Fields.Add(&core.TextField{
		Name:     "username",
		Required: false,
		Max:      45,
	})
	// SQL: email varchar(45) null
	collection.Fields.Add(&core.EmailField{
		Name:     "email",
		Required: false,
	})
	// SQL: password varchar(500) null
	collection.Fields.Add(&core.TextField{
		Name:     "password",
		Required: false,
		Max:      500,
	})
	// SQL: passwordRaw varchar(100) null
	collection.Fields.Add(&core.TextField{
		Name:     "passwordRaw",
		Required: false,
		Max:      100,
	})
	// SQL: date_created datetime default current_timestamp() null
	collection.Fields.Add(&core.AutodateField{
		Name:        "date_created",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// SQL: date_updated datetime default current_timestamp() null
	collection.Fields.Add(&core.AutodateField{
		Name:        "date_updated",
		OnCreate:    true,
		OnUpdate:    true,
		Presentable: true,
	})
	// Unique constraints:
	collection.AddIndex("unique_email_profile", true, "email", "")
	collection.AddIndex("unique_username_profile", true, "username", "")
	// Additional index on idRole:
	collection.AddIndex("index_idRole_profileTable", false, "idRole", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Profile' collection with SQL mapping")
	return nil
}

// SQL: create table RemoveUserData
// (
//     id_removeuserdata int auto_increment primary key,
//     id_profile        int not null,
//     id_session        int not null,
//     timestamp         timestamp default current_timestamp() null
// );
func setupRemoveUserDataCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("RemoveUserData"); err == nil {
		log.Println("Collection 'RemoveUserData' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("RemoveUserData")
	// id_removeuserdata is handled internally.
	// SQL: id_profile int not null, reference Profile.
	profileCollection, err := app.FindCollectionByNameOrId("Profile")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "id_profile",
		Required:     true,
		CascadeDelete: false,
		CollectionId: profileCollection.Id,
	})
	// SQL: id_session int not null, reference Session.
	sessionCollection, err := app.FindCollectionByNameOrId("Session")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "id_session",
		Required:     true,
		CascadeDelete: false,
		CollectionId: sessionCollection.Id,
	})
	// SQL: timestamp timestamp default current_timestamp() null
	collection.Fields.Add(&core.AutodateField{
		Name:        "timestamp",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'RemoveUserData' collection with SQL mapping")
	return nil
}

// SQL: create table Role
// (
//     idRole int auto_increment primary key,
//     role   varchar(20) not null
// )
// charset = utf8;
func setupRoleCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Role"); err == nil {
		log.Println("Collection 'Role' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Role")
	// SQL: role varchar(20) not null
	collection.Fields.Add(&core.TextField{
		Name:     "role",
		Required: true,
		Max:      20,
	})
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Role' collection with SQL mapping")
	return nil
}

// SQL: create table SearchType
// (
//     idSearchType int auto_increment primary key,
//     type         varchar(20) not null
// )
// charset = utf8;
func setupSearchTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("SearchType"); err == nil {
		log.Println("Collection 'SearchType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("SearchType")
	// SQL: type varchar(20) not null
	collection.Fields.Add(&core.TextField{
		Name:     "type",
		Required: true,
		Max:      20,
	})
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'SearchType' collection with SQL mapping")
	return nil
}

// SQL: create table SelectRange
// (
//     idInteraction int not null,
//     idSuggestion  int not null,
//     rowvalues     varchar(500) null
// )
// charset = utf8;
// create index fk_Click_Interaction2_idx on SelectRange (idInteraction);
// create index fk_Click_Suggestion2_idx on SelectRange (idSuggestion);
func setupSelectRangeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("SelectRange"); err == nil {
		log.Println("Collection 'SelectRange' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("SelectRange")
	// SQL: idInteraction int not null, reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idInteraction",
		Required:     true,
		CascadeDelete: false,
		CollectionId: interactionCollection.Id,
	})
	// SQL: idSuggestion int not null, reference Suggestions.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idSuggestion",
		Required:     true,
		CascadeDelete: false,
		CollectionId: suggestionsCollection.Id,
	})
	// SQL: rowvalues varchar(500) null
	collection.Fields.Add(&core.TextField{
		Name:     "rowvalues",
		Required: false,
		Max:      500,
	})
	// Indexes
	collection.AddIndex("fk_Click_Interaction2_idx", false, "idInteraction", "")
	collection.AddIndex("fk_Click_Suggestion2_idx", false, "idSuggestion", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'SelectRange' collection with SQL mapping")
	return nil
}

// SQL: create table Session
// (
//     idSession int auto_increment primary key,
//     idProfile int null,
//     start     datetime default current_timestamp() null,
//     end       datetime null
// )
// charset = utf8;
// create index fk_Session_Profile1_idx on Session (idProfile);
func setupSessionCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Session"); err == nil {
		log.Println("Collection 'Session' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Session")
	// SQL: idSession handled internally.
	// SQL: idProfile int null, reference Profile.
	profileCollection, err := app.FindCollectionByNameOrId("Profile")
	if err == nil {
		collection.Fields.Add(&core.RelationField{
			Name:         "idProfile",
			Required:     false,
			CascadeDelete: false,
			CollectionId: profileCollection.Id,
		})
	} else {
		collection.Fields.Add(&core.NumberField{
			Name:     "idProfile",
			Required: false,
		})
	}
	// SQL: start datetime default current_timestamp() null
	collection.Fields.Add(&core.AutodateField{
		Name:        "start",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// SQL: end datetime null
	collection.Fields.Add(&core.DateField{
		Name:        "end",
		Presentable: true,
	})
	// Index on idProfile:
	collection.AddIndex("fk_Session_Profile1_idx", false, "idProfile", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Session' collection with SQL mapping")
	return nil
}

// SQL: create table SuggestionType
// (
//     idSuggestionType int auto_increment primary key,
//     idDataType       int not null,
//     name             varchar(45) null,
//     isActive         tinyint(1) default 1 not null,
//     regex            varchar(150) default '.*' not null,
//     makesRowUnique   tinyint(1) default 0 null,
//     canBeBlank       tinyint(1) default 0 not null,
//     isFreeEdit       tinyint(1) default 1 not null,
//     isDate           tinyint(1) default 0 not null,
//     isLink           tinyint(1) default 0 not null,
//     isCurrency       tinyint(1) default 0 not null,
//     isEditable       int default 1 not null,
//     isPrivate        int default 0 not null,
//     columnOrder      int null
// )
// charset = utf8;
// create index fk_SuggestionType_DataType1_idx on SuggestionType (idDataType);
func setupSuggestionTypeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("SuggestionType"); err == nil {
		log.Println("Collection 'SuggestionType' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("SuggestionType")
	// SQL: idDataType int not null, reference DataType.
	dataTypeCollection, err := app.FindCollectionByNameOrId("DataType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idDataType",
		Required:     true,
		CascadeDelete: false,
		CollectionId: dataTypeCollection.Id,
	})
	// SQL: name varchar(45) null
	collection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: false,
		Max:      45,
	})
	// SQL: isActive tinyint(1) default 1 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isActive",
		Required: true,
	})
	// SQL: regex varchar(150) default '.*' not null
	collection.Fields.Add(&core.TextField{
		Name:     "regex",
		Required: true,
		Max:      150,
	})
	// SQL: makesRowUnique tinyint(1) default 0 null
	collection.Fields.Add(&core.BoolField{
		Name:     "makesRowUnique",
		Required: false,
	})
	// SQL: canBeBlank tinyint(1) default 0 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "canBeBlank",
		Required: true,
	})
	// SQL: isFreeEdit tinyint(1) default 1 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isFreeEdit",
		Required: true,
	})
	// SQL: isDate tinyint(1) default 0 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isDate",
		Required: true,
	})
	// SQL: isLink tinyint(1) default 0 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isLink",
		Required: true,
	})
	// SQL: isCurrency tinyint(1) default 0 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isCurrency",
		Required: true,
	})
	// SQL: isEditable int default 1 not null
	collection.Fields.Add(&core.NumberField{
		Name:     "isEditable",
		Required: true,
	})
	// SQL: isPrivate int default 0 not null
	collection.Fields.Add(&core.NumberField{
		Name:     "isPrivate",
		Required: true,
	})
	// SQL: columnOrder int null
	collection.Fields.Add(&core.NumberField{
		Name:     "columnOrder",
		Required: false,
	})
	// Index on idDataType
	collection.AddIndex("fk_SuggestionType_DataType1_idx", false, "idDataType", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'SuggestionType' collection with SQL mapping")
	return nil
}

// SQL: create table CopyColumn
// (
//     idInteraction    int not null,
//     idSuggestionType int not null,
//     primary key (idInteraction, idSuggestionType),
//     constraint _fk_idInteraction_CopyColumn_alsdfh12356 foreign key (idInteraction) references Interaction (idInteraction),
//     constraint _fk_idSuggestionType_CopyColumn_alsdfh12356 foreign key (idSuggestionType) references SuggestionType (idSuggestionType)
// )
// charset = utf8;
func setupCopyColumnCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("CopyColumn"); err == nil {
		log.Println("Collection 'CopyColumn' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("CopyColumn")
	// SQL: idInteraction int not null, reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idInteraction",
		Required:     true,
		CascadeDelete: false,
		CollectionId: interactionCollection.Id,
	})
	// SQL: idSuggestionType int not null, reference SuggestionType.
	suggestionTypeCollection, err := app.FindCollectionByNameOrId("SuggestionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idSuggestionType",
		Required:     true,
		CascadeDelete: false,
		CollectionId: suggestionTypeCollection.Id,
	})
	// Composite primary key as unique index:
	collection.AddIndex("primary_1", true, "idInteraction, idSuggestionType", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'CopyColumn' collection with SQL mapping")
	return nil
}

// SQL: create table Search
// (
//     idInteraction    int not null,
//     idSuggestionType int not null,
//     idSearchType     int default 3 not null,
//     isPartial        tinyint(1) default 1 not null,
//     isMulti          int(1) default 0 not null,
//     isFromUrl        int(1) default 0 not null,
//     value            varchar(150) null,
//     matchedValues    blob null,
//     constraint fk_Search_Interaction1 foreign key (idInteraction) references Interaction (idInteraction),
//     constraint fk_id_search_type_1 foreign key (idSearchType) references SearchType (idSearchType),
//     constraint search_ibfk_1 foreign key (idSuggestionType) references SuggestionType (idSuggestionType)
// )
// charset = utf8;
// create index _fk_idSuggestionType_12835gv on Search (idSuggestionType);
// create index _index_idSearchType_182356 on Search (idSearchType);
// create index fk_Search_Interaction1_idx on Search (idInteraction);
func setupSearchCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Search"); err == nil {
		log.Println("Collection 'Search' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Search")
	// SQL: idInteraction int not null, reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idInteraction",
		Required:     true,
		CascadeDelete: false,
		CollectionId: interactionCollection.Id,
	})
	// SQL: idSuggestionType int not null, reference SuggestionType.
	suggestionTypeCollection, err := app.FindCollectionByNameOrId("SuggestionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idSuggestionType",
		Required:     true,
		CascadeDelete: false,
		CollectionId: suggestionTypeCollection.Id,
	})
	// SQL: idSearchType int default 3 not null, reference SearchType.
	searchTypeCollection, err := app.FindCollectionByNameOrId("SearchType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idSearchType",
		Required:     true,
		CascadeDelete: false,
		CollectionId: searchTypeCollection.Id,
	})
	// SQL: isPartial tinyint(1) default 1 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isPartial",
		Required: true,
	})
	// SQL: isMulti int(1) default 0 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isMulti",
		Required: true,
	})
	// SQL: isFromUrl int(1) default 0 not null
	collection.Fields.Add(&core.BoolField{
		Name:     "isFromUrl",
		Required: true,
	})
	// SQL: value varchar(150) null
	collection.Fields.Add(&core.TextField{
		Name:     "value",
		Required: false,
		Max:      150,
	})
	// SQL: matchedValues blob null (using TextField)
	collection.Fields.Add(&core.TextField{
		Name:     "matchedValues",
		Required: false,
	})
	// Indexes:
	collection.AddIndex("_fk_idSuggestionType_12835gv_1", false, "idSuggestionType", "")
	collection.AddIndex("_index_idSearchType_182356_1", false, "idSearchType", "")
	collection.AddIndex("fk_Search_Interaction1_idx_1", false, "idInteraction", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Search' collection with SQL mapping")
	return nil
}

// SQL: create table SearchMulti
// (
//     idInteraction    int not null,
//     idSuggestionType int not null,
//     idSearchType     int default 3 not null,
//     value            varchar(150) null,
//     primary key (idInteraction, idSuggestionType),
//     constraint SearchMulti_ibfk_11231 foreign key (idSuggestionType) references SuggestionType (idSuggestionType),
//     constraint fk_SearchMulti_Interaction112312 foreign key (idInteraction) references Interaction (idInteraction),
//     constraint fk_id_Search_type_112312 foreign key (idSearchType) references SearchType (idSearchType)
// )
// charset = utf8;
// create index _fk_idSuggestionType_12835gv on SearchMulti (idSuggestionType);
// create index _index_idSearchType_182356 on SearchMulti (idSearchType);
// create index fk_SearchMulti_Interaction1_idx on SearchMulti (idInteraction);
func setupSearchMultiCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("SearchMulti"); err == nil {
		log.Println("Collection 'SearchMulti' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("SearchMulti")
	// SQL: idInteraction int not null, reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idInteraction",
		Required:     true,
		CascadeDelete: false,
		CollectionId: interactionCollection.Id,
	})
	// SQL: idSuggestionType int not null, reference SuggestionType.
	suggestionTypeCollection, err := app.FindCollectionByNameOrId("SuggestionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idSuggestionType",
		Required:     true,
		CascadeDelete: false,
		CollectionId: suggestionTypeCollection.Id,
	})
	// SQL: idSearchType int default 3 not null, reference SearchType.
	searchTypeCollection, err := app.FindCollectionByNameOrId("SearchType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "idSearchType",
		Required:     true,
		CascadeDelete: false,
		CollectionId: searchTypeCollection.Id,
	})
	// SQL: value varchar(150) null
	collection.Fields.Add(&core.TextField{
		Name:     "value",
		Required: false,
		Max:      150,
	})
	// Composite primary key on (idInteraction, idSuggestionType)
	collection.AddIndex("primary_2", true, "idInteraction, idSuggestionType", "")
	// Additional indexes:
	collection.AddIndex("_fk_idSuggestionType_12835gv_2", false, "idSuggestionType", "")
	collection.AddIndex("_index_idSearchType_182356_2", false, "idSearchType", "")
	collection.AddIndex("fk_SearchMulti_Interaction1_idx_2", false, "idInteraction", "")
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'SearchMulti' collection with SQL mapping")
	return nil
}

// SQL: create table Sort
// (
//     idInteraction    int not null,
//     idSuggestionType int not null,
//     isAsc            tinyint(1) default 1 not null,
//     isTrigger        tinyint(1) default 1 not null,
//     isMulti          tinyint(1) default 0 not null,
//     primary key (idInteraction, idSuggestionType),
//     constraint _fk_idInteraction_1827365 foreign key (idInteraction) references Interaction (idInteraction),
//     constraint _fk_idSuggestionType_2827365 foreign key (idSuggestionType) references SuggestionType (idSuggestionType)
// )
// charset = utf8;
// create index fk_Sort_Interaction1_idx on Sort (idInteraction);
// create index fk_Sort_SuggestionType1_idx on Sort (idSuggestionType);
func setupSortCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Sort"); err == nil {
		log.Println("Collection 'Sort' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Sort")

	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})

	// idSuggestionType: reference SuggestionType.
	suggestionTypeCollection, err := app.FindCollectionByNameOrId("SuggestionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestionType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionTypeCollection.Id,
	})

	// Boolean fields for isAsc, isTrigger, isMulti.
	collection.Fields.Add(&core.BoolField{
		Name:     "isAsc",
		Required: true,
	})
	collection.Fields.Add(&core.BoolField{
		Name:     "isTrigger",
		Required: true,
	})
	collection.Fields.Add(&core.BoolField{
		Name:     "isMulti",
		Required: true,
	})

	// Composite primary key unique index.
	collection.AddIndex("primary_3", true, "idInteraction, idSuggestionType", "")
	// Additional indexes.
	collection.AddIndex("fk_Sort_Interaction1_idx", false, "idInteraction", "")
	collection.AddIndex("fk_Sort_SuggestionType1_idx", false, "idSuggestionType", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Sort' collection with SQL mapping")
	return nil
}

// SQL: create table SuggestionTypeValues
// (
//     idSuggestionType int not null,
//     value            varchar(150) null,
//     active           tinyint(1) default 1 not null,
//     constraint PRIMARY_id_and_value unique (idSuggestionType, value),
//     constraint idSuggestionType unique (idSuggestionType, value)
// )
// charset = utf8;
// create index fk_idSuggestioType_1827635 on SuggestionTypeValues (idSuggestionType);
func setupSuggestionTypeValuesCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("SuggestionTypeValues"); err == nil {
		log.Println("Collection 'SuggestionTypeValues' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("SuggestionTypeValues")

	// idSuggestionType: reference SuggestionType.
	suggestionTypeCollection, err := app.FindCollectionByNameOrId("SuggestionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestionType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionTypeCollection.Id,
	})

	// value: varchar(150) null.
	collection.Fields.Add(&core.TextField{
		Name:     "value",
		Required: false,
		Max:      150,
	})

	// active: tinyint(1) default 1 not null.
	collection.Fields.Add(&core.BoolField{
		Name:     "active",
		Required: true,
	})

	// Composite unique index on (idSuggestionType, value).
	collection.AddIndex("PRIMARY_id_and_value", true, "idSuggestionType, value", "")
	// Additional index on idSuggestionType.
	collection.AddIndex("fk_idSuggestioType_1827635", false, "idSuggestionType", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'SuggestionTypeValues' collection with SQL mapping")
	return nil
}

// SQL: create table UniqueId
// (
//     idUniqueID int auto_increment primary key,
//     active     tinyint(1) default 1 not null,
//     notes      varchar(500) null
// )
// charset = utf8;
func setupUniqueIdCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("UniqueId"); err == nil {
		log.Println("Collection 'UniqueId' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("UniqueId")

	// active: tinyint(1) default 1 not null.
	collection.Fields.Add(&core.BoolField{
		Name:     "active",
		Required: true,
	})
	// notes: varchar(500) null.
	collection.Fields.Add(&core.TextField{
		Name:     "notes",
		Required: false,
		Max:      500,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'UniqueId' collection with SQL mapping")
	return nil
}

// SQL: create table Comments
// (
//     idComment     int auto_increment primary key,
//     idInteraction int not null,
//     idUniqueID    int not null,
//     comment       longtext not null,
//     voteUp        int default 0 not null,
//     voteDown      int default 0 not null,
//     constraint Comments___fk_idInteraction_ksjdfba87aidsb foreign key (idInteraction) references Interaction (idInteraction),
//     constraint Comments___fk_iduniqid_oq83eyfgqwuyofhba foreign key (idUniqueID) references UniqueId (idUniqueID)
// );
func setupCommentsCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Comments"); err == nil {
		log.Println("Collection 'Comments' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Comments")

	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idUniqueID: reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idUniqueID",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  uniqueIdCollection.Id,
	})
	// comment: longtext not null.
	collection.Fields.Add(&core.TextField{
		Name:     "comment",
		Required: true,
		// No max set to emulate longtext.
	})
	// voteUp: int default 0 not null.
	collection.Fields.Add(&core.NumberField{
		Name:     "voteUp",
		Required: true,
	})
	// voteDown: int default 0 not null.
	collection.Fields.Add(&core.NumberField{
		Name:     "voteDown",
		Required: true,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Comments' collection with SQL mapping")
	return nil
}

// SQL: create table CommentVote
// (
//     idCommentVote int auto_increment primary key,
//     idInteraction int not null,
//     idComment     int not null,
//     vote          varchar(20) not null,
//     selected      tinyint null,
//     constraint table_name___fk_comments_akhsdfashjld foreign key (idComment) references Comments (idComment),
//     constraint table_name___fk_interaction_akdhfa foreign key (idInteraction) references Interaction (idInteraction),
//     constraint voteType check (`vote` in ('voteUp', 'voteUp-deselect', 'voteDown', 'voteDown-desel)
// );
func setupCommentVoteCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("CommentVote"); err == nil {
		log.Println("Collection 'CommentVote' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("CommentVote")

	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idComment: reference Comments.
	commentsCollection, err := app.FindCollectionByNameOrId("Comments")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idComment",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  commentsCollection.Id,
	})
	// vote: varchar(20) not null.
	collection.Fields.Add(&core.TextField{
		Name:     "vote",
		Required: true,
		Max:      20,
	})
	// selected: tinyint null.
	collection.Fields.Add(&core.NumberField{
		Name:     "selected",
		Required: false,
	})
	// Note: The CHECK constraint for "vote" is not supported here.

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'CommentVote' collection with SQL mapping")
	return nil
}

// SQL: create table CommentsView
// (
//     idCommentsView int auto_increment primary key,
//     idUniqueID     int null,
//     idInteraction  int null,
//     constraint CommentsView___fk_comments_alksdhfga1231 foreign key (idUniqueID) references UniqueId (idUniqueID),
//     constraint CommentsView___fk_interaction_aljhsdfg5123 foreign key (idInteraction) references Interaction (idInteraction)
// );
func setupCommentsViewCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("CommentsView"); err == nil {
		log.Println("Collection 'CommentsView' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("CommentsView")

	// idUniqueID: optional, reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err == nil {
		collection.Fields.Add(&core.RelationField{
			Name:          "idUniqueID",
			Required:      false,
			CascadeDelete: false,
			CollectionId:  uniqueIdCollection.Id,
		})
	} else {
		collection.Fields.Add(&core.NumberField{
			Name:     "idUniqueID",
			Required: false,
		})
	}
	// idInteraction: optional, reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err == nil {
		collection.Fields.Add(&core.RelationField{
			Name:          "idInteraction",
			Required:      false,
			CascadeDelete: false,
			CollectionId:  interactionCollection.Id,
		})
	} else {
		collection.Fields.Add(&core.NumberField{
			Name:     "idInteraction",
			Required: false,
		})
	}

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'CommentsView' collection with SQL mapping")
	return nil
}

// SQL: create table Databaits
// (
//     idDatabait             int auto_increment primary key,
//     idInteraction          int not null,
//     idUniqueID             int null,
//     idDatabaitTemplateType int not null,
//     idDatabaitCreateType   int not null,
//     databait               varchar(1500) not null,
//     columns                varchar(1500) null,
//     vals                   varchar(1500) null,
//     notes                  varchar(5000) not null,
//     created                timestamp default current_timestamp() not null,
//     closed                 timestamp default '0000-00-00 00:00:00' not null,
//     nextAction             int null,
//     constraint _fk_databaits_nextAction_b6345das foreign key (nextAction) references DatabaitNextAction (idDatabaitNextAction),
//     constraint _fk_idDatabaitCreateType_b6345das foreign key (idDatabaitCreateType) references DatabaitCreateType (idDatabaitCreateType),
//     constraint _fk_idDatabaitTemplateType_b6345das foreign key (idDatabaitTemplateType) references DatabaitTemplateType (idDatabaitTemplateType),
//     constraint _fk_idInteraction_Databaits_a6a3344 foreign key (idInteraction) references Interaction (idInteraction),
//     constraint _fk_idUniqueID_Databaits_b111edss foreign key (idUniqueID) references UniqueId (idUniqueID)
// )
// charset = utf8;
func setupDatabaitsCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Databaits"); err == nil {
		log.Println("Collection 'Databaits' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Databaits")

	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idUniqueID: optional, reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err == nil {
		collection.Fields.Add(&core.RelationField{
			Name:          "idUniqueID",
			Required:      false,
			CascadeDelete: false,
			CollectionId:  uniqueIdCollection.Id,
		})
	} else {
		collection.Fields.Add(&core.NumberField{
			Name:     "idUniqueID",
			Required: false,
		})
	}
	// idDatabaitTemplateType: reference DatabaitTemplateType.
	databaitTemplateTypeCollection, err := app.FindCollectionByNameOrId("DatabaitTemplateType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idDatabaitTemplateType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  databaitTemplateTypeCollection.Id,
	})
	// idDatabaitCreateType: reference DatabaitCreateType.
	databaitCreateTypeCollection, err := app.FindCollectionByNameOrId("DatabaitCreateType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idDatabaitCreateType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  databaitCreateTypeCollection.Id,
	})
	// databait: varchar(1500) not null.
	collection.Fields.Add(&core.TextField{
		Name:     "databait",
		Required: true,
		Max:      1500,
	})
	// columns: varchar(1500) null.
	collection.Fields.Add(&core.TextField{
		Name:     "columns",
		Required: false,
		Max:      1500,
	})
	// vals: varchar(1500) null.
	collection.Fields.Add(&core.TextField{
		Name:     "vals",
		Required: false,
		Max:      1500,
	})
	// notes: varchar(5000) not null.
	collection.Fields.Add(&core.TextField{
		Name:     "notes",
		Required: true,
		Max:      5000,
	})
	// created: timestamp default current_timestamp() not null.
	collection.Fields.Add(&core.AutodateField{
		Name:        "created",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// closed: timestamp default '0000-00-00 00:00:00' not null.
	collection.Fields.Add(&core.AutodateField{
		Name:        "closed",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// nextAction: optional, reference DatabaitNextAction.
	databaitNextActionCollection, err := app.FindCollectionByNameOrId("DatabaitNextAction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "nextAction",
		Required:      false,
		CascadeDelete: false,
		CollectionId:  databaitNextActionCollection.Id,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Databaits' collection with SQL mapping")
	return nil
}

// SQL: create table DatabaitVisit
// (
//     idInteraction int not null,
//     idDatabait    int not null,
//     source        varchar(200) null,
//     constraint _unique_id_interaction_databaitvisit unique (idInteraction),
//     constraint _fk_idDatabait_DatabaitVisit_b123gda foreign key (idDatabait) references Databaits (idDatabait),
//     constraint _fk_idInteraction_DatabaitVisit_asdhjk16341 foreign key (idInteraction) references Interaction (idInteraction)
// )
// charset = utf8;
func setupDatabaitVisitCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("DatabaitVisit"); err == nil {
		log.Println("Collection 'DatabaitVisit' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("DatabaitVisit")

	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idDatabait: reference Databaits.
	databaitsCollection, err := app.FindCollectionByNameOrId("Databaits")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idDatabait",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  databaitsCollection.Id,
	})
	// source: varchar(200) null.
	collection.Fields.Add(&core.TextField{
		Name:     "source",
		Required: false,
		Max:      200,
	})
	// Unique constraint on idInteraction.
	collection.AddIndex("_unique_id_interaction_databaitvisit", true, "idInteraction", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'DatabaitVisit' collection with SQL mapping")
	return nil
}

// SQL: create table Edit_DelRow
// (
//     idUniqueID int not null,
//     idEdit     int not null,
//     comment    varchar(250) not null,
//     constraint "ALTER TABLE `Edit_DelRow` ADD  CONSTRAINT `_fk_edit_akjdhaas` FO" foreign key (idEdit) references Edit (idEdit),
//     constraint _fk_UniqueID_1027836asda foreign key (idUniqueID) references UniqueId (idUniqueID)
// )
// charset = utf8;
func setupEditDelRowCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Edit_DelRow"); err == nil {
		log.Println("Collection 'Edit_DelRow' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Edit_DelRow")

	// idUniqueID: reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idUniqueID",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  uniqueIdCollection.Id,
	})
	// idEdit: reference Edit.
	editCollection, err := app.FindCollectionByNameOrId("Edit")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idEdit",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  editCollection.Id,
	})
	// comment: varchar(250) not null.
	collection.Fields.Add(&core.TextField{
		Name:     "comment",
		Required: true,
		Max:      250,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Edit_DelRow' collection with SQL mapping")
	return nil
}

// SQL: create table HelpUs
// (
//     idHelpUs      int auto_increment primary key,
//     idInteraction int not null,
//     idUniqueID    int not null,
//     helpUsType    varchar(200) not null,
//     question      varchar(2500) not null,
//     answer        varchar(2500) null,
//     start         timestamp default current_timestamp() not null,
//     answered      timestamp default '0000-00-00 00:00:00' not null,
//     showAnother   timestamp default '0000-00-00 00:00:00' not null,
//     closed        timestamp default '0000-00-00 00:00:00' not null,
//     constraint HelpUs___fk_interaction_idInteraction foreign key (idInteraction) references Interaction (idInteraction),
//     constraint HelpUs___fk_interaction_idUniqueID foreign key (idUniqueID) references UniqueId (idUniqueID)
// );
func setupHelpUsCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("HelpUs"); err == nil {
		log.Println("Collection 'HelpUs' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("HelpUs")

	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idUniqueID: reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idUniqueID",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  uniqueIdCollection.Id,
	})
	// helpUsType: varchar(200) not null.
	collection.Fields.Add(&core.TextField{
		Name:     "helpUsType",
		Required: true,
		Max:      200,
	})
	// question: varchar(2500) not null.
	collection.Fields.Add(&core.TextField{
		Name:     "question",
		Required: true,
		Max:      2500,
	})
	// answer: varchar(2500) null.
	collection.Fields.Add(&core.TextField{
		Name:     "answer",
		Required: false,
		Max:      2500,
	})
	// start: timestamp default current_timestamp() not null.
	collection.Fields.Add(&core.AutodateField{
		Name:        "start",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// answered: timestamp default '0000-00-00 00:00:00' not null.
	collection.Fields.Add(&core.AutodateField{
		Name:        "answered",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// showAnother: timestamp default '0000-00-00 00:00:00' not null.
	collection.Fields.Add(&core.AutodateField{
		Name:        "showAnother",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})
	// closed: timestamp default '0000-00-00 00:00:00' not null.
	collection.Fields.Add(&core.AutodateField{
		Name:        "closed",
		OnCreate:    true,
		OnUpdate:    false,
		Presentable: true,
	})

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'HelpUs' collection with SQL mapping")
	return nil
}

// SQL: create table Suggestions
// (
//     idSuggestion     int auto_increment primary key,
//     idSuggestionType int not null,
//     idUniqueID       int not null,
//     idProfile        int default 2 not null,
//     suggestion       varchar(1500) default '' not null,
//     active           tinyint(1) default 1 not null,
//     confidence       bigint(255) null,
//     last_updated     datetime default current_timestamp() not null on update current_timestamp(),
//     constraint idSuggestion unique (idSuggestion),
//     constraint fk_idSuggestionType_123687 foreign key (idSuggestionType) references SuggestionType (idSuggestionType),
//     constraint fk_unique_row_id foreign key (idUniqueID) references UniqueId (idUniqueID)
// )
// charset = utf8;
// Also, add these indexes:
// create index fk_Suggestion_SuggestionType1_idx on Suggestions (idSuggestionType);
// create index fk_Suggestion_UniqueID_idx on Suggestions (idUniqueID);
// create index idSuggestion_2 on Suggestions (idSuggestion);
func setupSuggestionsCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Suggestions"); err == nil {
		log.Println("Collection 'Suggestions' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Suggestions")

	// idSuggestion is handled internally.
	// idSuggestionType: reference SuggestionType.
	suggestionTypeCollection, err := app.FindCollectionByNameOrId("SuggestionType")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestionType",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionTypeCollection.Id,
	})
	// idUniqueID: reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idUniqueID",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  uniqueIdCollection.Id,
	})
	// idProfile: reference Profile.
	profileCollection, err := app.FindCollectionByNameOrId("Profile")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idProfile",
		Required:      true, // default 2 is assumed to be set by application logic
		CascadeDelete: false,
		CollectionId:  profileCollection.Id,
	})
	// suggestion: text field (max 1500)
	collection.Fields.Add(&core.TextField{
		Name:     "suggestion",
		Required: true,
		Max:      1500,
	})
	// active: boolean field
	collection.Fields.Add(&core.BoolField{
		Name:     "active",
		Required: true,
	})
	// confidence: number field, optional.
	collection.Fields.Add(&core.NumberField{
		Name:     "confidence",
		Required: false,
	})
	// last_updated: autodate field (on create and update)
	collection.Fields.Add(&core.AutodateField{
		Name:        "last_updated",
		OnCreate:    false,
		OnUpdate:    true,
		Presentable: true,
	})

	// Add explicit indexes per SQL:
	collection.AddIndex("fk_Suggestion_SuggestionType1_idx", false, "idSuggestionType", "")
	collection.AddIndex("fk_Suggestion_UniqueID_idx", false, "idUniqueID", "")
	collection.AddIndex("idSuggestion_2", false, "id", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Suggestions' collection with SQL mapping")
	return nil
}

// SQL: create table Copy
// (
//     idInteraction int not null,
//     idSuggestion  int not null,
//     primary key (idInteraction, idSuggestion),
//     constraint _fk_idInteraction_4447654 foreign key (idInteraction) references Interaction (idInteraction),
//     constraint _fk_idSuggestion_4417654 foreign key (idSuggestion) references Suggestions (idSuggestion)
// )
// charset = utf8;
func setupCopyCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Copy"); err == nil {
		log.Println("Collection 'Copy' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Copy")
	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idSuggestion: reference Suggestions.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})
	// Composite primary key as unique index.
	collection.AddIndex("primary_4", true, "idInteraction, idSuggestion", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Copy' collection with SQL mapping")
	return nil
}

// SQL: create table Edit_NewRow
// (
//     idEdit       int not null,
//     idSuggestion int not null,
//     isCorrect    tinyint(1) default 2 not null,
//     constraint idEdit unique (idEdit, idSuggestion),
//     constraint _fk_idEdit_from_Edit_NewRow_asdkl123 foreign key (idEdit) references Edit (idEdit),
//     constraint _fk_idSuggestion_from_Edit_NewRow_asdkl123 foreign key (idSuggestion) references Suggestions (idSuggestion)
// )
// charset = utf8;
// create index _index_edit_suggestion_idEdit_agsdh1872dg on Edit_NewRow (idEdit);
// create index _index_edit_suggestion_idSuggestion_agsdh1872dg on Edit_NewRow (idSuggestion);
func setupEditNewRowCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Edit_NewRow"); err == nil {
		log.Println("Collection 'Edit_NewRow' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Edit_NewRow")
	// idEdit: reference Edit.
	editCollection, err := app.FindCollectionByNameOrId("Edit")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idEdit",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  editCollection.Id,
	})
	// idSuggestion: reference Suggestions.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})
	// isCorrect: number field (default 2)
	collection.Fields.Add(&core.NumberField{
		Name:     "isCorrect",
		Required: true,
	})
	// Unique composite index on (idEdit, idSuggestion)
	collection.AddIndex("idEdit_2", true, "idEdit, idSuggestion", "")
	// Additional indexes:
	collection.AddIndex("_index_edit_suggestion_idEdit_agsdh1872dg_2", false, "idEdit", "")
	collection.AddIndex("_index_edit_suggestion_idSuggestion_agsdh1872dg_2", false, "idSuggestion", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Edit_NewRow' collection with SQL mapping")
	return nil
}

// SQL: create table Paste
// (
//     idInteraction         int not null primary key,
//     pasteValue            varchar(1500) not null,
//     copyCellIdSuggestion  int null,
//     copyCellValue         varchar(1500) null,
//     pasteCellIdSuggestion int not null,
//     pasteCellValue        varchar(1500) not null,
//     constraint _fk_copy_id_suggestion_197823 foreign key (copyCellIdSuggestion) references Suggestions (idSuggestion),
//     constraint _fk_paste_id_suggestion_197823 foreign key (pasteCellIdSuggestion) references Suggestions (idSuggestion)
// )
// charset = utf8;
func setupPasteCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Paste"); err == nil {
		log.Println("Collection 'Paste' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Paste")
	// idInteraction: reference Interaction (primary key)
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// pasteValue: text field (max 1500)
	collection.Fields.Add(&core.TextField{
		Name:     "pasteValue",
		Required: true,
		Max:      1500,
	})
	// copyCellIdSuggestion: optional relation to Suggestions.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "copyCellIdSuggestion",
		Required:      false,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})
	// copyCellValue: optional text field (max 1500)
	collection.Fields.Add(&core.TextField{
		Name:     "copyCellValue",
		Required: false,
		Max:      1500,
	})
	// pasteCellIdSuggestion: relation to Suggestions.
	collection.Fields.Add(&core.RelationField{
		Name:          "pasteCellIdSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})
	// pasteCellValue: text field (max 1500)
	collection.Fields.Add(&core.TextField{
		Name:     "pasteCellValue",
		Required: true,
		Max:      1500,
	})
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Paste' collection with SQL mapping")
	return nil
}

// SQL: create table SearchGoogle
// (
//     IdInteraction int not null,
//     idUniqueID    int not null,
//     idSuggestion  int not null,
//     searchValues  varchar(3000) not null,
//     constraint IdInteraction unique (IdInteraction),
//     constraint _fk_idInteraction_SearchGoogle_a645das foreign key (IdInteraction) references Interaction (idInteraction),
//     constraint _fk_idSuggestion_SearchGoogle_a645das foreign key (idSuggestion) references Suggestions (idSuggestion),
//     constraint _fk_idUniqueID_SearchGoogle_a645das foreign key (idUniqueID) references UniqueId (idUniqueID)
// )
// charset = utf8;
func setupSearchGoogleCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("SearchGoogle"); err == nil {
		log.Println("Collection 'SearchGoogle' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("SearchGoogle")
	// IdInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "IdInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// idUniqueID: reference UniqueId.
	uniqueIdCollection, err := app.FindCollectionByNameOrId("UniqueId")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idUniqueID",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  uniqueIdCollection.Id,
	})
	// idSuggestion: reference Suggestions.
	suggestionsCollection, err := app.FindCollectionByNameOrId("Suggestions")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idSuggestion",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  suggestionsCollection.Id,
	})
	// searchValues: text field (max 3000)
	collection.Fields.Add(&core.TextField{
		Name:     "searchValues",
		Required: true,
		Max:      3000,
	})
	// Unique index on IdInteraction.
	collection.AddIndex("IdInteraction", true, "IdInteraction", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'SearchGoogle' collection with SQL mapping")
	return nil
}

// SQL: create table ViewChange
// (
//     idInteraction int not null,
//     viewname      varchar(50) not null,
//     primary key (idInteraction, viewname)
// )
// charset = utf8;
func setupViewChangeCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("ViewChange"); err == nil {
		log.Println("Collection 'ViewChange' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("ViewChange")
	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// viewname: text field (max 50)
	collection.Fields.Add(&core.TextField{
		Name:     "viewname",
		Required: true,
		Max:      50,
	})
	// Composite primary key as unique index.
	collection.AddIndex("primary_5", true, "idInteraction, viewname", "")

	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'ViewChange' collection with SQL mapping")
	return nil
}

// SQL: create table Visit
// (
//     idVisit       int auto_increment primary key,
//     idInteraction int not null,
//     source        varchar(500) null,
//     searchCol     varchar(50) null,
//     searchVal     varchar(500) null,
//     constraint Visit_idVisit_uindex unique (idVisit),
//     constraint Visit___fk_interaction foreign key (idInteraction) references Interaction (idInteraction)
// );
func setupVisitCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("Visit"); err == nil {
		log.Println("Collection 'Visit' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("Visit")
	// idVisit is handled internally.
	// idInteraction: reference Interaction.
	interactionCollection, err := app.FindCollectionByNameOrId("Interaction")
	if err != nil {
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "idInteraction",
		Required:      true,
		CascadeDelete: false,
		CollectionId:  interactionCollection.Id,
	})
	// source: varchar(500) null.
	collection.Fields.Add(&core.TextField{
		Name:     "source",
		Required: false,
		Max:      500,
	})
	// searchCol: varchar(50) null.
	collection.Fields.Add(&core.TextField{
		Name:     "searchCol",
		Required: false,
		Max:      50,
	})
	// searchVal: varchar(500) null.
	collection.Fields.Add(&core.TextField{
		Name:     "searchVal",
		Required: false,
		Max:      500,
	})
	// Unique constraint on idVisit is already the primary key.
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'Visit' collection with SQL mapping")
	return nil
}

// ------------------------ sessions ------------------------
// SQL: create table sessions
// (
//     session_id varchar(128) collate utf8mb4_bin not null primary key,
//     expires    int(11) unsigned not null,
//     data       text collate utf8mb4_bin null
// )
// charset = utf8;
func setupSessionsCollection(app *pocketbase.PocketBase) error {
	if _, err := app.FindCollectionByNameOrId("sessions"); err == nil {
		log.Println("Collection 'sessions' already exists. Skipping creation.")
		return nil
	}
	collection := core.NewBaseCollection("sessions")
	// session_id: primary key, use TextField.
	collection.Fields.Add(&core.TextField{
		Name:     "session_id",
		Required: true,
		Max:      128,
	})
	// expires: number field.
	collection.Fields.Add(&core.NumberField{
		Name:     "expires",
		Required: true,
	})
	// data: text field.
	collection.Fields.Add(&core.TextField{
		Name:     "data",
		Required: false,
	})
	if err := app.Save(collection); err != nil {
		return err
	}
	log.Println("Created 'sessions' collection with SQL mapping")
	return nil
}
