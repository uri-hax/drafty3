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
const monoFont =
  "ui-monospace, SFMono-Regular, Menlo, monospace";

// interface for the multi select modal props
interface MultiSelectModalProps {
  isOverlayVisible: boolean;
  setIsOverlayVisible: React.Dispatch<React.SetStateAction<boolean>>;

  optionsList: string[];

  selectedOptions: string[];
  setSelectedOptions: React.Dispatch<React.SetStateAction<string[]>>;

  multiple?: boolean;
  title?: string;

  handleSaveOptions: () => void;

  column?: string;
}

// component for the multi select modal - allows user to select from dropdown with multiple select and save the values on edit
const MultiSelectModal: React.FC<MultiSelectModalProps> = ({
  isOverlayVisible,
  setIsOverlayVisible,
  optionsList,
  selectedOptions,
  setSelectedOptions,
  handleSaveOptions,
  multiple = true,
  title = "Select Value(s)",
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

      {/* dropdown for selecting values from the options list and update value on change */}
      <Autocomplete
        multiple
        options={optionsList}
        value={selectedOptions}
        onChange={(_, newValue) => {
          if (multiple) {
            setSelectedOptions(newValue);
          } else {
            // single-select: clamp to one value
            setSelectedOptions(newValue.slice(0, 1));
          }
        }}
        sx={{
          "& .MuiInputBase-input": {
            fontFamily: monoFont,
          },
          "& .MuiAutocomplete-tag": {
            fontFamily: monoFont,
          },
          "& .MuiAutocomplete-listbox": {
            fontFamily: monoFont,
          },
        }}
        // show prompt conditionally based on multiple vs single select
        getOptionLabel={option => option}
        renderInput={params => (
          <TextField
            {...params}
            placeholder={multiple ? "Select values…" : "Select value…"}
          />
        )}
      />

      {/* save button to save the selected values with click to call prop handle save */}
      <Box sx={{ display: "flex", justifyContent: "flex-end", marginTop: 2 }}>
        <Button
          variant="contained"
          onClick={handleSaveOptions}
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
export default MultiSelectModal;