// src/interfaces/ColumnData.ts

/*
  Define interaces and types related to column data, including the structure of the data and the configuration for each column.
*/

export interface ColumnData {
  [key: string]: string | string[];
}

export type ColumnType = "string" | "string[]";
export type EditType = "free_text" | "dropdown" | "dropdown_free_text" | "multi_select";

export type ColumnConfig = {
  type: ColumnType;
  edit?: EditType;
  width?: string; 
};