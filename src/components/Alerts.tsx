import React from 'react';
import MuiAlert, { type AlertProps } from '@mui/material/Alert';

// component for alerts - used for error messages and notifications
const Alert = React.forwardRef<HTMLDivElement, AlertProps>(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

// export the component
export default Alert;