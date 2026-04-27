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
import { checkBackendHealth } from '../lib/health';
import FreeTextModal from './FreeTextModal';
import { datasetFiles, datasetLabels } from "../config/AppConfig";
import {
  getColumnMatches,
  valueToString,
  sortRows,
  handleAlertSnackbarClose,
  handleContributionSnackbarClose,
  handleHomePage,
  handleData,
  handleEditHistory
} from "../utils/gridHelper";


export default function App() {
  // create portal div if it doesn't exist to render modals into
  useEffect(() => {
    const portalDiv = document.getElementById('portal');

    if (!portalDiv) {
      const newPortalDiv = document.createElement('div');
      newPortalDiv.id = 'portal';
      document.body.appendChild(newPortalDiv);
    }
  }, []);

  const [session, setSession] = useState<BackendSession | null>(null);
  
  const [backendAvailable, setBackendAvailable] = useState<boolean>(false);

  // make sure we have a session before doing anything else
  useEffect(() => {
    if (!backendAvailable) return;

    ensureSession()
      .then((data) => setSession(data.session))
      .catch(console.error);
  }, [backendAvailable]);

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

  const [alertSnackbarOpen, setAlertSnackbarOpen] = useState<boolean>(false);
  const [contributionSnackbarOpen, setContributionSnackbarOpen] = useState<boolean>(false);

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

  const [flashDowntimeMessage, setFlashDowntimeMessage] = useState(false);

  // get the dataset we're in from the url
  const base = import.meta.env.BASE_URL;
  const dataset = window.location.pathname
    .replace(base, "")
    .split("/")
    .filter(Boolean)[0];

  // gather data from csv and yaml, get suggestion type ids, and apply any filters from the url and record the searches if so
  useEffect(() => {
    const fetchData = async () => {
      try {
        // only do certain things if backend is healthy and available
        const isHealthy = await checkBackendHealth();
        setBackendAvailable(isHealthy);
        console.log("Backend health:", isHealthy);

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

        if (isHealthy) {
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
            if (!pair) continue;

            const [name, id] = pair;
            idMap[name] = id;
          }

          setSuggestionTypeIds(idMap);
        }
  
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

        if (isHealthy) {
          for (const [colKey, value] of Object.entries(initialFilters)) {
            const trimmed = value.trim();
            if (!trimmed) continue;

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

  // re-apply filters and sorting whenever relevant state changes
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

  // update url when filters change
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

  // define widths for columns based on schema config
  const customWidths: Record<string, string> = {};
    for (const key of Object.keys(columnSchema)) {
      if (columnSchema[key]?.width) {
        customWidths[key] = columnSchema[key].width;
      }
    }

  // handle the filter change on a column based on the search input and record the search interaction
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

    if (backendAvailable) {
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
    }
  };

  // determine if all required fields for adding a new row are filled and valid
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

  // determine if the delete comment is filled or not
  const isDeleteCommentFilled = React.useMemo(() => {
    return deleteComment.trim() !== "";
  }, [deleteComment]);

  // get the suggestion type values for a given column to show in the dropdowns on editing
  const getOptionsForCol = (colKey: string) =>
    suggestionTypeValues[colKey] ?? optionsLists[colKey] ?? [];

  // ensure suggestion type values are loaded for a given column by fetching from the backend based on the suggestion type id if not already loaded
  const ensureSuggestionValuesLoaded = (colKey: string) => {
    if (suggestionTypeValues[colKey]) {
      return;
    }

    const id = suggestionTypeIds[colKey];
    if (!id) {
      return;
    }

    if (backendAvailable) {
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
    }
  };

  // handle when a cell is activated for editing by showing the appropriate overlay based on the column config for that cell
  const onCellEditorActivated = (cell: Item) => {
    if (!backendAvailable) {
      showDowntimeMessage();
      return;
    }

    const [col, row] = cell;
    const colKey = columns[col].id as string;

    const rowObj = filteredData[row];

    setIsMultiOverlayVisible(false);
    setIsTextOverlayVisible(false);
    setIsDropdownOverlayVisible(false);
    setIsDropdownFreeTextOverlayVisible(false);

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

    // fallback is free text
    setTextDraft(cellVal != null ? String(cellVal) : "");
    setIsTextOverlayVisible(true);
  };

  // handle saving an edit from the multi select modal by updating the data and filtered data states, closing the modal, and recording the edit interaction
  const handleSaveMulti = () => {
    if (editingCell && backendAvailable) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = [...selectedMultiOptions];
      setData(updatedData);
      setFilteredData(updatedData);
      setIsMultiOverlayVisible(false);
      setEditingCell(null);

      const rowObj = filteredData[row];
      const rowId = Number(rowObj["idUniqueID"]);
      const columnId = suggestionTypeIds[colKey];
      const suggestion = JSON.stringify(selectedMultiOptions);

      recordCellEdit(columnId, rowId, suggestion);

      setContributionSnackbarOpen(true);
    }
  };

  // handle saving an edit from the free text modal by updating the data and filtered data states, closing the modal, and recording the edit interaction
  const handleSaveText = () => {
    if (editingCell && backendAvailable) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = textDraft
      setData(updatedData);
      setFilteredData(updatedData);
      setIsTextOverlayVisible(false);
      setEditingCell(null);

      const rowObj = filteredData[row];
      const rowId = Number(rowObj["idUniqueID"]);
      const columnId = suggestionTypeIds[colKey];

      recordCellEdit(columnId, rowId, textDraft);

      setContributionSnackbarOpen(true);
    }
  };

  // handle saving an edit from the dropdown modal by updating the data and filtered data states, closing the modal, and recording the edit interaction
  const handleSaveDropdown = () => {
    if (editingCell && backendAvailable) {
      const { row, colKey } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex][colKey] = textDraft;
      setData(updatedData);
      setFilteredData(updatedData);
      setIsDropdownOverlayVisible(false);
      setEditingCell(null);

      const rowObj = filteredData[row];
      const rowId = Number(rowObj["idUniqueID"]);
      const columnId = suggestionTypeIds[colKey];

      recordCellEdit(columnId, rowId, textDraft);

      setContributionSnackbarOpen(true);
    }
  };

  // handle saving an edit from the dropdown + free text modal by updating the data and filtered data states, closing the modal, and recording the edit interaction
  const handleSaveDropdownFreeText = () => {
    if (editingCell && backendAvailable) {
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

      const rowObj = filteredData[row];
      const rowId = Number(rowObj["idUniqueID"]);
      const columnId = suggestionTypeIds[colKey];

      recordCellEdit(columnId, rowId, finalValue);

      setContributionSnackbarOpen(true);
    }
  };

  // old handle cell edit function for free text edits before modals (not currently used)
  const onCellEdited = (cell: Item, newValue: EditableGridCell | BubbleCell) => {
    if (!backendAvailable) {
      showDowntimeMessage();
      return;
    }

    const [col, row] = cell;
    const colKey = columns[col].id as string;
    const colType = columnSchema[colKey].type || 'string';
    const updatedData = [...data];
    const actualRowIndex = data.indexOf(filteredData[row]);

    let suggestion = "";

    if (colType === "string[]") {
      if (newValue.kind === "bubble") {
        const bubbleValue = newValue.data as string[];
        updatedData[actualRowIndex][colKey] = bubbleValue;
        suggestion = JSON.stringify(bubbleValue);
      }
    } else {
      if (newValue.kind === "text") {
        const textValue = newValue.data as string;
        updatedData[actualRowIndex][colKey] = textValue;
        suggestion = textValue;
      }
    }

    setData(updatedData);
    setFilteredData(updatedData);

    const rowObj = filteredData[row];
    const rowId = Number(rowObj["idUniqueID"]);
    const columnId = suggestionTypeIds[colKey];

    recordCellEdit(columnId, rowId, suggestion);

    setContributionSnackbarOpen(true);
  };

  // handle when the grid selection changes by recording the cell click interaction if a cell is selected and updating the grid selection state
  const onGridSelectionChange = (newSelection: GridSelection) => {
    const cell = newSelection.current?.cell;

    if (cell && backendAvailable) {
      const [col, row] = cell;
      const colKey = columns[col].id as string;
      const rowObj = filteredData[row];

      if (rowObj) {
        const rowId = Number(rowObj["idUniqueID"]);
        const columnId = suggestionTypeIds[colKey];
        recordCellClick(columnId, rowId, rowObj);
      }
    }

    setGridSelection(newSelection);
  };

  // handle deleting a row by showing the delete confirmation footer if a row is selected or a cell is selected, and showing an alert if no row or cell is selected
  const handleDeleteRow = () => {
    if (!backendAvailable) {
      showDowntimeMessage();
      return;
    }

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
      setAlertSnackbarOpen(true);
      setIsDeletingRow(false);
      return;
    }

    setPendingDeleteRows(selectedRows);
    setDeleteComment("");
    setIsDeletingRow(true);
    setIsAddingRow(false);
  };

  // handle confirming the deletion of a row by updating the data and filtered data states to remove the deleted row, resetting the grid selection, recording the delete interaction with the provided comment, and showing a snackbar confirming the contribution
  const handleDeleteRowConfirm = () => {
    if (!backendAvailable) {
      showDowntimeMessage();
      return;
    }

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

    // only delete one row at a time for now
    const rowObj = pendingDeleteRows[0];
    const rowId = Number(rowObj["idUniqueID"]);
    recordRowDelete(rowId, deleteComment);

    setIsDeletingRow(false);
    setDeleteComment("");
    setPendingDeleteRows([]);

    setContributionSnackbarOpen(true);
  };

  // handle canceling the deletion of a row by resetting the relevant state
  const handleDeleteRowCancel = () => {
    setIsDeletingRow(false);
    setDeleteComment("");
    setPendingDeleteRows([]);
  };

  // handle confirming the addition of a row by validating that all required fields are filled, sending the add row interaction data to the backend to be added to the database and getting back the new row id, updating the data and filtered data states to include the new row, resetting the new row form, and showing a snackbar confirming the contribution
  const handleAddRowConfirm = () => {
    if (!backendAvailable) {
      showDowntimeMessage();
      return;
    }

    if (!allFieldsFilled) {
      return;
    }

    const cells = columns.map((col) => {
      const colKey = col.id as string;
      const rawValue = newRowData[colKey];

      let suggestion = "";

      if (Array.isArray(rawValue)) {
        suggestion = JSON.stringify(rawValue);
      } 
      else if (rawValue != null) {
        suggestion = String(rawValue);
      }

      return {
        IDSuggestionType: suggestionTypeIds[colKey],
        Suggestion: suggestion,
        Active: 1, // default
        Confidence: 0, // default
      };
    });

    recordRowAdd(
      cells,
      async (res) => {
        const responseData = await res.json();
        const newRowId = Number(responseData.idUniqueID);

        const rowToAdd: ColumnData = {
          ...newRowData,
          idUniqueID: String(newRowId),
        };

        setData((prev) => [...prev, rowToAdd]);
        setFilteredData((prev) => [...prev, rowToAdd]);

        const resetObj: ColumnData = {};
        for (const key of Object.keys(columnSchema)) {
          resetObj[key] = columnSchema[key].type === "string[]" ? [] : "";
        }

        setNewRowData(resetObj);
        setIsAddingRow(false);

        setContributionSnackbarOpen(true);
      },
      (err) => {
        console.error("Failed to add row:", err);
      }
    );
  };

  // handle sorting when a column header is clicked by updating the sort column and direction state
  const handleHeaderSort = (colKey: string) => {
    if (sortColKey !== colKey) {
      setSortColKey(colKey);
      setSortDir("asc");
    } 
    else {
      setSortDir(sortDir === "asc" ? "desc" : "asc");
    }
  };

  const showDowntimeMessage = () => {
    setFlashDowntimeMessage(true);

    setTimeout(() => {
      setFlashDowntimeMessage(false);
    }, 600);
  };

  // get the label for the current dataset to show in the UI
  const datasetLabel = datasetLabels[dataset];

  // render the app with the action buttons at the top, filter bar below that if there are columns, the data grid taking up the rest of the space, and various modals and snackbars as needed
  return (
    <div className="App" style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      <ActionButtons
        backendAvailable={backendAvailable}
        flashDowntimeMessage={flashDowntimeMessage}
        showDowntimeMessage={showDowntimeMessage}
        datasetLabel={datasetLabel}
        handleHomePage={() => handleHomePage(base)}
        handleData={() => handleData(base, dataset)}
        handleEditHistory={() => handleEditHistory(base, dataset)}
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
        column={editingCell?.colKey}
      />

      <FreeTextModal
        isOverlayVisible={isTextOverlayVisible}
        setIsOverlayVisible={setIsTextOverlayVisible}
        title={"Edit Value"}
        value={textDraft}
        setValue={setTextDraft}
        column={editingCell?.colKey}
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
        column={editingCell?.colKey}
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
        column={editingCell?.colKey}
      />

      <Snackbar
        open={alertSnackbarOpen}
        autoHideDuration={5000}
        onClose={() => handleAlertSnackbarClose(setAlertSnackbarOpen)}
      >
        <Alert
          onClose={() => handleAlertSnackbarClose(setAlertSnackbarOpen)}
          severity="warning"
        >
          Please select a cell or row first.
        </Alert>
      </Snackbar>

      <Snackbar
        open={contributionSnackbarOpen}
        autoHideDuration={5000}
        onClose={() => handleContributionSnackbarClose(setContributionSnackbarOpen)}
      >
        <Alert
          onClose={() => handleContributionSnackbarClose(setContributionSnackbarOpen)}
          severity="success"
          sx={{
            backgroundColor: "#0b89ff",
            color: "#fff",
          }}
        >
          Thanks for contributing! You'll see your edits fully reflect within 15 minutes.
        </Alert>
      </Snackbar>

      {isAddingRow && backendAvailable && (
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

      {isDeletingRow && backendAvailable && (
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
