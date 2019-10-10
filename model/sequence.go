package model

import "fmt"
import "math/rand"
import "time"


// --

type ActionId string

// --

type AnimationId string

var ShowAlways AnimationId = ""

// --

type ActionOccurrence struct {
	Action ActionId
	Ago float32
}

// --

type Step struct{
	startAtCounter float32
	safeStart bool
	safeStop bool
	animations []AnimationId
	actions []ActionId
	transformSteps map[*Transform]*TransformStep
}

func NewStep(startAtCounter float32, safeStart, safeStop bool, animations []AnimationId, actions []ActionId, transformSteps map[*Transform]*TransformStep) *Step {
	return &Step{
		startAtCounter: startAtCounter, 
		safeStart: safeStart,
		safeStop: safeStop,
		animations: animations,
		actions: actions,
		transformSteps: transformSteps,
	}
}

// --

type SequenceId string

type Sequence struct{
	// Config
	counterLimit float32
	counterFactor float32
	loopable bool
	multiple bool // If true, every start->stop is a different animation, taken randomly
	speedFactorWhenStopping float32
	steps []*Step

	// Runtime
	safeStarts []int
	transforms map[*Transform]interface{}

	rand *rand.Rand

	counterValue float32
	currentStep int
	actions []*ActionOccurrence
	running bool
	stopping bool
}

func New(counterLimit float32, counterFactor float32, loopable bool, multiple bool, speedFactorWhenStopping float32, steps []*Step) (*Sequence, error) {
	counterLimitEffective := counterLimit*counterFactor
	transforms := map[*Transform]interface{}{}
	safeStarts := make([]int, 0, 5)
	var anySafeStop, anyStartZero bool
	var last float32 = -1
	for i, s := range steps {
		if s.safeStart {
			safeStarts = append(safeStarts, i)
		}
		if s.startAtCounter == 0 {
			anyStartZero = true
		}
		if s.safeStop {
			anySafeStop = true
		}
		for t, ts := range s.transformSteps {
			transforms[t] = nil
			if t == nil {
				return nil, fmt.Errorf("step %d has TransformStep %v with nil target", i, ts)
			}
			if ts.Set == nil {
				return nil, fmt.Errorf("step %d has TransformStep %v without Set", i, ts)
			}
			if ts.Set.Position != nil && t.Position == nil {
				return nil, fmt.Errorf("step %d has TransformStep %v with Set.Position but no target.Position", i, ts)
			}
			if ts.Set.Rotation != nil && t.Rotation == nil {
				return nil, fmt.Errorf("step %d has TransformStep %v with Set.Rotation but no target.Rotation", i, ts)
			}
			if ts.Set.Scale != nil && t.Scale == nil {
				return nil, fmt.Errorf("step %d has TransformStep %v with Set.Scale but no target.Scale", i, ts)
			}
		}
		if s.startAtCounter < last {
			return nil, fmt.Errorf("step %d has startAtCounter %v, which can't be lesser than previous one, %v", i, s.startAtCounter, last)
		}
		if s.startAtCounter < 0 {
			return nil, fmt.Errorf("step %d has startAtCounter %v, which can't be lesser than 0", i, s.startAtCounter)
		}
		if s.startAtCounter > counterLimit {
			return nil, fmt.Errorf("step %d has startAtCounter %v, which can't be greater than counterLimit %v", i, s.startAtCounter, counterLimit)
		}
		last = s.startAtCounter
		s.startAtCounter = s.startAtCounter*counterFactor
	}
	if len(steps) == 0 {
		return nil, fmt.Errorf("steps must have at least a step")
	}
	if len(safeStarts) == 0 {
		return nil, fmt.Errorf("steps must have at least a step marked safeStart")
	}
	if !anySafeStop {
		return nil, fmt.Errorf("steps must have at least a step marked safeStop")
	}
	if !anyStartZero {
		return nil, fmt.Errorf("steps must have at least a step with startAtCounter 0")
	}
	if speedFactorWhenStopping <= 0 {
		return nil, fmt.Errorf("speedFactorWhenStoping must be greater than 0")
	}
	if counterFactor <= 0 {
		return nil, fmt.Errorf("counterFactor must be greater than 0")
	}

	seq := &Sequence{
		counterLimit: counterLimitEffective,
		counterFactor: counterFactor,
		steps: steps,
		loopable: loopable,
		multiple: multiple,
		speedFactorWhenStopping: speedFactorWhenStopping,

		safeStarts: safeStarts,
		transforms: transforms,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
	seq.Reset()
	return seq, nil
}

func NewTimer(seconds float32) *Sequence{
	s, err := New(seconds, 1, false, false, 1, []*Step{NewStep(0, true, true, nil, nil, nil)})
        if err != nil {
                panic(fmt.Sprintf("Critical error. Failed to create timer: %s", err))
        }
    return s
}
var newTimerTest = NewTimer(1)

func (s *Sequence) CounterLimit() float32 {
	return s.counterLimit
}

func (s *Sequence) CounterFactor() float32 {
	return s.counterFactor
}

func (s *Sequence) Counter() float32 {
	return s.counterValue
}

func (s *Sequence) CurrentStep() int {
	return s.currentStep
}

func (s *Sequence) Running() bool {
	return s.running
}

func (s *Sequence) Stopping() bool {
	return s.stopping
}

func (s *Sequence) Start() {
	if s.running {
		s.stopping = false
		return
	} else {
		s.Reset()
		s.setStep(s.safeStarts[s.rand.Intn(len(s.safeStarts))])
		s.running = true
		s.Add(0)
	}
}

func (s *Sequence) Stop() {
	s.stopping = true
}

func (s *Sequence) stopAndWaitForReset() {
	s.running = false
}

func (s *Sequence) Add(c float32) {
	if !s.running {
		return
	}
	if s.stopping {
		c = c*s.speedFactorWhenStopping
	}
	var lastC float32 = 0
	for c > 0.0000000001 && s.running {
		var nextHit float32
		var hitLimit bool
		var hitStep bool
		if s.currentStep == len(s.steps)-1 {
			nextHit = s.counterLimit
			hitLimit = true
		} else {
			nextHit = s.steps[s.currentStep+1].startAtCounter
			hitLimit = false
		}
		if s.counterValue + c < nextHit {
			nextHit = s.counterValue + c
			hitStep = false
		} else {
			hitStep = true
		}
		addNow := nextHit - s.counterValue
		s.counterValue += addNow
		c -= addNow
		for _, a := range s.actions {
			a.Ago += addNow
		}
		if hitStep {
			if hitLimit {
				s.setStep(0)
			} else {
				s.setStep(s.currentStep + 1)
			}
		}
		// In case we get some epsilon around, get rid of it instead of looping forever.
		lastC = c
		if c >= lastC {
			c = 0
		}
	}
	var nextHit float32
	if s.currentStep == len(s.steps)-1 {
		nextHit = s.counterLimit
	} else {
		nextHit = s.steps[s.currentStep+1].startAtCounter
	}
	cur := s.steps[s.currentStep]
	cur.ApplyTransforms(float32(s.counterValue-cur.startAtCounter)/float32(nextHit-cur.startAtCounter))
}

func (s *Sequence) PollActions() []*ActionOccurrence {
	ret := s.actions
	s.actions = []*ActionOccurrence{}
	return ret
}

func (s *Sequence) IsActiveAnimation(v AnimationId) bool {
	if v == ShowAlways {
		return true
	}
	for _, a := range s.steps[s.currentStep].animations {
		if a == v {
			return true
		}
	}
	return false
}

func (s *Sequence) Reset() {
	s.counterValue = 0
	s.actions = []*ActionOccurrence{}
	s.running = false
	s.stopping = false
}

func (s *Sequence) setStep(id int) {
	oldStep := s.steps[s.currentStep]
	oldStep.ApplyTransforms(1)
	s.currentStep = id
	step := s.steps[id]
	step.ApplyTransforms(0)
	s.counterValue = step.startAtCounter
	for _, a := range step.actions {
		s.actions = append(s.actions, &ActionOccurrence{Action: a, Ago: 0})
	}
	if step.safeStop && s.stopping {
		s.stopAndWaitForReset()
	} else if !s.loopable && id == 0 {
		s.stopAndWaitForReset()
	} else if s.multiple && step.safeStop && step.safeStart {
		s.stopAndWaitForReset()
	}
}