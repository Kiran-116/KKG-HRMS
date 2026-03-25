import React, { useEffect, useState } from 'react';
import { authService, User } from '../services/authService';
import { useAuth } from '../contexts/AuthContext';

const formatDate = (value?: string): string => {
  if (!value) return '-';
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return '-';
  return parsed.toLocaleDateString();
};

const Profile: React.FC = () => {
  const { user: authUser } = useAuth();
  const [profile, setProfile] = useState<User | null>(authUser);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const data = await authService.getMe();
        setProfile(data);
      } catch (error) {
        console.error('Failed to load profile:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  if (loading) {
    return (
      <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6">
        <p className="text-gray-600">Loading profile...</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6 md:p-8">
        <div className="flex items-center justify-between gap-4 mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">My Profile</h1>
            <p className="text-gray-600 mt-1">View your account details</p>
          </div>
          <div className="h-14 w-14 rounded-full bg-indigo-100 text-indigo-700 flex items-center justify-center text-xl font-semibold">
            {profile?.name?.charAt(0)?.toUpperCase() || 'U'}
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <ProfileField label="Name" value={profile?.name} />
          <ProfileField label="Email" value={profile?.email} />
          <ProfileField label="Role" value={profile?.role} />
          <ProfileField label="Department" value={profile?.department} />
          <ProfileField label="Designation" value={profile?.designation} />
          <ProfileField label="Joining Date" value={formatDate(profile?.joining_date)} />
          <ProfileField label="Status" value={profile?.is_active ? 'Active' : 'Inactive'} />
          <ProfileField label="Member Since" value={formatDate(profile?.created_at)} />
        </div>
      </div>
    </div>
  );
};

interface ProfileFieldProps {
  label: string;
  value?: string;
}

const ProfileField: React.FC<ProfileFieldProps> = ({ label, value }) => {
  return (
    <div className="rounded-lg border border-gray-100 bg-gray-50 px-4 py-3">
      <p className="text-xs font-medium uppercase tracking-wide text-gray-500">{label}</p>
      <p className="mt-1 text-sm font-semibold text-gray-900">{value || '-'}</p>
    </div>
  );
};

export default Profile;
