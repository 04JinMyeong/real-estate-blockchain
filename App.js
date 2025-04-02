import React, { useRef } from 'react';
import './App.css';
import logo from './logo.svg';
import PropertyForm from './PropertyForm';
import PropertyList from './PropertyList';

function App() {
  const listRef = useRef();

  const handleRegister = () => {
    if (listRef.current && listRef.current.fetchBlocks) {
      listRef.current.fetchBlocks();
    }
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial' }}>
      <h1>🏗️ C++ 블록체인 부동산 시스템 (React)</h1>
      <PropertyForm onRegister={handleRegister} />
      <hr />
      <PropertyList ref={listRef} />
    </div>
  );
}

export default App;
