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

import React from 'react';
import { Modal, Button, TextField, Autocomplete } from '@mui/material';

interface MultiSelectModalProps {
  isOverlayVisible: boolean;
  setIsOverlayVisible: React.Dispatch<React.SetStateAction<boolean>>;
  selectedOptions: string[];
  setSelectedOptions: React.Dispatch<React.SetStateAction<string[]>>;
  handleSaveOptions: () => void;
  optionsList: string[];
  title?: string;
}

const MultiSelectModal: React.FC<MultiSelectModalProps> = ({
  isOverlayVisible,
  setIsOverlayVisible,
  selectedOptions,
  setSelectedOptions,
  handleSaveOptions,
  optionsList,
  title = "Select Values",
}) => (
  <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
    <div
      className="overlay-content"
      style={{
        background: "white",
        padding: "20px",
        borderRadius: "8px",
        margin: "50px auto",
        width: "400px",
      }}
    >
      <h3>{title}</h3>
      <Autocomplete
        multiple
        options={optionsList}
        getOptionLabel={(option) => option}
        value={selectedOptions}
        onChange={(event, newValue) => {
          setSelectedOptions(newValue);
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            variant="outlined"
            placeholder="Select Values"
          />
        )}
        style={{ marginBottom: "20px" }}
      />
      <Button variant="contained" color="primary" onClick={handleSaveOptions}>
        Save
      </Button>
    </div>
  </Modal>
);

export default MultiSelectModal;
