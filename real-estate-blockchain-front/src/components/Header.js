// src/components/Header.js
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Header.css';

const Header = ({ user, onLogout }) => {
  const navigate = useNavigate();
  const [scrolled, setScrolled] = useState(false);

  // 스크롤에 따라 scrolled 상태 업데이트
  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50);
    };
    window.addEventListener('scroll', handleScroll);
    // clean-up
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // 로그아웃 시 토큰·역할·이메일 정리 후 부모 상태 초기화
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    localStorage.removeItem('email');
    onLogout();
    navigate('/login');
  };

  // role에 따라 마이페이지 이동
  const handleMyPage = () => {
  // user?.role이 agent면, agent 마이페이지로 이동
  if (user?.role === 'agent') {
    navigate('/agent/mypage');
  } else if (user?.role === 'user') {
    navigate('/user/mypage');
  } else {
    navigate('/mypage'); // fallback
  }
};

  return (
    <header className={`navbar ${scrolled ? 'scrolled' : ''}`}>
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
                className="user-email"
                style={{ cursor: 'pointer' }}
              >
                {user.email}님
              </li>
              <li onClick={handleLogout} style={{ cursor: 'pointer' }}>
                로그아웃
              </li>
            </>
          ) : (
            <>
              <li
                onClick={() => navigate('/login')}
                style={{ cursor: 'pointer' }}
              >
                로그인
              </li>
              <li
                onClick={() => navigate('/signup')}
                style={{ cursor: 'pointer' }}
              >
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
