import './App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import { DataEditor, GridCell, GridCellKind, GridColumn, Item, EditableGridCell, BubbleCell, GridSelection, CompactSelection } from "@glideapps/glide-data-grid";
import Papa, { ParseResult } from 'papaparse';
import { Modal, Button, Checkbox, TextField, Snackbar } from '@mui/material';
import MuiAlert, { AlertProps } from '@mui/material/Alert';

// Define the structure of the data for each professor
interface Professor {
  FullName: string;
  University: string;
  JoinYear: string;
  SubField: string[];
  Bachelors: string;
  Doctorate: string;
}

// List of possible options for the SubField filter
const optionsList = [
  "Artificial Intelligence",
  "Software Engineering",
  "Computer Security",
  "Databases",
  "Cryptography",
  "Programming Languages",
];

// Define custom column widths for each column
const columnWidths: { [key: string]: number } = {
  FullName: 200,
  University: 250,
  JoinYear: 100,
  SubField: 250,
  Bachelors: 350,
  Doctorate: 300,
};

// Custom alert component for displaying messages in the snackbar
const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

// Main App component
export default function App() {
  // State variables for managing columns, data, filters, and UI elements
  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [data, setData] = useState<Professor[]>([]);
  const [filteredData, setFilteredData] = useState<Professor[]>([]);
  const [columnFilters, setColumnFilters] = useState<{ [key: string]: string }>({});
  const [isOverlayVisible, setIsOverlayVisible] = useState(false);
  const [editingCell, setEditingCell] = useState<{ row: number; col: number } | null>(null);
  const [selectedOptions, setSelectedOptions] = useState<string[]>([]);
  const [gridSelection, setGridSelection] = useState<GridSelection>({
    current: undefined,
    rows: CompactSelection.empty(),
    columns: CompactSelection.empty(),
  });
  const [snackbarOpen, setSnackbarOpen] = useState(false);

  // Fetch and parse CSV data when the component loads
  useEffect(() => {
    const fetchCsvData = async () => {
      try {
        const response = await fetch('/csprofessors.csv');
        const csvData = await response.text();

        // Parse the CSV data and set the columns and data for the grid
        Papa.parse(csvData, {
          header: true,
          complete: (results: ParseResult<Professor>) => {
            const parsedData = results.data.filter(row => Object.values(row).some(value => value !== null && value !== ""));

            const gridColumns: GridColumn[] = Object.keys(parsedData[0] || {})
              .filter((key) => key !== "UniqueId") 
              .map((key) => ({
                title: key,
                width: columnWidths[key] || 150,
              }));

            setColumns(gridColumns);
            setData(parsedData);
            setFilteredData(parsedData);
          },
          skipEmptyLines: true,
        });
      } catch (error) {
        console.error('Error fetching the CSV file:', error);
      }
    };

    fetchCsvData();
  }, []);

  // Apply filters to the data when filters or data change
  useEffect(() => {
    const applyFilters = () => {
      const filtered = data.filter((row) => {
        return columns.every((col) => {
          const colTitle = col.title;
          const filterValue = columnFilters[colTitle];
          if (!filterValue) return true;

          const cellValue = row[colTitle as keyof Professor];
          return cellValue?.toString().toLowerCase().includes(filterValue.toLowerCase());
        });
      });
      setFilteredData(filtered);
    };

    applyFilters();
  }, [columnFilters, data, columns]);

  // Update the column filter values
  const handleColumnFilterChange = (colTitle: string, value: string) => {
    setColumnFilters((prevFilters) => ({
      ...prevFilters,
      [colTitle]: value,
    }));
  };

  // Generate data cells for the grid, including handling bubble cells for the SubField column
  const getData = ([col, row]: Item): GridCell => {
    const cellData = filteredData[row]?.[columns[col].title as keyof Professor];

    if (columns[col].title === "SubField") {
      const bubbleData = typeof cellData === 'string' ? cellData.split(',').map(item => item.trim()) : cellData;

      return {
        kind: GridCellKind.Bubble,
        data: Array.isArray(bubbleData) ? bubbleData : [],
        allowOverlay: true,
      };
    }

    const textData = Array.isArray(cellData) ? cellData.join(', ') : cellData || "";

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

    if (columns[col].title === "SubField") {
      setEditingCell({ row, col });
      const actualRowIndex = data.indexOf(filteredData[row]);
      setSelectedOptions(Array.isArray(data[actualRowIndex]?.SubField) ? data[actualRowIndex].SubField : []);
      setIsOverlayVisible(true);
    }
  };

  // Handle cell edits and update the data accordingly
  const onCellEdited = ([col, row]: Item, newValue: EditableGridCell | BubbleCell) => {
    const updatedData = [...data];
    const key = columns[col].title as keyof Professor;

    const actualRowIndex = data.indexOf(filteredData[row]);

    if (newValue.kind === GridCellKind.Text) {
      if (typeof newValue.data === 'string' && key !== 'SubField') {
        updatedData[actualRowIndex][key] = newValue.data as Professor[typeof key];
      }
    } else if (newValue.kind === GridCellKind.Bubble) {
      if (Array.isArray(newValue.data) && key === 'SubField') {
        updatedData[actualRowIndex][key] = newValue.data as string[] as Professor[typeof key];
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

  // Handle changes to checkbox options for the SubField column
  const handleOptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    setSelectedOptions((prev) =>
      event.target.checked ? [...prev, value] : prev.filter((opt) => opt !== value)
    );
  };

  // Delete the selected rows or cell and update the data
  const handleDeleteRow = () => {
    if (gridSelection.rows.length > 0) {
      const selectedRowIndices = gridSelection.rows.toArray();
      const selectedRows = selectedRowIndices.map(index => filteredData[index]);

      setData(prevData => prevData.filter(item => !selectedRows.includes(item)));
      setFilteredData(prevFilteredData => prevFilteredData.filter(item => !selectedRows.includes(item)));
      setGridSelection({
        current: undefined,
        rows: CompactSelection.empty(),
        columns: CompactSelection.empty(),
      });
    } else if (gridSelection.current && gridSelection.current.cell) {
      const { cell } = gridSelection.current;
      const [, row] = cell;
      const rowData = filteredData[row];
      setData(prevData => prevData.filter(item => item !== rowData));
      setFilteredData(prevFilteredData => prevFilteredData.filter(item => item !== rowData));
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

  return (
    <div className="App">
      {/* Render filter text fields for each column */}
      <div style={{ display: "flex", justifyContent: "space-evenly", padding: "10px", flexWrap: "wrap" }}>
        {columns.map((col) => (
          <TextField
            key={col.title}
            label={`Search ${col.title}`}
            variant="outlined"
            size="small"
            value={columnFilters[col.title] || ""}
            onChange={(e) => handleColumnFilterChange(col.title, e.target.value)}
            style={{ marginBottom: "20px", width: columnWidths[col.title] }}
          />
        ))}
      </div>

      {/* Button to delete selected rows or cells */}
      <div style={{ padding: "10px" }}>
        <Button variant="contained" color="primary" onClick={handleDeleteRow}>
          Delete Row
        </Button>
      </div>

      {/* Render the DataEditor with filtered data */}
      {columns.length > 0 && filteredData.length > 0 && (
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
          width="200vw"
          height="80vh"
        />
      )}

      {/* Modal for selecting SubField options */}
      <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
        <div className="overlay-content" style={{ background: "white", padding: "20px", borderRadius: "10px", margin: "50px auto", width: "300px" }}>
          <h3>Select Options</h3>
          {optionsList.map((option) => (
            <div key={option}>
              <Checkbox value={option} checked={selectedOptions.includes(option)} onChange={handleOptionChange} />
              <label>{option}</label>
            </div>
          ))}
          <Button variant="contained" color="primary" onClick={handleSaveOptions}>Save</Button>
        </div>
      </Modal>

      {/* Snackbar for alerting the user if no selection is made */}
      <Snackbar open={snackbarOpen} autoHideDuration={3000} onClose={handleSnackbarClose}>
        <Alert onClose={handleSnackbarClose} severity="warning">
          Please select a cell or row first.
        </Alert>
      </Snackbar>
    </div>
  );
}
