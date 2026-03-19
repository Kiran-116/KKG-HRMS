import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { toast } from 'react-toastify';
import { employeeService, CreateEmployeeRequest, UpdateEmployeeRequest } from '../../services/employeeService';

const EmployeeFormPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreateEmployeeRequest>({
    name: '',
    email: '',
    role: 'employee',
    department: '',
    designation: '',
    joining_date: '',
    salary: 0,
  });

  useEffect(() => {
    const loadEmployee = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const emp = await employeeService.getById(id);
        // Convert possible ISO date to yyyy-MM-dd for input[type=date]
        const joiningDate =
          emp.joining_date ? new Date(emp.joining_date).toISOString().slice(0, 10) : '';

        setFormData({
          name: emp.name || '',
          email: emp.email || '',
          role: (emp.role as 'admin' | 'employee') || 'employee',
          department: emp.department || '',
          designation: emp.designation || '',
          joining_date: joiningDate,
          salary: emp.salary ?? 0,
        });
      } catch (error) {
        // silent; page still shows empty form if fetch fails
      } finally {
        setLoading(false);
      }
    };

    loadEmployee();
  }, [id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'salary' ? parseFloat(value) || 0 : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    setLoading(true);
    try {
      if (id) {
        const updateData: UpdateEmployeeRequest = {
          name: formData.name,
          email: formData.email,
          role: formData.role as 'admin' | 'employee',
          department: formData.department,
          designation: formData.designation,
          joining_date: formData.joining_date,
          salary: formData.salary,
        };
        await employeeService.update(id, updateData);
        toast.success('Employee updated successfully!');
      } else {
        const createData: CreateEmployeeRequest = {
          name: formData.name,
          email: formData.email,
          role: formData.role as 'admin' | 'employee',
          department: formData.department,
          designation: formData.designation,
          joining_date: formData.joining_date,
          salary: formData.salary,
        };
        await employeeService.create(createData);
        toast.success('Employee created successfully! A magic link has been sent to their email to set their password.');
      }
      navigate('/employees');
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Failed to save employee');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">
        {id ? 'Edit Employee' : 'Add New Employee'}
      </h1>

      <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-6 space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700">Name *</label>
          <input
            type="text"
            name="name"
            required
            value={formData.name}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Email *</label>
          <input
            type="email"
            name="email"
            required
            value={formData.email}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>


        <div>
          <label className="block text-sm font-medium text-gray-700">Role</label>
          <select
            name="role"
            value={formData.role}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          >
            <option value="employee">Employee</option>
            <option value="admin">Admin</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Department</label>
          <input
            type="text"
            name="department"
            value={formData.department}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Designation</label>
          <input
            type="text"
            name="designation"
            value={formData.designation}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Joining Date</label>
          <input
            type="date"
            name="joining_date"
            value={formData.joining_date}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">Salary</label>
          <input
            type="number"
            name="salary"
            min="0"
            value={formData.salary}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
          />
        </div>

        <div className="flex space-x-4">
          <button
            type="submit"
            disabled={loading}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:opacity-50"
          >
            {loading ? 'Saving...' : 'Save'}
          </button>
          <button
            type="button"
            onClick={() => navigate('/employees')}
            className="px-4 py-2 border rounded-md"
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default EmployeeFormPage;
