// src/components/MultiSelectModal.tsx

/*
  A generic modal for selecting multiple string values from a given options list.
  - Uses Autocomplete with multiple selection.
  - Triggered by a boolean state, and returns selected options to the parent.
  - Data-agnostic: can be used for any string[] column, not just a specific one.

  Requires:
  - isOverlayVisible & setIsOverlayVisible: controls modal visibility
  - selectedOptions & setSelectedOptions: current selection state
  - handleSaveOptions: callback to apply the changes
  - optionsList: array of possible string values
  - title: optional modal title for clarity
*/
import React from "react";
import {
  Modal,
  Button,
  TextField,
  Autocomplete,
  Box,
  Typography,
} from "@mui/material";

const monoFont =
  "ui-monospace, SFMono-Regular, Menlo, monospace";

interface MultiSelectModalProps {
  isOverlayVisible: boolean;
  setIsOverlayVisible: React.Dispatch<React.SetStateAction<boolean>>;

  optionsList: string[];

  selectedOptions: string[];
  setSelectedOptions: React.Dispatch<React.SetStateAction<string[]>>;

  multiple?: boolean;
  title?: string;

  handleSaveOptions: () => void;
}

const MultiSelectModal: React.FC<MultiSelectModalProps> = ({
  isOverlayVisible,
  setIsOverlayVisible,
  optionsList,
  selectedOptions,
  setSelectedOptions,
  handleSaveOptions,
  multiple = true,
  title = "Select Value(s)",
}) => (
  <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
    <Box
      sx={{
        backgroundColor: "#ffffff",
        color: "#47494d",
        fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
        width: 420,
        margin: "64px auto",
        padding: 2,
        borderRadius: "6px",
        border: "1px solid #e5e7eb",
        boxShadow: "0 8px 24px rgba(0,0,0,0.08)",
      }}
    >
      <Typography
        sx={{
          fontSize: 14,
          fontWeight: 600,
          marginBottom: 1.5,
        }}
      >
        {title}
      </Typography>

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
        getOptionLabel={option => option}
        renderInput={params => (
          <TextField
            {...params}
            placeholder={multiple ? "Select values…" : "Select value…"}
          />
        )}
      />

      <Box sx={{ display: "flex", justifyContent: "flex-end", marginTop: 2 }}>
        <Button
          variant="contained"
          onClick={handleSaveOptions}
          sx={{
            fontFamily: "ui-monospace, SFMono-Regular, Menlo, monospace",
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


export default MultiSelectModal;
