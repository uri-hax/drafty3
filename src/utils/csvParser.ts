import Papa, { type ParseResult } from 'papaparse';
import type { GridColumn } from "@glideapps/glide-data-grid";
import type { ColumnData } from '../interfaces/ColumnData';
import { generateOptionsLists, generateColumnWidths } from './constants';
import yaml from 'js-yaml'; 

/*
  Attempt to fetch and parse a YAML schema file that defines the column types.
  If the file is not provided or fails to load, return an empty schema.
*/

const fetchYamlSchema = async (yamlFilePath?: string): Promise<Record<string, string>> => {
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

    const columnSchema: Record<string, string> = {};

    for (const [key, value] of Object.entries(schema)) {
      if (typeof value === 'string') {
        columnSchema[key] = value;
      } 
      else {
        console.warn(`Non-string type for column "${key}" in schema. Defaulting to 'string'.`);
        columnSchema[key] = 'string';
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
  customWidths: Record<string, string> = {},
  csvFilePath: string, 
  yamlSchemaFilePath?: string
): Promise<{ 
  gridColumns: GridColumn[]; 
  parsedData: ColumnData[]; 
  optionsLists: Record<string, string[]>; 
  columnSchema: Record<string, string>; 
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
              const colType = columnSchema[key] || 'string'; 
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
          const columnWidths = generateColumnWidths(columnKeys, customWidths);

          const gridColumns: GridColumn[] = columnKeys.map((key) => ({
            id: key,
            title: key,
            width: parseFloat(columnWidths[key].replace('%', '')) * gridWidth / 100,
          }));

          resolve({ gridColumns, parsedData, optionsLists, columnSchema });
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

/*
  Use the PocketBase client instance (pb) to load and pivot Suggestion data

export const fetchPocketbaseData = async (
  pb: PocketBase,
  gridWidth: number,
  customWidths: Record<string, string> = {},
  yamlSchemaFilePath?: string
): Promise<{
  gridColumns: GridColumn[];
  parsedData: ColumnData[];
  optionsLists: Record<string, string[]>;
  columnSchema: Record<string, string>;
}> => {
  const columnSchema = await fetchYamlSchema(yamlSchemaFilePath);

  const suggestionTypes = await pb.collection("SuggestionType").getFullList({
    filter: "isActive=true",
    sort: "columnOrder",
  });

  const columnsMeta: [string, string][] = suggestionTypes.map((t: any) => [t.id, t.name]);

  const suggestions = await pb.collection("Suggestions").getFullList({
    filter: "active=true",
  });

  const grouped = new Map<string, any[]>();
  for (const s of suggestions) {
    if (!grouped.has(s.idUniqueID)) grouped.set(s.idUniqueID, []);
    grouped.get(s.idUniqueID)!.push(s);
  }

  const parsedData: ColumnData[] = Array.from(grouped.entries()).map(([uniqueId, records]) => {
    const row: ColumnData = { idUniqueID: uniqueId };

    for (const [typeId, name] of columnsMeta) {
      const candidates = records.filter((r) => r.idSuggestionType === typeId);
      let val = "";

      if (candidates.length) {
        candidates.sort((a, b) => b.confidence - a.confidence);
        val = candidates[0].suggestion;
      }

      const colType = columnSchema[name] || "string";
      let processed: string | string[] = val || "";

      if (typeof processed === "string") {
        if (colType === "string[]") {
          processed = processed.includes(",")
            ? processed.split(",").map((s) => s.trim())
            : [processed.trim()];
        } else {
          processed = processed.trim();
        }
      }

      row[name] = processed;
    }

    return row;
  });

  const columnKeys = [...columnsMeta.map(([, name]) => name)];
  const columnWidths = generateColumnWidths(columnKeys, customWidths);
  const gridColumns: GridColumn[] = columnKeys.map((key) => ({
    id: key,
    title: key,
    width: (parseFloat(columnWidths[key].replace("%", "")) * gridWidth) / 100,
  }));

  const optionsLists = generateOptionsLists(parsedData);

  return { gridColumns, parsedData, optionsLists, columnSchema };
};
*/