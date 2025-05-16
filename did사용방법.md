
````markdown
# 🏠 Real Estate Blockchain (DID + VC 기반 공인중개사 인증 시스템)

이 프로젝트는 Hyperledger Fabric과 Go를 기반으로  
**공인중개사에게 DID를 발급하고, 자격증명서(VC)를 서명하여 제공하는 블록체인 기반 인증 시스템**입니다.  
사용자는 중개사의 VC를 검증하여, 허위 매물을 방지할 수 있습니다.

---

## 🧭 주요 기능

- ✅ 공인중개사 회원가입 시 **DID 자동 생성**
- ✅ 중개사 자격번호 기반 **VC(Verifiable Credential) 발급**
- ✅ 발급된 VC에 **Ed25519 개인키로 디지털 서명**
- ✅ 사용자(또는 앱)가 **VC JSON 검증** 가능
- ✅ `vc.json` 파일 생성 및 서명 검증 CLI 포함

---

## 📁 프로젝트 구조 요약

| 폴더 | 설명 |
|------|------|
| `cmd/` | CLI 테스트 도구 (`generate_keys`, `vc_test`, `verify_test`) |
| `vc/` | VC 발급 (`issuer.go`) + VC 검증 (`verify.go`) |
| `did/` | DID 생성 로직 |
| `crypto/` | 키 로딩, 서명, 서명 검증 유틸 |
| `keystore/` | 발급자 공개키/개인키 저장 |
| `handler/` | 공인중개사 HTTP 요청 처리 핸들러 |
| `models/` | 중개사 / 사용자 모델 정의 |
| `routes/` | 라우팅 설정 |
| `main.go` | Gin 서버 진입점 |

---

## ⚙️ 실행 방법

### 1. 레포지토리 클론

```bash
git clone https://github.com/04JinMyeong/real-estate-blockchain.git
cd real-estate-blockchain
git checkout feature/did-vc
cd go-backend
````

### 2. 의존성 설치

```bash
go mod tidy
```

### 3. 키 쌍 생성

```bash
go run cmd/generate_keys/main.go
```

생성 결과:

* `keystore/issuer_public.key`
* `keystore/issuer_private.key`

---

## 🔐 환경변수 설정

```bash
export PUBLIC_KEY_PATH=keystore/issuer_public.key
export PRIVATE_KEY_PATH=keystore/issuer_private.key
export ISSUER_DID=did:realestate:<여기에 본인이 생성한 DID 입력>
```

---

## 📤 VC 발급

```bash
go run cmd/vc_test/main.go
```

실행 결과:

* 콘솔에 VC JSON 출력
* `vc.json` 파일 자동 저장됨

---

## 🔍 VC 검증

```bash
go run cmd/verify_test/main.go vc.json
```

출력 예시:

```
✔ VC 검증 성공: 이 자격 증명서는 유효합니다.
```

또는

```
✘ VC 검증 실패: 위조 또는 변조된 VC입니다.
```

---

## 🧠 기술 개요

* 공개키 기반 **DID 생성**: `did:realestate:<hash>`
* `credentialSubject`에 중개사 정보 포함
* 서명 방식: **Ed25519 (RFC 8032)**
* `proof.jws`에 서명 추가
* JSON 기반 VC는 W3C 표준 스키마 준수

---

## 🙋 사용 예시 흐름

```plaintext
[공인중개사 회원가입]
  ↓
[서버: DID 생성 → DB 저장]
  ↓
[VC 생성 → Ed25519 서명 → JSON 반환]
  ↓
[사용자에게 QR / JSON 전달]
  ↓
[사용자: VC 검증 → 신뢰 판단]
```

---

## 🧪 주요 테스트 CLI 요약

| 명령              | 기능              |
| --------------- | --------------- |
| `generate_keys` | 키 쌍 생성          |
| `vc_test`       | VC 발급 및 JSON 저장 |
| `verify_test`   | VC JSON 검증      |

---

## 📝 참고

* [DID 표준 - W3C](https://www.w3.org/TR/did-core/)
* [Verifiable Credentials (VC) 표준](https://www.w3.org/TR/vc-data-model/)
* [Ed25519 서명](https://datatracker.ietf.org/doc/html/rfc8032)

```
