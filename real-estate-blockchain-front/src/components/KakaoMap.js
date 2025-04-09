import React, { useEffect } from 'react';

const KakaoMap = () => {
  useEffect(() => {
    const script = document.createElement('script');
    script.src =
      'https://dapi.kakao.com/v2/maps/sdk.js?appkey=4f98b430d3ac0e982a3c3bd31b8b410d&autoload=false';
    script.async = true;

    script.onload = () => {
      window.kakao.maps.load(() => {
        const container = document.getElementById('map');
        const options = {
          center: new window.kakao.maps.LatLng(33.450701, 126.570667),
          level: 3,
        };
        new window.kakao.maps.Map(container, options);
      });
    };

    document.head.appendChild(script);
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
