// src/components/Login.js
import React, { useState } from 'react';
import axios from 'axios';
import './Auth.css';

const Login = ({ onLogin }) => {
  const [form, setForm] = useState({ email: '', password: '' });

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();
    try {
      const res = await axios.post('http://localhost:3001/api/auth/login', {
        email: form.email,
        password: form.password
      });
      console.log('✅ 로그인 응답:', res.data); // 여기서 email 있는지 확인!
      const token = res.data.token;
      localStorage.setItem('token', token); // 저장
      alert('✅ 로그인 성공!');
      onLogin?.(res.data.email);  // ✅ 이 줄이 핵심!
    } catch (err) {
      console.error('❌ 로그인 실패:', err.response || err.message);
      alert('❌ 로그인 실패: ' + (err.response?.data?.message || err.message));
    }
  };
  

  return (
    <div className="auth-container">
      <h2>로그인</h2>
      <form onSubmit={handleSubmit}>
        <input
          name="email"
          type="email"
          value={form.email}
          onChange={handleChange}
          placeholder="이메일"
          required
        />
        <input
          name="password"
          type="password"
          value={form.password}
          onChange={handleChange}
          placeholder="비밀번호"
          required
        />
        <button type="submit">로그인</button>
      </form>
    </div>
  );
};

export default Login;
