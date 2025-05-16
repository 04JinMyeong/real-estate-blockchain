// src/components/Wave.js
import React from 'react';

const Wave = () => (
 <svg
  xmlns="http://www.w3.org/2000/svg"
  viewBox="0 0 1440 320"
  style={{
    position: 'absolute',
    bottom: 0,
    left: 0,
    width: '100%',
    height: '200px',
    zIndex: 0, // Wave를 Hero의 가장 뒤로 (텍스트가 Wave 위로 올라오지 않음)
    pointerEvents: 'none',
  }}
>

    <path
      fill="#0099ff"
      fillOpacity="1"
      d="M0,256L60,213.3C120,171,240,85,360,58.7C480,32,600,64,720,112C840,160,960,224,1080,218.7C1200,213,1320,139,1380,101.3L1440,64L1440,0L1380,0C1320,0,1200,0,1080,0C960,0,840,0,720,0C600,0,480,0,360,0C240,0,120,0,60,0L0,0Z"
    />
  </svg>
);

export default Wave;
