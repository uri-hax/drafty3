// src/components/App.tsx

/*
  This component orchestrates the entire application:
  - Fetches CSV data and optional YAML schema using fetchCsvData, producing:
    * parsedData: dynamically typed rows (ColumnData)
    * gridColumns: columns derived from the CSV headers (excluding 'UniqueId')
    * optionsLists: unique option lists for 'string[]' columns
    * columnSchema: defines the type (string or string[]) for each column
  - Handles filtering rows by user-entered text queries in each column
  - Provides add/delete functionality for rows
  - Integrates a multi-select modal (MultiSelectModal) for editing string[] columns
  - Uses DataGridWrapper to render and edit data in a data-agnostic manner
  - Applies a custom column width configuration
  - Uses query parameters to preserve filters

  This file completes the data-agnostic refactor by dynamically handling columns, schema, and data.
*/

import './App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import { Snackbar, Button } from '@mui/material';
import { CompactSelection, type BubbleCell, type EditableGridCell, type GridSelection, type Item, type GridColumn } from "@glideapps/glide-data-grid";
import { fetchCsvData } from '../utils/csvParser';
import type { ColumnData } from '../interfaces/ColumnData';
import useWindowWidth from '../hooks/useWindow';
import Alert from './Alerts';
import FilterBar from './Filter';
import ActionButtons from './ActionButtons';
import DataGridWrapper from './DataGridWrapper';
import MultiSelectModal from './MultiSelectModal';
import AddRowFooter from './AddRowFooter';
import { ensureSession, type BackendSession } from "../lib/sessions";
import { recordCellClick, recordCellEdit, recordColumnSearch, recordRowAdd, recordRowDelete } from "../lib/interactions"

const customWidths: Record<string, string> = {
  FullName: "15%",
  University: "20%",
  JoinYear: "5%",
  SubField: "18%",
  Bachelors: "20%",
  Doctorate: "20%"
};

export default function App() {
  useEffect(() => {
    const portalDiv = document.getElementById('portal');

    if (!portalDiv) {
      const newPortalDiv = document.createElement('div');
      newPortalDiv.id = 'portal';
      document.body.appendChild(newPortalDiv);
    }
  }, []);

  const [session, setSession] = useState<BackendSession | null>(null);

  useEffect(() => {
    ensureSession().then(setSession).catch(console.error);
  }, []);

  const gridWidth = useWindowWidth();

  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [data, setData] = useState<ColumnData[]>([]);
  const [filteredData, setFilteredData] = useState<ColumnData[]>([]);
  const [optionsLists, setOptionsLists] = useState<Record<string, string[]>>({});
  const [columnSchema, setColumnSchema] = useState<Record<string, string>>({});
  const [columnFilters, setColumnFilters] = useState<Record<string, string>>({});
  const [isAddingRow, setIsAddingRow] = useState<boolean>(false);
  const [newRowData, setNewRowData] = useState<ColumnData>({});
  const [isOverlayVisible, setIsOverlayVisible] = useState<boolean>(false);
  const [selectedOptions, setSelectedOptions] = useState<string[]>([]);
  const [editingCell, setEditingCell] = useState<{ row: number; colKey: string } | null>(null);

  const [gridSelection, setGridSelection] = useState<GridSelection>({
    current: undefined,
    rows: CompactSelection.empty(),
    columns: CompactSelection.empty(),
  });

  const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
  // const originBase = new URL(import.meta.env.BASE_URL || '/', window.location.origin);
  // const url = (p: string) => new URL(p.replace(/^\//, ''), originBase).href;

  useEffect(() => {
    const fetchData = async () => {
      try {
        const { gridColumns, parsedData, optionsLists, columnSchema } =
          await fetchCsvData(
            gridWidth,
            customWidths,
            '/drafty3/suggestions.csv',
            '/drafty3/csprofessors.yaml'
          );
  
        console.log("Grid Columns:", gridColumns);
        console.log("Parsed Data:", parsedData);
        console.log("Options Lists:", optionsLists);
        console.log("Column Schema:", columnSchema);
  
        setColumns(gridColumns);
        setData(parsedData);
        setFilteredData(parsedData);
        setOptionsLists(optionsLists);
        setColumnSchema(columnSchema);
  
        const params = new URLSearchParams(window.location.search);
        const initialFilters: Record<string, string> = {};
        const columnKeys = Object.keys(parsedData[0] || {});
  
        columnKeys.forEach((name) => {
          const value = params.get(name);
          if (value) {
            initialFilters[name] = value;
          }
        });
  
        setColumnFilters(initialFilters);

        const urlFilterCount = Object.values(initialFilters).filter(
          (v) => v && v.trim() !== ""
        ).length;

        for (const [colKey, value] of Object.entries(initialFilters)) {
          const trimmed = value.trim();
          if (!trimmed) {
            continue;
          }

          const matchedValues = getColumnMatches(colKey, value, parsedData);
          const lowerSearch = trimmed.toLowerCase();
          const hasExactMatch = matchedValues.some(
            (m) => m.toLowerCase() === lowerSearch
          );
          const isPartial = !hasExactMatch;
          const isMulti = urlFilterCount > 1;
          const isFromURL = true;

          recordColumnSearch(
            value,
            matchedValues,
            isPartial,
            isMulti,
            isFromURL
          );
        }

        const initialNewRowData: ColumnData = {};
        for (const key of columnKeys) {
          initialNewRowData[key] =
            columnSchema[key] === "string[]" ? [] : "";
        }
  
        setNewRowData(initialNewRowData);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };
  
    fetchData();
  }, [gridWidth]);
 
  useEffect(() => {
    if (columns.length === 0 || data.length === 0) return;

    const filtered = data.filter((row) =>
      columns.every((col) => {
        const colKey = col.id as string;
        const filterValue = columnFilters[colKey];

        if (!filterValue) return true;

        const cellValue = row[colKey];

        if (Array.isArray(cellValue)) {
          return cellValue.some((val) =>
            val.toString().toLowerCase().includes(filterValue.toLowerCase())
          );
        }

        return cellValue?.toString().toLowerCase().includes(filterValue.toLowerCase());
      })
    );

    setFilteredData(filtered);
  }, [columnFilters, data, columns]);

  useEffect(() => {
    if (columns.length === 0) {
      return;
    }

    const params = new URLSearchParams();

    for (const [key, value] of Object.entries(columnFilters)) {
      if (value) {
        params.set(key, value);
      }
    }

    window.history.replaceState(null, '', '?' + params.toString());
  }, [columnFilters, columns]);

  const getColumnMatches = (colKey: string, value: string, rows: ColumnData[]): string[] => {
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

  const handleColumnFilterChange = (colKey: string, value: string) => {
    const nextFilters: Record<string, string> = {
      ...columnFilters,
      [colKey]: value,
    };

    setColumnFilters(nextFilters);

    const trimmed = value.trim();
    if (!trimmed) {
      return;
    }

    const nonEmptyFilterCount = Object.values(nextFilters).filter(
      (v) => v && v.trim() !== ""
    ).length;
    const isMulti = nonEmptyFilterCount > 1;

    const matchedValues = getColumnMatches(colKey, value, data);

    const lowerSearch = trimmed.toLowerCase();
    const hasExactMatch = matchedValues.some(
      (m) => m.toLowerCase() === lowerSearch
    );
    const isPartial = !hasExactMatch;

    const isFromURL = false; 

    recordColumnSearch(
      value,
      matchedValues,
      isPartial,
      isMulti,
      isFromURL
    );
  };

  const allFieldsFilled = React.useMemo(() => {
    if (!Object.keys(columnSchema).length || !Object.keys(newRowData).length) return false;

    return Object.keys(columnSchema).every((key) => {
      const colType = columnSchema[key];
      const val = newRowData[key];

      if (colType === 'string[]') {
        return Array.isArray(val) && val.length > 0;
      } 
      else {
        return typeof val === 'string' && val.trim() !== '';
      }
    });
  }, [columnSchema, newRowData]);

  const onCellActivated = (cell: Item) => {
    if (!session) {
      return;
    }

    const [col, row] = cell;
    console.log("Column: ", col, " Row: ", row);

    const actualRowIndex = data.indexOf(filteredData[row]);
    const rowData = data[actualRowIndex];
    const idSuggestion = actualRowIndex;
    recordCellClick(idSuggestion, rowData);

    const colKey = columns[col].id as string;
    const colType = columnSchema[colKey] || "string";
    console.log(cell)

    if (colType === 'string[]') {
      setEditingCell({ row, colKey });
      const actualRowIndex = data.indexOf(filteredData[row]);
      console.log("Column: ", col, " Actual Row: ", actualRowIndex);
      const cellData = data[actualRowIndex][colKey];
      console.log(cellData);
      setSelectedOptions(Array.isArray(cellData) ? cellData : []);
      setIsOverlayVisible(true);
    }
  };

  const handleSaveOptions = () => {
    if (editingCell) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = [...selectedOptions];
      setData(updatedData);
      setFilteredData(updatedData);
      setIsOverlayVisible(false);
      setEditingCell(null);
    }
  };

  const onCellEdited = (cell: Item, newValue: EditableGridCell | BubbleCell) => {
    const [col, row] = cell;
    const colKey = columns[col].id as string;
    const colType = columnSchema[colKey] || 'string';
    const updatedData = [...data];
    const actualRowIndex = data.indexOf(filteredData[row]);

    if (colType === 'string[]') {
      if (newValue.kind === 'bubble') {
        updatedData[actualRowIndex][colKey] = newValue.data as string[];
      }
    } 
    else {
      if (newValue.kind === 'text') {
        updatedData[actualRowIndex][colKey] = newValue.data as string;
      }
    }

    setData(updatedData);
    setFilteredData(updatedData);

    recordCellEdit();
  };

  const onGridSelectionChange = (newSelection: GridSelection) => {
    setGridSelection(newSelection);
  };

  const handleDeleteRow = () => {
    if (gridSelection.rows.length > 0) {
      const selectedRowIndices = gridSelection.rows.toArray();
      const selectedRows = selectedRowIndices.map((index) => filteredData[index]);

      setData((prev) => prev.filter((row) => !selectedRows.includes(row)));
      setFilteredData((prev) => prev.filter((row) => !selectedRows.includes(row)));
      setGridSelection({
        current: undefined,
        rows: CompactSelection.empty(),
        columns: CompactSelection.empty(),
      });

      recordRowDelete();
    } 
    else if (gridSelection.current && gridSelection.current.cell) {
      const [, row] = gridSelection.current.cell;
      const rowData = filteredData[row];

      setData((prev) => prev.filter((r) => r !== rowData));
      setFilteredData((prev) => prev.filter((r) => r !== rowData));
      setGridSelection({
        current: undefined,
        rows: CompactSelection.empty(),
        columns: CompactSelection.empty(),
      });

      recordRowDelete();
    } 
    else {
      setSnackbarOpen(true);
    }
  };

  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  const handleAddRowConfirm = () => {
    if (!allFieldsFilled) return;

    const updatedData = [...data, newRowData];

    const idSuggestion = updatedData.length - 1;
    recordRowAdd(idSuggestion);

    setData(updatedData);
    setFilteredData(updatedData);
    const resetObj: ColumnData = {};

    for (const key of Object.keys(columnSchema)) {
      resetObj[key] = columnSchema[key] === 'string[]' ? [] : '';
    }

    setNewRowData(resetObj);
    setIsAddingRow(false);
  };

  const handleData = () => {
    window.location.href = "/drafty3/csprofs";
  }

  const handleEditHistory = () => {
    window.location.href = "/drafty3/edit-history";
  };

  return (
    <div className="App" style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      <ActionButtons
        handleData={handleData}
        handleEditHistory={handleEditHistory}
        setIsAddingRow={setIsAddingRow}
        handleDeleteRow={handleDeleteRow}
      />

      {columns.length > 0 ? (
        <FilterBar
          columns={columns}
          columnFilters={columnFilters}
          handleColumnFilterChange={handleColumnFilterChange}
          columnWidths={customWidths}
        />
      ) : (
        <div></div>
      )}

      <div className="grid-container" style={{ flexGrow: 1 }}>
        <DataGridWrapper
          columns={columns}
          filteredData={filteredData}
          onCellEdited={onCellEdited}
          onCellActivated={onCellActivated}
          gridSelection={gridSelection}
          onGridSelectionChange={onGridSelectionChange}
          gridWidth={gridWidth}
          columnSchema={columnSchema}
        />
      </div>

      <MultiSelectModal
        isOverlayVisible={isOverlayVisible}
        setIsOverlayVisible={setIsOverlayVisible}
        selectedOptions={selectedOptions}
        setSelectedOptions={setSelectedOptions}
        handleSaveOptions={handleSaveOptions}
        optionsList={editingCell ? (optionsLists[editingCell.colKey] || []) : []}
        title="Edit Column Values"
      />

      <Snackbar open={snackbarOpen} autoHideDuration={3000} onClose={handleSnackbarClose}>
        <Alert onClose={handleSnackbarClose} severity="warning">
          Please select a cell or row first.
        </Alert>
      </Snackbar>

      {isAddingRow && (
        <AddRowFooter
          columnKeys={Object.keys(columnSchema)}
          newRowData={newRowData}
          setNewRowData={setNewRowData}
          optionsLists={optionsLists}
          handleAddRowConfirm={handleAddRowConfirm}
          setIsAddingRow={setIsAddingRow}
          allFieldsFilled={allFieldsFilled}
          columnSchema={columnSchema}
        />
      )}
    </div>
  );
}
