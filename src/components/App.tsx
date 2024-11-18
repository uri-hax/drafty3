// components/App.tsx
import '../../public/App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import {
  DataEditor,
  type GridCell,
  GridCellKind,
  type GridColumn,
  type Item,
  type EditableGridCell,
  type BubbleCell,
  type GridSelection,
  CompactSelection,
} from "@glideapps/glide-data-grid";
import { Snackbar, Button, Modal, TextField, Autocomplete } from '@mui/material';

import { validKeys, type Professor, type ProfessorKey } from '../interfaces/Professor';
import { columnWidths, optionsList } from '../utils/constants';
import { fetchCsvData } from '../utils/csvParser';
import useWindowWidth from '../hooks/useWindow';

import Alert from './Alerts';
import FilterBar from './Filter';
import ActionButtons from './ActionButtons';
import ProfessorGrid from './ProfessorGrid';
import SubFieldModal from './SubFieldModal';
import AddRowFooter from './AddRowFooter';

// Get initial filters from URL parameters
const params = new URLSearchParams(window.location.search);
const initialFilters: { [key in ProfessorKey]?: string } = {};
validKeys.forEach((name) => {
  const value = params.get(name);
  if (value) {
    initialFilters[name] = value;
  }
});

export default function App() {
  // State variables for managing columns, data, filters, and UI elements
  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [data, setData] = useState<Professor[]>([]);
  const [filteredData, setFilteredData] = useState<Professor[]>([]);
  const [columnFilters, setColumnFilters] = useState<{ [key in ProfessorKey]?: string }>(initialFilters);
  const [isOverlayVisible, setIsOverlayVisible] = useState<boolean>(false);
  const [editingCell, setEditingCell] = useState<{ row: number; col: number } | null>(null);
  const [selectedOptions, setSelectedOptions] = useState<string[]>([]);
  const [gridSelection, setGridSelection] = useState<GridSelection>({
    current: undefined,
    rows: CompactSelection.empty(),
    columns: CompactSelection.empty(),
  });
  const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
  const [isAddingRow, setIsAddingRow] = useState<boolean>(false);

  // Use custom hook to get window width
  const gridWidth = useWindowWidth();

  // State for new row data
  const [newRowData, setNewRowData] = useState<Professor>({
    FullName: "",
    University: "",
    JoinYear: "",
    SubField: [],
    Bachelors: "",
    Doctorate: "",
  });

  const allFieldsFilled =
    newRowData.FullName.trim() !== "" &&
    newRowData.University.trim() !== "" &&
    newRowData.JoinYear.trim() !== "" &&
    newRowData.SubField.length > 0 &&
    newRowData.Bachelors.trim() !== "" &&
    newRowData.Doctorate.trim() !== "";

  // Fetch and parse CSV data when the component loads
  useEffect(() => {
    const fetchData = async () => {
      try {
        const { gridColumns, parsedData } = await fetchCsvData(gridWidth);
        setColumns(gridColumns);
        setData(parsedData);
        setFilteredData(parsedData);
      } catch (error) {
        console.error('Error fetching the CSV file:', error);
      }
    };
    fetchData();
  }, [gridWidth]);

  // Apply filters to the data when filters or data change
  useEffect(() => {
    const applyFilters = () => {
      const filtered = data.filter((row) => {
        return columns.every((col) => {
          const colKey = col.id as ProfessorKey;
          const filterValue = columnFilters[colKey];
          if (!filterValue) return true;

          const cellValue = row[colKey];
          if (Array.isArray(cellValue)) {
            return cellValue.some((val) =>
              val.toString().toLowerCase().includes(filterValue.toLowerCase())
            );
          }
          return cellValue?.toString().toLowerCase().includes(filterValue.toLowerCase());
        });
      });
      setFilteredData(filtered);
    };

    applyFilters();
  }, [columnFilters, data, columns]);

  // Update the URL parameters when filters change
  useEffect(() => {
    const params = new URLSearchParams();
    validKeys.forEach((key) => {
      if (columnFilters[key]) {
        params.set(key, columnFilters[key] as string);
      }
    });
    window.history.replaceState(null, '', '?' + params.toString());
  }, [columnFilters]);

  // Update the column filter values
  const handleColumnFilterChange = (colKey: ProfessorKey, value: string) => {
    setColumnFilters((prevFilters) => ({
      ...prevFilters,
      [colKey]: value,
    }));
  };

  // Generate data cells for the grid, including handling bubble cells for the SubField column
  const getData = ([col, row]: Item): GridCell => {
    const colKey = columns[col].id as ProfessorKey;
    const cellData = filteredData[row]?.[colKey];

    if (colKey === "SubField") {
      const bubbleData = Array.isArray(cellData) ? cellData : [];
      return {
        kind: GridCellKind.Bubble,
        data: bubbleData,
        allowOverlay: true,
      };
    }

    const textData = cellData !== undefined && cellData !== null ? String(cellData) : "";
    return {
      kind: GridCellKind.Text,
      data: textData,
      displayData: textData,
      allowOverlay: true,
    };
  };

  // Handle selection changes in the grid
  const onGridSelectionChange = (newSelection: GridSelection) => {
    setGridSelection(newSelection);
  };

  // Handle activating a cell for editing, especially for SubField bubble cells
  const onCellActivated = (cell: Item) => {
    const [col, row] = cell;
    const colKey = columns[col].id as ProfessorKey;

    if (colKey === "SubField") {
      setEditingCell({ row, col });
      const actualRowIndex = data.indexOf(filteredData[row]);
      setSelectedOptions(Array.isArray(data[actualRowIndex]?.SubField) ? data[actualRowIndex].SubField : []);
      setIsOverlayVisible(true);
    }
  };

  // Handle cell edits and update the data accordingly
  const onCellEdited = ([col, row]: Item, newValue: EditableGridCell | BubbleCell) => {
    const updatedData = [...data];
    const key = columns[col].id as ProfessorKey;
    const actualRowIndex = data.indexOf(filteredData[row]);

    if (key === "SubField") {
      if (newValue.kind === GridCellKind.Bubble) {
        updatedData[actualRowIndex]["SubField"] = newValue.data as string[];
      } else {
        console.error("Expected BubbleCell for SubField");
      }
    } else {
      if (newValue.kind === GridCellKind.Text) {
        updatedData[actualRowIndex][key] = newValue.data as string;
      } else {
        console.error(`Expected TextCell for ${key}`);
      }
    }

    setData(updatedData);
    setFilteredData(updatedData);
  };

  // Save selected options for bubble cells (SubField column)
  const handleSaveOptions = () => {
    if (editingCell) {
      const { row } = editingCell;
      const updatedData = [...data];
      const actualRowIndex = data.indexOf(filteredData[row]);
      updatedData[actualRowIndex].SubField = [...selectedOptions];
      setData(updatedData);
      setFilteredData(updatedData);
      setIsOverlayVisible(false);
    }
  };

  // Delete the selected rows or cell and update the data
  const handleDeleteRow = () => {
    if (gridSelection.rows.length > 0) {
      const selectedRowIndices = gridSelection.rows.toArray();
      const selectedRows = selectedRowIndices.map((index) => filteredData[index]);

      setData((prevData) => prevData.filter((item) => !selectedRows.includes(item)));
      setFilteredData((prevFilteredData) => prevFilteredData.filter((item) => !selectedRows.includes(item)));
      setGridSelection({
        current: undefined,
        rows: CompactSelection.empty(),
        columns: CompactSelection.empty(),
      });
    } else if (gridSelection.current && gridSelection.current.cell) {
      const { cell } = gridSelection.current;
      const [, row] = cell;
      const rowData = filteredData[row];
      setData((prevData) => prevData.filter((item) => item !== rowData));
      setFilteredData((prevFilteredData) => prevFilteredData.filter((item) => item !== rowData));
      setGridSelection({
        current: undefined,
        rows: CompactSelection.empty(),
        columns: CompactSelection.empty(),
      });
    } else {
      setSnackbarOpen(true);
    }
  };

  // Handle closing the snackbar when a row or cell is not selected before deletion
  const handleSnackbarClose = () => {
    setSnackbarOpen(false);
  };

  // Handle adding a new row
  const handleAddRowConfirm = () => {
    if (!allFieldsFilled) return; // Ensure all fields are filled

    setData((prevData) => [...prevData, newRowData]);
    setFilteredData((prevFilteredData) => [...prevFilteredData, newRowData]);

    // Reset newRowData to empty values
    setNewRowData({
      FullName: "",
      University: "",
      JoinYear: "",
      SubField: [],
      Bachelors: "",
      Doctorate: "",
    });

    // Hide the sticky footer
    setIsAddingRow(false);
  };

  return (
    <div className="App" style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      {/* Render filter text fields for each column */}
      <FilterBar
        columns={columns}
        columnFilters={columnFilters}
        handleColumnFilterChange={handleColumnFilterChange}
        columnWidths={columnWidths}
      />

      {/* Action Buttons */}
      <ActionButtons
        handleDeleteRow={handleDeleteRow}
        setIsAddingRow={setIsAddingRow}
      />

      {/* Render the DataEditor with filtered data */}
      <div className="grid-container" style={{ flexGrow: 1 }}>
        <DataEditor
          columns={columns}
          getCellContent={getData}
          rows={filteredData.length}
          onCellEdited={onCellEdited}
          rowMarkers="number"
          onCellActivated={onCellActivated}
          onGridSelectionChange={onGridSelectionChange}
          gridSelection={gridSelection}
          showSearch={false}
          width={gridWidth}
          height="100%"
        />
      </div>

      {/* Modal for selecting SubField options in the grid cell editor */}
      <SubFieldModal
        isOverlayVisible={isOverlayVisible}
        setIsOverlayVisible={setIsOverlayVisible}
        selectedOptions={selectedOptions}
        setSelectedOptions={setSelectedOptions}
        handleSaveOptions={handleSaveOptions}
      />

      {/* Snackbar for alerting the user if no selection is made */}
      <Snackbar open={snackbarOpen} autoHideDuration={3000} onClose={handleSnackbarClose}>
        <Alert onClose={handleSnackbarClose} severity="warning">
          Please select a cell or row first.
        </Alert>
      </Snackbar>

      {/* Conditionally render the Sticky Footer for Adding New Row */}
      {isAddingRow && (
        <AddRowFooter
          validKeys={validKeys}
          newRowData={newRowData}
          setNewRowData={setNewRowData}
          optionsList={optionsList}
          handleAddRowConfirm={handleAddRowConfirm}
          setIsAddingRow={setIsAddingRow}
          allFieldsFilled={allFieldsFilled}
        />
      )}
    </div>
  );
}
