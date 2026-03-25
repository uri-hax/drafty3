// src/components/TextInputModal.tsx

// A generic modal for editing a single string value.
// - Uses a TextField for input.
// - Triggered by a boolean state, and returns the new value to the parent.
// - Data-agnostic: can be used for any string cell, not just a specific one.

import React from "react";
import { Modal, Button, TextField, Box, Typography } from "@mui/material";

const monoFont = "ui-monospace, SFMono-Regular, Menlo, monospace";

interface FreeTextModalProps {
  isOverlayVisible: boolean;
  setIsOverlayVisible: React.Dispatch<React.SetStateAction<boolean>>;

  title?: string;
  column?: string;

  value: string;
  setValue: React.Dispatch<React.SetStateAction<string>>;

  handleSave: () => void;
}

const FreeTextModal: React.FC<FreeTextModalProps> = ({
  isOverlayVisible,
  setIsOverlayVisible,
  title = "Edit Value",
  value,
  setValue,
  handleSave,
  column,
}) => (
  <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
    <Box
      sx={{
        backgroundColor: "#ffffff",
        color: "#47494d",
        fontFamily: monoFont,
        "& *": { fontFamily: monoFont },
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
          fontFamily: monoFont,
        }}
      >
        {column}: {title}
      </Typography>

      <TextField
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder="Enter value…"
        fullWidth
        autoFocus
        sx={{
          "& .MuiInputBase-input": { fontFamily: monoFont },
        }}
      />

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

export default FreeTextModal;
