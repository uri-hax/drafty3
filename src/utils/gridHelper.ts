import type { ColumnData } from '../interfaces/ColumnData';

export const getColumnMatches = (colKey: string, value: string, rows: ColumnData[]): string[] => {
  const search = value.trim().toLowerCase();
  if (!search) {
    return [];
  }

  const matches: string[] = [];

  for (const row of rows) {
    const cellValue = row[colKey];

    if (Array.isArray(cellValue)) {
      for (const v of cellValue) {
        const str = v?.toString() ?? "";
        if (str.toLowerCase().includes(search)) {
          matches.push(str);
        }
      }
    } 
    else if (cellValue != null) {
      const str = cellValue.toString();
      if (str.toLowerCase().includes(search)) {
        matches.push(str);
      }
    }
  }

  return Array.from(new Set(matches));
};

export const valueToString = (v: unknown): string => {
  if (v == null) {
    return "";
  }
  if (Array.isArray(v)) {
    return [...v].sort().join(", ");
  }
  return String(v);
};

export const sortRows = (
  rows: ColumnData[],
  colKey: string,
  dir: "asc" | "desc"
): ColumnData[] => {
  return [...rows].sort((r1, r2) => {
    const a = valueToString(r1[colKey]);
    const b = valueToString(r2[colKey]);

    if (dir === "asc") {
      return a.localeCompare(b);
    } 
    else {
      return b.localeCompare(a);
    }
  });
};

export const handleAlertSnackbarClose = (
  setAlertSnackbarOpen: React.Dispatch<React.SetStateAction<boolean>>
) => {
  setAlertSnackbarOpen(false);
};

export const handleContributionSnackbarClose = (
  setContributionSnackbarOpen: React.Dispatch<React.SetStateAction<boolean>>
) => {
  setContributionSnackbarOpen(false);
};

export const handleHomePage = (base: string) => {
  window.location.href = `${base}`;
};

export const handleData = (base: string, dataset: string) => {
  window.location.href = `${base}${dataset}`;
};

export const handleEditHistory = (base: string, dataset: string) => {
  window.location.href = `${base}${dataset}/history`;
};