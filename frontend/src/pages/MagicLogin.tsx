import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { authService } from '../services/authService';

const MagicLogin: React.FC = () => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const handleMagicLogin = async () => {
      if (!token) {
        setError('Invalid magic link. Please contact your administrator.');
        setLoading(false);
        return;
      }

      try {
        const response = await authService.magicLogin(token);
        
        // Store tokens
        localStorage.setItem('access_token', response.access_token);
        localStorage.setItem('refresh_token', response.refresh_token);
        localStorage.setItem('user', JSON.stringify(response.user));
        
        // Check if password change is required
        if (response.must_change_password) {
          navigate('/set-password');
        } else {
          navigate('/dashboard');
        }
      } catch (err: any) {
        setError(err.response?.data?.message || 'Invalid or expired magic link. Please contact your administrator.');
        setLoading(false);
      }
    };

    handleMagicLogin();
  }, [token, navigate]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Signing you in...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Magic Link Error
            </h2>
          </div>
          <div className="rounded-md bg-red-50 p-4">
            <div className="text-sm text-red-800">{error}</div>
          </div>
          <div className="text-center">
            <button
              onClick={() => navigate('/login')}
              className="text-indigo-600 hover:text-indigo-500 font-medium"
            >
              Go to Login
            </button>
          </div>
        </div>
      </div>
    );
  }

  return null;
};

export default MagicLogin;
