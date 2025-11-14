import type { ColumnData } from "../interfaces/ColumnData";

// base API call for interactions
export function recordInteraction(
  url: string,
  sessionID: number,
  data: Record<string, any>,
  onSuccess?: (res: Response) => void,
  onError?: (res: Response | Error) => void
) {
  fetch(url, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      IDSession: sessionID,
      ...data,
    }),
  })
    .then(res => {
      if (res.ok) {
        onSuccess?.(res);
      }
      else {
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
  sessionID: number,
  idSuggestion: number,
  rowValues: ColumnData
) {
  const rowValuesString = JSON.stringify(rowValues);

  recordInteraction(
    "/api/clicks",
    sessionID,
    {
      IDInteraction: 1,
      IDSuggestion: idSuggestion,
      RowValues: rowValuesString,
    }
  );
}

// NEED CLARIFICATION

// column search
// export function recordColumnSearch()

// cell edit
// export function recordCellEdit()

// row add
// export function recordRowAdd()

// row delete
// export function recordRowDelete()