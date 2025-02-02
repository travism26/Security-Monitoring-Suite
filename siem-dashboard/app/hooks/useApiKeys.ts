import { useState, useEffect } from 'react';
import { apiService } from '../services/api';

interface ApiKey {
  id: string;
  name: string;
  key: string;
  createdAt: string;
}

export function useApiKeys() {
  const [apiKeys, setApiKeys] = useState<ApiKey[]>([]);

  useEffect(() => {
    fetchApiKeys();
  }, []);

  const fetchApiKeys = async () => {
    try {
      const keys = await apiService.getApiKeys();
      setApiKeys(keys);
    } catch (error) {
      console.error('Failed to fetch API keys:', error);
    }
  };

  const createApiKey = async (name: string) => {
    try {
      const newKey = await apiService.createApiKey(name);
      setApiKeys([...apiKeys, newKey]);
      return newKey;
    } catch (error) {
      console.error('Failed to create API key:', error);
      throw error;
    }
  };

  const deleteApiKey = async (id: string) => {
    try {
      await apiService.deleteApiKey(id);
      setApiKeys(apiKeys.filter(key => key.id !== id));
    } catch (error) {
      console.error('Failed to delete API key:', error);
      throw error;
    }
  };

  return { apiKeys, createApiKey, deleteApiKey };
}

