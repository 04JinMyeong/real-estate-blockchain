// src/components/Signup.js
import React, { useState } from 'react';
import axios from 'axios';
import './Auth.css';

const Signup = ({ onSignup }) => {
  const [form, setForm] = useState({
    email: '',
    password: '',
    confirmPassword: ''
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
      const res = await axios.post('http://localhost:3001/api/auth/signup', {
        email: form.email,
        password: form.password
      });
      alert('✅ 회원가입 성공: ' + res.data.message);
      onSignup?.();
    } catch (err) {
      alert('❌ 회원가입 실패: ' + (err.response?.data?.message || err.message));
      console.error(err);
    }
  };

  return (
    <div className="auth-container">
      <h2>회원가입</h2>
      <form onSubmit={handleSubmit}>
        <input type="email" name="email" value={form.email} onChange={handleChange} placeholder="이메일" required />
        <input type="password" name="password" value={form.password} onChange={handleChange} placeholder="비밀번호" required />
        <input type="password" name="confirmPassword" value={form.confirmPassword} onChange={handleChange} placeholder="비밀번호 확인" required />
        <button type="submit">회원가입</button>
      </form>
    </div>
  );
};

export default Signup;
