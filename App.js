import logo from './logo.svg';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;
import React, { useRef } from 'react';
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
