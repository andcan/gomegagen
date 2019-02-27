package generators

import "k8s.io/gengo/namer"

func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public":  namer.NewPublicNamer(0),
		"private": namer.NewPrivateNamer(0),
	}
}

func DefaultNameSystem() string {
	return "public"
}
