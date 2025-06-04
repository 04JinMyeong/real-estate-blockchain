// src/components/PropertyDetail.js
import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import axios from 'axios';
import './PropertyDetail.css';

// (1) 정적 파일 라우터 기본 URL
//     실제 서버 주소/포트나 ngrok 주소로 변경해주세요.
const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';

/**
 * UTC 문자열을 KST(한국 시간) 형식으로 바꿔줍니다.
 * 백엔드가 "YYYY-MM-DD HH:mm:ss" 형식으로 내려준다고 가정.
 */
function toKST(utcStr) {
  if (!utcStr) return '-';
  // "YYYY-MM-DD HH:mm:ss" → "YYYY-MM-DDTHH:mm:ssZ" 형태로 변환 (UTC로 간주)
  const iso = utcStr.replace(' ', 'T') + 'Z';
  const d = new Date(iso);

  return d.toLocaleString('ko-KR', {
    year:   'numeric',
    month:  '2-digit',
    day:    '2-digit',
    hour:   '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
    timeZone: 'Asia/Seoul'
  });
}

export default function PropertyDetail() {
  const { id } = useParams();
  const [property, setProperty] = useState(null);
  const [history,  setHistory]  = useState([]);
  const [loading,  setLoading]  = useState(true);

  useEffect(() => {
    // 1) 매물 단건 조회 (admin 고정)
    axios
      .get(`${API_URL}/property/${id}?user=admin`, {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      })
      .then(res => {
        let p = res.data.property;
        if (typeof p === 'string') {
          try { p = JSON.parse(p); } catch {}
        }
        console.log('🛰️ 서버에서 내려준 property →', p);
        setProperty(p);
      })
      .catch(err => {
        console.error('❌ 상세 조회 실패:', err.response?.data || err.message);
        setProperty(null);
      });

    // 2) 이력 조회
    axios
      .get(`${API_URL}/property/history?id=${encodeURIComponent(id)}&user=admin`, {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      })
      .then(res => {
        let h = res.data.history || res.data;
        if (typeof h === 'string') {
          try { h = JSON.parse(h); } catch {}
        }
        console.log('🛰️ 서버에서 내려준 history →', h);
        setHistory(h);
      })
      .catch(err => {
        console.error('❌ 이력 조회 실패:', err.response?.data || err.message);
        setHistory([]);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [id]);

  if (loading) return <p>로딩 중…</p>;
  if (!property) return <p>해당 매물을 찾을 수 없습니다.</p>;

  // —————————————————————————
  // 3) 이미지 URL 결정
const rawImg = property.photoUrl || '';
let imgSrc = null;

if (rawImg) {
  if (rawImg.startsWith('https://2094-165-229-229-106.ngrok-free.app') || rawImg.startsWith('https://2094-165-229-229-106.ngrok-free.app')) {
    imgSrc = rawImg.replace('https://2094-165-229-229-106.ngrok-free.app', API_URL).replace('https://2094-165-229-229-106.ngrok-free.app', API_URL);
  } else if (rawImg.startsWith(API_URL)) {
    imgSrc = rawImg;
  } else if (rawImg.startsWith('http://') || rawImg.startsWith('https://')) {
    imgSrc = rawImg;
  } else {
    let normalized = rawImg.replace(/\\/g, '/');
    if (!normalized.startsWith('/')) normalized = '/' + normalized;
    imgSrc = `${API_URL}${normalized}`;
  }
}
console.log('🖼️ 실제로 사용하는 imgSrc →', imgSrc);
  // —————————————————————————

  // 4) 가격 / 소유자 / 등록일(UTC→KST) / 상태
  const lastPriceEntry = property.priceHistory?.slice(-1)[0] || {};
  const priceNum = lastPriceEntry.price
    ? Number(lastPriceEntry.price).toLocaleString()
    : '-';

  const regUTC = lastPriceEntry.date || lastPriceEntry.timestamp || '';
  const regDateKST = regUTC ? toKST(regUTC) : '-';

  return (
    <div
      className="property-detail-card"
      style={{
        maxWidth: 800,
        margin: '2rem auto',
        padding: '1.5rem',
        background: '#fff',
        borderRadius: 8,
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}
    >
      <Link
        to="/map"
        style={{
          textDecoration: 'none',
          color: '#555',
          marginBottom: '1rem',
          display: 'inline-block'
        }}
      >
        ← 목록으로 돌아가기
      </Link>

      {/* — 이미지 영역 — */}
      <div
        style={{
          width: '100%',
          height: 0,
          paddingBottom: '56.25%', // 16:9 비율
          position: 'relative',
          marginBottom: '1rem',
          border: '2px dashed #ccc',
          borderRadius: 8,
          overflow: 'hidden',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center'
        }}
      >
        {imgSrc ? (
          <img
            src={imgSrc}
            alt="매물 사진"
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: '100%',
              objectFit: 'cover'
            }}
          />
        ) : (
          <div style={{ color: '#888', fontSize: 16 }}>
            <span role="img" aria-label="camera" style={{ marginRight: 8 }}>
              📷
            </span>
            사진이 없습니다.
          </div>
        )}
      </div>

      {/* — 기본 정보 — */}
      <h2 style={{ marginBottom: '1rem', fontSize: 24, color: '#333' }}>
        {property.address}
      </h2>

      <div
        className="detail-info"
        style={{
          display: 'grid',
          gridTemplateColumns: 'auto auto',
          gap: '1rem 2rem',
          marginBottom: '1rem'
        }}
      >
        <div>
          <strong style={{ display: 'block', marginBottom: 4 }}>💰 가격</strong>
          <p style={{ margin: 0, fontSize: 18, color: '#444' }}>
            {priceNum}원
          </p>
        </div>
        <div>
          <strong style={{ display: 'block', marginBottom: 4 }}>👤 소유자</strong>
          <p style={{ margin: 0, fontSize: 18, color: '#444' }}>
            {property.ownerHistory?.slice(-1)[0]?.owner || '-'}
          </p>
        </div>
        <div>
          <strong style={{ display: 'block', marginBottom: 4 }}>
            📅 등록일 (KST)
          </strong>
          <p style={{ margin: 0, fontSize: 18, color: '#444' }}>
            {regDateKST}
          </p>
        </div>
        <div>
          <strong style={{ display: 'block', marginBottom: 4 }}>상태</strong>
          <p
            style={{
              margin: 0,
              fontSize: 18,
              color: property.reservedBy ? 'red' : '#444'
            }}
          >
            {property.reservedBy ? '예약됨' : '예약 가능'}
          </p>
        </div>
      </div>

      {/* — 설명 — */}
      <div className="description" style={{ marginBottom: '1rem' }}>
        <strong style={{ display: 'block', marginBottom: 4 }}>📝 설명</strong>
        <p style={{ margin: 0, color: '#666' }}>
          {property.description || '설명 없음'}
        </p>
      </div>

      {/* — 이력 (History) — */}
      {history.length > 0 && (
        <div className="history" style={{ marginTop: '1.5rem' }}>
          <strong style={{ display: 'block', marginBottom: 8 }}>
            🕒 매물 이력
          </strong>
          <ul
            style={{
              margin: 0,
              paddingLeft: '1.2rem',
              color: '#555',
              lineHeight: 1.6
            }}
          >
            {history.map((h, idx) => {
              const entry = typeof h === 'string' ? JSON.parse(h) : h;
              const dtUTC = entry.date || entry.timestamp || '';
              const dateStr = dtUTC ? toKST(dtUTC) : '-';
              const lp = entry.priceHistory?.slice(-1)[0] || {};
              const pPrice = lp.price
                ? Number(lp.price).toLocaleString()
                : '-';
              const owner = entry.ownerHistory?.slice(-1)[0]?.owner || '-';

              return (
                <li key={idx}>
                  {dateStr} : 가격 {pPrice}원, 소유자 {owner}
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
}
