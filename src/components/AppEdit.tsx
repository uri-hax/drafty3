import './App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import { type GridColumn } from "@glideapps/glide-data-grid";
import { fetchCsvData } from '../utils/csvParser';
import type { ColumnData } from '../interfaces/ColumnData';
import useWindowWidth from '../hooks/useWindow';
import ActionButtons from './ActionButtons';
import DataGridWrapper from './DataGridWrapper';

const customWidths: Record<string, string> = {
  When: "10%",
  EditedBy: "10%",
  Action: "10%",
  WhoWasEdited: "20%",
  Column: "10%",
  ChangedFrom: "20%",
  ChangedTo: "20%"
};

export default function App() {
  useEffect(() => {
    const portalDiv = document.getElementById('portal');

    if (!portalDiv) {
      const newPortalDiv = document.createElement('div');
      newPortalDiv.id = 'portal';
      document.body.appendChild(newPortalDiv);
    }
  }, []);

  const gridWidth = useWindowWidth();

  const [columns, setColumns] = useState<GridColumn[]>([]);
  const [data, setData] = useState<ColumnData[]>([]);
  const [filteredData, setFilteredData] = useState<ColumnData[]>([]);
  const [columnSchema, setColumnSchema] = useState<Record<string, string>>({});
  const [columnFilters, setColumnFilters] = useState<Record<string, string>>({});

  useEffect(() => {
    const fetchData = async () => {
      try {
        const { gridColumns, parsedData, columnSchema } = await fetchCsvData(
          gridWidth, 
          customWidths,
          '/drafty3/edit-history.csv',
        ); 

        setColumns(gridColumns);
        setData(parsedData);
        setFilteredData(parsedData);
        setColumnSchema(columnSchema);

        const params = new URLSearchParams(window.location.search);
        const initialFilters: Record<string, string> = {};
        const columnKeys = Object.keys(parsedData[0] || {});

        columnKeys.forEach((name) => {
          const value = params.get(name);
          if (value) {
            initialFilters[name] = value;
          }
        });

        setColumnFilters(initialFilters);
      } 
      catch (error) {
        console.error('Error fetching CSV/YAML data:', error);
      }
    };

    fetchData();
  }, [gridWidth]);
 
  useEffect(() => {
    if (columns.length === 0 || data.length === 0) return;

    const filtered = data.filter((row) =>
      columns.every((col) => {
        const colKey = col.id as string;
        const filterValue = columnFilters[colKey];

        if (!filterValue) return true;

        const cellValue = row[colKey];

        if (Array.isArray(cellValue)) {
          return cellValue.some((val) =>
            val.toString().toLowerCase().includes(filterValue.toLowerCase())
          );
        }

        return cellValue?.toString().toLowerCase().includes(filterValue.toLowerCase());
      })
    );

    setFilteredData(filtered);
  }, [columnFilters, data, columns]);

  useEffect(() => {
    if (columns.length === 0) return;

    const params = new URLSearchParams();

    for (const [key, value] of Object.entries(columnFilters)) {
      if (value) {
        params.set(key, value);
      }
    }

    window.history.replaceState(null, '', '?' + params.toString());
  }, [columnFilters, columns]);

  const handleHomePage = () => {
    window.location.href = "/drafty3/";
  }

  const handleData = () => {
    window.location.href = "/drafty3/csprofs";
  }

  const handleEditHistory = () => {
    window.location.href = "/drafty3/edit-history";
  };

  return (
    <div className="App" style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      <ActionButtons
        handleHomePage={handleHomePage}
        handleData={handleData}
        handleEditHistory={handleEditHistory}
      />

      <div className="grid-container" style={{ flexGrow: 1 }}>
        <DataGridWrapper
          columns={columns}
          filteredData={filteredData}
          gridWidth={gridWidth}
          columnSchema={columnSchema}
        />
      </div>
    </div>
  );
}
