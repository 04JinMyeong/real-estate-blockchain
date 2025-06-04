// src/components/PropertyList.js
import React, {
  useEffect,
  useState,
  forwardRef,
  useImperativeHandle,
  useCallback
} from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
//import './PropertyList.css';

const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';

const PropertyList = forwardRef(({ user, mode = 'all', onReserve }, ref) => {
  const [properties, setProperties] = useState([]);
  const navigate = useNavigate();

  const fetchProperties = useCallback(async () => {
    try {
      let url = mode === 'my'
        ? (user?.username
            ? `${API_URL}/my-properties?user=${user.username}`
            : '')
        : `${API_URL}/properties?user=admin`;

      if (!url) {
        setProperties([]);
        return;
      }

      const res = await axios.get(url, {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      });
      const arr = Array.isArray(res.data.properties)
        ? res.data.properties
        : Array.isArray(res.data)
        ? res.data
        : [];
      setProperties(arr);
    } catch (err) {
      setProperties([]);
      console.error('❌ 매물 조회 실패:', err.response?.data || err.message);
    }
  }, [user, mode]);

  useImperativeHandle(ref, () => ({ fetchProperties }));

  useEffect(() => {
    fetchProperties();
  }, [fetchProperties]);

  const handleReserve = async (property) => {
  if (!user?.username) {
    alert('로그인 후 예약 가능합니다.');
    return;
  }
  // expiresAt: 12시간 뒤(초 단위 Unix 타임)
  const expiresAt = Math.floor(Date.now() / 1000) + 12 * 3600;

  const payload = {
    user: user.username,
    id: property.id,
    expiresAt
  };
  const token = localStorage.getItem('token');
  try {
    const res = await axios.post(
      `${API_URL}/reserve-property`,
      payload,
      {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
          'ngrok-skip-browser-warning': 'true'
        }
      }
    );
    // 성공 메시지/유효기한 응답
    alert(res.data.message || '✅ 예약이 완료되었습니다.');
    if (res.data.expiresAt) {
      const date = new Date(res.data.expiresAt * 1000);
      alert(
        '✅ 예약 유효기한: ' +
        date.toLocaleString('ko-KR', { hour12: false })
      );
    }
    fetchProperties();
    onReserve?.();
  } catch (err) {
    console.error('❌ reserve-property error response:', err.response);
    const msg =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message;
    alert('예약 실패: ' + msg);
  }
};


  return (
    <div className="property-list">
      {properties.length === 0 ? (
        <p>
          {mode === 'my'
            ? '📭 등록한 매물이 없습니다.'
            : '📭 등록된 매물이 없습니다.'}
        </p>
      ) : (
        <div className="property-grid">
          {properties.map((p, i) => (
            <div key={p.id || i} className="property-card">
              <h4>{p.address}</h4>
              <p>👤 소유자: {p.ownerHistory?.slice(-1)[0]?.owner || '-'}</p>
              <p>
                💰 가격:{' '}
                {p.priceHistory?.slice(-1)[0]?.price.toLocaleString() || '-'}원
              </p>
              <div className="property-actions" style={{ display: 'flex', gap: '8px' }}>
                {/* 상세보기 버튼을 추가했습니다. */}
                <button onClick={() => navigate(`/properties/${p.id}`)}>
                  상세보기
                </button>

                {p.reservedBy ? (
                  <span style={{ color: 'red', lineHeight: '32px' }}>예약됨</span>
                ) : (
                  <button onClick={() => handleReserve(p)}>
                    예약하기
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
});

export default PropertyList;
