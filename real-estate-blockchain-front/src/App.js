import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';
import './App.css';

import Header from './components/Header';
import HomeLanding from './components/HomeLanding';
import MapView from './components/MapView';
import Login from './components/Login';
import SignupSelect from './components/SignupSelect';
import SignupUser from './components/SignupUser';
import SignupAgent from './components/SignupAgent';
import MyPage from './components/MyPage';

function App() {
  const [user, setUser] = useState(null);

  const handleLogin = (email) => {
    setUser({ email });
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
        setUser({ email: decoded.email });
      } catch (err) {
        console.error('토큰 디코딩 실패:', err);
        localStorage.removeItem('token');
      }
    }
  }, []);

  return (
    <Router>
      <div className="app">
        <Header user={user} onLogout={handleLogout} />

        <Routes>
          <Route path="/" element={<HomeLanding />} />
          <Route path="/map" element={<MapView user={user} />} />
          <Route path="/login" element={<Login onLogin={handleLogin} />} />
          <Route path="/signup" element={<SignupSelect />} />
          <Route path="/signup/user" element={<SignupUser />} />
          <Route path="/signup/agent" element={<SignupAgent />} />
          <Route path="/mypage" element={<MyPage user={user} />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
