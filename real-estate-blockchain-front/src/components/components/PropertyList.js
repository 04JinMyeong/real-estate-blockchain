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
      const res = await axios.get('https://bb52-219-251-84-31.ngrok-free.app/properties?user=TestUser9', {
        headers: { 'ngrok-skip-browser-warning': 'true' }
      });
      console.log('📦 응답:', res.data);
      setBlocks(res.data);
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
    <div className="property-list">
      <h3>📋 등록된 매물</h3>
      {blocks.length === 0 ? (
        <p>등록된 매물이 없습니다.</p>
      ) : (
        <div className="property-grid">
          {blocks.map((block, index) => {
            const latestPriceObj = block.priceHistory?.[block.priceHistory.length - 1];
            const latestPrice = latestPriceObj ? Number(latestPriceObj.price) : NaN;
            const displayPrice = isNaN(latestPrice) ? '정보 없음' : latestPrice.toLocaleString();

            const latestOwnerObj = block.ownerHistory?.[block.ownerHistory.length - 1];
            const displayOwner = latestOwnerObj?.owner || '정보 없음';

            return (
              <div className="property-card" key={index}>
                <h4>{block.address}</h4>
                <p>💰 가격: {displayPrice} 원</p>
                <p>👤 소유자: {displayOwner}</p>
                <p>🆔 매물 ID: {block.id}</p>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
});

export default PropertyList;
