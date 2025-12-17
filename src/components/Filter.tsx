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
import { TextField, InputAdornment, IconButton } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import ClearIcon from '@mui/icons-material/Clear';
import type { GridColumn } from '@glideapps/glide-data-grid';

interface FilterBarProps {
  columns: GridColumn[];
  columnFilters: Record<string, string>;
  handleColumnFilterChange: (colKey: string, value: string) => void;

  // % widths used by Glide (e.g. "12.5%" or 12.5)
  columnWidths?: Record<string, string | number>;
}

const FilterBar: React.FC<FilterBarProps> = ({
  columns,
  columnFilters,
  handleColumnFilterChange,
  columnWidths = {},
}) => {
  const containerRef = React.useRef<HTMLDivElement>(null);
  const [containerWidthPx, setContainerWidthPx] = React.useState(0);

  // Measure the SAME container Glide uses
  React.useLayoutEffect(() => {
    if (!containerRef.current) return;

    const ro = new ResizeObserver(([entry]) => {
      setContainerWidthPx(entry.contentRect.width);
    });

    ro.observe(containerRef.current);
    return () => ro.disconnect();
  }, []);

  const defaultWidthPercent = 100 / columns.length;

  const percentToPx = (value: string | number, containerPx: number) => {
    const percent =
      typeof value === 'number'
        ? value
        : parseFloat(value.toString().replace('%', ''));

    return Math.round((percent / 100) * containerPx);
  };

  return (
    <div
      ref={containerRef}
      style={{
        display: 'flex',
        gap: '0px',
        overflowX: 'hidden',
        padding: 0,
        background: '#f0f0f0',
      }}
    >
      {columns.map((col) => {
        const colKey = col.id as string;
        const colWidthValue =
          columnWidths[colKey] ?? defaultWidthPercent;

        const colWidthPx =
          containerWidthPx > 0
            ? percentToPx(colWidthValue, containerWidthPx)
            : 150;

        return (
          <TextField
            key={colKey}
            variant="outlined"
            size="small"
            value={columnFilters[colKey] || ''}
            onChange={(e) =>
              handleColumnFilterChange(colKey, e.target.value)
            }

            InputProps={{
              startAdornment: !columnFilters[colKey] && (
                <InputAdornment
                  position="start"
                  sx={{
                    position: 'absolute',
                    left: 8,
                    pointerEvents: 'none',
                    opacity: 0.5,
                  }}
                >
                  <SearchIcon sx={{ fontSize: 14 }} />
                </InputAdornment>
              ),
              endAdornment: columnFilters[colKey] && (
                <InputAdornment
                  position="end"
                  sx={{
                    position: 'absolute',
                    right: 6,
                  }}
                >
                  <IconButton
                    size="small"
                    onClick={() => handleColumnFilterChange(colKey, '')}
                    sx={{
                      padding: 0,
                      fontSize: 12,
                      color: '#999',

                      '&:hover': {
                        color: '#444',
                        background: 'transparent',
                      },
                    }}
                  >
                    <ClearIcon fontSize="inherit" />
                  </IconButton>
                </InputAdornment>
              ),
            }}

            sx={{
              width: `${colWidthPx}px`,
              minWidth: `${colWidthPx}px`,
              maxWidth: `${colWidthPx}px`,
              boxSizing: 'border-box',

              // give room for the icon (NOT layout width)
              '& .MuiInputBase-input': {
                paddingLeft: '12px',
                fontSize: '12px',
              },

              // MUI internals
              '& .MuiOutlinedInput-root': {
                boxSizing: 'border-box',
                borderRadius: '4px',
                background: '#fff',
              },

              '& fieldset': {
                borderColor: '#e5e5e7',
              },
            }}
          />
        );
      })}
    </div>
  );
};

export default FilterBar;
