package generators

import (
	"fmt"
	"io"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

const (
	structMatcherDecl = `
func Match$.|public$(s *$.|matcher$) types.GomegaMatcher {
	return &$.|private${
		$.|matcher$: s,
	}
}
	
type $.|matcher$ struct {
	$range .Members$
		$.Name$ types.GomegaMatcher
	$end$
}

type $.|private$ struct {
	*$.|matcher$

	// State.
	failures []error
}`
	structMatcherMatch = `
func (m *$.|private$) Match(actual interface{}) (success bool, err error) {
	var act *$.|raw$
	kind := reflect.TypeOf(actual).Kind()
	switch kind {
	
	case reflect.Struct:
		a, ok := actual.($.|raw$)
		if !ok {
			return false, fmt.Errorf("%v is type %T, expected $.|public$", actual, actual) 
		}
		act = &a
	
	case reflect.Ptr:
		a, ok := actual.(*$.|raw$)
		if !ok {
			return false, fmt.Errorf("%v is type %T, expected *$.|public$", actual, actual) 
		}
		act = a
	
	default:
		return false, fmt.Errorf("%v is type %T, expected struct or pointer", actual, actual)
	}

	m.failures = m.matchFields(act)
	if len(m.failures) > 0 {
		return false, nil
	}
	return true, nil
}`
	structMatcherMatchFields = `
func (m *$.|private$) matchFields(actual *$.|raw$) (errs []error) {
	var fieldName string
	
	err := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic checking %+v: %v\n%s", actual, r, debug.Stack())
			}
		}()
		
		$range .Members$	
			if nil != m.$.Name$ {
				fieldName = "$.Name$"
				match, err := m.$.Name$.Match(actual.$.Name$)
				if err != nil {
					return err
				} else if !match {
					if nesting, ok := m.$.Name$.(errorsutil.NestingMatcher); ok {
						return errorsutil.AggregateError(nesting.Failures())
					}
					return errors.New(m.$.Name$.FailureMessage(actual.$.Name$))
				}
			}
		$end$
		
		return nil
	}()
	if err != nil {
		errs = append(errs, errorsutil.Nest("."+fieldName, err))
	}

	return errs
}`
	structMatcherFailureMessage = `
func (m *$.|private$) FailureMessage(actual interface{}) (message string) {
	failures := make([]string, len(m.failures))
	for i := range m.failures {
		failures[i] = m.failures[i].Error()
	}
	return format.Message(reflect.TypeOf(actual).Name(),
		fmt.Sprintf("to match fields: {\n%v\n}\n", strings.Join(failures, "\n")))
}
	
func (m *$.|private$) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match fields")
}`
	structMatcherFailures = `
func (m *$.|private$) Failures() []error {
	return m.failures
}`
)

func NewStructMatcherGenerator(filename string, targetType *types.Type, imports namer.ImportTracker) generator.Generator {
	return &structMatcherGenerator{
		DefaultGen: generator.DefaultGen{
			OptionalName: filename,
		},
		targetType: targetType,
		imports:    imports,
	}
}

type structMatcherGenerator struct {
	generator.DefaultGen
	targetType *types.Type
	imports    namer.ImportTracker
}

func (g *structMatcherGenerator) Filter(ctx *generator.Context, t *types.Type) bool {
	return t.Name.Name == g.targetType.Name.Name && t.Name.Package == g.targetType.Name.Package
}

func (g *structMatcherGenerator) Namers(*generator.Context) namer.NameSystems {
	n := namer.NewRawNamer("", g.imports)

	m := namer.NewPublicNamer(0)
	m.Suffix = "Matcher"
	return map[string]namer.Namer{
		"raw":     n,
		"matcher": m,
	}
}

func (g *structMatcherGenerator) GenerateType(ctx *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, ctx, "$", "$").
		Do(structMatcherDecl, t).
		Do(structMatcherMatch, t).
		Do(structMatcherMatchFields, t).
		Do(structMatcherFailureMessage, t).
		Do(structMatcherFailures, t)

	return sw.Error()
}

func (g *structMatcherGenerator) Imports(*generator.Context) []string {
	return []string{
		fmt.Sprintf(`%s "%s"`, g.imports.LocalNameOf(g.targetType.Name.Package), g.targetType.Name.Package),
		"errors",
		"fmt",
		"reflect",
		"runtime/debug",
		"strings",
		"github.com/onsi/gomega/format",
		`errorsutil "github.com/onsi/gomega/gstruct/errors"`,
		"github.com/onsi/gomega/types",
	}
}
