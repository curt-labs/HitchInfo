package database

import (
	"errors"
	"expvar"
	"github.com/ziutek/mymysql/autorc"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"log"
	"reflect"
	"runtime"
	"strings"
)

type T struct{}

// prepared statements go here
var (
	Statements = make(map[string]*autorc.Stmt, 0)
)

func GetFunctionName(i interface{}) string {
	arr := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	if strings.Contains(arr, "database.") {
		strArr := strings.Split(arr, "database.")
		return strArr[1]
	}
	return arr
}

func PrepareAll() error {

	BindDatabase()

	adminChan := make(chan int)
	devChan := make(chan int)
	pcdbChan := make(chan int)
	vcdbChan := make(chan int)

	go PrepareAdmin(adminChan)
	go PrepareCurtDev(devChan)
	go PreparePcdb(pcdbChan)
	go PrepareVcdb(vcdbChan)

	log.Println("Executing Prepared Statements...")

	<-adminChan
	log.Printf("AdminDB Statements Completed.............\033[32;1m[OK]\033[0m")
	<-pcdbChan
	log.Printf("PCDB Statements Completed.............\033[32;1m[OK]\033[0m")
	<-vcdbChan
	log.Printf("VCDB Statements Completed.............\033[32;1m[OK]\033[0m")
	<-devChan
	log.Printf("CurtDev Statements Completed.............\033[32;1m[OK]\033[0m")
	return nil
}

// Prepare all MySQL statements
func PrepareAdmin(ch chan int) {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["getID"] = "select LAST_INSERT_ID() AS id"
	UnPreparedStatements["authenticateUserStmt"] = "select * from user where username = ? and encpassword = ? and isActive = 1"
	UnPreparedStatements["getUserByIDStmt"] = "select * from user where id=?"
	UnPreparedStatements["getUserByUsernameStmt"] = "select * from user where username=?"
	UnPreparedStatements["getUserByEmailStmt"] = "select * from user where email=?"
	UnPreparedStatements["allUserStmt"] = "select * from user"
	UnPreparedStatements["getAllModulesStmt"] = "select * from module order by module"
	UnPreparedStatements["userModulesStmt"] = "select module.* from module inner join user_module on module.id = user_module.moduleID where user_module.userID = ? order by module"
	UnPreparedStatements["setUserPasswordStmt"] = "update user set encpassword = ? where id = ?"
	UnPreparedStatements["registerUserStmt"] = "insert into user (username,email,fname,lname,isActive,superUser) VALUES (?,?,?,?,0,0)"
	UnPreparedStatements["getAllUserStmt"] = "select * from user order by fname, lname"
	UnPreparedStatements["setUserStatusStmt"] = "update user set isActive = ? WHERE id = ?"
	UnPreparedStatements["clearUserModuleStmt"] = "delete from user_module WHERE userid = ?"
	UnPreparedStatements["deleteUserStmt"] = "delete from user WHERE id = ?"
	UnPreparedStatements["addUserStmt"] = "insert into user (username,email,fname,lname,biography,photo,isActive,superUser) VALUES (?,?,?,?,?,?,?,?)"
	UnPreparedStatements["updateUserStmt"] = "update user set username=?, email=?, fname=?, lname=?, biography=?, photo=?, isActive=?, superUser=? WHERE id = ?"
	UnPreparedStatements["addModuleToUserStmt"] = "insert into user_module (userID,moduleID) VALUES (?,?)"

	if !AdminDb.Raw.IsConnected() {
		AdminDb.Raw.Connect()
	}

	AdminDb.Register("SET NAMES latin1")

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareAdminStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
	ch <- 1
}

func PrepareCurtDev(ch chan int) {

	chn := make(chan int, 0)

	// TESTING
	var t T
	val := reflect.ValueOf(&t)
	for i := 0; i < val.NumMethod(); i++ {
		m := val.Method(i)
		if iface := m.Interface(); iface != nil {
			go func(method func(), c chan int) {
				method()
				c <- 1
			}(iface.(func()), chn)
		}
	}

	for i := 0; i < val.NumMethod(); i++ {
		m := val.Method(i)
		if iface := m.Interface(); iface != nil {
			<-chn
			// me := iface.(func())
			log.Printf(" ~ \033[32;1m [OK]\033[0m %s Statements Completed", GetFunctionName(iface))
		}
	}

	ch <- 1
}

func PrepareVcdb(ch chan int) {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["vcdb_GetMakes"] = `select distinct ma.MakeID, ma.MakeName
																						from Make as ma
																						join BaseVehicle as bv on ma.MakeID = bv.MakeID
																						join Model as mo on bv.ModelID = mo.ModelID
																						where mo.VehicleTypeID = 5 || mo.VehicleTypeID = 6 || mo.VehicleTypeID = 7
																						order by MakeName`
	UnPreparedStatements["vcdb_GetMakeByName"] = `select ma.MakeID, ma.MakeName from Make as ma where LOWER(ma.MakeName) = ? limit 1`
	UnPreparedStatements["vcdb_GetMake"] = `select ma.MakeID, ma.MakeName from Make as ma where ma.MakeID = ? limit 1`
	UnPreparedStatements["vcdb_GetModels"] = `select mo.ModelID, mo.ModelName from Model mo order by mo.ModelName`
	UnPreparedStatements["vcdb_GetModelsByYearMake"] = `select mo.ModelID, mo.ModelName from Model mo
														join BaseVehicle bv on mo.ModelID = bv.ModelID
														where bv.YearID = ? && bv.MakeID = ?
														group by mo.ModelID
														order by mo.ModelName`
	UnPreparedStatements["vcdb_GetModelsByYear"] = `select mo.ModelID, mo.ModelName from Model mo
													join BaseVehicle bv on mo.ModelID = bv.ModelID
													where bv.YearID = ?
													group by mo.ModelID
													order by mo.ModelName`
	UnPreparedStatements["vcdb_GetModelsByMake"] = `select mo.ModelID, mo.ModelName from Model mo
													join BaseVehicle bv on mo.ModelID = bv.ModelID
													where bv.MakeID = ?
													group by mo.ModelID
													order by mo.ModelName`
	UnPreparedStatements["vcdb_GetBaseVehicles"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, m.MakeName, ma.ModelName, cv.ID, r.RegionAbbr, cv.ConfigID,
													(
														select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
														join CurtDev.vcdb_Vehicle as cv1 on vp.VehicleID = cv1.ID
														join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
														where cbv1.AAIABaseVehicleID = bv.BaseVehicleID && cv1.SubModelID = 0
													) as parts
													from BaseVehicle bv
													join Make m on bv.MakeID = m.MakeID
													join Model ma on bv.ModelID = ma.ModelID
													join Vehicle v on bv.BaseVehicleID = v.BaseVehicleID
													join Region r on v.RegionID = r.RegionID
													left join CurtDev.vcdb_Vehicle cv on v.VehicleID = cv.ID
													group by bv.BaseVehicleID`
	UnPreparedStatements["vcdb_GetBaseVehiclesByYearMakeModel"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, m.MakeName, ma.ModelName, cv.ID, r.RegionAbbr, cv.ConfigID,
																	(
																		select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																		join CurtDev.vcdb_Vehicle as cv1 on vp.VehicleID = cv1.ID
																		join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																		where cbv1.AAIABaseVehicleID = bv.BaseVehicleID && cv1.SubModelID = 0
																	) as parts
																	from BaseVehicle bv
																	join Make m on bv.MakeID = m.MakeID
																	join Model ma on bv.ModelID = ma.ModelID
																	join Vehicle v on bv.BaseVehicleID = v.BaseVehicleID
																	join Region r on v.RegionID = r.RegionID
																	left join CurtDev.vcdb_Vehicle cv on v.VehicleID = cv.ID
																	where bv.YearID = ? && bv.MakeID = ? && bv.ModelID = ?
																	group by bv.BaseVehicleID`
	UnPreparedStatements["vcdb_GetBaseVehiclesByYearMake"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, m.MakeName, ma.ModelName, cv.ID, r.RegionAbbr, cv.ConfigID,
																(
																	select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																	join CurtDev.vcdb_Vehicle as cv1 on vp.VehicleID = cv1.ID
																	join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																	where cbv1.AAIABaseVehicleID = bv.BaseVehicleID && cv1.SubModelID = 0
																) as parts
																from BaseVehicle bv
																join Make m on bv.MakeID = m.MakeID
																join Model ma on bv.ModelID = ma.ModelID
																join Vehicle v on bv.BaseVehicleID = v.BaseVehicleID
																join Region r on v.RegionID = r.RegionID
																left join CurtDev.vcdb_Vehicle cv on v.VehicleID = cv.ID
																where bv.YearID = ? && bv.MakeID = ?
																group by bv.BaseVehicleID`
	UnPreparedStatements["vcdb_GetBaseVehiclesByMakeModel"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, m.MakeName, ma.ModelName, cv.ID, r.RegionAbbr, cv.ConfigID,
																(
																	select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																	join CurtDev.vcdb_Vehicle as cv1 on vp.VehicleID = cv1.ID
																	join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																	where cbv1.AAIABaseVehicleID = bv.BaseVehicleID && cv1.SubModelID = 0
																) as parts
																from BaseVehicle bv
																join Make m on bv.MakeID = m.MakeID
																join Model ma on bv.ModelID = ma.ModelID
																join Vehicle v on bv.BaseVehicleID = v.BaseVehicleID
																join Region r on v.RegionID = r.RegionID
																left join CurtDev.vcdb_Vehicle cv on v.VehicleID = cv.ID
																where bv.MakeID = ? && bv.ModelID = ?
																group by bv.BaseVehicleID`
	UnPreparedStatements["vcdb_GetBaseVehiclesByYear"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, m.MakeName, ma.ModelName, cv.ID, r.RegionAbbr, cv.ConfigID,
															(
																select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																join CurtDev.vcdb_Vehicle as cv1 on vp.VehicleID = cv1.ID
																join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																where cbv1.AAIABaseVehicleID = bv.BaseVehicleID && cv1.SubModelID = 0
															) as parts
															from BaseVehicle bv
															join Make m on bv.MakeID = m.MakeID
															join Model ma on bv.ModelID = ma.ModelID
															join Vehicle v on bv.BaseVehicleID = v.BaseVehicleID
															join Region r on v.RegionID = r.RegionID
															left join CurtDev.vcdb_Vehicle cv on v.VehicleID = cv.ID
															where bv.YearID = ?
															group by bv.BaseVehicleID, bv.YearID, bv.MakeID, bv.ModelID`
	UnPreparedStatements["vcdb_GetBaseVehicleByVehicleId"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, cv.ID, r.RegionAbbr, cv.ConfigID,
																														(
																															select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																															join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
																															join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																															where cbv1.AAIABaseVehicleID = bv.BaseVehicleID && cv1.SubModelID = 0
																														) as parts
																														from Vehicle v
																														join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
																														join Make ma on bv.MakeID = ma.MakeID
																														join Model mo on bv.ModelID = mo.ModelID
																														join Region r on v.RegionID = r.RegionID
																														join Submodel s on v.SubmodelID = s.SubmodelID
																														left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
																														left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
																														where v.VehicleID = ? limit 1`
	UnPreparedStatements["vcdb_GetVehicles"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, ma.MakeName, mo.ModelName, s.SubmodelName, cv.ID, r.RegionAbbr, cv.ConfigID,
													(
														select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
														join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
														join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
														where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
													) as parts,
													(
														select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
														join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
														join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
														where vca.VehicleConfigID = cv.ConfigID
														order by cat.sort
													) as config_values,
													(
														select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
														join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
														join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
														where vca.VehicleConfigID = cv.ConfigID
														order by cat.sort
													) as config_types
													from Vehicle v
													join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
													join Make ma on bv.MakeID = ma.MakeID
													join Model mo on bv.ModelID = mo.ModelID
													join Region r on v.RegionID = r.RegionID
													join Submodel s on v.SubmodelID = s.SubmodelID
													left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
													left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
													group by v.VehicleID, bv.BaseVehicleID
													order by cv.ID desc`
	UnPreparedStatements["vcdb_GetVehiclesByPart"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, s.SubmodelName, s.SubmodelID, cv.ID, r.RegionAbbr, cv.ConfigID,
													(
														select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
														join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
														join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
														where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
													) as parts,
													(
														select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
														join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
														join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
														where vca.VehicleConfigID = cv.ConfigID
														order by cat.sort
													) as config_values,
													(
														select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
														join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
														join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
														where vca.VehicleConfigID = cv.ConfigID
														order by cat.sort
													) as config_types
													from Vehicle v
													join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
													join Make ma on bv.MakeID = ma.MakeID
													join Model mo on bv.ModelID = mo.ModelID
													join Region r on v.RegionID = r.RegionID
													join Submodel s on v.SubmodelID = s.SubmodelID
													left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
													left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
													left join CurtDev.vcdb_VehiclePart as vp on cv.ID = vp.VehicleID
													where vp.PartNumber = ?
													group by v.VehicleID, bv.BaseVehicleID
													order by cv.ID desc`
	UnPreparedStatements["vcdb_GetVehiclesByYearMakeModel"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, s.SubmodelName, s.SubmodelID, cv.ID, r.RegionAbbr, cv.ConfigID,
																	(
																		select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																		join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
																		join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																		where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
																	) as parts,
																	(
																		select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
																		join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																		join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																		where vca.VehicleConfigID = cv.ConfigID
																		order by cat.sort
																	) as config_values,
																	(
																		select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
																		join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																		join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																		where vca.VehicleConfigID = cv.ConfigID
																		order by cat.sort
																	) as config_types
																	from Vehicle v
																	join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
																	join Make ma on bv.MakeID = ma.MakeID
																	join Model mo on bv.ModelID = mo.ModelID
																	join Region r on v.RegionID = r.RegionID
																	join Submodel s on v.SubmodelID = s.SubmodelID
																	left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
																	left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
																	where bv.YearID = ? && bv.MakeID = ? && bv.ModelID = ?
																	group by v.VehicleID, bv.BaseVehicleID
																	order by cv.ID desc`
	UnPreparedStatements["vcdb_GetVehiclesByYearMake"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, s.SubmodelName, s.SubmodelID, cv.ID, r.RegionAbbr, cv.ConfigID,
																(
																	select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																	join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
																	join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																	where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
																) as parts,
																(
																	select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
																	join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																	join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																	where vca.VehicleConfigID = cv.ConfigID
																	order by cat.sort
																) as config_values,
																(
																	select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
																	join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																	join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																	where vca.VehicleConfigID = cv.ConfigID
																	order by cat.sort
																) as config_types
																from Vehicle v
																join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
																join Make ma on bv.MakeID = ma.MakeID
																join Model mo on bv.ModelID = mo.ModelID
																join Region r on v.RegionID = r.RegionID
																join Submodel s on v.SubmodelID = s.SubmodelID
																left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
																left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
																where bv.YearID = ? && bv.MakeID = ?
																group by v.VehicleID, bv.BaseVehicleID
																order by cv.ID desc`
	UnPreparedStatements["vcdb_GetVehiclesByMakeModel"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, s.SubmodelName, s.SubmodelID, cv.ID, r.RegionAbbr, cv.ConfigID,
															(
																select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
																join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
															) as parts,
															(
																select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
																join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																join CurtDev.vcdb_Vehicle cv2 on vca.VehicleConfigID = cv2.ConfigID
																where cv2.ID = cv.ID
																order by cat.sort
															) as config_values,
															(
																select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
																join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																join CurtDev.vcdb_Vehicle cv2 on vca.VehicleConfigID = cv2.ConfigID
																where cv2.ID = cv.ID
																order by cat.sort
															) as config_types
															from Vehicle v
															join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
															join Make ma on bv.MakeID = ma.MakeID
															join Model mo on bv.ModelID = mo.ModelID
															join Region r on v.RegionID = r.RegionID
															join Submodel s on v.SubmodelID = s.SubmodelID
															left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
															left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
															where bv.MakeID = ? && bv.ModelID = ?
															order by cv.ID desc`
	UnPreparedStatements["vcdb_GetVehiclesByYear"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, s.SubmodelName, s.SubmodelID, cv.ID, r.RegionAbbr, cv.ConfigID,
															(
																select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
																join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
															) as parts,
															(
																select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
																join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																where vca.VehicleConfigID = cv.ConfigID
																order by cat.sort
															) as config_values,
															(
																select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
																join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																where vca.VehicleConfigID = cv.ConfigID
																order by cat.sort
															) as config_types
															from Vehicle v
															join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
															join Make ma on bv.MakeID = ma.MakeID
															join Model mo on bv.ModelID = mo.ModelID
															join Region r on v.RegionID = r.RegionID
															join Submodel s on v.SubmodelID = s.SubmodelID
															left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
															left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
															where bv.YearID = ?
															group by v.VehicleID, bv.BaseVehicleID
															order by cv.ID desc`
	UnPreparedStatements["vcdb_GetVehicleById"] = `select v.VehicleID, bv.YearID, bv.MakeID, bv.ModelID, bv.BaseVehicleID, ma.MakeName, mo.ModelName, s.SubmodelName, s.SubmodelID, cv.ID, r.RegionAbbr, cv.ConfigID,
															(
																select group_concat(distinct vp.PartNumber SEPARATOR '|') from CurtDev.vcdb_VehiclePart as vp
																join CurtDev.vcdb_Vehicle cv1 on vp.VehicleID = cv1.ID
																join CurtDev.BaseVehicle cbv1 on cv1.BaseVehicleID = cbv1.ID
																where cbv1.AAIABaseVehicleID = bv.BaseVehicleID
															) as parts,
															(
																select group_concat(ca.value SEPARATOR '|') from CurtDev.ConfigAttribute ca
																join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																where vca.VehicleConfigID = cv.ConfigID
																order by cat.sort
															) as config_values,
															(
																select group_concat(cat.name SEPARATOR '|') from CurtDev.ConfigAttribute ca
																join CurtDev.ConfigAttributeType cat on ca.ConfigAttributeTypeID = cat.ID
																join CurtDev.VehicleConfigAttribute vca on ca.ID = vca.AttributeID
																where vca.VehicleConfigID = cv.ConfigID
																order by cat.sort
															) as config_types
															from Vehicle v
															join BaseVehicle bv on v.BaseVehicleID = bv.BaseVehicleID
															join Make ma on bv.MakeID = ma.MakeID
															join Model mo on bv.ModelID = mo.ModelID
															join Region r on v.RegionID = r.RegionID
															join Submodel s on v.SubmodelID = s.SubmodelID
															left join CurtDev.BaseVehicle cbv on bv.BaseVehicleID = cbv.AAIABaseVehicleID
															left join CurtDev.vcdb_Vehicle cv on cbv.ID = cv.BaseVehicleID
															where v.VehicleID = ? && s.SubmodelID = ? && cv.ConfigID = ? limit 1`

	if !Vcdb.Raw.IsConnected() {
		Vcdb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareVcdbStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
	//log.Println(Statements["vcdb_GetModelsByMake"])
	ch <- 1
}

func PreparePcdb(ch chan int) {

	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetPartTerminology"] = `select * from Parts where PartTerminologyID = ?`

	if !Pcdb.Raw.IsConnected() {
		Pcdb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PreparePcdbStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
	ch <- 1
}

func (t *T) WebsiteStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Website Statements
	UnPreparedStatements["getAllSiteContentStmt"] = "select * from SiteContent WHERE active = 1 order by page_title"
	UnPreparedStatements["getPrimaryMenuStmt"] = "select * from Menu where isPrimary = 1"
	UnPreparedStatements["getMenuByIDStmt"] = "select * from Menu where menuID = ?"
	UnPreparedStatements["getMenuItemsStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC
												INNER JOIN Menu AS M ON MSC.menuID = M.menuID
												LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
												WHERE MSC.menuID = ?`
	UnPreparedStatements["getAllMenuItemLinksStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC
													INNER JOIN Menu AS M ON MSC.menuID = M.menuID
													LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
													where SC.page_title is null
													order by MSC.menuTitle`
	UnPreparedStatements["GetContentRevisionsStmt"] = "select * from SiteContentRevision WHERE contentID = ?"
	UnPreparedStatements["GetAllMenusStmt"] = "select * from Menu where active = 1"
	UnPreparedStatements["UpdateMenuStmt"] = "Update Menu Set menu_name = ?, requireAuthentication = ?, showOnSitemap = ?, display_name = ? where menuID = ?"
	UnPreparedStatements["AddMenuStmt"] = `INSERT INTO Menu (menu_name,display_name,requireAuthentication,showOnSitemap,isPrimary,active,sort) VALUES (?,?,?,?,0,1,1)`
	UnPreparedStatements["getInsertedMenuID"] = "select LAST_INSERT_ID() FROM Menu AS id LIMIT 1"
	UnPreparedStatements["deleteMenuStmt"] = "Update Menu set active = 0 WHERE menuID = ?"
	UnPreparedStatements["clearPrimaryMenuStmt"] = "update Menu set isPrimary = 0"
	UnPreparedStatements["setPrimaryMenuStmt"] = "update Menu set isPrimary = 1 WHERE menuID = ?"
	UnPreparedStatements["getMenuSortStmt"] = "select menuSort from Menu_SiteContent WHERE menuID = ? order by menuSort DESC"
	UnPreparedStatements["addMenuContentItemStmt"] = "INSERT INTO Menu_SiteContent (menuID,contentID,menuSort,parentID) VALUES (?,?,?,0)"
	UnPreparedStatements["addMenuLinkItemStmt"] = "INSERT INTO Menu_SiteContent (menuID,menuTitle,menuLink,linkTarget,menuSort,contentID,parentID) VALUES (?,?,?,?,?,0,0)"
	UnPreparedStatements["updateMenuLinkItemStmt"] = "UPDATE Menu_SiteContent set menuTitle = ?, menuLink = ?, linkTarget = ? WHERE menuContentID = ?"
	UnPreparedStatements["deleteMenuLinkItemStmt"] = "delete from Menu_SiteContent WHERE menuContentID = ?"
	UnPreparedStatements["getMenuItemStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC
												INNER JOIN Menu AS M ON MSC.menuID = M.menuID
												LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
												WHERE MSC.menuContentID = ?`
	UnPreparedStatements["getMenuItemByContentIDStmt"] = `select MSC.menuContentID, MSC.menuID, MSC.menuSort, MSC.menuTitle, MSC.menuLink, MSC.parentID, MSC.linkTarget, SC.* from Menu_SiteContent AS MSC
												INNER JOIN Menu AS M ON MSC.menuID = M.menuID
												LEFT JOIN SiteContent AS SC ON MSC.contentID = SC.contentID
												WHERE SC.contentID = ?`
	UnPreparedStatements["GetMenuParentsStmt"] = "select * from Menu_SiteContent where parentID = ? AND menuID = ? order by menuSort"
	UnPreparedStatements["DeleteMenuItemStmt"] = "delete from Menu_SiteContent where menuContentID = ?"
	UnPreparedStatements["clearPrimaryContentStmt"] = "update SiteContent set isPrimary = 0 WHERE isPrimary = 1"
	UnPreparedStatements["setPrimaryContentStmt"] = "update SiteContent set isPrimary = 1 WHERE contentID = ?"
	UnPreparedStatements["getContentStmt"] = "select * from SiteContent WHERE contentID = ?"
	UnPreparedStatements["deleteContentStmt"] = "update SiteContent set active = 0 WHERE contentID = ?"
	UnPreparedStatements["checkContentStmt"] = `select M.menu_name FROM Menu AS M
												INNER JOIN Menu_SiteContent AS MSC ON M.menuID = MSC.menuID
												WHERE MSC.contentID = ?`
	UnPreparedStatements["addContentStmt"] = `insert into SiteContent (page_title,content_type,createdDate,lastModified,meta_title,meta_description,keywords,isPrimary,published,active,slug,requireAuthentication,canonical)
											  VALUES (?,"",?,?,?,?,?,0,?,1,?,?,?)`
	UnPreparedStatements["updateContentStmt"] = `update SiteContent set page_title = ?, meta_title = ?, meta_description = ?,keywords = ?, published = ?, slug = ?, requireAuthentication = ?, canonical = ?
											     where contentID = ?`
	UnPreparedStatements["addContentRevisionStmt"] = `insert into SiteContentRevision (contentID,content_text,createdOn,active)
													  VALUES (?,?,?,?)`
	UnPreparedStatements["updateContentRevisionStmt"] = `update SiteContentRevision Set content_text = ? WHERE revisionID = ?`
	UnPreparedStatements["copyContentRevisionStmt"] = `insert into SiteContentRevision (contentID,content_text,createdOn,active)
													   (select contentID,content_text,?,0 from SiteContentRevision WHERE revisionID = ?)`
	UnPreparedStatements["getRevisionContentIDStmt"] = `select contentID from SiteContentRevision WHERE revisionID = ?`
	UnPreparedStatements["deactivateContentRevisionStmt"] = `update SiteContentRevision set active = 0 WHERE contentID = ? and active = 1`
	UnPreparedStatements["activateContentRevisionStmt"] = `update SiteContentRevision set active = 1 WHERE revisionID = ?`
	UnPreparedStatements["deleteContentRevisionStmt"] = `delete from SiteContentRevision WHERE revisionID = ?`

	UnPreparedStatements["UpdateMenuItemSortStmt"] = `Update Menu_SiteContent Set menuSort = ? where menuContentID = ?`
	UnPreparedStatements["GetMenuItemsWithSortGreaterOrEqualStmt"] = `select menuContentID from Menu_SiteContent where menuSort >= ? && menuID = ? order by menuSort`
	UnPreparedStatements["GetMenuItemsWithSortLessStmt"] = `select menuContentID from Menu_SiteContent where menuSort < ? && menuID = ? order by menuSort`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
	return
}

func (t *T) ContactStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Contact Manager Statements
	UnPreparedStatements["getAllContactsStmt"] = `select * from Contact limit ?,?`
	UnPreparedStatements["getContactCountStmt"] = `select count(contactID) as count from Contact`
	UnPreparedStatements["getContactStmt"] = `select * from Contact WHERE contactID = ?`
	UnPreparedStatements["getAllContactTypesStmt"] = `select * from ContactType`
	UnPreparedStatements["getContactTypeStmt"] = `select * from ContactType where contactTypeID = ? limit 0,1`
	UnPreparedStatements["getAllContactReceiversStmt"] = `select * from ContactReceiver`
	UnPreparedStatements["getContactReceiversStmt"] = `select * from ContactReceiver WHERE contactReceiverID = ?`
	UnPreparedStatements["getReceiverContactTypesStmt"] = `select CT.* from ContactType AS CT
														   INNER JOIN ContactReceiver_ContactType AS CR ON CT.contactTypeID = CR.contactTypeID
														   WHERE CR.contactReceiverID = ?`
	UnPreparedStatements["clearReceiverTypesStmt"] = `delete from ContactReceiver_ContactType WHERE contactReceiverID = ?`
	UnPreparedStatements["addReceiverTypeStmt"] = `insert into ContactReceiver_ContactType (contactReceiverID,contactTypeID) VALUES (?,?)`
	UnPreparedStatements["addContactReceiverStmt"] = `INSERT INTO ContactReceiver (first_name,last_name,email) VALUES (?,?,?)`
	UnPreparedStatements["updateContactReceiverStmt"] = `update ContactReceiver SET first_name = ?, last_name = ?, email = ? where contactReceiverID = ?`
	UnPreparedStatements["deleteContactReceiverStmt"] = `delete from ContactReceiver where contactReceiverID = ?`
	UnPreparedStatements["addContactTypeStmt"] = `insert into ContactType (name) VALUE (?)`
	UnPreparedStatements["updateContactTypeStmt"] = `update ContactType set name = ? where contactTypeID = ?`
	UnPreparedStatements["deleteContactTypeStmt"] = `delete from ContactType WHERE contactTypeID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) FAQStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// FAQ Manager Statements
	UnPreparedStatements["GetAllFAQStmt"] = `select * from FAQ`
	UnPreparedStatements["GetFAQStmt"] = `select * from FAQ WHERE faqID = ?`
	UnPreparedStatements["UpdateFAQStmt"] = `UPDATE FAQ SET question = ?, answer = ? WHERE faqID = ?`
	UnPreparedStatements["AddFAQStmt"] = `INSERT INTO FAQ (question, answer) VALUES (?,?)`
	UnPreparedStatements["DeleteFAQStmt"] = `DELETE FROM FAQ WHERE faqID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) NewsStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// News Manager Statements
	UnPreparedStatements["GetAllNewsStmt"] = `select * from NewsItem WHERE active = 1`
	UnPreparedStatements["GetNewsItemStmt"] = `select * from NewsItem WHERE newsItemID = ?`
	UnPreparedStatements["UpdateNewsItemStmt"] = `Update NewsItem SET title = ?, lead = ?, content = ?, publishStart = ?, publishEnd = ?, slug = ? WHERE newsItemID = ?`
	UnPreparedStatements["AddNewsItemStmt"] = `INSERT INTO NewsItem (title,lead,content,publishStart,publishEnd,slug,active) VALUES (?,?,?,?,?,?,1)`
	UnPreparedStatements["DeleteNewsItemStmt"] = `Update NewsItem SET active = 0 WHERE newsItemID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) MiscStatements() {
	UnPreparedStatements := make(map[string]string, 0)
	//Miscellaneous Statements

	UnPreparedStatements["GetAllContentTypesStmt"] = `SELECT * FROM ContentType`
	UnPreparedStatements["GetContentTypeStmt"] = `SELECT * FROM ContentType WHERE cTypeID = ?`
	UnPreparedStatements["UpdateContentTypeStmt"] = "UPDATE ContentType SET type = ?, allowHTML = ? WHERE cTypeID = ?"
	UnPreparedStatements["AddContentTypeStmt"] = `INSERT into ContentType(type,allowHTML) VALUES(?,?)`
	UnPreparedStatements["DeleteContentTypeStmt"] = `DELETE FROM ContentType WHERE cTypeID = ?`

	UnPreparedStatements["GetAllVideoTypesStmt"] = `SELECT * FROM videoType`
	UnPreparedStatements["GetVideoTypeStmt"] = `SELECT * FROM videoType WHERE vTypeID = ?`
	UnPreparedStatements["UpdateVideoTypeStmt"] = "UPDATE videoType SET name = ?, icon = ? WHERE vTypeID = ?"
	UnPreparedStatements["AddVideoTypeStmt"] = `INSERT into videoType(name,icon) VALUES(?,?)`
	UnPreparedStatements["DeleteVideoTypeStmt"] = `DELETE FROM videoType WHERE vTypeID = ?`

	UnPreparedStatements["GetAllMeasureUnitsStmt"] = `SELECT * FROM UnitOfMeasure`
	UnPreparedStatements["GetMeasureUnitStmt"] = `SELECT * FROM UnitOfMeasure WHERE ID = ?`
	UnPreparedStatements["UpdateMeasureUnitStmt"] = `UPDATE UnitOfMeasure SET name = ?, code = ? WHERE ID = ?`
	UnPreparedStatements["AddMeasureUnitStmt"] = `INSERT UnitOfMeasure(name,code) VALUES(?,?)`
	UnPreparedStatements["DeleteMeasureUnitSmt"] = `DELETE FROM UnitOfMeasure WHERE ID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) VideoStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Video Manager Statements
	UnPreparedStatements["GetAllVideosStmt"] = `select * from Video order by sort`
	UnPreparedStatements["GetVideoStmt"] = `select * from Video where videoID = ?`
	UnPreparedStatements["DeleteVideoStmt"] = `DELETE FROM Video where videoID = ?`
	UnPreparedStatements["GetLastVideoSortStmt"] = `select sort from Video Order By sort desc`
	UnPreparedStatements["AddVideoStmt"] = `INSERT INTO Video (embed_link,dateAdded,sort,title,description,youtubeID,watchpage,screenshot) VALUES (?,?,?,?,?,?,?,?)`
	UnPreparedStatements["UpdateVideoSortStmt"] = `Update Video Set sort = ? where videoID = ?`
	UnPreparedStatements["GetVideosWithSortGreaterOrEqualStmt"] = `select videoID from Video where sort >= ? order by sort`
	UnPreparedStatements["GetVideosWithSortLessStmt"] = `select videoID from Video where sort < ? order by sort`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) LandingPageStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	//Landing Page Manager Statements
	UnPreparedStatements["GetActiveLandingPagesStmt"] = `SELECT * from LandingPage WHERE endDate > UTC_TIMESTAMP()`
	UnPreparedStatements["GetPastLandingPagesStmt"] = `SELECT * from LandingPage WHERE endDate < UTC_TIMESTAMP()`
	UnPreparedStatements["GetLandingPageImagesStmt"] = `SELECT * from LandingPageImages WHERE landingPageID = ?`
	UnPreparedStatements["GetLandingPageDataStmt"] = `SELECT * from LandingPageData WHERE landingPageID = ?`
	UnPreparedStatements["AddLandingPageStmt"] = `INSERT INTO LandingPage (name,startDate,endDate,url,pageContent,linkClasses,conversionID,conversionLabel,newWindow,menuPosition) VALUES (?,?,?,?,?,?,?,?,?,?)`
	UnPreparedStatements["UpdateLandingPageStmt"] = `UPDATE LandingPage SET name = ?, startDate = ?, endDate = ?, url = ?, pageContent = ?, linkClasses = ?, conversionID = ?, conversionLabel = ?, newWindow = ?, menuPosition = ? WHERE id = ?`
	UnPreparedStatements["GetLandingPageStmt"] = `SELECT * from LandingPage WHERE id = ?`
	UnPreparedStatements["AddLandingPageDataStmt"] = `INSERT INTO LandingPageData (landingPageID,dataKey,dataValue) VALUES (?,?,?)`
	UnPreparedStatements["DeleteLandingPageDataStmt"] = `DELETE FROM LandingPageData WHERE id = ?`
	UnPreparedStatements["TruncateLandingPageDataStmt"] = `DELETE FROM LandingPageData WHERE landingPageID = ?`
	UnPreparedStatements["DeleteLandingPageStmt"] = `DELETE FROM LandingPage WHERE id = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) TestimonialStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	//Testimonial Manager Statements
	UnPreparedStatements["GetAllTestimonialsStmt"] = `SELECT * from Testimonial WHERE active = 1 AND approved = ?`
	UnPreparedStatements["GetTestimonialStmt"] = `select * from Testimonial where testimonialID = ? limit 1`
	UnPreparedStatements["DeleteTestimonialStmt"] = `UPDATE Testimonial SET active = 0 WHERE testimonialID = ?`
	UnPreparedStatements["SetTestimonialApprovalStmt"] = `UPDATE Testimonial SET approved = ? WHERE testimonialID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) SalesRepStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Sales Reps
	UnPreparedStatements["GetAllSalesRepsStmt"] = `SELECT salesRepID, name, code, (select COUNT(cust_id) from Customer WHERE Customer.salesRepID = SalesRepresentative.salesRepID) AS customercount from SalesRepresentative`
	UnPreparedStatements["GetSalesRepStmt"] = `SELECT salesRepID, name, code, (select COUNT(cust_id) from Customer WHERE Customer.salesRepID = SalesRepresentative.salesRepID) AS customercount from SalesRepresentative WHERE salesRepID = ?`
	UnPreparedStatements["UpdateSalesRepStmt"] = `UPDATE SalesRepresentative set name = ?, code = ? WHERE salesRepID = ?`
	UnPreparedStatements["AddSalesRepStmt"] = `INSERT INTO SalesRepresentative (name,code) VALUES (?,?)`
	UnPreparedStatements["DeleteSalesRepStmt"] = `DELETE FROM SalesRepresentative WHERE salesRepID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) BlogStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Blog
	UnPreparedStatements["GetAllPostsStmt"] = `SELECT * from BlogPosts WHERE active = 1`
	UnPreparedStatements["GetPostStmt"] = `SELECT * from BlogPosts WHERE blogPostID = ?`
	UnPreparedStatements["GetAllBlogCategoriesStmt"] = `SELECT * from BlogCategories WHERE active = 1`
	UnPreparedStatements["GetBlogCategoryStmt"] = `SELECT * from BlogCategories WHERE blogCategoryID = ?`
	UnPreparedStatements["AddBlogCategoryStmt"] = `INSERT INTO BlogCategories (name,slug,active) VALUES (?,?,1)`
	UnPreparedStatements["UpdateBlogCategoryStmt"] = `Update BlogCategories set name = ?, slug = ? WHERE blogCategoryID = ?`
	UnPreparedStatements["DeleteBlogCategoryStmt"] = `Update BlogCategories set active = 0 WHERE blogCategoryID = ?`
	UnPreparedStatements["GetPostCategoriesStmt"] = `SELECT BC.* from BlogCategories AS BC
													 INNER JOIN BlogPost_BlogCategory AS BPBC ON BC.blogCategoryID = BPBC.blogCategoryID
													 WHERE blogPostID = ? AND BC.active = 1`
	UnPreparedStatements["GetPostCommentsStmt"] = `SELECT * from Comments WHERE blogPostID = ? AND active = 1`
	UnPreparedStatements["AddPostStmt"] = `INSERT INTO BlogPosts (post_title,slug,post_text,createdDate,userID,meta_title,meta_description,keywords,active) VALUES (?,?,?,UTC_TIMESTAMP(),?,?,?,?,1)`
	UnPreparedStatements["UpdatePostStmt"] = `UPDATE BlogPosts SET post_title = ?, slug = ?, post_text = ?, userID = ?, meta_title = ?, meta_description = ?, keywords = ? WHERE blogPostID = ?`
	UnPreparedStatements["PublishPostStmt"] = `UPDATE BlogPosts SET publishedDate = UTC_TIMESTAMP() WHERE blogPostID = ?`
	UnPreparedStatements["UnPublishPostStmt"] = `UPDATE BlogPosts SET publishedDate = null WHERE blogPostID = ?`
	UnPreparedStatements["ClearPostCategoriesStmt"] = `DELETE FROM BlogPost_BlogCategory WHERE blogPostID = ?`
	UnPreparedStatements["AddPostCategoryStmt"] = `INSERT INTO BlogPost_BlogCategory (blogPostID,blogCategoryID) VALUES (?,?)`
	UnPreparedStatements["DeletePostStmt"] = `UPDATE BlogPosts SET active = 0 WHERE blogPostID = ?`
	UnPreparedStatements["GetBlogCommentsStmt"] = `SELECT * from Comments WHERE active = 1 AND approved = 0`
	UnPreparedStatements["GetBlogCommentStmt"] = `SELECT * from Comments WHERE commentID = ?`
	UnPreparedStatements["ApproveBlogCommentStmt"] = `UPDATE Comments set approved = 1 WHERE commentID = ?`
	UnPreparedStatements["DeleteBlogCommentStmt"] = `UPDATE Comments set active = 0 WHERE commentID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) ForumStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetAllForumGroupsStmt"] = `SELECT * from ForumGroup`
	UnPreparedStatements["GetForumGroupStmt"] = `SELECT * from ForumGroup WHERE forumGroupID = ?`
	UnPreparedStatements["GetForumGroupTopicsStmt"] = `SELECT FT.* from ForumTopic as FT
												   INNER JOIN ForumGroup as FG ON FT.TopicGroupID = FG.forumGroupID
												   WHERE FT.TopicGroupID = ? AND FT.active = 1`
	UnPreparedStatements["AddForumGroupStmt"] = `INSERT INTO ForumGroup (name,description,createdDate) VALUES (?,?,UTC_TIMESTAMP())`
	UnPreparedStatements["UpdateForumGroupStmt"] = `UPDATE ForumGroup set name = ?, description = ? WHERE forumGroupID = ?`
	UnPreparedStatements["DeleteForumGroupStmt"] = `DELETE from ForumGroup WHERE forumGroupID = ?`

	UnPreparedStatements["GetAllForumTopicsStmt"] = `SELECT * from ForumTopic WHERE active = 1`
	UnPreparedStatements["GetForumTopicStmt"] = `SELECT * from ForumTopic WHERE topicID = ?`
	UnPreparedStatements["AddForumTopicStmt"] = `INSERT into ForumTopic(TopicGroupID, name, description, image, createdDate, active, closed) VALUES (?,?,?,?,UTC_TIMESTAMP(),?,?)`
	UnPreparedStatements["UpdateForumTopicStmt"] = `UPDATE ForumTopic SET TopicGroupID = ?, name = ?, description = ?, image = ?, active = ?, closed = ? WHERE topicID = ?`
	UnPreparedStatements["DeleteForumTopicStmt"] = `UPDATE ForumTopic SET active = 0 WHERE topicID = ?`

	UnPreparedStatements["GetAllForumThreadsStmt"] = `SELECT * from ForumThread WHERE active = 1`
	UnPreparedStatements["GetForumThreadStmt"] = `SELECT * from ForumThread where threadID = ? AND active = 1`
	UnPreparedStatements["GetForumThreadsByTopic"] = `SELECT * from ForumThread where topicID = ? AND active = 1 ORDER BY createdDate DESC`
	UnPreparedStatements["DeleteForumThreadStmt"] = `UPDATE ForumThread SET active = 0 WHERE threadID = ?`
	UnPreparedStatements["UpdateForumThreadStmt"] = `UPDATE ForumThread SET topicID = ?, active = ?, closed = ? WHERE threadID = ?`
	UnPreparedStatements["AddForumThreadStmt"] = `INSERT into ForumThread(topicID,createdDate,active,closed) VALUES(?,UTC_TIMESTAMP(), ?, ?)`

	UnPreparedStatements["GetAllForumPostsStmt"] = `SELECT * from ForumPost WHERE active = 1 ORDER BY createdDate DESC`
	UnPreparedStatements["GetAllForumBaseForumPostsStmt"] = `SELECT * from ForumPost WHERE active = 1 AND parentID = 0 ORDER BY createdDate DESC`
	UnPreparedStatements["GetAllForumReplyPostsStmt"] = `SELECT * from ForumPost where parentID <> 0 AND active = 1`

	UnPreparedStatements["GetForumPostStmt"] = `SELECT * from ForumPost WHERE postID = ? AND active = 1`
	UnPreparedStatements["GetForumPostsByPostStmt"] = `SELECT * from ForumPost WHERE parentID = ? AND active = 1 ORDER BY sticky, createdDate DESC`
	UnPreparedStatements["GetForumPostsByThreadStmt"] = `SELECT * from ForumPost WHERE threadID = ? AND active = 1 AND flag = 0 ORDER BY sticky, parentID, createdDate DESC`
	UnPreparedStatements["AddForumPostStmt"] = `INSERT INTO ForumPost(parentID,threadID,createdDate,title,post,name,email,company,notify,approved,active,IPAddress,flag,sticky) VALUES(?,?,UTC_TIMESTAMP(),?,?,?,?,?,?,?,1,?,?,?)`
	UnPreparedStatements["UpdateForumPostStmt"] = `UPDATE ForumPost SET parentID = ?, threadID = ?, title = ?, post = ?, name = ?, email = ?, company = ?, notify = ?, approved = ?, IPAddress = ?, flag = ?, sticky = ? WHERE postID = ?`
	UnPreparedStatements["DeleteForumPostStmt"] = `UPDATE ForumPost SET active = 0 WHERE postID = ?`

	UnPreparedStatements["ApproveForumPostStmt"] = `UPDATE ForumPost SET approved = ? WHERE postID = ?`
	UnPreparedStatements["FlagForumPostStmt"] = `UPDATE ForumPost SET active = ?, flag = ? WHERE postID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) PartStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Part Statements
	UnPreparedStatements["PartVideoStmt"] = `select pv.pVideoID, pv.partID, pv.video, pv.vTypeID, vt.name,pv.isPrimary, vt.icon from PartVideo as pv
												join videoType vt on pv.vTypeID = vt.vTypeID
												where pv.partID = ?
												order by pv.isPrimary desc`

	UnPreparedStatements["GetAllParts"] = `select  p.*, c.class from Part as p
											left join Class as c on p.classID = c.classID`
	UnPreparedStatements["GetPartsByPage"] = `select  p.*, c.class from Part as p
											left join Class as c on p.classID = c.classID
											limit ?,?`
	UnPreparedStatements["GetAllPartsByStatus"] = `select  p.*, c.class from Part as p
													left join Class as c on p.classID = c.classID
													where status in (?)`
	UnPreparedStatements["GetPartCount"] = `select count(distinct partID) from Part`
	UnPreparedStatements["GetPartCountByStatus"] = `select count(distinct partID) from Part
											where status in (?)`
	UnPreparedStatements["GetPartStmt"] = `select  p.*, c.class, pa.value as upc from Part as p
											left join Class as c on p.classID = c.classID
											left join PartAttribute as pa on p.partID = pa.partID
											where p.partID = ? && (pa.field = 'UPC' || pa.field is null) limit 1`
	UnPreparedStatements["GetAllClasses"] = `select * from Class order by class`
	UnPreparedStatements["PartExistsStmt"] = `select partID from Part where partID = ?`
	UnPreparedStatements["UpdatePartStmt"] = `update Part set status = ?, dateModified = CURDATE(), shortDesc = ?,
												priceCode = ?, classID = ?, featured = ?
												where partID = ?`
	UnPreparedStatements["InsertPartStmt"] = `insert into Part (partID, status, dateAdded, dateModified, shortDesc, priceCode, classID, featured)
												values(?,?,CURDATE(),CURDATE(),?,?,?,?)`
	UnPreparedStatements["DistinctPartNumberStmt"] = `select distinct partID from Part order by partID`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) PartAttributeStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Part Attribute Statements
	UnPreparedStatements["GetAttributesByPart"] = `select pa.pAttrID, pa.field,pa.value,pa.sort,pa.partID from PartAttribute as pa
													where pa.partID = ?
													order by pa.sort`
	UnPreparedStatements["UpdatePartAttributeStmt"] = `update PartAttribute as pa
														set pa.field = ?, pa.value = ?, pa.sort = ?
														where pa.pAttrID = ?`
	UnPreparedStatements["UpdatePartAttributeByFieldStmt"] = `update PartAttribute as pa
														set pa.value = ?, pa.sort = ?
														where pa.field = ? && pa.partID = ?`
	UnPreparedStatements["InsertPartAttributeStmt"] = `insert into PartAttribute(partID, field, value, sort)
														values(?,?,?,?)`
	UnPreparedStatements["DeletePartAttributeStmt"] = `delete from PartAttribute where pAttrID = ?`
	UnPreparedStatements["DistinctAttributeValuesStmt"] = `select distinct value from PartAttribute order by value`
	UnPreparedStatements["DistinctAttributeFieldsStmt"] = `select distinct field from PartAttribute order by field`
	UnPreparedStatements["UPCExistStmt"] = `select value from PartAttribute where field = 'UPC' && partID = ? limit 1`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) PartContentStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Part Content Statements
	UnPreparedStatements["GetContentTypeIdStmt"] = `select cTypeID from ContentType where type = ? limit 1`
	UnPreparedStatements["CheckWhileSuppliesLastStmt"] = `select c.contentID from Content as c
															join ContentBridge as cb on c.contentID = cb.contentID
															where cb.partID = ? && c.text = 'While supplies last' limit 1`
	UnPreparedStatements["RemoveWhileSuppliesLastStmt"] = `delete from ContentBridge where contentID = ? and partID = ?`
	UnPreparedStatements["GetWhileSuppliesLastStmt"] = `select contentID from Content where text = 'While supplies last' limit 1`
	UnPreparedStatements["InsertWhileSuppliesLastStmt"] = `insert into ContentBridge (contentID, partID) values(?,?)`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) PartPriceStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Part Price Statements
	UnPreparedStatements["GetPricesByPartStmt"] = `select * from Price where partID = ? order by priceType`
	UnPreparedStatements["UpdatePartPriceStmt"] = `update Price
													set partID = ?, priceType = ?, price = ?, enforced = ?, dateModified = CURDATE()
													where priceID = ?`
	UnPreparedStatements["InsertPartPriceStmt"] = `insert into Price(partID, priceType, price, enforced, dateModified)
													values(?,?,?,?,CURDATE())`
	UnPreparedStatements["DeletePartPriceStmt"] = `delete from Price where priceID = ?`
	UnPreparedStatements["DistinctPriceTypeStmt"] = `select distinct priceType from Price order by priceType`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) RelatedPartStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Related Part Statements
	UnPreparedStatements["GetRelatedPartsStmt"] = `select * from RelatedPart where partID = ? order by relatedID`
	UnPreparedStatements["InsertRelatedPartStmt"] = `insert into RelatedPart(partID, relatedID, rTypeID)
													values(?,?, 0)`
	UnPreparedStatements["DeleteRelatedPartStmt"] = `delete from RelatedPart where relPartID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) PartPackageStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetPartPackagesStmt"] = `select pp.partID, pp.ID, pp.height as height, pp.length as length, pp.width as width, pp.weight as weight, pp.quantity as quantity,
													um_dim.code as dimensionUnit, um_dim.name as dimensionUnitLabel, um_dim.ID as dimensionUnitID,
													um_wt.code as weightUnit, um_wt.name as weightUnitLabel, um_wt.ID as weightUnitID,
													um_pkg.code as packageUnit, um_pkg.name as packageUnitLabel, um_pkg.ID as packageUnitID,
													pt.ID as typeID, pt.name as typeName
													from PartPackage as pp
													join UnitOfMeasure as um_dim on pp.dimensionUOM = um_dim.ID
													join UnitOfMeasure as um_wt on pp.weightUOM = um_wt.ID
													join UnitOfMeasure as um_pkg on pp.packageUOM = um_pkg.ID
													join PackageType as pt on pp.typeID = pt.ID
													where pp.partID = ?`
	UnPreparedStatements["GetPartPackageStmt"] = `select pp.partID, pp.ID, pp.height as height, pp.length as length, pp.width as width, pp.weight as weight, pp.quantity as quantity,
													um_dim.code as dimensionUnit, um_dim.name as dimensionUnitLabel, um_dim.ID as dimensionUnitID,
													um_wt.code as weightUnit, um_wt.name as weightUnitLabel, um_wt.ID as weightUnitID,
													um_pkg.code as packageUnit, um_pkg.name as packageUnitLabel, um_pkg.ID as packageUnitID,
													pt.ID as typeID, pt.name as typeName
													from PartPackage as pp
													join UnitOfMeasure as um_dim on pp.dimensionUOM = um_dim.ID
													join UnitOfMeasure as um_wt on pp.weightUOM = um_wt.ID
													join UnitOfMeasure as um_pkg on pp.packageUOM = um_pkg.ID
													join PackageType as pt on pp.typeID = pt.ID
													where pp.ID = ?`
	UnPreparedStatements["InsertPartPackageStmt"] = `insert into PartPackage
														(partID, height, width, length, weight, dimensionUOM, weightUOM, packageUOM, quantity, typeID)
														values (?,?,?,?,?,?,?,?,?,?)`
	UnPreparedStatements["UpdatePartPackageStmt"] = `update PartPackage
														set height = ?, width = ?, length = ?, weight = ?,
														dimensionUOM = ?, weightUOM = ?, packageUOM = ?,
														quantity = ?, typeID = ?
														where ID = ?`
	UnPreparedStatements["DistinctUnitOfMeasureStmt"] = `select ID, name, code from UnitOfMeasure order by name`
	UnPreparedStatements["DistinctPackageTypeStmt"] = `select ID, name from PackageType order by name`
	UnPreparedStatements["DeletePartPackageStmt"] = `delete from PartPackage where ID = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) CustomerStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Customer Locations
	UnPreparedStatements["GetAllLocationsByCustomerCountStmt"] = `SELECT COUNT(*) from CustomerLocations where cust_id = ?`
	UnPreparedStatements["GetAllCustomerLocationsStmt"] = `SELECT * from CustomerLocations`
	UnPreparedStatements["GetAllCustomerLocationsByPageStmt"] = `SELECT * from CustomerLocations LIMIT ?, ?`
	UnPreparedStatements["GetCustomerLocationsStmt"] = `SELECT * from CustomerLocations WHERE cust_id = ?`
	UnPreparedStatements["GetCustomerLocationsByPageStmt"] = `SELECT * from CustomerLocations WHERE cust_id = ? LIMIT ?, ?`

	UnPreparedStatements["GetCustomerLocationsNoGeoStmt"] = `SELECT * from CustomerLocations WHERE cust_id = ? AND (latitude = 0 OR longitude = 0)`
	UnPreparedStatements["GetCustomerLocationStmt"] = `SELECT * from CustomerLocations WHERE locationID = ?`

	UnPreparedStatements["UpdateCustomerLocationStmt"] = `UPDATE CustomerLocations SET name = ?, address = ?, city = ?,
	stateID = ?, email = ?, phone = ?, fax = ?, latitude = ?, longitude = ?, cust_id = ?, contact_person = ?,
	isprimary = ?, postalCode = ?, ShippingDefault = ? WHERE locationID = ?`
	UnPreparedStatements["AddCustomerLocationStmt"] = `INSERT INTO CustomerLocations (name, address, city, stateID, email, phone, fax, latitude, longitude, cust_id, contact_person, isprimary, postalCode, ShippingDefault) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	UnPreparedStatements["DeleteCustomerLocationStmt"] = `DELETE FROM CustomerLocations WHERE locationID = ?`

	// Customer Users
	UnPreparedStatements["GetAllCustomerUsersStmt"] = `SELECT * from CustomerUser`
	UnPreparedStatements["GetAllCustomerUsersByPageStmt"] = `SELECT * from CustomerUser LIMIT ?, ?`
	UnPreparedStatements["GetAllCustomerUsersCountStmt"] = `SELECT COUNT(*) from CustomerUser`
	UnPreparedStatements["GetCustomerUsersStmt"] = `SELECT * from CustomerUser WHERE cust_id = ?`
	UnPreparedStatements["GetCustomerUsersByPageStmt"] = `SELECT * from CustomerUser WHERE cust_id = ? LIMIT ?, ?`
	UnPreparedStatements["GetCustomerUserStmt"] = `SELECT * from CustomerUser WHERE id = ?`
	UnPreparedStatements["GetUserKeysStmt"] = `select AK.*, AKT.type, AKT.date_added AS typeDateAdded
														from ApiKey AK
														INNER JOIN ApiKeyType AKT ON AK.type_id = AKT.id
														where AK.user_id = ?`
	UnPreparedStatements["GetCustomerUserKeysStmt"] = `select AK.*, AKT.type, AKT.date_added AS typeDateAdded
														from ApiKey AK
														INNER JOIN ApiKeyType AKT ON AK.type_id = AKT.id
														INNER JOIN CustomerUser CU on AK.user_id = CU.id
														where CU.cust_ID = ?`
	UnPreparedStatements["GetAllCustomerUserKeysStmt"] = `select AK.*, AKT.type, AKT.date_added AS typeDateAdded
														from ApiKey AK
														INNER JOIN ApiKeyType AKT ON AK.type_id = AKT.id`
	// Customer
	UnPreparedStatements["GetAllSimpleCustomersStmt"] = `SELECT cust_id, name, customerID FROM Customer`
	UnPreparedStatements["GetAllCustomersStmt"] = `SELECT *,
	                                                (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                                (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                                (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount FROM Customer`
	UnPreparedStatements["GetAllCustomersByPageStmt"] = `SELECT *,
	                                                (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                                (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                                (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
	                                                FROM Customer LIMIT ?, ?`
	UnPreparedStatements["GetAllParentCustomersStmt"] = `SELECT *,
                                                    (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
                                                    (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
                                                    (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
                                                    FROM Customer WHERE parentID <> 0`
	UnPreparedStatements["GetAllParentCustomersByPageStmt"] = `SELECT *,
                                                    (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
                                                    (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
                                                    (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
                                                    FROM Customer WHERE parentID <> 0 LIMIT ?, ?`
	UnPreparedStatements["GetAllCustomersByRepStmt"] = `SELECT *,
	                                                (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                                (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                                (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
	                                                FROM Customer WHERE salesRepID = ?`
	UnPreparedStatements["GetAllCustomersByRepByPageStmt"] = `SELECT *,
	                                                (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                                (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                                (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
	                                                FROM Customer WHERE salesRepID = ? LIMIT ?, ?`
	UnPreparedStatements["GetAllCustomersByDealerTypeStmt"] = `SELECT *,
	                                                (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                                (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                                (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
	                                                FROM Customer WHERE dealer_type = ?`
	UnPreparedStatements["GetAllCustomersByDealerTypeByPageStmt"] = `SELECT *,
	                                                (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                                (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                                (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
	                                                FROM Customer WHERE dealer_type = ? LIMIT ?, ?`
	UnPreparedStatements["GetAllCustomersCountStmt"] = `SELECT COUNT(*) from Customer`
	UnPreparedStatements["GetAllCustomersByRepCountStmt"] = `SELECT COUNT(*) FROM Customer WHERE salesRepID = ?`
	UnPreparedStatements["GetAllCustomersByDealerTypeCountStmt"] = `SELECT COUNT(*) FROM Customer WHERE dealer_type = ?`
	UnPreparedStatements["GetCustomerStmt"] = `SELECT *,
	                                               (SELECT COUNT(locationID) FROM CustomerLocations WHERE cust_id = Customer.cust_id ) AS locationCount,
	                                               (SELECT COUNT(id) FROM CustomerUser WHERE cust_ID = Customer.cust_id ) AS userCount,
	                                               (SELECT COUNT(id) FROM WebProperties WHERE cust_ID = Customer.cust_id ) AS propertyCount
	                                               FROM Customer WHERE cust_id = ?`
	UnPreparedStatements["UpdateCustomerStmt"] = `UPDATE Customer SET name = ?, email = ?, address = ?, address2 = ?, city = ?, stateID = ?, postal_code = ?, phone = ?, fax = ?,
												  contact_person = ?, dealer_type = ?, tier = ?, website = ?, searchURL = ?, eLocalURL = ?, logo = ?, customerID = ?,
												  parentID = ?, isDummy = ?, mCodeID = ?, salesRepID = ?, showWebsite = ? WHERE cust_id = ?`
	UnPreparedStatements["AddCustomerStmt"] = `INSERT INTO Customer (name,email,address,address2,city,stateID,postal_code,phone,fax,contact_person,dealer_type,tier,website,searchURL,eLocalURL,logo,customerID,parentID,isDummy,mCodeID,salesRepID,showWebsite)
											   VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	UnPreparedStatements["DeleteCustomerStmt"] = `DELETE FROM Customer WHERE cust_id = ?`
	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) GeoStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	// Locations
	UnPreparedStatements["GetAllCountriesStmt"] = `SELECT * from Country`
	UnPreparedStatements["GetCountryStmt"] = `SELECT * from Country WHERE countryID = ?`
	UnPreparedStatements["GetStatesByCountryStmt"] = `SELECT * from States WHERE countryID = ?`
	UnPreparedStatements["GetAllStatesStmt"] = `SELECT * from States`
	UnPreparedStatements["GetStateStmt"] = `SELECT * from States WHERE stateID = ?`

	// Dealer Types
	UnPreparedStatements["GetAllDealerTypesStmt"] = `SELECT * From DealerTypes`
	UnPreparedStatements["GetDealerTypeStmt"] = `SELECT * from DealerTypes WHERE dealer_type = ?`

	// Dealer Tier
	UnPreparedStatements["GetAllDealerTiersStmt"] = `SELECT * from DealerTiers`
	UnPreparedStatements["GetDealerTierStmt"] = `SELECT * from DealerTiers WHERE ID = ?`

	// Mapics Codes
	UnPreparedStatements["GetAllMapicsCodesStmt"] = `SELECT * from MapixCode`
	UnPreparedStatements["GetMapicsCodeStmt"] = `SELECT * from MapixCode WHERE mCodeID = ?`

	// MapIcons
	UnPreparedStatements["GetMapIconsStmt"] = `SELECT * from MapIcons`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) TechServicesStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetAllTechNewsStmt"] = `select * from TechNews order by dateModified`
	UnPreparedStatements["GetTechNewsStmt"] = `select * from TechNews where id = ? order by dateModified limit 1`
	UnPreparedStatements["InsertTechNewsStmt"] = `insert into TechNews(pageContent, showDealers, showPublic, dateModified, displayOrder, active, title, subTitle)
													values (?,?,?,?,?,?,?,?)`
	UnPreparedStatements["UpdateTechNewsStmt"] = `update TechNews
													set pageContent = ?, showDealers = ?, showPublic = ?, dateModified = ?, displayOrder = ?, active = ?, title = ?, subTitle = ?
													where id = ?`
	UnPreparedStatements["DeleteTechNewsStmt"] = `delete from TechNews where id = ?`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) FileManagerStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetAllFilesStmt"] = `select
												f.fileID, f.name, f.path, f.height, f.width, f.size, f.createdDate, f.fileGalleryID,
												e.fileExtID, e.fileExt, e.fileExtIcon,
												t.fileTypeID, t.fileType
												from File as f
												left join FileExt as e on f.fileExtID = e.fileExtID
												left join FileType as t on e.fileTypeID = t.fileTypeID
												order by f.name, f.createdDate desc limit 10`
	UnPreparedStatements["GetFileStmt"] = `select
												f.fileID, f.name, f.path, f.size, f.height, f.width, f.createdDate, f.fileGalleryID,
												e.fileExtID, e.fileExt, e.fileExtIcon,
												t.fileTypeID, t.fileType
												from File as f
												left join FileExt as e on f.fileExtID = e.fileExtID
												left join FileType as t on e.fileTypeID = t.fileTypeID
												where f.fileID = ?
												order by f.name, f.createdDate desc`
	UnPreparedStatements["DeleteFileStmt"] = `delete from File where fileID = ?`
	UnPreparedStatements["GetAllGalleriesStmt"] = `select fileGalleryID, name, description, parentID from FileGallery as fg
													order by name`
	UnPreparedStatements["GetGalleriesStmt"] = `select fileGalleryID, name, description, parentID
														from FileGallery as fg
														where parentID = ?
														order by name`
	UnPreparedStatements["GetGalleryFilesStmt"] = `select
													f.fileID, f.name, f.size, f.path, f.height, f.width,
													f.createdDate, f.fileGalleryID,
													e.fileExtID, e.fileExt, e.fileExtIcon,
													t.fileTypeID, t.fileType
													from File as f
													left join FileExt as e on f.fileExtID = e.fileExtID
													left join FileType as t on e.fileTypeID = t.fileTypeID
													where f.fileGalleryID = ?
													order by f.name, f.createdDate desc`
	UnPreparedStatements["GetGalleryStmt"] = `select fileGalleryID, name, description, parentID
												from FileGallery as fg
												where fileGalleryID = ?
												order by name`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) VehicleStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["GetVehicleYearsStmt"] = `select year from Year`
	UnPreparedStatements["GetVehicleMakesStmt"] = `select ma.makeID, ma.make from Make as ma
													left join YearMake as ym on ma.makeID = ym.makeID
													left join Year as y on ym.yearID = y.yearID
													where y.year = ?
													order by ma.make`
	UnPreparedStatements["GetAllVehicleMakesStmt"] = `select ma.makeID, ma.make from Make as ma
													order by ma.make`
	UnPreparedStatements["GetVehicleModelsStmt"] = `select mo.modelID, mo.model from Model as mo
													left join MakeModel as mm on mo.modelID = mm.modelID
													left join Make as ma on mm.makeID = ma.makeID
													left join YearMake as ym on ma.makeID = ym.makeID
													left join Year as y on ym.yearID = y.yearID
													where y.year = ? && ma.make = ?
													order by mo.model`
	UnPreparedStatements["GetAllVehicleModelsStmt"] = `select mo.modelID, mo.model from Model as mo order by mo.model`
	UnPreparedStatements["GetVehicleStylesStmt"] = `select s.styleID,s.style from Style as s
													left join ModelStyle as sm on s.styleID = sm.styleID
													left join Model as mo on sm.modelID = mo.modelID
													left join MakeModel as mm on mo.modelID = mm.modelID
													left join Make as ma on mm.makeID = ma.makeID
													left join YearMake as ym on ma.makeID = ym.makeID
													left join Year as y on ym.yearID = y.yearID
													where y.year = ? && ma.make = ? && mo.model = ?
													order by mo.model`
	UnPreparedStatements["GetAllVehicleStylesStmt"] = `select s.styleID, s.style from Style as s order by s.style`
	UnPreparedStatements["GetVehiclesByPartStmt"] = `select distinct v.vehicleID, v.yearID, v.makeID, v.modelID, v.styleID, y.year, ma.make, mo.model, s.style,
													(
														select count(vp2.vPartID) from VehiclePart as vp2
														where vp2.vehicleID = v.vehicleID
													) as partCount from Vehicle as v
													left join Style as s on v.styleID = s.styleID
													left join Model as mo on v.modelID = mo.modelID
													left join Make as ma on v.makeID = ma.makeID
													left join Year as y on v.yearID = y.yearID
													left join VehiclePart as vp on v.vehicleID = vp.vehicleID
													where vp.partID = ?
													order by y.year, ma.make,mo.model,s.style`
	UnPreparedStatements["GetVehiclesStmt"] = `select distinct v.vehicleID, v.yearID, v.makeID, v.modelID, v.styleID, y.year, ma.make, mo.model, s.style,
													(
														select count(vp2.vPartID) from VehiclePart as vp2
														where vp2.vehicleID = v.vehicleID
													) as partCount from Vehicle as v
													left join Style as s on v.styleID = s.styleID
													left join Model as mo on v.modelID = mo.modelID
													left join Make as ma on v.makeID = ma.makeID
													left join Year as y on v.yearID = y.yearID
													left join VehiclePart as vp on v.vehicleID = vp.vehicleID
													order by y.year, ma.make,mo.model,s.style`
	UnPreparedStatements["GetVehiclesByMakeModelStmt"] = `select distinct v.vehicleID, v.yearID, v.makeID, v.modelID, v.styleID, y.year, ma.make, mo.model, s.style,
													(
														select count(vp2.vPartID) from VehiclePart as vp2
														where vp2.vehicleID = v.vehicleID
													) as partCount from Vehicle as v
													left join Style as s on v.styleID = s.styleID
													left join Model as mo on v.modelID = mo.modelID
													left join Make as ma on v.makeID = ma.makeID
													left join Year as y on v.yearID = y.yearID
													left join VehiclePart as vp on v.vehicleID = vp.vehicleID
													where ma.make = ? && mo.model = ?
													order by y.year, ma.make,mo.model,s.style`
	UnPreparedStatements["GetVehiclesByMakeStmt"] = `select distinct v.vehicleID, v.yearID, v.makeID, v.modelID, v.styleID, y.year, ma.make, mo.model, s.style,
													(
														select count(vp2.vPartID) from VehiclePart as vp2
														where vp2.vehicleID = v.vehicleID
													) as partCount from Vehicle as v
													left join Style as s on v.styleID = s.styleID
													left join Model as mo on v.modelID = mo.modelID
													left join Make as ma on v.makeID = ma.makeID
													left join Year as y on v.yearID = y.yearID
													left join VehiclePart as vp on v.vehicleID = vp.vehicleID
													where ma.make = ?
													order by y.year, ma.make,mo.model,s.style`
	UnPreparedStatements["GetVehicleByIDStmt"] = `select distinct v.vehicleID, v.yearID, v.makeID, v.modelID, v.styleID, y.year, ma.make, mo.model, s.style,
													(
														select count(vp2.vPartID) from VehiclePart as vp2
														where vp2.vehicleID = v.vehicleID
													) as partCount from Vehicle as v
													left join Style as s on v.styleID = s.styleID
													left join Model as mo on v.modelID = mo.modelID
													left join Make as ma on v.makeID = ma.makeID
													left join Year as y on v.yearID = y.yearID
													where v.vehicleID = ?
													order by y.year, ma.make,mo.model,s.style limit 1`
	UnPreparedStatements["GetVehicleByIDWithPartStmt"] = `select distinct v.vehicleID, v.yearID, v.makeID, v.modelID, v.styleID, y.year, ma.make, mo.model, s.style,
													(
														select count(vp2.vPartID) from VehiclePart as vp2
														where vp2.vehicleID = v.vehicleID
													) as partCount,
													vp.installTime, vp.drilling, vp.exposed from Vehicle as v
													left join Style as s on v.styleID = s.styleID
													left join Model as mo on v.modelID = mo.modelID
													left join Make as ma on v.makeID = ma.makeID
													left join Year as y on v.yearID = y.yearID
													left join VehiclePart as vp on v.vehicleID = vp.vehicleID && vp.partID = ?
													where v.vehicleID = ?
													order by y.year, ma.make,mo.model,s.style limit 1`
	UnPreparedStatements["GetDistinctDrillingNote"] = `select distinct vp.drilling from VehiclePart vp order by vp.drilling`
	UnPreparedStatements["GetDistinctExposedNote"] = `select distinct vp.exposed from VehiclePart vp order by vp.exposed`
	UnPreparedStatements["GetDistinctInstallTime"] = `select distinct vp.installTime from VehiclePart vp order by vp.installTime`
	UnPreparedStatements["InsertVehiclePartStmt"] = `insert into VehiclePart (partID, vehicleID, drilling, exposed, installTime)
														values(?,?,?,?,?)`
	UnPreparedStatements["UpdateVehiclePartStmt"] = `update VehiclePart
														set drilling = ?, exposed = ?, installTime = ?
														where partID = ? && vehicleID = ?`
	UnPreparedStatements["DeleteVehiclePartStmt"] = `delete from VehiclePart where vehicleID = ? && partID = ?`
	UnPreparedStatements["GetVehiclePartsStmt"] = `select distinct partID from VehiclePart where vehicleID = ? order by partID`
	UnPreparedStatements["InsertVehicleStmt"] = `insert into Vehicle(yearID,makeID, modelID, styleID)
													values(
													(select yearID from Year where year = ?),
													(select ma.makeID from Make ma where ma.make = ?),
													(select modelID from Model where model = ?),
													(select styleID from Style where style = ?))`
	UnPreparedStatements["UpdateVehicleStmt"] = `update Vehicle v
													set v.yearID = (select yearID from Year where year = ?),
													v.makeID = (select makeID from Make where make = ?),
													v.modelID = (select modelID from Model where model = ?),
													v.styleID = (select styleID from Style where style = ?)
													where v.vehicleID = ?`
	UnPreparedStatements["DeleteVehicleStmt"] = `delete from Vehicle where vehicleID = ?`

	UnPreparedStatements["GetYearStmt"] = `select yearID from Year where year = ? limit 1`
	UnPreparedStatements["GetMakeStmt"] = `select makeID from Make where make = ? limit 1`
	UnPreparedStatements["GetModelStmt"] = `select modelID from Model where model = ? limit 1`
	UnPreparedStatements["GetStyleStmt"] = `select styleID from Style where style = ? limit 1`

	UnPreparedStatements["InsertYearStmt"] = `insert into Year(year) values(?)`
	UnPreparedStatements["InsertMakeStmt"] = `insert into Make(make) values(?)`
	UnPreparedStatements["InsertModelStmt"] = `insert into Model(model) values(?)`
	UnPreparedStatements["InsertStyleStmt"] = `insert into Style(style, aaiaID) values(?,0)`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func (t *T) AcesStatements() {
	UnPreparedStatements := make(map[string]string, 0)

	UnPreparedStatements["aces_GetConfigurationOptionsByMakeModel"] = `select cat.name, ca.value from vcdb_Vehicle as v
																																			join BaseVehicle as bv on v.BaseVehicleID = bv.ID
																																			join VehicleConfigAttribute as vca on v.ConfigID = vca.VehicleConfigID
																																			join ConfigAttribute as ca on vca.AttributeID = ca.ID
																																			join ConfigAttributeType as cat on ca.ConfigAttributeTypeID = cat.ID
																																			join vcdb.BaseVehicle as vbv on bv.AAIABaseVehicleID = vbv.BaseVehicleID
																																			where vbv.MakeID = ? && vbv.ModelID = ?
																																			group by ca.value`
	UnPreparedStatements["aces_GetMakes"] = `select ma.ID, ma.MakeName, ma.AAIAMakeID from vcdb_Make as ma
																						order by ma.MakeName`
	UnPreparedStatements["aces_GetMakeByName"] = `select ma.ID, ma.AAIAMakeID, ma.MakeName from vcdb_Make as ma where LOWER(ma.MakeName) = ? limit 1`
	UnPreparedStatements["aces_GetMake"] = `select ma.ID, ma.AAIAMakeID, ma.MakeName from vcdb_Make as ma where ma.ID = ? limit 1`
	UnPreparedStatements["aces_GetMakeByAAIA"] = `select ma.ID, ma.AAIAMakeID, ma.MakeName from vcdb_Make as ma where ma.AAIAMakeID = ? limit 1`
	UnPreparedStatements["aces_AddMake"] = `insert into vcdb_Make(AAIAMakeID, MakeName) values(?,?)`
	UnPreparedStatements["aces_UpdateMake"] = `update vcdb_Make
																							set MakeName = ?
																							where ID = ?`
	UnPreparedStatements["aces_RemoveMake"] = `delete from vcdb_Make where ID = ?`
	UnPreparedStatements["aces_BaseVehicleCountByMake"] = `select count(ID) as count from BaseVehicle where MakeID = ?`
	UnPreparedStatements["aces_GetYear"] = `select YearID from vcdb_Year order by YearID`
	UnPreparedStatements["aces_GetYear"] = `select YearID from vcdb_Year where YearID = ? limit 1`
	UnPreparedStatements["aces_AddYear"] = `insert into vcdb_Year(YearID) values(?)`
	UnPreparedStatements["aces_DeleteYear"] = `delete from vcdb_Year where YearID = ?`
	UnPreparedStatements["aces_UpdateBaseVehicleMake"] = `update BaseVehicle as bv
																												set bv.MakeID = ?
																												where bv.MakeID = ?`
	UnPreparedStatements["aces_GetConfigurationOptionsByType"] = `select distinct ca.ID, ca.value from ConfigAttribute as ca
																																join ConfigAttributeType as cat on ca.ConfigAttributeTypeID = cat.ID
																																where LOWER(cat.name) = LOWER(?)`

	if !CurtDevDb.Raw.IsConnected() {
		CurtDevDb.Raw.Connect()
	}

	c := make(chan int)

	for stmtname, stmtsql := range UnPreparedStatements {
		go PrepareCurtDevStatement(stmtname, stmtsql, c)
	}

	for _, _ = range UnPreparedStatements {
		<-c
	}
}

func PrepareAdminStatement(name string, sql string, ch chan int) {
	stmt, err := AdminDb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
}

func PreparePcdbStatement(name string, sql string, ch chan int) {
	stmt, err := Pcdb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
}

func PrepareVcdbStatement(name string, sql string, ch chan int) {
	stmt, err := Vcdb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
}

func PrepareCurtDevStatement(name string, sql string, ch chan int) {
	stmt, err := CurtDevDb.Prepare(sql)
	if err == nil {
		Statements[name] = stmt
	} else {
		log.Println(err)
	}
	ch <- 1
}

func GetStatement(key string) (stmt *autorc.Stmt, err error) {
	stmt, ok := Statements[key]
	if !ok {
		qry := expvar.Get(key)
		if qry == nil {
			err = errors.New("Invalid query reference")
		}
	}
	return
}
