package gen

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gen/internal/generate"
	tmpl "gorm.io/gen/internal/template"
	"gorm.io/gen/internal/utils/pools"
)

func (g *Generator) getModelProtoOutputPath() (outPath string, err error) {
	if strings.Contains(g.ModelProtoPkgPath, string(os.PathSeparator)) {
		outPath, err = filepath.Abs(g.ModelProtoPkgPath)
		if err != nil {
			return "", fmt.Errorf("cannot parse model pkg path: %w", err)
		}
	} else {
		outPath = filepath.Join(filepath.Dir(g.OutPath), g.ModelProtoPkgPath)
	}
	return outPath + string(os.PathSeparator), nil
}

// generateModelProtoFile generate model structures and save to file
func (g *Generator) generateModelProtoFile() error {
	if len(g.models) == 0 {
		return nil
	}

	modelProtoOutPath, err := g.getModelProtoOutputPath()
	if err != nil {
		return err
	}

	if err = os.MkdirAll(modelProtoOutPath, os.ModePerm); err != nil {
		return fmt.Errorf("create model proto pkg path(%s) fail: %s", modelProtoOutPath, err)
	}

	errChan := make(chan error)
	pool := pools.NewPool(concurrent)
	for _, data := range g.models {
		if data == nil || !data.Generated {
			continue
		}
		pool.Wait()
		go func(data *generate.QueryStructMeta) {
			defer pool.Done()

			var buf bytes.Buffer
			err := render(tmpl.ModelProto, &buf, data)
			if err != nil {
				errChan <- err
				return
			}

			modelFile := modelProtoOutPath + data.FileName + ".gen.proto"
			err = g.protoOutput(modelFile, buf.Bytes())
			if err != nil {
				errChan <- err
				return
			}

			g.info(fmt.Sprintf("generate model proto file(table <%s> -> {%s.%s}): %s", data.TableName, data.StructInfo.Package, data.StructInfo.Type, modelFile))
		}(data)
	}
	select {
	case err = <-errChan:
		return err
	case <-pool.AsyncWaitAll():
		g.fillModelPkgPath(modelProtoOutPath)
	}
	return nil
}

// output format and output
func (g *Generator) protoOutput(fileName string, content []byte) error {
	return ioutil.WriteFile(fileName, content, 0640)
}
