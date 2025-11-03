package main

type (
	importList struct {
		imports map[string]importItem
	}

	importItem struct {
		path    string
		total   uint
		aliases map[string]uint
	}
)

func (l *importList) add(path string) {
	l.addAliased(path, "")
}

func (l *importList) addAliased(path string, alias string) {
	if len(l.imports) == 0 {
		l.imports = make(map[string]importItem)
	}

	item, ok := l.imports[path]
	if !ok {
		item = importItem{
			path: path,
		}
	}

	item.total++

	if alias != "" {
		if len(item.aliases) == 0 {
			item.aliases = make(map[string]uint)
		}
		item.aliases[alias]++
	}

	l.imports[path] = item
}
