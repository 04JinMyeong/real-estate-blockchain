// src/components/SignupSelect.js
import React from 'react';
import { useNavigate } from 'react-router-dom';

const SignupSelect = () => {
  const navigate = useNavigate();

  return (
    <div className="auth-container">
      <h2>회원 유형 선택</h2>
      <button onClick={() => navigate('/signup/user')}>🙋‍♂️ 일반 사용자로 가입</button>
      <button onClick={() => navigate('/signup/agent')}>🏢 부동산 업자로 가입</button>
    </div>
  );
};

export default SignupSelect;