<div align=center>
	<img src="docs/images/favicon.ico"/>
</div>

[English](./README.md) | [中文](./README_CN.md) | 한국어

> 이 저장소는 [Autumn-27/ScopeSentry](https://github.com/Autumn-27/ScopeSentry)를 포크하여 한국어 UI를 추가한 버전입니다. 웹 UI의 기본 언어는 한국어이며, 우측 상단 언어 메뉴에서 English / 简体中文 으로 전환할 수 있습니다.

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Autumn-27/ScopeSentry-Scan)

## 소개
Scope Sentry는 자산 매핑, 서브도메인 열거, 정보 유출 탐지, 취약점 스캔, 디렉터리 스캔, 서브도메인 탈취, 크롤러, 페이지 모니터링 등의 기능을 갖춘 도구입니다. 여러 노드를 구성하여 원하는 노드에서 스캔 작업을 실행할 수 있습니다. 새로운 취약점이 등장하면 관심 자산에 관련 컴포넌트가 있는지 빠르게 확인할 수 있습니다.

분산 구현 참고 글: [https://mp.weixin.qq.com/s/xfgRxUjljoQ8KzacblktxA](https://mp.weixin.qq.com/s/xfgRxUjljoQ8KzacblktxA)

서버 추천: [lightnode](https://www.lightnode.com/?inviteCode=CQ11JU&promoteWay=LINK)

## Discord:

[https://discord.gg/GWVwSBBm48](https://discord.gg/GWVwSBBm48)

## 기술 스택
서버: go

스캐너: go

프론트엔드: vue - vue-element-plus-admin

## 웹사이트

- 공식 웹사이트: [https://www.scope-sentry.top](https://www.scope-sentry.top/en/)
- Github: [https://github.com/Autumn-27/ScopeSentry](https://github.com/Autumn-27/ScopeSentry)
- 스캐너 소스 코드: [https://github.com/Autumn-27/ScopeSentry-Scan](https://github.com/Autumn-27/ScopeSentry-Scan)
- UI 소스 코드: [https://github.com/Autumn-27/ScopeSentry-UI](https://github.com/Autumn-27/ScopeSentry-UI)
- 플러그인 마켓: [플러그인 마켓](https://plugin.scope-sentry.top/en)
- 플러그인 템플릿:[https://github.com/Autumn-27/ScopeSentry-Plugin-Template](https://github.com/Autumn-27/ScopeSentry-Plugin-Template)

## 설치
```
git clone https://github.com/jesse3315/ScopeSentry.git
cd ScopeSentry
# .env 파일에서 MongoDB와 Redis 계정 비밀번호를 변경하세요.
docker compose -f single-host-deployment.yml up -d --build
```
실행하면 mongodb, redis, scope-sentry(서버), scopesentry-scan(스캐너) 네 개의 컨테이너가 생성됩니다. 기본적으로 스캔 노드 1개가 포함됩니다.

> **한국어판 서버 이미지:** `scope-sentry` 서버는 원본 이미지(`autumn27/scopesentry`)를 받지 않고, 이 저장소의 소스로 `Dockerfile.source`를 사용해 직접 빌드합니다(`scopesentry-ko:local`). 첫 실행 시 프론트엔드와 Go 빌드 때문에 몇 분 정도 걸립니다. 소스를 업데이트(`git pull`)한 뒤에는 `--build` 옵션을 붙여 다시 실행하세요. 스캔 노드(`scopesentry-scan`)는 화면이 없으므로 원본 이미지를 그대로 사용합니다.

초기 사용자 비밀번호와 플러그인 2차 인증 비밀번호 확인

```docker logs scope-sentry```

**노드 추가 (선택 사항)**

```
git clone https://github.com/Autumn-27/ScopeSentry-Scan.git
cd ScopeSentry-Scan/build
# .env 파일에서 MongoDB와 Redis 연결 정보를 수정하세요. NodeName은 노드 이름이며 노드마다 고유해야 합니다(비워두면 무작위로 생성되며 웹 화면에서 이름을 변경할 수 있습니다).
docker-compose -f scan-docker-compose.yml up -d
```

## 플러그인 흐름도

<img src="流程图.svg"/>

## 현재 기능
- 플러그인 시스템 (확장을 통해 어떤 도구든 추가 가능)
- 서브도메인 열거
- 서브도메인 탈취 탐지
- 포트 스캔
- 자산 식별
- 디렉터리 스캔
- 취약점 스캔
- 민감 정보 유출 탐지
- URL 추출
- 크롤러
- 페이지 모니터링
- 사용자 정의 웹 핑거프린트
- POC 가져오기
- 자산 그룹화
- 다중 노드 스캔
- Webhook

## 예정
- 취약한 비밀번호 크래킹

## 설치 안내

자세한 설치 방법은 [공식 웹사이트](https://www.scope-sentry.top)를 참고하세요.

## 커뮤니티

Discord:

[https://discord.gg/agsYdAyN](https://discord.gg/agsYdAyN)


## 스크린샷

### 로그인

![alt text](docs/images/login.png)

### 홈 대시보드
![alt text](docs/images/index-en.png)

## 플러그인 시스템
![alt text](docs/images/plugin-cn.png)
![alt text](docs/images/plugin-m-en.png)
## 자산 데이터
### 자산
![alt text](docs/images/asset-en.png)
![alt text](docs/images/asset-s-en.png)
![alt text](docs/images/asset-s2-en.png)

### 빠른 구문 검색:
![alt text](docs/images/search.gif)

## 루트 도메인
![alt text](docs/images/rootdomain-cn.png)

### 서브도메인
![alt text](docs/images/subdomain-en.png)

### 서브도메인 탈취
![alt text](docs/images/subt-en.png)

### APP
![alt text](docs/images/app-cn.png)

### 미니 프로그램
![alt text](docs/images/mp-cn.png)

### URL
![alt text](docs/images/craw-cn.png)

### 크롤러
![alt text](docs/images/craw-en.png)

### 민감 정보
![alt text](docs/images/sns-cn.png)

### 디렉터리 스캔
![alt text](docs/images/dir-cn.png)

### 취약점
![alt text](docs/images/vul-en.png)

### 페이지 모니터링
![alt text](docs/images/page-en.png)
![alt text](docs/images/page-change.png)
## 프로젝트

![](docs/images/project-cn.png)

## 프로젝트 자산 집계
### 패널 - 개요
![](docs/images/project-dsh.png)
### 서브도메인
![docs/images/project-subdomain.png)
### 포트
![](docs/images/project-port.png)
### 서비스
![](docs/images/project-server.png)

## 작업

![](docs/images/create-task-en.png)

## 작업 진행 상황

![](docs/images/task-pg-en.png)

## 노드

![](docs/images/node-cn.png)


## 라이선스

이 프로젝트의 모든 브랜치는 AGPL-3.0을 따르며, 다음 추가 조항을 준수해야 합니다:
1. 이 소프트웨어를 상업적으로 사용하려면 별도의 상업용 라이선스가 필요합니다.
2. 기업, 조직 및 영리 단체는 이 소프트웨어를 사용, 배포 또는 수정하기 전에 상업용 라이선스를 취득해야 합니다.
3. 개인 및 비영리 단체는 AGPL-3.0 조건에 따라 이 소프트웨어를 자유롭게 사용할 수 있습니다.
4. 상업용 라이선스 문의: rainy-autumn@outlook.com

> 참고: 라이선스의 법적 효력은 원문(영문/중문)을 기준으로 합니다.
