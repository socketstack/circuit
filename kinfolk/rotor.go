// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package kinfolk

import (
	"sync"

	"github.com/gocircuit/circuit/kit/lang"
)

type Reservoir interface {
	Add(XID)
	Scrub(XID) bool
	ScrubRandom() (XID, bool)
	View() []XID
	Len() int
	Choose() XID
}

// Neighborhood is a set of perm cross-interfaces.
type Neighborhood struct {
	sync.Mutex
	open map[lang.ReceiverID]XID
}

// NewNeighborhood creates a new rotor.
func NewNeighborhood() *Neighborhood {
	return &Neighborhood{
		open: make(map[lang.ReceiverID]XID),
	}
}

func (rtr *Neighborhood) Add(xid XID) {
	rtr.Lock()
	defer rtr.Unlock()
	rtr.open[xid.ID] = xid
}

func (rtr *Neighborhood) Scrub(xid XID) bool {
	rtr.Lock()
	defer rtr.Unlock()
	if xid.ID == 0 {
		panic("missig unique receiver id")
	}
	_, ok := rtr.open[xid.ID]
	delete(rtr.open, xid.ID)
	return ok
}

func (rtr *Neighborhood) ScrubRandom() (XID, bool) {
	rtr.Lock()
	defer rtr.Unlock()
	for hid, xid := range rtr.open {
		delete(rtr.open, hid)
		return xid, true
	}
	return XID{}, false
}

// View returns a list of all XIDs in the rotor.
func (rtr *Neighborhood) View() []XID {
	rtr.Lock()
	defer rtr.Unlock()
	open := make([]XID, 0, len(rtr.open))
	for _, xid := range rtr.open {
		open = append(open, xid)
	}
	return open
}

// Len returns the number of XIDs in the rotor.
func (rtr *Neighborhood) Len() int {
	rtr.Lock()
	defer rtr.Unlock()
	return len(rtr.open)
}

// Choose returns a randomly chosen XID.
func (rtr *Neighborhood) Choose() XID {
	rtr.Lock()
	defer rtr.Unlock()
	for _, xid := range rtr.open {
		return xid
	}
	return XID{}
}