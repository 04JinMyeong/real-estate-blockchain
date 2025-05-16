import React, { useState } from 'react';
import axios from 'axios';

function PropertyForm({ user, onRegister }) {
  const [form, setForm] = useState({
    address: '',
    owner: '',
    price: '',
  });

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();
    try {
      await axios.post('https://1af7-165-229-229-137.ngrok-free.app/add-property', {
        user: user.email,
        address: form.address,
        owner: form.owner,
        price: form.price,
      });
      alert('📌 등록 완료!');
      onRegister();
      setForm({ address: '', owner: '', price: '' });
    } catch (err) {
      alert('⚠️ 등록 실패');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h3>🏠 매물 등록</h3>
      <input name="address" placeholder="주소" value={form.address} onChange={handleChange} required />
      <input name="owner" placeholder="소유자" value={form.owner} onChange={handleChange} required />
      <input name="price" type="number" placeholder="가격" value={form.price} onChange={handleChange} required />
      <button type="submit">등록</button>
    </form>
  );
}

export default PropertyForm;
