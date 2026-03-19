import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface AttendanceChartProps {
  data: Array<{
    date: string;
    present: number;
    absent: number;
  }>;
}

const AttendanceChart: React.FC<AttendanceChartProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="flex items-center justify-center h-64 text-gray-500">
        No attendance data available
      </div>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart
        data={data}
        margin={{ top: 10, right: 20, bottom: 10, left: 0 }}
        barCategoryGap={12}
        barGap={6}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis
          dataKey="date"
          tickFormatter={(value) =>
            new Date(value).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
          }
        />
        <YAxis allowDecimals={false} />
        <Tooltip
          labelFormatter={(value) => new Date(value as string).toLocaleDateString()}
          formatter={(value: number, name: string) => [
            value,
            name === 'present' ? 'Present' : 'Absent',
          ]}
        />
        <Legend />
        {/* Grouped bars: Present and Absent side-by-side per day */}
        <Bar dataKey="present" fill="#10b981" name="Present" maxBarSize={24} />
        <Bar dataKey="absent" fill="#ef4444" name="Absent" maxBarSize={24} />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default AttendanceChart;
