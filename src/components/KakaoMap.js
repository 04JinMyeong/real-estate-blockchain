import React, { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom'; // 추가!

const API_URL            = 'https://2094-165-229-229-106.ngrok-free.app';
const KAKAO_JS_KEY       = '4f98b430d3ac0e982a3c3bd31b8b410d';
const KAKAO_REST_API_KEY = 'c32790578348e6ef920eca53f2c23382';
const DEFAULT_CENTER     = { lat: 37.5665, lng: 126.9780 };
const handleReserve = async (property) => {
    // 실제 로그인 유저 정보 필요하면 props로 받아서 user.username 사용 (아래처럼 property.createdBy도 임시로 가능)
    const expiresAt = Math.floor(Date.now() / 1000) + 12 * 3600;
    const payload = {
      user: property.createdBy, // 혹은 로그인한 유저 정보로 교체!
      id: property.id,
      expiresAt
    };
    const token = localStorage.getItem('token');
    try {
      const res = await fetch(
        `${API_URL}/reserve-property`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            'ngrok-skip-browser-warning': 'true'
          },
          body: JSON.stringify(payload)
        }
      );
      const data = await res.json();
      alert(data.message || '✅ 예약이 완료되었습니다.');
      // 예약 성공 후 목록 갱신 원하면 아래 둘 중 하나 사용:
      // 방법1: 전체 새로고침 (간단)
      window.location.reload();
      // 방법2: setProperties로 state만 갱신 (원하면 따로 안내 가능)
    } catch (err) {
      alert('예약 실패: ' + err.message);
    }
  };

export default function KakaoMap() {
  const [properties, setProperties] = useState([]);
  const mapRef = useRef(null);
  const navigate = useNavigate(); // 추가

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
        setProperties(list);
      })
      .catch(() => setProperties([]));
  }, []);

  useEffect(() => {
    const container = mapRef.current;
    if (!container) return;

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

    let openedInfoWindow = null;

    const initMap = async () => {
      const map = new window.kakao.maps.Map(container, {
        center: new window.kakao.maps.LatLng(DEFAULT_CENTER.lat, DEFAULT_CENTER.lng),
        level: 5,
      });

      if (!properties.length) return;

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
      if (!valid.length) return;

      map.setCenter(new window.kakao.maps.LatLng(valid[0].lat, valid[0].lng));

      valid.forEach(p => {
        const pos = new window.kakao.maps.LatLng(p.lat, p.lng);
        const marker = new window.kakao.maps.Marker({ map, position: pos });

        const owner = p.ownerHistory?.at(-1)?.owner ?? '-';
        const price = p.priceHistory?.at(-1)?.price?.toLocaleString() ?? '-';

        // 항상 상세보기 버튼이 보이도록
        const reservedOrButton = p.reservedBy
          ? `<span style="color:red; margin-right:10px;">예약됨</span>`
          : `<button class="map-reserve-btn" data-id="${p.id}">예약하기</button>`;

        // 항상 상세보기 버튼 포함
        const content = `
          <div style="
            width: 240px;
            background: #fff;
            border-radius: 15px;
            box-shadow: 0 2px 12px rgba(40,40,40,0.14);
            padding: 18px 20px 12px 20px;
            border: 1px solid #ececec;
            font-family: 'Pretendard','Noto Sans KR',sans-serif;
          ">
            <div style="font-size:18px; font-weight:bold; margin-bottom:8px; color:#212121;">
              ${p.address}
            </div>
            <div style="font-size:15px; margin-bottom:5px;">
              👤 <span style="font-weight:bold;">${owner}</span>
            </div>
            <div style="font-size:15px; margin-bottom:5px;">
              💰 <span style="color:#E8B100; font-weight:bold;">${price} 원</span>
            </div>
            <div style="font-size:15px; margin-top:7px; display:flex; gap:8px; align-items:center;">
              <button class="map-detail-btn" data-id="${p.id}">상세보기</button>
              ${reservedOrButton}
            </div>
          </div>
        `;

        const iw = new window.kakao.maps.InfoWindow({ content });

        window.kakao.maps.event.addListener(marker, 'click', () => {
          if (openedInfoWindow === iw) {
            iw.close();
            openedInfoWindow = null;
          } else {
            if (openedInfoWindow) openedInfoWindow.close();
            iw.open(map, marker);
            openedInfoWindow = iw;

            // InfoWindow가 DOM에 그려진 뒤 버튼 이벤트 바인딩
            setTimeout(() => {
              // 상세보기 버튼
              const detailBtn = document.querySelector(`.map-detail-btn[data-id="${p.id}"]`);
              if (detailBtn) {
                detailBtn.onclick = (e) => {
                  e.preventDefault();
                  navigate(`/properties/${p.id}`);
                };
              }

              // 예약 버튼
              if (!p.reservedBy) {
                const reserveBtn = document.querySelector(`.map-reserve-btn[data-id="${p.id}"]`);
                if (reserveBtn) {
      reserveBtn.onclick = (e) => {
        e.preventDefault();
        handleReserve(p);
                  };
                }
              }
            }, 100);
          }
        });
      });
    };

    if (typeof window.kakao === 'undefined' || !window.kakao.maps) {
      const script = document.createElement('script');
      script.src = `https://dapi.kakao.com/v2/maps/sdk.js?appkey=${KAKAO_JS_KEY}&libraries=services&autoload=false`;
      script.async = true;
      script.onload = () => {
        window.kakao.maps.load(initMap);
      };
      document.head.appendChild(script);
    } else {
      window.kakao.maps.load(initMap);
    }

    return () => {
      if (container) container.innerHTML = '';
    };
  }, [properties, navigate]);

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
