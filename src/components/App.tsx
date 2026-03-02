/*
  This component orchestrates the entire application:
  - Fetches CSV data and optional YAML schema using fetchCsvData, producing:
    * parsedData: dynamically typed rows (ColumnData)
    * gridColumns: columns derived from the CSV headers (excluding 'UniqueId')
    * optionsLists: unique option lists for 'string[]' columns
    * columnSchema: defines the type (string or string[]) and width and edit type for each column
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
import { Snackbar } from '@mui/material';
import { CompactSelection, type BubbleCell, type EditableGridCell, type GridSelection, type Item, type GridColumn } from "@glideapps/glide-data-grid";
import { fetchCsvData } from '../utils/csvParser';
import type { ColumnData } from '../interfaces/ColumnData';
import useWindowWidth from '../hooks/useWindow';
import Alert from './Alerts';
import FilterBar from './Filter';
import ActionButtons from './ActionButtons';
import DataGridWrapper from './DataGridWrapper';
import MultiSelectModal from './MultiSelectModal';
import TextInputModal from "./FreeTextModal";
import DropdownModal from "./DropdownModal";
import DropdownFreeTextModal from "./DropdownFreeTextModal";
import AddRowFooter from './AddRowFooter';
import DeleteRowFooter from './DeleteRowFooter';
import { ensureSession, type BackendSession } from "../lib/sessions";
import { recordCellClick, recordCellEdit, recordColumnSearch, recordRowAdd, recordRowDelete } from "../lib/interactions"
import type { ColumnConfig } from '../interfaces/ColumnData';
import { getColumnId, getSuggestionTypeValues } from "../lib/edits";

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
  const [columnSchema, setColumnSchema] = useState<Record<string, ColumnConfig>>({});
  const [columnFilters, setColumnFilters] = useState<Record<string, string>>({});
  const [isAddingRow, setIsAddingRow] = useState<boolean>(false);
  const [newRowData, setNewRowData] = useState<ColumnData>({});
  const [isMultiOverlayVisible, setIsMultiOverlayVisible] = useState<boolean>(false);
  const [selectedMultiOptions, setSelectedMultiOptions] = useState<string[]>([]);
  const [editingCell, setEditingCell] = useState<{ row: number; colKey: string } | null>(null);

  const [gridSelection, setGridSelection] = useState<GridSelection>({
    current: undefined,
    rows: CompactSelection.empty(),
    columns: CompactSelection.empty(),
  });

  type SortDir = "asc" | "desc";

  const [sortColKey, setSortColKey] = useState<string | null>(null);
  const [sortDir, setSortDir] = useState<SortDir>("asc");

  const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
  // const originBase = new URL(import.meta.env.BASE_URL || '/', window.location.origin);
  // const url = (p: string) => new URL(p.replace(/^\//, ''), originBase).href;

  const [isDeletingRow, setIsDeletingRow] = useState<boolean>(false);
  const [deleteComment, setDeleteComment] = useState<string>("");
  const [pendingDeleteRows, setPendingDeleteRows] = useState<ColumnData[]>([]);

  const [isTextOverlayVisible, setIsTextOverlayVisible] = useState<boolean>(false);
  const [textDraft, setTextDraft] = useState<string>("");

  const [isDropdownOverlayVisible, setIsDropdownOverlayVisible] = useState<boolean>(false);
  const [isDropdownFreeTextOverlayVisible, setIsDropdownFreeTextOverlayVisible] = useState<boolean>(false);

  const [freeTextAltDraft, setFreeTextAltDraft] = useState<string>("");

  const [suggestionTypeIds, setSuggestionTypeIds] = useState<Record<string, number>>({});
  const [suggestionTypeValues, setSuggestionTypeValues] = useState<Record<string, string[]>>({});

  const [activeOptionsList, setActiveOptionsList] = useState<string[]>([]);

  const datasetFiles: Record<string, { csv: string; yaml: string; }> = {
    csprofs: {
      csv: "suggestions.csv",
      yaml: "csprofessors.yaml",
    },
    students: {
      csv: "students.csv",
      yaml: "students.yaml",
    },
  };

  const base = import.meta.env.BASE_URL;
  const dataset = window.location.pathname
    .replace(base, "")
    .split("/")
    .filter(Boolean)[0];

  useEffect(() => {
    const fetchData = async () => {
      try {
        const { gridColumns, parsedData, optionsLists, columnSchema } =
          await fetchCsvData(
            gridWidth,
            `${base}${datasetFiles[dataset].csv}`,
            `${base}${datasetFiles[dataset].yaml}`
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

        const schemaColumnNames = Object.keys(columnSchema);

        const idPairs = await Promise.all(
          schemaColumnNames.map(
            (name) =>
              new Promise<[string, number] | null>((resolve) => {
                getColumnId(
                  name,
                  (res) => resolve([name, res.idSuggestionType]),
                  () => resolve(null) 
                );
              })
          )
        );

        const idMap: Record<string, number> = {};
        for (const pair of idPairs) {
          if (!pair) {
            continue;
          }
          const [name, id] = pair;
          idMap[name] = id;
        }

        setSuggestionTypeIds(idMap);

        console.log("SuggestionType IDs:", idMap);
  
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
            columnSchema[key]?.type === "string[]" ? [] : "";
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

    const finalRows = sortColKey ? sortRows(filtered, sortColKey, sortDir) : filtered;

    setFilteredData(finalRows);
  }, [columnFilters, data, columns, sortColKey, sortDir]);

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

  const customWidths: Record<string, string> = {};
    for (const key of Object.keys(columnSchema)) {
      if (columnSchema[key]?.width) {
        customWidths[key] = columnSchema[key].width;
      }
    }

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
      const colType = columnSchema[key].type;
      const val = newRowData[key];

      if (colType === 'string[]') {
        return Array.isArray(val) && val.length > 0;
      } 
      else {
        return typeof val === 'string' && val.trim() !== '';
      }
    });
  }, [columnSchema, newRowData]);

  const isDeleteCommentFilled = React.useMemo(() => {
    return deleteComment.trim() !== "";
  }, [deleteComment]);

  const getOptionsForCol = (colKey: string) =>
    suggestionTypeValues[colKey] ?? optionsLists[colKey] ?? [];

  const ensureSuggestionValuesLoaded = (colKey: string) => {
    if (suggestionTypeValues[colKey]) {
      return;
    }

    const id = suggestionTypeIds[colKey];
    if (!id) {
      return;
    }

    getSuggestionTypeValues(
      id,
      (rows) => {
        const values = rows.map((r) => r.Value);
        setSuggestionTypeValues((prev) => ({
          ...prev,
          [colKey]: values,
        }));

        setActiveOptionsList(values);
      },
      (err) => {
        console.error(`Failed to load SuggestionTypeValues for ${colKey} (id=${id})`, err);
        setSuggestionTypeValues((prev) => ({
          ...prev,
          [colKey]: [],
        }));

        setActiveOptionsList([]);
      }
    );
  };

  const onCellEditorActivated = (cell: Item) => {
    if (!session) {
      return;
    }

    const [col, row] = cell;
    const rowObj = filteredData[row];
    const rowId = Number(rowObj["idUniqueID"]);
    // id is now row index * num columns + col index
    const cellId = rowId * columns.length + col;
    recordCellClick(cellId, rowObj);

    setIsMultiOverlayVisible(false);
    setIsTextOverlayVisible(false);
    setIsDropdownOverlayVisible(false);
    setIsDropdownFreeTextOverlayVisible(false);

    const colKey = columns[col].id as string;
    setActiveOptionsList(getOptionsForCol(colKey));
    const config = columnSchema[colKey];
    const colType = config?.type ?? "string";
    const editType =
      config?.edit ?? (colType === "string[]" ? "multi_select" : "free_text");

    setEditingCell({ row, colKey });
    if (
      editType === "multi_select" ||
      editType === "dropdown" ||
      editType === "dropdown_free_text"
    ) {
      ensureSuggestionValuesLoaded(colKey);
    }
    const cellVal = rowObj[colKey];

    // multi select modal
    if (colType === "string[]" && editType === "multi_select") {
      setSelectedMultiOptions(Array.isArray(cellVal) ? cellVal : []);
      setIsMultiOverlayVisible(true);
      return;
    }

    // free text modal
    if (colType === "string" && editType === "free_text") {
      setTextDraft(cellVal != null ? String(cellVal) : "");
      setIsTextOverlayVisible(true);
      return;
    }

    // dropdown modal
    if (colType === "string" && editType === "dropdown") {
      setTextDraft(cellVal != null ? String(cellVal) : "");
      setIsDropdownOverlayVisible(true);
      return;
    }

    // dropdown + free text modal
    if (colType === "string" && editType === "dropdown_free_text") {
      setTextDraft(cellVal != null ? String(cellVal) : ""); 
      setFreeTextAltDraft(""); 
      setIsDropdownFreeTextOverlayVisible(true);
      return;
    }

    // fallback -> free text
    setTextDraft(cellVal != null ? String(cellVal) : "");
    setIsTextOverlayVisible(true);
  };

  const handleSaveMulti = () => {
    if (editingCell) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = [...selectedMultiOptions];
      setData(updatedData);
      setFilteredData(updatedData);
      setIsMultiOverlayVisible(false);
      setEditingCell(null);

      recordCellEdit();
    }
  };

  const handleSaveText = () => {
    if (editingCell) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = textDraft
      setData(updatedData);
      setFilteredData(updatedData);
      setIsTextOverlayVisible(false);
      setEditingCell(null);

      recordCellEdit();
    }
  };

  const handleSaveDropdown = () => {
    if (editingCell) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = textDraft;
      setData(updatedData);
      setFilteredData(updatedData);
      setIsDropdownOverlayVisible(false);
      setEditingCell(null);

      recordCellEdit();
    }
  };

  const handleSaveDropdownFreeText = () => {
    if (editingCell) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);

      const finalValue =
        freeTextAltDraft.trim() !== "" ? freeTextAltDraft.trim() : textDraft;

      updatedData[actualRowIndex][colKey] = finalValue;
      setData(updatedData);
      setFilteredData(updatedData);
      setIsDropdownFreeTextOverlayVisible(false);
      setEditingCell(null);

      recordCellEdit();
    }
  };

  const onCellEdited = (cell: Item, newValue: EditableGridCell | BubbleCell) => {
    const [col, row] = cell;
    const colKey = columns[col].id as string;
    const colType = columnSchema[colKey].type || 'string';
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
    const cell = newSelection.current?.cell;

    if (cell) {
      const [col, row] = cell;

      const rowObj = filteredData[row];
      if (rowObj) {
        const rowId = Number(rowObj["idUniqueID"]);
        // id is now row index * num columns + col index
        const cellId = rowId * columns.length + col;
        recordCellClick(cellId, rowObj);
      }
    }

    setGridSelection(newSelection);
  };

  const handleDeleteRow = () => {
    let selectedRows: ColumnData[] = [];

    if (gridSelection.rows.length > 0) {
      const selectedRowIndices = gridSelection.rows.toArray();
      selectedRows = selectedRowIndices
        .map((index) => filteredData[index])
        .filter(Boolean);
    } 
    else if (gridSelection.current?.cell) {
      const [, row] = gridSelection.current.cell;
      const rowData = filteredData[row];
      if (rowData) {
        selectedRows = [rowData];
      }
    }

    if (selectedRows.length === 0) {
      setSnackbarOpen(true);
      setIsDeletingRow(false);
      return;
    }

    setPendingDeleteRows(selectedRows);
    setDeleteComment("");
    setIsDeletingRow(true);
    setIsAddingRow(false);
  };

  const handleDeleteRowConfirm = () => {
    if (!isDeleteCommentFilled) {
      return;
    }
    if (pendingDeleteRows.length === 0) {
      setIsDeletingRow(false);
      setDeleteComment("");
      return;
    }

    setData((prev) => prev.filter((row) => !pendingDeleteRows.includes(row)));
    setFilteredData((prev) => prev.filter((row) => !pendingDeleteRows.includes(row)));

    setGridSelection({
      current: undefined,
      rows: CompactSelection.empty(),
      columns: CompactSelection.empty(),
    });

    recordRowDelete(deleteComment);

    setIsDeletingRow(false);
    setDeleteComment("");
    setPendingDeleteRows([]);
  };

  const handleDeleteRowCancel = () => {
    setIsDeletingRow(false);
    setDeleteComment("");
    setPendingDeleteRows([]);
  };

  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  const handleAddRowConfirm = () => {
    if (!allFieldsFilled) {
      return;
    }

    // fix works assuming idUniqueID increments by 1
    const lastRow = data[data.length - 1];
    const newRowId = Number(lastRow["idUniqueID"]) + 1;

    const rowToAdd: ColumnData = {
      ...newRowData,
      idUniqueID: String(newRowId),
    };

    // id is now row index * num columns + col index
    const cellId = newRowId * columns.length + 0;
    recordRowAdd(cellId);

    setData((prev) => [...prev, rowToAdd]);

    const resetObj: ColumnData = {};
    for (const key of Object.keys(columnSchema)) {
      resetObj[key] = columnSchema[key].type === "string[]" ? [] : "";
    }

    setNewRowData(resetObj);
    setIsAddingRow(false);
  };

  const valueToString = (v: unknown): string => {
    if (v == null) {
      return "";
    }
    if (Array.isArray(v)) {
      return [...v].sort().join(", ");
    }
    return String(v);
  };

  const sortRows = (
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

  const handleHeaderSort = (colKey: string) => {
    if (sortColKey !== colKey) {
      setSortColKey(colKey);
      setSortDir("asc");
    } 
    else {
      setSortDir(sortDir === "asc" ? "desc" : "asc");
    }
  };

  const handleHomePage = () => {
    window.location.href = `${base}`;
  }

  const handleData = () => {
    window.location.href = `${base}${datasetFiles[dataset].csv}`;
  }

  const handleEditHistory = () => {
    window.location.href = `${base}${dataset}/history`;
  };

  const datasetLabels: Record<string, string> = {
    csprofs: "CS Professors",
    students: "Students",
  };

  const datasetLabel = datasetLabels[dataset];

  return (
    <div className="App" style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      <ActionButtons
        datasetLabel={datasetLabel}
        handleHomePage={handleHomePage}
        handleData={handleData}
        handleEditHistory={handleEditHistory}
        setIsAddingRow={setIsAddingRow}
        handleDeleteRow={handleDeleteRow}
        setIsDeletingRow={setIsDeletingRow}
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
          onCellEditorActivated={onCellEditorActivated}
          gridSelection={gridSelection}
          onGridSelectionChange={onGridSelectionChange}
          gridWidth={gridWidth}
          columnSchema={columnSchema}

          onHeaderSort={handleHeaderSort}
          sortColKey={sortColKey}
          sortDir={sortDir}
        />
      </div>

      <MultiSelectModal
        isOverlayVisible={isMultiOverlayVisible}
        setIsOverlayVisible={setIsMultiOverlayVisible}
        handleSaveOptions={handleSaveMulti}
        optionsList={activeOptionsList}
        multiple={true}
        selectedOptions={selectedMultiOptions}
        setSelectedOptions={setSelectedMultiOptions}
        title = "Select Value(s)"
      />

      <TextInputModal
        isOverlayVisible={isTextOverlayVisible}
        setIsOverlayVisible={setIsTextOverlayVisible}
        title={"Edit Value"}
        value={textDraft}
        setValue={setTextDraft}
        handleSave={handleSaveText}
      />

      <DropdownModal
        isOverlayVisible={isDropdownOverlayVisible}
        setIsOverlayVisible={setIsDropdownOverlayVisible}
        optionsList={activeOptionsList}
        title={"Select Value"}
        value={textDraft}
        setValue={setTextDraft}
        handleSave={handleSaveDropdown}
      />

      <DropdownFreeTextModal
        isOverlayVisible={isDropdownFreeTextOverlayVisible}
        setIsOverlayVisible={setIsDropdownFreeTextOverlayVisible}
        optionsList={activeOptionsList}
        title={"Select Value"}
        dropdownValue={textDraft}
        setDropdownValue={setTextDraft}
        draftValue={freeTextAltDraft}
        setDraftValue={setFreeTextAltDraft}
        handleSave={handleSaveDropdownFreeText}
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

      {isDeletingRow && (
        <DeleteRowFooter
          comment={deleteComment}
          setComment={setDeleteComment}
          handleDeleteRowConfirm={handleDeleteRowConfirm}
          onCancel={handleDeleteRowCancel}
          isCommentFilled={isDeleteCommentFilled}
        />
      )}
    </div>
  );
}
