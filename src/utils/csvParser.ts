import Papa, { type ParseResult } from 'papaparse';
import type { GridColumn } from "@glideapps/glide-data-grid";
import type { ColumnData } from '../interfaces/ColumnData';
import { generateOptionsLists, generateColumnWidths } from './constants';
import yaml from 'js-yaml'; 
import type { EditType, ColumnConfig } from '../interfaces/ColumnData';

type YamlSchema = Record<string, ColumnConfig>;

/*
  Attempt to fetch and parse a YAML schema file that defines the column types, edit types, and widths.
  If the file is not provided or fails to load, return an empty schema.
*/

const fetchYamlSchema = async (yamlFilePath?: string): Promise<YamlSchema> => {
  if (!yamlFilePath) {
    return {}; 
  }

  try {
    const response = await fetch(yamlFilePath);

    if (!response.ok) {
      console.warn(`Could not load YAML schema from ${yamlFilePath}, defaulting to strings.`);
      return {};
    }

    const yamlText = await response.text();
    const schema = yaml.load(yamlText);

    if (typeof schema !== 'object' || schema === null) {
      console.warn(`Invalid YAML schema format at ${yamlFilePath}. Defaulting to strings.`);
      return {};
    }

    const columnSchema: YamlSchema = {};

    for (const [key, value] of Object.entries(schema)) {
      if (typeof value === 'object' && value !== null && 'type' in value && 'edit' in value) {
        if (!['string', 'string[]'].includes(value.type) && !['free_text', 'dropdown', 'dropdown_free_text', 'multi_select'].includes(value.edit)) {
          console.log(`Invalid type or edit value for column "${key}" in YAML schema.`);
        }
        else {
          columnSchema[key] = value as ColumnConfig;
        }
      } 
      else {
        console.warn(`Defaulting to 'string' and 'free_text' for column "${key}".`);
        columnSchema[key] = {
          type: 'string',
          edit: 'free_text',
        };
      }
    }

    return columnSchema;
  } 
  catch (error) {
    console.warn(`Error loading YAML schema: ${error}. Defaulting to strings.`);
    return {};
  }
};

/*
  Fetch and parse CSV data from the specified file, optionally using a YAML schema to determine column types.

  gridWidth - The width of the grid for calculating column widths.
  csvFilePath - The path to the CSV file.
  customWidths - Optional custom percentages for specific columns.
  yamlSchemaFilePath - Optional path to a YAML file that defines column types.
  
  returns An object containing dynamically generated grid columns, parsed data rows, and options lists.
*/

export const fetchCsvData = async (
  gridWidth: number,
  csvFilePath: string, 
  yamlSchemaFilePath?: string
): Promise<{ 
  gridColumns: GridColumn[]; 
  parsedData: ColumnData[]; 
  optionsLists: Record<string, string[]>; 
  columnSchema: YamlSchema; 
}> => {
  try {
    const columnSchema = await fetchYamlSchema(yamlSchemaFilePath);
    const response = await fetch(csvFilePath);
    const csvData = await response.text();

    return new Promise((resolve, reject) => {
      Papa.parse(csvData, {
        header: true,
        transformHeader: (header) => header.trim(),
        complete: (results: ParseResult<ColumnData>) => {
          /*
          if (!results.data || results.errors.length > 0) {
            reject(new Error('Error parsing CSV data.'));
            return;
          }
          */
          const filteredRows = results.data.filter((row) => 
            Object.values(row).some((value) => value !== null && value !== "")
          );

          const parsedData: ColumnData[] = filteredRows.map((row) => {
            const processedRow: ColumnData = {};

            for (const [key, value] of Object.entries(row)) {
              const colType = columnSchema[key]?.type || 'string';
              let processedValue: string | string[] = value || "";

              if (typeof processedValue === 'string') {
                if (colType === 'string[]') {
                  processedValue = processedValue.includes(",")
                    ? processedValue.split(",").map((s) => s.trim())
                    : [processedValue.trim()];
                } 
                else {
                  processedValue = processedValue.trim();
                }
              }

              processedRow[key] = processedValue;
            }

            return processedRow;
          });

          const optionsLists = generateOptionsLists(parsedData);
          const columnKeys = Object.keys(parsedData[0] || {}).filter((key) => key !== 'idUniqueID');

          for (const key of columnKeys) {
            if (!columnSchema[key]) {
              columnSchema[key] = { type: "string" };
            }
          }

          const customWidths: Record<string, string> = {};
          for (const key of columnKeys) {
            if (columnSchema[key]?.width) {
              if (!/^\d+%$/.test(columnSchema[key].width)) {
                console.log(`Invalid width format for column "${key}".`);
              } 
              else {
                customWidths[key] = columnSchema[key].width;
              }
            }
          }
          const columnWidths = generateColumnWidths(columnKeys, customWidths);

          const gridColumns: GridColumn[] = columnKeys.map((key) => ({
            id: key,
            title: key,
            width: parseFloat(columnWidths[key].replace('%', '')) * gridWidth / 100,
          }));

          resolve({ gridColumns, parsedData, optionsLists, columnSchema});
        },

        skipEmptyLines: true,
      });
    });
  } 
  catch (error) {
    console.error('Error fetching the CSV file:', error);
    return { gridColumns: [], parsedData: [], optionsLists: {}, columnSchema: {} };
  }
};
