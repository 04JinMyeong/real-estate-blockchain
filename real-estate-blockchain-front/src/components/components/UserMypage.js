import React from 'react';
import './MyPage.css';

const UserMypage = ({ user }) => {
  return (
    <div className="mypage-container">
      <h2>마이페이지</h2>
      <p><strong>아이디:</strong> {user?.id}</p>

      <div className="mypage-card">
        <h2>📌 내 예약 목록</h2>
        <p>여기에 예약한 매물 정보 또는 찜한 매물 등을 표시합니다.</p>
      </div>
    </div>
  );
};

export default UserMypage;
