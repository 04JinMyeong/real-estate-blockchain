// src/components/PropertyForm.js
import React, { useState } from 'react';
import axios from 'axios';
import './PropertyForm.css';

// 실제 백엔드 주소(ngrok 등)로 교체하세요.
//const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';
const API_URL = 'http://localhost:8080'; // 로컬 개발용

export default function PropertyForm({ user, onRegister }) {
  console.log('▶ PropertyForm received user:', user);
  // 1) 매물 정보: 주소 · 소유자 · 가격
  const [form, setForm] = useState({
    address: '',
    owner: '',
    price: ''
  });

  // 2) 선택된 사진 파일
  const [photoFile, setPhotoFile] = useState(null);

  // 3) 업로드 중 표시
  const [uploading, setUploading] = useState(false);

  // 입력값 핸들러
  const handleChange = e => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  // 파일 선택 핸들러
  const handleFileChange = e => {
    if (e.target.files && e.target.files[0]) {
      setPhotoFile(e.target.files[0]);
    }
  };

  // 폼 제출
  const handleSubmit = async e => {
    e.preventDefault();

    // 유효성 검사
    if (!form.address.trim() || !form.owner.trim() || !form.price) {
      alert('❗ 주소, 소유자, 가격을 모두 입력해주세요.');
      return;
    }

    setUploading(true);

    try {
      let finalPhotoUrl = '';

      // 1단계: 사진 업로드 (photoFile이 있을 때만)
      if (photoFile) {
        console.log('👉 1단계: 사진 업로드 시작');

        // **중요: 백엔드는 `photo`라는 필드명(FormFile("photo"))으로 가져갑니다.**
        const photoData = new FormData();
        photoData.append('photo', photoFile);
        // ^^^ 여기 key를 "photo"로 반드시 맞춰야 합니다.

        const photoRes = await axios.post(
          `${API_URL}/upload-photo`,
          photoData,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
              // ngrok 경고가 필요하면 여기에 추가: 'ngrok-skip-browser-warning': 'true'
            }
          }
        );

        // 백엔드에서 { "photoUrl": "http://localhost:8080/uploads/파일명.jpg" } 형태로 내려준다고 가정
        finalPhotoUrl = photoRes.data.photoUrl;
        console.log('✅ 사진 업로드 성공, photoUrl:', finalPhotoUrl);
      }

      // 2단계: 매물 등록
      console.log('👉 2단계: 매물 등록 시작');
      console.log('   payload →', {
        user: user?.username || user?.id,
        address: form.address,
        owner: form.owner,
        price: form.price,
        photoUrl: finalPhotoUrl
      });

      const addRes = await axios.post(
        `${API_URL}/add-property`,
        {
          user: user?.username,     // ← 여기를 꼭 추가!
          address: form.address,
          owner: form.owner,
          price: form.price,
          photoUrl: finalPhotoUrl
        },
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      console.log('✅ 매물 등록 성공:', addRes.data);
      alert('🏠 매물이 성공적으로 등록되었습니다.');

      // 입력 초기화
      setForm({ address: '', owner: '', price: '' });
      setPhotoFile(null);

      // 부모에게 등록 완료 알리고 목록 새로고침 요청
      onRegister();
    } catch (err) {
      console.error('❌ 매물 등록 오류:', err.response?.data || err.message);
      alert(
        '⚠️ 등록 중 오류가 발생했습니다.\n' +
        (err.response?.data?.error || err.response?.data?.detail || err.message)
      );
    } finally {
      setUploading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="property-form">
      <div className="form-group">
        <label>주소</label>
        <input
          name="address"
          type="text"
          placeholder="예) 서울특별시 강남구 테헤란로 123"
          value={form.address}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label>소유자</label>
        <input
          name="owner"
          type="text"
          placeholder="예) 홍길동"
          value={form.owner}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label>가격 (원)</label>
        <input
          name="price"
          type="number"
          placeholder="예) 500000000"
          value={form.price}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label>매물 사진 (선택)</label>
        {/* 파일 업로드용 input: accept 이미지 포맷 */}
        <input
          type="file"
          accept="image/*"
          onChange={handleFileChange}
        />
      </div>

      <button
        type="submit"
        className="btn-submit"
        disabled={uploading}
      >
        {uploading ? '등록 중…' : '매물 등록'}
      </button>
    </form>
  );
}
