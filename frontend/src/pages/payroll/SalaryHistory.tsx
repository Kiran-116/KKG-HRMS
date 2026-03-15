import React, { useEffect, useState } from 'react';
import { salaryService, Salary } from '../../services/salaryService';

const SalaryHistoryPage: React.FC = () => {
  const [salaries, setSalaries] = useState<Salary[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  useEffect(() => {
    loadSalaries();
  }, [page]);

  const loadSalaries = async () => {
    try {
      const data = await salaryService.getMySalary(page, 10);
      setSalaries(data.salaries);
    } catch (error) {
      console.error('Failed to load salaries:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-12">Loading salary history...</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">My Salary History</h1>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Month/Year</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Base Salary</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Bonus</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Deductions</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Net Salary</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {salaries?.map((salary) => (
              <tr key={salary.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {salary.month}/{salary.year}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${salary.base_salary.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${salary.bonus.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  ${salary.deductions.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900">
                  ${salary.net_salary.toLocaleString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {salaries?.length === 0 && (
        <div className="text-center py-12 text-gray-500">No salary records found</div>
      )}
    </div>
  );
};

export default SalaryHistoryPage;
