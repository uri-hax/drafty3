import React from "react";
import DataEditor, { GridCellKind, type GridCell, type GridColumn, type Item } from "@glideapps/glide-data-grid";
import "@glideapps/glide-data-grid/dist/index.css";

const columns: GridColumn[] = [
  { title: "ID", id: "id" },
  { title: "Name", id: "name" },
];

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