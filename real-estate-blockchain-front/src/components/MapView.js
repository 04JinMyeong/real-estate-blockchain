// src/components/MapView.js
import React from 'react';
import KakaoMap from './KakaoMap';
import PropertyList from '../PropertyList';
import PropertyForm from '../PropertyForm';

const MapView = () => {
  return (
    <div style={{ display: 'flex', height: 'calc(100vh - 60px)' }}>
      <div style={{ flex: 8 }}>
        <KakaoMap />
      </div>
      <div
        style={{
          flex: 2,
          overflowY: 'auto',
          borderLeft: '1px solid #ddd',
          backgroundColor: '#fff',
        }}
      >
        <PropertyList />
      </div>
    </div>
  );
};

export default MapView;
