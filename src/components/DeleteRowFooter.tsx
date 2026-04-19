/*
  Displays a footer UI when deleting a row from the dataset.
  - Prompts the user for a deletion reason
  - Confirm button is disabled until comment is non-empty
  - Cancel button closes the footer and resets the comment
*/

import React from "react";
import { TextField, Button, IconButton } from "@mui/material";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";

// font style for the footer
const monoFont =
  "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace";

// interface for the delete row footer props
  interface DeleteRowFooterProps {
  comment: string;
  setComment: React.Dispatch<React.SetStateAction<string>>;
  handleDeleteRowConfirm: () => void;
  onCancel: () => void;
  isCommentFilled: boolean;
}

// component for the delete row footer - prompts user for deletion reason and has confirm/cancel actions
const DeleteRowFooter: React.FC<DeleteRowFooterProps> = ({
  comment,
  setComment,
  handleDeleteRowConfirm,
  onCancel,
  isCommentFilled,
}) => (
  <div
    style={{
      position: "fixed",
      bottom: 0,
      left: 0,
      right: 0,
      backgroundColor: "white",
      padding: "10px",
      display: "flex",
      alignItems: "center",
      zIndex: 1000,
      borderTop: "1px solid #ccc",
      overflowX: "auto",
      fontFamily: monoFont,
    }}
  >
    {/* input - field for comment */}
    <div
      style={{
        display: "flex",
        flex: 1,
        alignItems: "center",
        overflowX: "auto",
      }}
    >
      <div
        style={{
          margin: "5px",
          minWidth: "200px",
          fontFamily: monoFont,
          whiteSpace: "nowrap",
        }}
      >
        Why are you deleting this? :
      </div>

      <TextField
        label="Comment"
        placeholder="Comment"
        size="small"
        value={comment}
        onChange={(e) => setComment(e.target.value)}
        sx={{
          margin: "5px",
          minWidth: "350px",
          flex: 1,
          fontFamily: monoFont,

          "& .MuiInputBase-root": {
            fontFamily: monoFont,
          },
          "& .MuiInputBase-input": {
            fontFamily: monoFont,
          },
          "& .MuiInputLabel-root": {
            fontFamily: monoFont,
          },
        }}
      />
    </div>

    {/* actions - confirm button if comment is filled and option for cancel button */}
    <div style={{ display: "flex", alignItems: "center", marginLeft: "auto" }}>
      <IconButton
        onClick={handleDeleteRowConfirm}
        disabled={!isCommentFilled}
        sx={{ fontFamily: monoFont }}
        aria-label="Confirm Delete Row"
      >
        <CheckCircleIcon />
      </IconButton>

      <Button
        variant="contained"
        onClick={() => {
          onCancel();
        }}
        sx={{
          marginLeft: "10px",
          fontFamily: monoFont,
          textTransform: "none",
        }}
      >
        Cancel
      </Button>
    </div>
  </div>
);

// export the component
export default DeleteRowFooter;