import React, { useState, useEffect } from 'react';
import './App.css';
import Header from './components/Header';
import HomeLanding from './components/HomeLanding';
import MapView from './components/MapView';
import Login from './components/Login';
import Signup from './components/Signup';
import { jwtDecode } from 'jwt-decode';                 // ✅ jwtDecode 추가
import MyPage from './components/MyPage';

function App() {
  const [view, setView] = useState('home');
  const [user, setUser] = useState(null);

  const handleLogin = (email) => {
    setUser({ email });       // ✅ 사용자 상태 저장
    setView('home');          // 홈으로 이동
  };

  const handleLogout = () => {
    localStorage.removeItem('token'); // ✅ 토큰 삭제
    setUser(null);
    setView('home');
  };

   // ✅ 로그인 유지
   useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      try {
        const decoded = jwtDecode(token);
        setUser({ email: decoded.email });
      } catch (err) {
        console.error('토큰 디코딩 실패:', err);
        localStorage.removeItem('token'); // 토큰이 이상하면 제거
      }
    }
  }, []);

  return (
    <div className="app">
      <Header onNavigate={setView} user={user} onLogout={handleLogout} />

      {view === 'home' && <HomeLanding onStart={() => setView('map')} />}
      {view === 'map' && <MapView user={user} />}
      {view === 'login' && <Login onLogin={handleLogin} />}
      {view === 'signup' && <Signup onSignup={() => setView('home')} />}
      {view === 'mypage' && <MyPage user={user} />}
    </div>
  );
}

export default App;
