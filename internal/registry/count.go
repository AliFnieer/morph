package registry

func Count() int {
	count := 0
	type p struct{ f, t string }
	seen := map[p]bool{}
	for _, c := range converters {
		for _, f := range c.From() {
			for _, t := range c.To() {
				pr := p{f, t}
				if seen[pr] {
					continue
				}
				seen[pr] = true
				count++
			}
		}
	}
	return count
}
