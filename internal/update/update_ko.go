package update

import (
	"context"
	"regexp"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// 이미 설치된 환경에서도 내장 플러그인 설명을 한국어로 바꿉니다.
// 값에 한자(중국어)가 남아 있는 경우에만 덮어쓰므로, 사용자가 직접 수정한 설명은 유지됩니다.
var hanPattern = regexp.MustCompile(`\p{Han}`)

func UpdateKoPluginText() {
	coll := mongodb.DB.Collection("plugins")
	ctx := context.Background()
	for _, p := range constants.Plugins {
		var doc struct {
			Help         string `bson:"help"`
			Introduction string `bson:"introduction"`
		}
		if err := coll.FindOne(ctx, bson.M{"hash": p.Hash}).Decode(&doc); err != nil {
			continue
		}
		set := bson.M{}
		if hanPattern.MatchString(doc.Help) && !hanPattern.MatchString(p.Help) {
			set["help"] = p.Help
		}
		if hanPattern.MatchString(doc.Introduction) && !hanPattern.MatchString(p.Introduction) {
			set["introduction"] = p.Introduction
		}
		if len(set) > 0 {
			_, _ = coll.UpdateOne(ctx, bson.M{"hash": p.Hash}, bson.M{"$set": set})
		}
	}
}

// 민감 정보 규칙 중 중국어 이름(클라우드 업체 등)을 영어로 바꿉니다.
// 스캔 결과(SensitiveResult)는 규칙 이름을 sid 로 저장하므로, 규칙과 함께 기존 결과의 sid 도 바꿔
// 같은 규칙이 이전/새 이름 두 그룹으로 나뉘지 않게 합니다. 이전 이름의 규칙이 남아 있을 때만 실행됩니다.
var sensitiveRuleRename = map[string]string{
	"腾讯云":   "Tencent Cloud",
	"亚马逊云":  "Amazon Cloud",
	"阿里云":   "Alibaba Cloud",
	"华为云":   "Huawei Cloud",
	"百度云":   "Baidu Cloud",
	"京东云":   "JD Cloud",
	"青云":    "QingCloud",
	"金山云":   "Kingsoft Cloud",
	"联通云":   "China Unicom Cloud",
	"移动云":   "China Mobile Cloud",
	"电信云":   "China Telecom Cloud",
	"一云通":   "Yiyuntong Cloud",
	"用友云":   "Yonyou Cloud",
	"南大通用云": "GBASE Cloud",
	"敏感信息":  "Sensitive Information",
}

func UpdateSensitiveRuleNames() {
	ctx := context.Background()
	rules := mongodb.DB.Collection("SensitiveRule")
	results := mongodb.DB.Collection("SensitiveResult")
	for oldName, newName := range sensitiveRuleRename {
		res, err := rules.UpdateMany(ctx, bson.M{"name": oldName}, bson.M{"$set": bson.M{"name": newName}})
		if err != nil || res.ModifiedCount == 0 {
			continue
		}
		_, _ = results.UpdateMany(ctx, bson.M{"sid": oldName}, bson.M{"$set": bson.M{"sid": newName}})
	}
}
