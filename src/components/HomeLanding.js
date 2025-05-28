// src/components/HomeLanding.js
import React from 'react';
import './HomeLanding.css';

const HomeLanding = () => {
  return (
    <section
      className="landing-container"
      style={{
        backgroundImage: "url('/background.jpg')",
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        height: '100vh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        position: 'relative',
      }}
    >
      <div className="landing-text">
        <h1>어떤 집을 찾고 계세요?</h1>
        <p>블록체인 기반 매물 정보를 지도에서 확인하세요.</p>
        {/* <button className="start-button" onClick={onStart}>지도</button> */}
        {/* ✅ 버튼 완전히 삭제하거나 주석 처리 */}
      </div>
    </section>
  );
};

export default HomeLanding;
