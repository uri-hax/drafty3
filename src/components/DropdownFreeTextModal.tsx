import React from "react";
import {
  Modal,
  Button,
  TextField,
  Autocomplete,
  Box,
  Typography,
} from "@mui/material";

// set font for the modal
const monoFont = "ui-monospace, SFMono-Regular, Menlo, monospace";

// interface for the dropdown free text modal props
interface DropdownFreeTextModalProps {
  isOverlayVisible: boolean;
  setIsOverlayVisible: React.Dispatch<React.SetStateAction<boolean>>;

  optionsList: string[];

  title?: string;
  column?: string;

  dropdownValue: string;
  setDropdownValue: React.Dispatch<React.SetStateAction<string>>;

  draftValue: string;
  setDraftValue: React.Dispatch<React.SetStateAction<string>>;

  handleSave: () => void;
}

// component for the dropdown free text modal - allows user to select from dropdown or enter free text and save the value on edit
const DropdownFreeTextModal: React.FC<DropdownFreeTextModalProps> = ({
  isOverlayVisible,
  setIsOverlayVisible,
  optionsList,
  title = "Select Value",
  dropdownValue,
  setDropdownValue,
  draftValue,
  setDraftValue,
  handleSave,
  column,
}) => (
  <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
    <Box
      sx={{
        backgroundColor: "#ffffff",
        color: "#47494d",
        fontFamily: monoFont,
        "& *": {
          fontFamily: monoFont,
        },
        width: 420,
        margin: "64px auto",
        padding: 2,
        borderRadius: "6px",
        border: "1px solid #e5e7eb",
        boxShadow: "0 8px 24px rgba(0,0,0,0.08)",
      }}
    >
      {/* title with column name to display */}
      <Typography
        sx={{
          fontSize: 14,
          fontWeight: 600,
          marginBottom: 1.5,
          fontFamily: monoFont,
        }}
      >
        {column}: {title}
      </Typography>

      {/* dropdown for selecting values from the options list and update dropdown value */}
      <Autocomplete
        options={optionsList}
        value={dropdownValue || null}
        onChange={(_, newValue) => setDropdownValue(newValue ?? "")}
        sx={{
          "& .MuiInputBase-input": { fontFamily: monoFont },
          "& .MuiAutocomplete-listbox": { fontFamily: monoFont },
        }}
        getOptionLabel={(option) => option}
        renderInput={(params) => (
          <TextField {...params} placeholder="Select value…" />
        )}
      />

      {/* prompt to enter free text if necessary */}
      <Typography
        sx={{
          fontSize: 13,
          fontWeight: 600,
          marginTop: 1.75,
          marginBottom: 1,
          fontFamily: monoFont,
        }}
      >
        Not here? Type value here:
      </Typography>

      {/* text field for entering free text value with the current draft value and update on change */}
      <TextField
        value={draftValue}
        onChange={(e) => setDraftValue(e.target.value)}
        placeholder="Enter value…"
        fullWidth
        sx={{
          "& .MuiInputBase-input": { fontFamily: monoFont },
        }}
      />

      {/* save button to save the selected or entered value with click to call prop handle save */}
      <Box sx={{ display: "flex", justifyContent: "flex-end", marginTop: 2 }}>
        <Button
          variant="contained"
          onClick={handleSave}
          sx={{
            fontFamily: monoFont,
            textTransform: "none",
            fontSize: 13,
            backgroundColor: "#2a9cff",
            boxShadow: "none",
            "&:hover": { backgroundColor: "#1f86e6" },
          }}
        >
          Save
        </Button>
      </Box>
    </Box>
  </Modal>
);

// export the component
export default DropdownFreeTextModal;