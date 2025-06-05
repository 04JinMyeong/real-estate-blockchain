// src/components/Login.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';  // named export
import axios from 'axios';
import './Auth.css';

//const AUTH_API = 'https://2094-165-229-229-106.ngrok-free.app/login';
const AUTH_API = 'http://localhost:8080/login';

export default function Login({ onLogin }) {
  const [form, setForm] = useState({ id: '', password: '', vc: '' });
  const [vcFileName, setVcFileName] = useState(''); // 업로드된 VC 파일 이름 표시용 상태
  const navigate = useNavigate();

  const handleChange = e => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };


  // 2. VC 파일 선택 시 파일 내용을 읽어 form.vc 상태에 저장하는 함수
  const handleVCFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setVcFileName(file.name);
      const reader = new FileReader();
      reader.onload = (event) => {
        try {
          const fileContent = event.target.result;
          setForm(prev => ({ ...prev, vc: fileContent }));
          console.log('[Login.js] VC file content loaded into form state. Length:', fileContent.length);
        } catch (error) {
          console.error('[Login.js] Error processing VC file content:', error);
          alert('VC 파일을 처리하는 중 문제가 발생했습니다.');
          setForm(prev => ({ ...prev, vc: '' }));
          setVcFileName('');
          if (e.target) e.target.value = null; // 파일 선택 input 초기화 (선택적)
        }
      };
      reader.onerror = (error) => {
        console.error('[Login.js] Error reading VC file:', error);
        alert('VC 파일을 읽는 중 오류가 발생했습니다.');
        setForm(prev => ({ ...prev, vc: '' }));
        setVcFileName('');
        if (e.target) e.target.value = null; // 파일 선택 input 초기화 (선택적)
      };
      reader.readAsText(file);
    } else {
      console.log('[Login.js] No VC file selected or selection cancelled.');
      setForm(prev => ({ ...prev, vc: '' }));
      setVcFileName('');
    }
  };


  // 변경된 handleSubmit (제가 제안드린 버전)
  const handleSubmit = async e => {
    e.preventDefault();
    console.log('[Login.js] handleSubmit 함수 실행됨.'); // 함수 실행 시작 로그

    // 아이디, 비밀번호 입력 여부 확인
    if (!form.id || !form.password) {
      alert('아이디와 비밀번호를 입력해주세요.');
      return;
    }
    // VC 파일 첨부 여부 확인 (필수라고 가정)
    if (!form.vc) {
      alert('VC 파일을 선택해주세요.');
      console.log('[Login.js] handleSubmit: VC 파일 누락됨.');
      return;
    }

    // 백엔드로 보낼 데이터(페이로드) 구성
    const payload = {
      id: form.id,
      password: form.password,
      vc: form.vc, // ⬅️ VC 정보(파일에서 읽은 문자열) 포함
    };
    // VC 내용이 길 수 있으므로, 일부 또는 길이만 로깅
    console.log('[Login.js] 👉 로그인 요청 데이터(페이로드). ID:', payload.id, ', VC 길이:', payload.vc ? payload.vc.length : 0);

    try {
      console.log(`[Login.js] API 요청 시도: POST ${AUTH_API}`);
      const res = await axios.post(AUTH_API, payload, { // payload 변수 사용
        headers: {
          'Content-Type': 'application/json',
        }
      });

      console.log('[Login.js] ✅ 로그인 성공 응답 상태 코드:', res.status);
      console.log('[Login.js] ✅ 로그인 성공 응답 데이터:', res.data);
      // 백엔드 응답에서 did도 받을 수 있도록 가정하고 추가
      const { token, role, user: userId, did } = res.data;

      localStorage.setItem('token', token);

      // onLogin 콜백에 전달할 사용자 정보에 did도 포함
      let userInfoToPass = {
        email: form.id,
        role: role,
        username: userId,
        did: did, // did 추가
      };

      try {
        const decoded = jwtDecode(token);
        console.log('[Login.js] 🔓 JWT 디코딩 완료:', decoded);
        // JWT 클레임에 따라 userInfoToPass 업데이트 가능
        if (decoded.user_id) userInfoToPass.username = decoded.user_id;
        if (decoded.did) userInfoToPass.did = decoded.did;
      } catch (err) {
        console.warn('[Login.js] ⚠️ JWT 디코딩 실패:', err);
      }

      onLogin(userInfoToPass);
      console.log('[Login.js] onLogin 콜백 함수 실행됨.');

      if (role === 'agent') {
        navigate('/agent/mypage');
      } else {
        navigate(role === 'user' ? '/user/mypage' : '/');
      }
      console.log('[Login.js] 페이지 이동 완료.');

    } catch (err) {
      // 에러 로깅 상세화
      console.error('[Login.js] ❌ 로그인 API 호출 실패.');
      if (err.response) {
        console.error('[Login.js] ❌ 로그인 오류 상태 코드:', err.response.status);
        console.error('[Login.js] ❌ 로그인 오류 응답 데이터:', err.response.data);
        console.error('[Login.js] ❌ 로그인 오류 응답 헤더:', err.response.headers);
      } else if (err.request) {
        console.error('[Login.js] ❌ 로그인 요청에 대한 응답 없음:', err.request);
      } else {
        console.error('[Login.js] ❌ 로그인 요청 설정 중 오류:', err.message);
      }
      const errorMessage = err.response?.data?.error || err.response?.data?.message || err.message || '로그인에 실패했습니다. 네트워크 연결 또는 서버 상태를 확인해주세요.';
      alert(errorMessage);
    }
  };

  return (
    <div className="auth-container">
      <form onSubmit={handleSubmit} className="auth-form">
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
        <div style={{ marginTop: '15px', marginBottom: '15px' }}>
          <label htmlFor="vcFile" style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
            VC 파일 첨부:
          </label>
          <input
            type="file"
            id="vcFile"
            name="vcFile" // 이 name은 form 상태와 직접 연결되지 않음
            accept=".json" // JSON 파일만 선택하도록 유도
            onChange={handleVCFileChange}
            style={{ width: '100%', padding: '8px', boxSizing: 'border-box' }}
            required // VC 파일 첨부를 필수로 만듦
          />
          {vcFileName && (
            <p style={{ marginTop: '5px', fontSize: '0.9em', color: '#555' }}>
              선택된 파일: {vcFileName}
            </p>
          )}
        </div>
        <button type="submit">로그인</button>
      </form>
    </div>
  );
}
