// src/components/Login.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';  // ← named import로 변경
import axios from 'axios';
import './Auth.css';

const AUTH_API = 'https://1af7-165-229-229-137.ngrok-free.app';

const Login = ({ onLogin }) => {
  const [form, setForm] = useState({ id: '', password: '' });
  const navigate = useNavigate();

  const handleChange = e => {
    setForm(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async e => {
    e.preventDefault();
    try {
      // 1) 로그인 요청
      const res = await axios.post(`${AUTH_API}/auth/login`, {
        email: form.id,
        password: form.password
      });

      // 2) 토큰 저장
      const token = res.data.token;
      localStorage.setItem('token', token);

      // 3) 디코딩해서 email, role 추출
      const { email, role } = jwtDecode(token);
      localStorage.setItem('role', role);

      // 4) 상위 컴포넌트에 로그인 정보 전달
      onLogin({ email, role });

      // 5) role에 따라 페이지 이동
      if (role === 'agent') {
        navigate('/mypage');
      } else {
        navigate('/user/mypage');
      }
    } catch (err) {
      console.error('로그인 오류:', err);
      alert('로그인에 실패했습니다. 다시 시도해 주세요.');
    }
  };

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          name="id"
          type="email"
          placeholder="이메일"
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
};

export default Login;
