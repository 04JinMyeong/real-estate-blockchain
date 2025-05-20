import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import nacl from 'tweetnacl'; // Ed25519 키 생성용
import { encodeBase64 } from 'tweetnacl-util'; // Base64 인코딩용
import './Auth.css'; // 기존 CSS 파일

const BASE_API_URL = 'http://localhost:8080'; // ⬅️ 로컬 백엔드 서버 주소로 변경!;
const DID_SIGNUP_API_ENDPOINT = `${BASE_API_URL}/api/brokers/register-with-did`; // 새 엔드포인트

const SignupAgent = () => {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    platform_username: '', // 백엔드 SignUpBrokerRequest의 PlatformUsername 필드와 매칭
    platform_password: '', // 백엔드 SignUpBrokerRequest의 PlatformPassword 필드와 매칭
    confirmPassword: '',   // 비밀번호 확인용
    email: '',
    full_name: '',       // 백엔드 SignUpBrokerRequest의 FullName 필드와 매칭
    license_number: '',  // 백엔드 SignUpBrokerRequest의 LicenseNumber 필드와 매칭
    office_address: ''   // 백엔드 SignUpBrokerRequest의 OfficeAddress 필드와 매칭
  });
  const [generatedPrivateKey, setGeneratedPrivateKey] = useState(''); // 생성된 개인키 임시 저장 (사용자 안내용)

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();
    setGeneratedPrivateKey(''); // 이전 개인키 정보 초기화

    if (form.platform_password !== form.confirmPassword) {
      alert('❌ 비밀번호가 일치하지 않습니다.');
      return;
    }

    // 필수 입력 필드 검증 (예시)
    if (!form.platform_username || !form.platform_password || !form.email || !form.full_name || !form.license_number) {
      alert('❌ 필수 입력 항목을 모두 채워주세요.');
      return;
    }

    try {
      // --- DID 발급을 위한 키 쌍 생성 (Ed25519) ---
      const keyPair = nacl.sign.keyPair(); // { publicKey: Uint8Array, secretKey: Uint8Array }
      const publicKeyBytes = keyPair.publicKey;
      const privateKeyBytes = keyPair.secretKey;

      // 공개키를 Base64 문자열로 인코딩 (서버 전송용)
      const agentPublicKeyBase64 = encodeBase64(publicKeyBytes);

      // 개인키를 Base64 문자열로 인코딩 (사용자에게 안전하게 전달/보관 안내용)
      const agentPrivateKeyBase64 = encodeBase64(privateKeyBytes);
      setGeneratedPrivateKey(agentPrivateKeyBase64); // 상태에 저장하여 사용자에게 보여줄 준비

      console.log("프론트엔드에서 생성된 공개키 (Base64):", agentPublicKeyBase64);
      // console.log("프론트엔드에서 생성된 개인키 (Base64):", agentPrivateKeyBase64); // 개발 중 확인용, 실제로는 바로 사용자에게 안전하게 안내

      // 백엔드로 전송할 데이터 구성
      const registrationData = {
        platform_username: form.platform_username,
        platform_password: form.platform_password,
        email: form.email,
        full_name: form.full_name,
        license_number: form.license_number,
        office_address: form.office_address,
        agent_public_key: agentPublicKeyBase64 // 백엔드가 받을 필드명 agent_public_key
      };

      // Axios를 사용하여 백엔드 API 호출
      const res = await axios.post(DID_SIGNUP_API_ENDPOINT, registrationData, {
        headers: {
          'Content-Type': 'application/json',
          'ngrok-skip-browser-warning': 'true' // ngrok 사용 시 필요할 수 있음
        }
      });

      // 성공 시 서버 응답에서 DID 정보 확인
      const issuedDID = res.data.did; // 백엔드 응답 형식에 따라 조정
      alert(`✅ 회원가입 및 DID 발급 성공!\n발급된 DID: ${issuedDID}\n\n[매우 중요]\n아래 표시된 개인키를 즉시 복사하여 안전한 곳에 보관하십시오.\n이 개인키는 분실 시 복구할 수 없으며, 서비스 이용에 반드시 필요합니다.\n\n개인키 (Base64):\n${agentPrivateKeyBase64}`);
      // 실제 서비스에서는 alert 대신 모달 창 등을 사용하여 개인키를 보여주고
      // 복사 버튼, 파일 저장 기능 등을 제공하는 것이 훨씬 안전하고 사용자 친화적입니다.

      navigate('/'); // 회원가입 성공 후 리디렉션 (예: 로그인 페이지 또는 대시보드)

    } catch (err) {
      console.error('❌ 회원가입 또는 DID 발급 중 에러:', err.response || err);
      let errorMessage = '회원가입 또는 DID 발급에 실패했습니다.';
      if (err.response && err.response.data && err.response.data.error) {
        errorMessage += `\n사유: ${err.response.data.error}`;
      } else if (err.message) {
        errorMessage += `\n사유: ${err.message}`;
      }
      alert(errorMessage);
      setGeneratedPrivateKey(''); // 오류 발생 시 개인키 정보도 초기화
    }
  };

  return (
    <div className="auth-container">
      <h2>🏢 부동산업자 회원가입 (DID 발급)</h2>
      <form onSubmit={handleSubmit}>
        <input
          name="platform_username" // form 상태의 키와 일치
          type="text"
          placeholder="사용자 아이디 (플랫폼 로그인용)"
          value={form.platform_username}
          onChange={handleChange}
          required
        />
        <input
          name="email"
          type="email"
          placeholder="이메일 주소"
          value={form.email}
          onChange={handleChange}
          required
        />
        <input
          name="platform_password" // form 상태의 키와 일치
          type="password"
          placeholder="비밀번호"
          value={form.platform_password}
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
        <hr />
        <input
          name="full_name" // form 상태의 키와 일치
          type="text"
          placeholder="이름 (예: 홍길동)"
          value={form.full_name}
          onChange={handleChange}
          required
        />
        <input
          name="license_number" // form 상태의 키와 일치
          type="text"
          placeholder="공인중개사 면허번호"
          value={form.license_number}
          onChange={handleChange}
          required
        />
        <input
          name="office_address" // form 상태의 키와 일치
          type="text"
          placeholder="사무실 주소 (선택 입력)"
          value={form.office_address}
          onChange={handleChange}
        />
        <button type="submit">회원가입 및 DID 발급</button>
      </form>
      {generatedPrivateKey && (
        <div className="private-key-display">
          <h3>[중요] 발급된 개인키 (안전하게 보관하세요!):</h3>
          <textarea value={generatedPrivateKey} readOnly rows="4" cols="50" />
          <button onClick={() => navigator.clipboard.writeText(generatedPrivateKey)}>
            개인키 복사
          </button>
          <p style={{ color: 'red', fontWeight: 'bold' }}>
            이 개인키는 다시 표시되지 않습니다. 반드시 지금 복사하여 안전한 곳에 저장하세요.
          </p>
        </div>
      )}
    </div>
  );
};

export default SignupAgent;