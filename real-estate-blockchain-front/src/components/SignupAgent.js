import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './Auth.css';

const SignupAgent = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
    licenseNum: '',
    phone: ''
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
   // 필수 입력 체크
    if (!form.name || !form.licenseNum || !form.phone) {
      alert('❌ 이름, 자격번호, 전화번호를 모두 입력해주세요.');
      return;
     }

    try {
      await axios.post('http://localhost:3001/api/auth/signup', {
        email: form.email,
        password: form.password,
        role: 'agent'
      });

      
     // 2) Go-Backend에 공인중개사 등록 + DID/VC 발급 요청
     const vcRes = await axios.post('http://localhost:8080/api/broker/register', {
       name: form.name,
       licenseNum: form.licenseNum,
       phone: form.phone
     });
     const { did, vc } = vcRes.data;

     // 3) 받은 DID/VC를 로컬에 저장
     localStorage.setItem('brokerDid', did);
     localStorage.setItem('brokerVC', JSON.stringify(vc));

     alert(`✅ 회원가입 & VC 발급 성공!\nDID: ${did}`);
     // 4) 발급받은 VC 확인 페이지로 이동
     navigate('/my-vc');
    } catch (err) {
      alert('❌ 회원가입 또는 VC 발급 실패: ' +
        (err.response?.data?.message || err.message));
    }
  };


  return (
    <div className="auth-container">
      <h2>🏢 부동산 업자 회원가입</h2>
      <form onSubmit={handleSubmit}>
      <input
         type="text"
         name="name"
         value={form.name}
         onChange={handleChange}
         placeholder="이름"
         required
       />
       <input
         type="text"
         name="licenseNum"
         value={form.licenseNum}
         onChange={handleChange}
         placeholder="자격증 번호 (예: 서울-123456)"
         required
       />
       <input
         type="tel"
        name="phone"
         value={form.phone}
         onChange={handleChange}
         placeholder="전화번호"
         required
       />
        <input type="email" name="email" value={form.email} onChange={handleChange} placeholder="이메일" required />
        <input type="password" name="password" value={form.password} onChange={handleChange} placeholder="비밀번호" required />
        <input type="password" name="confirmPassword" value={form.confirmPassword} onChange={handleChange} placeholder="비밀번호 확인" required />
        <button type="submit">회원가입</button>
      </form>
    </div>
  );
};

export default SignupAgent;
