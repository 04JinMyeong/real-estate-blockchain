// src/components/PropertyForm.js
import React, { useState } from 'react';
import axios from 'axios';
import './PropertyForm.css';

// 반드시 현재 동작 중인 ngrok 주소(또는 배포된 API 주소)로 바꿔주세요.
const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';

export default function PropertyForm({ user, onRegister }) {
  const [form, setForm] = useState({
    address: '',
    owner: '',
    price: ''
  });
  const [photoFile, setPhotoFile] = useState(null);
  const [uploading, setUploading] = useState(false);

  // 입력값 변경 핸들러
  const handleChange = e => {
    const { name, value } = e.target;
    setForm(prev => ({ ...prev, [name]: value }));
  };

  // 파일 선택 핸들러
  const handleFileChange = e => {
    if (e.target.files && e.target.files[0]) {
      console.log('📄 선택된 파일:', e.target.files[0]);
      setPhotoFile(e.target.files[0]);
    }
  };

  // 폼 제출 핸들러
  const handleSubmit = async e => {
    e.preventDefault();

    // 1) 필수 항목 검사
    if (!form.address.trim() || !form.owner.trim() || !form.price) {
      alert('❗ 주소, 소유자, 가격을 모두 입력해주세요.');
      return;
    }

    setUploading(true);

    try {
      let photoUrl = '';

      // ─── 2) 사진 업로드 단계 (photoFile이 있을 때만) ───
      if (photoFile) {
        console.log('👉 1단계: 사진 업로드 시작');
        console.log('   📄 업로드할 photoFile:', photoFile);

        const photoData = new FormData();
        // 백엔드가 c.FormFile("photo")로 받으므로, key를 반드시 "photo"로 맞춥니다.
        photoData.append('photo', photoFile, photoFile.name);

        // FormData 내부 확인 (디버깅용)
        for (let [key, value] of photoData.entries()) {
          console.log(`   📥 FormData entry -> ${key}:`, value);
        }

        // Content-Type 헤더를 직접 지정하지 않으면 Axios가 boundary를 자동으로 붙여줍니다.
        const photoRes = await axios.post(
          `${API_URL}/upload-photo`,
          photoData
        );

        // 서버에서 { "photoUrl": "http://…/uploads/xxx.jpg" } 형태로 내려옴
        photoUrl = photoRes.data.photoUrl;
        console.log('✅ 사진 업로드 성공, photoUrl:', photoUrl);
      }

      // ─── 3) 매물 등록 단계 ───
      console.log('👉 2단계: 매물 등록 시작 → payload:', {
        user:     user.username || user.id,
        address:  form.address,
        owner:    form.owner,
        price:    form.price,
        photoUrl  // upload-photo에서 받은 URL (없으면 빈 문자열)
      });

      const addRes = await axios.post(
        `${API_URL}/add-property`,
        {
          user:     user.username || user.id,
          address:  form.address,
          owner:    form.owner,
          price:    form.price,
          photoUrl
        }
      );

      console.log('✅ 매물 등록 성공:', addRes.data);
      alert('🏠 매물이 성공적으로 등록되었습니다.');

      // 폼 초기화
      setForm({ address: '', owner: '', price: '' });
      setPhotoFile(null);
      onRegister();
    } catch (err) {
      console.error(
        '❌ 매물 등록 오류:',
        err.response?.data || err.message
      );
      alert(
        '⚠️ 오류가 발생했습니다:\n' +
        (err.response?.data?.error ||
         err.response?.data?.detail ||
         err.message)
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
