import React, { useState, useEffect } from 'react';
import { Container, Grid, Typography, CircularProgress } from '@mui/material';
import MenuItemCard from '../components/MenuItemCard';
import { useNavigate } from 'react-router-dom';
import { useAuthenticatedFetch } from '../utils/api';

const Menu = () => {
  const [menuItems, setMenuItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const fetchWithAuth = useAuthenticatedFetch();

  useEffect(() => {
    let mounted = true; // For cleanup

    const fetchMenuItems = async () => {
      try {
        const response = await fetchWithAuth('/api/authenticated/v1/menu');
        if (mounted && response && response.items) {
          setMenuItems(response.items);
        }
      } catch (err) {
        if (mounted) {
          setError(err.message);
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    fetchMenuItems();

    // Cleanup function
    return () => {
      mounted = false;
    };
  }, []);

  if (loading) {
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
        <CircularProgress />
      </Container>
    );
  }

  if (error) {
    return (
      <Container sx={{ py: 4 }}>
        <Typography color="error">{error}</Typography>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Today's Menu
      </Typography>
      <Grid container spacing={3}>
        {menuItems.map((item) => (
          <Grid item xs={12} sm={6} md={4} key={item.id}>
            <MenuItemCard item={item} />
          </Grid>
        ))}
      </Grid>
    </Container>
  );
};

export default Menu;