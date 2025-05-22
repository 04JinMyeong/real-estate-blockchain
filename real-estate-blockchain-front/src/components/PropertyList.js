import React, { useEffect, useState, forwardRef, useImperativeHandle, useCallback } from 'react';
import axios from 'axios';

const API_URL = 'https://1af7-165-229-229-137.ngrok-free.app';

const PropertyList = forwardRef(({ user, mode = 'all', onReserve }, ref) => {
  const [properties, setProperties] = useState([]);

  const fetchProperties = useCallback(async () => {
    try {
      let url = `${API_URL}/properties`;
      if (mode === 'my' && user?.id) {
        url = `${API_URL}/properties?user=${user.id}`;
      }

      const res = await axios.get(url, {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      });

      setProperties(res.data);
    } catch (err) {
      console.error('❌ 매물 조회 실패:', err.response || err.message);
    }
  }, [user, mode]);

  useImperativeHandle(ref, () => ({
    fetchProperties,
  }));

  useEffect(() => {
    fetchProperties();
  }, [user, mode, fetchProperties]);

  // 예약 버튼 클릭 시
  const handleReserve = async (property) => {
    if (!user) {
      alert('로그인 후 예약 가능합니다.');
      return;
    }
    const token = localStorage.getItem('token');
    try {
      // 예약 API 요청
      await axios.post(
        `${API_URL}/api/reservations`,
        {
          propertyId: property.id,
          userId: user.id,
          date: new Date().toISOString().slice(0, 10), // 임시: 오늘 날짜
          time: "10:00",
          notes: ""
        },
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          }
        }
      );
      alert('✅ 예약이 완료되었습니다.');
      fetchProperties(); // 예약 후 리스트 갱신
      if (onReserve) onReserve(); // 내 예약목록도 새로고침
    } catch (err) {
      alert('예약 실패: ' + (err.response?.data?.message || err.message));
    }
  };

  return (
    <div className="property-list">
      {properties.length === 0 ? (
        <p>{mode === 'my' ? '📭 등록한 매물이 없습니다.' : '📭 등록된 매물이 없습니다.'}</p>
      ) : (
        <div className="property-grid">
          {properties.map((property, index) => (
            <div key={index} className="property-card">
              <h4>{property.address}</h4>
              <p>💰 가격: {property.price.toLocaleString()} 원</p>
              <p>👤 소유자: {property.owner}</p>
              <p>🆔 매물 ID: {property.id}</p>
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
