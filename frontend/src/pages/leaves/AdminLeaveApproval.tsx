import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { leaveService, Leave } from '../../services/leaveService';
import ConfirmModal from '../../components/ConfirmModal';

const AdminLeaveApprovalPage: React.FC = () => {
  const [leaves, setLeaves] = useState<Leave[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    type: 'approve' | 'reject';
    leaveId: string | null;
  }>({
    isOpen: false,
    type: 'approve',
    leaveId: null,
  });

  useEffect(() => {
    loadLeaves();
  }, [page, statusFilter]);

  const loadLeaves = async () => {
    try {
      const data = await leaveService.getAll(page, 10, statusFilter || undefined);
      setLeaves(data.leaves);
    } catch (error) {
      console.error('Failed to load leaves:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleApproveClick = (id: string) => {
    setConfirmModal({
      isOpen: true,
      type: 'approve',
      leaveId: id,
    });
  };

  const handleRejectClick = (id: string) => {
    setConfirmModal({
      isOpen: true,
      type: 'reject',
      leaveId: id,
    });
  };

  const handleConfirm = async () => {
    if (!confirmModal.leaveId) return;

    try {
      if (confirmModal.type === 'approve') {
        await leaveService.approve(confirmModal.leaveId);
        toast.success('Leave request approved successfully!');
      } else {
        await leaveService.reject(confirmModal.leaveId);
        toast.success('Leave request rejected');
      }
      loadLeaves();
    } catch (error: any) {
      toast.error(
        error.response?.data?.message ||
        `Failed to ${confirmModal.type === 'approve' ? 'approve' : 'reject'} leave`
      );
    } finally {
      setConfirmModal({ isOpen: false, type: 'approve', leaveId: null });
    }
  };

  const handleCancel = () => {
    setConfirmModal({ isOpen: false, type: 'approve', leaveId: null });
  };

  if (loading) {
    return <div className="text-center py-12">Loading leaves...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Leave Requests</h1>
        <select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
          className="px-4 py-2 border rounded-md"
        >
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Employee Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Start Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">End Date</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Reason</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {leaves?.map((leave) => (
              <tr key={leave.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {leave.user_name || leave.user_id}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {new Date(leave.start_date).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {new Date(leave.end_date).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{leave.reason}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs rounded-full ${
                    leave.status === 'approved' ? 'bg-green-100 text-green-800' :
                    leave.status === 'rejected' ? 'bg-red-100 text-red-800' :
                    'bg-yellow-100 text-yellow-800'
                  }`}>
                    {leave.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                  {leave.status === 'pending' && (
                    <>
                      <button
                        onClick={() => handleApproveClick(leave.id)}
                        className="text-green-600 hover:text-green-900"
                      >
                        Approve
                      </button>
                      <button
                        onClick={() => handleRejectClick(leave.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        Reject
                      </button>
                    </>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <ConfirmModal
        isOpen={confirmModal.isOpen}
        title={confirmModal.type === 'approve' ? 'Approve Leave Request' : 'Reject Leave Request'}
        message={
          confirmModal.type === 'approve'
            ? 'Are you sure you want to approve this leave request?'
            : 'Are you sure you want to reject this leave request?'
        }
        confirmText={confirmModal.type === 'approve' ? 'Approve' : 'Reject'}
        cancelText="Cancel"
        type={confirmModal.type === 'reject' ? 'danger' : 'info'}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
      />
    </div>
  );
};

export default AdminLeaveApprovalPage;
