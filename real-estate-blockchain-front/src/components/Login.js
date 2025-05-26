import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';
import axios from 'axios';
import './Auth.css';

const AUTH_API = 'https://252f-219-251-84-31.ngrok-free.app/login';

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
      const res = await axios.post(AUTH_API, {
        username: form.id,
        password: form.password
      });

      // 2) 토큰 저장
      const token = res.data.token;
      localStorage.setItem('token', token);

      // 3) 토큰에서 정보 추출
      let email = form.id;
      let role = 'user';
      let username = form.id; // 기본값(입력값)
      try {
        const decoded = jwtDecode(token);
        email = decoded.email || form.id;
        role = decoded.role || 'user';
        username = decoded.username || decoded.email || form.id;
      } catch (e) {
        // 토큰 구조 불분명 시 fallback
      }
      localStorage.setItem('role', role);

      // 4) onLogin에 username 포함!
      onLogin({ email, role, username });

      // 5) 페이지 분기
      if (role === 'agent') {
        navigate('/agent/mypage');
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
          type="text"
          placeholder="아이디(이메일 또는 username)"
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
