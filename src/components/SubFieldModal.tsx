// components/SubFieldModal.tsx
import React from 'react';
import { Modal, Button, TextField, Autocomplete } from '@mui/material';
import { optionsList } from '../utils/constants';

interface SubFieldModalProps {
  isOverlayVisible: boolean;
  setIsOverlayVisible: React.Dispatch<React.SetStateAction<boolean>>;
  selectedOptions: string[];
  setSelectedOptions: React.Dispatch<React.SetStateAction<string[]>>;
  handleSaveOptions: () => void;
}

const SubFieldModal: React.FC<SubFieldModalProps> = ({
  isOverlayVisible,
  setIsOverlayVisible,
  selectedOptions,
  setSelectedOptions,
  handleSaveOptions,
}) => (
  <Modal open={isOverlayVisible} onClose={() => setIsOverlayVisible(false)}>
    <div
      className="overlay-content"
      style={{
        background: "white",
        padding: "20px",
        borderRadius: "10px",
        margin: "50px auto",
        width: "400px",
      }}
    >
      <h3>Select SubFields</h3>
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
            placeholder="Select SubFields"
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

export default SubFieldModal;
