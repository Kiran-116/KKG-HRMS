import React, { useState, useEffect } from 'react';
import { attendanceService, Attendance } from '../../services/attendanceService';
import { formatTimeIST, formatDateWithWeekdayIST, getCurrentTimeIST } from '../../utils/timeUtils';

const AttendanceCheckInPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [todayAttendance, setTodayAttendance] = useState<Attendance | null>(null);
  const [checkingStatus, setCheckingStatus] = useState(false);
  const [currentTime, setCurrentTime] = useState<string>(getCurrentTimeIST());

  useEffect(() => {
    checkTodayStatus();
    
    // Update clock every second
    const interval = setInterval(() => {
      setCurrentTime(getCurrentTimeIST());
    }, 1000);
    
    return () => clearInterval(interval);
  }, []);

  const checkTodayStatus = async () => {
    try {
      setCheckingStatus(true);
      const today = new Date().toISOString().split('T')[0];
      // Get user's attendance records and find today's record
      const data = await attendanceService.getMyAttendance(1, 50);
      if (data.attendances && data.attendances.length > 0) {
        const todayRecord = data.attendances.find(
          (att) => att.date.split('T')[0] === today
        );
        if (todayRecord) {
          setTodayAttendance(todayRecord);
        } else {
          setTodayAttendance(null);
        }
      } else {
        setTodayAttendance(null);
      }
    } catch (error) {
      // If no attendance found, that's okay - user hasn't checked in today
      setTodayAttendance(null);
    } finally {
      setCheckingStatus(false);
    }
  };

  const handleCheckIn = async () => {
    setLoading(true);
    setMessage(null);
    try {
      const attendance = await attendanceService.checkIn();
      setTodayAttendance(attendance);
      setMessage({ type: 'success', text: 'Checked in successfully!' });
      // Clear message after 3 seconds
      setTimeout(() => setMessage(null), 3000);
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
      const attendance = await attendanceService.checkOut();
      setTodayAttendance(attendance);
      setMessage({ type: 'success', text: 'Checked out successfully!' });
      // Clear message after 3 seconds
      setTimeout(() => setMessage(null), 3000);
    } catch (error: any) {
      setMessage({ type: 'error', text: error.response?.data?.message || 'Failed to check out' });
    } finally {
      setLoading(false);
    }
  };

  const hasCheckedIn = todayAttendance?.check_in != null;
  const hasCheckedOut = todayAttendance?.check_out != null;

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Attendance</h1>

      <div className="bg-white rounded-lg shadow p-6 space-y-6">
        <div className="text-center">
          <p className="text-lg text-gray-600 mb-4">
            {formatDateWithWeekdayIST(new Date().toISOString())}
          </p>
          <p className="text-2xl font-bold text-gray-900">
            {currentTime}
          </p>
        </div>

        {hasCheckedIn && !hasCheckedOut && (
          <div className="p-4 rounded-md bg-yellow-50 text-yellow-800 border border-yellow-200">
            <p className="font-medium">Already checked in for this date</p>
            {todayAttendance?.check_in && (
              <p className="text-sm mt-1">
                Check-in time: {formatTimeIST(todayAttendance.check_in)} IST
              </p>
            )}
          </div>
        )}

        {hasCheckedIn && hasCheckedOut && (
          <div className="p-4 rounded-md bg-blue-50 text-blue-800 border border-blue-200">
            <p className="font-medium">Already checked out for this date</p>
            {todayAttendance?.check_in && todayAttendance?.check_out && (
              <div className="text-sm mt-1 space-y-1">
                <p>Check-in: {formatTimeIST(todayAttendance.check_in)} IST</p>
                <p>Check-out: {formatTimeIST(todayAttendance.check_out)} IST</p>
              </div>
            )}
          </div>
        )}

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
            disabled={loading || checkingStatus || hasCheckedIn}
            className="w-full px-4 py-3 bg-green-600 text-white rounded-md hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {checkingStatus ? 'Loading...' : hasCheckedIn ? '✓ Already Checked In' : 'Check In'}
          </button>

          <button
            onClick={handleCheckOut}
            disabled={loading || checkingStatus || !hasCheckedIn || hasCheckedOut}
            className="w-full px-4 py-3 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {checkingStatus ? 'Loading...' : hasCheckedOut ? '✓ Already Checked Out' : 'Check Out'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default AttendanceCheckInPage;
