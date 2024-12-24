package builder

type Builder struct {
	workflow *Workflow
}

func (b *Builder) Imports() []string {
	imports := []string{b.workflow.Manager.Package}

	if len(b.workflow.Transitions) == 0 {
		return imports
	}

	for _, trans := range b.workflow.Transitions {
		if len(trans.Guards) == 0 {
			continue
		}

		for _, g := range trans.Guards {
			if g.Package == "" {
				continue
			}

			imports = append(imports, g.Package)
		}
	}

	return imports
}
