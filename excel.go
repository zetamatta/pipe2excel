package main

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/pkg/errors"
)

type Book struct {
	*ole.IDispatch
}

func NewBook() (*Book, error) {
	ole.CoInitializeEx(0, 0)

	_excel, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		return nil, errors.Wrap(err, "NewBook")
	}
	excel, err := _excel.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, errors.Wrap(err, "NewBook")
	}
	defer excel.Release()

	oleutil.PutProperty(excel, "Visible", true)

	_workbooks, err := oleutil.GetProperty(excel, "Workbooks")
	if err != nil {
		return nil, errors.Wrap(err, "NewBook")
	}
	workbooks := _workbooks.ToIDispatch()
	_workbook, err := oleutil.CallMethod(workbooks, "Add", nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewBook")
	}
	return &Book{_workbook.ToIDispatch()}, nil
}

type Sheet struct {
	*ole.IDispatch
}

type Sheets struct {
	*ole.IDispatch
}

func (b *Book) Item(name interface{}) (*Sheet, error) {
	_worksheet, err := oleutil.GetProperty(b.IDispatch, "Worksheets", name)
	if err != nil {
		return nil, errors.Wrap(err, "(*Book).Item")
	}
	return &Sheet{_worksheet.ToIDispatch()}, nil
}

func (b *Book) Sheets() (*Sheets, error) {
	_sheets, err := oleutil.GetProperty(b.IDispatch, "Sheets")
	if err != nil {
		return nil, errors.Wrap(err, "(*Book).Sheets")
	}
	return &Sheets{_sheets.ToIDispatch()}, nil
}

func (st *Sheets) Count() (int, error) {
	value, err := oleutil.GetProperty(st.IDispatch, "Count")
	if err != nil {
		return -1, errors.Wrap(err, "(*Sheets).Count")
	}
	return int(value.Val), nil
}

func (b *Book) Add() (*Sheet, error) {
	sheets, err := b.Sheets()
	if err != nil {
		return nil, errors.Wrap(err, "(*Book).Add")
	}
	defer sheets.Release()
	_worksheet, err := oleutil.CallMethod(sheets.IDispatch, "Add", nil)
	if err != nil {
		return nil, errors.Wrap(err, "(*Book).Add")
	}
	return &Sheet{_worksheet.ToIDispatch()}, nil
}

func (s *Sheet) SetName(name string) error {
	if _, err := oleutil.PutProperty(s.IDispatch, "Name", name); err != nil {
		return errors.Wrap(err, "(*Sheet).SetName")
	}
	return nil
}
