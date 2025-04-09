import React from 'react';
import PropertyForm from '../PropertyForm';

const MyPage = ({ user }) => {
  return (
    <div style={{ padding: '2rem' }}>
      <h2>📄 마이페이지</h2>
      <p><strong>이메일:</strong> {user?.email}</p>

      <hr />

      <PropertyForm user={user} />
    </div>
  );
};

export default MyPage; // ✅ 이거 없으면 App에서 에러남!
