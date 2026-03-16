import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface AttendanceChartProps {
  data: Array<{
    date: string;
    present: number;
    absent?: number;
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
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis 
          dataKey="date" 
          tickFormatter={(value) => new Date(value).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
        />
        <YAxis />
        <Tooltip 
          labelFormatter={(value) => new Date(value).toLocaleDateString()}
        />
        <Legend />
        <Line 
          type="monotone" 
          dataKey="present" 
          stroke="#10b981" 
          strokeWidth={2}
          name="Present"
        />
        {data[0]?.absent !== undefined && (
          <Line 
            type="monotone" 
            dataKey="absent" 
            stroke="#ef4444" 
            strokeWidth={2}
            name="Absent"
          />
        )}
      </LineChart>
    </ResponsiveContainer>
  );
};

export default AttendanceChart;
