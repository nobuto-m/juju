// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"
	"strconv"
	"time"

	"github.com/juju/errors"
	"gopkg.in/juju/names.v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"

	"github.com/juju/juju/core/migration"
)

// This file contains functionality for managing the state documents
// used by Juju to track model migrations.

// ModelMigration represents the state of an migration attempt for a
// model.
type ModelMigration interface {
	// Id returns a unique identifier for the model migration.
	Id() string

	// ModelUUID returns the UUID for the model being migrated.
	ModelUUID() string

	// Attempt returns the migration attempt identifier. This
	// increments for each migration attempt for the model.
	Attempt() (int, error)

	// StartTime returns the time when the migration was started.
	StartTime() time.Time

	// SuccessTime returns the time when the migration reached
	// SUCCESS.
	SuccessTime() time.Time

	// EndTime returns the time when the migration reached DONE or
	// REAPFAILED.
	EndTime() time.Time

	// Phase returns the migration's phase.
	Phase() (migration.Phase, error)

	// PhaseChangedTime returns the time when the migration's phase
	// last changed.
	PhaseChangedTime() time.Time

	// StatusMessage returns human readable text about the current
	// progress of the migration.
	StatusMessage() string

	// InitiatedBy returns username the initiated the migration.
	InitiatedBy() string

	// TargetInfo returns the details required to connect to the
	// migration's target controller.
	TargetInfo() (*migration.TargetInfo, error)

	// SetPhase sets the phase of the migration. An error will be
	// returned if the new phase does not follow the current phase or
	// if the migration is no longer active.
	SetPhase(nextPhase migration.Phase) error

	// SetStatusMessage sets some human readable text about the
	// current progress of the migration.
	SetStatusMessage(text string) error

	// Refresh updates the contents of the ModelMigration from the
	// underlying state.
	Refresh() error
}

// modelMigration is an implementation of ModelMigration.
type modelMigration struct {
	st        *State
	doc       modelMigDoc
	statusDoc modelMigStatusDoc
}

// modelMigDoc holds parameters of a migration attempt for a
// model. These are written into migrationsC.
type modelMigDoc struct {
	// Id holds migration document key. It has the format
	// "uuid:sequence".
	Id string `bson:"_id"`

	// The UUID of the model being migrated.
	ModelUUID string `bson:"model-uuid"`

	// InitiatedBy holds the username of the user that triggered the
	// migration. It should be in "user@domain" format.
	InitiatedBy string `bson:"initiated-by"`

	// TargetController holds the UUID of the target controller.
	TargetController string `bson:"target-controller"`

	// TargetAddrs holds the host:port values for the target API
	// server.
	TargetAddrs []string `bson:"target-addrs"`

	// TargetCACert holds the certificate to validate the target API
	// server's TLS certificate.
	TargetCACert string `bson:"target-cacert"`

	// TargetAuthTag holds a string representation of the tag to
	// authenticate to the target controller with.
	TargetAuthTag string `bson:"target-entity"`

	// TargetPassword holds the password to use with TargetAuthTag
	// when authenticating.
	TargetPassword string `bson:"target-password"`
}

// modelMigStatusDoc tracks the progress of a migration attempt for a
// model. These are written into migrationsStatusC.
//
// There is exactly one document in migrationsStatusC for each
// document in migrationsC. Separating them allows for watching
// for new model migrations without being woken up for each model
// migration status change.
type modelMigStatusDoc struct {
	// These are the same as the ids as migrationsC.
	// "uuid:sequence".
	Id string `bson:"_id"`

	// StartTime holds the time the migration started (stored as per
	// UnixNano).
	StartTime int64 `bson:"start-time"`

	// StartTime holds the time the migration reached the SUCCESS
	// phase (stored as per UnixNano).
	SuccessTime int64 `bson:"success-time"`

	// EndTime holds the time the migration reached a terminal (end)
	// phase (stored as per UnixNano).
	EndTime int64 `bson:"end-time"`

	// Phase holds the current migration phase. This should be one of
	// the string representations of the core/migrations.Phase
	// constants.
	Phase string `bson:"phase"`

	// PhaseChangedTime holds the time that Phase last changed (stored
	// as per UnixNano).
	PhaseChangedTime int64 `bson:"phase-changed-time"`

	// StatusMessage holds a human readable message about the
	// migration's progress.
	StatusMessage string `bson:"status-message"`
}

// Id implements ModelMigration.
func (mig *modelMigration) Id() string {
	return mig.doc.Id
}

// ModelUUID implements ModelMigration.
func (mig *modelMigration) ModelUUID() string {
	return mig.doc.ModelUUID
}

// Attempt implements ModelMigration.
func (mig *modelMigration) Attempt() (int, error) {
	attempt, err := strconv.Atoi(mig.st.localID(mig.doc.Id))
	if err != nil {
		// This really shouldn't happen.
		return -1, errors.Errorf("invalid migration id: %v", mig.doc.Id)
	}
	return attempt, nil
}

// StartTime implements ModelMigration.
func (mig *modelMigration) StartTime() time.Time {
	return unixNanoToTime0(mig.statusDoc.StartTime)
}

// SuccessTime implements ModelMigration.
func (mig *modelMigration) SuccessTime() time.Time {
	return unixNanoToTime0(mig.statusDoc.SuccessTime)
}

// EndTime implements ModelMigration.
func (mig *modelMigration) EndTime() time.Time {
	return unixNanoToTime0(mig.statusDoc.EndTime)
}

// Phase implements ModelMigration.
func (mig *modelMigration) Phase() (migration.Phase, error) {
	phase, ok := migration.ParsePhase(mig.statusDoc.Phase)
	if !ok {
		return phase, errors.Errorf("invalid phase in DB: %v", mig.statusDoc.Phase)
	}
	return phase, nil
}

// PhaseChangedTime implements ModelMigration.
func (mig *modelMigration) PhaseChangedTime() time.Time {
	return unixNanoToTime0(mig.statusDoc.PhaseChangedTime)
}

// StatusMessage implements ModelMigration.
func (mig *modelMigration) StatusMessage() string {
	return mig.statusDoc.StatusMessage
}

// InitiatedBy implements ModelMigration.
func (mig *modelMigration) InitiatedBy() string {
	return mig.doc.InitiatedBy
}

// TargetInfo implements ModelMigration.
func (mig *modelMigration) TargetInfo() (*migration.TargetInfo, error) {
	authTag, err := names.ParseUserTag(mig.doc.TargetAuthTag)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &migration.TargetInfo{
		ControllerTag: names.NewModelTag(mig.doc.TargetController),
		Addrs:         mig.doc.TargetAddrs,
		CACert:        mig.doc.TargetCACert,
		AuthTag:       authTag,
		Password:      mig.doc.TargetPassword,
	}, nil
}

// SetPhase implements ModelMigration.
func (mig *modelMigration) SetPhase(nextPhase migration.Phase) error {
	now := GetClock().Now().UnixNano()

	phase, err := mig.Phase()
	if err != nil {
		return errors.Trace(err)
	}

	if nextPhase == phase {
		return nil // Already at that phase. Nothing to do.
	}
	if !phase.CanTransitionTo(nextPhase) {
		return errors.Errorf("illegal phase change: %s -> %s", phase, nextPhase)
	}

	nextDoc := mig.statusDoc
	nextDoc.Phase = nextPhase.String()
	nextDoc.PhaseChangedTime = now
	update := bson.M{
		"phase":              nextDoc.Phase,
		"phase-changed-time": now,
	}
	if nextPhase == migration.SUCCESS {
		nextDoc.SuccessTime = now
		update["success-time"] = now
	}
	var ops []txn.Op
	if nextPhase.IsTerminal() {
		nextDoc.EndTime = now
		update["end-time"] = now
		ops = append(ops, txn.Op{
			C:      migrationsActiveC,
			Id:     mig.doc.ModelUUID,
			Assert: txn.DocExists,
			Remove: true,
		})
	}

	ops = append(ops, txn.Op{
		C:      migrationsStatusC,
		Id:     mig.statusDoc.Id,
		Update: bson.M{"$set": update},
		// Ensure phase hasn't changed underneath us
		Assert: bson.M{"phase": mig.statusDoc.Phase},
	})

	if err := mig.st.runTransaction(ops); err == txn.ErrAborted {
		return errors.New("phase already changed")
	} else if err != nil {
		return errors.Annotate(err, "failed to update phase")
	}

	mig.statusDoc = nextDoc
	return nil
}

// SetStatusMessage implements ModelMigration.
func (mig *modelMigration) SetStatusMessage(text string) error {
	ops := []txn.Op{{
		C:      migrationsStatusC,
		Id:     mig.statusDoc.Id,
		Update: bson.M{"$set": bson.M{"status-message": text}},
		Assert: txn.DocExists,
	}}
	if err := mig.st.runTransaction(ops); err != nil {
		return errors.Annotate(err, "failed to set migration status")
	}
	mig.statusDoc.StatusMessage = text
	return nil
}

// Refresh implements ModelMigration.
func (mig *modelMigration) Refresh() error {
	// Only the status document is updated. The modelMigDoc is static
	// after creation.
	statusColl, closer := mig.st.getCollection(migrationsStatusC)
	defer closer()
	var statusDoc modelMigStatusDoc
	err := statusColl.FindId(mig.doc.Id).One(&statusDoc)
	if err == mgo.ErrNotFound {
		return errors.NotFoundf("migration status")
	} else if err != nil {
		return errors.Annotate(err, "migration status lookup failed")
	}

	mig.statusDoc = statusDoc
	return nil
}

// ModelMigrationSpec holds the information required to create a
// ModelMigration instance.
type ModelMigrationSpec struct {
	InitiatedBy names.UserTag
	TargetInfo  migration.TargetInfo
}

// Validate returns an error if the ModelMigrationSpec contains bad
// data. Nil is returned otherwise.
func (spec *ModelMigrationSpec) Validate() error {
	if !names.IsValidUser(spec.InitiatedBy.Id()) {
		return errors.NotValidf("InitiatedBy")
	}
	return spec.TargetInfo.Validate()
}

// CreateModelMigration initialises state that tracks a model
// migration. It will return an error if there is already a
// model migration in progress.
func (st *State) CreateModelMigration(spec ModelMigrationSpec) (ModelMigration, error) {
	if st.IsController() {
		return nil, errors.New("controllers can't be migrated")
	}
	if err := spec.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	if err := checkTargetController(st, spec.TargetInfo.ControllerTag); err != nil {
		return nil, errors.Trace(err)
	}

	now := GetClock().Now().UnixNano()
	modelUUID := st.ModelUUID()
	var doc modelMigDoc
	var statusDoc modelMigStatusDoc
	buildTxn := func(int) ([]txn.Op, error) {
		model, err := st.Model()
		if err != nil {
			return nil, errors.Annotate(err, "failed to load model")
		}
		if model.Life() != Alive {
			return nil, errors.New("model is not alive")
		}

		if isActive, err := st.IsModelMigrationActive(); err != nil {
			return nil, errors.Trace(err)
		} else if isActive {
			return nil, errors.New("already in progress")
		}

		seq, err := st.sequence("modelmigration")
		if err != nil {
			return nil, errors.Trace(err)
		}

		id := fmt.Sprintf("%s:%d", modelUUID, seq)
		doc = modelMigDoc{
			Id:               id,
			ModelUUID:        modelUUID,
			InitiatedBy:      spec.InitiatedBy.Id(),
			TargetController: spec.TargetInfo.ControllerTag.Id(),
			TargetAddrs:      spec.TargetInfo.Addrs,
			TargetCACert:     spec.TargetInfo.CACert,
			TargetAuthTag:    spec.TargetInfo.AuthTag.String(),
			TargetPassword:   spec.TargetInfo.Password,
		}
		statusDoc = modelMigStatusDoc{
			Id:               id,
			StartTime:        now,
			Phase:            migration.QUIESCE.String(),
			PhaseChangedTime: now,
		}
		return []txn.Op{{
			C:      migrationsC,
			Id:     doc.Id,
			Assert: txn.DocMissing,
			Insert: &doc,
		}, {
			C:      migrationsStatusC,
			Id:     statusDoc.Id,
			Assert: txn.DocMissing,
			Insert: &statusDoc,
		}, {
			C:      migrationsActiveC,
			Id:     modelUUID,
			Assert: txn.DocMissing,
			Insert: bson.M{"id": doc.Id},
		}, model.assertActiveOp(),
		}, nil
	}
	if err := st.run(buildTxn); err != nil {
		return nil, errors.Annotate(err, "failed to create migration")
	}

	return &modelMigration{
		doc:       doc,
		statusDoc: statusDoc,
		st:        st,
	}, nil
}

func checkTargetController(st *State, targetControllerTag names.ModelTag) error {
	currentController, err := st.ControllerModel()
	if err != nil {
		return errors.Annotate(err, "failed to load existing controller model")
	}
	if targetControllerTag == currentController.ModelTag() {
		return errors.New("model already attached to target controller")
	}
	return nil
}

// GetModelMigration returns the most recent ModelMigration for a
// model (if any).
func (st *State) GetModelMigration() (ModelMigration, error) {
	migColl, closer := st.getCollection(migrationsC)
	defer closer()

	query := migColl.Find(bson.M{"model-uuid": st.ModelUUID()})
	query = query.Sort("-_id").Limit(1)
	var doc modelMigDoc
	err := query.One(&doc)
	if err == mgo.ErrNotFound {
		return nil, errors.NotFoundf("migration")
	} else if err != nil {
		return nil, errors.Annotate(err, "migration lookup failed")
	}

	statusColl, closer := st.getCollection(migrationsStatusC)
	defer closer()
	var statusDoc modelMigStatusDoc
	err = statusColl.FindId(doc.Id).One(&statusDoc)
	if err != nil {
		return nil, errors.Annotate(err, "failed to find status document")
	}

	return &modelMigration{
		doc:       doc,
		statusDoc: statusDoc,
		st:        st,
	}, nil
}

// IsModelMigrationActive return true if a migration is in progress for
// the model associated with the State.
func (st *State) IsModelMigrationActive() (bool, error) {
	active, closer := st.getCollection(migrationsActiveC)
	defer closer()
	n, err := active.FindId(st.ModelUUID()).Count()
	if err != nil {
		return false, errors.Trace(err)
	}
	return n > 0, nil
}

func unixNanoToTime0(i int64) time.Time {
	if i == 0 {
		return time.Time{}
	}
	return time.Unix(0, i)
}
