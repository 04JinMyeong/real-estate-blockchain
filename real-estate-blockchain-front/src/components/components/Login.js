import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './Auth.css';

const AUTH_API = 'https://1af7-165-229-229-137.ngrok-free.app'; // ✨ 서버 주소 (로컬 또는 실제 서버)

const Login = ({ onLogin }) => {
  const [form, setForm] = useState({ id: '', password: '' });
  const navigate = useNavigate();

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();

    try {
      const res = await axios.post(`${AUTH_API}/login`, {
        id: form.id,
        password: form.password
      }, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      console.log('✅ 로그인 응답:', res.data);
      const token = res.data.token;
      localStorage.setItem('token', token);
      alert('✅ 로그인 성공!');
      onLogin?.(form.id);
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
          name="id"
          type="text"
          value={form.id}
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
