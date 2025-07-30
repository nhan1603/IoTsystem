import { useNavigate } from 'react-router-dom';
import { useCallback } from 'react';


export const createAuthenticatedRequest = async (url, options = {}) => {
  const token = localStorage.getItem('token');
  if (!token) {
    throw new Error('No authentication token found');
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': token,
      'Content-Type': 'application/json',
    },
  });

  if (response.status === 401) {
    localStorage.removeItem('token');
    throw new Error('unauthorized');
  }

  const data = await response.json();
  
  if (!response) {
    throw new Error(data.error || `HTTP error! status: ${response.status}`);
  }

  return data;
};

export const useAuthenticatedFetch = () => {
  const navigate = useNavigate();

  // Memoize the fetchWithAuth function
  return useCallback(async (url, options = {}) => {
    try {
      return await createAuthenticatedRequest(url, options);
    } catch (error) {
      if (error.message === 'unauthorized' || error.message === 'No authentication token found') {
        localStorage.removeItem('token');
        localStorage.removeItem('cart');
        navigate('/login');
      }
      throw error;
    }
  }, [navigate]);
}; 