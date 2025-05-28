import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const SignupUser = () => {
  const navigate = useNavigate();
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
      // API 명세상: /register 엔드포인트, role 제외
      await axios.post('https://252f-219-251-84-31.ngrok-free.app/register', {
        username: form.email,
        password: form.password
        // role: 'user'  // 명세에 없음, 프론트에서만 관리!
      });
      alert('✅ 일반 사용자 회원가입 성공!');
      // 프론트 분기용 role 저장 (원하면)
      // localStorage.setItem('role', 'user');
      navigate('/'); // 회원가입 후 홈으로 이동
    } catch (err) {
      if (err.response?.status === 409) {
        alert('❌ 이미 등록된 계정입니다.');
      } else {
        alert('❌ 회원가입 실패: ' + (err.response?.data?.message || err.message));
      }
    }
  };

  return (
    <div className="auth-container">
      <h2>🙋‍♂️ 일반 사용자 회원가입</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="email"
          name="email"
          value={form.email}
          onChange={handleChange}
          placeholder="이메일"
          required
        />
        <input
          type="password"
          name="password"
          value={form.password}
          onChange={handleChange}
          placeholder="비밀번호"
          required
        />
        <input
          type="password"
          name="confirmPassword"
          value={form.confirmPassword}
          onChange={handleChange}
          placeholder="비밀번호 확인"
          required
        />
        <button type="submit">회원가입</button>
      </form>
    </div>
  );
};

export default SignupUser;
