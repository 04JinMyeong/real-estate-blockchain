import React, { useState } from 'react';
import axios from 'axios';
import './VcIssueModal.css';

const MOCK_ISSUER_API_ENDPOINT = 'http://localhost:8083/issue-vc';

function VcIssueModal({ onClose }) {
    const [formData, setFormData] = useState({ name: '', id: '', did: '' });
    const [isLoading, setIsLoading] = useState(false);
    const [issuedVC, setIssuedVC] = useState(null);
    const [error, setError] = useState('');

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');
        setIssuedVC(null);

        try {
            const response = await axios.post(MOCK_ISSUER_API_ENDPOINT, formData);
            setIssuedVC(response.data.vc);
        } catch (err) {
            setError(`VC 발급 실패: ${err.response?.data?.error || err.message}`);
        } finally {
            setIsLoading(false);
        }
    };

    // --- ▼▼▼ 1. 클립보드 복사 기능 구현 ▼▼▼ ---
    const copyToClipboard = () => {
        if (!issuedVC) return;
        const vcString = JSON.stringify(issuedVC, null, 2);
        navigator.clipboard.writeText(vcString).then(() => {
            alert('✅ VC가 클립보드에 복사되었습니다!');
        }).catch(err => {
            console.error('클립보드 복사 실패:', err);
            alert('❌ 클립보드 복사에 실패했습니다.');
        });
    };

    // --- ▼▼▼ 2. 파일 다운로드 기능 구현 ▼▼▼ ---
    const downloadVCAsFile = () => {
        if (!issuedVC) return;
        const vcString = JSON.stringify(issuedVC, null, 2);
        const blob = new Blob([vcString], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = 'my-real-estate-vc.json'; // 다운로드될 파일 이름
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    };

    return (
        <div className="vc-modal-overlay">
            <div className="vc-modal-content">
                <button className="vc-modal-close" onClick={onClose}>&times;</button>
                <h2>📜 자격증명(VC) 발급</h2>

                {!issuedVC ? (
                    <form onSubmit={handleSubmit}>
                        <p>VC를 발급받기 위해 정보를 입력해주세요.</p>
                        <input name="name" type="text" value={formData.name} onChange={handleChange} placeholder="이름" required />
                        <input name="id" type="text" value={formData.id} onChange={handleChange} placeholder="아이디" required />
                        <input name="did" type="text" value={formData.did} onChange={handleChange} placeholder="발급받은 DID" required />
                        <button type="submit" disabled={isLoading}>{isLoading ? '발급 중...' : 'VC 발급'}</button>
                        {error && <p style={{ color: 'red' }}>{error}</p>}
                    </form>
                ) : (
                    <div className="issued-info-display">
                        <h3>[중요] 발급된 VC 정보</h3>
                        <p>아래 VC 정보를 복사하거나 파일로 다운로드하여 안전하게 보관하세요. **로그인 시 필요합니다.**</p>
                        <textarea value={JSON.stringify(issuedVC, null, 2)} readOnly rows="10" />

                        {/* --- ▼▼▼ 버튼 UI 수정 ▼▼▼ --- */}
                        <div className="button-group">
                            <button onClick={copyToClipboard}>VC 텍스트 복사</button>
                            <button onClick={downloadVCAsFile}>VC 파일로 다운로드</button>
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}

export default VcIssueModal;