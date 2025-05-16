// src/components/Header.js
import React from 'react';
import { useNavigate } from 'react-router-dom';
import './Header.css';

const Header = ({ user, onLogout }) => {
  const navigate = useNavigate();

  return (
    <header className="navbar">
      <div className="navbar-logo" onClick={() => navigate('/')}>
        <img src="/truehome-logo.png" alt="TrueHome Logo" className="logo-img" />
      </div>
      <nav className="navbar-menu">
        <ul className="navbar-list">
          <li onClick={() => navigate('/')}>홈</li>
          <li onClick={() => navigate('/map')}>매물</li>
          {user ? (
            <>
              <li onClick={() => navigate('/mypage')} style={{ fontWeight: 'bold' }}>
                {user.email}님
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
