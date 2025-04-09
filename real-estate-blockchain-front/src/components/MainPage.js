import React from 'react';
import './MainPage.css';

const MainPage = () => {
  return (
    <div className="main-container">
      <div className="hero-section">
        <div className="hero-content">
          <h1>모든 부동산 이력을 블록체인에 담다</h1>
          <h2>거짓 없는 매물,검증된 거래의 시작 - TrueHome</h2>
          <p>블록체인 기술을 기반으로 모든 매물의 등록 정보, 가격 변동, 거래 이력을 투명하고 안전하게 관리합니다.</p>
          <p>한 번 블록체인에 기록된 정보는 누구도 임의로 수정하거나 조작할 수 없으며,</p>
          <p>언제 어디서나 그 기록을 열람하고 검증할 수 있습니다.</p>
          <button className="learn-more">LEARN MORE</button>
        </div>
        <div className="hero-image"></div>
      </div>
    </div>
  );
};

export default MainPage; 