// components/ProfessorGrid.tsx
import React from 'react';
import {
  DataEditor,
  type GridCell,
  GridCellKind,
  type GridColumn,
  type Item,
  type EditableGridCell,
  type BubbleCell,
  type GridSelection,
} from "@glideapps/glide-data-grid";
import type { Professor, ProfessorKey } from '../interfaces/Professor';

interface ProfessorGridProps {
  columns: GridColumn[];
  filteredData: Professor[];
  onCellEdited: ([col, row]: Item, newValue: EditableGridCell | BubbleCell) => void;
  onCellActivated: (cell: Item) => void;
  gridSelection: GridSelection;
  onGridSelectionChange: (newSelection: GridSelection) => void;
  gridWidth: number;
}

const ProfessorGrid: React.FC<ProfessorGridProps> = ({
  columns,
  filteredData,
  onCellEdited,
  onCellActivated,
  gridSelection,
  onGridSelectionChange,
  gridWidth,
}) => {
  // Generate data cells for the grid, including handling bubble cells for the SubField column
  const getData = ([col, row]: Item): GridCell => {
    const colKey = columns[col].id as ProfessorKey;
    const cellData = filteredData[row]?.[colKey];

    if (colKey === "SubField") {
      const bubbleData = Array.isArray(cellData) ? cellData : [];
      return {
        kind: GridCellKind.Bubble,
        data: bubbleData,
        allowOverlay: true,
      };
    }

    const textData = cellData !== undefined && cellData !== null ? String(cellData) : "";
    return {
      kind: GridCellKind.Text,
      data: textData,
      displayData: textData,
      allowOverlay: true,
    };
  };

  return (
    <DataEditor
      columns={columns}
      getCellContent={getData}
      rows={filteredData.length}
      onCellEdited={onCellEdited}
      rowMarkers="number"
      onCellActivated={onCellActivated}
      onGridSelectionChange={onGridSelectionChange}
      gridSelection={gridSelection}
      showSearch={false}
      width={gridWidth}
      height="100%"
    />
  );
};

export default ProfessorGrid;
