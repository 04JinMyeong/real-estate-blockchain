import React, { useRef } from 'react';
import PropertyForm from './PropertyForm';
import PropertyList from './PropertyList';
import './MyPage.css';

const MyPage = ({ user }) => {
  const propertyListRef = useRef();

  const handleRegister = () => {
    console.log('📦 등록 완료 후 리스트 새로고침!');
    propertyListRef.current?.fetchProperties();
  };

  return (
    <div className="mypage-container">
      <h2>마이페이지</h2>
      <p><strong>아이디:</strong> {user?.id}</p>

      <div className="mypage-card">
        <h2>🏠 매물 등록</h2>
        <PropertyForm user={user} onRegister={handleRegister} />
      </div>

      <div className="mypage-card" style={{ marginTop: '2rem' }}>
        <h2>📋 내 매물 목록</h2>
        <PropertyList ref={propertyListRef} user={user} mode="my" />
      </div>
    </div>
  );
};

export default MyPage;
