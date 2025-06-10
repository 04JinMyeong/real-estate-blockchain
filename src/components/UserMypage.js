import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import PropertyList from './PropertyList';

//const API_URL = 'https://2094-165-229-229-106.ngrok-free.app';
const API_URL = 'http://localhost:8080'; // 로컬 개발용

const UserMypage = ({ user }) => {
  const [myReservations, setMyReservations] = useState([]);

  // user?.username으로 구조 통일!
  const fetchReservations = useCallback(async () => {
    if (!user?.username) return;
    const token = localStorage.getItem('token');
    try {
      const res = await axios.get(
        `${API_URL}/api/reservations?userId=${user.username}`,
        { headers: { 'Authorization': `Bearer ${token}` } }
      );
      setMyReservations(Array.isArray(res.data) ? res.data : []);
    } catch (err) {
      console.error('내 예약 목록 불러오기 실패:', err);
    }
  }, [user]);

  useEffect(() => {
    fetchReservations();
  }, [fetchReservations]);

  return (
    <div className="mypage-container">
      <h2>마이페이지</h2>
      <p><strong>아이디:</strong> {user?.username || user?.email || '-'}</p>
      <div className="mypage-card">
        <h2>📌 내 예약 목록</h2>
        {myReservations.length === 0 ? (
          <p>예약한 매물이 없습니다.</p>
        ) : (
          <ul>
            {(Array.isArray(myReservations) ? myReservations : []).map(resv => (
              <li key={resv.reservationId}>
                [매물ID: {resv.propertyId}] {resv.date} {resv.time} - {resv.status}
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* 내가 등록한 매물 리스트도 보여주기! */}
      <div className="mypage-card">
        <h2>🏠 내가 등록한 매물 목록</h2>
        <PropertyList user={user} mode="my" />
      </div>
    </div>
  );
};

export default UserMypage;
