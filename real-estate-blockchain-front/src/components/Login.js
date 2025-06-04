// 📄 src/components/Login.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';  // named export
import axios from 'axios';
import './Auth.css'; // Login.js와 동일한 폴더에 Auth.css가 있다고 가정

// 🚨🚨🚨 사용자님의 실제 백엔드 로그인 API 주소로 수정해주세요! 🚨🚨🚨
const AUTH_API = 'http://localhost:8080/login'; // 로컬 백엔드 서버 주소

export default function Login({ onLogin }) {
  const [form, setForm] = useState({ id: '', password: '', vc: '' });
  const [vcFileName, setVcFileName] = useState('');
  const navigate = useNavigate();

  const handleChange = e => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  const handleVCFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setVcFileName(file.name);
      const reader = new FileReader();
      reader.onload = (event) => {
        try {
          const fileContent = event.target.result;
          setForm(prev => ({ ...prev, vc: fileContent }));
          console.log('[Login.js] VC file content loaded into form state. Length:', fileContent.length); // 로그 상세화
        } catch (jsonError) {
          console.error('[Login.js] Error parsing VC file as JSON (though not parsing here):', jsonError); // 로그 메시지 수정
          alert('VC 파일을 읽는 중 문제가 발생했습니다. 파일 형식을 확인해주세요.'); // 사용자 메시지 개선
          setForm(prev => ({ ...prev, vc: '' }));
          setVcFileName('');
          e.target.value = null;
        }
      };
      reader.onerror = (error) => {
        console.error('[Login.js] Error reading VC file:', error);
        alert('VC 파일을 읽는 중 오류가 발생했습니다.');
        setForm(prev => ({ ...prev, vc: '' }));
        setVcFileName('');
        e.target.value = null;
      };
      reader.readAsText(file);
    } else {
      console.log('[Login.js] No VC file selected or selection cancelled.'); // 로그 추가
      setForm(prev => ({ ...prev, vc: '' }));
      setVcFileName('');
    }
  };

  const handleSubmit = async e => {
    e.preventDefault();
    console.log('[Login.js] handleSubmit triggered.'); // 함수 시작 로그

    if (!form.vc) {
      alert('VC 파일을 선택해주세요.');
      console.log('[Login.js] handleSubmit: VC is missing.');
      return;
    }

    const payload = {
      id: form.id,
      password: form.password,
      vc: form.vc,
    };
    // VC 내용이 매우 길 수 있으므로, 전체를 로깅하는 대신 일부 또는 길이만 로깅하는 것을 고려
    console.log('[Login.js] 👉 Login payload. ID:', payload.id, ', VC length:', payload.vc.length);

    try {
      console.log(`[Login.js] Attempting to POST to: ${AUTH_API}`);
      const res = await axios.post(AUTH_API, payload, {
        headers: {
          'Content-Type': 'application/json',
        }
      });

      console.log('[Login.js] ✅ login success response status:', res.status); // 상태 코드 로그
      console.log('[Login.js] ✅ login success response data:', res.data); // 응답 데이터 전체 로그
      const { token, role, user: userId, did } = res.data;

      localStorage.setItem('token', token);

      let userInfoToPass = {
        email: form.id,
        role: role,
        username: userId,
        did: did,
      };

      try {
        const decoded = jwtDecode(token);
        console.log('[Login.js] 🔓 decoded JWT:', decoded);
      } catch (err) {
        console.warn('[Login.js] ⚠️ JWT decode failed:', err);
      }

      onLogin(userInfoToPass);
      console.log('[Login.js] onLogin callback executed.');

      if (role === 'agent') {
        navigate('/agent/mypage');
      } else {
        navigate('/user/mypage'); // 일반 사용자 페이지가 있다면
      }
      console.log('[Login.js] Navigated to mypage.');

    } catch (err) {
      console.error('[Login.js] ❌ Login API call failed.');
      if (err.response) {
        console.error('[Login.js] ❌ login error status:', err.response.status);
        console.error('[Login.js] ❌ login error response data:', err.response.data);
        console.error('[Login.js] ❌ login error response headers:', err.response.headers);
      } else if (err.request) {
        console.error('[Login.js] ❌ No response received for login request:', err.request);
      } else {
        console.error('[Login.js] ❌ Error setting up login request:', err.message);
      }
      const errorMessage = err.response?.data?.error || err.response?.data?.message || err.message || '로그인에 실패했습니다. 네트워크 연결 또는 서버 상태를 확인해주세요.';
      alert(errorMessage);
    }
  };

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
        <h2>로그인</h2>
        <input
          name="id"
          type="text"
          placeholder="아이디"
          value={form.id}
          onChange={handleChange}
          required
          autoComplete="username"
        />
        <input
          name="password"
          type="password"
          placeholder="비밀번호"
          value={form.password}
          onChange={handleChange}
          required
          autoComplete="current-password"
        />
        <div style={{ marginTop: '15px', marginBottom: '15px' }}>
          <label htmlFor="vcFile" style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
            VC 파일 첨부:
          </label>
          <input
            type="file"
            id="vcFile"
            name="vcFile"
            accept=".json"
            onChange={handleVCFileChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
          />
          {vcFileName && (
            <p style={{ marginTop: '5px', fontSize: '0.9em', color: '#555' }}>
              선택된 파일: {vcFileName}
            </p>
          )}
        </div>
        <button type="submit" style={{ marginTop: '15px' }}>로그인</button>
      </form>
    </div>
  );
}