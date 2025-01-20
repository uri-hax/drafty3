// src/components/AddRowFooter.tsx

/*
  Displays a footer UI when adding a new row to the dataset.
  - Dynamically generates input fields for each column, based on the column schema:
    * string columns: TextField input
    * string[] columns: Autocomplete multiple-select input
  - Includes a confirm button (checks if all fields are filled) and a cancel button to discard changes.
  - Relies on parent-provided:
    * columnKeys: array of column names
    * columnSchema: defines column types ('string' or 'string[]')
    * optionsLists: map of column name -> unique string values (for Autocomplete)
    * newRowData and setNewRowData: state for building a new row before adding it
    * allFieldsFilled: boolean indicating if all inputs have values
*/

import React from 'react';
import { TextField, Button, IconButton, Autocomplete } from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import type { ColumnData } from '../interfaces/ColumnData';

interface AddRowFooterProps {
  columnKeys: string[];
  newRowData: ColumnData;
  setNewRowData: React.Dispatch<React.SetStateAction<ColumnData>>;
  optionsLists: Record<string, string[]>;
  handleAddRowConfirm: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
  allFieldsFilled: boolean;
  columnSchema: Record<string, string>; 
}

const AddRowFooter: React.FC<AddRowFooterProps> = ({
  columnKeys,
  newRowData,
  setNewRowData,
  optionsLists,
  handleAddRowConfirm,
  setIsAddingRow,
  allFieldsFilled,
  columnSchema,
}) => (
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
      {columnKeys.map((key) => {
        const colType = columnSchema[key] || 'string';

        if (colType === 'string[]') {
          const currentValues = Array.isArray(newRowData[key]) ? newRowData[key] as string[] : [];
          return (
            <Autocomplete
              key={key}
              multiple
              options={optionsLists[key] || []}
              getOptionLabel={(option) => option}
              value={currentValues}
              onChange={(event, newValue) => {
                setNewRowData((prevData) => ({
                  ...prevData,
                  [key]: newValue,
                }));
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  variant="outlined"
                  label={key}
                  placeholder={`Select ${key}`}
                  size="small"
                  style={{ margin: "5px", minWidth: "150px" }}
                />
              )}
              style={{ margin: "5px", minWidth: "150px" }}
            />
          );
        } 
        else {
          const currentValue = typeof newRowData[key] === 'string' ? (newRowData[key] as string) : '';

          return (
            <TextField
              key={key}
              label={key}
              variant="outlined"
              size="small"
              value={currentValue}
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

    <div style={{ display: "flex", alignItems: "center", marginLeft: "auto" }}>
      <IconButton
        color="primary"
        onClick={handleAddRowConfirm}
        disabled={!allFieldsFilled}
      >
        <CheckCircleIcon />
      </IconButton>

      <Button
        variant="contained"
        color="primary"
        onClick={() => {
          setIsAddingRow(false);
          setNewRowData(() => {
            const resetObj: ColumnData = {};

            for (const key of columnKeys) {
              resetObj[key] = columnSchema[key] === 'string[]' ? [] : '';
            }

            return resetObj;
          });
        }}
        
        style={{ marginLeft: "10px" }}
      >
        Cancel
      </Button>
    </div>
  </div>
);

export default AddRowFooter;
