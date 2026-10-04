package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/app"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/asset"
	assetCommon "github.com/Autumn-27/ScopeSentry/internal/services/assets/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/crawler"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/dirscan"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/ip"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/mp"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/page_monitoring"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/root_domain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/sensitive"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/subdomain"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/url"
	"github.com/Autumn-27/ScopeSentry/internal/services/assets/vulnerability"
	"github.com/Autumn-27/ScopeSentry/internal/services/node"
	"github.com/Autumn-27/ScopeSentry/internal/services/plugin"
	"github.com/Autumn-27/ScopeSentry/internal/services/project"
	taskCommon "github.com/Autumn-27/ScopeSentry/internal/services/task/common"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/task"
	"github.com/Autumn-27/ScopeSentry/internal/services/task/template"
	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type deps struct {
	projectService    project.Service
	taskService       task.Service
	taskCommonService taskCommon.Service
	templateService   template.Service
	assetService      asset.Service
	rootDomainService root_domain.Service
	subdomainService  subdomain.Service
	appService        app.Service
	mpService         mp.Service
	urlService        url.Service
	sensitiveService  sensitive.Service
	dirscanService    dirscan.Service
	crawlerService    crawler.Service
	vulnService       vulnerability.Service
	ipService         ip.Service
	pageMonService    page_monitoring.Service
	commonService     assetCommon.Service
	nodeService       node.Service
	pluginService     plugin.Service
}

var d = &deps{
	projectService:    project.NewService(),
	taskService:       task.NewService(),
	taskCommonService: taskCommon.NewService(),
	templateService:   template.NewService(),
	assetService:      asset.NewService(),
	rootDomainService: root_domain.NewService(),
	subdomainService:  subdomain.NewService(),
	appService:        app.NewService(),
	mpService:         mp.NewService(),
	urlService:        url.NewService(),
	sensitiveService:  sensitive.NewService(),
	dirscanService:    dirscan.NewService(),
	crawlerService:    crawler.NewService(),
	vulnService:       vulnerability.NewService(),
	ipService:         ip.NewService(),
	pageMonService:    page_monitoring.NewService(),
	commonService:     assetCommon.NewService(),
	nodeService:       node.NewService(),
	pluginService:     plugin.NewService(),
}

func registerTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects",
		Description: "태그별로 그룹화된 프로젝트 목록 조회",
	}, listProjects)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_projects_data",
		Description: "프로젝트 목록을 페이지네이션으로 조회 (검색 지원)",
	}, listProjectsData)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_project",
		Description: "프로젝트 ID로 프로젝트 상세 정보 조회",
	}, getProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_project",
		Description: "새 프로젝트 생성. name과 target은 필수, tag, template, node 등은 선택",
	}, createProject)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tasks",
		Description: "스캔 작업 목록을 페이지네이션으로 조회",
	}, listTasks)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_task",
		Description: "작업 ID로 작업 상세 정보 조회",
	}, getTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_scan_templates",
		Description: "스캔 템플릿 목록을 페이지네이션으로 조회",
	}, listScanTemplates)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_scan_template",
		Description: "템플릿 ID로 스캔 템플릿 상세 정보 조회",
	}, getScanTemplate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_plugin_modules",
		Description: "스캔 템플릿의 전체 모듈명(파이프라인 단계) 조회. 스캔 템플릿은 이 모듈들로 구성되며, 각 모듈에는 여러 플러그인이 연결됩니다.",
	}, listPluginModules)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_plugins",
		Description: "사용 가능한 스캔 플러그인 조회, module로 필터링 가능. 각 플러그인의 hash, name, module, 기본 parameter를 반환합니다. 스캔 템플릿 생성 시 모듈 필드는 플러그인 hash를 참조합니다.",
	}, listPlugins)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_scan_template",
		Description: createScanTemplateToolDesc,
	}, createScanTemplate)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_scan_task",
		Description: createScanTaskToolDesc,
	}, createScanTask)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_assets",
		Description: listAssetsToolDesc,
	}, listAssets)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "count_assets",
		Description: countAssetsToolDesc,
	}, countAssets)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_asset_detail",
		Description: "자산 상세 정보 조회, asset_type은 asset 또는 vulnerability 지원",
	}, getAssetDetail)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "add_asset_tag",
		Description: "자산에 태그 추가",
	}, addAssetTag)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_nodes",
		Description: "스캔 노드 목록 조회, online_only=true이면 온라인 노드만 반환",
	}, listNodes)
}

type listProjectsDataInput struct {
	Search    string `json:"search,omitempty" jsonschema:"프로젝트 이름 퍼지 검색 키워드"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"페이지 번호, 1부터 시작, 기본값 1"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"페이지당 항목 수, 기본값 20"`
}

func listProjects(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectsByTag(c)
	if err != nil {
		return errorResult("프로젝트 목록 조회 실패", err)
	}
	return jsonToolResult(ginH{"list": result})
}

func listProjectsData(ctx context.Context, _ *mcp.CallToolRequest, input listProjectsDataInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectsData(c, input.Search, input.PageIndex, input.PageSize)
	if err != nil {
		return errorResult("프로젝트 데이터 조회 실패", err)
	}
	return jsonToolResult(result)
}

type getProjectInput struct {
	ID string `json:"id" jsonschema:"프로젝트 MongoDB ObjectID"`
}

func getProject(ctx context.Context, _ *mcp.CallToolRequest, input getProjectInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id는 비워 둘 수 없습니다", nil)
	}
	c := ginContext(ctx)
	result, err := d.projectService.GetProjectContent(c, input.ID)
	if err != nil {
		return errorResult("프로젝트 상세 정보 조회 실패", err)
	}
	if result == nil {
		return errorResult("프로젝트가 존재하지 않습니다", nil)
	}
	return jsonToolResult(result)
}

type createProjectInput struct {
	Name           string   `json:"name" jsonschema:"프로젝트 이름, 필수"`
	Tag            string   `json:"tag,omitempty" jsonschema:"프로젝트 태그, 그룹화에 사용"`
	Target         string   `json:"target" jsonschema:"스캔 대상, 필수. 여러 줄 또는 쉼표로 구분, 도메인/IP/URL 등 지원"`
	Template       string   `json:"template,omitempty" jsonschema:"연결할 스캔 템플릿 ID"`
	Node           []string `json:"node,omitempty" jsonschema:"지정할 스캔 노드 이름 목록"`
	AllNode        bool     `json:"allNode,omitempty" jsonschema:"모든 노드 사용 여부"`
	Ignore         string   `json:"ignore,omitempty" jsonschema:"무시할 대상 목록, 형식은 target과 동일"`
	Duplicates     string   `json:"duplicates,omitempty" jsonschema:"중복 제거 전략"`
	ScheduledTasks bool     `json:"scheduledTasks,omitempty" jsonschema:"예약 스캔 활성화 여부"`
	Hour           int      `json:"hour,omitempty" jsonschema:"예약 스캔 간격(시간), scheduledTasks 활성화 시에만 유효"`
}

func createProject(ctx context.Context, _ *mcp.CallToolRequest, input createProjectInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" || input.Target == "" {
		return errorResult("name과 target은 비워 둘 수 없습니다", nil)
	}
	p := &models.Project{
		Name:           input.Name,
		Tag:            input.Tag,
		Target:         input.Target,
		Template:       input.Template,
		Node:           input.Node,
		AllNode:        input.AllNode,
		Ignore:         input.Ignore,
		Duplicates:     input.Duplicates,
		ScheduledTasks: input.ScheduledTasks,
		Hour:           input.Hour,
		Tp:             "project",
	}
	c := ginContext(ctx)
	if err := d.projectService.AddProject(c, p); err != nil {
		return errorResult("프로젝트 생성 실패", err)
	}
	return jsonToolResult(ginH{"success": true, "message": "프로젝트 생성 성공"})
}

type listTasksInput struct {
	Search    string `json:"search,omitempty" jsonschema:"작업 이름 퍼지 검색"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"페이지 번호, 1부터 시작, 기본값 1"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"페이지당 항목 수, 기본값 20"`
}

func listTasks(ctx context.Context, _ *mcp.CallToolRequest, input listTasksInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	c := ginContext(ctx)
	tasks, total, err := d.taskService.List(c, input.Search, input.PageIndex, input.PageSize)
	if err != nil {
		return errorResult("작업 목록 조회 실패", err)
	}
	return jsonToolResult(ginH{"list": tasks, "total": total})
}

type getTaskInput struct {
	ID string `json:"id" jsonschema:"작업 MongoDB ObjectID"`
}

func getTask(ctx context.Context, _ *mcp.CallToolRequest, input getTaskInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id는 비워 둘 수 없습니다", nil)
	}
	c := ginContext(ctx)
	result, err := d.taskService.GetTaskDetail(c, input.ID)
	if err != nil {
		return errorResult("작업 상세 정보 조회 실패", err)
	}
	if result == nil {
		return errorResult("작업이 존재하지 않습니다", nil)
	}
	return jsonToolResult(result)
}

type listScanTemplatesInput struct {
	Query     string `json:"query,omitempty" jsonschema:"템플릿 이름 퍼지 검색"`
	PageIndex int    `json:"pageIndex,omitempty" jsonschema:"페이지 번호, 1부터 시작, 기본값 1"`
	PageSize  int    `json:"pageSize,omitempty" jsonschema:"페이지당 항목 수, 기본값 20"`
}

func listScanTemplates(ctx context.Context, _ *mcp.CallToolRequest, input listScanTemplatesInput) (*mcp.CallToolResult, any, error) {
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}
	result, err := d.templateService.List(ctx, input.PageIndex, input.PageSize, input.Query)
	if err != nil {
		return errorResult("템플릿 목록 조회 실패", err)
	}
	return jsonToolResult(result)
}

type getScanTemplateInput struct {
	ID string `json:"id" jsonschema:"스캔 템플릿 MongoDB ObjectID"`
}

func getScanTemplate(ctx context.Context, _ *mcp.CallToolRequest, input getScanTemplateInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id는 비워 둘 수 없습니다", nil)
	}
	c := ginContext(ctx)
	result, err := d.templateService.Detail(c, input.ID)
	if err != nil {
		return errorResult("템플릿 상세 정보 조회 실패", err)
	}
	return jsonToolResult(result)
}

type createScanTemplateInput struct {
	Name         string                       `json:"name" jsonschema:"템플릿 이름, 필수"`
	Modules      map[string][]string          `json:"modules,omitempty" jsonschema:"모듈에서 플러그인 hash 목록으로의 매핑. key는 모듈명(list_plugin_modules로 조회), value는 해당 모듈에서 활성화할 플러그인 hash 배열(list_plugins로 조회, 배열 순서대로 실행). 예: {\"SubdomainScan\":[\"d60ba73c...\"]}"`
	Parameters   map[string]map[string]string `json:"parameters,omitempty" jsonschema:"선택, 플러그인 실행 파라미터 재정의. 구조는 모듈명->플러그인 hash->파라미터 문자열. 제공하지 않으면 플러그인 기본 파라미터를 자동 사용"`
	VulList      []string                     `json:"vullist,omitempty" jsonschema:"선택, nuclei POC 템플릿 ID 목록, VulnerabilityScan에서 nuclei를 사용할 때만 유효"`
	TemplateJSON string                       `json:"template_json,omitempty" jsonschema:"선택, 전체 ScanTemplate JSON. 제공 시 modules보다 우선하며, 고급 사용자 정의에 사용"`
}

func createScanTemplate(ctx context.Context, _ *mcp.CallToolRequest, input createScanTemplateInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" && input.TemplateJSON == "" {
		return errorResult("name 또는 template_json 중 최소 하나는 제공해야 합니다", nil)
	}

	var tmpl *models.ScanTemplate

	switch {
	case input.TemplateJSON != "":
		tmpl = &models.ScanTemplate{}
		if err := json.Unmarshal([]byte(input.TemplateJSON), tmpl); err != nil {
			return errorResult("template_json 형식이 유효하지 않습니다", err)
		}
		if tmpl.Name == "" {
			tmpl.Name = input.Name
		}
	case len(input.Modules) > 0:
		built, err := buildTemplateFromModules(ctx, input)
		if err != nil {
			return errorResult(err.Error(), nil)
		}
		tmpl = built
	default:
		tmpl = &models.ScanTemplate{Name: input.Name}
	}

	if tmpl.Name == "" {
		return errorResult("템플릿 이름은 비워 둘 수 없습니다", nil)
	}

	id, err := d.templateService.Save(ctx, "", tmpl)
	if err != nil {
		return errorResult("템플릿 생성 실패", err)
	}
	return jsonToolResult(ginH{"success": true, "id": id, "message": "템플릿 생성 성공, create_scan_task의 template 파라미터에 사용할 수 있습니다"})
}

// buildTemplateFromModules "모듈->플러그인 hash 목록"으로 스캔 템플릿을 구성하고 플러그인 기본 파라미터를 자동으로 채움
func buildTemplateFromModules(ctx context.Context, input createScanTemplateInput) (*models.ScanTemplate, error) {
	validModule := make(map[string]bool, len(constants.PLUGINSMODULES))
	for _, m := range constants.PLUGINSMODULES {
		validModule[m] = true
	}

	templateMap := map[string]any{
		"name": input.Name,
	}
	paramMap := map[string]map[string]string{}

	for module, hashes := range input.Modules {
		if !validModule[module] {
			return nil, fmt.Errorf("유효하지 않은 모듈명: %s (list_plugin_modules로 유효한 모듈 조회)", module)
		}
		if len(hashes) == 0 {
			continue
		}
		templateMap[module] = hashes

		modParams := map[string]string{}
		for _, hash := range hashes {
			// 호출자가 제공한 파라미터 재정의를 우선 사용
			if input.Parameters != nil {
				if mp, ok := input.Parameters[module]; ok {
					if v, ok := mp[hash]; ok {
						modParams[hash] = v
						continue
					}
				}
			}
			// 그렇지 않으면 플러그인 기본 파라미터로 채움
			plg, err := d.pluginService.GetPluginByHash(ctx, hash)
			if err != nil || plg == nil {
				return nil, fmt.Errorf("플러그인 hash가 존재하지 않습니다: %s (모듈 %s)", hash, module)
			}
			if plg.Module != module {
				return nil, fmt.Errorf("플러그인 %s(hash=%s)는 모듈 %s에 속하므로 모듈 %s에 넣을 수 없습니다", plg.Name, hash, plg.Module, module)
			}
			modParams[hash] = plg.Parameter
		}
		if len(modParams) > 0 {
			paramMap[module] = modParams
		}
	}

	if len(paramMap) > 0 {
		templateMap["Parameters"] = paramMap
	}
	if len(input.VulList) > 0 {
		templateMap["vullist"] = input.VulList
	}

	data, err := json.Marshal(templateMap)
	if err != nil {
		return nil, fmt.Errorf("템플릿 구성 실패: %w", err)
	}
	tmpl := &models.ScanTemplate{}
	if err := json.Unmarshal(data, tmpl); err != nil {
		return nil, fmt.Errorf("템플릿 구성 실패: %w", err)
	}
	return tmpl, nil
}

type createScanTaskInput struct {
	Name           string              `json:"name" jsonschema:"작업 이름, 필수, 중복 불가"`
	Target         string              `json:"target,omitempty" jsonschema:"스캔 대상, targetSource가 general일 때 필수, 여러 줄 또는 쉼표로 구분"`
	Node           []string            `json:"node" jsonschema:"스캔을 실행할 노드 이름 목록, 필수"`
	Template       string              `json:"template,omitempty" jsonschema:"스캔 템플릿 ObjectID, 스캔 실행에 필수"`
	AllNode        bool                `json:"allNode,omitempty" jsonschema:"모든 온라인 노드 자동 추가 여부"`
	Ignore         string              `json:"ignore,omitempty" jsonschema:"무시할 대상 목록, 형식은 target과 동일"`
	Duplicates     string              `json:"duplicates,omitempty" jsonschema:"중복 제거 전략, 예: None"`
	Project        []string            `json:"project,omitempty" jsonschema:"연결할 프로젝트 ObjectID 목록, targetSource가 project일 때 필수"`
	TargetSource   string              `json:"targetSource,omitempty" jsonschema:"대상 출처, general project asset RootDomain subdomain UrlScan 및 해당 Source 접미사 중 선택, 기본값 general, 자세한 내용은 도구 description 참조"`
	TargetTp       string              `json:"targetTp,omitempty" jsonschema:"Source 출처일 때의 선택 방식, search 또는 select"`
	Search         string              `json:"search,omitempty" jsonschema:"자산 저장소에서 대상을 필터링하는 search 표현식, 문법은 list_assets와 동일, 자세한 내용은 도구 description 참조"`
	Filter         map[string][]string `json:"filter,omitempty" jsonschema:"정확 필터, search와 조합 가능, filter.project는 프로젝트 ObjectID"`
	TargetNumber   int                 `json:"targetNumber,omitempty" jsonschema:"search 모드에서 대상 수 상한, 0은 제한 없음"`
	TargetIds      []string            `json:"targetIds,omitempty" jsonschema:"select 모드에서 선택된 자산 ObjectID 목록"`
	BindProject    string              `json:"bindProject,omitempty" jsonschema:"바인딩할 프로젝트 ObjectID, 스캔 결과 귀속 대상"`
	ScheduledTasks bool                `json:"scheduledTasks,omitempty" jsonschema:"예약 작업으로 생성할지 여부"`
	Hour           int                 `json:"hour,omitempty" jsonschema:"예약 작업 간격-시간"`
	Minute         int                 `json:"minute,omitempty" jsonschema:"예약 작업 간격-분"`
	Day            int                 `json:"day,omitempty" jsonschema:"예약 작업 간격-일"`
	Week           int                 `json:"week,omitempty" jsonschema:"예약 작업 간격-주 (weekly 주기)"`
	CycleType      string              `json:"cycleType,omitempty" jsonschema:"주기 유형, 예: nhours daily weekly"`
}

func createScanTask(ctx context.Context, _ *mcp.CallToolRequest, input createScanTaskInput) (*mcp.CallToolResult, any, error) {
	if input.Name == "" || len(input.Node) == 0 {
		return errorResult("name과 node는 비워 둘 수 없습니다", nil)
	}
	c := ginContext(ctx)
	exists, err := d.taskService.CheckTaskNameExists(c, input.Name)
	if err != nil {
		return errorResult("작업 이름 확인 실패", err)
	}
	if exists {
		return errorResult("작업 이름이 이미 존재합니다", nil)
	}

	targetSource := input.TargetSource
	if targetSource == "" {
		targetSource = "general"
	}

	taskModel := &models.Task{
		Name:           input.Name,
		Target:         input.Target,
		Node:           input.Node,
		Template:       input.Template,
		AllNode:        input.AllNode,
		Ignore:         input.Ignore,
		Duplicates:     input.Duplicates,
		Project:        input.Project,
		TargetSource:   targetSource,
		TargetTp:       input.TargetTp,
		Search:         input.Search,
		TargetNumber:   input.TargetNumber,
		TargetIds:      input.TargetIds,
		BindProject:    input.BindProject,
		ScheduledTasks: input.ScheduledTasks,
		Hour:           input.Hour,
		Minute:         input.Minute,
		Day:            input.Day,
		Week:           input.Week,
		CycleType:      input.CycleType,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		taskModel.Filter = filter
	}
	taskID, err := d.taskCommonService.Insert(ctx, taskModel)
	if err != nil {
		return errorResult("스캔 작업 생성 실패", err)
	}
	return jsonToolResult(ginH{"success": true, "id": taskID})
}

type listAssetsInput struct {
	AssetType        string              `json:"asset_type" jsonschema:"자산 유형, 필수. 예: asset, RootDomain, subdomain, app, mp, UrlScan, SensitiveResult, DirScanResult, crawler, vulnerability, PageMonitoring, IPAsset, SubdomainTakerResult"`
	PageIndex        int                 `json:"pageIndex,omitempty" jsonschema:"페이지 번호, 1부터 시작, 기본값 1"`
	PageSize         int                 `json:"pageSize,omitempty" jsonschema:"페이지당 항목 수, 기본값 20"`
	SearchExpression string              `json:"search,omitempty" jsonschema:"검색 표현식, SQL 아님, 문법은 도구 description 참조"`
	Filter           map[string][]string `json:"filter,omitempty" jsonschema:"정확 필터 JSON, search와 조합 가능. filter.project는 프로젝트 ObjectID(list_projects로 조회, 프로젝트 이름 아님); filter.task는 작업 이름. 자세한 내용은 도구 description 참조"`
	Sort             map[string]string   `json:"sort,omitempty" jsonschema:"정렬. UrlScan/DirScanResult만 length 지원: ascending은 오름차순, 그 외 값은 내림차순"`
	Sid              string              `json:"sid,omitempty" jsonschema:"민감 정보 규칙 이름, asset_type이 SensitiveResult일 때만 유효"`
}

type countAssetsInput struct {
	AssetType        string              `json:"asset_type" jsonschema:"자산 유형, 필수. 값은 list_assets와 동일"`
	SearchExpression string              `json:"search,omitempty" jsonschema:"검색 표현식, 문법은 list_assets와 동일"`
	Filter           map[string][]string `json:"filter,omitempty" jsonschema:"정확 필터 JSON, 문법은 list_assets와 동일; filter.project는 프로젝트 ObjectID"`
	Sid              string              `json:"sid,omitempty" jsonschema:"민감 정보 규칙 이름, asset_type이 SensitiveResult일 때만 유효"`
}

func countAssets(ctx context.Context, _ *mcp.CallToolRequest, input countAssetsInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}

	query := models.SearchRequest{
		Index:            index,
		SearchExpression: input.SearchExpression,
		Sid:              input.Sid,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		query.Filter = filter
	}

	total, err := d.commonService.TotalData(ctx, &query)
	if err != nil {
		return errorResult("자산 수 집계 실패", err)
	}
	return jsonToolResult(ginH{"total": total})
}

func listAssets(ctx context.Context, _ *mcp.CallToolRequest, input listAssetsInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}
	if input.PageIndex <= 0 {
		input.PageIndex = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 20
	}

	query := models.SearchRequest{
		PageIndex:        input.PageIndex,
		PageSize:         input.PageSize,
		Index:            index,
		SearchExpression: input.SearchExpression,
		Sort:             input.Sort,
		Sid:              input.Sid,
	}
	if len(input.Filter) > 0 {
		filter := make(map[string][]interface{}, len(input.Filter))
		for k, vals := range input.Filter {
			items := make([]interface{}, len(vals))
			for i, v := range vals {
				items[i] = v
			}
			filter[k] = items
		}
		query.Filter = filter
	}

	c := ginContext(ctx)
	data, err := queryAssets(ctx, c, index, query)
	if err != nil {
		return errorResult("자산 조회 실패", err)
	}
	return jsonToolResult(data)
}

type getAssetDetailInput struct {
	AssetType string `json:"asset_type" jsonschema:"자산 유형. 상세 조회는 asset 또는 vulnerability 지원"`
	ID        string `json:"id" jsonschema:"자산 MongoDB ObjectID; vulnerability 유형은 hash 값 전달"`
}

func getAssetDetail(ctx context.Context, _ *mcp.CallToolRequest, input getAssetDetailInput) (*mcp.CallToolResult, any, error) {
	if input.ID == "" {
		return errorResult("id는 비워 둘 수 없습니다", nil)
	}
	c := ginContext(ctx)
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}

	switch index {
	case "asset":
		result, err := d.assetService.GetAssetByID(c, input.ID)
		if err != nil {
			return errorResult("자산 상세 정보 조회 실패", err)
		}
		if result == nil {
			return errorResult("자산이 존재하지 않습니다", nil)
		}
		return jsonToolResult(result)
	case "vulnerability":
		result, err := d.vulnService.GetVulnerabilityDetailByHash(c, input.ID)
		if err != nil {
			return errorResult("취약점 상세 정보 조회 실패", err)
		}
		return jsonToolResult(result)
	default:
		return errorResult("asset_type은 asset 또는 vulnerability의 상세 조회만 지원합니다", nil)
	}
}

type addAssetTagInput struct {
	AssetType string `json:"asset_type" jsonschema:"자산 유형, list_assets의 asset_type과 동일"`
	ID        string `json:"id" jsonschema:"자산 MongoDB ObjectID"`
	Tag       string `json:"tag" jsonschema:"추가할 태그 이름"`
}

func addAssetTag(ctx context.Context, _ *mcp.CallToolRequest, input addAssetTagInput) (*mcp.CallToolResult, any, error) {
	index, err := normalizeAssetIndex(input.AssetType)
	if err != nil {
		return errorResult(err.Error(), nil)
	}
	if input.ID == "" || input.Tag == "" {
		return errorResult("id와 tag는 비워 둘 수 없습니다", nil)
	}
	c := ginContext(ctx)
	req := &models.TagRequest{Type: index, ID: input.ID, Tag: input.Tag}
	if err := d.commonService.AddTag(c, req); err != nil {
		return errorResult("태그 추가 실패", err)
	}
	return jsonToolResult(ginH{"success": true})
}

type listNodesInput struct {
	OnlineOnly bool `json:"online_only,omitempty" jsonschema:"true이면 온라인 노드만 반환, 기본값 false는 전체 반환"`
}

func listNodes(ctx context.Context, _ *mcp.CallToolRequest, input listNodesInput) (*mcp.CallToolResult, any, error) {
	result, err := d.nodeService.GetNodeData(ctx, input.OnlineOnly)
	if err != nil {
		return errorResult("노드 목록 조회 실패", err)
	}
	return jsonToolResult(ginH{"list": result})
}

func listPluginModules(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	return jsonToolResult(ginH{"modules": constants.PLUGINSMODULES})
}

type listPluginsInput struct {
	Module string `json:"module,omitempty" jsonschema:"모듈명으로 필터링, 비워 두면 전체 스캔 플러그인 반환. 모듈명은 list_plugin_modules로 조회"`
	Search string `json:"search,omitempty" jsonschema:"플러그인 이름 퍼지 검색 (module 미지정 시에만 적용)"`
}

type pluginBrief struct {
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	Module    string `json:"module"`
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
	Status    bool   `json:"status"`
}

func toPluginBriefs(plugins []models.Plugin) []pluginBrief {
	briefs := make([]pluginBrief, 0, len(plugins))
	for _, p := range plugins {
		// 서버 측 플러그인은 스캔 파이프라인에 참여하지 않으므로 건너뜀
		if p.Type == "server" {
			continue
		}
		briefs = append(briefs, pluginBrief{
			Hash:      p.Hash,
			Name:      p.Name,
			Module:    p.Module,
			Parameter: p.Parameter,
			Type:      p.Type,
			Status:    p.Status,
		})
	}
	return briefs
}

func listPlugins(ctx context.Context, _ *mcp.CallToolRequest, input listPluginsInput) (*mcp.CallToolResult, any, error) {
	c := ginContext(ctx)

	if input.Module != "" {
		plugins, err := d.pluginService.ListByModule(c, input.Module)
		if err != nil {
			return errorResult("플러그인 목록 조회 실패", err)
		}
		return jsonToolResult(ginH{"list": toPluginBriefs(plugins)})
	}

	resp, err := d.pluginService.List(c, &models.PluginListRequest{
		PageIndex: 1,
		PageSize:  500,
		Search:    input.Search,
	})
	if err != nil {
		return errorResult("플러그인 목록 조회 실패", err)
	}
	return jsonToolResult(ginH{"list": toPluginBriefs(resp.List), "total": resp.Total})
}

func queryAssets(ctx context.Context, c *gin.Context, index string, query models.SearchRequest) (any, error) {
	switch index {
	case "asset":
		list, err := d.assetService.GetAssets(c, query)
		return ginH{"list": list}, err
	case "RootDomain":
		return d.rootDomainService.GetRootDomainData(c, query)
	case "subdomain":
		list, err := d.subdomainService.GetSubdomains(ctx, query)
		return ginH{"list": list}, err
	case "app":
		return d.appService.GetAppData(c, query)
	case "mp":
		return d.mpService.GetMPData(c, query)
	case "UrlScan":
		list, err := d.urlService.GetURLs(ctx, query)
		return ginH{"list": list}, err
	case "SensitiveResult":
		list, err := d.sensitiveService.GetSensitiveInfo(ctx, query)
		return ginH{"list": list}, err
	case "DirScanResult":
		list, err := d.dirscanService.List(ctx, query)
		return ginH{"list": list}, err
	case "crawler":
		list, err := d.crawlerService.GetCrawlers(ctx, query)
		return ginH{"list": list}, err
	case "vulnerability":
		return d.vulnService.GetVulnerabilities(ctx, query)
	case "PageMonitoring":
		return d.pageMonService.GetResult(c, query)
	case "IPAsset":
		return d.ipService.GetIPAssets(c, query)
	case "SubdomainTakerResult":
		return querySubdomainTaker(ctx, query)
	default:
		return nil, fmt.Errorf("지원되지 않는 자산 유형: %s", index)
	}
}

func querySubdomainTaker(ctx context.Context, query models.SearchRequest) (any, error) {
	takerService := subdomain.NewTakerService()
	c := ginContext(ctx)
	return takerService.GetSubdomainTakerData(c, query)
}

func normalizeAssetIndex(assetType string) (string, error) {
	assetType = strings.TrimSpace(assetType)
	if assetType == "" {
		return "", fmt.Errorf("asset_type은 비워 둘 수 없습니다")
	}
	aliases := map[string]string{
		"asset":                "asset",
		"web":                  "asset",
		"rootdomain":           "RootDomain",
		"root_domain":          "RootDomain",
		"root-domain":          "RootDomain",
		"subdomain":            "subdomain",
		"app":                  "app",
		"mp":                   "mp",
		"miniprogram":          "mp",
		"mini_program":         "mp",
		"url":                  "UrlScan",
		"urlscan":              "UrlScan",
		"sensitive":            "SensitiveResult",
		"sensitiveresult":      "SensitiveResult",
		"sensitive_result":     "SensitiveResult",
		"dirscan":              "DirScanResult",
		"dir_scan":             "DirScanResult",
		"dirscanresult":        "DirScanResult",
		"directory":            "DirScanResult",
		"crawler":              "crawler",
		"vulnerability":        "vulnerability",
		"vuln":                 "vulnerability",
		"pagemonitoring":       "PageMonitoring",
		"page_monitoring":      "PageMonitoring",
		"ip":                   "IPAsset",
		"ipasset":              "IPAsset",
		"subdomaintaker":       "SubdomainTakerResult",
		"subdomain_taker":      "SubdomainTakerResult",
		"subdomaintakerresult": "SubdomainTakerResult",
	}
	key := strings.ToLower(assetType)
	if v, ok := aliases[key]; ok {
		return v, nil
	}
	if _, ok := aliases[strings.ReplaceAll(key, "-", "_")]; ok {
		return aliases[strings.ReplaceAll(key, "-", "_")], nil
	}
	// MongoDB 컬렉션 이름을 직접 전달하는 것도 허용
	valid := []string{"asset", "RootDomain", "subdomain", "app", "mp", "UrlScan",
		"SensitiveResult", "DirScanResult", "crawler", "vulnerability",
		"PageMonitoring", "IPAsset", "SubdomainTakerResult"}
	for _, v := range valid {
		if v == assetType {
			return v, nil
		}
	}
	return "", fmt.Errorf("지원되지 않는 asset_type: %s", assetType)
}

type ginH map[string]any
