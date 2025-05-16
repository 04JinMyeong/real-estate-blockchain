import React, { useState } from 'react';
import PropertyForm from './PropertyForm';
import PropertyList from './PropertyList';
import EditProfileModal from './EditProfileModal';
import './MyPage.css';

const MyPage = ({ user }) => {
  const [showEditModal, setShowEditModal] = useState(false);
  const [currentUser, setCurrentUser] = useState(user);

  const handleRegister = () => {
    console.log('📦 매물 등록 후 처리 실행됨');
  };

  const openEditModal = () => setShowEditModal(true);
  const closeEditModal = () => setShowEditModal(false);

  const handleSaveProfile = (updatedUser) => {
    setCurrentUser(updatedUser);
    setShowEditModal(false);
    // 여기에 서버에 프로필 수정 API 호출 코드 추가 가능
  };

  return (
    <div className="mypage-container">
      <div className="profile-section">
        <img src="/profile-icon.svg" alt="프로필" className="profile-img" />
        <div>
          <h2>마이페이지</h2>
          <p><strong>이메일:</strong> {currentUser?.email}</p>
          <button onClick={openEditModal}>프로필 수정</button>
        </div>
      </div>

      <div className="mypage-card">
        <h2>🏠 매물 등록</h2>
        <PropertyForm user={currentUser} onRegister={handleRegister} />
      </div>

      <div className="mypage-card" style={{ marginTop: '2rem' }}>
        <h2>📋 내 매물 목록</h2>
        <PropertyList user={currentUser} />
      </div>

      {showEditModal && (
        <EditProfileModal
          user={currentUser}
          onClose={closeEditModal}
          onSave={handleSaveProfile}
        />
      )}
    </div>
  );
};

export default MyPage;
