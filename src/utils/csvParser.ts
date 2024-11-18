// utils/csvParser.ts
import Papa, { type ParseResult } from 'papaparse';
import type { GridColumn } from "@glideapps/glide-data-grid";
import type { Professor } from '../interfaces/Professor';
import { validKeys } from '../interfaces/Professor';
import { columnWidths } from './constants';

export const fetchCsvData = async (gridWidth: number): Promise<{ gridColumns: GridColumn[]; parsedData: Professor[] }> => {
  try {
    const response = await fetch('/csprofessors.csv');
    const csvData = await response.text();

    // Parse the CSV data and set the columns and data for the grid
    return new Promise<{ gridColumns: GridColumn[]; parsedData: Professor[] }>((resolve) => {
      Papa.parse(csvData, {
        header: true,
        transformHeader: (header) => header.trim(),
        complete: (results: ParseResult<{ [key: string]: string }>) => {
          const parsedData = results.data
            .filter((row) => Object.values(row).some((value) => value !== null && value !== ""))
            .map((row) => {
              const professor: Professor = {
                FullName: row["FullName"] || "",
                University: row["University"] || "",
                JoinYear: row["JoinYear"] || "",
                SubField: row["SubField"] ? row["SubField"].split(',').map((s) => s.trim()) : [],
                Bachelors: row["Bachelors"] || "",
                Doctorate: row["Doctorate"] || "",
              };
              return professor;
            });

          // Create grid columns using validKeys
          const gridColumns: GridColumn[] = validKeys.map((key) => {
            let width = 150;
            const colWidth = columnWidths[key];
            if (typeof colWidth === 'string' && colWidth.endsWith('%')) {
              const percent = parseFloat(colWidth) / 100;
              width = gridWidth * percent;
            }
            return {
              id: key,
              title: key,
              width: width,
            };
          });

          resolve({ gridColumns, parsedData });
        },
        skipEmptyLines: true,
      });
    });
  } catch (error) {
    console.error('Error fetching the CSV file:', error);
    return { gridColumns: [], parsedData: [] };
  }
};