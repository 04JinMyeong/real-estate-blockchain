// src/App.js
import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';
import './App.css';
import backgroundImage from './background.jpg';

import Header from './components/Header';
// import HomeLanding from './components/HomeLanding';
import MapView from './components/MapView';
import Login from './components/Login';
import SignupSelect from './components/SignupSelect';
import SignupUser from './components/SignupUser';
import SignupAgent from './components/SignupAgent';
import AgentMypage from './components/AgentMypage';
import MainPage from './components/MainPage';
import UserMypage from './components/UserMypage';
import AOS from 'aos';
import 'aos/dist/aos.css';

function App() {

  const [darkMode, setDarkMode]=useState(false);

  useEffect(() => {
    if (darkMode) {
      document.body.classList.add('dark-mode');
    } else {
      document.body.classList.remove('dark-mode');
    }
  }, [darkMode]);

  useEffect(() => {
    AOS.init({
      duration: 1000,
      once: false,
      mirror: true,
      offset: 60,
      easing: "ease-in-out"
    });
  }, []);

  const [user, setUser] = useState(null);

  const handleLogin = (userInfo) => {
    setUser(userInfo); // email, role 모두 user에 저장
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decoded = jwtDecode(token);
        setUser({ email: decoded.email, role: decoded.role }); // role까지
      } catch (err) {
        console.error('토큰 디코딩 실패:', err);
        localStorage.removeItem('token');
      }
    }
  }, []);

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
        {/* 다크모드 토글 버튼 (상단에 고정) */}
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
          {darkMode ? '☀️ ' : '🌙 '}
        </button>

        <Header user={user} onLogout={handleLogout} />

        <main style={{ flex: 1 }}>
          <Routes>
            <Route path="/" element={<MainPage />} />
            <Route path="/map" element={<MapView user={user} />} />
            <Route path="/login" element={<Login onLogin={handleLogin} />} />
            <Route path="/signup" element={<SignupSelect />} />
            <Route path="/signup/user" element={<SignupUser />} />
            <Route path="/signup/agent" element={<SignupAgent />} />

            <Route path="/agent/mypage" element={<AgentMypage user={user} />} />
            <Route path="/user/mypage" element={<UserMypage user={user} />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
