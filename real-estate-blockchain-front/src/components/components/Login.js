import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './Auth.css';

const AUTH_API = 'http://165.229.125.72:8080'; // ✨ 인증 서버 주소

const Login = ({ onLogin }) => {
  const [form, setForm] = useState({ email: '', password: '' });
  const navigate = useNavigate();

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();

    // ✅ 테스트 계정 우선 처리 (아이디 비교 방식으로)
    if (form.email === 'test@test.com' && form.password === '1234') {
      alert('✅ 테스트 계정 로그인 성공!');
      localStorage.setItem('token', 'TEST_TOKEN');
      onLogin?.('test@test.com');
      navigate('/mypage');
      return;
    }

    try {
      const res = await axios.post(`${AUTH_API}/api/auth/login`, {
        email: form.email,
        password: form.password
      });
      console.log('✅ 로그인 응답:', res.data);
      const token = res.data.token;
      localStorage.setItem('token', token);
      alert('✅ 로그인 성공!');
      onLogin?.(res.data.email);
      navigate('/mypage');
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
          type="text"
          value={form.email}
          onChange={handleChange}
          placeholder="아이디"
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
