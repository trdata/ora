// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
#include "version.h"
*/
import "C"
import (
	"unsafe"
)

type bndUint8 struct {
	stmt      *Stmt
	ocibnd    *C.OCIBind
	ociNumber [1]C.OCINumber
}

func (bnd *bndUint8) bind(value uint8, position int, stmt *Stmt) error {
	bnd.stmt = stmt
	r := C.OCINumberFromInt(
		bnd.stmt.ses.srv.env.ocierr, //OCIError            *err,
		unsafe.Pointer(&value),      //const void          *inum,
		1, //uword               inum_length,
		C.OCI_NUMBER_UNSIGNED, //uword               inum_s_flag,
		&bnd.ociNumber[0])     //OCINumber           *number );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}
	r = C.OCIBINDBYPOS(
		bnd.stmt.ocistmt, //OCIStmt      *stmtp,
		&bnd.ocibnd,
		bnd.stmt.ses.srv.env.ocierr,       //OCIError     *errhp,
		C.ub4(position),                   //ub4          position,
		unsafe.Pointer(&bnd.ociNumber[0]), //void         *valuep,
		C.LENGTH_TYPE(C.sizeof_OCINumber), //sb8          value_sz,
		C.SQLT_VNU,                        //ub2          dty,
		nil,                               //void         *indp,
		nil,                               //ub2          *alenp,
		nil,                               //ub2          *rcodep,
		0,                                 //ub4          maxarr_len,
		nil,                               //ub4          *curelep,
		C.OCI_DEFAULT)                     //ub4          mode );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}
	return nil
}

func (bnd *bndUint8) setPtr() error {
	return nil
}

func (bnd *bndUint8) close() (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = errR(value)
		}
	}()

	stmt := bnd.stmt
	bnd.stmt = nil
	bnd.ocibnd = nil
	stmt.putBnd(bndIdxUint8, bnd)
	return nil
}
