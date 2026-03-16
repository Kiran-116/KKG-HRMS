import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface PayrollChartProps {
  data: Array<{
    month: number;
    year: number;
    total: number;
  }>;
}

const PayrollChart: React.FC<PayrollChartProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        No payroll data available
      </div>
    );
  }

  const chartData = data.map(item => ({
    label: `${item.month}/${item.year}`,
    total: item.total,
  })).reverse(); // Reverse to show oldest to newest

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="label" />
        <YAxis />
        <Tooltip 
          formatter={(value: number) => `$${value.toLocaleString()}`}
        />
        <Legend />
        <Bar dataKey="total" fill="#6366f1" name="Total Payroll" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default PayrollChart;
