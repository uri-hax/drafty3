// src/components/Alerts.tsx

/*
  Custom alert component for displaying messages within a snackbar or similar UI.
  - Uses Material UI's Alert with a filled variant and defined ref forwarding.
  - Data-agnostic, simply provides a styled alert component.
*/

import React from 'react';
import MuiAlert, { type AlertProps } from '@mui/material/Alert';

const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export default Alert;
