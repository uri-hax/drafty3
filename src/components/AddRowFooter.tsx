/*
  Displays a footer UI when adding a new row to the dataset.
  - Dynamically generates input fields for each column, based on the column schema:
    * string columns: TextField input
    * string[] columns: Autocomplete multiple-select input
  - Includes a confirm button (checks if all fields are filled) and a cancel button to discard changes.
  - Relies on parent-provided:
    * columnKeys: array of column names
    * columnSchema: defines column types ('string' or 'string[]')
    * optionsLists: map of column name -> unique string values (for Autocomplete)
    * newRowData and setNewRowData: state for building a new row before adding it
    * allFieldsFilled: boolean indicating if all inputs have values
*/

import React from 'react';
import {
  TextField,
  Button,
  IconButton,
  Autocomplete,
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import type { ColumnData } from '../interfaces/ColumnData';

const monoFont =
  'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace';

interface AddRowFooterProps {
  columnKeys: string[];
  newRowData: ColumnData;
  setNewRowData: React.Dispatch<React.SetStateAction<ColumnData>>;
  optionsLists: Record<string, string[]>;
  handleAddRowConfirm: () => void;
  setIsAddingRow: React.Dispatch<React.SetStateAction<boolean>>;
  allFieldsFilled: boolean;
  columnSchema: Record<string, string>;
}

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
    {/* Inputs */}
    <div style={{ display: 'flex', flex: 1, alignItems: 'center', overflowX: 'auto' }}>
      {columnKeys.map((key) => {
        const colType = columnSchema[key] || 'string';

        // --- string[] columns ---
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

        // --- string columns ---
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

    {/* Actions */}
    <div style={{ display: 'flex', alignItems: 'center', marginLeft: 'auto' }}>
      <IconButton
        onClick={handleAddRowConfirm}
        disabled={!allFieldsFilled}
        sx={{ fontFamily: monoFont }}
        aria-label="Confirm Add Row"
      >
        <CheckCircleIcon />
      </IconButton>

      <Button
        variant="contained"
        onClick={() => {
          setIsAddingRow(false);
          setNewRowData(() => {
            const reset: ColumnData = {};
            for (const key of columnKeys) {
              reset[key] = columnSchema[key] === 'string[]' ? [] : '';
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

export default AddRowFooter;
