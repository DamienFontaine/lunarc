// +build integration

// Copyright (c) - Damien Fontaine <damien.fontaine@lineolia.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package mongo

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Chdir("../../testdata")
	os.Exit(m.Run())
}

func TestNewMongoNormal(t *testing.T) {
	mongo, err := NewMongo("config.yml", "staging")
	if err != nil {
		t.Fatalf("NewMongo must realize a success connection but %v", err)
	}
	if mongo.Client == nil {
		t.Fatalf("NewMongo must have a session")
	}
	if mongo.Database == nil {
		t.Fatalf("NewMongo must have a Database")
	}
}

func TestDisconnectNormal(t *testing.T) {
	mongo, err := NewMongo("config.yml", "staging")
	if err != nil {
		t.Fatalf("NewMongo must realize a success connection but %v", err)
	}
	err = mongo.Disconnect()
	if err != nil {
		t.Fatalf("Mongo must disconnect but %v", err)
	}
}

func TestNewMongoWithBadPort(t *testing.T) {
	mongo, err := NewMongo("config.yml", "stagingBadPort")
	if err == nil {
		t.Fatalf("NewMongo must'nt realize a success connection")
	}
	if mongo != nil {
		t.Fatalf("NewMongo must return nil")
	}
}
