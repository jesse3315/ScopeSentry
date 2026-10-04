---

## name: scopesentry-mcp
description: ScopeSentry MCP를 통해 보안 스캔 플랫폼(프로젝트, 작업, 템플릿, 자산, 노드)을 관리합니다. 사용자가 ScopeSentry, MCP, API Key, 스캔 작업, 자산 조회를 언급할 때 사용합니다.

# ScopeSentry MCP 사용 가이드

**이미 배포된 ScopeSentry 인스턴스**를 사용하는 사용자를 위한 가이드입니다. Cursor(또는 다른 MCP 클라이언트)로 플랫폼에 연결하며, 로컬 소스 코드는 필요하지 않습니다.

## 1. 준비 작업

### 1.1 서비스 접근 가능 여부 확인

- 기본 Web 인터페이스: `http://<호스트>`
- MCP 엔드포인트: `http://<호스트>/mcp` (앞단에 리버스 프록시나 프런트엔드 프록시가 있다면 실제 `/mcp` 주소를 기준으로 합니다)

### 1.2 API Key 생성

1. 브라우저에서 ScopeSentry Web 인터페이스에 로그인
2. **구성 → API 키 관리** 페이지에서 키를 생성 (또는 관리자가 제공한 API로 생성)
3. 반환된 `ssk_...` 문자열을 저장 (**한 번만 표시됨**)

### 1.3 Cursor MCP 구성

Cursor → Settings → MCP → 서버 추가:

```json
{
  "mcpServers": {
    "scopesentry": {
      "url": "http://<호스트>:8082/mcp",
      "headers": {
        "X-API-Key": "ssk_발급받은키"
      }
    }
  }
}
```

다음 방식도 사용할 수 있습니다: `Authorization: Bearer ssk_발급받은키`

구성을 마친 뒤 MCP를 재시작하거나 Cursor를 다시 로드하고, 도구 목록에 `list_projects`, `list_assets` 등이 나타나는지 확인합니다.

---

## 2. 도구 목록


| 도구                     | 용도                |
| ---------------------- | ----------------- |
| `list_projects`        | 태그별로 그룹화된 프로젝트 트리(프로젝트 ID 포함) |
| `list_projects_data`   | 페이지 단위 프로젝트 목록, 이름으로 검색 가능     |
| `get_project`          | 프로젝트 상세              |
| `create_project`       | 프로젝트 생성              |
| `list_tasks`           | 스캔 작업 목록            |
| `get_task`             | 작업 상세              |
| `list_scan_templates`  | 스캔 템플릿 목록            |
| `get_scan_template`    | 템플릿 상세              |
| `list_plugin_modules`  | 스캔 파이프라인 모듈 이름          |
| `list_plugins`         | 사용 가능한 플러그인(hash, 기본 파라미터 포함) |
| `create_scan_template` | 스캔 템플릿 생성            |
| `create_scan_task`     | 스캔 작업 생성            |
| `list_assets`          | 각종 자산 조회(페이지 단위 목록)       |
| `count_assets`         | 자산 수 집계(`/api/assets/common/total`) |
| `get_asset_detail`     | 자산 또는 취약점 상세           |
| `add_asset_tag`        | 자산에 태그 추가           |
| `list_nodes`           | 스캔 노드 목록            |


각 도구의 파라미터는 MCP 도구 설명(schema)을 기준으로 합니다. `list_assets` / `count_assets`의 search, filter 문법은 동일하므로 자산을 조회하기 전에 `list_assets` description을 먼저 읽어 보세요.

"총 몇 건인지"가 필요할 때는 `count_assets`(Web 페이지네이션 총계 API에 해당)를 사용하세요. 총계를 세기 위해 `list_assets`로 페이지를 반복해서 넘길 필요가 없습니다.

---

## 3. 자주 쓰는 워크플로

### 3.1 프로젝트별 자산 조회

사용자나 컨텍스트에 **이미 프로젝트 조건이 있는** 경우, `filter.project`를 우선 지정해 범위를 좁히세요. 여러 프로젝트에 걸친 데이터가 너무 많아 응답이 느려지는 것을 방지합니다. 명확한 프로젝트가 없다면 프로젝트 필터를 강제로 지정하지 않아도 됩니다.

1. `list_projects` 또는 `list_projects_data`로 대상 프로젝트의 **ObjectID**(`id` / `children[].value`)를 가져옴
2. `list_assets`에 `filter.project`를 전달 (**반드시 ID여야 하며, 프로젝트 이름을 쓰면 안 됨**)

```json
{
  "asset_type": "asset",
  "pageIndex": 1,
  "pageSize": 20,
  "search": "domain=^example.com",
  "filter": {
    "project": ["<프로젝트 ObjectID>"]
  }
}
```

### 3.2 스캔 작업 생성

1. `list_nodes`로 온라인 노드 이름을 가져옴
2. `list_scan_templates` 또는 `create_scan_template`로 템플릿 **ObjectID**를 가져옴
3. `create_scan_task`: `name`, `node`는 필수이며, `template`에는 템플릿 ID를 입력 (템플릿 이름은 불가)

**대상 소스 `targetSource` (Web과 동일):**

| targetSource | 설명 | 필수 파라미터 |
| --- | --- | --- |
| `general` | 대상을 직접 입력 | `target` |
| `project` | 프로젝트에서 대상을 읽음 | `project`(프로젝트 ObjectID 배열) |
| `asset` | Web 자산 저장소에서 검색 | `search`; 선택 `project`, `filter`, `targetNumber` |
| `RootDomain` | 루트 도메인 저장소에서 검색 | `search`; 선택 `project`, `filter`, `targetNumber` |
| `subdomain` | 서브도메인 저장소에서 검색 | `search`; 선택 `project`, `filter`, `targetNumber` |
| `UrlScan` | URL 스캔 결과에서 검색 | `search`; 선택 `project`, `filter`, `targetNumber` |
| `*Source`(예: `subdomainSource`) | 자산 페이지의 "선택/검색"으로 생성 | `targetTp=search`이면 `search`, `targetTp=select`이면 `targetIds` 사용 |

**예시 - 루트 도메인 직접 스캔:**

```json
{
  "name": "example-서브도메인수집",
  "node": ["node-1"],
  "template": "<템플릿 ObjectID>",
  "targetSource": "general",
  "target": "example.com\nfoo.com",
  "project": ["<프로젝트 ObjectID>"]
}
```

**예시 - 서브도메인 저장소에서 이어서 스캔 (이전 작업 이름으로 필터링):**

```json
{
  "name": "example-포트및취약점",
  "node": ["node-1"],
  "template": "<후속 모듈 템플릿 ObjectID>",
  "targetSource": "subdomain",
  "search": "task==\"example-서브도메인수집\"",
  "project": ["<프로젝트 ObjectID>"]
}
```

### 3.3 루트 도메인 전체 정보 수집 (2단계 권장)

입력이 **루트 도메인**이고 **전체 정보 수집**을 하려는 경우, 전체 파이프라인을 한 번에 실행하지 말고 두 번에 나눠 스캔하는 것을 권장합니다.

**이유:** 분산 작업은 **개별 대상** 단위로 분배됩니다. 루트 도메인을 대상으로 하면, 해당 루트 도메인을 할당받은 노드에서 발견된 서브도메인도 같은 노드에서 후속 모듈이 계속 실행되므로 부하 불균형, 속도 저하, 오류가 발생하기 쉽습니다.

**모범 사례:**

1. **1단계 - 서브도메인 수집만**
   - `targetSource`: `general`
   - `target`: 모든 루트 도메인(여러 줄)
   - 템플릿: `SubdomainScan`, `SubdomainSecurity`만 활성화 (서브도메인 스캔 + 서브도메인 탈취)
   - `get_task`로 작업 완료를 기다림

2. **2단계 - 후속 모듈**
   - `targetSource`: `subdomain`
   - `search`: `task=="<1단계 작업 이름>"` (작업 이름 정확히 일치)
   - 선택적으로 `project`로 범위 축소
   - 템플릿: 포트 스캔, 자산 매핑, 취약점 스캔 등 (SubdomainScan 제외 가능)
   - 서브도메인이 독립 대상으로 각 노드에 분배되어 병렬 효율이 더 높음

Web 인터페이스의 "서브도메인" 자산 페이지에서 작업 이름으로 필터링한 뒤 "서브도메인에서 작업 생성"을 사용해도 같은 효과를 얻을 수 있습니다.

```mermaid
flowchart LR
  A[루트 도메인 목록] --> B[1단계: general + SubdomainScan]
  B --> C[서브도메인 저장]
  C --> D[2단계: subdomain + task==1단계 작업 이름]
  D --> E[포트/자산/취약점 등 모듈]
```

### 3.4 스캔 템플릿 생성

1. `list_plugin_modules` → 모듈 이름 목록
2. `list_plugins`(`module`로 필터링 가능) → 각 플러그인의 `hash`와 기본 `parameter`
3. `create_scan_template`: `modules`로 "모듈 → 플러그인 hash 배열"을 지정

---

## 4. 자산 조회 (`list_assets` / `count_assets`)

`count_assets`는 `list_assets`와 같은 `asset_type`, `search`, `filter`를 사용하며 `{ "total": N }`을 반환합니다. Web의 `/api/assets/common/total`에 해당합니다.

```json
{
  "asset_type": "subdomain",
  "search": "task==\"작업이름\"",
  "filter": {"project": ["<프로젝트 ObjectID>"]}
}
```

**성능 권장 사항 (`list_assets` / `count_assets` 공통):** 프로젝트 조건이 있으면 `filter.project`로 범위를 우선 좁히세요. `search`에서 인덱스가 있는 필드는 가능한 한 `==` 완전 일치나 `^` 접두사 일치를 사용하고([4.3](#43-search-검색-표현식) 참고), 응답을 느리게 만드는 광범위한 `=` 퍼지 검색은 피하세요. 프로젝트 컨텍스트가 없으면 프로젝트 필터를 강제하지 않습니다.

`filter.project`를 지원하는 유형은 [4.4](#44-filter-정확-필터) 표를 참고하세요.

### 4.1 자산 유형 `asset_type`

`asset`, `RootDomain`, `subdomain`, `app`, `mp`, `UrlScan`, `SensitiveResult`, `DirScanResult`, `crawler`, `vulnerability`, `PageMonitoring`, `IPAsset`, `SubdomainTakerResult`

별칭 예시: `web`→asset, `vuln`→vulnerability, `ip`→IPAsset, `url`→UrlScan

### 4.2 파라미터 설명


| 파라미터                       | 설명                                      |
| ------------------------ | --------------------------------------- |
| `pageIndex` / `pageSize` | 페이지네이션, 기본값 1 / 20                            |
| `search`                 | 검색 표현식(다음 절 참고)                              |
| `filter`                 | 정확 필터 JSON(다음 절 참고)                          |
| `sort`                   | UrlScan, DirScanResult만 `length` 기준 정렬 지원 |
| `sid`                    | SensitiveResult 전용: 민감 정보 규칙 이름                |


`search`와 `filter`는 **함께 사용할 수 있습니다**.

### 4.3 search 검색 표현식

자체 DSL(**SQL이 아님**):


| 연산자  | 의미   | 인덱스 | 예시                          |
| ---- | ---- | ---- | --------------------------- |
| `=`  | 퍼지 일치(regex) | 인덱스 미사용 | `domain=example`            |
| `==` | 정확 일치(완전 일치) | **인덱스 사용** | `port==443`                 |
| `!=` | 제외   | — | `port!="80"`                |
| `&&` | AND    | — | `domain==example.com && port==443` |
| `||` | OR    | — | `title=admin || body=login` |


**인덱스와 연산자:** `domain`, `ip`, `port`, `title` 등의 필드에는 인덱스가 있지만, **`==` 완전 일치** 또는 **값이 `^`로 시작하는 접두사 일치**(예: `domain=^example.com`)만 인덱스를 사용합니다. **`=`는 regex 퍼지 일치로 변환되어 인덱스를 사용할 수 없으므로** 데이터가 많으면 느려지기 쉽습니다.

**모든 유형 공통 search 필드:** `tag`, `task`(작업 이름), `rootDomain`

**project는 search에 쓸 수 없습니다** (무효이거나 `&&`와 조합 시 오류 발생). 프로젝트 필터링은 `filter.project`를 사용하세요.

**유형별 주요 search 필드:**


| asset_type           | 필드                                                                                  |
| -------------------- | ----------------------------------------------------------------------------------- |
| asset                | domain, ip, port, service, app, title, statuscode, icon, banner, type, body, header |
| RootDomain           | domain, icp, company                                                                |
| subdomain            | domain, ip, type, value                                                             |
| app                  | name, icp, company, category, description, url, apk                                 |
| mp                   | name, icp, company, category, description, url                                      |
| UrlScan              | url, input, source, resultId, type                                                  |
| SensitiveResult      | url, sname, body, info, md5                                                         |
| DirScanResult        | url, statuscode, redirect, length                                                   |
| vulnerability        | url, vulname, matched, request, response, level                                     |
| crawler              | url, method, body, resultId                                                         |
| PageMonitoring       | url, hash, diff, response                                                           |
| IPAsset              | ip, domain, port, service, webServer, app                                           |
| SubdomainTakerResult | domain, value, type, response                                                       |


**search 예시:**

- `domain==www.example.com && port==443` (완전 일치, 인덱스 사용)
- `domain=^example.com` (접두사 일치, 인덱스 사용)
- `ip==192.168.1.1`
- `task=="작업이름"`
- `level==high` (vulnerability)
- `statuscode==200` (DirScanResult)

부분 포함 검색이 꼭 필요할 때만 `=`를 사용하세요. 예: `title=admin` (인덱스 미사용이므로 프로젝트 등의 조건과 함께 범위를 좁히는 것이 좋습니다).

### 4.4 filter 정확 필터

JSON 객체: 같은 key의 여러 값은 **OR**, 서로 다른 key는 **AND**입니다.

**프로젝트 조건이 있으면 `project`를 우선 사용:** 사용자나 컨텍스트에서 프로젝트가 명확하고 asset_type이 `project`를 지원하면 지정해서 범위를 좁히세요. 프로젝트 정보가 없으면 강제하지 않습니다.


| filter key   | 의미       | 값 설명                                                     |
| ------------ | -------- | -------------------------------------------------------- |
| `project`    | 소속 프로젝트     | **ObjectID**, `list_projects` / `list_projects_data`로 조회 |
| `task`       | 출처 작업     | **작업 이름**, `list_tasks`의 `name`                         |
| `port`       | 포트       | 예: `"443"`                                                |
| `service`    | 서비스/프로토콜    | 예: `"https"`                                              |
| `app`        | 애플리케이션 핑거프린트     | 예: `"Nginx"`                                              |
| `icon`       | 아이콘 hash  |                                                          |
| `statuscode` | HTTP 상태 코드 | 주로 asset에 사용                                               |
| `status`     | 상태       | UrlScan/DirScan HTTP 코드; 취약점/민감 정보 처리 상태                       |
| `level`      | 취약점 등급     | critical / high / medium / low / info                    |
| `type`       | 유형       | 예: 서브도메인 레코드 유형 A, CNAME                                         |
| `color`      | 민감 정보 규칙 색상   | SensitiveResult                                          |
| `sname`      | 민감 정보 규칙 이름    | SensitiveResult                                          |
| `tags`       | 태그       |                                                          |


**유형별 사용 가능한 filter key:**


| asset_type                            | filter key                                                      |
| ------------------------------------- | --------------------------------------------------------------- |
| asset                                 | project, port, service, app, icon, statuscode, type, task, tags |
| RootDomain                            | project, tags                                                   |
| subdomain                             | project, type, task, tags                                       |
| app / mp                              | project, tags                                                   |
| UrlScan                               | status, tags                                                    |
| DirScanResult                         | status, tags                                                    |
| SensitiveResult                       | status, color, sname, tags                                      |
| crawler                               | project, task, tags                                             |
| vulnerability                         | project, level, status, task, tags                              |
| PageMonitoring / SubdomainTakerResult | tags                                                            |
| IPAsset                               | project, port, service, app                                     |


**filter 예시:**

```json
{"project": ["<프로젝트 ObjectID>"], "port": ["443"]}
```

**조합 조회 예시:**

```json
{
  "asset_type": "asset",
  "search": "domain=^baidu && port==443",
  "filter": {"project": ["<프로젝트 ObjectID>"]},
  "pageIndex": 1,
  "pageSize": 10
}
```

**주의:**

- 프로젝트 조건이 있으면 `filter.project`를 우선 지정(지원되는 경우). 프로젝트 컨텍스트가 없으면 강제하지 않아도 됨
- `filter.project`에 프로젝트 표시 이름을 넣지 말 것
- 알고 있는 값은 `==`, 접두사는 `^` 사용. 큰 테이블에 `=` 퍼지 일치를 남용하지 말 것
- UrlScan의 HTTP 상태는 `filter.status` 사용. DirScanResult는 search에서 `statuscode==200` 사용 가능
- SensitiveResult를 규칙 이름으로 조회: `search`에서 `sname=규칙이름`, 또는 `filter.sname`

### 4.5 정렬 sort

**UrlScan**, **DirScanResult**만 지원:

```json
{"length": "ascending"}
```

다른 유형은 `sort`를 무시하고 기본적으로 시간순으로 정렬합니다.

---

## 5. 스캔 템플릿 모듈 이름

`TargetHandler`, `SubdomainScan`, `SubdomainSecurity`, `PortScanPreparation`, `PortScan`, `PortFingerprint`, `AssetMapping`, `AssetHandle`, `URLScan`, `WebCrawler`, `URLSecurity`, `DirScan`, `VulnerabilityScan`, `PassiveScan`

---

## 6. 문제 해결


| 현상        | 조치                                                 |
| --------- | -------------------------------------------------- |
| MCP 도구가 없음   | URL, API Key, ScopeSentry 실행 여부 확인                    |
| 401 / 403 | API Key를 다시 생성하거나 교체                                    |
| 자산이 조회되지 않음     | `filter.project`가 ObjectID인지 확인. search에 project를 쓰지 말 것 |
| 템플릿/작업 생성 실패 | `template`은 반드시 템플릿 ObjectID여야 함. `node`에는 온라인 노드 이름 입력            |
| 조회가 매우 느림/멈춤   | 프로젝트가 있으면 `filter.project` 추가. search에서 인덱스 필드는 `==` 또는 `^` 접두사를 사용하고 `=` 사용 최소화. `pageSize` 축소 |


---

