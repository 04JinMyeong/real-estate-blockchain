// src/components/Header.js
import React from 'react';
import './Header.css';

const Header = ({ onNavigate, user, onLogout }) => {
  return (
    <header className="navbar">
      <div className="navbar-logo" onClick={() => onNavigate('home')}>
        <img src="/truehome-logo.png" alt="TrueHome Logo" className="logo-img" />
      </div>
      <nav className="navbar-menu">
        <ul className="navbar-list">
          <li onClick={() => onNavigate('home')}>홈</li>
          <li onClick={() => onNavigate('map')}>매물</li>
          {user ? (
          <>
         <li onClick={() => onNavigate('mypage')} style={{ fontWeight: 'bold' }}>
          {user.email}님
         </li>
         <li onClick={onLogout}>로그아웃</li>
         </>
         ) : (
         <>
         <li onClick={() => onNavigate('login')}>로그인</li>
         <li onClick={() => onNavigate('signup')}>회원가입</li>
         </>
         )}
        </ul>
      </nav>
    </header>
  );
};

export default Header;
