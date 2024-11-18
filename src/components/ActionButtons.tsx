// components/ActionButtons.tsx
import React from 'react';
import { Button } from '@mui/material';

interface ActionButtonsProps {
  handleDeleteRow: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
}

const ActionButtons: React.FC<ActionButtonsProps> = ({ handleDeleteRow, setIsAddingRow }) => (
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
);

export default ActionButtons;