import React from 'react';
import {
  Container,
  Paper,
  Typography,
  Button,
  Box
} from '@mui/material';
import { CheckCircle } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const OrderSuccess = () => {
  const navigate = useNavigate();

  return (
    <Container maxWidth="sm" sx={{ py: 8 }}>
      <Paper elevation={3} sx={{ p: 4, textAlign: 'center' }}>
        <CheckCircle sx={{ fontSize: 64, color: 'success.main', mb: 2 }} />
        <Typography variant="h5" gutterBottom>
          Order Placed Successfully!
        </Typography>
        <Typography color="text.secondary" sx={{ mb: 4 }}>
          Thank you for your order. You will receive a confirmation email shortly.
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', gap: 2 }}>
          <Button
            variant="contained"
            onClick={() => navigate('/')}
          >
            Return to Menu
          </Button>
          <Button
            variant="outlined"
            onClick={() => navigate('/orders')}
          >
            View Orders
          </Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default OrderSuccess; 