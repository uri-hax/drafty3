// components/Alerts.tsx
import React from 'react';
import MuiAlert, { type AlertProps } from '@mui/material/Alert';

// Custom alert component for displaying messages in the snackbar
const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export default Alert;