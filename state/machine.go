// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/txn"
	"launchpad.net/juju-core/constraints"
	"launchpad.net/juju-core/errors"
	"launchpad.net/juju-core/state/api/params"
	"launchpad.net/juju-core/state/presence"
	"launchpad.net/juju-core/utils"
	"strings"
	"time"
)

// An InstanceId is a provider-specific identifier associated with an
// instance (physical or virtual machine allocated in the provider).
type InstanceId string

// Machine represents the state of a machine.
type Machine struct {
	st  *State
	doc machineDoc
	annotator
}

// MachineJob values define responsibilities that machines may be
// expected to fulfil.
type MachineJob int

const (
	_ MachineJob = iota
	JobHostUnits
	JobManageEnviron
	JobServeAPI
)

var jobNames = []string{
	JobHostUnits:     "JobHostUnits",
	JobManageEnviron: "JobManageEnviron",
	JobServeAPI:      "JobServeAPI",
}

func (job MachineJob) String() string {
	j := int(job)
	if j <= 0 || j >= len(jobNames) {
		return fmt.Sprintf("<unknown job %d>", j)
	}
	return jobNames[j]
}

// machineDoc represents the internal state of a machine in MongoDB.
// Note the correspondence with MachineInfo in state/api/params.
type machineDoc struct {
	Id            string `bson:"_id"`
	Nonce         string
	Series        string
	ContainerType string
	InstanceId    InstanceId
	Principals    []string
	Life          Life
	Tools         *Tools `bson:",omitempty"`
	TxnRevno      int64  `bson:"txn-revno"`
	Jobs          []MachineJob
	PasswordHash  string
	Clean         bool
}

func newMachine(st *State, doc *machineDoc) *Machine {
	machine := &Machine{
		st:  st,
		doc: *doc,
	}
	machine.annotator = annotator{
		globalKey: machine.globalKey(),
		tag:       machine.Tag(),
		st:        st,
	}
	return machine
}

// Id returns the machine id.
func (m *Machine) Id() string {
	return m.doc.Id
}

// Series returns the operating system series running on the machine.
func (m *Machine) Series() string {
	return m.doc.Series
}

// ContainerType returns the type of container hosting this machine.
func (m *Machine) ContainerType() ContainerType {
	return ContainerType(m.doc.ContainerType)
}

// machineGlobalKey returns the global database key for the identified machine.
func machineGlobalKey(id string) string {
	return "m#" + id
}

// globalKey returns the global database key for the machine.
func (m *Machine) globalKey() string {
	return machineGlobalKey(m.doc.Id)
}

// MachineTag returns the tag for the
// machine with the given id.
func MachineTag(id string) string {
	return fmt.Sprintf("machine-%s", id)
}

// Tag returns a name identifying the machine that is safe to use
// as a file name.  The returned name will be different from other
// Tag values returned by any other entities from the same state.
func (m *Machine) Tag() string {
	return MachineTag(m.Id())
}

// Life returns whether the machine is Alive, Dying or Dead.
func (m *Machine) Life() Life {
	return m.doc.Life
}

// Jobs returns the responsibilities that must be fulfilled by m's agent.
func (m *Machine) Jobs() []MachineJob {
	return m.doc.Jobs
}

// AgentTools returns the tools that the agent is currently running.
// It returns an error that satisfies IsNotFound if the tools have not yet been set.
func (m *Machine) AgentTools() (*Tools, error) {
	if m.doc.Tools == nil {
		return nil, errors.NotFoundf("agent tools for machine %v", m)
	}
	tools := *m.doc.Tools
	return &tools, nil
}

// SetAgentTools sets the tools that the agent is currently running.
func (m *Machine) SetAgentTools(t *Tools) (err error) {
	defer utils.ErrorContextf(&err, "cannot set agent tools for machine %v", m)
	if t.Series == "" || t.Arch == "" {
		return fmt.Errorf("empty series or arch")
	}
	ops := []txn.Op{{
		C:      m.st.machines.Name,
		Id:     m.doc.Id,
		Assert: notDeadDoc,
		Update: D{{"$set", D{{"tools", t}}}},
	}}
	if err := m.st.runner.Run(ops, "", nil); err != nil {
		return onAbort(err, errDead)
	}
	tools := *t
	m.doc.Tools = &tools
	return nil
}

// SetMongoPassword sets the password the agent responsible for the machine
// should use to communicate with the state servers.  Previous passwords
// are invalidated.
func (m *Machine) SetMongoPassword(password string) error {
	return m.st.setMongoPassword(m.Tag(), password)
}

// SetPassword sets the password for the machine's agent.
func (m *Machine) SetPassword(password string) error {
	hp := utils.PasswordHash(password)
	ops := []txn.Op{{
		C:      m.st.machines.Name,
		Id:     m.doc.Id,
		Assert: notDeadDoc,
		Update: D{{"$set", D{{"passwordhash", hp}}}},
	}}
	if err := m.st.runner.Run(ops, "", nil); err != nil {
		return fmt.Errorf("cannot set password of machine %v: %v", m, onAbort(err, errDead))
	}
	m.doc.PasswordHash = hp
	return nil
}

// PasswordValid returns whether the given password is valid
// for the given machine.
func (m *Machine) PasswordValid(password string) bool {
	return utils.PasswordHash(password) == m.doc.PasswordHash
}

// Destroy sets the machine lifecycle to Dying if it is Alive. It does
// nothing otherwise. Destroy will fail if the machine has principal
// units assigned, or if the machine has JobManageEnviron.
// If the machine has assigned units, Destroy will return
// a HasAssignedUnitsError.
func (m *Machine) Destroy() error {
	return m.advanceLifecycle(Dying)
}

// EnsureDead sets the machine lifecycle to Dead if it is Alive or Dying.
// It does nothing otherwise. EnsureDead will fail if the machine has
// principal units assigned, or if the machine has JobManageEnviron.
// If the machine has assigned units, EnsureDead will return
// a HasAssignedUnitsError.
func (m *Machine) EnsureDead() error {
	return m.advanceLifecycle(Dead)
}

type HasAssignedUnitsError struct {
	MachineId string
	UnitNames []string
}

func (e *HasAssignedUnitsError) Error() string {
	return fmt.Sprintf("machine %s has unit %q assigned", e.MachineId, e.UnitNames[0])
}

func IsHasAssignedUnitsError(err error) bool {
	_, ok := err.(*HasAssignedUnitsError)
	return ok
}

// Containers returns the container ids belonging to a parent machine.
// TODO(wallyworld): move this method to a service
func (m *Machine) Containers() ([]string, error) {
	var mc machineContainers
	err := m.st.containerRefs.FindId(m.Id()).One(&mc)
	if err == nil {
		return mc.Children, nil
	}
	if err == mgo.ErrNotFound {
		return nil, errors.NotFoundf("container info for machine %v", m.Id())
	}
	return nil, err
}

type HasContainersError struct {
	MachineId    string
	ContainerIds []string
}

func (e *HasContainersError) Error() string {
	return fmt.Sprintf("machine %s is hosting containers %q", e.MachineId, strings.Join(e.ContainerIds, ","))
}

func IsHasContainersError(err error) bool {
	_, ok := err.(*HasContainersError)
	return ok
}

// advanceLifecycle ensures that the machine's lifecycle is no earlier
// than the supplied value. If the machine already has that lifecycle
// value, or a later one, no changes will be made to remote state. If
// the machine has any responsibilities that preclude a valid change in
// lifecycle, it will return an error.
func (original *Machine) advanceLifecycle(life Life) (err error) {
	containers, err := original.Containers()
	if err != nil {
		return err
	}
	if len(containers) > 0 {
		return &HasContainersError{
			MachineId:    original.doc.Id,
			ContainerIds: containers,
		}
	}
	m := original
	defer func() {
		if err == nil {
			// The machine's lifecycle is known to have advanced; it may be
			// known to have already advanced further than requested, in
			// which case we set the latest known valid value.
			if m == nil {
				life = Dead
			} else if m.doc.Life > life {
				life = m.doc.Life
			}
			original.doc.Life = life
		}
	}()
	// op and
	op := txn.Op{
		C:      m.st.machines.Name,
		Id:     m.doc.Id,
		Update: D{{"$set", D{{"life", life}}}},
	}
	advanceAsserts := D{
		{"jobs", D{{"$nin", []MachineJob{JobManageEnviron}}}},
		{"$or", []D{
			{{"principals", D{{"$size", 0}}}},
			{{"principals", D{{"$exists", false}}}},
		}},
	}
	// 3 attempts: one with original data, one with refreshed data, and a final
	// one intended to determine the cause of failure of the preceding attempt.
	for i := 0; i < 3; i++ {
		// If the transaction was aborted, grab a fresh copy of the machine data.
		// We don't write to original, because the expectation is that state-
		// changing methods only set the requested change on the receiver; a case
		// could perhaps be made that this is not a helpful convention in the
		// context of the new state API, but we maintain consistency in the
		// face of uncertainty.
		if i != 0 {
			if m, err = m.st.Machine(m.doc.Id); errors.IsNotFoundError(err) {
				return nil
			} else if err != nil {
				return err
			}
		}
		// Check that the life change is sane, and collect the assertions
		// necessary to determine that it remains so.
		switch life {
		case Dying:
			if m.doc.Life != Alive {
				return nil
			}
			op.Assert = append(advanceAsserts, isAliveDoc...)
		case Dead:
			if m.doc.Life == Dead {
				return nil
			}
			op.Assert = append(advanceAsserts, notDeadDoc...)
		default:
			panic(fmt.Errorf("cannot advance lifecycle to %v", life))
		}
		// Check that the machine does not have any responsibilities that
		// prevent a lifecycle change.
		for _, j := range m.doc.Jobs {
			if j == JobManageEnviron {
				// (NOTE: When we enable multiple JobManageEnviron machines,
				// the restriction will become "there must be at least one
				// machine with this job".)
				return fmt.Errorf("machine %s is required by the environment", m.doc.Id)
			}
		}
		if len(m.doc.Principals) != 0 {
			return &HasAssignedUnitsError{
				MachineId: m.doc.Id,
				UnitNames: m.doc.Principals,
			}
		}
		// Run the transaction...
		if err := m.st.runner.Run([]txn.Op{op}, "", nil); err != txn.ErrAborted {
			return err
		}
		// ...and retry on abort.
	}
	// In very rare circumstances, the final iteration above will have determined
	// no cause of failure, and attempted a final transaction: if this also failed,
	// we can be sure that the machine document is changing very fast, in a somewhat
	// surprising fashion, and that it is sensible to back off for now.
	return fmt.Errorf("machine %s cannot advance lifecycle: %v", m, ErrExcessiveContention)
}

// Remove removes the machine from state. It will fail if the machine is not
// Dead.
func (m *Machine) Remove() (err error) {
	defer utils.ErrorContextf(&err, "cannot remove machine %s", m.doc.Id)
	if m.doc.Life != Dead {
		return fmt.Errorf("machine is not dead")
	}
	ops := []txn.Op{
		{
			C:      m.st.machines.Name,
			Id:     m.doc.Id,
			Assert: txn.DocExists,
			Remove: true,
		},
		removeStatusOp(m.st, m.globalKey()),
		removeConstraintsOp(m.st, m.globalKey()),
		annotationRemoveOp(m.st, m.globalKey()),
	}
	ops = append(ops, removeContainerRefOps(m.st, m.Id())...)
	// The only abort conditions in play indicate that the machine has already
	// been removed.
	return onAbort(m.st.runner.Run(ops, "", nil), nil)
}

// Refresh refreshes the contents of the machine from the underlying
// state. It returns an error that satisfies IsNotFound if the machine has
// been removed.
func (m *Machine) Refresh() error {
	doc := machineDoc{}
	err := m.st.machines.FindId(m.doc.Id).One(&doc)
	if err == mgo.ErrNotFound {
		return errors.NotFoundf("machine %v", m)
	}
	if err != nil {
		return fmt.Errorf("cannot refresh machine %v: %v", m, err)
	}
	m.doc = doc
	return nil
}

// AgentAlive returns whether the respective remote agent is alive.
func (m *Machine) AgentAlive() (bool, error) {
	return m.st.pwatcher.Alive(m.globalKey())
}

// WaitAgentAlive blocks until the respective agent is alive.
func (m *Machine) WaitAgentAlive(timeout time.Duration) (err error) {
	defer utils.ErrorContextf(&err, "waiting for agent of machine %v", m)
	ch := make(chan presence.Change)
	m.st.pwatcher.Watch(m.globalKey(), ch)
	defer m.st.pwatcher.Unwatch(m.globalKey(), ch)
	for i := 0; i < 2; i++ {
		select {
		case change := <-ch:
			if change.Alive {
				return nil
			}
		case <-time.After(timeout):
			return fmt.Errorf("still not alive after timeout")
		case <-m.st.pwatcher.Dead():
			return m.st.pwatcher.Err()
		}
	}
	panic(fmt.Sprintf("presence reported dead status twice in a row for machine %v", m))
}

// SetAgentAlive signals that the agent for machine m is alive.
// It returns the started pinger.
func (m *Machine) SetAgentAlive() (*presence.Pinger, error) {
	p := presence.NewPinger(m.st.presence, m.globalKey())
	err := p.Start()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// InstanceId returns the provider specific instance id for this machine
// and whether it has been set.
func (m *Machine) InstanceId() (InstanceId, bool) {
	return m.doc.InstanceId, m.doc.InstanceId != ""
}

// Units returns all the units that have been assigned to the machine.
func (m *Machine) Units() (units []*Unit, err error) {
	defer utils.ErrorContextf(&err, "cannot get units assigned to machine %v", m)
	pudocs := []unitDoc{}
	err = m.st.units.Find(D{{"machineid", m.doc.Id}}).All(&pudocs)
	if err != nil {
		return nil, err
	}
	for _, pudoc := range pudocs {
		units = append(units, newUnit(m.st, &pudoc))
		docs := []unitDoc{}
		err = m.st.units.Find(D{{"principal", pudoc.Name}}).All(&docs)
		if err != nil {
			return nil, err
		}
		for _, doc := range docs {
			units = append(units, newUnit(m.st, &doc))
		}
	}
	return units, nil
}

// SetProvisioned sets the provider specific machine id and nonce for
// this machine. Once set, the instance id cannot be changed.
func (m *Machine) SetProvisioned(id InstanceId, nonce string) (err error) {
	defer utils.ErrorContextf(&err, "cannot set instance id of machine %q", m)

	if id == "" || nonce == "" {
		return fmt.Errorf("instance id and nonce cannot be empty")
	}
	notSetYet := D{{"instanceid", ""}, {"nonce", ""}}
	ops := []txn.Op{{
		C:      m.st.machines.Name,
		Id:     m.doc.Id,
		Assert: append(isAliveDoc, notSetYet...),
		Update: D{{"$set", D{{"instanceid", id}, {"nonce", nonce}}}},
	}}

	if err = m.st.runner.Run(ops, "", nil); err == nil {
		m.doc.InstanceId = id
		m.doc.Nonce = nonce
		return nil
	} else if err != txn.ErrAborted {
		return err
	} else if alive, err := isAlive(m.st.machines, m.doc.Id); err != nil {
		return err
	} else if !alive {
		return errNotAlive
	}
	return fmt.Errorf("already set")
}

// CheckProvisioned returns true if the machine was provisioned with
// the given nonce.
func (m *Machine) CheckProvisioned(nonce string) bool {
	return m.doc.Nonce == nonce && m.doc.InstanceId != ""
}

// String returns a unique description of this machine.
func (m *Machine) String() string {
	return m.doc.Id
}

// Constraints returns the exact constraints that should apply when provisioning
// an instance for the machine.
func (m *Machine) Constraints() (constraints.Value, error) {
	return readConstraints(m.st, m.globalKey())
}

// SetConstraints sets the exact constraints to apply when provisioning an
// instance for the machine. It will fail if the machine is Dead, or if it
// is already provisioned.
func (m *Machine) SetConstraints(cons constraints.Value) (err error) {
	defer utils.ErrorContextf(&err, "cannot set constraints")
	notProvisioned := D{{"instanceid", ""}}
	ops := []txn.Op{
		{
			C:      m.st.machines.Name,
			Id:     m.doc.Id,
			Assert: append(isAliveDoc, notProvisioned...),
		},
		setConstraintsOp(m.st, m.globalKey(), cons),
	}
	// 3 attempts is enough to push the ErrExcessiveContention case out of the
	// realm of plausibility: it implies local state indicating unprovisioned,
	// and remote state indicating provisioned (reasonable); but which changes
	// back to unprovisioned and then to provisioned again with *very* specific
	// timing in the course of this loop.
	for i := 0; i < 3; i++ {
		if m.doc.Life != Alive {
			return errNotAlive
		}
		if m.doc.InstanceId != "" {
			return fmt.Errorf("machine is already provisioned")
		}
		if err := m.st.runner.Run(ops, "", nil); err != txn.ErrAborted {
			return err
		}
		if m, err = m.st.Machine(m.doc.Id); err != nil {
			return err
		}
	}
	return ErrExcessiveContention
}

// Status returns the status of the machine.
func (m *Machine) Status() (status params.Status, info string, err error) {
	doc, err := getStatus(m.st, m.globalKey())
	if err != nil {
		return "", "", err
	}
	status = doc.Status
	info = doc.StatusInfo
	return
}

// SetStatus sets the status of the machine.
func (m *Machine) SetStatus(status params.Status, info string) error {
	if status == params.StatusError && info == "" {
		panic("machine error status with no info")
	}
	if status == params.StatusPending {
		panic("machine status cannot be set to pending")
	}
	doc := statusDoc{
		Status:     status,
		StatusInfo: info,
	}
	ops := []txn.Op{{
		C:      m.st.machines.Name,
		Id:     m.doc.Id,
		Assert: notDeadDoc,
	},
		updateStatusOp(m.st, m.globalKey(), doc),
	}
	if err := m.st.runner.Run(ops, "", nil); err != nil {
		return fmt.Errorf("cannot set status of machine %q: %v", m, onAbort(err, errNotAlive))
	}
	return nil
}

// Clean returns true if the machine does not have any deployed units or containers.
func (m *Machine) Clean() bool {
	return m.doc.Clean
}
