import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
const API_URL = 'https://1af7-165-229-229-137.ngrok-free.app';

const UserMypage = ({ user }) => {
  const [myReservations, setMyReservations] = useState([]);

  const fetchReservations = useCallback(async () => {
    if (!user?.id) return;
    const token = localStorage.getItem('token');
    try {
      const res = await axios.get(
        `${API_URL}/api/reservations?userId=${user.id}`,
        { headers: { 'Authorization': `Bearer ${token}` } }
      );
      setMyReservations(res.data);
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
      <p><strong>아이디:</strong> {user?.id}</p>
      <div className="mypage-card">
        <h2>📌 내 예약 목록</h2>
        {myReservations.length === 0 ? (
          <p>예약한 매물이 없습니다.</p>
        ) : (
          <ul>
            {myReservations.map(resv => (
              <li key={resv.reservationId}>
                [매물ID: {resv.propertyId}] {resv.date} {resv.time} - {resv.status}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};
export default UserMypage;
