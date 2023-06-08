package template

// Model used as a variable because it cannot load template file after packed, params still can pass file
const Model = NotEditMark + `
package {{.StructInfo.Package}}

import (
	"context"
	"encoding/json"
	"time"

	xtime "gitlab.datahunter.cn/common/kratos/pkg/time"
	bm "gitlab.datahunter.cn/common/kratos/pkg/net/http/blademaster"
	"gitlab.datahunter.cn/common/kratos/pkg/net/metadata"

	pb "{{.ProjectName}}/api"
	"github.com/jinzhu/copier"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gopkg.in/mgo.v2/bson"
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

{{if .TableName -}}const TableName{{.ModelStructName}} = "{{.TableName}}"{{- end}}

// {{.ModelStructName}} {{.StructComment}}
type {{.ModelStructName}} struct {
    {{range .Fields}}
    {{if .MultilineComment -}}
	/*
{{.ColumnComment}}
    */
	{{end -}}
    {{.Name}} {{.Type}} ` + "`{{.Tags}}` " +
	"{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}" +
	`{{end}}
}

func (m *{{.ModelStructName}}) BeforeCreate_(ctx context.Context) *{{.ModelStructName}} {
	{{if ExistsField "ID" .Fields}}
	m.ID = bson.NewObjectId().Hex()
	{{end}}	
	{{if ExistsField "CreatedTime" .Fields}}
	m.CreatedTime = xtime.Millisecond()
	{{end}}	
	{{if ExistsField "CreatedID" .Fields}}
	m.CreatedID = metadata.String(ctx, metadata.UserID)
	{{end}}
	{{if ExistsField "CreatedName" .Fields}}
	m.CreatedName = metadata.String(ctx, metadata.UserName)
	{{end}}
	return m
}

func (m *{{.ModelStructName}}) BeforeUpdate_(ctx context.Context) *{{.ModelStructName}} {
	{{if ExistsField "UpdatedTime" .Fields}}
	m.UpdatedTime = xtime.Millisecond()
	{{end}}
	{{if ExistsField "UpdatedID" .Fields}}
	m.UpdatedID = metadata.String(ctx, metadata.UserID)
	{{end}}
	{{if ExistsField "UpdatedName" .Fields}}
	m.UpdatedName = metadata.String(ctx, metadata.UserName)
	{{end}}
	return m
}

func (m *{{.ModelStructName}}) BeforeDelete_(ctx context.Context) *{{.ModelStructName}} {
	{{if ExistsField "DeletedTime" .Fields}}
	{{if .FieldSoftDelete}}
	m.DeletedTime = field_type.DeletedTime(xtime.Millisecond())
	{{else}}
	m.DeletedTime = xtime.Millisecond()
	{{end}}
	{{end}}
	{{if ExistsField "DeletedID" .Fields}}
	m.DeletedID = metadata.String(ctx, metadata.UserID)
	{{end}}
	{{if ExistsField "DeletedName" .Fields}}
	m.DeletedName = metadata.String(ctx, metadata.UserName)
	{{end}}
	return m
}

func (m *{{.ModelStructName}}) ToPb() *pb.{{.ModelStructName}} {
	to := &pb.{{.ModelStructName}}{}
	copier.Copy(to, m)
	return to
}

func (m *{{.ModelStructName}}) ToModel(u *pb.{{.ModelStructName}}) *{{.ModelStructName}} {
	copier.Copy(m, u)
	return m
}

type {{.ModelStructName}}s []*{{.ModelStructName}}

func (ms {{.ModelStructName}}s) ToPb() []*pb.{{.ModelStructName}} {
	list := make([]*pb.{{.ModelStructName}}, len(ms))
	for i, m := range ms {
		list[i] = m.ToPb()
	}
	return list
}
`

// ModelMethod model struct DIY method
const ModelMethod = `

{{if .Doc -}}// {{.DocComment -}}{{end}}
func ({{.GetBaseStructTmpl}}){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){{.Body}}
`
