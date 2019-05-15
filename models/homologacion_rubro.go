package models

import (
	"errors"
	"strconv"
	"fmt"
	"reflect"
	"strings"
	"github.com/udistrital/utils_oas/request"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

// HomologacionRubro ...
type HomologacionRubro struct {
	Id               int64  `orm:"auto"`
	CodigoHomologado string `orm:"size(128)"`
	NombreHomologado string `orm:"size(128)"`
	Organizacion     int
	Vigencia         float64
}

func init() {
	orm.RegisterModel(new(HomologacionRubro))
}

// AddHomologacionRubro insert a new HomologacionRubro into database and returns
// last inserted Id on success.
func AddHomologacionRubro(m *HomologacionRubro) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHomologacionRubroByID retrieves HomologacionRubro by Id. Returns error if
// Id doesn't exist
func GetHomologacionRubroByID(id int64) (v *HomologacionRubro, err error) {
	o := orm.NewOrm()
	v = &HomologacionRubro{Id: id}
	if err = o.QueryTable(new(HomologacionRubro)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHomologacionRubro retrieves all HomologacionRubro matches certain condition. Returns empty list if
// no records exist
func GetAllHomologacionRubro(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(HomologacionRubro))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []HomologacionRubro
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateHomologacionRubro ...
// updates HomologacionRubro by Id and returns error if the record to be updated doesn't exist
func UpdateHomologacionRubroByID(m *HomologacionRubro) (err error) {
	o := orm.NewOrm()
	v := HomologacionRubro{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHomologacionRubro deletes HomologacionRubro by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHomologacionRubro(id int64) (err error) {
	o := orm.NewOrm()
	v := HomologacionRubro{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&HomologacionRubro{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func formatoHomologacionF(tipo int) (function func(data interface{}, params ...interface{}) interface{}){
	switch tipo {
	case 1:
		return getOrganizacion
	case 2:
	  return getOrganizacionDisponibilidad
	default:
		return nil
	}
}


func getOrganizacionDisponibilidad (data interface{},params ...interface{}) (res interface{}) {
	var organizacion interface{}
	var rubroHomolMap map[string]interface{}
	var cntRubroHomologado map[string]interface{}

	rubroHomolMap = data.(map[string]interface{})

	idOrganizacion := strconv.FormatFloat(rubroHomolMap["Organizacion"].(float64),'f',-1,64)
		if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?limit=-1&query=Id:"+idOrganizacion, &organizacion); err == nil {
			rubroHomolMap["Organizacion"]= organizacion
			res = rubroHomolMap
		}else {
			beego.Error("error",err)
		}
		if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/rubro_homologado/GetRecordsNumberRubroHomologadoById/"+strconv.FormatFloat(rubroHomolMap["Id"].(float64),'f',-1,64),&cntRubroHomologado) ; err == nil {
		if cntRubroHomologado["Body"].(float64)>= 1 {
				rubroHomolMap["Disponibilidad"]=false
			}else{
				rubroHomolMap["Disponibilidad"]=true
			}
		}else{
			beego.Error("error",err)
		}
	return
}

func getOrganizacion (data interface{},params ...interface{}) (res interface{}) {
	var organizacion interface{}
	var rubroHomolMap map[string]interface{}

	rubroHomolMap = data.(map[string]interface{})["RubroHomologado"].(map[string]interface{})

	idOrganizacion := strconv.FormatFloat(rubroHomolMap["Organizacion"].(float64),'f',-1,64)
		if err := request.GetJson(beego.AppConfig.String("coreOrganizacionService")+"organizacion?limit=-1&query=Id:"+idOrganizacion, &organizacion); err == nil {
			rubroHomolMap["Organizacion"]= organizacion
			res = rubroHomolMap
		}

	return
}

// GetOrganizacionRubroHomologado ...
// getValue from organizacion for all rows
func GetOrganizacionRubroHomologado(rubroHomol interface{},params ...interface{})(res interface{}){
 var tipo int
	if rubroHomol.(map[string]interface{})["RubroHomologado"]!= nil{
		tipo = 1
	}else{
		tipo = 2
	}
	beego.Error("rubro homologado",rubroHomol)
 if function := formatoHomologacionF(tipo); function != nil{
 		res = function(rubroHomol,params)
		beego.Error("respuesta",res)
 } else {
	 res = rubroHomol
	 beego.Error("respuesta vacio",res)
 }
	return
}