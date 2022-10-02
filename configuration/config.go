package configuration

import "xi/src/syncer"

func Packages () []syncer.Package {
	return []syncer.Package {
		{Name: "neofetch"},
		{Name: "bat"},
	}
}
