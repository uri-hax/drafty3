package main

import (
    "log"

    "github.com/pocketbase/pocketbase"
)

func main() {
    app := pocketbase.New()

    if err := app.Bootstrap(); err != nil {
        log.Fatal(err)
    }

    // Independent tables (no dependencies)
    if err := setupDataTypeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupDatabaitCreateTypeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupDatabaitNextActionCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupDatabaitTemplateTypeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupInteractionTypeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupRoleCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSearchTypeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupUniqueIdCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupEntryTypeCollection(app); err != nil {
        log.Fatal(err)
    }

    // Tables that depend on the above
    if err := setupProfileCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSessionCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSuggestionTypeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSuggestionTypeValuesCollection(app); err != nil {
        log.Fatal(err)
    }
    // Interaction is referenced by many tables so create it now.
    if err := setupInteractionCollection(app); err != nil {
        log.Fatal(err)
    }
    // Now create tables that reference Interaction.
    if err := setupCopyColumnCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSuggestionsCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupCopyCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupDoubleClickCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupEditCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupEditSuggestionCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupEditNewRowCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupAliasCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupClickCollection(app); err != nil {
        log.Fatal(err)
    }
    
    // Databait-related tables
    if err := setupDatabaitsCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupDatabaitTweetCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupDatabaitVisitCollection(app); err != nil {
        log.Fatal(err)
    }

    // Other tables that reference Interaction, Suggestions, etc.
    if err := setupSelectRangeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSearchCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSearchMultiCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSortCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupCommentsCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupCommentVoteCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupCommentsViewCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupRemoveUserDataCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupHelpUsCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupPasteCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSearchGoogleCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupViewChangeCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupVisitCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupSessionsCollection(app); err != nil {
        log.Fatal(err)
    }
    if err := setupEditDelRowCollection(app); err != nil {
        log.Fatal(err)
    }

    //csv imports
    // if err := suggestionsCSVImport(app); err != nil {
    //     log.Fatal(err)
    // }

    // if err := uniqueIdsCSVImport(app); err != nil {
    //     log.Fatal(err)
    // }

    // if err := suggestionTypeCSVImport(app); err != nil {
    //     log.Fatal(err)
    // }

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}