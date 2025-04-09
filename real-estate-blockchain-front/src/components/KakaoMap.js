import React, { useEffect } from 'react';

const KakaoMap = () => {
  useEffect(() => {
    const loadMap = () => {
      window.kakao.maps.load(() => {
        const container = document.getElementById('map');
        const options = {
          center: new window.kakao.maps.LatLng(37.5665, 126.9780), // 서울 시청 좌표
          level: 3,
        };
        const map = new window.kakao.maps.Map(container, options);
        
        // 마커 생성
        const marker = new window.kakao.maps.Marker({
          position: map.getCenter()
        });
        marker.setMap(map);
      });
    };

    // 카카오맵 SDK가 이미 로드되어 있는지 확인
    if (window.kakao && window.kakao.maps) {
      loadMap();
    } else {
      const script = document.createElement('script');
      script.src = 'https://dapi.kakao.com/v2/maps/sdk.js?appkey=4f98b430d3ac0e982a3c3bd31b8b410d&autoload=false';
      script.async = true;
      script.onload = loadMap;
      document.head.appendChild(script);
    }

    // cleanup function
    return () => {
      const mapContainer = document.getElementById('map');
      if (mapContainer) {
        mapContainer.innerHTML = '';
      }
    };
  }, []);

  return (
    <div
      id="map"
      style={{
        width: '100%',
        height: '100%',
      }}
    ></div>
  );
};

export default KakaoMap;
