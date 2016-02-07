package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/orcalabs/plugin-sdk/go/plugin/message"
	"github.com/orcalabs/plugin-sdk/go/plugin/parameter"
)

// Types of Start Messages
const (
	TriggerStart = "trigger_start"
	ActionStart  = "action_start"
)

// init initializes the plugin package, collecting parameters from the command line
// Functionality comes from param.go
func init() {
	// defaults to stdin
	parameter.Stdin = parameter.NewParamSet(os.Stdin)

	// check for params after the double dash
	// in the command string
	for i, argv := range os.Args {
		if argv == "--" {
			arg := os.Args[i+1]
			buf := bytes.NewBufferString(arg)
			parameter.Stdin = parameter.NewParamSet(buf)
			break
		}
	}
}

// Plugin holds the common Plugin information
type Plugin struct {
	Name     string
	triggers map[string]Triggerable
	actions  map[string]Actionable
}

// Init initializes a Plugin
func (p *Plugin) Init(name string) {
	p.Name = name
	p.triggers = map[string]Triggerable{}
	p.actions = map[string]Actionable{}
}

// Run runs a Plugin
func (p *Plugin) Run() error {
	m := message.Message{}

	parameter.Param("type", &m.Type)
	parameter.Param("version", &m.Version)
	parameter.Param("body", &m.Body)
	err := parameter.Parse()

	if err != nil {
		return fmt.Errorf("Unable to deserialize message: %+v", err)
	}

	switch m.Type {
	case TriggerStart:
		start := message.TriggerStart{}
		err = m.UnmarshalBody(&start)
		if err != nil {
			return err
		}

		// lookup the trigger that matches
		trigger, err := p.LookupTrigger(start.Trigger)

		if err != nil {
			return err
		}

		task := &triggerTask{
			message:    &start,
			trigger:    trigger,
			dispatcher: triggerDispatcher(),
		}
		return task.Run()
	case ActionStart:
		start := message.ActionStart{}
		err = m.UnmarshalBody(&start)
		if err != nil {
			return err
		}

		// lookup the Action that matches
		action, err := p.LookupAction(start.Action)

		if err != nil {
			return err
		}

		task := &actionTask{
			message:    &start,
			action:     action,
			dispatcher: actionDispatcher(),
		}
		return task.Run()
	default:
		return fmt.Errorf("Unexpected message type: %s", m.Type)
	}
}

// AddTrigger adds triggers to the map of Plugins triggers
func (p Plugin) AddTrigger(trigger Triggerable) error {

	if trigger.Name() == "" {
		return errors.New("No Name() was found for the trigger.")
	}

	p.triggers[trigger.Name()] = trigger
	return nil
}

// LookupTrigger triggers on the given trigger if it exists
func (p Plugin) LookupTrigger(trigger string) (Triggerable, error) {
	if t, ok := p.triggers[trigger]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("Failed to LookupTrigger() with %s. Trigger not valid with plugin: %s.", trigger, p.Name)
}

// Triggers returns the map of Triggerables in the Plugin
func (p Plugin) Triggers() map[string]Triggerable {
	return p.triggers
}

// AddAction adds triggers to the map of Plugins triggers
func (p Plugin) AddAction(action Actionable) error {
	if action.Name() == "" {
		return errors.New("No Name() was found for the action.")
	}

	p.actions[action.Name()] = action
	return nil
}

// LookupAction will fetch the given action if it exists
func (p Plugin) LookupAction(action string) (Actionable, error) {
	if a, ok := p.actions[action]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("Failed to LookupAction() with action: %s. Action not valid with plugin: %s.", action, p.Name)
}

// Actions returns the map of Actionables in the Plugin
func (p Plugin) Actions() map[string]Actionable {
	return p.actions
}
