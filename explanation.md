##explanation

###1. App.js
-전체 앱을 구성하는 메인 컴포넌트
-KakaoMap, PropertyForm, PropertyList 컴포넌트를 조합
-PropertyForm에서 등록 후 PropertyList의 fetchBlocks 메서드 호출로 목록 갱신

###2. PropertyForm.js
-매물 등록 폼
-입력값을 axios.post('http://localhost:3000/addProperty')로 전송 → 백엔드에서 블록체인 처리 예상

###3. PropertyList.js
-등록된 매물들을 http://localhost:3000/getAllBlocks를 통해 블록 리스트 형태로 불러옴
-ref를 통해 외부에서 fetchBlocks 호출 가능

###4. KakaoMap.js
-Kakao 지도 연동
-index.html에서 Kakao SDK를 불러온 뒤, 지도 렌더링
