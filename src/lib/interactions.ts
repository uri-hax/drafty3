import type { ColumnData } from "../interfaces/ColumnData";
import { getAPI } from "./api";

// base API call for interactions
export function recordInteraction(
  url: string,
  data: Record<string, any>,
  onSuccess?: (res: Response) => void,
  onError?: (res: Response | Error) => void
) {
  fetch(url, {
    method: "POST",
    // sends the cookie so backend can read the session
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then(res => {
      if (res.ok) {
        onSuccess?.(res);
      } else {
        onError?.(res);
      }
    })
    .catch(err => {
      console.error("Network error recording interaction:", err);
      onError?.(err);
    });
}

// cell click
export function recordCellClick(
  idSuggestionType: number,
  idUniqueID: number,
  rowValues: ColumnData
) {
  const rowValuesString = JSON.stringify(rowValues);

  recordInteraction(
    `${getAPI()}/clicks`,
    {
      IDInteractionType: 1, // placeholder
      IDSuggestionType: idSuggestionType,
      IDUniqueID:       idUniqueID,
      RowValues:        rowValuesString,
    }
  );
}

// cell edit
export function recordCellEdit(
  idSuggestionType: number,
  idUniqueID: number,
  suggestion: string,
) {
  recordInteraction(
    `${getAPI()}/edits`,
    {
      IDInteractionType: 2, // placeholder
      IDEntryType: 1, // placeholder
      Mode: "normal", // default
      IsCorrect: 2, // default
      IDSuggestionType: idSuggestionType,
      IDUniqueID: idUniqueID,
      Suggestion: suggestion,
      Active: 1, // default
    }
  );
}

// column search
export function recordColumnSearch(
  value: string,
  matchedValues: any,
  isPartial: boolean,
  isMulti: boolean,
  isFromURL: boolean
) {
  recordInteraction(
    `${getAPI()}/searches`,
    {
      IDInteractionType: 2, // placeholder
      IDSuggestionType: 1, // placeholder
      IDSearchType: 1, // placeholder

      IsPartial:  isPartial ? 1 : 0,
      IsMulti:    isMulti ? 1 : 0,
      IsFromURL:  isFromURL ? 1 : 0,

      Value: value,
      MatchedValues: JSON.stringify(matchedValues),
    }
  );
}

// row add
export function recordRowAdd(
  idUniqueID: number,
  cells: {
    IDSuggestionType: number;
    Suggestion: string;
    Active: number; 
    Confidence: number; 
  }[]
) {
  recordInteraction(`${getAPI()}/editnewrows`, {
    IDInteractionType: 3, // placeholder
    IDEntryType: 2, // placeholder
    IDUniqueID: idUniqueID,
    Mode: "normal", // default
    IsCorrect: 2, // default
    Cells: cells,
  });
}

// row delete
export function recordRowDelete(
  idUniqueID: number,
  comment: string
) {
  recordInteraction(
    `${getAPI()}/editdelrows`,
    {
      IDInteractionType: 4, // placeholder
      IDEntryType: 3, // placeholder
      IDUniqueID: idUniqueID,
      Comment: comment,
      Mode: "normal", // default
      IsCorrect: 2, // default
    }
  );
}