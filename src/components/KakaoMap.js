// src/components/KakaoMap.js
import React, { useEffect, useRef, useState } from 'react';

// 백엔드 API 주소
const API_URL            = 'https://2094-165-229-229-106.ngrok-free.app';
// 카카오 API 키 (본인 키로 바꿔주세요)
const KAKAO_JS_KEY       = '4f98b430d3ac0e982a3c3bd31b8b410d';
const KAKAO_REST_API_KEY = 'c32790578348e6ef920eca53f2c23382';
// 지도 기본 중심 (서울 시청)
const DEFAULT_CENTER     = { lat: 37.5665, lng: 126.9780 };

export default function KakaoMap() {
  const [properties, setProperties] = useState([]);
  const mapRef = useRef(null);

  // 1) 백엔드에서 매물 목록 가져오기 (한 번만)
  useEffect(() => {
    fetch(`${API_URL}/properties?user=admin`, {
      headers: { 'ngrok-skip-browser-warning': 'true' }
    })
      .then(r => r.json())
      .then(data => {
        const list = Array.isArray(data.properties)
          ? data.properties
          : Array.isArray(data)
            ? data
            : [];
        console.log('🔍 fetched properties:', list);
        setProperties(list);
      })
      .catch(err => {
        console.error('❌ fetch properties failed:', err);
        setProperties([]);
      });
  }, []);

  // 2) properties 변경 시마다 map 초기화 & 마커 그리기
  useEffect(() => {
    const container = mapRef.current;
    if (!container) return;

    // 주소→위경도 변환 함수
    const geocode = async (addr) => {
      const res = await fetch(
        `https://dapi.kakao.com/v2/local/search/address.json?query=${encodeURIComponent(addr)}`,
        { headers: { Authorization: `KakaoAK ${KAKAO_REST_API_KEY}` } }
      );
      const json = await res.json();
      if (!json.documents?.length) throw new Error('no docs');
      const { x, y } = json.documents[0];
      return { lat: parseFloat(y), lng: parseFloat(x) };
    };

    // 실제 지도+마커 그리는 initMap
    const initMap = async () => {
      console.log('✅ kakao.maps is ready, initializing map');
      const map = new window.kakao.maps.Map(container, {
        center: new window.kakao.maps.LatLng(DEFAULT_CENTER.lat, DEFAULT_CENTER.lng),
        level: 5,
      });

      if (!properties.length) {
        console.warn('⚠️ no properties — showing default map');
        return;
      }

      const coordsList = await Promise.all(
        properties.map(async p => {
          try {
            const c = await geocode(p.address);
            return { ...p, ...c };
          } catch {
            return null;
          }
        })
      );
      const valid = coordsList.filter(Boolean);
      if (!valid.length) {
        console.warn('⚠️ all geocode failed');
        return;
      }

      // 첫 매물 중심으로 이동
      map.setCenter(new window.kakao.maps.LatLng(valid[0].lat, valid[0].lng));

      valid.forEach(p => {
        const pos = new window.kakao.maps.LatLng(p.lat, p.lng);
        const marker = new window.kakao.maps.Marker({ map, position: pos });

        const owner   = p.ownerHistory?.at(-1)?.owner ?? '-';
        const price   = p.priceHistory?.at(-1)?.price?.toLocaleString() ?? '-';
        const reserved= p.reservedBy
          ? `<span style="color:red">예약됨</span>`
          : `<button style="margin-top:8px;">예약하기</button>`;

        const content = `
          <div class="property-card" style="width:200px;">
            <h4>${p.address}</h4>
            <p>👤 소유자: ${owner}</p>
            <p>💰 가격: ${price} 원</p>
            <p>${reserved}</p>
          </div>
        `;
        const iw = new window.kakao.maps.InfoWindow({ content });
        window.kakao.maps.event.addListener(marker, 'click', () => iw.open(map, marker));
      });
    };

    // SDK 로드 (autoload=false 로, maps.load(initMap) 반드시 사용)
    if (typeof window.kakao === 'undefined' || !window.kakao.maps) {
      const script = document.createElement('script');
      script.src = `https://dapi.kakao.com/v2/maps/sdk.js?appkey=${KAKAO_JS_KEY}&libraries=services&autoload=false`;
      script.async = true;
      script.onload = () => {
        console.log('📥 Kakao SDK loaded, calling kakao.maps.load');
        window.kakao.maps.load(initMap);
      };
      document.head.appendChild(script);
    } else {
      window.kakao.maps.load(initMap);
    }

    return () => {
      if (container) container.innerHTML = '';
    };
  }, [properties]);

  return (
    <div style={{ width: '100%', height: '100%' }}>
      <div
        ref={mapRef}
        style={{
          width: '100%',
          height: '100%',
          borderRadius: '12px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.06)'
        }}
      />
    </div>
  );
}
