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

func (m *{{.ModelStructName}}) BeforeCreate(ctx context.Context) *{{.ModelStructName}} {
	{{if ExistsField "ID" .Fields}}
	m.ID = bson.NewObjectId().Hex()
	{{end}}	
	{{if ExistsField "CreatedTime" .Fields}}
	m.CreatedTime = xtime.Millisecond()
	{{if or (ExistsField "CreatedID" .Fields) (ExistsField "CreatedName" .Fields)}}
	if c, ok := ctx.(*bm.Context); ok {
		{{if ExistsField "CreatedID" .Fields}}
		if id, ok := c.Keys[metadata.UserID]; ok {
			m.CreatedID = id.(string)	
		}
		{{end}}
		{{if ExistsField "CreatedName" .Fields}}
		if name, ok := c.Keys[metadata.UserName]; ok {
			m.CreatedName = name.(string)
		}
		{{end}}
	}
	{{end}}
	{{end}}
	return m
}

func (m *{{.ModelStructName}}) BeforeUpdate(ctx context.Context) *{{.ModelStructName}} {
	{{if ExistsField "UpdatedTime" .Fields}}
	m.UpdatedTime = xtime.Millisecond()
	{{if or (ExistsField "UpdatedID" .Fields) (ExistsField "UpdatedName" .Fields)}}
	if c, ok := ctx.(*bm.Context); ok {
		{{if ExistsField "UpdatedID" .Fields}}
		if id, ok := c.Keys[metadata.UserID]; ok {
			m.UpdatedID = id.(string)	
		}
		{{end}}
		{{if ExistsField "UpdatedName" .Fields}}
		if name, ok := c.Keys[metadata.UserName]; ok {
			m.UpdatedName = name.(string)
		}
		{{end}}
	}
	{{end}}		
	{{end}}
	return m
}

func (m *{{.ModelStructName}}) BeforeDelete(ctx context.Context) *{{.ModelStructName}} {
	{{if ExistsField "DeletedTime" .Fields}}
	m.DeletedTime = xtime.Millisecond()
	{{if or (ExistsField "DeletedID" .Fields) (ExistsField "DeletedName" .Fields)}}
	if c, ok := ctx.(*bm.Context); ok {
		{{if ExistsField "DeletedID" .Fields}}
		if id, ok := c.Keys[metadata.UserID]; ok {
			m.DeletedID = id.(string)	
		}
		{{end}}
		{{if ExistsField "DeletedName" .Fields}}
		if name, ok := c.Keys[metadata.UserName]; ok {
			m.DeletedName = name.(string)
		}
		{{end}}
	}
	{{end}}
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
