// src/components/MapView.js
import React from 'react';
import KakaoMap from './KakaoMap';
import PropertyList from './PropertyList';
import './MapView.css';

const MapView = () => {
  return (
    <div className="mapview-container">
      <div className="mapview-map">
        <KakaoMap />
      </div>
      <div className="mapview-list">
        <PropertyList />
      </div>
    </div>
  );
};

export default MapView;
