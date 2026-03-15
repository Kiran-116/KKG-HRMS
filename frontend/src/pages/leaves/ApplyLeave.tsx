import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { leaveService, ApplyLeaveRequest } from '../../services/leaveService';

const ApplyLeavePage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<ApplyLeaveRequest>({
    start_date: '',
    end_date: '',
    reason: '',
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await leaveService.apply(formData);
      alert('Leave application submitted successfully!');
      navigate('/leaves/me');
    } catch (error: any) {
      alert(error.response?.data?.message || 'Failed to apply for leave');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">Apply for Leave</h1>

      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6 space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700">Start Date *</label>
          <input
            type="date"
            name="start_date"
            required
            value={formData.start_date}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">End Date *</label>
          <input
            type="date"
            name="end_date"
            required
            value={formData.end_date}
            onChange={handleChange}
            min={formData.start_date}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Reason *</label>
          <textarea
            name="reason"
            required
            rows={4}
            value={formData.reason}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
            placeholder="Please provide a reason for your leave..."
          />
        </div>

        <div className="flex space-x-4">
          <button
            type="submit"
            disabled={loading}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:opacity-50"
          >
            {loading ? 'Submitting...' : 'Submit Application'}
          </button>
          <button
            type="button"
            onClick={() => navigate('/leaves/me')}
            className="px-4 py-2 border rounded-md"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default ApplyLeavePage;
