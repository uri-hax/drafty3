import React from 'react';
import { DataEditor, GridCellKind, CompactSelection } from "@glideapps/glide-data-grid";
import { Button } from '@mui/material';

interface EditHistoryEntry {
    When: string;
    Edited_By: string;
    Who_Was_Edit: string;
    Column: string;
    Changed_From: string;
    Changed_To: string;
}

const editHistoryData: EditHistoryEntry[] = [
    {
        When: "2021-09-01 10:00:00",
        Edited_By: "Alice",
        Who_Was_Edit: "Bob",
        Column: "FullName",
        Changed_From: "Bob",
        Changed_To: "Alice",
    },
    // Add more entries here
];

const validKeys =  ["When", "Edited_By", "Who_Was_Edit", "Column", "Changed_From", "Changed_To"];
type EditKeys = typeof validKeys[number];

const columnWidths: { [key in EditKeys]: string } = {
    When: '15%',
    Edited_By: '20%',
    Who_Was_Edit: '5%',
    Column: '20%',
    Changed_From: '20%',
    Changed_To: '20%',
};




export default AppEdit;