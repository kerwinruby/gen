package gen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gen/internal/generate"
	tmpl "gorm.io/gen/internal/template"
	"gorm.io/gen/internal/utils/pools"
)

func (g *Generator) getModelServiceOutputPath() (outPath string, err error) {
	if strings.Contains(g.ServicePkgPath, string(os.PathSeparator)) {
		outPath, err = filepath.Abs(g.ServicePkgPath)
		if err != nil {
			return "", fmt.Errorf("cannot parse service pkg path: %w", err)
		}
	} else {
		outPath = filepath.Join(filepath.Dir(g.OutPath), g.ServicePkgPath)
	}
	return outPath + string(os.PathSeparator), nil
}

// generateModelProtoFile generate model structures and save to file
func (g *Generator) generateModelServiceFile() error {
	if len(g.models) == 0 {
		return nil
	}

	serviceOutPath, err := g.getModelServiceOutputPath()
	if err != nil {
		return err
	}

	if err = os.MkdirAll(serviceOutPath, os.ModePerm); err != nil {
		return fmt.Errorf("create model service pkg path(%s) fail: %s", serviceOutPath, err)
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
			err := render(tmpl.ModelService, &buf, data)
			if err != nil {
				errChan <- err
				return
			}

			modelFile := serviceOutPath + "service_" + data.FileName + ".gen.crud.go"
			err = g.protoOutput(modelFile, buf.Bytes())
			if err != nil {
				errChan <- err
				return
			}

			g.info(fmt.Sprintf("generate model service file(table <%s> -> {%s.%s}): %s", data.TableName, data.StructInfo.Package, data.StructInfo.Type, modelFile))
		}(data)
	}
	select {
	case err = <-errChan:
		return err
	case <-pool.AsyncWaitAll():
		g.fillModelPkgPath(serviceOutPath)
	}
	return nil
}
