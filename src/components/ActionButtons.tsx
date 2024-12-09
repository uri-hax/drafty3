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
  handleDeleteRow: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
}

const ActionButtons: React.FC<ActionButtonsProps> = ({ handleDeleteRow, setIsAddingRow }) => (
  <div style={{ padding: "10px" }}>
    <div style={{ display: "flex", justifyContent: "flex-start", gap: "10px" }}>
      <Button variant="contained" color="primary" onClick={handleDeleteRow}>
        Delete Row
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={() => setIsAddingRow(true)}
        style={{ marginLeft: "10px" }}
      >
        Add Row
      </Button>
    </div>
  </div>
);

export default ActionButtons;
