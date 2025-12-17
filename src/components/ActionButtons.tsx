// src/components/ActionButtons.tsx

/*
  Provides action buttons for adding and deleting rows.
  - Delete Row: Triggers a callback to remove the currently selected row(s).
  - Add Row: Opens the interface for adding a new row to the dataset.

  This component is data-agnostic and simply provides UI actions that the parent component can handle.
*/

import React from 'react';
import { Button } from '@mui/material';

interface ActionButtonsProps {
  handleDeleteRow?: () => void;
  setIsAddingRow?: React.Dispatch<React.SetStateAction<boolean>>;
  handleEditHistory: () => void;
  handleData: () => void;
}

const ActionButtons: React.FC<ActionButtonsProps> = ({
  handleDeleteRow,
  setIsAddingRow,
  handleEditHistory,
  handleData,
}) => (
  <div style={{ padding: "0.5em", backgroundColor: "#1976D2" }}>
    <div style={{ display: "flex", alignItems: "center", gap: "0.5em" }}>

      <Button
        variant="contained"
        color="primary"
        onClick={handleData}
        style={{
          borderColor: "white",
          borderWidth: "1px",
          borderStyle: "solid",
          fontSize: "0.75em"
        }}
      >
        Drafty
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={handleData}
        style={{
          borderColor: "white",
          borderWidth: "1px",
          borderStyle: "solid",
          fontSize: "0.75em"
        }}
      >
        CS Professors
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={handleEditHistory}
        style={{
          borderColor: "white",
          borderWidth: "1px",
          borderStyle: "solid",
          fontSize: "0.75em"
        }}
      >
        Edit History
      </Button>

      {setIsAddingRow && (
        <Button
          variant="contained"
          color="primary"
          onClick={() => setIsAddingRow(true)}
          style={{
            borderColor: "white",
            borderWidth: "1px",
            borderStyle: "solid",
            fontSize: "0.75em"
          }}
        >
          Add Row
        </Button>
      )}

      {handleDeleteRow && (
        <Button
          variant="contained"
          color="primary"
          onClick={handleDeleteRow}
          style={{
            borderColor: "white",
            borderWidth: "1px",
            borderStyle: "solid",
            fontSize: "0.75em"
          }}
        >
          Delete Row
        </Button>
      )}
    </div>
  </div>
);

export default ActionButtons;
