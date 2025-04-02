import React, {
    useEffect,
    useState,
    forwardRef,
    useImperativeHandle,
  } from 'react';
  import axios from 'axios';
  
  const PropertyList = forwardRef((props, ref) => {
    const [blocks, setBlocks] = useState([]);
  
    const fetchBlocks = async () => {
      try {
        const res = await axios.get('http://localhost:8080/getAllBlocks');
        const parsed = typeof res.data === 'string' ? JSON.parse(res.data) : res.data;
        setBlocks(parsed);
      } catch (err) {
        console.error('불러오기 실패:', err);
      }
    };
  
    useImperativeHandle(ref, () => ({
      fetchBlocks,
    }));
  
    useEffect(() => {
      fetchBlocks();
    }, []);
  
    return (
      <div>
        <h3>📋 등록된 매물</h3>
        <ul>
          {blocks.map((block) => (
            <li key={block.index}>
              {block.property.id} | {block.property.location} | {block.property.price} | {block.property.owner}
            </li>
          ))}
        </ul>
      </div>
    );
  });
  
  export default PropertyList;
  