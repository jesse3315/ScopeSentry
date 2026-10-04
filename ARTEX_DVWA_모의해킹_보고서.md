# ARTEX 자율 침투 테스트 콘솔 — DVWA 모의해킹 수행 보고서

> 본 보고서는 ARTEX 플랫폼의 **대화(Auto Agent) 흐름을 중심으로**, 유료 모델
> **OpenRouter GLM 5.3 Flash**로 수행한 테스트용 DVWA 사이트(`http://100.11.*.*`)
> 모의해킹의 전 과정을 정리한 것이다. 1차 목적은 ARTEX의 **작업 생성 → 워커 실행 →
> 라이브 재검증 → 발견 등록(report_finding) → 상세보고서**에 이르는 오케스트레이션
> 로직이 **엔드투엔드로 정상 동작**하는지 검증하는 것이다.

---

## 1. 개요

| 항목 | 내용 |
|------|------|
| 테스트 일시 | 2026-10-04 (KST) |
| 대상 사이트 | 테스트용 DVWA v1.10 *Development* — `http://100.11.*.*` |
| 대상 스택 | Apache/2.4.25 (Debian) · PHP · MySQL · 보안레벨 **low** |
| 목적 | ARTEX 구동 로직·기능의 엔드투엔드 검증 + DVWA 취약점 점검 |
| 작업 ID | Task #1(지휘) → Task #2(실행) → **Task #3(등록 마감, 완료)** |
| 수행 주체 | ARTEX 운영 도우미 **Auto**(대화 Agent) + 작업 워커(planner/worker) |
| **사용 모델** | **OpenRouter GLM 5.3 Flash (유료)** |
| 확정 취약점 | **15건 등록 완료** — Critical 3 / High 2 / Medium 8 / Low 2 |
| 최종 결과 | **전 과정 정상 완료** — 라이브 재검증 → 15건 report_finding 등록(confirmed) → 상세보고서 5건 등록 |

> **무료 티어는 초기 테스트에서 잠깐 사용**했을 뿐이다. 최초 탐색 세션에서 무료 모델의
> 분당 한도(429, rate limit)로 작업엔진이 일시 중단된 적이 있으나, 이는 테스트 환경의
> 일시적 현상이었다. **본 수행은 유료 GLM 5.3 Flash로 전환하여 진행**했고, 레이트리밋
> 중단 없이 작업 생성부터 발견 등록·상세보고서까지 **전 파이프라인을 완주**했다.

본 보고서는 유료 GLM 5.3 Flash로 수행한 **완주 세션**을 기준으로 한다. 대상 DVWA는 보안
교육·검증 전용 환경이며, 업로드된 테스트 웹셸은 테스트 특성상 보존했다(삭제 가능).

---

## 2. ARTEX 콘솔/플랫폼 구성

테스트가 수행된 ARTEX 콘솔(`버전 · dev`)은 **대화·작업·발견·트래픽·도구 실행·LLM 녹화·
자산·차단 규칙·자산 차단** 등 침투 테스트 운용 기능을 단일 콘솔로 제공한다. 본 세션은
대화(Auto)에서 작업을 오케스트레이션하고, 작업 워커가 자율 실행하며, 발견이 플랫폼에
정식 등록되는 전 흐름을 사용했다.

---

## 3. 수행 과정 — ARTEX 대화 중심

사용자는 대화(Auto) 세션에 다음을 지시했다.

> 테스트용 DVWA사이트를 모의해킹하라. ARTEX 돌아가는 로직과 기능을 테스트하기 위함이다.
> 작업은 녹화하고 사용하는 도구 등 기록하라. 목표 `http://100.11.*.*`

### 3.1 대화 장면 재구성 — Auto의 독백 🎬

그림 1~3의 대화에는 Auto 에이전트의 사고(reasoning)가 영문으로 기록된다. 아래는 각
메시지를 장면처럼 드라마틱하게 풀어 쓴 해석으로, 유료 GLM 5.3 Flash가 작업을 어떻게
이어받아 끝까지 밀어붙였는지를 보여준다.

- **지시를 받다** — *"The user wants to test a DVWA site … for testing the ARTEX platform's logic."*
  → 명령이 떨어졌다. **"이 DVWA를 뚫어, ARTEX가 어디까지 해내는지 보여라."** 이번엔 중간에
  멈추지 않는다 — 유료 엔진을 달았다.
- **기존 진척을 발견하다** — *"There are existing paused tasks."*
  → 백지가 아니다. 멈춰 선 작업 두 개(①ARTEX 테스트 ②레이트리밋 우회 서브태스크)가 남아 있다.
  **"처음부터 다시 하지 않는다. 멈춘 지점에서 이어 붙인다."**
- **증거를 되짚다** — *"Evidence files exist in conv-1. … Evidence confirmed for the critical/high findings."*
  → 이전 세션이 남긴 **증거 47개**를 펼친다. 핵심 4건(SQLi·RCE·웹셸·LFI)이 멀쩡히 살아 있다.
  **"기록은 거짓말하지 않는다. 그러나 나는 직접 눈으로 확인한다."**
- **라이브로 다시 찌른다** — *"Login successful, security=low confirmed. … All 5 re-verified live."*
  → `admin/password` 로그인 성공, 보안레벨 low 확인. 그리고 **지금 이 순간** 다섯 급소를 다시
  찌른다 — 전부 재현. **"과거의 증거가 아니라, 살아 있는 증거다."**
- **전리품을 등록하다** — *"Asset 3 confirmed … Now register the 15 findings."*
  → 자산 id=3 확인. **"이제 15개의 취약점을 플랫폼에 정식으로 새긴다."** 작전의 마지막 장(章)이다.

### 3.2 작업 인계 및 기존 진척 확인

Auto는 `list_tasks`로 기존 작업을 확인하고, 멈춰 있던 작업 1·2의 그래프와 취약점 현황
(`get_task_graph`, `list_task_findings`)을 조회해 **이전 진척을 이어받았다**. 또한 대상에 대한
간단한 라이브 정찰(`curl -sI`)로 DVWA 로그인 페이지·보안 쿠키를 재확인했다.

*(그림 1 — 대화 시작: 사용자 지시 + 기존 paused 작업 인계, 모델 = OpenRouter GLM 5.3 Flash)*

### 3.3 증거 검증 + 라이브 재검증

Auto는 이전 세션의 증거 47개 파일을 검증(`grep`/`ls`)해 핵심 4건(SQLi 덤프·RCE·업로드 쉘·LFI)을
확인한 뒤, **신선한 증거 확보를 위해 라이브 재검증**을 수행했다. `admin/password` 로그인 성공,
보안레벨 low 확인 후 핵심 취약점을 실제로 다시 재현했다.

*(그림 2 — 증거 검증 → 라이브 재검증: 로그인/보안레벨/취약점 재현)*

**라이브 재검증 결과 (본 세션, curl 기반):**

| 검증 항목 | 결과 |
|-----------|------|
| 로그인(admin/password) | 성공(302 + Welcome 문구) |
| 보안레벨 low 설정 | 성공(`<em>low</em>`) |
| SQLi UNION 덤프 | 재현 — 5계정 MD5 해시 동일 |
| Command Injection | 재현 — `uid=33(www-data)` |
| 업로드 웹셸 | 재현 — 여전히 활성(`cmd=id` 응답) |
| Reflected XSS | 재현 — 무인코딩 반사 |
| LFI `/etc/passwd` | 재현 — 전체 열람 |
| 디렉터리 리스팅 | 재현 — `artex_shell.php` 노출 |

### 3.4 취약점 등록 — report_finding 15건

라이브 재검증 완료 후 Auto는 플랫폼에 취약점을 등록했다. 먼저 `list_assets`로 자산 id=3
(DVWA 서비스)을 확인하고, 15건을 `report_finding`으로 등록했다.

*(그림 3 — 라이브 재검증 완료 → 자산 확인 → report_finding 15건 등록, report_finding 도구 트리거 연속 실행)*

> **가드레일 검증** — 대화 컨텍스트에서 `report_finding`을 직접 호출하면 **"작업 컨텍스트
> 필요"로 정상 거부**된다. ARTEX 설계대로 `spawn_task`(작업 3) 생성 후 `add_task_hint`로 등록
> 규격을 인계하고, **작업 워커가 report_finding을 수행**하도록 흐름이 강제된다(가드레일 동작 확인).

### 3.5 ARTEX 자동 파이프라인 (작업 3 워커, 713초, done)

- **작업 3 생성**: `spawn_task(parent_ref=1, source_task_ids=[1,2], seed_first_intent=true, timeout 1800s)`
  — 작업 1·2의 goal·hints·intent 결과를 **상속**(related_tasks에 source_task_id 표기).
- **힌트 인계**: `add_task_hint`로 등록 규격 4건 일괄 주입(반환 ids 13~16) — Critical 3 / High 2 / Medium 8 / Low 2.
- **planner 동작**: 의도 생성 → 목표를 2건으로 분해(①15건 report_finding 등록 ②Critical/High 상세보고서).
- **worker 실행**: 증거 검증 → 상세보고서 파일 작성(`/app/data/tasks/3/i12/report_01~05.md`) →
  **report_finding으로 15건 등록 완료** → 작업 상태 **done**(713초), 전건 `state=confirmed`.
- **기획자 실시간 교정 훅**: 워커가 힌트 규격 외 방향(CSP 헤더 라이브 확인)으로 확장하려 하자
  **Bash 호출을 차단하고 힌트 규격 반영을 지시** — 방향 유지 제어가 실시간으로 동작함을 확인.
- **상세보고서 등록**: `update_finding_report`로 5건(노드 19~23) 등록(2.3~11.2KB).

---

## 4. 발견 취약점 요약 (15건, 등록 완료)

가장 심각한 공격 체인은 두 갈래다.

- `SQL Injection → DB 계정·해시 덤프 → 크랙 → 관리자/일반계정 로그인 재현`
- `파일 업로드 / 명령 인젝션 → www-data 권한 원격 코드 실행(RCE)`

| # | 취약점 | 등급 | 핵심 증거 |
|---|--------|------|-----------|
| 1 | SQL Injection — users 전체 덤프 | **Critical** | `UNION SELECT user,password FROM users` → 5계정 MD5, 크랙·실로그인 재현 |
| 2 | Command Injection → RCE | **Critical** | `;id` → `uid=33(www-data)` |
| 3 | Unrestricted File Upload → RCE | **Critical** | 웹셸 업로드·실행 → www-data (재검증 시점에도 활성) |
| 4 | File Inclusion (LFI + php://filter) | High | `/etc/passwd` 열람, 소스 base64 유출 |
| 5 | Blind SQL Injection (시간기반) | High | `SLEEP(5)` → 5.003s vs 0.003s |
| 6 | Stored XSS | Medium | 방명록 영구 저장 |
| 7 | Reflected XSS | Medium | name 파라미터 무인코딩 반사 |
| 8 | DOM XSS | Medium | `document.write` 무새니제이션 |
| 9 | CSRF (비밀번호 변경) | Medium | 토큰·현비밀번호 검증 없음 |
| 10 | Weak Session ID | Medium | `dvwaSession=1,2,3` 순차 예측가능 |
| 11 | Brute Force 무차단 | Medium | 실패 무제한·무지연, admin:password 유효 |
| 12 | CSP Bypass | Medium | 신뢰불가 외부 소스(pastebin) 허용 |
| 13 | JS 클라이언트 토큰 변조 | Medium | 서버 미검증 |
| 14 | 디렉터리 리스팅 | Low | uploads/ 등 Index 노출 |
| 15 | 설정정보/경로 노출 | Low | captcha가 설정파일 절대경로 노출 |

**크랙된 계정 (SQLi → 크랙 → 로그인 체인 검증 완료)**: `admin/password` · `gordonb/abc123`(재로그인 성공) ·
`1337/charley` · `pablo/letmein` · `smithy/password` — 무솔트 MD5 저장.

---

## 5. ARTEX 플랫폼 검증 결과

이번 유료 GLM 5.3 Flash 세션으로 **오케스트레이션부터 발견 등록까지 전 파이프라인이
엔드투엔드로 검증**되었다.

| 기능 | 테스트 내용 | 결과 |
|------|-------------|------|
| 작업 오케스트레이션 | `list_tasks`로 부모(1)·서브(2) 상태 식별 | ✓ paused 식별·인계 |
| 그래프/결과 조회 | `get_task_graph`, `list_task_findings` | ✓ 힌트·의도·진척 확인 |
| 자산 관리 | `list_assets(id=3)` — DVWA 서비스 자산 | ✓ 기술스택 지문 포함 |
| **취약점 등록 가드레일** | 대화 컨텍스트 `report_finding` 직접 호출 | ✓ 정상 거부(작업 컨텍스트 필요) |
| 서브태스크 생성·상속 | `spawn_task(source_task_ids=[1,2])` | ✓ 작업 3 생성·내용 상속 |
| 힌트 인계 | `add_task_hint` 4건(ids 13~16) | ✓ 등급별 규격 전달 |
| **워커 자동 등록** | 작업 3 워커가 `report_finding` 15건 등록 | ✓ 전건 confirmed, done(713초) |
| **기획자 실시간 교정 훅** | 워커의 힌트 외 확장 시도 차단 | ✓ Bash 차단 + 규격 반영 지시 |
| 상세보고서 | `update_finding_report`(노드 19~23) | ✓ 5건 등록(2.3~11.2KB) |

**개선 제안(관찰된 이슈)**: `update_finding_report`에 독립 레코드 id(1~5) 전달 시 "기록을 찾을 수
없음" 오류가 발생하고, 탐색 노드 id(19~23)로만 정상 갱신되었다 — 문서와 실동작 간 **id 규격
불일치**를 기록(개선 제안).

---

## 6. 대상 보안 권고 (DVWA 맥락, 일반 권고)

1. **SQLi** — Prepared Statement(PDO) + 최소권한 DB 계정
2. **명령 인젝션** — `shell_exec` 제거, 검증된 인자만 처리, 셸 호출 금지
3. **파일 업로드** — 확장자/MIME 화이트리스트, 업로드 경로 PHP 실행 차단, 실행권한 제거
4. **XSS/CSRF/토큰** — 출력 인코딩, 서버사이드 CSRF 토큰, 토큰 서버 검증
5. **세션** — 세션 토큰 CSPRNG 랜덤화, 로그인 실패 차단/지연(lockout+backoff)
6. **정보 노출** — 디렉터리 리스팅 비활성화(`Options -Indexes`), 오류 메시지 절대경로 마스킹
7. **잔존 웹셸 제거** — `/hackable/uploads/artex_shell.php` 즉시 삭제 및 업로드 디렉터리 정화

---

## 7. 결론

이번 세션은 ARTEX가 **단일 대화 지시만으로 침투 테스트 작업을 자율 생성·실행·추적하고,
라이브 재검증을 거쳐 발견을 플랫폼에 정식 등록(report_finding 15건)하기까지 전 파이프라인을
엔드투엔드로 완주**함을 확인했다. 특히 유료 **GLM 5.3 Flash**로 수행하여 **레이트리밋 중단 없이**
작업 생성·상속·힌트 인계·워커 자동 등록·기획자 실시간 교정 훅·상세보고서 등록이 모두 정상
동작했다. 그 결과 Critical 3건을 포함한 **15건 취약점을 confirmed 상태로 등록**했다.

> 무료 티어에서 관측됐던 레이트리밋 중단은 **초기 테스트의 일시적 현상**이었으며, 유료
> 모델 전환으로 해소됨을 이번 완주로 실증했다. 관찰된 `update_finding_report` id 규격
> 불일치만 개선 과제로 남는다.

### 7.1 방어 관점 — 관제센터는 ARTEX/LLM 해킹을 어떤 특징으로 탐지하는가

공격 과정에서 드러난 ARTEX의 행위·트래픽을 방어(SOC) 시각으로 뒤집어 보면,
**취약점이 아니라 "공격 주체(자율 에이전트)"를 탐지**해야 한다는 결론에 이른다. LLM이 응답을
읽고 다음 수를 조립하는 이상 페이로드 시그니처는 빗나가기 쉬우며, **아키텍처가 강제하는
구조적 지문**이 더 견고한 단서가 된다. 관제가 주목할 핵심 특징을 **내구성(공격자가 바꾸기
어려운 정도)** 순으로 정리하면 다음과 같다.

| 탐지 특징 | 관제가 무엇을 보는가 | 내구성 |
|-----------|----------------------|--------|
| **적응형 다단계 체인** | 민감·대용량 응답 → 동일 IP의 **논리적 후속(인증 성공·권한 상승)**이 초~분 내 발생 — **모델 등급과 무관**한 LLM 고유의 자율 체인 | **견고 · Critical** |
| **단일 프록시 집중 egress** | 한 소스 IP에서 정찰·인증·주입·업로드 등 **서로 다른 벡터 ≥ 3종**이 한 세션창에 혼재 | 중간 |
| **혼재된 클라이언트 정체성** | UA 리터럴 `artex-enrich/1.0`, Go·curl·모바일 Chromium의 **UA·JA3 혼재**, **모바일 UA ↔ 데이터센터 IP 불일치** | 중간·취약 |
| **저소음·고정밀** | 요청 수는 적으나 **민감 엔드포인트 적중률이 비정상적으로 높은** 세션 | 중간 |
| **LLM 호출 리듬** | tool-call마다 추론 지연이 끼는 버스트+갭(스캐너의 연속 폭주와 구별). 단 **수 분 단위 침묵은 모델 쿼터 한도에 좌우**되므로 유료 모델에서는 약해진다 | 가변 |

> **유의 — 레이트 리밋은 '만능 차단책'이 아니다** — 본 세션이 바로 그 증거다. **유료 GLM 5.3
> Flash에서는 레이트리밋 중단 없이 공격이 완주**되었다. 즉 방어측의 속도 제한은 공격자의
> 토큰 비용·소요 시간을 끌어올리는 **심층 방어의 한 겹일 뿐, 결정적 차단책이 아니다.** 모델
> 등급과 무관하게 남는 가장 견고한 단서는 **"민감 응답 → 논리적 후속"의 적응형 체인(의미
> 상관)**이며, 탐지는 여기에 다지표 상관과 인증 강화·자격증명 위생 같은 심층 방어를 더해
> 구성해야 한다. 세부 탐지룰·SIEM 상관·플레이북은 별도 문서
> 「ARTEX 에이전트 탐지·관제센터 대응방안」에 정리되어 있다.

---

*본 문서는 ARTEX 콘솔의 테스트 세션 로그(`pentest_log.md`, `pentest_report.md`)와 콘솔
스크린샷(그림 1~3)을 근거로, 유료 GLM 5.3 Flash 수행 세션을 Auto 대화 흐름 중심으로 정리한
보고서입니다. 대상 사이트 IP는 보안상 `100.11.*.*`로 마스킹하였으며, 대상 DVWA는 보안
교육·검증 전용 환경입니다.*
