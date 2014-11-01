// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"unsafe"
)

type float64SliceBind struct {
	env        *Environment
	ocibnd     *C.OCIBind
	ociNumbers []C.OCINumber
}

func (b *float64SliceBind) bindOra(values []Float64, position int, ocistmt *C.OCIStmt) error {
	float64Values := make([]float64, len(values))
	nullInds := make([]C.sb2, len(values))
	for n, _ := range values {
		if values[n].IsNull {
			nullInds[n] = C.sb2(-1)
		} else {
			float64Values[n] = values[n].Value
		}
	}
	return b.bind(float64Values, nullInds, position, ocistmt)
}

func (b *float64SliceBind) bind(values []float64, nullInds []C.sb2, position int, ocistmt *C.OCIStmt) error {
	if nullInds == nil {
		nullInds = make([]C.sb2, len(values))
	}
	alenp := make([]C.ub4, len(values))
	rcodep := make([]C.ub2, len(values))
	b.ociNumbers = make([]C.OCINumber, len(values))
	for n, _ := range values {
		alenp[n] = C.ub4(C.sizeof_OCINumber)
		r := C.OCINumberFromReal(
			b.env.ocierr,               //OCIError            *err,
			unsafe.Pointer(&values[n]), //const void          *rnum,
			8,                //uword               rnum_length,
			&b.ociNumbers[n]) //OCINumber           *number );
		if r == C.OCI_ERROR {
			return b.env.ociError()
		}
	}
	r := C.OCIBindByPos2(
		ocistmt,                          //OCIStmt      *stmtp,
		(**C.OCIBind)(&b.ocibnd),         //OCIBind      **bindpp,
		b.env.ocierr,                     //OCIError     *errhp,
		C.ub4(position),                  //ub4          position,
		unsafe.Pointer(&b.ociNumbers[0]), //void         *valuep,
		C.sb8(C.sizeof_OCINumber),        //sb8          value_sz,
		C.SQLT_VNU,                       //ub2          dty,
		unsafe.Pointer(&nullInds[0]),     //void         *indp,
		&alenp[0],                        //ub4          *alenp,
		&rcodep[0],                       //ub2          *rcodep,
		0,                                //ub4          maxarr_len,
		nil,                              //ub4          *curelep,
		C.OCI_DEFAULT)                    //ub4          mode );
	if r == C.OCI_ERROR {
		return b.env.ociError()
	}
	r = C.OCIBindArrayOfStruct(
		b.ocibnd,
		b.env.ocierr,
		C.ub4(C.sizeof_OCINumber), //ub4         pvskip,
		C.ub4(C.sizeof_sb2),       //ub4         indskip,
		C.ub4(C.sizeof_ub4),       //ub4         alskip,
		C.ub4(C.sizeof_ub2))       //ub4         rcskip
	if r == C.OCI_ERROR {
		return b.env.ociError()
	}
	return nil
}

func (b *float64SliceBind) setPtr() error {
	return nil
}

func (b *float64SliceBind) close() {
	defer func() {
		recover()
	}()
	b.ocibnd = nil
	b.env.float64SliceBindPool.Put(b)
}
