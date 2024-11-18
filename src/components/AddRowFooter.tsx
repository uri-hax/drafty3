// components/AddRowFooter.tsx
import React from 'react';
import { TextField, Button, IconButton, Autocomplete } from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import type { Professor, ProfessorKey } from '../interfaces/Professor';
import { optionsList } from '../utils/constants';

interface AddRowFooterProps {
  validKeys: readonly ProfessorKey[];
  newRowData: Professor;
  setNewRowData: React.Dispatch<React.SetStateAction<Professor>>;
  optionsList: string[];
  handleAddRowConfirm: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
  allFieldsFilled: boolean;
}

const AddRowFooter: React.FC<AddRowFooterProps> = ({
  validKeys,
  newRowData,
  setNewRowData,
  optionsList,
  handleAddRowConfirm,
  setIsAddingRow,
  allFieldsFilled,
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
);

export default AddRowFooter;