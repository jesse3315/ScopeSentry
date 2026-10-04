package mcp

// 자산 조회 관련 도구 설명 (Tool.Description에 기록되어 MCP 클라이언트에 표시됨)
const listAssetsToolDesc = `ScopeSentry 자산 목록을 조회합니다.

[엔드포인트와 요청 본문]
각 자산의 POST 엔드포인트와 동일 (예: /api/assets/asset), body는 models.SearchRequest:
- pageIndex, pageSize: 페이지네이션
- search: 검색 표현식 문자열 (프런트엔드 Csearch 검색창의 searchParams에 대응)
- filter: 정확 필터 객체 (프런트엔드의 프로젝트/작업 드롭다운, 통계 사이드바 클릭, 테이블 열 필터에 대응하며 병합되어 전달됨)
- sort: 정렬 (대부분의 유형에서 무효)
- sid: SensitiveResult 전용, 민감 정보 규칙 이름

MCP의 asset_type은 백엔드 Index로 매핑되며 (예: asset→"asset"), 이후 helper.GetSearchQuery를 호출해 MongoDB 쿼리를 생성합니다.

[프런트엔드 파라미터 역할 분담 (Csearch.vue + 각 자산 페이지)]
- search: 사용자가 검색창에 입력하는 DSL, 예: domain=baidu && port==443
- filter.project: ElTreeSelect 선택값, 프로젝트 ObjectID 배열 (프로젝트 이름 아님)
- filter.task: ElSelect 선택값, 작업 이름 배열 (작업 ID 아님); 동적 태그 task=이름 으로 filter에 기록할 수도 있음
- filter는 다음에서도 올 수 있음: 통계 사이드바 클릭(port/service/app/icon), 테이블 열 필터(statuscode/level/type/status 등)
- 주의: 각 자산 페이지의 searchKeywordsData에 project 힌트가 있지만, 백엔드 SearchToMongoDB에 project 키워드가 등록되어 있지 않으므로 project는 filter로만 지정 가능

asset_type 선택 가능 값:
- asset: Web/포트 자산
- RootDomain: 루트 도메인
- subdomain: 서브도메인
- app: 모바일 앱
- mp: 미니 프로그램
- UrlScan: URL
- SensitiveResult: 민감 정보
- DirScanResult: 디렉터리 스캔
- crawler: 크롤러
- vulnerability: 취약점
- PageMonitoring: 페이지 모니터링
- IPAsset: IP 집계 자산
- SubdomainTakerResult: 서브도메인 탈취

[search 검색 표현식] (SQL 아님, 커스텀 DSL이며 helper.SearchToMongoDB가 파싱)
- field=value : 부분 일치 (대소문자 구분 없음)
- field=="값" : 정확 일치, 공백이 포함된 값은 큰따옴표로 감싸야 함
- field!="값" : 제외 일치
- expr1 && expr2 : AND
- expr1 || expr2 : OR
- (expr) : 그룹
- 공통 search 필드 (모든 유형, SearchToMongoDB가 주입): tag→tags, task→taskName (작업 이름), rootDomain
- project는 search를 지원하지 않으며 filter.project로만 지정 가능 (값은 프로젝트 ObjectID, list_projects / list_projects_data로 조회)

각 asset_type에서 사용 가능한 search 키워드 → MongoDB 필드:
- asset: domain→host, ip, port, service, app→technologies, title, statuscode, icon→faviconmmh3, banner→metadata, type, body, header→rawheaders
- RootDomain: domain, icp, company
- subdomain: domain→host, ip, type, value
- app: name, icp, company, category, description, url, apk
- mp: name, icp, company, category, description, url
- UrlScan: url→output, input, source, resultId, type→outputtype (statuscode search 없음; HTTP 상태 코드는 filter.status만 가능)
- SensitiveResult: url, sname→sid, body, info→match, md5
- DirScanResult: url, statuscode→status, redirect→msg, length
- vulnerability: url, vulname, matched, request, response, level
- crawler: url, method, body, resultId
- PageMonitoring: url, hash, diff, response
- IPAsset: ip, domain→ports.server.domain, port→ports.port, service→ports.server.service, webServer→ports.server.webServer, app→ports.server.technologies
- SubdomainTakerResult: domain→input, value, type→cname, response

search 예시:
- domain=baidu && port==443
- port==443 && service=nginx
- title="관리자 로그인" || body=admin

[filter 정확 필터] (search와 조합 가능; helper.GetSearchQuery가 $and 조건으로 추가)
JSON 객체이며 key는 필터 차원, value는 문자열 배열 (같은 key의 여러 값은 OR, 서로 다른 key 간에는 AND).
filterKeyCache에 정의된 key만 적용되며 나머지 key는 무시됩니다.

filter 요점:
- project: filter 전용, 값은 프로젝트 ObjectID 배열 (DB에는 ID로 저장; 일부 엔드포인트는 표시 시 프로젝트 이름으로 변환)
- task: search (task=="작업 이름") 또는 filter.task (작업 이름 배열) 사용 가능; 값은 list_tasks의 name이며 작업 ID 아님
- 같은 key의 여러 값은 OR, 서로 다른 key는 AND; filterKeyCache에 없는 key는 무시됨

전역 filter key → MongoDB 필드:
- project→project, port→port, service→service, app→technologies
- icon→faviconmmh3, statuscode→statuscode, status→status
- level→level, type→type, color→color, tags→tags
- task→taskName, sname→sid

각 asset_type에서 사용 가능한 filter (나열되지 않은 key는 해당 컬렉션에서 무효이거나 필드가 없음):
- asset: project, port, service, app, icon, statuscode, type, task, tags
- RootDomain: project, tags
- subdomain: project, type, task, tags
- app / mp: project, tags
- UrlScan: status (HTTP 상태 코드, statuscode 아님), tags
- SensitiveResult: status (1 미처리/2 처리 중/3 무시/4 의심/5 확인/6 처리 완료), color, sname, tags
- DirScanResult: status (HTTP 상태 코드), tags
- crawler: project, task, tags
- vulnerability: project, level (critical/high/medium/low/info/unknown), status (1-6), task, tags
- PageMonitoring: tags
- IPAsset: project, port, service, app (중첩된 ports 필드, 집계 쿼리 사용)
- SubdomainTakerResult: tags

filter 예시 (project는 반드시 ObjectID여야 하며 프로젝트 이름은 사용 불가):
- asset: {"project":["<프로젝트ObjectID>"],"port":["443"]}
- subdomain: {"project":["<프로젝트ObjectID>"],"type":["A"]}
- vulnerability: {"project":["<프로젝트ObjectID>"],"level":["high"]}

search + filter 조합 예시:
- search: domain=baidu && port==443, filter: {"project":["<프로젝트ObjectID>"]}

주의:
- search에 project=... 또는 project=="..."를 쓰지 마세요 (무효이거나 &&와 조합 시 오류 발생)
- DirScanResult의 HTTP 상태 코드는 search에서만 statuscode로 사용 (예: statuscode==200)
- UrlScan에는 statuscode search가 없음; 상태 필터는 filter.status 사용
- SensitiveResult 규칙 이름 필터: search에서 sname=규칙이름, 또는 filter.sname

[sort 정렬]
- UrlScan / DirScanResult: {"length":"ascending"} 오름차순 지원, 그 외 값 ("descending", "-1" 포함)은 내림차순
- 그 외 유형: 서버에서 time 또는 _id 기준으로 고정 정렬, sort 파라미터는 보통 무효

[sid 파라미터]
SensitiveResult 전용: 민감 정보 규칙 이름(sid)을 전달하며, 규칙별로 매칭 상세를 펼치는 데 사용 (프런트엔드에서 규칙 이름 클릭에 대응)`

// createScanTemplateToolDesc 스캔 템플릿 생성 설명
const createScanTemplateToolDesc = `스캔 템플릿을 생성합니다.

스캔 템플릿은 여러 [모듈](파이프라인 단계)로 구성되며, 각 모듈에는 여러 [플러그인]이 연결됩니다.
모듈 필드는 플러그인의 hash를 참조합니다 (플러그인 id나 플러그인 이름이 아님).

권장 생성 절차:
1. list_plugin_modules로 전체 모듈 이름 조회
2. list_plugins (module로 필터 가능)로 각 모듈에서 사용 가능한 플러그인의 hash와 기본 parameter 조회
3. modules 파라미터로 모듈->플러그인hash 목록 지정 (같은 모듈 내에서는 배열 순서대로 실행)
4. 이 도구는 플러그인 기본 파라미터로 Parameters를 자동으로 채움; 사용자 정의 파라미터가 필요하면 parameters로 덮어쓰기

파라미터 설명:
- name: 템플릿 이름, 필수
- modules: {모듈이름: [플러그인hash,...]}, 예: {"SubdomainScan":["d60ba73c..."],"PortScan":["..."]}
- parameters: 선택, {모듈이름: {플러그인hash: 파라미터 문자열}}, 기본 파라미터를 덮어씀
- vullist: 선택, nuclei POC 템플릿 ID 목록 (VulnerabilityScan에서 nuclei 사용 시)
- template_json: 선택, 완전한 ScanTemplate JSON, 우선순위가 가장 높으며 고급 사용자 정의용

생성 성공 시 템플릿 id를 반환하며, create_scan_task의 template 파라미터에 바로 사용할 수 있습니다.

주요 모듈: TargetHandler, SubdomainScan, SubdomainSecurity, PortScanPreparation,
PortScan, PortFingerprint, AssetMapping, AssetHandle, URLScan, WebCrawler,
URLSecurity, DirScan, VulnerabilityScan, PassiveScan`

// createScanTaskToolDesc 스캔 작업 생성 설명 (Web의 /api/task/add 및 common.Insert 로직과 동일)
const createScanTaskToolDesc = `스캔 작업을 생성합니다. name, node는 필수; template은 스캔 템플릿 ObjectID (list_scan_templates로 조회).

[대상 소스 targetSource] 작업 대상을 어떻게 해석할지 결정 (internal/services/task/common/common.go에 대응):
- general: target 필드를 직접 사용 (여러 줄/쉼표로 구분된 도메인, IP, URL)
- project: 연결된 프로젝트에서 대상을 읽음, project (프로젝트 ObjectID 배열) 입력, target 불필요
- asset: Web 자산 DB에서 검색해 선택, search 필요; 선택 project, filter, targetNumber
- RootDomain: 루트 도메인 DB에서 검색해 선택, search 필요; 선택 project, filter, targetNumber
- subdomain: 서브도메인 DB에서 검색해 선택, search 필요; 선택 project, filter, targetNumber
- UrlScan: URL 스캔 결과에서 검색해 선택, search 필요; 선택 project, filter, targetNumber
- assetSource / RootDomainSource / subdomainSource / UrlScanSource: 해당 자산 페이지에서 작업 생성
  - targetTp=search: search + filter + project + targetNumber로 대상 필터링
  - targetTp=select: targetIds로 자산 ObjectID 목록 지정

[search] 검색 표현식, 문법은 list_assets와 동일 (예: task=="어떤작업이름", domain=^example.com).
서브도메인에서 이어서 스캔하는 예시: targetSource=subdomain, search=task=="서브도메인 수집 작업 이름"

[filter] 정확 필터 JSON, search와 조합 가능; filter.project는 프로젝트 ObjectID.
[targetNumber] search 모드에서의 대상 수 상한, 0은 제한 없음.
[targetIds] select 모드에서 선택한 자산 ObjectID 목록.

[기타 파라미터]
- allNode: 온라인 상태인 모든 노드를 자동 추가
- ignore / duplicates: 무시할 대상, 중복 제거 전략
- bindProject: 프로젝트 바인딩 (결과 귀속)
- scheduledTasks + cycleType/hour/minute/day/week: 예약 작업

루트 도메인 전체 정보 수집은 2단계 권장: 먼저 general + SubdomainScan/SubdomainSecurity만으로 루트 도메인 스캔;
완료 후 subdomain + search=task=="이전 작업 이름"으로 후속 모듈 작업 생성.`

// countAssetsToolDesc 자산 수 통계 (POST /api/assets/common/total에 대응)
const countAssetsToolDesc = `조건에 맞는 자산 수를 집계합니다 (Web 페이지네이션의 "총 N건"과 동일한 엔드포인트 /api/assets/common/total).

파라미터는 list_assets의 search/filter와 동일하지만 페이지네이션 없이 total만 반환합니다.
asset_type 값은 list_assets와 동일 (asset, RootDomain, subdomain, vulnerability 등).

예시: 특정 프로젝트의 서브도메인 수 집계
{"asset_type":"subdomain","filter":{"project":["<프로젝트ObjectID>"]}}

예시: 특정 작업에서 생성된 Web 자산 집계
{"asset_type":"asset","search":"task==\"어떤작업이름\""}`
