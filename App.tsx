import './App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import { DataEditor, GridCell, GridCellKind, GridColumn, Item, EditableGridCell, TextCell, BubbleCell } from "@glideapps/glide-data-grid";
import Papa, { ParseResult } from 'papaparse';
import { Modal, Button, Checkbox, TextField } from '@mui/material';

interface Professor {
  FullName: string;
  University: string;
  JoinYear: string;
  SubField: string[];
  Bachelors: string;
  Doctorate: string;
}

// List of bubble options available for the SubField column
const optionsList = [
  "Artificial Intelligence",
  "Software Engineering",
  "Computer Security",
  "Databases",
  "Cryptography",
  "Programming Languages",
];

// Pre-defined column widths for each data field
const columnWidths: { [key: string]: number } = {
  FullName: 200,
  University: 250,
  JoinYear: 100,
  SubField: 250,
  Bachelors: 350,
  Doctorate: 300,
};

export default function App() {
  // States for managing columns, data, filtered data, and UI states
  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [data, setData] = useState<Professor[]>([]);
  const [filteredData, setFilteredData] = useState<Professor[]>([]);
  const [columnFilters, setColumnFilters] = useState<{ [key: string]: string }>({}); // Track search values for each column
  const [isOverlayVisible, setIsOverlayVisible] = useState(false);
  const [editingCell, setEditingCell] = useState<{ row: number; col: number } | null>(null);
  const [selectedOptions, setSelectedOptions] = useState<string[]>([]);

  // Fetch and parse CSV data when the component mounts
  useEffect(() => {
    const fetchCsvData = async () => {
      try {
        const response = await fetch('/csprofessors.csv');
        const csvData = await response.text();

        // Parse CSV data using PapaParse
        Papa.parse(csvData, {
          header: true, // Use first row as headers
          complete: (results: ParseResult<Professor>) => {
            const parsedData = results.data.filter(row => Object.values(row).some(value => value !== null && value !== ""));

            // Map CSV headers to grid columns
            const gridColumns: GridColumn[] = Object.keys(parsedData[0] || {}).map((key) => ({
              title: key,
              width: columnWidths[key] || 150,
            }));

            setColumns(gridColumns);
            setData(parsedData);
            setFilteredData(parsedData);
          },
          skipEmptyLines: true, // Ignore empty lines
        });
      } catch (error) {
        console.error('Error fetching the CSV file:', error);
      }
    };

    fetchCsvData();
  }, []);

  // Apply filters whenever `columnFilters` or `data` changes
  useEffect(() => {
    const applyFilters = () => {
      // Start with the original data and apply each column filter individually
      const filtered = data.filter((row) => {
        return columns.every((col) => {
          const colTitle = col.title;
          const filterValue = columnFilters[colTitle];
          if (!filterValue) return true; // No filter applied for this column, so include all rows

          const cellValue = row[colTitle as keyof Professor];
          return cellValue?.toString().toLowerCase().includes(filterValue.toLowerCase());
        });
      });
      setFilteredData(filtered); // Update filtered data based on column filters
    };

    applyFilters(); // Apply filters whenever columnFilters or data changes
  }, [columnFilters, data, columns]);

  // Handle changes to column-specific search inputs
  const handleColumnFilterChange = (colTitle: string, value: string) => {
    setColumnFilters((prevFilters) => ({
      ...prevFilters,
      [colTitle]: value, // Update the filter value for the specific column
    }));
  };

  const getData = ([col, row]: Item): GridCell => {
    const cellData = filteredData[row]?.[columns[col].title as keyof Professor];

    // Check if the cell corresponds to the "SubField" column to render it as a BubbleCell
    if (columns[col].title === "SubField") {
      // Convert cellData to an array if it's a string or use an empty array if undefined
      const bubbleData = typeof cellData === 'string' ? cellData.split(',').map(item => item.trim()) : cellData;

      return {
        kind: GridCellKind.Bubble,
        data: Array.isArray(bubbleData) ? bubbleData : [],
        allowOverlay: true,
      };
    }

    // Ensure that the data for TextCell is always a string
    const textData = Array.isArray(cellData) ? cellData.join(', ') : cellData || "";

    return {
      kind: GridCellKind.Text,
      data: textData, // Assign the correctly formatted text data here
      displayData: textData, // Display the formatted text data
      allowOverlay: true,
    };
  };

  // Handle activation of cell (e.g., click to open modal for editing SubField)
  const onCellActivated = (cell: Item) => {
    const [col, row] = cell;

    if (columns[col].title === "SubField") {
      setEditingCell({ row, col });
      setSelectedOptions(Array.isArray(data[row]?.SubField) ? data[row].SubField : []);
      setIsOverlayVisible(true);
    }
  };

  const onCellEdited = ([col, row]: Item, newValue: EditableGridCell | BubbleCell) => {
    const updatedData = [...data]; // Create a copy of the current data
    const key = columns[col].title as keyof Professor; // Get the column key

    // Check the type of newValue and the expected type of the column
    if (newValue.kind === GridCellKind.Text) {
      // Check if the data is of type `string` and the field is not `SubField` (which expects `string[]`)
      if (typeof newValue.data === 'string' && key !== 'SubField') {
        updatedData[row][key] = newValue.data as Professor[typeof key]; // Assign the new value safely
      }
    } else if (newValue.kind === GridCellKind.Bubble) {
      // Ensure that we are updating the `SubField` column and the data is an array
      if (Array.isArray(newValue.data) && key === 'SubField') {
        updatedData[row][key] = newValue.data as string[] as Professor[typeof key]; // Assign the array of strings safely
      }
    }

    setData(updatedData);
    setFilteredData(updatedData); // Update filtered data to reflect changes
  };

  const handleSaveOptions = () => {
    if (editingCell) {
      const { row, col } = editingCell;
      const updatedData = [...data];
      updatedData[row].SubField = [...selectedOptions];
      setData(updatedData);
      setFilteredData(updatedData);
      setIsOverlayVisible(false);
    }
  };

  const handleOptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const value = event.target.value;
    setSelectedOptions((prev) =>
      event.target.checked ? [...prev, value] : prev.filter((opt) => opt !== value)
    );
  };

  return (
    <div className="App">
      {/* Column-Specific Search Bars */}
      <div style={{ display: "flex", justifyContent: "space-evenly", padding: "10px" }}>
        {columns.map((col) => (
          <TextField
            key={col.title}
            label={`Search ${col.title}`}
            variant="outlined"
            size="small"
            value={columnFilters[col.title] || ""}
            onChange={(e) => handleColumnFilterChange(col.title, e.target.value)} // Handle changes to search values for each column
            style={{ marginBottom: "20px", width: columnWidths[col.title] }} // Adjust width for each search bar
          />
        ))}
      </div>

      {/* Data Grid Display */}
      {columns.length > 0 && filteredData.length > 0 && (
        <DataEditor
          columns={columns}
          getCellContent={getData}
          rows={filteredData.length}
          onCellEdited={onCellEdited}
          rowMarkers="number"
          onCellActivated={onCellActivated}
          showSearch={false} // Disable built-in search
          width="200vw"
          height="80vh"
        />
      )}

      {/* Modal for Editing Bubble Cells */}
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
    </div>
  );
}
