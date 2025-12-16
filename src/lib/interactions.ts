import type { ColumnData } from "../interfaces/ColumnData";

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
  idSuggestion: number,
  rowValues: ColumnData
) {
  const rowValuesString = JSON.stringify(rowValues);

  recordInteraction(
    "/api/clicks",
    {
      IDInteractionType: 1, // placeholder
      IDSuggestion:      idSuggestion,
      RowValues:         rowValuesString,
    }
  );
}

// cell edit
export function recordCellEdit(
) {
  recordInteraction(
    "/api/edits",
    {
      IDInteractionType: 2, // placeholder
      IDEntryType:       1, // placeholder
      // Mode and IsCorrect will use defaults for now
      Mode: "normal",
      IsCorrect: 2,
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
    "/api/searches",
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
  idSuggestion: number
) {
  recordInteraction(
    "/api/editnewrows",
    {
      IDInteractionType: 3, // placeholder
      IDEntryType: 2,       // placeholder
      IDSuggestion: idSuggestion,
      // Mode and IsCorrect will use defaults for now
      Mode: "normal",
      IsCorrect: 2,
    }
  );
}

// row delete
export function recordRowDelete() {
  recordInteraction(
    "/api/editdelrows",
    {
      IDInteractionType: 4, // placeholder
      IDEntryType: 3,       // placeholder
      Comment: "",          // placeholder
      // Mode and IsCorrect will use defaults for now
      Mode: "normal",
      IsCorrect: 2,
    }
  );
}