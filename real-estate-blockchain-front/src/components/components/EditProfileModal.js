import React, { useState } from 'react';
import './EditProfileModal.css';

const EditProfileModal = ({ user, onClose, onSave }) => {
  const [email, setEmail] = useState(user?.email || '');
  const [name, setName] = useState(user?.name || '');

  const handleSubmit = e => {
    e.preventDefault();
    onSave({ email, name });
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={e => e.stopPropagation()}>
        <h2>프로필 수정</h2>
        <form onSubmit={handleSubmit}>
          <label>
            이메일:
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
            />
          </label>
          <label>
            이름:
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              required
            />
          </label>
          <div className="modal-buttons">
            <button type="submit">저장</button>
            <button type="button" onClick={onClose}>취소</button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditProfileModal;
