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
// import './PropertyList.css';

//const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';
const API_URL = 'http://localhost:8080'; // 로컬 개발용

// 남은 초 → "H:MM:SS" 포맷 변환 함수
function formatLeftTime(sec) {
  if (sec <= 0) return "0:00";
  const h = Math.floor(sec / 3600);
  const m = Math.floor((sec % 3600) / 60);
  const s = sec % 60;
  return h > 0
    ? `${h}:${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`
    : `${m}:${String(s).padStart(2, "0")}`;
}

const PropertyList = forwardRef(({ user, mode = 'all', onReserve }, ref) => {
  const [properties, setProperties] = useState([]);
  const [now, setNow] = useState(Date.now()); // 실시간 갱신용

  const navigate = useNavigate();

  // 1초마다 현재 시각 갱신 (남은 시간 표시용)
  useEffect(() => {
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(timer);
  }, []);

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
    // 12시간 뒤(초 단위)
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
          {properties.map((p, i) => {
            // 예약 남은 시간 계산
            let leftSeconds = null;
            if (p.reservedBy && p.expiresAt) {
              leftSeconds = p.expiresAt - Math.floor(now / 1000);
            }

            return (
              <div key={p.id || i} className="property-card">
                {/* --- 사진 미리보기 --- */}
                {p.photoUrl ? (
                  <div style={{
                    width: "100%",
                    height: 160,
                    marginBottom: 8,
                    overflow: "hidden",
                    borderRadius: 8,
                    background: "#eee",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center"
                  }}>
                    <img
                      src={p.photoUrl}
                      alt="매물사진"
                      style={{
                        width: "100%",
                        height: "100%",
                        objectFit: "cover",
                        display: "block"
                      }}
                      onError={e => { e.target.style.display = "none"; }}
                    />
                  </div>
                ) : (
                  <div style={{
                    width: "100%",
                    height: 160,
                    marginBottom: 8,
                    background: "#eee",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    color: "#bbb"
                  }}>
                    <span>사진 없음</span>
                  </div>
                )}

                <h4>{p.address}</h4>
                <p>👤 소유자: {p.ownerHistory?.slice(-1)[0]?.owner || '-'}</p>
                <p>
                  💰 가격:{' '}
                  {p.priceHistory?.slice(-1)[0]?.price?.toLocaleString() || '-'}원
                </p>
                <div className="property-actions" style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
                  <button onClick={() => navigate(`/properties/${p.id}`)}>
                    상세보기
                  </button>
                  {p.reservedBy && leftSeconds > 0 ? (
                    <>
                      <span style={{ color: 'red', lineHeight: '32px' }}>예약됨</span>
                      <span style={{ color: '#555', marginLeft: 8 }}>
                        남은 시간: {formatLeftTime(leftSeconds)}
                      </span>
                    </>
                  ) : (
                    <button onClick={() => handleReserve(p)}>
                      예약하기
                    </button>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
});

export default PropertyList;
