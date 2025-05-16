import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const AUTH_API = 'https://1af7-165-229-229-137.ngrok-free.app'; // ⚡️ 서버 주소

const SignupAgent = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    id: '',
    password: '',
    confirmPassword: '',
    email: ''
  });

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();

    if (form.password !== form.confirmPassword) {
      alert('❌ 비밀번호가 일치하지 않습니다.');
      return;
    }

    try {
      const res = await axios.post(`${AUTH_API}/signup`, {
        id: form.id,
        password: form.password,
        email: form.email
      }, {
        headers: {
          'Content-Type': 'application/json'
        }
      });

      alert('✅ 회원가입 + 자동 등록 성공!');
      navigate('/');
    } catch (err) {
      console.error('❌ 에러:', err.response || err.message);
      alert('❌ 회원가입 실패: ' + (err.response?.data?.message || err.message));
    }
  };

  return (
    <div className="auth-container">
      <h2>부동산업자 회원가입</h2>
      <form onSubmit={handleSubmit}>
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
        <input
          name="confirmPassword"
          type="password"
          placeholder="비밀번호 확인"
          value={form.confirmPassword}
          onChange={handleChange}
          required
        />
        <input
          name="email"
          type="email"
          placeholder="이메일"
          value={form.email}
          onChange={handleChange}
          required
        />
        <button type="submit">회원가입</button>
      </form>
    </div>
  );
};

export default SignupAgent;
