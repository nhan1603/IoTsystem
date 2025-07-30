import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  Typography,
  Box,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  CircularProgress,
  Chip
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { useAuthenticatedFetch } from '../utils/api';

const Orders = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const fetchWithAuth = useAuthenticatedFetch();

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const response = await fetchWithAuth('/api/authenticated/v1/orders');
        if (response.success) {
          setOrders(response.data);
        } else {
          throw new Error('Failed to fetch orders');
        }
      } catch (err) {
        setError('Failed to load orders. Please try again later.');
        console.error('Error fetching orders:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [fetchWithAuth]);

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case 'paid':
        return 'success';
      case 'pending':
        return 'warning';
      case 'cancelled':
        return 'error';
      default:
        return 'default';
    }
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="md" sx={{ py: 4 }}>
      <Typography variant="h4" gutterBottom>
        My Orders
      </Typography>

      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      {orders.length === 0 ? (
        <Paper sx={{ p: 3, textAlign: 'center' }}>
          <Typography>No orders found.</Typography>
        </Paper>
      ) : (
        orders.map((order) => (
          <Accordion key={order.id} sx={{ mb: 2 }}>
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%', alignItems: 'center' }}>
                <Typography>
                  Order #{order.id} - {formatDate(order.created_at)}
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <Chip
                    label={order.status}
                    color={getStatusColor(order.status)}
                    size="small"
                  />
                  <Typography>
                    £{order.total_amount.toFixed(2)}
                  </Typography>
                </Box>
              </Box>
            </AccordionSummary>
            <AccordionDetails>
              {order.items.map((item) => (
                <Box
                  key={item.id}
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    py: 1,
                    borderBottom: '1px solid #eee'
                  }}
                >
                  <Box>
                    <Typography>{item.name}</Typography>
                    <Typography variant="body2" color="text.secondary">
                      Quantity: {item.quantity} × £{item.unit_price.toFixed(2)}
                    </Typography>
                  </Box>
                  <Typography>
                    £{item.subtotal.toFixed(2)}
                  </Typography>
                </Box>
              ))}
            </AccordionDetails>
          </Accordion>
        ))
      )}
    </Container>
  );
};

export default Orders; 