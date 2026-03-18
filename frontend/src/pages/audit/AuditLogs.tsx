import React, { useEffect, useState } from 'react';
import { auditService, AuditLog } from '../../services/auditService';
import { formatTimeIST, formatDateIST } from '../../utils/timeUtils';

const AuditLogsPage: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [actionFilter, setActionFilter] = useState('');

  useEffect(() => {
    loadLogs();
  }, [page, actionFilter]);

  const loadLogs = async () => {
    try {
      const data = await auditService.getAll(page, 10, undefined, actionFilter || undefined);
      setLogs(data.audit_logs);
    } catch (error) {
      console.error('Failed to load audit logs:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-12">Loading audit logs...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Audit Logs</h1>
        <select
          value={actionFilter}
          onChange={(e) => setActionFilter(e.target.value)}
          className="px-4 py-2 border rounded-md"
        >
          <option value="">All Actions</option>
          <option value="CREATE_EMPLOYEE">Create Employee</option>
          <option value="UPDATE_EMPLOYEE">Update Employee</option>
          <option value="CREATE_SALARY">Create Salary</option>
          <option value="UPDATE_SALARY">Update Salary</option>
          <option value="APPLY_LEAVE">Apply Leave</option>
          <option value="APPROVE_LEAVE">Approve Leave</option>
          <option value="REJECT_LEAVE">Reject Leave</option>
          <option value="CHECKIN">Check In</option>
          <option value="CHECKOUT">Check Out</option>
          <option value="LOGIN">Login</option>
          <option value="REGISTER">Register</option>
          <option value="UPLOAD_DOCUMENT">Upload Document</option>
          <option value="DELETE_DOCUMENT">Delete Document</option>
        </select>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Timestamp</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">User Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Entity Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Changes</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">IP Address</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {logs?.map((log) => (
              <tr key={log.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <div>{formatDateIST(log.created_at)}</div>
                  <div className="text-xs text-gray-500">{formatTimeIST(log.created_at)} IST</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {log.user_name || log.user_id || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs rounded-full ${
                    log.action.startsWith('CREATE_') || log.action === 'APPLY_LEAVE' || log.action === 'REGISTER' || log.action === 'UPLOAD_DOCUMENT' ? 'bg-green-100 text-green-800' :
                    log.action.startsWith('UPDATE_') || log.action === 'APPROVE_LEAVE' || log.action === 'LOGIN' ? 'bg-blue-100 text-blue-800' :
                    log.action.startsWith('DELETE_') || log.action === 'REJECT_LEAVE' ? 'bg-red-100 text-red-800' :
                    log.action === 'CHECKIN' || log.action === 'CHECKOUT' ? 'bg-purple-100 text-purple-800' :
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {log.action.replace(/_/g, ' ')}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {log.entity_type}
                </td>
                <td className="px-6 py-4 text-sm text-gray-700 max-w-md">
                  <div className="truncate" title={log.description}>
                    {log.description || '-'}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {log.ip_address || '-'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {logs?.length === 0 && !loading && (
        <div className="text-center py-12 text-gray-500">
          <p className="text-lg font-medium mb-2">No audit logs found</p>
          <p className="text-sm">Audit logs will appear here when actions are performed in the system.</p>
        </div>
      )}

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
          disabled={logs?.length < 10}
          className="px-4 py-2 border rounded-md disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default AuditLogsPage;
