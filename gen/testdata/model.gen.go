package testdata

import (
	sqlx "database/sql"

	easyorm "github.com/JrMarcco/easy-orm"
)

const (
	ModelId       = "Id"
	ModelAge      = "Age"
	ModelUsername = "Username"
	ModelAddress  = "Address"
)

func ModelIdEq(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Eq(val)
}

func ModelIdNe(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Ne(val)
}

func ModelIdGt(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Gt(val)
}

func ModelIdGe(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Ge(val)
}

func ModelIdLt(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Lt(val)
}

func ModelIdLe(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Le(val)
}

func ModelAgeEq(val *int32) easyorm.Predicate {
	return easyorm.Col("Age").Eq(val)
}

func ModelAgeNe(val *int32) easyorm.Predicate {
	return easyorm.Col("Age").Ne(val)
}

func ModelAgeGt(val *int32) easyorm.Predicate {
	return easyorm.Col("Age").Gt(val)
}

func ModelAgeGe(val *int32) easyorm.Predicate {
	return easyorm.Col("Age").Ge(val)
}

func ModelAgeLt(val *int32) easyorm.Predicate {
	return easyorm.Col("Age").Lt(val)
}

func ModelAgeLe(val *int32) easyorm.Predicate {
	return easyorm.Col("Age").Le(val)
}

func ModelUsernameEq(val string) easyorm.Predicate {
	return easyorm.Col("Username").Eq(val)
}

func ModelUsernameNe(val string) easyorm.Predicate {
	return easyorm.Col("Username").Ne(val)
}

func ModelUsernameGt(val string) easyorm.Predicate {
	return easyorm.Col("Username").Gt(val)
}

func ModelUsernameGe(val string) easyorm.Predicate {
	return easyorm.Col("Username").Ge(val)
}

func ModelUsernameLt(val string) easyorm.Predicate {
	return easyorm.Col("Username").Lt(val)
}

func ModelUsernameLe(val string) easyorm.Predicate {
	return easyorm.Col("Username").Le(val)
}

func ModelAddressEq(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Address").Eq(val)
}

func ModelAddressNe(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Address").Ne(val)
}

func ModelAddressGt(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Address").Gt(val)
}

func ModelAddressGe(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Address").Ge(val)
}

func ModelAddressLt(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Address").Lt(val)
}

func ModelAddressLe(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Address").Le(val)
}

const (
	SubModelId      = "Id"
	SubModelName    = "Name"
	SubModelEmail   = "Email"
	SubModelBalance = "Balance"
)

func SubModelIdEq(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Eq(val)
}

func SubModelIdNe(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Ne(val)
}

func SubModelIdGt(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Gt(val)
}

func SubModelIdGe(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Ge(val)
}

func SubModelIdLt(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Lt(val)
}

func SubModelIdLe(val uint64) easyorm.Predicate {
	return easyorm.Col("Id").Le(val)
}

func SubModelNameEq(val string) easyorm.Predicate {
	return easyorm.Col("Name").Eq(val)
}

func SubModelNameNe(val string) easyorm.Predicate {
	return easyorm.Col("Name").Ne(val)
}

func SubModelNameGt(val string) easyorm.Predicate {
	return easyorm.Col("Name").Gt(val)
}

func SubModelNameGe(val string) easyorm.Predicate {
	return easyorm.Col("Name").Ge(val)
}

func SubModelNameLt(val string) easyorm.Predicate {
	return easyorm.Col("Name").Lt(val)
}

func SubModelNameLe(val string) easyorm.Predicate {
	return easyorm.Col("Name").Le(val)
}

func SubModelEmailEq(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Email").Eq(val)
}

func SubModelEmailNe(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Email").Ne(val)
}

func SubModelEmailGt(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Email").Gt(val)
}

func SubModelEmailGe(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Email").Ge(val)
}

func SubModelEmailLt(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Email").Lt(val)
}

func SubModelEmailLe(val *sqlx.NullString) easyorm.Predicate {
	return easyorm.Col("Email").Le(val)
}

func SubModelBalanceEq(val float64) easyorm.Predicate {
	return easyorm.Col("Balance").Eq(val)
}

func SubModelBalanceNe(val float64) easyorm.Predicate {
	return easyorm.Col("Balance").Ne(val)
}

func SubModelBalanceGt(val float64) easyorm.Predicate {
	return easyorm.Col("Balance").Gt(val)
}

func SubModelBalanceGe(val float64) easyorm.Predicate {
	return easyorm.Col("Balance").Ge(val)
}

func SubModelBalanceLt(val float64) easyorm.Predicate {
	return easyorm.Col("Balance").Lt(val)
}

func SubModelBalanceLe(val float64) easyorm.Predicate {
	return easyorm.Col("Balance").Le(val)
}
