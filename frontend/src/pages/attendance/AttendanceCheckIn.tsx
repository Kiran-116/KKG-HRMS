import React, { useState } from 'react';
import { attendanceService } from '../../services/attendanceService';

const AttendanceCheckInPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [hasCheckedIn, setHasCheckedIn] = useState(false);
  const [hasCheckedOut, setHasCheckedOut] = useState(false);

  const handleCheckIn = async () => {
    setLoading(true);
    setMessage(null);
    try {
      await attendanceService.checkIn();
      setMessage({ type: 'success', text: 'Checked in successfully!' });
      setHasCheckedIn(true);
    } catch (error: any) {
      setMessage({ type: 'error', text: error.response?.data?.message || 'Failed to check in' });
    } finally {
      setLoading(false);
    }
  };

  const handleCheckOut = async () => {
    setLoading(true);
    setMessage(null);
    try {
      await attendanceService.checkOut();
      setMessage({ type: 'success', text: 'Checked out successfully!' });
      setHasCheckedOut(true);
    } catch (error: any) {
      setMessage({ type: 'error', text: error.response?.data?.message || 'Failed to check out' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Attendance</h1>

      <div className="bg-white rounded-lg shadow p-6 space-y-6">
        <div className="text-center">
          <p className="text-lg text-gray-600 mb-4">
            {new Date().toLocaleDateString('en-US', {
              weekday: 'long',
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}
          </p>
          <p className="text-2xl font-bold text-gray-900">
            {new Date().toLocaleTimeString()}
          </p>
        </div>

        {message && (
          <div
            className={`p-4 rounded-md ${
              message.type === 'success' ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'
            }`}
          >
            {message.text}
          </div>
        )}

        <div className="space-y-4">
          <button
            onClick={handleCheckIn}
            disabled={loading || hasCheckedIn}
            className="w-full px-4 py-3 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {hasCheckedIn ? '✓ Already Checked In' : 'Check In'}
          </button>

          <button
            onClick={handleCheckOut}
            disabled={loading || hasCheckedOut || !hasCheckedIn}
            className="w-full px-4 py-3 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {hasCheckedOut ? '✓ Already Checked Out' : 'Check Out'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default AttendanceCheckInPage;
