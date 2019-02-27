package generators

import (
	"github.com/iancoleman/strcase"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
	"k8s.io/klog"
	"path"
)

func Packages(_ *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		klog.Fatalf("Failed loading boilerplate: %v", err)
	}

	customArgs := arguments.CustomArgs.(*CustomArgs)

	return generator.Packages{
		&generator.DefaultPackage{
			PackageName: path.Base(arguments.OutputPackagePath),
			PackagePath: arguments.OutputPackagePath,
			HeaderText:  boilerplate,
			PackageDocumentation: []byte(
				`// auto-generated gomega matchers.
`),
			GeneratorFunc: func(ctx *generator.Context) (generators []generator.Generator) {
				it := generator.NewImportTracker()
				for _, t := range ctx.Order {
					if t.Kind == types.Struct {
						fileName := strcase.ToSnake(t.Name.Name)
						if customArgs.ShouldGenerate(t) {
							generators = append(generators, NewStructMatcherGenerator(fileName, t, it))
						}
					}
				}

				return generators
			},
		},
	}
}
