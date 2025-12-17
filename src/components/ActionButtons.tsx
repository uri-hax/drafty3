import React from 'react';
import {
  Button,
  Tooltip,
  type SxProps,
  type Theme,
} from '@mui/material';

import HistoryIcon from '@mui/icons-material/History';
import TableViewIcon from '@mui/icons-material/TableView';
import AddBoxRoundedIcon from '@mui/icons-material/AddBoxRounded';
import IndeterminateCheckBoxRoundedIcon from
  '@mui/icons-material/IndeterminateCheckBoxRounded';

const primaryGridButtonSx: SxProps<Theme> = {
  backgroundColor: '#0b89ff',
  border: '1px solid #0b89ff',
  boxShadow: 'none',
  textTransform: 'none',
  fontFamily: 'monospace',
  color: '#fff',

  '&:hover': {
    backgroundColor: '#0a7be6',
    boxShadow: 'none',
  },
};

const primaryGridButtonLargeSx: SxProps<Theme> = {
  ...primaryGridButtonSx,
  fontSize: '1.2em',
  fontWeight: 600,
  padding: '6px 14px',
};

const gridActionButtonSx: SxProps<Theme> = {
  ...primaryGridButtonSx,
  fontSize: '0.9em',
  padding: '4px 10px',
  minHeight: 32,

  '& .MuiButton-startIcon': {
    marginRight: '6px',
  },
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
  <div style={{ padding: '0.5em', backgroundColor: '#0b89ff' }}>
    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5em' }}>

      {/* Primary brand button */}
      <Button
        variant="contained"
        onClick={handleHomePage}
        sx={primaryGridButtonLargeSx}
      >
        Drafty
      </Button>

      {/* Data view */}
      <Tooltip title="CS Professors" arrow>
        <Button
          onClick={handleData}
          sx={gridActionButtonSx}
          startIcon={<TableViewIcon fontSize="small" />}
        >
          CS Professors
        </Button>
      </Tooltip>

      {/* Edit history */}
      <Tooltip title="Edit History" arrow>
        <Button
          onClick={handleEditHistory}
          sx={gridActionButtonSx}
          startIcon={<HistoryIcon fontSize="small" />}
        >
          Edit History
        </Button>
      </Tooltip>

      {/* Add row */}
      {setIsAddingRow && (
        <Tooltip title="Add Row" arrow>
          <Button
            onClick={() => setIsAddingRow(true)}
            sx={gridActionButtonSx}
            startIcon={<AddBoxRoundedIcon fontSize="small" />}
          >
            Add Row
          </Button>
        </Tooltip>
      )}

      {/* Delete row */}
      {handleDeleteRow && (
        <Tooltip title="Delete Row" arrow>
          <Button
            onClick={handleDeleteRow}
            sx={gridActionButtonSx}
            startIcon={
              <IndeterminateCheckBoxRoundedIcon fontSize="small" />
            }
          >
            Delete Row
          </Button>
        </Tooltip>
      )}
    </div>
  </div>
);

export default ActionButtons;
