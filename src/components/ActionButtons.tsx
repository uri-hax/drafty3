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

// custom style for primary brand button and action buttons
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

// custom style for larger primary brand button
const primaryGridButtonLargeSx: SxProps<Theme> = {
  ...primaryGridButtonSx,
  fontSize: '1.2em',
  fontWeight: 600,
  padding: '6px 14px',
};

// custom style for action buttons
const gridActionButtonSx: SxProps<Theme> = {
  ...primaryGridButtonSx,
  fontSize: '0.9em',
  padding: '4px 10px',
  minHeight: 32,

  '& .MuiButton-startIcon': {
    marginRight: '6px',
  },
};

// interface for action buttons props
interface ActionButtonsProps {
  datasetLabel: string;
  handleDeleteRow?: () => void;
  setIsDeletingRow?: React.Dispatch<React.SetStateAction<boolean>>;
  setIsAddingRow?: React.Dispatch<React.SetStateAction<boolean>>;
  handleEditHistory: () => void;
  handleData: () => void;
  handleHomePage: () => void;
}

// component for action buttons - appear at top of grid page in header
const ActionButtons: React.FC<ActionButtonsProps> = ({
  datasetLabel,
  handleDeleteRow,
  setIsDeletingRow,
  setIsAddingRow,
  handleEditHistory,
  handleData,
  handleHomePage,
}) => (
  <div style={{ padding: '0.5em', backgroundColor: '#0b89ff' }}>
    <div style={{ display: 'flex', alignItems: 'center', gap: '0.5em' }}>

      {/* Primary brand button - goes back to home page */}
      <Button
        variant="contained"
        onClick={handleHomePage}
        sx={primaryGridButtonLargeSx}
      >
        Drafty
      </Button>

      {/* Data view - go to dataset page */}
      <Tooltip title={datasetLabel} arrow>
        <Button
          onClick={handleData}
          sx={gridActionButtonSx}
          startIcon={<TableViewIcon fontSize="small" />}
        >
          {datasetLabel}
        </Button>
      </Tooltip>

      {/* Edit history - go to edit history page */}
      <Tooltip title="Edit History" arrow>
        <Button
          onClick={handleEditHistory}
          sx={gridActionButtonSx}
          startIcon={<HistoryIcon fontSize="small" />}
        >
          Edit History
        </Button>
      </Tooltip>

      {/* Add row - pull up add row modal and handle it */}
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

      {/* Delete row - pull up delete row modal and handle it */}
      {handleDeleteRow && setIsDeletingRow && (
        <Tooltip title="Delete Row" arrow>
          <Button
            onClick={() => {
              setIsDeletingRow(true);
              handleDeleteRow();
            }}
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

// export the component
export default ActionButtons;