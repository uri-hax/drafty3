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
  <div style={{ padding: "10px", backgroundColor: "dodgerblue" }}>
    <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
      <div style={{ padding: "10px", color: "white", fontSize: 20, fontWeight: "bold" }}>
        Drafty
      </div>

      <Button
        variant="contained"
        color="primary"
        onClick={handleData}
        style={{
          borderColor: "white",
          borderWidth: "2px",
          borderStyle: "solid",
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
          borderWidth: "2px",
          borderStyle: "solid",
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
            borderWidth: "2px",
            borderStyle: "solid",
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
            borderWidth: "2px",
            borderStyle: "solid",
          }}
        >
          Delete Row
        </Button>
      )}
    </div>
  </div>
);

export default ActionButtons;
