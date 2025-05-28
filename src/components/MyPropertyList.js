import React, { useEffect, useState } from 'react';
import axios from 'axios';

const API_URL = 'https://252f-219-251-84-31.ngrok-free.app';

const MyPropertyList = ({ user }) => {
  const [properties, setProperties] = useState([]);

  const fetchMyProperties = async () => {
    try {
      const res = await axios.get(`${API_URL}/properties?user=${user.id}`, {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      });
      setProperties(res.data);
    } catch (err) {
      console.error('❌ 내 매물 조회 실패:', err.response || err.message);
    }
  };

  useEffect(() => {
    if (user?.id) {
      fetchMyProperties();
    }
  }, [user]);

  return (
    <div className="property-list">
      {properties.length === 0 ? (
        <p>📭 등록한 매물이 없습니다.</p>
      ) : (
        <div className="property-grid">
          {properties.map((property, index) => (
            <div key={index} className="property-card">
              <h4>{property.address}</h4>
              <p>💰 가격: {property.price.toLocaleString()} 원</p>
              <p>👤 소유자: {property.owner}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default MyPropertyList;
