import React from 'react';
import {
  TextField,
  Button,
  IconButton,
  Autocomplete,
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import type { ColumnConfig, ColumnData } from '../interfaces/ColumnData';

// set fonts
const monoFont =
  'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace';

// interface for add row footer props
interface AddRowFooterProps {
  columnKeys: string[];
  newRowData: ColumnData;
  setNewRowData: React.Dispatch<React.SetStateAction<ColumnData>>;
  optionsLists: Record<string, string[]>;
  handleAddRowConfirm: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
  allFieldsFilled: boolean;
  columnSchema: Record<string, ColumnConfig>;
}

// component for add row footer - appears at bottom of page when adding a new row contains inputs for new row
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
      position: 'fixed',
      bottom: 0,
      left: 0,
      right: 0,
      backgroundColor: 'white',
      padding: '10px',
      display: 'flex',
      alignItems: 'center',
      zIndex: 1000,
      borderTop: '1px solid #ccc',
      overflowX: 'auto',
      fontFamily: monoFont,
    }}
  >
    {/* inputs - dynamically generated field inputs based on column schema */}
    <div style={{ display: 'flex', flex: 1, alignItems: 'center', overflowX: 'auto' }}>
      {columnKeys.map((key) => {
        const colType = columnSchema[key].type || 'string';

        // string[] columns - use autocomplete with multiple select
        if (colType === 'string[]') {
          const currentValues = Array.isArray(newRowData[key])
            ? (newRowData[key] as string[])
            : [];

          return (
            <Autocomplete
              key={key}
              multiple
              options={optionsLists[key] || []}
              value={currentValues}
              onChange={(_, newValue) => {
                setNewRowData((prev) => ({
                  ...prev,
                  [key]: newValue,
                }));
              }}
              sx={{
                margin: '5px',
                minWidth: '150px',
                fontFamily: monoFont,

                '& .MuiChip-root': {
                  fontFamily: monoFont,
                },
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label={key}
                  placeholder={`Select ${key}`}
                  size="small"
                  sx={{
                    fontFamily: monoFont,

                    '& .MuiInputBase-root': {
                      fontFamily: monoFont,
                    },
                    '& .MuiInputBase-input': {
                      fontFamily: monoFont,
                    },
                    '& .MuiInputLabel-root': {
                      fontFamily: monoFont,
                    },
                  }}
                />
              )}
            />
          );
        }

        // string columns - use regular text field
        const currentValue =
          typeof newRowData[key] === 'string'
            ? (newRowData[key] as string)
            : '';

        return (
          <TextField
            key={key}
            label={key}
            size="small"
            value={currentValue}
            onChange={(e) =>
              setNewRowData((prev) => ({
                ...prev,
                [key]: e.target.value,
              }))
            }
            sx={{
              margin: '5px',
              minWidth: '150px',
              fontFamily: monoFont,

              '& .MuiInputBase-root': {
                fontFamily: monoFont,
              },
              '& .MuiInputBase-input': {
                fontFamily: monoFont,
              },
              '& .MuiInputLabel-root': {
                fontFamily: monoFont,
              },
            }}
          />
        );
      })}
    </div>

    {/* actions - confirm button if all fields are filled and option for cancel button */}
    <div style={{ display: 'flex', alignItems: 'center', marginLeft: 'auto' }}>
      <Button
        variant="contained"
        onClick={handleAddRowConfirm}
        disabled={!allFieldsFilled}
        sx={{
          fontFamily: monoFont,
          textTransform: 'none',
          backgroundColor: allFieldsFilled ? '#0b89ff' : undefined,
        }}
      >
        Submit
      </Button>

      <Button
        variant="contained"
        onClick={() => {
          setIsAddingRow(false);
          setNewRowData(() => {
            const reset: ColumnData = {};
            for (const key of columnKeys) {
              reset[key] = columnSchema[key].type === 'string[]' ? [] : '';
            }
            return reset;
          });
        }}
        sx={{
          marginLeft: '10px',
          fontFamily: monoFont,
          textTransform: 'none',
        }}
      >
        Cancel
      </Button>
    </div>
  </div>
);

// export the component
export default AddRowFooter;