// src/components/VcIssueModal.js

import React, { useState } from "react";
import "./VcIssueModal.css";

function VcIssueModal({ onClose }) {
    // 1) 입력 폼에 필요한 state 선언
    const [name, setName] = useState("");
    const [phoneNumber, setPhoneNumber] = useState("");
    const [licenseNumber, setLicenseNumber] = useState("");
    const [gender, setGender] = useState("");
    const [did, setDid] = useState("");

    // 2) 폼 제출 핸들러
    const handleSubmit = (e) => {
        e.preventDefault();

        // 빈 칸 체크
        if (!name || !phoneNumber || !licenseNumber || !gender || !did) {
            alert("모든 정보를 입력해주세요.");
            return;
        }

        // 실제 발급 로직(예: axios/fetch)을 여기서 처리할 수도 있음
        // 예시: axios.post("/api/vc/issue", { name, email, did })
        //       .then(response => { ... }).catch(err => { ... });

        // 지금은 단순 알림
        alert("발급이 완료되었습니다!");
        onClose(); // 모달 닫기
    };

    return (
        <div className="vc-modal-overlay">
            <div className="vc-modal-content">
                {/* 닫기 버튼(×) */}
                <button className="vc-modal-close" onClick={onClose}>
                    &times;
                </button>

                <h2>VC 발급 정보 입력</h2>
                <form className="vc-issue-form" onSubmit={handleSubmit}>
                    {/* (3) 이름(Name) */}
                    <label htmlFor="vc-name">ID</label>
                    <input
                        id="vc-name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="홍길동"
                    />

                    {/* (4) 전화번호(Phone Number) */}
                    <label htmlFor="vc-phone">PassWord</label>
                    <input
                        id="vc-phone"
                        type="text"
                        value={phoneNumber}
                        onChange={(e) => setPhoneNumber(e.target.value)}
                        placeholder="010-1234-5678"
                    />

                    {/* (5) 자격증번호(License Number) */}
                    <label htmlFor="vc-license">자격증번호</label>
                    <input
                        id="vc-license"
                        type="text"
                        value={licenseNumber}
                        onChange={(e) => setLicenseNumber(e.target.value)}
                        placeholder="ABCD-1234"
                    />

                    {/* (6) 성별(Gender) */}
                    <label htmlFor="vc-gender">성별</label>
                    <input
                        id="vc-gender"
                        type="text"
                        value={gender}
                        onChange={(e) => setGender(e.target.value)}
                        placeholder="남자 / 여자 / 기타"
                    />

                    {/* (7) DID (기존) */}
                    <label htmlFor="vc-did">DID</label>
                    <input
                        id="vc-did"
                        type="text"
                        value={did}
                        onChange={(e) => setDid(e.target.value)}
                        placeholder="did:example:123456789"
                    />

                    {/* (8) 제출 버튼 */}
                    <button type="submit" className="vc-submit-button">
                        지금 VC발급받기
                    </button>
                </form>
            </div>
        </div>
    );
}

export default VcIssueModal;
