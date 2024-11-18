// components/Filter.tsx
import React from 'react';
import { TextField } from '@mui/material';
import type { GridColumn } from "@glideapps/glide-data-grid";
import type { ProfessorKey } from '../interfaces/Professor';

interface FilterBarProps {
  columns: GridColumn[];
  columnFilters: { [key in ProfessorKey]?: string };
  handleColumnFilterChange: (colKey: ProfessorKey, value: string) => void;
  columnWidths: { [key in ProfessorKey]: string };
}

const FilterBar: React.FC<FilterBarProps> = ({
  columns,
  columnFilters,
  handleColumnFilterChange,
  columnWidths,
}) => (
  <div
    style={{
      display: "flex",
      justifyContent: "space-evenly",
      padding: "10px",
      flexWrap: "wrap",
    }}
  >
    {columns.map((col) => {
      const colKey = col.id as ProfessorKey;
      return (
        <TextField
          key={colKey}
          label={`Search ${col.title}`}
          variant="outlined"
          size="small"
          value={columnFilters[colKey] || ""}
          onChange={(e) => handleColumnFilterChange(colKey, e.target.value)}
          style={{ marginBottom: "20px", width: columnWidths[colKey] }}
        />
      );
    })}
  </div>
);

export default FilterBar;