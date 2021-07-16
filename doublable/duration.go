package doublable

import (
	"time"
)

type Duration struct {
	Min      time.Duration
	Max      time.Duration
	duration time.Duration
}

func (d *Duration) Double() {
	d.duration = d.withinRange(d.duration * 2)
}

func (d *Duration) Halve() {
	d.duration = d.withinRange(d.duration / 2)
}

func (d *Duration) Minimize() {
	d.duration = d.Min
}

func (d *Duration) Maximize() {
	d.duration = d.Max
}

func (d *Duration) Duration() time.Duration {
	return d.withinRange(d.duration)
}

func (d *Duration) withinRange(duration time.Duration) time.Duration {
	if duration < d.Min {
		return d.Min
	} else if duration > d.Max {
		return d.Max
	} else {
		return duration
	}
}
