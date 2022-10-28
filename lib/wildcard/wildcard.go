package wildcard

const (
	normal     = iota
	all        //*
	any        //?
	setSymbol  //[]
	rangSymbol //[a-b]
	negSymbol  //[^a]
)

type item struct {
	character byte
	set       map[byte]bool
	typeCode  int
}

func (i *item) contains(c byte) bool {
	if i.typeCode == setSymbol {
		ok := i.set[c]
		return ok
	} else if i.typeCode == rangSymbol {
		ok := i.set[c]
		if ok {
			var (
				min uint8 = 255
				max uint8 = 0
			)
			for k := range i.set {
				if min > k {
					min = k
				}
				if max < k {
					max = k
				}
			}
			return c >= min && c <= max
		}
		return !ok
	} else {
		ok := i.set[c]
		return !ok
	}
}

type Pattern struct {
	items []*item
}

func compilePattern(src string) *Pattern {
	items := make([]*item, 0)
	escape := false
	inSet := false
	var set map[byte]bool
	for _, v := range src {
		c := byte(v)
		if escape {
			items = append(items, &item{typeCode: normal, character: c})
			escape = false
		} else if c == '*' {
			items = append(items, &item{typeCode: all})
		} else if c == '?' {
			items = append(items, &item{typeCode: any})
		} else if c == '\\' {
			escape = true
		} else if c == '[' {
			if !inSet {
				inSet = true
				set = make(map[byte]bool)
			} else {
				set[c] = true
			}
		}
	}
	return &Pattern{
		items: items,
	}
}
