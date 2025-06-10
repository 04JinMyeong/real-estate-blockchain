import React, { useEffect, useState, useRef } from 'react';
import PropertyList from './PropertyList';
import KakaoMap from './KakaoMap';
import './MapView.css';

//const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';
const API_URL = 'http://localhost:8080'; // 로컬 개발용

const MapView = ({ user }) => {
  const [properties, setProperties] = useState([]);
  const propertyListRef = useRef();

  // 1️⃣ PropertyList와 동일한 fetchProperties 함수
  const fetchProperties = async () => {
    try {
      const res = await fetch(`${API_URL}/properties?user=admin`);
      const data = await res.json();
      setProperties(Array.isArray(data.properties) ? data.properties : data);
    } catch (e) {
      setProperties([]);
    }
  };

  useEffect(() => {
    fetchProperties();
  }, []);

  return (
    <div className="mapview-container">
      <div className="mapview-map">
        {/* 2️⃣ PropertyList에서 받아온 properties를 KakaoMap에도 그대로 전달 */}
        <KakaoMap properties={properties} />
      </div>
      <div className="mapview-list">
        <PropertyList
          ref={propertyListRef}
          user={user}
          mode="all"
        // PropertyList가 내부적으로 fetchProperties를 쓴다면, 따로 setProperties는 안 건드려도 됨
        // properties={properties} // (만약 PropertyList를 props 기반으로 리팩토링하면 추가)
        />
      </div>
    </div>
  );
};

export default MapView;
