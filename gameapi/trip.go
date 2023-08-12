package gameapi

type Trip struct {
	cities []City
	roads  []Road
	index  int
}

func NewTrip(cities []City, roads []Road) Trip {
	return Trip{cities: cities, roads: roads}
}

func (t *Trip) Cities() []City {
	return t.cities
}

func (t *Trip) Roads() []Road {
	return t.roads
}

func (t *Trip) Next() {
	t.index++
}

func (t *Trip) Src() City {
	return t.cities[0]
}

func (t *Trip) Current() City {
	if t.Arrived() {
		return t.Dst()
	}
	return t.cities[t.index]
}

func (t *Trip) Dst() City {
	return t.cities[len(t.cities)-1]
}

func (t *Trip) Arrived() bool {
	return t.index >= len(t.cities)
}

func (t *Trip) String() string {
	s := ""
	for i, city := range t.cities {
		if t.index == i {
			s += "[" + city.Name() + "]" + "-"
		} else {
			s += city.Name() + "-"
		}
	}
	return s[:len(s)-1]
}
