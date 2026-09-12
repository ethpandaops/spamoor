package frametxfuzz

import (
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// coverage is what a run reports.
//
// The scenario deliberately does not judge outcomes. Whether a frame should have failed,
// what an instruction should have returned, whether a shape ought to propagate -- every
// one of those is a reading of a specification that is still moving, and a fuzzer that
// enshrines its author's reading turns a client disagreement into a false alarm about
// the client that disagreed. On a network of more than one node a genuine disagreement
// splits the chain on its own, which is the signal worth having.
//
// So what is tracked here is what was reached: which dimensions each generated
// transaction exercised, how many landed, and what the chain said about the ones it
// refused. A dimension with a zero count is the actionable result -- it means the
// generator never got there.
type coverage struct {
	mutex  sync.Mutex
	logger logrus.FieldLogger
	seed   string

	// dimensions counts how many transactions exercised each dimension.
	dimensions map[string]uint64

	// generated, submitted and confirmed track the stream.
	generated uint64
	confirmed uint64
	refused   uint64

	// invalidSent and invalidAccepted track the transactions built to be refused. An
	// accepted one is worth looking at, but it is reported as an observation rather
	// than as a verdict: what a client must reject is exactly what is under discussion.
	invalidSent     uint64
	invalidAccepted uint64

	// reasons remembers what the chain said about each violation, so a reason that
	// changes between runs or between clients is visible.
	reasons map[string]string
}

// newCoverage returns an empty tracker.
func newCoverage(logger logrus.FieldLogger, seed string) *coverage {
	return &coverage{
		logger:     logger,
		seed:       seed,
		dimensions: map[string]uint64{},
		reasons:    map[string]string{},
	}
}

// record notes the dimensions a generated transaction exercised.
func (c *coverage) record(dimensions []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.generated++

	seen := map[string]bool{}

	for _, dimension := range dimensions {
		if seen[dimension] {
			continue
		}

		seen[dimension] = true
		c.dimensions[dimension]++
	}
}

// confirmedOne notes a transaction that landed.
func (c *coverage) confirmedOne() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.confirmed++
}

// refusedOne notes a transaction the chain would not take, with what it said.
func (c *coverage) refusedOne(recipe *Recipe, reason string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.refused++

	key := recipe.refusalKey()

	// Numbers are stripped before comparing: an index or a gas figure inside the
	// message varies between transactions without the reason having changed.
	normalized := stripNumbers(reason)

	previous, seen := c.reasons[key]
	c.reasons[key] = normalized

	if seen && previous != normalized {
		c.logger.WithField("reproduce", c.reproduceLine(recipe)).Infof(
			"%q was refused with a different reason than before: %q, previously %q", key, reason, previous)
	}
}

// invalidSubmitted notes a transaction built to be refused, and whether it was.
func (c *coverage) invalidSubmitted(recipe *Recipe, accepted bool) {
	c.mutex.Lock()
	c.invalidSent++

	if accepted {
		c.invalidAccepted++
	}
	c.mutex.Unlock()

	if !accepted {
		return
	}

	c.logger.WithField("reproduce", c.reproduceLine(recipe)).Infof(
		"the chain accepted a transaction carrying the %q violation | recipe %s", recipe.Invalid, recipe)
}

// stripNumbers removes digit runs so two messages that differ only in an index or a gas
// figure compare equal.
func stripNumbers(reason string) string {
	var out strings.Builder

	skipping := false

	for _, r := range reason {
		if r >= '0' && r <= '9' {
			if !skipping {
				out.WriteByte('#')
				skipping = true
			}

			continue
		}

		skipping = false

		out.WriteRune(r)
	}

	return out.String()
}

// reproduceLine renders the command that regenerates one transaction.
func (c *coverage) reproduceLine(recipe *Recipe) string {
	return "--payload-seed 0x" + c.seed + " --tx-id-offset " + uint64String(recipe.Index) + " -c 1"
}

// summary renders the end-of-run report.
func (c *coverage) summary() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	parts := []string{
		uint64String(c.generated) + " generated",
		uint64String(c.confirmed) + " confirmed",
		uint64String(c.refused) + " refused",
	}

	if c.invalidSent > 0 {
		parts = append(parts, uint64String(c.invalidSent)+" invalid sent ("+
			uint64String(c.invalidAccepted)+" accepted)")
	}

	names := make([]string, 0, len(c.dimensions))
	for name := range c.dimensions {
		names = append(names, name)
	}

	sort.Strings(names)

	covered := make([]string, 0, len(names))
	for _, name := range names {
		covered = append(covered, name+"="+uint64String(c.dimensions[name]))
	}

	if len(covered) > 0 {
		parts = append(parts, "covered ["+strings.Join(covered, " ")+"]")
	}

	return strings.Join(parts, ", ")
}

// refusals renders what the chain said about each shape it would not take, one line per
// shape. It is the other half of the coverage report: knowing a combination was reached
// is only useful alongside what happened to it.
func (c *coverage) refusals() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := make([]string, 0, len(c.reasons))
	for key := range c.reasons {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+": "+c.reasons[key])
	}

	return lines
}

// uint64String renders a number without pulling in a formatter.
func uint64String(v uint64) string {
	if v == 0 {
		return "0"
	}

	digits := [20]byte{}
	i := len(digits)

	for v > 0 {
		i--
		digits[i] = byte('0' + v%10)
		v /= 10
	}

	return string(digits[i:])
}
