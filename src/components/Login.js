// src/components/Login.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';  // named export
import axios from 'axios';
import './Auth.css';

const AUTH_API = 'https://2094-165-229-229-106.ngrok-free.app/login';

export default function Login({ onLogin }) {
  const [form, setForm] = useState({ id: '', password: '' });
  const navigate = useNavigate();

  const handleChange = e => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async e => {
    e.preventDefault();

    // 1️⃣ 디버그: 보내기 직전 payload
    console.log('👉 Login payload:', {
      id: form.id,
      password: form.password
    });

    try {
      const res = await axios.post(
        AUTH_API,
        {
          id:       form.id,
          password: form.password
        },
        {
          headers: {
            'Content-Type': 'application/json',
            'ngrok-skip-browser-warning': 'true'
          }
        }
      );

      // 2️⃣ 디버그: 성공 응답
      console.log('✅ login success response:', res.data);

      const { message, token, role, user: userId } = res.data;
      console.log('message:', message);
      console.log('token:', token);
      console.log('role:', role, 'user:', userId);

      // 3️⃣ 토큰 저장
      localStorage.setItem('token', token);

      // 4️⃣ JWT 디코딩 (선택)
      let email    = form.id;
      let username = userId;
      try {
        const decoded = jwtDecode(token);
        console.log('🔓 decoded JWT:', decoded);
        email    = decoded.email    || email;
        username = decoded.username || username;
      } catch (err) {
        console.warn('⚠️ JWT decode failed:', err);
      }

      // 5️⃣ 로그인 정보 전달
      onLogin({ email, role, username });

      // 6️⃣ role 분기
      if (role === 'agent') {
        navigate('/agent/mypage');
      } else {
        navigate('/user/mypage');
      }
    } catch (err) {
      // 7️⃣ 디버그: 에러 응답까지 모두 로그
      console.error('❌ login error response.data:', err.response?.data);
      console.error('❌ login error status:', err.response?.status);
      alert(err.response?.data?.error || err.message || '로그인에 실패했습니다.');
    }
  };

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          name="id"
          type="text"
          placeholder="아이디"
          value={form.id}
          onChange={handleChange}
          required
        />
        <input
          name="password"
          type="password"
          placeholder="비밀번호"
          value={form.password}
          onChange={handleChange}
          required
        />
        <button type="submit">로그인</button>
      </form>
    </div>
  );
}
