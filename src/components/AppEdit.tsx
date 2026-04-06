import './App.css';
import "@glideapps/glide-data-grid/dist/index.css";
import React, { useState, useEffect } from 'react';
import { type GridColumn } from "@glideapps/glide-data-grid";
import { fetchCsvData } from '../utils/csvParser';
import type { ColumnData } from '../interfaces/ColumnData';
import useWindowWidth from '../hooks/useWindow';
import ActionButtons from './ActionButtons';
import DataGridWrapper from './DataGridWrapper';
import type { ColumnConfig } from '../interfaces/ColumnData';
import { editFiles, datasetLabels } from "../config/AppConfig";
import { 
  handleHomePage,
  handleData,
  handleEditHistory
} from "../utils/gridHelper";

export default function App() {
  // get the dataset we're in from the url
  const base = import.meta.env.BASE_URL;
  const dataset = window.location.pathname
    .replace(base, "")
    .split("/")
    .filter(Boolean)[0];

  // create portal div if it doesn't exist to render modals into
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
  const [columnSchema, setColumnSchema] = useState<Record<string, ColumnConfig>>({});
  const [columnFilters, setColumnFilters] = useState<Record<string, string>>({});

  useEffect(() => {
    // gather data from csv and yaml, and apply any filters from the url
    const fetchData = async () => {
      try {
        const { gridColumns, parsedData, columnSchema } = await fetchCsvData(
          gridWidth, 
          `${base}${editFiles[dataset]}`,
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
 
  // re-apply filters and sorting whenever relevant state changes
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

  // update url when filters change
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

  // get the label for the current dataset to show in the UI
  const datasetLabel = datasetLabels[dataset];

  // render the app with the action buttons at the top and the data grid taking up the rest of the space
  return ( 
    <div className="App" style={{ display: "flex", flexDirection: "column", height: "100vh" }}>
      <ActionButtons
        handleHomePage={() => handleHomePage(base)}
        handleData={() => handleData(base, dataset)}
        handleEditHistory={() => handleEditHistory(base, dataset)}
        datasetLabel={datasetLabel}
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
