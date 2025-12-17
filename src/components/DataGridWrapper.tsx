// src/components/DataGridWrapper.tsx

/*
  A wrapper around the DataEditor from @glideapps/glide-data-grid that:
  - Receives dynamically generated columns and filtered data rows.
  - Maps cell data to either Text or Bubble cells based on column types from columnSchema.
  - Handles edits and selections through parent callbacks.
  - Data-agnostic: No column assumptions.

  Requires:
  - columns, filteredData, columnSchema: to know how to render each cell
  - onCellEdited, onCellActivated, onGridSelectionChange: parent event handlers
  - gridSelection: current selection state
  - gridWidth: width of the grid in pixels
*/

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
import type { ColumnData } from '../interfaces/ColumnData';

const draftyOld = {
  //bgCell: "#0f172a",
  //bgCellHeader: "#020617",
  //bgCellHovered: "#1e293b",
  //bgCellSelected: "#334155",

  textDark: "#47494d",
  //textMedium: "#cbd5f5",
  //textLight: "#94a3b8",

  //borderColor: "#f0f0f0",
  accentColor: "#2a9cff",
  accentFg: "#18609eff",

  scrollbarThumb: "#475569",
  scrollbarTrack: "#020617",

  textBubble: "#ffffff",
  bgBubble: "#2a9cff", // will render ~#0b89ff glide darkens it
  bgBubbleSelected: "#2a9cff",

  fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
};


interface DataGridWrapperProps {
  columns: GridColumn[];
  filteredData: ColumnData[];
  onCellEdited?: ([col, row]: Item, newValue: EditableGridCell | BubbleCell) => void;
  onCellActivated?: (cell: Item) => void;
  gridSelection?: GridSelection;
  onGridSelectionChange?: (newSelection: GridSelection) => void;
  gridWidth: number;
  columnSchema: Record<string, string>;
}

const DataGridWrapper: React.FC<DataGridWrapperProps> = ({
  columns,
  filteredData,
  onCellEdited,
  onCellActivated,
  gridSelection,
  onGridSelectionChange,
  gridWidth,
  columnSchema,
}) => {
  const editable = Boolean(onCellEdited);

  const getData = ([col, row]: Item): GridCell => {
    const colKey = columns[col].id as string;
    const cellData = filteredData[row]?.[colKey];
    const colType = columnSchema[colKey] || 'string';

    if (colType === 'string[]') {
      const bubbleData = Array.isArray(cellData) ? cellData as string[] : [];
      return {
        kind: GridCellKind.Bubble,
        data: bubbleData,
        allowOverlay: editable,
      };
    }

    const textData = (cellData !== undefined && cellData !== null) ? String(cellData) : "";
    
    return {
      kind: GridCellKind.Text,
      data: textData,
      displayData: textData,
      allowOverlay: editable,
    };
  };

  return (
    <DataEditor
      columns={columns}
      getCellContent={getData}
      rows={filteredData.length}
      onCellEdited={onCellEdited}
      rowMarkers="none"
      rowSelect="none"
      columnSelect="none"
      rangeSelect="cell"
      onCellActivated={onCellActivated}
      onGridSelectionChange={onGridSelectionChange}
      gridSelection={gridSelection}
      showSearch={false}
      width={gridWidth}
      height="100%"

      theme={draftyOld}

      // glide table hack for zebra striping
      getRowThemeOverride={(rowIndex: number) => ({
        bgCell: rowIndex % 2 === 0 ? "#ffffff" : "#f7f7f7", 
      })}
    />
  );
};

export default DataGridWrapper;
