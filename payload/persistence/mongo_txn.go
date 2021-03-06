// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package persistence

import (
	"github.com/juju/errors"
	jujutxn "github.com/juju/txn"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"

	"github.com/juju/juju/payload"
)

type payloadsTxnRunner interface {
	// Run runs the transaction generated by the provided factory
	// function. It may be retried several times.
	Run(transactions jujutxn.TransactionSource) error
}

type payloadsTransaction interface {
	checkAssertsAndUpdate(pq payloadsQueries) error
	ops() []txn.Op
}

type payloadsTransactions struct {
	queries payloadsQueries
	runner  payloadsTxnRunner
}

func (pt payloadsTransactions) run(ptxn payloadsTransaction) error {
	buildTxn := pt.newTxnSource(ptxn)
	if err := pt.runner.Run(buildTxn); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (pt payloadsTransactions) newTxnSource(ptxn payloadsTransaction) jujutxn.TransactionSource {
	return func(attempt int) ([]txn.Op, error) {
		// We always check the asserts manually before returning ops.
		if err := ptxn.checkAssertsAndUpdate(pt.queries); err != nil {
			return nil, errors.Trace(err)
		}
		// If the asserts check out then it was probably a transient
		// error, so try again.

		return ptxn.ops(), nil
	}
}

type upsertPayloadTxn struct {
	payload payload.FullPayloadInfo

	exists bool
}

func (utxn *upsertPayloadTxn) checkAssertsAndUpdate(pq payloadsQueries) error {
	utxn.exists = false
	query := payloadIDQuery(utxn.payload.Unit, utxn.payload.Name)
	_, err := pq.one(query)
	if err == nil {
		utxn.exists = true
	} else if !errors.IsNotFound(err) {
		return errors.Trace(err)
	}

	return nil
}

func (utxn upsertPayloadTxn) ops() []txn.Op {
	doc := newPayloadDoc(utxn.payload)
	var ops []txn.Op
	if utxn.exists {
		ops = append(ops, txn.Op{
			C:      payloadsC,
			Id:     doc.DocID,
			Assert: txn.DocExists,
			Remove: true,
		})
	}
	// TODO(ericsnow) Add unitPersistence.newEnsureAliveOp(pp.unit)?
	return append(ops, txn.Op{
		C:      payloadsC,
		Id:     doc.DocID,
		Assert: txn.DocMissing,
		Insert: doc,
	})
}

type setPayloadStatusTxn struct {
	unit   string
	name   string
	status string
}

func (stxn *setPayloadStatusTxn) checkAssertsAndUpdate(pq payloadsQueries) error {
	query := payloadIDQuery(stxn.unit, stxn.name)
	_, err := pq.one(query)
	if errors.IsNotFound(err) {
		return errors.Annotatef(payload.ErrNotFound, "(%s)", stxn.name)
	}
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (stxn setPayloadStatusTxn) ops() []txn.Op {
	id := payloadID(stxn.unit, stxn.name)
	updates := bson.D{
		{"state", stxn.status},
	}
	// TODO(ericsnow) Add unitPersistence.newEnsureAliveOp(pp.unit)?
	return []txn.Op{{
		C:      payloadsC,
		Id:     id,
		Assert: txn.DocExists,
		Update: bson.D{{"$set", updates}},
	}}
}

type removePayloadTxn struct {
	unit string
	name string
}

func (rtxn *removePayloadTxn) checkAssertsAndUpdate(pq payloadsQueries) error {
	query := payloadIDQuery(rtxn.unit, rtxn.name)
	_, err := pq.one(query)
	if errors.IsNotFound(err) {
		// Must have already beeen removed! The business logic
		// (i.e. state) can decide whether or not to ignore this.
		return errors.Annotatef(payload.ErrNotFound, "(%s)", rtxn.name)
	}
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (rtxn removePayloadTxn) ops() []txn.Op {
	id := payloadID(rtxn.unit, rtxn.name)
	// TODO(ericsnow) Add unitPersistence.newEnsureAliveOp(pp.unit)?
	return []txn.Op{{
		C:      payloadsC,
		Id:     id,
		Assert: txn.DocExists,
		Remove: true,
	}}
}
