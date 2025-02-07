// src/components/ActionButtons.tsx

/*
  Provides action buttons for adding and deleting rows.
  - Delete Row: Triggers a callback to remove the currently selected row(s).
  - Add Row: Opens the interface for adding a new row to the dataset.

  This component is data-agnostic and simply provides UI actions that the parent component can handle.
*/

import React from 'react';
import { Button, responsiveFontSizes } from '@mui/material';

interface ActionButtonsProps {
  handleDeleteRow: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
  handleEditHistory: () => void;
}

const ActionButtons: React.FC<ActionButtonsProps> = ({ handleDeleteRow, setIsAddingRow, handleEditHistory }) => (
  <div style={{ padding: "10px", backgroundColor: "dodgerblue" }}>
    <div style={{ display: "flex", justifyContent: "flex-start", gap: "10px" }}>
      <div style={{padding: "10px", color: "white", fontSize: 20, fontWeight: "bold"}}>
        Drafty
      </div>
      <Button 
        variant="contained"
        color="primary" 
        onClick={handleDeleteRow} 
        style={
          { marginLeft: "10px", 
            borderColor: "white", 
            borderWidth: "2px",
            borderStyle: "solid" }
        }
      >
        Delete Row
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={() => setIsAddingRow(true)}
        style={
          { marginLeft: "10px", 
            borderColor: "white", 
            borderWidth: "2px",
            borderStyle: "solid" }
        }
      >
        Add Row
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={handleEditHistory} 
        style={
          { marginLeft: "10px", 
            borderColor: "white", 
            borderWidth: "2px",
            borderStyle: "solid" }
        }
      >
          Edit History
      </Button>
    </div>
  </div>
);

export default ActionButtons;
