package pg

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/amahdian/golang-gin-boilerplate/domain/model/common"
	"github.com/amahdian/golang-gin-boilerplate/global"
	"github.com/amahdian/golang-gin-boilerplate/global/errs"
	"github.com/amahdian/golang-gin-boilerplate/pkg/logger"
	"github.com/gertd/go-pluralize"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var pluralizer = pluralize.NewClient()

type crudStg[M schema.Tabler] struct {
	db *gorm.DB

	tableName          string
	columnNames        []string
	tagToColumnNameMap map[string]string
}

func (stg *crudStg[M]) CreateOne(model M) error {
	err := stg.db.Create(model).Error
	return err
}

func (stg *crudStg[M]) CreateMany(models []M) error {
	if len(models) == 0 {
		return nil
	}
	err := stg.db.Create(&models).Error
	return err
}

func (stg *crudStg[M]) FindById(id int64) (model M, err error) {
	err = stg.db.First(&model, id).Error
	if err != nil {
		tableName := stg.getTableName()
		entryName := pluralizer.Singular(tableName)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			err = errs.Newf(errs.NotFound, nil, errs.DefaultEntryNotFoundMessage, entryName, id)
		default:
			err = errs.Wrapf(err, "failed to get %s", tableName)
		}
		return
	}

	return
}

func (stg *crudStg[M]) ListByIds(ids []int64) (models []M, err error) {
	if len(ids) == 0 {
		return make([]M, 0), nil
	}

	err = stg.db.Where("id in ?", ids).Find(&models).Error
	return
}

func (stg *crudStg[M]) UpdateOne(model M, saveAssociations bool) error {
	query := stg.db
	if saveAssociations {
		query = query.Session(&gorm.Session{FullSaveAssociations: true})
	}
	err := query.Debug().Save(model).Error
	return err
}

func (stg *crudStg[M]) UpdatePartial(model M, returnUpdated bool) error {
	query := stg.db.Model(model)
	if returnUpdated {
		query.Clauses(clause.Returning{})
	}
	err := query.Updates(model).Error
	return err
}

func (stg *crudStg[M]) UpdateMany(models []M) error {
	if len(models) == 0 {
		return nil
	}
	err := stg.db.Save(&models).Error
	return err
}

func (stg *crudStg[M]) ExistsById(id int64) (exists bool, err error) {
	_, err = stg.FindById(id)

	exists = true

	if err != nil {
		exists = false
		if errors.As(err, &errs.EntryNotFoundErr{}) {
			// Do not propagate the record not found error
			err = nil
		}
	}

	return
}

func (stg *crudStg[M]) DeleteById(id int64) error {
	var model M
	db := stg.db.Delete(&model, id)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected < 1 {
		tableName := stg.getTableName()
		entryName := pluralizer.Singular(tableName)
		return errs.NewEntryNotFoundErr(fmt.Sprintf(errs.DefaultEntryNotFoundMessage, entryName, id))
	}
	return nil
}

func (stg *crudStg[M]) DeleteByIds(ids []int64) error {
	var model M
	return stg.db.Where("id in ?", ids).Delete(&model).Error
}

func (stg *crudStg[M]) ListAll() (models []M, err error) {
	err = stg.db.Find(&models).Error
	return
}

func (stg *crudStg[M]) withPagination(pagination *common.Pagination, tableAlias ...string) gormScope {
	return func(db *gorm.DB) *gorm.DB {
		// disable pagination for internal API calls
		if pagination == common.InternalPagination() {
			return db
		}

		prefix := getTablePrefix(tableAlias...)

		if pagination == nil {
			logger.Debug("withPagination scope is called but no pagination info is provided")
			return db
		}

		// count the result in a separate query session to prevent polluting the original query
		db.Session(&gorm.Session{}).Count(&pagination.TotalCount)
		if pagination.OrderBy != "" {
			fieldName := pagination.OrderBy
			tagToColumnNameMap := stg.getTagToColumnNameMap()
			columnName, ok := tagToColumnNameMap[fieldName]
			if !ok {
				db.AddError(fmt.Errorf("invalid field %q for search condition", fieldName))
				return db
			}
			db.Order(fmt.Sprintf("%s%s %s", prefix, columnName, pagination.Order))
		}

		db.Limit(pagination.PageSize)

		if pagination.Page != 0 {
			offset := pagination.PageSize * pagination.Page
			db.Offset(offset)
		}

		return db
	}
}

func (stg *crudStg[M]) withSearchFilters(filters []*common.FieldFilter, tableAlias ...string) gormScope {
	return func(db *gorm.DB) *gorm.DB {
		prefix := getTablePrefix(tableAlias...)

		tagToColumnNamesMap := stg.getTagToColumnNameMap()

		for _, v := range filters {
			fieldName := v.FieldName
			columnName, ok := tagToColumnNamesMap[fieldName]
			if !ok {
				db.AddError(errs.NewInvalidSearchFieldErr(fieldName))
				return db
			}
			prefixedColumnName := prefix + columnName

			formattedColumnName := fmt.Sprintf("%s::TEXT", prefixedColumnName)
			if strings.Contains(columnName, global.DateColumnPostfix) {
				formattedColumnName = fmt.Sprintf("TO_CHAR(TO_TIMESTAMP(%s)::timestamp with time zone at time zone 'UTC', '%s')", prefixedColumnName, global.DateColumnFormat)
			}

			switch v.Condition {
			case common.SearchConditionContains:
				db.Where(fmt.Sprintf("%s ILIKE '%s'", formattedColumnName, "%"+v.Value+"%"))
			case common.SearchConditionEqual:
				db.Where(fmt.Sprintf("%s = '%s'", formattedColumnName, v.Value))
			case common.SearchConditionNotEqual:
				db.Not(fmt.Sprintf("%s = '%s'", formattedColumnName, v.Value))
			default:
				db.AddError(fmt.Errorf("unsupported search operation %q for field %q", v.Condition, fieldName))
				return db
			}
		}
		return db
	}
}

func (stg *crudStg[M]) withSearch(params *common.SearchParams, tableAlias ...string) gormScope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(
			stg.withSearchFilters(params.Filters, tableAlias...),
			stg.withPagination(params.Pagination, tableAlias...),
		)
	}
}

func (stg *crudStg[M]) updateElements(elements any) error {
	if reflect.TypeOf(elements).Kind() != reflect.Slice {
		return errs.Newf(errs.Internal, nil, "the provided argument is not an array")
	}

	if reflect.ValueOf(elements).Len() == 0 {
		return nil
	}

	tableName := stg.getTableName()

	var model M
	modelType := reflect.TypeOf(model)
	elementType := reflect.TypeOf(elements).Elem()

	if elementType == modelType {
		// since elementType and modelType are the same we don't need to do anything fancy. leave the rest to gorm
		err := stg.db.Save(elements).Error
		return err
	} else {
		// elementType and modelType differ meaning the elementType is probably a subset of the modelType.
		// therefore we only retrieve the columns needed by the elementType.
		elementInstance := reflect.Zero(elementType).Interface()
		elementSchema := getGormSchema(elementInstance)

		queryColumnNames := lo.Map(elementSchema.Fields, func(f *schema.Field, _ int) string {
			return f.DBName
		})
		if len(queryColumnNames) == 0 {
			return fmt.Errorf("could not find any column to query from %q table", tableName)
		}

		// must contain "id" column name
		if !lo.Contains(queryColumnNames, "id") {
			return fmt.Errorf("cannot perform batch update on table %q without %q column", tableName, "id")
		}

		modelColumnNames := stg.getColumnNames()
		if !lo.Every(modelColumnNames, queryColumnNames) {
			invalidColumnNames := lo.Without(queryColumnNames, modelColumnNames...)
			return fmt.Errorf("cannot perform batch update on table %q because of invalid columns: %q", tableName, strings.Join(invalidColumnNames, ", "))
		}

		err := stg.db.
			Table(tableName).
			Select(queryColumnNames).
			Save(elements).
			Error
		return err
	}
}

func (stg *crudStg[M]) getColumnNames() []string {
	if stg.columnNames != nil {
		return stg.columnNames
	}

	var model M
	s := getGormSchema(&model)

	columnNames := make([]string, 0)
	for _, field := range s.Fields {
		columnName := field.DBName
		columnNames = append(columnNames, columnName)
	}

	stg.columnNames = columnNames
	return columnNames
}

func (stg *crudStg[M]) getTagToColumnNameMap() map[string]string {
	if stg.tagToColumnNameMap != nil {
		return stg.tagToColumnNameMap
	}

	var model M
	s := getGormSchema(&model)

	// map json tag to column name
	m := make(map[string]string)
	for _, field := range s.Fields {
		columnName := field.DBName
		jsonTag := field.Tag.Get("json")
		jsonName := strings.Split(jsonTag, ",")[0]
		m[jsonName] = columnName
		// in some requests FE have sent the column name instead of json name for avoiding to err I added column name also
		m[columnName] = columnName
	}

	stg.tagToColumnNameMap = m
	return m
}

func (stg *crudStg[M]) getTableName() string {
	if stg.tableName != "" {
		return stg.tableName
	}

	// M is a pointer to the model struct
	var model M
	v := reflect.New(reflect.TypeOf(model).Elem())
	// the Tabler might be implemented pointer receiver

	if tabler, ok := v.Interface().(schema.Tabler); ok {
		stg.tableName = tabler.TableName()
	} else if tabler, ok = v.Elem().Interface().(schema.Tabler); ok {
		stg.tableName = tabler.TableName()
	}

	if stg.tableName == "" {
		panic(fmt.Sprintf("the model %v does not implement the tabler interface", model))
	}

	return stg.tableName
}
