import React, { useState } from 'react';
import axios from 'axios';

function PropertyForm({ onRegister }) {
  const [form, setForm] = useState({
    id: '',
    location: '',
    price: '',
    owner: '',
  });

  const handleChange = e => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async e => {
    e.preventDefault();
    try {
      await axios.post('http://localhost:3000/addProperty', {
        ...form,
        price: parseFloat(form.price)
      });
      alert('📌 등록 완료!');
      onRegister();
      setForm({ id: '', location: '', price: '', owner: '' });
    } catch (err) {
      alert('⚠️ 등록 실패');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h3>🏠 매물 등록</h3>
      <input name="id" placeholder="매물 ID" value={form.id} onChange={handleChange} required />
      <input name="location" placeholder="위치" value={form.location} onChange={handleChange} required />
      <input name="price" type="number" placeholder="가격" value={form.price} onChange={handleChange} required />
      <input name="owner" placeholder="소유자" value={form.owner} onChange={handleChange} required />
      <button type="submit">등록</button>
    </form>
  );
}

export default PropertyForm;
