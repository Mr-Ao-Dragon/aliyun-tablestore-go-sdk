package tablestore

import (
	"math"
	"io"
	"bytes"
	"reflect"
	"github.com/golang/protobuf/proto"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/tsprotocol"
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
)

const (
	maxTableGroupNameLength = 100
	maxTableNameLength = 100
	maxPrimaryKeyLength = 255
	maxPrimaryKeyNum = 16
	maxMultiDeleteRows = 100
)

type ColumnType int32

const (
	ColumnType_STRING ColumnType = 1
	ColumnType_INTEGER ColumnType = 2
	ColumnType_BOOLEAN ColumnType = 3
	ColumnType_DOUBLE ColumnType = 4
	ColumnType_BINARY ColumnType = 5
)

const (
	Version = "1.0"
	ApiVersion = "2015-12-31"
	xOtsDateFormat = "2006-01-02T15:04:05.123Z"
	xOtsInstanceName = "x-ots-instancename"
)

type ColumnValue struct {
	Type  ColumnType
	Value interface{}
}

func (cv *ColumnValue) writeCellValue(w io.Writer) {
	writeTag(w, TAG_CELL_VALUE)
	if cv == nil {
		writeRawLittleEndian32(w, 1)
		writeRawByte(w, VT_AUTO_INCREMENT)
		return
	}

	switch cv.Type {
	case ColumnType_STRING:
		v := cv.Value.(string)

		writeRawLittleEndian32(w, int32(LITTLE_ENDIAN_32_SIZE + 1 + len(v))) // length + type + value
		writeRawByte(w, VT_STRING)
		writeRawLittleEndian32(w, int32(len(v)))
		writeBytes(w, []byte(v))

	case ColumnType_INTEGER:
		v := cv.Value.(int64)
		writeRawLittleEndian32(w, int32(LITTLE_ENDIAN_64_SIZE + 1))
		writeRawByte(w, VT_INTEGER)
		writeRawLittleEndian64(w, v)
	case ColumnType_BOOLEAN:
		v := cv.Value.(bool)
		writeRawLittleEndian32(w, 2)
		writeRawByte(w, VT_BOOLEAN)
		writeBoolean(w, v)

	case ColumnType_DOUBLE:
		v := cv.Value.(float64)

		writeRawLittleEndian32(w, LITTLE_ENDIAN_64_SIZE + 1)
		writeRawByte(w, VT_DOUBLE)
		writeDouble(w, v)

	case ColumnType_BINARY:
		v := cv.Value.([]byte)

		writeRawLittleEndian32(w, int32(LITTLE_ENDIAN_32_SIZE + 1 + len(v))) // length + type + value
		writeRawByte(w, VT_BLOB)
		writeRawLittleEndian32(w, int32(len(v)))
		writeBytes(w, v)
	}
}

func (cv *ColumnValue) writeCellValueTobytes() []byte {
	var b bytes.Buffer
	w := &b

	writeTag(w, TAG_CELL_VALUE)

	switch cv.Type {
	case ColumnType_STRING:
		v := cv.Value.(string)

		writeRawLittleEndian32(w, int32(LITTLE_ENDIAN_32_SIZE + 1 + len(v))) // length + type + value
		writeRawByte(w, VT_STRING)
		writeRawLittleEndian32(w, int32(len(v)))
		writeBytes(w, []byte(v))

	case ColumnType_INTEGER:
		v := cv.Value.(int64)
		writeRawLittleEndian32(w, int32(LITTLE_ENDIAN_64_SIZE + 1))
		writeRawByte(w, VT_INTEGER)
		writeRawLittleEndian64(w, v)
	case ColumnType_BOOLEAN:
		v := cv.Value.(bool)
		writeRawLittleEndian32(w, 2)
		writeRawByte(w, VT_BOOLEAN)
		writeBoolean(w, v)

	case ColumnType_DOUBLE:
		v := cv.Value.(float64)

		writeRawLittleEndian32(w, LITTLE_ENDIAN_64_SIZE + 1)
		writeRawByte(w, VT_DOUBLE)
		writeDouble(w, v)

	case ColumnType_BINARY:
		v := cv.Value.([]byte)

		writeRawLittleEndian32(w, int32(LITTLE_ENDIAN_32_SIZE + 1 + len(v))) // length + type + value
		writeRawByte(w, VT_BLOB)
		writeRawLittleEndian32(w, int32(len(v)))
		writeBytes(w, v)
	}
	return b.Bytes()
}

func (cv *ColumnValue) writeCellValueWithoutLengthPrefix() []byte {
	var b bytes.Buffer
	w := &b
	switch cv.Type {
	case ColumnType_STRING:
		v := cv.Value.(string)

		writeRawByte(w, VT_STRING)
		writeRawLittleEndian32(w, int32(len(v)))
		writeBytes(w, []byte(v))

	case ColumnType_INTEGER:
		v := cv.Value.(int64)
		writeRawByte(w, VT_INTEGER)
		writeRawLittleEndian64(w, v)
	case ColumnType_BOOLEAN:
		v := cv.Value.(bool)
		writeRawByte(w, VT_BOOLEAN)
		writeBoolean(w, v)

	case ColumnType_DOUBLE:
		v := cv.Value.(float64)

		writeRawByte(w, VT_DOUBLE)
		writeDouble(w, v)

	case ColumnType_BINARY:
		v := cv.Value.([]byte)

		writeRawByte(w, VT_BLOB)
		writeRawLittleEndian32(w, int32(len(v)))
		writeBytes(w, v)
	}

	return b.Bytes()
}

func (cv *ColumnValue) getCheckSum(crc byte) byte {
	if cv == nil {
		return crc8Byte(crc, VT_AUTO_INCREMENT);
	}

	switch cv.Type {
	case ColumnType_STRING:
		v := cv.Value.(string)
		crc = crc8Byte(crc, VT_STRING)
		crc = crc8Int32(crc, int32(len(v)))
		crc = crc8Bytes(crc, []byte(v))
	case ColumnType_INTEGER:
		v := cv.Value.(int64)
		crc = crc8Byte(crc, VT_INTEGER)
		crc = crc8Int64(crc, v)
	case ColumnType_BOOLEAN:
		v := cv.Value.(bool)
		crc = crc8Byte(crc, VT_BOOLEAN)
		if v {
			crc = crc8Byte(crc, 0x1)
		} else {
			crc = crc8Byte(crc, 0x0)
		}

	case ColumnType_DOUBLE:
		v := cv.Value.(float64)
		crc = crc8Byte(crc, VT_DOUBLE)
		crc = crc8Int64(crc, int64(math.Float64bits(v)))
	case ColumnType_BINARY:
		v := cv.Value.([]byte)
		crc = crc8Byte(crc, VT_BLOB)
		crc = crc8Int32(crc, int32(len(v)))
		crc = crc8Bytes(crc, v)
	}

	return crc
}

type Column struct {
	Name         []byte
	Value        ColumnValue
	Type         byte
	Timestamp    int64
	HasType      bool
	HasTimestamp bool
	IgnoreValue  bool
}

func NewColumn(name []byte, value interface{}) *Column {

	v := &Column{}
	v.Name = name

	if value != nil {
		t := reflect.TypeOf(value)
		switch t.Kind() {
		case reflect.String:
			v.Value.Type = ColumnType_STRING

		case reflect.Int64:
			v.Value.Type = ColumnType_INTEGER

		case reflect.Bool:
			v.Value.Type = ColumnType_BOOLEAN

		case reflect.Float64:
			v.Value.Type = ColumnType_DOUBLE

		default:
			return nil
		}

		v.Value.Value = value
	}

	return v
}

func (c *Column) toPlainBufferCell(ignoreValue bool) *PlainBufferCell {
	cell := &PlainBufferCell{}
	cell.cellName = c.Name
	cell.ignoreValue = ignoreValue
	if ignoreValue == false {
		cell.cellValue = &c.Value
	}

	if c.HasType {
		cell.hasCellType = c.HasType
		cell.cellType = byte(c.Type)
	}

	if c.HasTimestamp {
		cell.hasCellTimestamp = c.HasTimestamp
		cell.cellTimestamp = c.Timestamp
	}

	return cell
}

type PrimaryKeyColumnInner struct {
	Name  []byte
	Type  tsprotocol.PrimaryKeyType
	Value interface{}
}

func NewPrimaryKeyColumnINF_MAX(name []byte) *PrimaryKeyColumnInner {
	v := &PrimaryKeyColumnInner{}
	v.Name = name
	v.Type = 0
	v.Value = "INF_MAX"

	return v
}

func NewPrimaryKeyColumnINF_MIN(name []byte) *PrimaryKeyColumnInner {
	v := &PrimaryKeyColumnInner{}
	v.Name = name
	v.Type = 0
	v.Value = "INF_MIN"

	return v
}

func NewPrimaryKeyColumnAuto_Increment(name []byte) *PrimaryKeyColumnInner {
	v := &PrimaryKeyColumnInner{}
	v.Name = name
	v.Type = 0
	v.Value = "AUTO_INCRMENT"
	return v
}

func NewPrimaryKeyColumn(name []byte, value interface{}, option PrimaryKeyOption) *PrimaryKeyColumnInner {

	if option == NONE {
		v := &PrimaryKeyColumnInner{}
		v.Name = name

		t := reflect.TypeOf(value)
		switch t.Kind() {
		case reflect.String:
			v.Type = tsprotocol.PrimaryKeyType_STRING

		case reflect.Int64:
			v.Type = tsprotocol.PrimaryKeyType_INTEGER

		default:
			return nil
		}

		v.Value = value

		return v
	} else if option == AUTO_INCREMENT {
		return NewPrimaryKeyColumnAuto_Increment(name)
	} else if option == MIN {
		return NewPrimaryKeyColumnINF_MIN(name)
	} else {
		return NewPrimaryKeyColumnINF_MAX(name)
	}
}

func (pkc *PrimaryKeyColumnInner) toColumnValue() *ColumnValue {
	switch pkc.Type {
	case tsprotocol.PrimaryKeyType_INTEGER:
		return &ColumnValue{ColumnType_INTEGER, pkc.Value}
	case tsprotocol.PrimaryKeyType_STRING:
		return &ColumnValue{ColumnType_STRING, pkc.Value}
	case tsprotocol.PrimaryKeyType_BINARY:
		return &ColumnValue{ColumnType_BINARY, pkc.Value}
	}

	return nil
}

func (pkc *PrimaryKeyColumnInner) toPlainBufferCell() *PlainBufferCell {
	cell := &PlainBufferCell{}
	cell.cellName = pkc.Name
	cell.cellValue = pkc.toColumnValue()
	return cell
}

func (pkc *PrimaryKeyColumnInner) isInfMin() bool {
	if pkc.Type == 0 && pkc.Value.(string) == "INF_MIN" {
		return true
	}

	return false
}

func (pkc *PrimaryKeyColumnInner) isInfMax() bool {
	if pkc.Type == 0 && pkc.Value.(string) == "INF_MAX" {
		return true
	}

	return false
}

func (pkc *PrimaryKeyColumnInner) isAutoInc() bool {
	if pkc.Type == 0 && pkc.Value.(string) == "AUTO_INCRMENT" {
		return true
	}
	return false
}

func (pkc *PrimaryKeyColumnInner) getCheckSum(crc byte) byte {
	if pkc.isInfMin() {
		return crc8Byte(crc, VT_INF_MIN)
	}
	if pkc.isInfMax() {
		return crc8Byte(crc, VT_INF_MAX)
	}
	if pkc.isAutoInc() {
		return crc8Byte(crc, VT_AUTO_INCREMENT)
	}

	return pkc.toColumnValue().getCheckSum(crc)
}

func (pkc *PrimaryKeyColumnInner) writePrimaryKeyColumn(w io.Writer) {
	writeTag(w, TAG_CELL)
	writeCellName(w, []byte(pkc.Name))
	if pkc.isInfMin() {
		writeTag(w, TAG_CELL_VALUE)
		writeRawLittleEndian32(w, 1)
		writeRawByte(w, VT_INF_MIN)
		return
	}
	if pkc.isInfMax() {
		writeTag(w, TAG_CELL_VALUE)
		writeRawLittleEndian32(w, 1)
		writeRawByte(w, VT_INF_MAX)
		return
	}
	if pkc.isAutoInc() {
		writeTag(w, TAG_CELL_VALUE)
		writeRawLittleEndian32(w, 1)
		writeRawByte(w, VT_AUTO_INCREMENT)
		return
	}

	pkc.toColumnValue().writeCellValue(w)
}

type PrimaryKey2 struct {
	primaryKey []*PrimaryKeyColumnInner
}

func (pk *PrimaryKey) Build(isDelete bool) []byte {
	var b bytes.Buffer
	writeHeader(&b)
	writeTag(&b, TAG_ROW_PK)

	rowChecksum := byte(0x0)
	var cellChecksum byte

	for _, column := range (pk.PrimaryKeys) {
		primaryKeyColumn := NewPrimaryKeyColumn([]byte(column.ColumnName), column.Value, column.PrimaryKeyOption)

		cellChecksum = crc8Bytes(byte(0x0), []byte(primaryKeyColumn.Name))
		cellChecksum = primaryKeyColumn.getCheckSum(cellChecksum)
		rowChecksum = crc8Byte(rowChecksum, cellChecksum)
		primaryKeyColumn.writePrimaryKeyColumn(&b)

		writeTag(&b, TAG_CELL_CHECKSUM)
		writeRawByte(&b, cellChecksum)
	}

	// 没有deleteMarker, 要与0x0做crc.
	if isDelete {
		writeTag(&b, TAG_DELETE_ROW_MARKER)
		rowChecksum = crc8Byte(rowChecksum, byte(0x1))
	} else {
		rowChecksum = crc8Byte(rowChecksum, byte(0x0))
	}
	writeTag(&b, TAG_ROW_CHECKSUM)
	writeRawByte(&b, rowChecksum)

	return b.Bytes()
}

type RowPutChange struct {
	primaryKey   []*PrimaryKeyColumnInner
	columnsToPut []*Column
}

type RowUpdateChange struct {
	primaryKey      []*PrimaryKeyColumnInner
	columnsToUpdate []*Column
}

func (rpc *RowPutChange) Build() []byte {
	pkCells := make([]*PlainBufferCell, len(rpc.primaryKey))
	for i, pkc := range (rpc.primaryKey) {
		pkCells[i] = pkc.toPlainBufferCell()
	}

	cells := make([]*PlainBufferCell, len(rpc.columnsToPut))
	for i, c := range (rpc.columnsToPut) {
		cells[i] = c.toPlainBufferCell(false)
	}

	row := &PlainBufferRow{pkCells, cells, false}
	var b bytes.Buffer
	row.writeRowWithHeader(&b)

	return b.Bytes()
}

func (ruc *RowUpdateChange) Build() []byte {
	pkCells := make([]*PlainBufferCell, len(ruc.primaryKey))
	for i, pkc := range (ruc.primaryKey) {
		pkCells[i] = pkc.toPlainBufferCell()
	}

	cells := make([]*PlainBufferCell, len(ruc.columnsToUpdate))
	for i, c := range (ruc.columnsToUpdate) {
		cells[i] = c.toPlainBufferCell(c.IgnoreValue)
	}

	row := &PlainBufferRow{pkCells, cells, false}
	var b bytes.Buffer
	row.writeRowWithHeader(&b)

	return b.Bytes()
}

const (
	MaxValue = "_get_range_max"
	MinValue = "_get_range_min"
)

func (comparatorType *ComparatorType) ConvertToPbComparatorType() tsprotocol.ComparatorType {
	switch *comparatorType {
	case CT_EQUAL:
		return tsprotocol.ComparatorType_CT_EQUAL
	case CT_NOT_EQUAL:
		return tsprotocol.ComparatorType_CT_NOT_EQUAL
	case CT_GREATER_THAN:
		return tsprotocol.ComparatorType_CT_GREATER_THAN
	case CT_GREATER_EQUAL:
		return tsprotocol.ComparatorType_CT_GREATER_EQUAL
	case CT_LESS_THAN:
		return tsprotocol.ComparatorType_CT_LESS_THAN
	default:
		return tsprotocol.ComparatorType_CT_LESS_EQUAL
	}
}

func (loType *LogicalOperator) ConvertToPbLoType() tsprotocol.LogicalOperator {
	switch *loType {
	case LO_NOT:
		return tsprotocol.LogicalOperator_LO_NOT
	case LO_AND:
		return tsprotocol.LogicalOperator_LO_AND
	default:
		return tsprotocol.LogicalOperator_LO_OR
	}
}

func (ft *FilterType) ConvertFilterType() tsprotocol.FilterType {
	switch *ft {
	case FT_SINGLE_COLUMN_VALUE:
		return tsprotocol.FilterType_FT_SINGLE_COLUMN_VALUE
	case FT_COMPOSITE_COLUMN_VALUE:
		return tsprotocol.FilterType_FT_COMPOSITE_COLUMN_VALUE
	default:
		return tsprotocol.FilterType_FT_COLUMN_PAGINATION
	}
}

func NewSingleColumnValueFilter(condition *SingleColumnCondition) *tsprotocol.SingleColumnValueFilter {
	filter := new(tsprotocol.SingleColumnValueFilter)

	comparatorType := condition.Comparator.ConvertToPbComparatorType()
	filter.Comparator = &comparatorType
	filter.ColumnName = condition.ColumnName
	col := NewColumn([]byte(*condition.ColumnName), condition.ColumnValue)
	filter.ColumnValue = col.toPlainBufferCell(false).cellValue.writeCellValueWithoutLengthPrefix() //writeCellValueTobytes()  //
	filter.FilterIfMissing = proto.Bool(condition.FilterIfMissing)
	filter.LatestVersionOnly = proto.Bool(condition.LatestVersionOnly)

	return filter
}

func NewCompositeFilter(filters []ColumnFilter, lo LogicalOperator) *tsprotocol.CompositeColumnValueFilter {
	ccvfilter := new(tsprotocol.CompositeColumnValueFilter)
	combinator := lo.ConvertToPbLoType()
	ccvfilter.Combinator = &combinator
	for _, cf := range (filters) {
		filter := cf.ToFilter()
		ccvfilter.SubFilters = append(ccvfilter.SubFilters, filter)
	}

	return ccvfilter
}

func NewPaginationFilter(filter *PaginationFilter) *tsprotocol.ColumnPaginationFilter {
	pageFilter := new(tsprotocol.ColumnPaginationFilter)
	pageFilter.Offset = proto.Int32(filter.Offset)
	pageFilter.Limit = proto.Int32(filter.Limit)
	return pageFilter
}

func getTableStoreDefaultConfig() *TSConfig {
	httpTimeout := &HTTPTimeout{ConnectionTimeout:time.Second * 15, RequestTimeout :time.Second * 30  }
	config := &TSConfig{RetryTimes: 10, HTTPTimeout: *httpTimeout, MaxRetryTime: time.Second * 5 }
	return config
}

func (otsClient *TableStoreClient) postReq(req *http.Request, url string) (body []byte, err error, statusCode int) {
	resp, err := otsClient.httpClient.Do(req)
	if err != nil {
		return nil, err, resp.StatusCode
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, resp.StatusCode
	}

	if (resp.StatusCode >= 200 && resp.StatusCode < 300) == false {
		return body, fmt.Errorf("get %s response status is %d", url, resp.StatusCode), resp.StatusCode
	}

	return body, nil, resp.StatusCode
}

func buildRowPutChange(primarykey *PrimaryKey, columns []DataColumn) *RowPutChange {
	row := new(RowPutChange)
	row.primaryKey = make([]*PrimaryKeyColumnInner, len(primarykey.PrimaryKeys))
	for i, p := range (primarykey.PrimaryKeys) {
		row.primaryKey[i] = NewPrimaryKeyColumn([]byte(p.ColumnName), p.Value, p.PrimaryKeyOption)
	}

	row.columnsToPut = make([]*Column, len(columns))
	for i, p := range (columns) {
		row.columnsToPut[i] = NewColumn([]byte(p.ColumnName), p.Value)
		if p.Timestamp != 0 {
			row.columnsToPut[i].HasTimestamp = true
			row.columnsToPut[i].Timestamp = p.Timestamp
		}
	}

	return row
}

func buildRowUpdateChange(primarykey *PrimaryKey, columns []ColumnToUpdate) *RowUpdateChange {
	row := new(RowUpdateChange)
	row.primaryKey = make([]*PrimaryKeyColumnInner, len(primarykey.PrimaryKeys))
	for i, p := range (primarykey.PrimaryKeys) {
		row.primaryKey[i] = NewPrimaryKeyColumn([]byte(p.ColumnName), p.Value, p.PrimaryKeyOption)
	}

	row.columnsToUpdate = make([]*Column, len(columns))
	for i, p := range (columns) {
		row.columnsToUpdate[i] = NewColumn([]byte(p.ColumnName), p.Value)
		row.columnsToUpdate[i].HasTimestamp = p.HasTimestamp
		row.columnsToUpdate[i].HasType = p.HasType
		row.columnsToUpdate[i].Type = p.Type
		row.columnsToUpdate[i].Timestamp = p.Timestamp
		row.columnsToUpdate[i].IgnoreValue = p.IgnoreValue
	}

	return row
}

func (condition *RowCondition) buildCondition() *tsprotocol.RowExistenceExpectation {
	switch condition.RowExistenceExpectation {
	case RowExistenceExpectation_IGNORE:
		return tsprotocol.RowExistenceExpectation_IGNORE.Enum()
	case RowExistenceExpectation_EXPECT_EXIST:
		return tsprotocol.RowExistenceExpectation_EXPECT_EXIST.Enum()
	case RowExistenceExpectation_EXPECT_NOT_EXIST:
		return tsprotocol.RowExistenceExpectation_EXPECT_NOT_EXIST.Enum()
	}

	// Todo : refine the error
	panic("Invalid input")
}

// build primary key for create table, put row, delete row and update row
func buildPrimaryKey(primaryKeyName string, value interface{}) *PrimaryKeyColumn {
	// Todo: validate the input
	return &PrimaryKeyColumn{ColumnName: primaryKeyName, Value:value, PrimaryKeyOption: NONE}
}

// Todo: consider block user pass any type
func (rowchange *PutRowChange) AddColumn(columnName string, value interface{}) {
	// Todo: validate the input
	column := &DataColumn{ColumnName: columnName, Value:value}
	rowchange.Columns = append(rowchange.Columns, *column)
}

func (rowchange *PutRowChange) AddColumnWithTimestamp(columnName string, value interface{}, timestamp int64) {
	// Todo: validate the input
	column := &DataColumn{ColumnName: columnName, Value:value}
	column.Timestamp = timestamp
	rowchange.Columns = append(rowchange.Columns, *column)
}

func (pk *PrimaryKey) AddPrimaryKeyColumn(primaryKeyName string, value interface{}) {
	pk.PrimaryKeys = append(pk.PrimaryKeys, buildPrimaryKey(primaryKeyName, value))
}

func (pk *PrimaryKey) AddPrimaryKeyColumnWithAutoIncrement(primaryKeyName string) {
	pk.PrimaryKeys = append(pk.PrimaryKeys, &PrimaryKeyColumn{ColumnName: primaryKeyName, PrimaryKeyOption: AUTO_INCREMENT })
}

func (pk *PrimaryKey) AddPrimaryKeyColumnWithMinValue(primaryKeyName string) {
	pk.PrimaryKeys = append(pk.PrimaryKeys, &PrimaryKeyColumn{ColumnName: primaryKeyName, PrimaryKeyOption: MIN })
}

// Only used for range query
func (pk *PrimaryKey) AddPrimaryKeyColumnWithMaxValue(primaryKeyName string) {
	pk.PrimaryKeys = append(pk.PrimaryKeys, &PrimaryKeyColumn{ColumnName: primaryKeyName, PrimaryKeyOption: MAX })
}

func (rowchange *PutRowChange) SetCondition(rowExistenceExpectation RowExistenceExpectation) {
	rowchange.Condition = &RowCondition{RowExistenceExpectation:rowExistenceExpectation}
}

func (rowchange *DeleteRowChange) SetCondition(rowExistenceExpectation RowExistenceExpectation) {
	rowchange.Condition = &RowCondition{RowExistenceExpectation:rowExistenceExpectation}
}

func (Criteria *SingleRowQueryCriteria) SetFilter(filter ColumnFilter) {
	Criteria.Filter = filter
}

func NewSingleColumnCondition(columnName string, comparator ComparatorType, value interface{}) *SingleColumnCondition {
	return &SingleColumnCondition{ColumnName: &columnName, Comparator: &comparator, ColumnValue:value}
}

func NewCompositeColumnCondition(lo LogicalOperator) *CompositeColumnValueFilter {
	return &CompositeColumnValueFilter{Operator: lo}
}

func (rowchange *PutRowChange) SetColumnCondition(condition ColumnFilter) {
	rowchange.Condition.ColumnCondition = condition
}

func (rowchange *UpdateRowChange) SetCondition(rowExistenceExpectation RowExistenceExpectation) {
	rowchange.Condition = &RowCondition{RowExistenceExpectation:rowExistenceExpectation}
}

func (rowchange *UpdateRowChange) SetColumnCondition(condition ColumnFilter) {
	rowchange.Condition.ColumnCondition = condition
}

func (rowchange *DeleteRowChange) SetColumnCondition(condition ColumnFilter) {
	rowchange.Condition.ColumnCondition = condition
}

func (meta *TableMeta) AddPrimaryKeyColumn(name string, keyType PrimaryKeyType) {
	meta.SchemaEntry = append(meta.SchemaEntry, &PrimaryKeySchema{Name:&name, Type: &keyType})
}

func (meta *TableMeta) AddPrimaryKeyColumnOption(name string, keyType PrimaryKeyType, keyOption PrimaryKeyOption) {
	meta.SchemaEntry = append(meta.SchemaEntry, &PrimaryKeySchema{Name:&name, Type: &keyType, Option: &keyOption})
}

func (rowchange *UpdateRowChange) PutColumn(columnName string, value interface{}) {
	// Todo: validate the input
	column := &ColumnToUpdate{ColumnName: columnName, Value: value}
	rowchange.Columns = append(rowchange.Columns, *column)
}

func (rowchange *UpdateRowChange) DeleteColumn(columnName string) {
	// Todo: validate the input
	column := &ColumnToUpdate{ColumnName: columnName, Value: nil, Type: DELETE_ALL_VERSION, HasType: true, IgnoreValue: true}
	rowchange.Columns = append(rowchange.Columns, *column)
}

func (rowchange *UpdateRowChange) DeleteColumnWithTimestamp(columnName string, timestamp int64) {
	// Todo: validate the input
	column := &ColumnToUpdate{ColumnName: columnName, Value: nil, Type: DELETE_ONE_VERSION, HasType: true, HasTimestamp: true, Timestamp: timestamp, IgnoreValue: true}
	rowchange.Columns = append(rowchange.Columns, *column)
}

func (rowchange *DeleteRowChange) Serialize() []byte {
	return rowchange.PrimaryKey.Build(true)
}

func (rowchange *PutRowChange) Serialize() []byte {
	row := buildRowPutChange(rowchange.PrimaryKey, rowchange.Columns)
	return row.Build()
}

func (rowchange *UpdateRowChange) Serialize() []byte {
	row := buildRowUpdateChange(rowchange.PrimaryKey, rowchange.Columns)
	return row.Build()
}

func (rowchange *DeleteRowChange) GetTableName() string {
	return rowchange.TableName
}

func (rowchange *PutRowChange) GetTableName() string {
	return rowchange.TableName
}

func (rowchange *UpdateRowChange) GetTableName() string {
	return rowchange.TableName
}

func (rowchange *DeleteRowChange) getOperationType() tsprotocol.OperationType {
	return tsprotocol.OperationType_DELETE
}

func (rowchange *PutRowChange) getOperationType() tsprotocol.OperationType {
	return tsprotocol.OperationType_PUT
}

func (rowchange *UpdateRowChange) getOperationType() tsprotocol.OperationType {
	return tsprotocol.OperationType_UPDATE
}

func (rowchange *DeleteRowChange) getCondition() *tsprotocol.Condition {
	condition := new(tsprotocol.Condition)
	condition.RowExistence = rowchange.Condition.buildCondition()
	if rowchange.Condition.ColumnCondition != nil {
		condition.ColumnCondition = rowchange.Condition.ColumnCondition.Serialize()
	}
	return condition
}

func (rowchange *UpdateRowChange) getCondition() *tsprotocol.Condition {
	condition := new(tsprotocol.Condition)
	condition.RowExistence = rowchange.Condition.buildCondition()
	if rowchange.Condition.ColumnCondition != nil {
		condition.ColumnCondition = rowchange.Condition.ColumnCondition.Serialize()
	}
	return condition
}

func (rowchange *PutRowChange) getCondition() *tsprotocol.Condition {
	condition := new(tsprotocol.Condition)
	condition.RowExistence = rowchange.Condition.buildCondition()
	if rowchange.Condition.ColumnCondition != nil {
		condition.ColumnCondition = rowchange.Condition.ColumnCondition.Serialize()
	}
	return condition
}

func (request *BatchWriteRowRequest) AddRowChange(change RowChange) {
	if request.RowChangesGroupByTable == nil {
		request.RowChangesGroupByTable = make(map[string][]RowChange)
	}
	request.RowChangesGroupByTable[change.GetTableName()] = append(request.RowChangesGroupByTable[change.GetTableName()], change)
}

func (direction Direction) ToDirection() tsprotocol.Direction {
	if direction == FORWARD {
		return tsprotocol.Direction_FORWARD
	} else {
		return tsprotocol.Direction_BACKWARD
	}
}