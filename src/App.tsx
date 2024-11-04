import './App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import {
  DataEditor,
  GridCell,
  GridCellKind,
  GridColumn,
  Item,
  EditableGridCell,
  BubbleCell,
  GridSelection,
  CompactSelection,
} from "@glideapps/glide-data-grid";
import Papa, { ParseResult } from 'papaparse';
import { Modal, Button, TextField, Snackbar, Autocomplete } from '@mui/material';
import MuiAlert, { AlertProps } from '@mui/material/Alert';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import IconButton from '@mui/material/IconButton';

// Define the structure of the data for each professor
interface Professor {
  FullName: string;
  University: string;
  JoinYear: string;
  SubField: string[];
  Bachelors: string;
  Doctorate: string;
}

// List of valid keys and corresponding TypeScript type
const validKeys = ["FullName", "University", "JoinYear", "SubField", "Bachelors", "Doctorate"] as const;
type ProfessorKey = typeof validKeys[number];

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
const columnWidths: { [key in ProfessorKey]: string } = {
  FullName: '15%',
  University: '20%',
  JoinYear: '5%',
  SubField: '20%',
  Bachelors: '20%',
  Doctorate: '20%',
};

// Custom alert component for displaying messages in the snackbar
const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

// Get initial filters from URL parameters
const params = new URLSearchParams(window.location.search);
const initialFilters: { [key in ProfessorKey]?: string } = {};
validKeys.forEach((name) => {
  const value = params.get(name);
  if (value) {
    initialFilters[name] = value;
  }
});

// Main App component
export default function App() {
  // State variables for managing columns, data, filters, and UI elements
  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [data, setData] = useState<Professor[]>([]);
  const [filteredData, setFilteredData] = useState<Professor[]>([]);
  const [columnFilters, setColumnFilters] = useState<{ [key in ProfessorKey]?: string }>(initialFilters);
  const [isOverlayVisible, setIsOverlayVisible] = useState(false);
  const [editingCell, setEditingCell] = useState<{ row: number; col: number } | null>(null);
  const [selectedOptions, setSelectedOptions] = useState<string[]>([]);
  const [gridSelection, setGridSelection] = useState<GridSelection>({
    current: undefined,
    rows: CompactSelection.empty(),
    columns: CompactSelection.empty(),
  });
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [gridWidth, setGridWidth] = useState(window.innerWidth);
  const [isAddingRow, setIsAddingRow] = useState(false);

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

  useEffect(() => {
    const handleResize = () => {
      setGridWidth(window.innerWidth);
    };
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  // Fetch and parse CSV data when the component loads
  useEffect(() => {
    const fetchCsvData = async () => {
      try {
        const response = await fetch('/csprofessors.csv');
        const csvData = await response.text();

        // Parse the CSV data and set the columns and data for the grid
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

  // Modify the handleAddRow function
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
      <div
        style={{
          display: "flex",
          justifyContent: "space-evenly",
          padding: "10px",
          flexWrap: "wrap",
        }}
      >
        {columns.map((col) => {
          const colKey = col.id as ProfessorKey;
          return (
            <TextField
              key={colKey}
              label={`Search ${col.title}`}
              variant="outlined"
              size="small"
              value={columnFilters[colKey] || ""}
              onChange={(e) => handleColumnFilterChange(colKey, e.target.value)}
              style={{ marginBottom: "20px", width: columnWidths[colKey] }}
            />
          );
        })}
      </div>

      {/* Button to delete selected rows or cells */}
      <div style={{ padding: "10px" }}>
        <Button variant="contained" color="primary" onClick={handleDeleteRow}>
          Delete Row
        </Button>

        {/* Button to add a new row */}
        <Button
          variant="contained"
          color="primary"
          onClick={() => setIsAddingRow(true)}
          style={{ marginLeft: "10px" }}
        >
          Add Row
        </Button>
      </div>

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
      <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
        <div
          className="overlay-content"
          style={{
            background: "white",
            padding: "20px",
            borderRadius: "10px",
            margin: "50px auto",
            width: "400px",
          }}
        >
          <h3>Select SubFields</h3>
          <Autocomplete
            multiple
            options={optionsList}
            getOptionLabel={(option) => option}
            value={selectedOptions}
            onChange={(event, newValue) => {
              setSelectedOptions(newValue);
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                variant="outlined"
                placeholder="Select SubFields"
              />
            )}
            style={{ marginBottom: "20px" }}
          />
          <Button variant="contained" color="primary" onClick={handleSaveOptions}>
            Save
          </Button>
        </div>
      </Modal>

      {/* Snackbar for alerting the user if no selection is made */}
      <Snackbar open={snackbarOpen} autoHideDuration={3000} onClose={handleSnackbarClose}>
        <Alert onClose={handleSnackbarClose} severity="warning">
          Please select a cell or row first.
        </Alert>
      </Snackbar>

      {/* Conditionally render the Sticky Footer for Adding New Row */}
      {isAddingRow && (
        <div
          style={{
            position: "fixed",
            bottom: 0,
            left: 0,
            right: 0,
            backgroundColor: "white",
            padding: "10px",
            display: "flex",
            alignItems: "center",
            zIndex: 1000,
            borderTop: "1px solid #ccc",
            flexWrap: "nowrap",
            overflowX: "auto",
          }}
        >
          <div style={{ display: "flex", flex: 1, alignItems: "center", overflowX: "auto" }}>
            {validKeys.map((key) => {
              if (key === "SubField") {
                // Render Autocomplete for SubField
                return (
                  <Autocomplete
                    key={key}
                    multiple
                    options={optionsList}
                    getOptionLabel={(option) => option}
                    value={newRowData.SubField}
                    onChange={(event, newValue) => {
                      setNewRowData((prevData) => ({
                        ...prevData,
                        SubField: newValue,
                      }));
                    }}
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        variant="outlined"
                        label="SubField"
                        placeholder="Select SubFields"
                        size="small"
                        style={{ margin: "5px", minWidth: "150px" }}
                      />
                    )}
                    style={{ margin: "5px", minWidth: "150px" }}
                  />
                );
              } else {
                return (
                  <TextField
                    key={key}
                    label={key}
                    variant="outlined"
                    size="small"
                    value={newRowData[key] || ""}
                    onChange={(e) =>
                      setNewRowData((prevData) => ({
                        ...prevData,
                        [key]: e.target.value,
                      }))
                    }
                    style={{ margin: "5px", minWidth: "150px" }}
                  />
                );
              }
            })}
          </div>

          {/* Buttons */}
          <div style={{ display: "flex", alignItems: "center", marginLeft: "auto" }}>
            {/* Add Button */}
            <IconButton
              color="primary"
              onClick={handleAddRowConfirm}
              disabled={!allFieldsFilled}
            >
              <CheckCircleIcon />
            </IconButton>

            {/* Cancel Button */}
            <Button
              variant="contained"
              color="primary"
              onClick={() => {
                setIsAddingRow(false);
                // Reset newRowData to empty values
                setNewRowData({
                  FullName: "",
                  University: "",
                  JoinYear: "",
                  SubField: [],
                  Bachelors: "",
                  Doctorate: "",
                });
              }}
              style={{ marginLeft: "10px" }}
            >
              Cancel
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
