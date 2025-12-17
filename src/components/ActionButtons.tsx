// src/components/ActionButtons.tsx

/*
  Provides action buttons for adding and deleting rows.
  - Delete Row: Triggers a callback to remove the currently selected row(s).
  - Add Row: Opens the interface for adding a new row to the dataset.

  This component is data-agnostic and simply provides UI actions that the parent component can handle.
*/

import React from 'react';
import { Button, type SxProps, type Theme } from '@mui/material';

const primaryGridButtonSx: SxProps<Theme> = {
  borderColor: '#13599f',
  borderWidth: 1,
  borderStyle: 'solid',
  fontSize: '0.75em',
  boxShadow: 'none',
  textTransform: 'none',
};

interface ActionButtonsProps {
  handleDeleteRow?: () => void;
  setIsAddingRow?: React.Dispatch<React.SetStateAction<boolean>>;
  handleEditHistory: () => void;
  handleData: () => void;
  handleHomePage: () => void;
}

const ActionButtons: React.FC<ActionButtonsProps> = ({
  handleDeleteRow,
  setIsAddingRow,
  handleEditHistory,
  handleData,
  handleHomePage,
}) => (
  <div style={{ padding: "0.5em", backgroundColor: "#1976D2" }}>
    <div style={{ display: "flex", alignItems: "center", gap: "0.5em" }}>

      <Button
        variant="contained"
        color="primary"
        onClick={handleHomePage}
        sx={primaryGridButtonSx}
      >
        Drafty
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={handleData}
        sx={primaryGridButtonSx}
      >
        CS Professors
      </Button>

      <Button
        variant="contained"
        color="primary"
        onClick={handleEditHistory}
        sx={primaryGridButtonSx}
      >
        Edit History
      </Button>

      {setIsAddingRow && (
        <Button
          variant="contained"
          color="primary"
          onClick={() => setIsAddingRow(true)}
          sx={primaryGridButtonSx}
        >
          Add Row
        </Button>
      )}

      {handleDeleteRow && (
        <Button
          variant="contained"
          color="primary"
          onClick={handleDeleteRow}
          sx={primaryGridButtonSx}
        >
          Delete Row
        </Button>
      )}
    </div>
  </div>
);

export default ActionButtons;
