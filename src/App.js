// src/App.js
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
// jwt-decode는 named export로 jwtDecode만 제공합니다.
import { jwtDecode } from 'jwt-decode';
import './App.css';
import backgroundImage from './background.jpg';
import Header from './components/Header';
import MapView from './components/MapView';
import Login from './components/Login';
import SignupSelect from './components/SignupSelect';
import SignupUser from './components/SignupUser';
import SignupAgent from './components/SignupAgent';
import AgentMypage from './components/AgentMypage';
import MainPage from './components/MainPage';
import UserMypage from './components/UserMypage';
import PropertyDetail from './components/PropertyDetail'; // 상세페이지 컴포넌트
import AOS from 'aos';
import 'aos/dist/aos.css';

function App() {
  const [darkMode, setDarkMode] = useState(false);
  const [user, setUser] = useState(null);

  // 다크 모드 토글
  useEffect(() => {
    document.body.classList.toggle('dark-mode', darkMode);
  }, [darkMode]);

  // AOS 초기화
  useEffect(() => {
    AOS.init({
      duration: 1000,
      once: false,
      mirror: true,
      offset: 60,
      easing: 'ease-in-out'
    });
  }, []);

  // 로그인된 사용자 상태 초기화
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token && (token.match(/\./g) || []).length === 2) {
      try {
        const decoded = jwtDecode(token);
        setUser({
          email: decoded.email,
          role: decoded.role,
          username: decoded.username || decoded.email
        });
      } catch (err) {
        console.error('토큰 디코딩 실패:', err);
        localStorage.removeItem('token');
      }
    } else {
      localStorage.removeItem('token');
    }
  }, []);

  const handleLogin = (userInfo) => {
    setUser(userInfo);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  return (
    <Router>
      <div
        className="App"
        style={{
          background: `url(${backgroundImage}) no-repeat center center fixed`,
          backgroundSize: 'cover',
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column'
        }}
      >
        {/* 다크 모드 토글 버튼 */}
        <button
          onClick={() => setDarkMode(dm => !dm)}
          style={{
            position: 'fixed',
            top: 300,
            left: 18,
            zIndex: 9999,
            background: darkMode ? '#191b2b' : '#efefef',
            color: darkMode ? '#fff' : '#111',
            border: '1px solid #bbb',
            padding: '8px 16px',
            borderRadius: '18px',
            fontWeight: 'bold',
            boxShadow: '0 2px 8px rgba(0,0,0,0.10)'
          }}
        >
          {darkMode ? '☀️' : '🌙'}
        </button>

        {/* 헤더 */}
        <Header user={user} onLogout={handleLogout} />

        {/* 메인 콘텐츠 */}
        <main style={{ flex: 1 }}>
          <Routes>
            <Route path="/" element={<MainPage user={user} />} />
            <Route path="/map" element={<MapView user={user} />} />
            <Route path="/login" element={<Login onLogin={handleLogin} />} />
            <Route path="/signup" element={<SignupSelect />} />
            <Route path="/signup/user" element={<SignupUser />} />
            <Route path="/signup/agent" element={<SignupAgent />} />
            <Route path="/agent/mypage" element={<AgentMypage user={user} />} />
            <Route path="/user/mypage" element={<UserMypage user={user} />} />
            {/* 매물 상세 페이지 라우트 */}
            <Route path="/properties/:id" element={<PropertyDetail user={user} />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
