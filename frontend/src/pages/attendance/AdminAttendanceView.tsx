import React, { useEffect, useState } from 'react';
import { attendanceService, Attendance } from '../../services/attendanceService';

const AdminAttendanceViewPage: React.FC = () => {
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [userIdFilter, setUserIdFilter] = useState('');
  const [dateFilter, setDateFilter] = useState('');

  useEffect(() => {
    loadAttendance();
  }, [page, userIdFilter, dateFilter]);

  const loadAttendance = async () => {
    try {
      const data = await attendanceService.getAll(
        page,
        10,
        userIdFilter || undefined,
        dateFilter || undefined
      );
      setAttendances(data.attendances);
    } catch (error) {
      console.error('Failed to load attendance:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'present':
        return 'bg-green-100 text-green-800';
      case 'late':
        return 'bg-yellow-100 text-yellow-800';
      case 'half_day':
        return 'bg-orange-100 text-orange-800';
      case 'absent':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return <div className="text-center py-12">Loading attendance...</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">All Attendance</h1>

      <div className="bg-white rounded-lg shadow p-4 flex space-x-4">
        <input
          type="text"
          placeholder="Filter by User ID"
          value={userIdFilter}
          onChange={(e) => setUserIdFilter(e.target.value)}
          className="px-4 py-2 border rounded-md"
        />
        <input
          type="date"
          value={dateFilter}
          onChange={(e) => setDateFilter(e.target.value)}
          className="px-4 py-2 border rounded-md"
        />
        <button
          onClick={() => {
            setUserIdFilter('');
            setDateFilter('');
          }}
          className="px-4 py-2 border rounded-md"
        >
          Clear Filters
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">User ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Check In</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Check Out</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {attendances?.map((attendance) => (
              <tr key={attendance.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{attendance.user_id}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {new Date(attendance.date).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {attendance.check_in ? new Date(attendance.check_in).toLocaleTimeString() : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {attendance.check_out ? new Date(attendance.check_out).toLocaleTimeString() : '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs rounded-full ${getStatusColor(attendance.status)}`}>
                    {attendance.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="flex justify-between">
        <button
          onClick={() => setPage(p => Math.max(1, p - 1))}
          disabled={page === 1}
          className="px-4 py-2 border rounded-md disabled:opacity-50"
        >
          Previous
        </button>
        <button
          onClick={() => setPage(p => p + 1)}
          disabled={attendances?.length < 10}
          className="px-4 py-2 border rounded-md disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default AdminAttendanceViewPage;
