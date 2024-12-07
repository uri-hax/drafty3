// src/utils/constants.ts
import type { ColumnData } from '../interfaces/ColumnData';

/*
  Generates a dynamic options list for each string[] column in the dataset.
  Each list will contain unique values for its column, but values may repeat across columns.
  parsedData - The parsed CSV data.
  returns An object where keys are column names and values are unique options for string[] columns.
*/

export const generateOptionsLists = (parsedData: ColumnData[]): Record<string, string[]> => {
  const optionsLists: Record<string, string[]> = {};

  parsedData.forEach((row) => {
    for (const [key, value] of Object.entries(row)) {
      if (Array.isArray(value)) {
        if (!optionsLists[key]) {
          optionsLists[key] = [];
        }
        // Add all values from the current row to the list
        optionsLists[key].push(...value);
      }
    }
  });

  // Remove duplicates within each column
  for (const key in optionsLists) {
    optionsLists[key] = Array.from(new Set(optionsLists[key]));
  }

  return optionsLists;
};

/*
  Generates column widths based on the number of columns and optional custom widths.
  columnKeys - Array of column names.
  customWidths - Optional custom percentages for specific columns.
  returns An object where keys are column names and values are column widths in percentages.
*/

export const generateColumnWidths = (
  columnKeys: string[],
  customWidths: Record<string, string> = {}
): Record<string, string> => {
  const defaultPercent = Math.floor(100 / columnKeys.length); 
  const columnWidths: Record<string, string> = {};

  columnKeys.forEach((key) => {
    columnWidths[key] = customWidths[key] ?? `${defaultPercent}%`; 
  });

  return columnWidths;
};