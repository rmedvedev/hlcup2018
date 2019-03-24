package main

var indexes = Indexes{}

type Index map[string][]uint32

func (index Index) Add(value string, id uint32) {
	index[value] = append(index[value], id)
}

type Indexes map[string]Index

type IndexOperation struct {
	Indexes   [][]uint32
	Condition string
}

func (indexes Indexes) NewIndex(field string) {
	indexes[field] = Index{}
}

func (indexes Indexes) GetIds(field string, value string) []uint32 {
	return indexes[field][value]
}

func (indexes Indexes) UpdateIndex(field string, id uint32, value string) {
	indexes[field].Add(value, id)
}

func rangeIndexes(keys *map[int]map[int]int, indexes []IndexOperation) uint32 {

	if indexes[0].Condition == "OR" {

		if _, ok := (*keys)[0]; !ok {
			(*keys)[0] = make(map[int]int)
		}

		for id, keysTmp := foreachSorted((*keys)[0], indexes[0].Indexes...); id > 0; id, keysTmp = foreachSorted(keysTmp, indexes[0].Indexes...) {
			(*keys)[0] = keysTmp
			fail := false

			if len(indexes) > 1 || len(indexes[0].Indexes) > 1 {
				for keyIndex, index := range indexes {
					orEnd := false
					if _, ok := (*keys)[keyIndex]; !ok {
						(*keys)[keyIndex] = make(map[int]int)
					}
					for keyIds, ids := range index.Indexes {
						//пропускаем первый индекс так как итерируем по нему
						if keyIndex == 0 {
							continue
						}

						//текущая позиция в массиве
						curIndex := (*keys)[keyIndex][keyIds]

						if curIndex > (len(ids)-1) && index.Condition == "OR" {
							fail = true
							break
						}

						for _, tmpId := range ids[curIndex:] {
							if id == tmpId {
								(*keys)[keyIndex][keyIds]++
								if index.Condition == "OR" {
									orEnd = true
								}
								break
							}
							if id < tmpId {
								if (*keys)[keyIndex][keyIds] >= (len(ids) - 1) {
									fail = true
									break
								}
								(*keys)[keyIndex][keyIds]++
								continue
							} else {
								if index.Condition == "OR" && keyIds != len(index.Indexes)-1 {
									break
								} else {
									fail = true
									break
								}
							}
						}

						if fail || orEnd {
							break
						}
					}

					if fail {
						break
					}
				}
			} else {
				return id
			}

			if fail {
				continue
			} else {
				return id
			}

		}

		return 0

	} else {
		if _, ok := (*keys)[0]; !ok {
			(*keys)[0] = make(map[int]int)
		}

		for _, id := range indexes[0].Indexes[0][(*keys)[0][0]:] {
			(*keys)[0][0]++
			fail := false
			if len(indexes) > 1 || len(indexes[0].Indexes) > 1 {
				for keyIndex, index := range indexes {
					orEnd := false
					if _, ok := (*keys)[keyIndex]; !ok {
						(*keys)[keyIndex] = make(map[int]int)
					}
					for keyIds, ids := range index.Indexes {
						//пропускаем первый индекс так как итерируем по нему
						if keyIndex == 0 && keyIds == 0 {
							continue
						}

						//текущая позиция в массиве
						curIndex := (*keys)[keyIndex][keyIds]

						if curIndex > (len(ids)-1) && index.Condition == "AND" {
							return 0
						}

						for _, tmpId := range ids[curIndex:] {
							if id == tmpId {
								(*keys)[keyIndex][keyIds]++
								if index.Condition == "OR" {
									orEnd = true
								}
								break
							}
							if id < tmpId {
								if (*keys)[keyIndex][keyIds] >= (len(ids) - 1) {
									fail = true
									break
								}
								(*keys)[keyIndex][keyIds]++
								continue
							} else {
								if index.Condition == "OR" && keyIds != len(index.Indexes)-1 {
									break
								} else {
									fail = true
									break
								}
							}
						}

						if fail || orEnd {
							break
						}
					}

					if fail {
						break
					}
				}
			} else {
				return id
			}

			if fail {
				continue
			} else {
				return id
			}
		}

		return 0
	}
}

func foreachSorted(keys map[int]int, arrays ...[]uint32) (uint32, map[int]int) {
	var id uint32
	var keyTmp int
	for key, ids := range arrays {
		if keys[key] > (len(ids) - 1) {
			continue
		}
		if id == 0 {
			id = ids[keys[key]]
			keyTmp = key
			continue
		}

		if id < ids[keys[key]] {
			id = ids[keys[key]]
			keyTmp = key
			continue
		}

		if id == ids[keys[key]] {
			keys[key]++
			continue
		}
	}

	if id > 0 {
		keys[keyTmp]++
	}

	return id, keys
}
