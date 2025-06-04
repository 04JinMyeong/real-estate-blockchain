import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Header.css';

const Header = ({ user, onLogout }) => {
  const navigate = useNavigate();
  const [scrolled, setScrolled] = useState(false);

  // 스크롤 감지
  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 50);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // 로그아웃
  const handleLogout = () => {
    localStorage.removeItem('token');
    onLogout();
    navigate('/login');
  };

  // 내 페이지 버튼 클릭
  const handleMyPage = () => {
    console.log('👤 Header user prop:', user);
    if (user?.role === 'agent') {
      console.log('🚀 Header navigation to /agent/mypage');
      navigate('/agent/mypage');
    } else if (user?.role === 'user') {
      console.log('🚀 Header navigation to /user/mypage');
      navigate('/user/mypage');
    } else {
      console.log('🚀 Header navigation fallback to /mypage');
      navigate('/mypage');
    }
  };

  return (
    <header className={`navbar ${scrolled ? 'scrolled' : ''}`}>
      <div className="navbar-logo" onClick={() => navigate('/')} style={{ cursor: 'pointer' }}>
        <img src="/truehome-logo.png" alt="TrueHome Logo" className="logo-img" />
      </div>

      <nav className="navbar-menu">
        <ul className="navbar-list">
          <li onClick={() => navigate('/')} style={{ cursor: 'pointer' }}>홈</li>
          <li onClick={() => navigate('/map')} style={{ cursor: 'pointer' }}>매물</li>

          {user ? (
            <>
              <li onClick={handleMyPage} className="user-email" style={{ cursor: 'pointer' }}>
                {user.email}님
              </li>
              <li onClick={handleLogout} style={{ cursor: 'pointer' }}>로그아웃</li>
            </>
          ) : (
            <>
              <li onClick={() => navigate('/login')} style={{ cursor: 'pointer' }}>로그인</li>
              <li onClick={() => navigate('/signup')} style={{ cursor: 'pointer' }}>회원가입</li>
            </>
          )}
        </ul>
      </nav>
    </header>
  );
};

export default Header;
