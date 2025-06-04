import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const SignupUser = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    id: '',
    email: '',
    password: '',
    confirmPassword: ''
  });

  // 입력 변경 핸들러 (이메일 입력 시 id 자동 완성)
  const handleChange = e => {
    const { name, value } = e.target;
    // 이메일 입력 시 id 자동 완성
    if (name === 'email') {
      const idPart = value.split('@')[0];
      setForm(prev => ({
        ...prev,
        email: value,
        id: prev.id ? prev.id : idPart // id 칸을 이미 사용자가 직접 입력했다면 덮어쓰지 않음
      }));
    } else {
      setForm(prev => ({ ...prev, [name]: value }));
    }
  };

  // 회원가입 제출
  const handleSubmit = async e => {
    e.preventDefault();
    if (form.password !== form.confirmPassword) {
      alert('❌ 비밀번호가 일치하지 않습니다.');
      return;
    }

    try {
      await axios.post('https://2094-165-229-229-106.ngrok-free.app/signup', {
        id: form.id,
        email: form.email,
        password: form.password,
        role: 'user' 
      });
      alert('✅ 일반 사용자 회원가입 성공!');
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
          type="text"
          name="id"
          value={form.id}
          onChange={handleChange}
          placeholder="아이디 (이메일 앞부분, 직접 수정 가능)"
          required
        />
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
