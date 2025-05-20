// src/components/Header.js
import React from 'react';
import { useNavigate } from 'react-router-dom';
import './Header.css';

const Header = ({ user, onLogout }) => {
  const navigate = useNavigate();

  // 로그아웃 시 토큰 제거 후 로그인 페이지로 이동
  const handleLogout = () => {
    onLogout();
    navigate('/login');
  };

  // role에 따라 에이전트 또는 일반 유저 마이페이지로 이동
  const handleMyPage = () => {
    const role = localStorage.getItem('role');
    if (role === 'agent') {
      navigate('/mypage');
    } else {
      navigate('/user/mypage');
    }
  };

  return (
    <header className="navbar">
      <div
        className="navbar-logo"
        onClick={() => navigate('/')}
        style={{ cursor: 'pointer' }}
      >
        <img
          src="/truehome-logo.png"
          alt="TrueHome Logo"
          className="logo-img"
        />
      </div>

      <nav className="navbar-menu">
        <ul className="navbar-list">
          <li onClick={() => navigate('/')} style={{ cursor: 'pointer' }}>
            홈
          </li>
          <li onClick={() => navigate('/map')} style={{ cursor: 'pointer' }}>
            매물
          </li>

          {user ? (
            <>
              <li
                onClick={handleMyPage}
                style={{ fontWeight: 'bold', cursor: 'pointer' }}
              >
                {user.email}님
              </li>
              <li onClick={handleLogout} style={{ cursor: 'pointer' }}>
                로그아웃
              </li>
            </>
          ) : (
            <>
              <li onClick={() => navigate('/login')} style={{ cursor: 'pointer' }}>
                로그인
              </li>
              <li onClick={() => navigate('/signup')} style={{ cursor: 'pointer' }}>
                회원가입
              </li>
            </>
          )}
        </ul>
      </nav>
    </header>
  );
};

export default Header;
