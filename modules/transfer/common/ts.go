package common

func AlignTs(ts int64, period int64) int64 {
	return ts - ts%period
}
