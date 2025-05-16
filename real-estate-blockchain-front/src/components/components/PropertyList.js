import React, { useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import axios from 'axios';

const API_URL = 'https://1af7-165-229-229-137.ngrok-free.app';

const PropertyList = forwardRef(({ user, mode = 'all' }, ref) => {
  const [properties, setProperties] = useState([]);

  const fetchProperties = async () => {
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
  };

  useImperativeHandle(ref, () => ({
    fetchProperties,
  }));

  useEffect(() => {
    fetchProperties();
  }, [user, mode]);

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
            </div>
          ))}
        </div>
      )}
    </div>
  );
});

export default PropertyList;
