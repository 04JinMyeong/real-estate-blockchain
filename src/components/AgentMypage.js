// src/components/AgentMypage.js  (혹은 MyPage.js)
import React, { useState, useRef, useEffect } from 'react';
import PropertyForm from './PropertyForm';
import PropertyList from './PropertyList';
import './MyPage.css';

export default function AgentMypage({ user }) {
  const [properties, setProperties] = useState([]);
  const didFetchRef = useRef(false);

  // ── 1) 최초 마운트 시 & user가 바뀔 때, 내 매물 목록을 가져오기 ─────────────────
  useEffect(() => {
    async function loadMyProperties() {
      if (!user) {
        setProperties([]);
        return;
      }
      try {
        const res = await fetch(
          // 반드시 현재 동작 중인 ngrok 주소(또는 배포된 API 주소)로 바꿔주세요.
          //`https://2094-165-229-229-106.ngrok-free.app/my-properties?user=${user.username}`
          `http://localhost:8080/my-properties?user=${user.username}` // 로컬 개발용
        );
        const data = await res.json();
        setProperties(
          Array.isArray(data.properties)
            ? data.properties
            : Array.isArray(data)
              ? data
              : []
        );
      } catch {
        setProperties([]);
      }
    }
    if (!didFetchRef.current) {
      loadMyProperties();
      didFetchRef.current = true;
    }
  }, [user]);

  // ── 2) PropertyForm이 등록에 성공한 뒤 호출되는 함수 ────────────────────────
  const refreshMyList = async () => {
    if (!user) return;
    try {
      const res = await fetch(
        // 반드시 현재 동작 중인 ngrok 주소(또는 배포된 API 주소)로 바꿔주세요.
        //`https://2094-165-229-229-106.ngrok-free.app/my-properties?user=${user.username}`
        `http://localhost:8080/my-properties?user=${user.username}` // 로컬 개발용
      );
      const data = await res.json();
      setProperties(
        Array.isArray(data.properties)
          ? data.properties
          : Array.isArray(data)
            ? data
            : []
      );
    } catch {
      setProperties([]);
    }
  };

  return (
    <div className="mypage-container">
      <h2>🏢 중개인 마이페이지</h2>
      <p><strong>아이디:</strong> {user?.username || user?.id}</p>

      <div className="mypage-card">
        <h2>🏠 매물 등록</h2>
        <PropertyForm user={user} onRegister={refreshMyList} />
      </div>

      <div className="mypage-card" style={{ marginTop: '2rem' }}>
        <h2>📋 내 매물 목록</h2>
        <PropertyList
          user={user}
          mode="my"
          onReserve={refreshMyList}
        // 위 컴포넌트처럼 ref를 사용해 직접 fetchProperties를 호출해도 되지만,
        // 여기서는 onReserve 콜백을 통해 목록을 갱신하도록 했습니다.
        />
      </div>
    </div>
  );
}
