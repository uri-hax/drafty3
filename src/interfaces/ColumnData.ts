// src/interfaces/ColumnData.ts

/*
  Represents a single row of data from any CSV file.
  The keys are the column names, and the values are either strings or arrays of strings.
*/

export interface ColumnData {
  [key: string]: string | string[];
}