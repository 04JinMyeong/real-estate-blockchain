import React, { useEffect, useState, forwardRef, useImperativeHandle, useCallback } from 'react';
import axios from 'axios';

const API_URL = 'https://252f-219-251-84-31.ngrok-free.app';

const PropertyList = forwardRef(({ user, mode = 'all', onReserve }, ref) => {
  const [properties, setProperties] = useState([]);

  const fetchProperties = useCallback(async () => {
    try {
      let url = '';
      if (mode === 'my') {
        if (!user || !user.username) {
          setProperties([]); // user 정보 없으면 빈 배열
          return;
        }
        url = `${API_URL}/my-properties?user=${user.username}`;
      } else {
        // 전체 매물: 반드시 실존 계정 넣어야 400/500 예방 (백엔드 특이 구현 대응)
        url = `${API_URL}/properties?user=TestUser13`;
      }

      const res = await axios.get(url, {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      });
      let arr = [];
      if (Array.isArray(res.data.properties)) {
        arr = res.data.properties;
      } else if (Array.isArray(res.data)) {
        arr = res.data;
      }
      setProperties(arr);
    } catch (err) {
      setProperties([]); // 에러시 빈 배열로
      console.error('❌ 매물 조회 실패:', err.response?.data || err.message);
    }
  }, [user, mode]);

  useImperativeHandle(ref, () => ({
    fetchProperties,
  }));

  useEffect(() => {
    fetchProperties();
  }, [user, mode, fetchProperties]);

  const handleReserve = async (property) => {
    console.log('PropertyList user prop:', user);
    console.log('reserve payload:', { user: user.username, id: property.id });
    if (!user || !user.username) {
      alert('로그인 후 예약 가능합니다.');
      return;
    }
    const token = localStorage.getItem('token');
    try {
      await axios.post(
        `${API_URL}/reserve-property`,
        {
          user: user.username,
          id: property.id,
          
        },
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          }
        }
      );
      alert('✅ 예약이 완료되었습니다.');
      fetchProperties();
      if (onReserve) onReserve();
    } catch (err) {
      alert('예약 실패: ' + (err.response?.data?.error || err.response?.data?.message || err.message));
    }
  };

  return (
    <div className="property-list">
      {Array.isArray(properties) && properties.length === 0 ? (
        <p>{mode === 'my' ? '📭 등록한 매물이 없습니다.' : '📭 등록된 매물이 없습니다.'}</p>
      ) : (
        <div className="property-grid">
          {Array.isArray(properties) && properties.map((property, index) => (
            <div key={property.id || index} className="property-card">
              <h4>{property.address}</h4>
              <p>
                👤 소유자:{' '}
                {
                  Array.isArray(property.ownerHistory) && property.ownerHistory.length > 0
                    ? property.ownerHistory[property.ownerHistory.length - 1].owner
                    : '-'
                }
              </p>
              <p>
                💰 가격:{' '}
                {
                  Array.isArray(property.priceHistory) && property.priceHistory.length > 0
                    ? Number(property.priceHistory[property.priceHistory.length - 1].price).toLocaleString()
                    : '-'
                } 원
              </p>
              <p>
                {property.reservedBy ? (
                  <span style={{ color: 'red' }}>예약됨</span>
                ) : (
                  <button
                    onClick={() => handleReserve(property)}
                    disabled={!!property.reservedBy}
                  >
                    예약하기
                  </button>
                )}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
});

export default PropertyList;
