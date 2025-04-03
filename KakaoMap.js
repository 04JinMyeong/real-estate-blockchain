import React, { useEffect } from 'react';

const KakaoMap = () => {
  useEffect(() => {
    const loadKakaoMap = () => {
      // SDK 완전히 로드되었을 때만 실행
      if (window.kakao && window.kakao.maps && window.kakao.maps.load) {
        window.kakao.maps.load(() => {
          const container = document.getElementById('map');
          const options = {
            center: new window.kakao.maps.LatLng(37.5665, 126.9780),
            level: 3,
          };
          new window.kakao.maps.Map(container, options);
        });
      } else {
        console.error("❌ Kakao SDK 미로딩 상태. 100ms 후 재시도");
        setTimeout(loadKakaoMap, 100); // 재시도
      }
    };

    loadKakaoMap();
  }, []);

  return (
    <div
      id="map"
      style={{ width: '100%', height: '500px', marginTop: '20px', backgroundColor: '#eee' }}
    ></div>
  );
};

export default KakaoMap;
