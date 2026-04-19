import React from 'react';
import { TextField, InputAdornment, IconButton } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import ClearIcon from '@mui/icons-material/Clear';
import type { GridColumn } from '@glideapps/glide-data-grid';

// interface for the filter bar props
interface FilterBarProps {
  columns: GridColumn[];
  columnFilters: Record<string, string>;
  handleColumnFilterChange: (colKey: string, value: string) => void;

  // % widths used by Glide (e.g. "12.5%" or 12.5)
  columnWidths?: Record<string, string | number>;
}

// component for the filter bar - renders a text input for each column to allow user to filter rows by column values
const FilterBar: React.FC<FilterBarProps> = ({
  columns,
  columnFilters,
  handleColumnFilterChange,
  columnWidths = {},
}) => {
  // ref and state to track container width for calculating column widths in px based on % widths from props
  const containerRef = React.useRef<HTMLDivElement>(null);
  const [containerWidthPx, setContainerWidthPx] = React.useState(0);

  // use layout effect to track container width changes and update state accordingly
  React.useLayoutEffect(() => {
    if (!containerRef.current) return;

    const ro = new ResizeObserver(([entry]) => {
      setContainerWidthPx(entry.contentRect.width);
    });

    ro.observe(containerRef.current);
    return () => ro.disconnect();
  }, []);

  // default width for columns without specified width in props
  const defaultWidthPercent = 100 / columns.length;

  // helper function to convert % width to px based on container width
  const percentToPx = (value: string | number, containerPx: number) => {
    const percent =
      typeof value === 'number'
        ? value
        : parseFloat(value.toString().replace('%', ''));

    return Math.round((percent / 100) * containerPx);
  };

  // return the component - a row of text inputs for each column
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
      {/* map over columns to render a text input for each column with search and clear icons and handle changes to update filters in parent component */}
      {columns.map((col) => {
        // get the column key and calculate the width in px based on props or default
        const colKey = col.id as string;
        const colWidthValue =
          columnWidths[colKey] ?? defaultWidthPercent;

        // calculate the width in px based on container width and column width value
        const colWidthPx =
          containerWidthPx > 0
            ? percentToPx(colWidthValue, containerWidthPx)
            : 150;

        // render the text input with adornments for search and clear icons and styles to set width and give room for icons without affecting layout width
        return (
          <TextField
            key={colKey}
            variant="outlined"
            size="small"
            value={columnFilters[colKey] || ''}
            // handle changes to update filters in parent component
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

            // styles to set the width of the input box based on calculated column width and give room for icons without affecting layout width
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

// export the component
export default FilterBar;