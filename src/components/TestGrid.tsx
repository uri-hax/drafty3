import React from "react";
import DataEditor, { GridCellKind, type GridCell, type GridColumn, type Item } from "@glideapps/glide-data-grid";
import "@glideapps/glide-data-grid/dist/index.css";

// test columns
const columns: GridColumn[] = [
  { title: "ID", id: "id" },
  { title: "Name", id: "name" },
];

// test content 
function getCellContent([col, row]: Item): GridCell {
  if (col === 0) {
    return {
      kind: GridCellKind.Text,
      allowOverlay: true,
      displayData: String(row),
      data: String(row),
    };
  }
  return {
    kind: GridCellKind.Text,
    allowOverlay: true,
    displayData: `Row ${row}`,
    data: `Row ${row}`,
  };
}

// test component to test the glide data grid with basic columns and content - can be used for testing features in isolation before integrating into the main data grid component
export default function TestGrid() {
  return (
    <div style={{ padding: 16 }}>
      <h3 style={{ margin: 0, marginBottom: 12 }}>Glide Minimal Test</h3>
      <div style={{ border: "1px solid #ccc", display: "inline-block" }}>
        <DataEditor
          width={800}
          height={500}
          rows={100}
          columns={columns}
          getCellContent={getCellContent}
        />
      </div>
    </div>
  );
}