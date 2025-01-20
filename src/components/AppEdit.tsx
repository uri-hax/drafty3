import React, { useState, useEffect } from "react";
import { DataEditor, GridCellKind } from "@glideapps/glide-data-grid";
import { Button } from "@mui/material";
import Papa from "papaparse";

interface EditHistoryEntry {
    When: string;
    Edited_By: string;
    Who_Was_Edit: string;
    Column: string;
    Action: string;
    Changed_From: string;
    Changed_To: string;
}

const editHistoryData: EditHistoryEntry[] = [
    {
      When: "2024-11-24 14:58",
      Edited_By: "anon10271783",
      Who_Was_Edit: "Cindy Ya Xiong (Georgia Institute of Technology)",
      Column: "FullName",
      Action: "del row",
      Changed_From: "Cindy Ya Xiong",
      Changed_To: "[row deleted]",
    },
    {
      When: "2024-11-24 12:20",
      Edited_By: "anon10271783",
      Who_Was_Edit: "Zirui Liu (University of Minnesota)",
      Column: "University",
      Action: "edit cell",
      Changed_From: "Arizona State University",
      Changed_To: "University of Minnesota",
    },
  ];
  
const validKeys =  ["When", "Edited_By", "Who_Was_Edit", "Column", "Changed_From", "Changed_To"];
type EditKeys = typeof validKeys[number];

const AppEdit = () => {
    const [data, setData] = useState<EditHistoryEntry[]>([]);
  
    useEffect(() => {
      const fetchData = async () => {
        try {
          const response = await fetch("/edit-history.csv");
          const csvData = await response.text();
          Papa.parse(csvData, {
            header: true,
            skipEmptyLines: true,
            complete: (result) => {
              const parsedData = result.data as EditHistoryEntry[];
              setData(parsedData);
            },
          });
        } 
        catch (error) {
          console.error("Error fetching the CSV file:", error);
        }
      };
      fetchData();
    }, []);
  
    return (
        <div style={{ padding: "0px" }}>
        <h1>Edit History</h1>
        {/* Header Buttons */}
        <div style={{ marginBottom: "20px" }}>
          <button
            onClick={() => (window.location.href = "/csprofs")}
            style={{
              marginRight: "10px",
              padding: "10px 20px",
              backgroundColor: "#007bff",
              color: "white",
              border: "none",
              borderRadius: "5px",
              cursor: "pointer",
            }}
          >
            CS Professors
          </button>
          <button
            onClick={() => (window.location.href = "/edit-history")}
            style={{
              padding: "10px 20px",
              backgroundColor: "#007bff",
              color: "white",
              border: "none",
              borderRadius: "5px",
              cursor: "pointer",
            }}
          >
            Edit History
          </button>
        </div>

        {/* Table */}
        <table
          style={{
            width: "100%",
            borderCollapse: "collapse",
            textAlign: "left",
          }}
        >
          <thead>
            <tr>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>When</th>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>Edited By</th>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>Who Was Edited</th>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>Column</th>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>Action</th>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>Changed From</th>
              <th style={{ border: "1px solid #ddd", padding: "8px", backgroundColor: "#f0f0f0" }}>Changed To</th>
            </tr>
          </thead>
          <tbody>
            {data.map((row, index) => (
              <tr key={index} style={{backgroundColor: index % 2 === 1 ? "#f9f9f9" : "#fff" }} >
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.When}</td>
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.Edited_By}</td>
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.Who_Was_Edit}</td>
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.Column}</td>
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.Action}</td>
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.Changed_From}</td>
                <td style={{ border: "1px solid #ddd", padding: "8px" }}>{row.Changed_To}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  };
  
  export default AppEdit;