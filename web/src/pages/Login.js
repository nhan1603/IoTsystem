import React, { useState, useEffect } from 'react';
import {
  Container,
  Paper,
  TextField,
  Button,
  Box,
  Alert,
  Tab,
  Tabs,
  CircularProgress
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useCart } from '../contexts/CartContext';

const Login = () => {
  const navigate = useNavigate();
  const { clearCart } = useCart();
  const [tab, setTab] = useState(0); // 0 for login, 1 for register
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true); // Start with loading true
  const [loginData, setLoginData] = useState({
    email: '',
    password: ''
  });
  const [registerData, setRegisterData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  });

  // Check token validity on component mount
  useEffect(() => {
    const verifyToken = async () => {
      const token = localStorage.getItem('token');
      if (!token) {
        setLoading(false);
        return;
      }

      try {
        // Make a request to verify the token
        const response = await fetch('/api/authenticated/v1/menu', {
          headers: {
            'Authorization': token,
          },
        });

        if (response.ok) {
          // Token is valid, redirect to home
          navigate('/');
        } else {
          // Token is invalid, remove it
          localStorage.removeItem('token');
          setLoading(false);
        }
      } catch (err) {
        // Error occurred, remove token and show login form
        localStorage.removeItem('token');
        setLoading(false);
      }
    };

    verifyToken();
  }, [navigate]);

  // Show loading spinner while verifying token
  if (loading) {
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <CircularProgress />
      </Container>
    );
  }

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await fetch('/api/public/v1/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: loginData.email,
          password: loginData.password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Login failed');
      }

      // Store the token
      localStorage.setItem('token', data.token);
      
      // Clear the cart when logging in
      clearCart();
      
      // Redirect to menu page
      navigate('/');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const validatePassword = (password) => {
    const minLength = 6;
    const maxLength = 15;
    const hasUpperCase = /[A-Z]/.test(password);
    const hasLowerCase = /[a-z]/.test(password);
    const hasNumbers = /\d/.test(password);
    const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password);

    const errors = [];
    
    if (password.length < minLength || password.length > maxLength) {
      errors.push('Password must be between 6 and 15 characters');
    }
    if (!hasUpperCase) {
      errors.push('Must contain at least one uppercase letter');
    }
    if (!hasLowerCase) {
      errors.push('Must contain at least one lowercase letter');
    }
    if (!hasNumbers) {
      errors.push('Must contain at least one number');
    }
    if (!hasSpecialChar) {
      errors.push('Must contain at least one special character');
    }

    return errors;
  };

  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');

    // Validate passwords match
    if (registerData.password !== registerData.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    // Validate password strength
    const passwordErrors = validatePassword(registerData.password);
    if (passwordErrors.length > 0) {
      setError(passwordErrors.join('\n'));
      return;
    }

    setLoading(true);

    try {
      const response = await fetch('/api/public/v1/user', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: registerData.username,
          email: registerData.email,
          password: registerData.password,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Registration failed');
      }

      // Switch to login tab after successful registration
      setTab(0);
      setLoginData({
        ...loginData,
        email: registerData.email // Pre-fill email for convenience
      });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="sm" sx={{ mt: 8 }}>
      <Paper elevation={3} sx={{ p: 4 }}>
        <Tabs
          value={tab}
          onChange={(e, newValue) => setTab(newValue)}
          centered
          sx={{ mb: 3 }}
        >
          <Tab label="Login" />
          <Tab label="Register" />
        </Tabs>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {tab === 0 ? (
          // Login Form
          <Box component="form" onSubmit={handleLogin}>
            <TextField
              fullWidth
              label="Email"
              type="email"
              margin="normal"
              required
              value={loginData.email}
              onChange={(e) => setLoginData({ ...loginData, email: e.target.value })}
            />
            <TextField
              fullWidth
              label="Password"
              type="password"
              margin="normal"
              required
              value={loginData.password}
              onChange={(e) => setLoginData({ ...loginData, password: e.target.value })}
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3 }}
              disabled={loading}
            >
              {loading ? 'Logging in...' : 'Login'}
            </Button>
          </Box>
        ) : (
          // Register Form
          <Box component="form" onSubmit={handleRegister}>
            <TextField
              fullWidth
              label="Username"
              margin="normal"
              required
              value={registerData.username}
              onChange={(e) => setRegisterData({ ...registerData, username: e.target.value })}
            />
            <TextField
              fullWidth
              label="Email"
              type="email"
              margin="normal"
              required
              value={registerData.email}
              onChange={(e) => setRegisterData({ ...registerData, email: e.target.value })}
            />
            <TextField
              fullWidth
              label="Password"
              type="password"
              margin="normal"
              required
              value={registerData.password}
              onChange={(e) => setRegisterData({ ...registerData, password: e.target.value })}
              helperText="Password must be 6-15 characters with uppercase, lowercase, number, and special character"
            />
            <TextField
              fullWidth
              label="Confirm Password"
              type="password"
              margin="normal"
              required
              value={registerData.confirmPassword}
              onChange={(e) => setRegisterData({ ...registerData, confirmPassword: e.target.value })}
              error={registerData.password !== registerData.confirmPassword}
              helperText={
                registerData.password !== registerData.confirmPassword
                  ? 'Passwords do not match'
                  : ''
              }
            />
            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3 }}
              disabled={loading}
            >
              {loading ? 'Registering...' : 'Register'}
            </Button>
          </Box>
        )}
      </Paper>
    </Container>
  );
};

export default Login; 