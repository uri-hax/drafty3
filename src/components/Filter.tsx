// src/components/Filter.tsx

/*
  A filter bar that allows the user to search within each column.
  - Renders one TextField per column, using column title and widths from parent props.
  - Updates filters using a provided callback function.
  - Data-agnostic: no assumptions about specific columns.

  Requires:
  - columns: an array of GridColumn objects
  - columnFilters: a record of column key -> filter string
  - handleColumnFilterChange: updates filter state externally
  - columnWidths: a record of column key -> CSS width string
*/

import React from 'react';
import { TextField, InputAdornment } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import type { GridColumn } from '@glideapps/glide-data-grid';

interface FilterBarProps {
  columns: GridColumn[];
  columnFilters: Record<string, string>;
  handleColumnFilterChange: (colKey: string, value: string) => void;
  columnWidths?: Record<string, string>;
}

const FilterBar: React.FC<FilterBarProps> = ({
  columns,
  columnFilters,
  handleColumnFilterChange,
  columnWidths = {},
}) => {
  const defaultWidth = `${100 / columns.length}%`; 

  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'flex-start',
        padding: '0px',
        overflowX: 'auto',
      }}
    >
      {columns.map((col) => {
        const colKey = col.id as string;

        return (
          <TextField
            key={colKey}
            variant="outlined"
            size="small"
            value={columnFilters[colKey] || ''}
            onChange={(e) => handleColumnFilterChange(colKey, e.target.value)}
            style={{
              width: columnWidths?.[colKey] || defaultWidth, 
              marginRight: '0px',
            }}
            InputProps={{
              startAdornment: !columnFilters[colKey] ? (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ) : null,
            }}
          />
        );
      })}
    </div>
  );
};

export default FilterBar;