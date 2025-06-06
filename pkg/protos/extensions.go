package protos

import "strconv"

func ParseJobScope(s string) (JobScope, error) {
	if js, ok := JobScope_value[s]; ok {
		return JobScope(js), nil
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	if _, ok := JobScope_name[int32(i)]; ok {
		return JobScope(i), nil
	}
	return 0, nil
}
