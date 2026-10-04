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
