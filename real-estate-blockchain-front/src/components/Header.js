// src/components/Header.js
import React from 'react';
import { useNavigate } from 'react-router-dom';
import './Header.css';

const Header = ({ user, onLogout }) => {
  const navigate = useNavigate();
  console.log('[Header.js] user prop:', user); // user prop 내용 확인

  // 표시할 사용자 이름 결정 (username 우선, 없으면 email, 그것도 없으면 id)
  // 로그인 시 Login.js에서 username에 사용자 ID (예: "TRUE")를 넣어주고 있음
  const displayName = user ? (user.username || user.email || user.id || '사용자') : '';

  return (
    <header className="navbar">
      <div className="navbar-logo" onClick={() => navigate('/')}>
        <img src="/truehome-logo.png" alt="TrueHome Logo" className="logo-img" />
      </div>
      <nav className="navbar-menu">
        <ul className="navbar-list">
          <li onClick={() => navigate('/')}>홈</li>
          <li onClick={() => navigate('/map')}>매물</li>
          {user && typeof displayName === 'string' ? ( // ◀◀◀ user 객체 존재 및 displayName이 문자열인지 확인
            <>
              <li onClick={() => navigate(user.role === 'agent' ? '/agent/mypage' : '/user/mypage')} style={{ fontWeight: 'bold' }}>
                {displayName}님
              </li>
              <li onClick={onLogout}>로그아웃</li>
            </>
          ) : (
            <>
              <li onClick={() => navigate('/login')}>로그인</li>
              <li onClick={() => navigate('/signup')}>회원가입</li>
            </>
          )}
        </ul>
      </nav>
    </header>
  );
};

export default Header;