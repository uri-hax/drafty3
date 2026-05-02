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
import type { ColumnConfig, ColumnData } from '../interfaces/ColumnData';

// styling for the grid
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
  textBubbleSelected: "#ffffff",
  bgBubble: "#2a9cff", // will render ~#0b89ff glide darkens it
  bgBubbleSelected: "#2a9cff",

  fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
};

// interface for data grid wrapper props
interface DataGridWrapperProps {
  columns: GridColumn[];
  filteredData: ColumnData[];
  onCellEdited?: ([col, row]: Item, newValue: EditableGridCell | BubbleCell) => void;
  onCellEditorActivated?: (cell: Item) => void;
  gridSelection?: GridSelection;
  onGridSelectionChange?: (newSelection: GridSelection) => void;
  gridWidth: number;
  columnSchema: Record<string, ColumnConfig>;

  onHeaderSort?: (colKey: string) => void;
  sortColKey?: string | null;
  sortDir?: "asc" | "desc";
}

// component for the data grid - wraps the glide data editor and handles rendering cells based on column schema and applying custom styles
const DataGridWrapper: React.FC<DataGridWrapperProps> = ({
  columns,
  filteredData,
  onCellEdited,
  onCellEditorActivated,
  gridSelection,
  onGridSelectionChange,
  gridWidth,
  columnSchema,

  onHeaderSort,
  sortColKey = null,
  sortDir = "asc",
}) => {
  // add sort indicators to column headers if sorting is applied and return the columns with the title and indicators for sorting 
  const columnsWithSort = columns.map((c) => {
    const key = c.id as string;
    const indicator = sortColKey === key ? (sortDir === "asc" ? " ▲" : " ▼") : "";
    return { ...c, title: `${c.title}${indicator}` };
  });

  // function to get the data for a given cell - handles rendering based on column schema and applies custom styling
  const getData = ([col, row]: Item): GridCell => {
    const colKey = columnsWithSort[col].id as string;
    const cellData = filteredData[row]?.[colKey];
    const colType = columnSchema[colKey].type || 'string';

    // if the column is a string array, render as a bubble cell 
    if (colType === 'string[]') {
      const bubbleData = Array.isArray(cellData) ? cellData as string[] : [];
      return {
        kind: GridCellKind.Bubble,
        data: bubbleData,
        allowOverlay: false,
      };
    }

    // default rendering as text cell for string and other types and handle unexpected data
    const textData = (cellData !== undefined && cellData !== null) ? String(cellData) : "";
    
    return {
      kind: GridCellKind.Text,
      data: textData,
      displayData: textData,
      allowOverlay: false,
    };
  };

  // render the data editor with the appropriate props and custom styling with zebra striping and handle header clicks for sorting 
  return (
    <DataEditor
      columns={columnsWithSort}
      getCellContent={getData}
      rows={filteredData.length}
      onCellEdited={onCellEdited}
      rowMarkers="none"
      rowSelect="none"
      columnSelect="none"
      rangeSelect="cell"
      onCellActivated={onCellEditorActivated}
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

      // disable select all - not entirely working yet
      keybindings={{
        selectAll: false,
      }}

      onHeaderClicked={(col) => {
        const colKey = columns[col].id as string;
        onHeaderSort?.(colKey);
      }}
    />
  );
};

// export the component
export default DataGridWrapper;