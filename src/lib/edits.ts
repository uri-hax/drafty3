import { getAPI } from "./api";

// base API call for edit-related GET requests.
export function getEditInformation<T>(
  url: string,
  onSuccess?: (data: T, res: Response) => void,
  onError?: (err: Response | Error) => void
) {
  fetch(url, {
    method: "GET",
    // sends cookie so backend can read the session
    credentials: "include", 
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then(async (res) => {
      if (!res.ok) {
        onError?.(res);
        return;
      }

      // try to parse JSON
      try {
        const data = (await res.json()) as T;
        onSuccess?.(data, res);
      } 
      catch (e) {
        onSuccess?.(null as unknown as T, res);
      }
    })
    .catch((err) => {
      console.error("Network error fetching edit information:", err);
      onError?.(err);
    });
}

// type for json response from getColumnId
export type GetColumnIdResponse = {
  idSuggestionType: number;
};

// get column id given column name
export function getColumnId(
  name: string,
  onSuccess?: (data: GetColumnIdResponse, res: Response) => void,
  onError?: (err: Response | Error) => void
) {
  getEditInformation<GetColumnIdResponse>(
    `${getAPI()}/suggestiontypes/${name}`,
    onSuccess,
    onError
  );
}

// type for SuggestionTypeValues rows
export type SuggestionTypeValue = {
  IDSuggestionType: number;
  Value: string;
  Active: number;
};

// get active values for a suggestion type id
export function getSuggestionTypeValues(
  id: number,
  onSuccess?: (data: SuggestionTypeValue[], res: Response) => void,
  onError?: (err: Response | Error) => void
) {
  getEditInformation<SuggestionTypeValue[]>(
    `${getAPI()}/suggestiontypevalues/${id}`,
    onSuccess,
    onError
  );
}