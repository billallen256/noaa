package noaa

import (
	"time"
)

type TimeSpan struct {
	Begin time.Time
	End   time.Time
}

func (ts TimeSpan) Hours() []time.Time {
	begin := ts.Begin.UTC()
	end := ts.End.UTC()

	if begin.After(end) {
		begin, end = end, begin
	}

	begin = begin.Round(time.Hour)
	end = end.Round(time.Hour)

	if begin == end {
		return []time.Time{begin}
	}

	numHours := int(end.Sub(begin).Hours())

	if numHours == 0 {
		return []time.Time{begin}
	}

	times := make([]time.Time, numHours)

	for i := 0; i < numHours; i++ {
		times[i] = begin.Add(time.Duration(i) * time.Hour)
	}

	return times
}
