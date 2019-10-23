package eel

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var AREA = make(map[string][]map[string]interface{})
var provinces = make([]*Province, 0)
var name2Province = make(map[string]*Province)
var id2Province = make(map[int]*Province)
var name2City = make(map[string]*City)
var id2City = make(map[int]*City)
//var cities = make([]*City, 0)
//var cityId2districts = make(map[string][]*District)
var id2District = make(map[int]*District)
var name2District = make(map[string]*District)

type Province struct {
	Id int
	Name string
	Cities []*City
}

// IsDGC 是否是直辖市
func (this *Province) IsDGC() bool {
	return len(this.Cities) == 1
}

func NewProvince(data map[string]interface{}) *Province{
	province := new(Province)
	province.Id, _ = strconv.Atoi(data["id"].(string))
	province.Name = data["name"].(string)
	province.Cities = make([]*City, 0)
	return province
}

type City struct {
	Id int
	ProvinceId int
	Name string
	ZipCode string
	Districts []*District
}

func NewCity(data map[string]interface{}) *City{
	city := new(City)
	city.Id, _ = strconv.Atoi(data["id"].(string))
	city.ProvinceId, _ = strconv.Atoi(data["province_id"].(string))
	city.Name = data["name"].(string)
	city.ZipCode = data["zip_code"].(string)
	city.Districts = make([]*District, 0)
	return city
}

type District struct {
	Id int
	CityId int
	Name string
}

func NewDistrict(data map[string]interface{}) *District{
	district := new(District)
	district.Id, _ = strconv.Atoi(data["id"].(string))
	district.CityId, _ = strconv.Atoi(data["city_id"].(string))
	district.Name = data["name"].(string)
	return district
}

type Area struct {
	Province *Province
	City *City
	District *District
}

type AreaService struct {
	ServiceBase
}

func NewAreaService() *AreaService {
	service := new(AreaService)
	return service
}

/*
  Province相关api
 */

func (this *AreaService) GetProvinces() []*Province {
	return provinces
}

func (this *AreaService) GetProvinceByName(name string) *Province{
	return name2Province[name]
}

func (this *AreaService) GetProvincesByNames(names []string) []*Province {
	provinces := make([]*Province, 0)
	for _, name := range names {
		if province, ok := name2Province[name]; ok {
			provinces = append(provinces, province)
		}
	}
	
	return provinces
}


func (this *AreaService) GetProvinceById(id int) *Province {
	return id2Province[id]
}

func (this *AreaService) GetProvincesByIds(ids []int) []*Province {
	provinces := make([]*Province, 0)
	for _, id := range ids {
		if province, ok := id2Province[id]; ok {
			provinces = append(provinces, province)
		}
	}
	
	return provinces
}

/*
  City相关api
  */

func (this *AreaService) GetCityByName(name string) *City{
	return name2City[name]
}

func (this *AreaService) GetCitiesByNames(names []string) []*City {
	cities := make([]*City, 0)
	for _, name := range names {
		if city, ok := name2City[name]; ok {
			cities = append(cities, city)
		}
	}
	
	return cities
}

func (this *AreaService) GetCityById(id int) *City{
	return id2City[id]
}

func (this *AreaService) GetCitiesByIds(ids []int) []*City {
	cities := make([]*City, 0)
	for _, id := range ids {
		if city, ok := id2City[id]; ok {
			cities = append(cities, city)
		}
	}
	
	return cities
}

func (this *AreaService) GetCitiesForProvince(provinceId int) []*City {
	if province, ok := id2Province[provinceId]; ok {
		return province.Cities
	} else {
		return make([]*City, 0)
	}
}

/*
 District相关api
 */
func (this *AreaService) GetDistrictByName(cityId int, name string) *District{
	name = fmt.Sprintf("%d_%s", cityId, name)
	return name2District[name]
}

func (this *AreaService) GetDistrictsByNames(cityId int, names []string) []*District{
	districts := make([]*District, 0)
	for _, name := range names {
		_name := fmt.Sprintf("%d_%s", cityId, name)
		if district, ok := name2District[_name]; ok {
			districts = append(districts, district)
		}
	}
	
	return districts
}

func (this *AreaService) GetDistrictById(id int) *District{
	return id2District[id]
}

func (this *AreaService) GetDistrictsByIds(ids []int) []*District {
	districts := make([]*District, 0)
	for _, id := range ids {
		if district, ok := id2District[id]; ok {
			districts = append(districts, district)
		}
	}
	
	return districts
}

func (this *AreaService) GetDistrictsForCity(cityId int) []*District {
	if city, ok := id2City[cityId]; ok {
		return city.Districts
	} else {
		return make([]*District, 0)
	}
}

/*
 Area相关api
 */

// GetAreaByName 根据area name(北京市 北京市 东城区)获得area
func (this *AreaService) GetAreaByName(name string) *Area {
	var items []string
	if strings.Contains(name, "-"){
		items = strings.Split(name, "-")
	} else {
		items = strings.Split(name, " ")
	}

	province := this.GetProvinceByName(items[0])
	city := this.GetCityByName(items[1])
	district := this.GetDistrictByName(city.Id, items[2])
	
	return &Area{
		Province: province,
		City: city,
		District: district,
	}
}

// GetAreaByCode 根据area code(1_1_1)获取area
func (this *AreaService) GetAreaByCode(code string) *Area {
	items := strings.Split(code, "_")
	
	provinceId, _ := strconv.Atoi(items[0])
	province := this.GetProvinceById(provinceId)
	
	cityId, _ := strconv.Atoi(items[1])
	city := this.GetCityById(cityId)
	
	districtId, _ := strconv.Atoi(items[2])
	district := this.GetDistrictById(districtId)
	
	return &Area{
		Province: province,
		City: city,
		District: district,
	}
}


func init(){
	json.Unmarshal([]byte(areaJSON), &AREA)

	for name, arrs := range AREA{
		switch name {
		case "PROVINCES":
			for _, data := range arrs {
				province := NewProvince(data)
				provinces = append(provinces, province)
				name2Province[province.Name] = province
				id2Province[province.Id] = province
			}

		case "CITIES":
			for _, data := range arrs {
				city := NewCity(data)
				name2City[city.Name] = city
				id2City[city.Id] = city
				//cities = append(cities, city)
				
				//向province.Cities中加入city
				if province, ok := id2Province[city.ProvinceId]; ok {
					province.Cities = append(province.Cities, city)
				}
			}
		case "DISTRICTS":
			for _, data := range arrs {
				district := NewDistrict(data)
				name := fmt.Sprintf("%d_%s", district.CityId, district.Name)
				name2District[name] = district
				id2District[district.Id] = district
				
				if city, ok := id2City[district.CityId]; ok {
					city.Districts = append(city.Districts, district)
				}
			}
		}
	}
}

const areaJSON = `{
	"PROVINCES": [
		{
			"id": "1",
			"name": "北京市"
		},
		{
			"id": "2",
			"name": "天津市"
		},
		{
			"id": "3",
			"name": "河北省"
		},
		{
			"id": "4",
			"name": "山西省"
		},
		{
			"id": "5",
			"name": "内蒙古自治区"
		},
		{
			"id": "6",
			"name": "辽宁省"
		},
		{
			"id": "7",
			"name": "吉林省"
		},
		{
			"id": "8",
			"name": "黑龙江省"
		},
		{
			"id": "9",
			"name": "上海市"
		},
		{
			"id": "10",
			"name": "江苏省"
		},
		{
			"id": "11",
			"name": "浙江省"
		},
		{
			"id": "12",
			"name": "安徽省"
		},
		{
			"id": "13",
			"name": "福建省"
		},
		{
			"id": "14",
			"name": "江西省"
		},
		{
			"id": "15",
			"name": "山东省"
		},
		{
			"id": "16",
			"name": "河南省"
		},
		{
			"id": "17",
			"name": "湖北省"
		},
		{
			"id": "18",
			"name": "湖南省"
		},
		{
			"id": "19",
			"name": "广东省"
		},
		{
			"id": "20",
			"name": "广西壮族自治区"
		},
		{
			"id": "21",
			"name": "海南省"
		},
		{
			"id": "22",
			"name": "重庆市"
		},
		{
			"id": "23",
			"name": "四川省"
		},
		{
			"id": "24",
			"name": "贵州省"
		},
		{
			"id": "25",
			"name": "云南省"
		},
		{
			"id": "26",
			"name": "西藏自治区"
		},
		{
			"id": "27",
			"name": "陕西省"
		},
		{
			"id": "28",
			"name": "甘肃省"
		},
		{
			"id": "29",
			"name": "青海省"
		},
		{
			"id": "30",
			"name": "宁夏回族自治区"
		},
		{
			"id": "31",
			"name": "新疆维吾尔自治区"
		},
		{
			"id": "32",
			"name": "香港特别行政区"
		},
		{
			"id": "33",
			"name": "澳门特别行政区"
		},
		{
			"id": "34",
			"name": "台湾省"
		}
	],
	"CITIES": [
		{
			"id": "1",
			"province_id": "1",
			"name": "北京市",
			"zip_code": "100000"
		},
		{
			"id": "2",
			"province_id": "2",
			"name": "天津市",
			"zip_code": "100000"
		},
		{
			"id": "3",
			"province_id": "3",
			"name": "石家庄市",
			"zip_code": "050000"
		},
		{
			"id": "4",
			"province_id": "3",
			"name": "唐山市",
			"zip_code": "063000"
		},
		{
			"id": "5",
			"province_id": "3",
			"name": "秦皇岛市",
			"zip_code": "066000"
		},
		{
			"id": "6",
			"province_id": "3",
			"name": "邯郸市",
			"zip_code": "056000"
		},
		{
			"id": "7",
			"province_id": "3",
			"name": "邢台市",
			"zip_code": "054000"
		},
		{
			"id": "8",
			"province_id": "3",
			"name": "保定市",
			"zip_code": "071000"
		},
		{
			"id": "9",
			"province_id": "3",
			"name": "张家口市",
			"zip_code": "075000"
		},
		{
			"id": "10",
			"province_id": "3",
			"name": "承德市",
			"zip_code": "067000"
		},
		{
			"id": "11",
			"province_id": "3",
			"name": "沧州市",
			"zip_code": "061000"
		},
		{
			"id": "12",
			"province_id": "3",
			"name": "廊坊市",
			"zip_code": "065000"
		},
		{
			"id": "13",
			"province_id": "3",
			"name": "衡水市",
			"zip_code": "053000"
		},
		{
			"id": "14",
			"province_id": "4",
			"name": "太原市",
			"zip_code": "030000"
		},
		{
			"id": "15",
			"province_id": "4",
			"name": "大同市",
			"zip_code": "037000"
		},
		{
			"id": "16",
			"province_id": "4",
			"name": "阳泉市",
			"zip_code": "045000"
		},
		{
			"id": "17",
			"province_id": "4",
			"name": "长治市",
			"zip_code": "046000"
		},
		{
			"id": "18",
			"province_id": "4",
			"name": "晋城市",
			"zip_code": "048000"
		},
		{
			"id": "19",
			"province_id": "4",
			"name": "朔州市",
			"zip_code": "036000"
		},
		{
			"id": "20",
			"province_id": "4",
			"name": "晋中市",
			"zip_code": "030600"
		},
		{
			"id": "21",
			"province_id": "4",
			"name": "运城市",
			"zip_code": "044000"
		},
		{
			"id": "22",
			"province_id": "4",
			"name": "忻州市",
			"zip_code": "034000"
		},
		{
			"id": "23",
			"province_id": "4",
			"name": "临汾市",
			"zip_code": "041000"
		},
		{
			"id": "24",
			"province_id": "4",
			"name": "吕梁市",
			"zip_code": "030500"
		},
		{
			"id": "25",
			"province_id": "5",
			"name": "呼和浩特市",
			"zip_code": "010000"
		},
		{
			"id": "26",
			"province_id": "5",
			"name": "包头市",
			"zip_code": "014000"
		},
		{
			"id": "27",
			"province_id": "5",
			"name": "乌海市",
			"zip_code": "016000"
		},
		{
			"id": "28",
			"province_id": "5",
			"name": "赤峰市",
			"zip_code": "024000"
		},
		{
			"id": "29",
			"province_id": "5",
			"name": "通辽市",
			"zip_code": "028000"
		},
		{
			"id": "30",
			"province_id": "5",
			"name": "鄂尔多斯市",
			"zip_code": "010300"
		},
		{
			"id": "31",
			"province_id": "5",
			"name": "呼伦贝尔市",
			"zip_code": "021000"
		},
		{
			"id": "32",
			"province_id": "5",
			"name": "巴彦淖尔市",
			"zip_code": "014400"
		},
		{
			"id": "33",
			"province_id": "5",
			"name": "乌兰察布市",
			"zip_code": "011800"
		},
		{
			"id": "34",
			"province_id": "5",
			"name": "兴安盟",
			"zip_code": "137500"
		},
		{
			"id": "35",
			"province_id": "5",
			"name": "锡林郭勒盟",
			"zip_code": "011100"
		},
		{
			"id": "36",
			"province_id": "5",
			"name": "阿拉善盟",
			"zip_code": "016000"
		},
		{
			"id": "37",
			"province_id": "6",
			"name": "沈阳市",
			"zip_code": "110000"
		},
		{
			"id": "38",
			"province_id": "6",
			"name": "大连市",
			"zip_code": "116000"
		},
		{
			"id": "39",
			"province_id": "6",
			"name": "鞍山市",
			"zip_code": "114000"
		},
		{
			"id": "40",
			"province_id": "6",
			"name": "抚顺市",
			"zip_code": "113000"
		},
		{
			"id": "41",
			"province_id": "6",
			"name": "本溪市",
			"zip_code": "117000"
		},
		{
			"id": "42",
			"province_id": "6",
			"name": "丹东市",
			"zip_code": "118000"
		},
		{
			"id": "43",
			"province_id": "6",
			"name": "锦州市",
			"zip_code": "121000"
		},
		{
			"id": "44",
			"province_id": "6",
			"name": "营口市",
			"zip_code": "115000"
		},
		{
			"id": "45",
			"province_id": "6",
			"name": "阜新市",
			"zip_code": "123000"
		},
		{
			"id": "46",
			"province_id": "6",
			"name": "辽阳市",
			"zip_code": "111000"
		},
		{
			"id": "47",
			"province_id": "6",
			"name": "盘锦市",
			"zip_code": "124000"
		},
		{
			"id": "48",
			"province_id": "6",
			"name": "铁岭市",
			"zip_code": "112000"
		},
		{
			"id": "49",
			"province_id": "6",
			"name": "朝阳市",
			"zip_code": "122000"
		},
		{
			"id": "50",
			"province_id": "6",
			"name": "葫芦岛市",
			"zip_code": "125000"
		},
		{
			"id": "51",
			"province_id": "7",
			"name": "长春市",
			"zip_code": "130000"
		},
		{
			"id": "52",
			"province_id": "7",
			"name": "吉林市",
			"zip_code": "132000"
		},
		{
			"id": "53",
			"province_id": "7",
			"name": "四平市",
			"zip_code": "136000"
		},
		{
			"id": "54",
			"province_id": "7",
			"name": "辽源市",
			"zip_code": "136200"
		},
		{
			"id": "55",
			"province_id": "7",
			"name": "通化市",
			"zip_code": "134000"
		},
		{
			"id": "56",
			"province_id": "7",
			"name": "白山市",
			"zip_code": "134300"
		},
		{
			"id": "57",
			"province_id": "7",
			"name": "松原市",
			"zip_code": "131100"
		},
		{
			"id": "58",
			"province_id": "7",
			"name": "白城市",
			"zip_code": "137000"
		},
		{
			"id": "59",
			"province_id": "7",
			"name": "延边朝鲜族自治州",
			"zip_code": "133000"
		},
		{
			"id": "60",
			"province_id": "8",
			"name": "哈尔滨市",
			"zip_code": "150000"
		},
		{
			"id": "61",
			"province_id": "8",
			"name": "齐齐哈尔市",
			"zip_code": "161000"
		},
		{
			"id": "62",
			"province_id": "8",
			"name": "鸡西市",
			"zip_code": "158100"
		},
		{
			"id": "63",
			"province_id": "8",
			"name": "鹤岗市",
			"zip_code": "154100"
		},
		{
			"id": "64",
			"province_id": "8",
			"name": "双鸭山市",
			"zip_code": "155100"
		},
		{
			"id": "65",
			"province_id": "8",
			"name": "大庆市",
			"zip_code": "163000"
		},
		{
			"id": "66",
			"province_id": "8",
			"name": "伊春市",
			"zip_code": "152300"
		},
		{
			"id": "67",
			"province_id": "8",
			"name": "佳木斯市",
			"zip_code": "154000"
		},
		{
			"id": "68",
			"province_id": "8",
			"name": "七台河市",
			"zip_code": "154600"
		},
		{
			"id": "69",
			"province_id": "8",
			"name": "牡丹江市",
			"zip_code": "157000"
		},
		{
			"id": "70",
			"province_id": "8",
			"name": "黑河市",
			"zip_code": "164300"
		},
		{
			"id": "71",
			"province_id": "8",
			"name": "绥化市",
			"zip_code": "152000"
		},
		{
			"id": "72",
			"province_id": "8",
			"name": "大兴安岭地区",
			"zip_code": "165000"
		},
		{
			"id": "73",
			"province_id": "9",
			"name": "上海市",
			"zip_code": "200000"
		},
		{
			"id": "74",
			"province_id": "10",
			"name": "南京市",
			"zip_code": "210000"
		},
		{
			"id": "75",
			"province_id": "10",
			"name": "无锡市",
			"zip_code": "214000"
		},
		{
			"id": "76",
			"province_id": "10",
			"name": "徐州市",
			"zip_code": "221000"
		},
		{
			"id": "77",
			"province_id": "10",
			"name": "常州市",
			"zip_code": "213000"
		},
		{
			"id": "78",
			"province_id": "10",
			"name": "苏州市",
			"zip_code": "215000"
		},
		{
			"id": "79",
			"province_id": "10",
			"name": "南通市",
			"zip_code": "226000"
		},
		{
			"id": "80",
			"province_id": "10",
			"name": "连云港市",
			"zip_code": "222000"
		},
		{
			"id": "81",
			"province_id": "10",
			"name": "淮安市",
			"zip_code": "223200"
		},
		{
			"id": "82",
			"province_id": "10",
			"name": "盐城市",
			"zip_code": "224000"
		},
		{
			"id": "83",
			"province_id": "10",
			"name": "扬州市",
			"zip_code": "225000"
		},
		{
			"id": "84",
			"province_id": "10",
			"name": "镇江市",
			"zip_code": "212000"
		},
		{
			"id": "85",
			"province_id": "10",
			"name": "泰州市",
			"zip_code": "225300"
		},
		{
			"id": "86",
			"province_id": "10",
			"name": "宿迁市",
			"zip_code": "223800"
		},
		{
			"id": "87",
			"province_id": "11",
			"name": "杭州市",
			"zip_code": "310000"
		},
		{
			"id": "88",
			"province_id": "11",
			"name": "宁波市",
			"zip_code": "315000"
		},
		{
			"id": "89",
			"province_id": "11",
			"name": "温州市",
			"zip_code": "325000"
		},
		{
			"id": "90",
			"province_id": "11",
			"name": "嘉兴市",
			"zip_code": "314000"
		},
		{
			"id": "91",
			"province_id": "11",
			"name": "湖州市",
			"zip_code": "313000"
		},
		{
			"id": "92",
			"province_id": "11",
			"name": "绍兴市",
			"zip_code": "312000"
		},
		{
			"id": "93",
			"province_id": "11",
			"name": "金华市",
			"zip_code": "321000"
		},
		{
			"id": "94",
			"province_id": "11",
			"name": "衢州市",
			"zip_code": "324000"
		},
		{
			"id": "95",
			"province_id": "11",
			"name": "舟山市",
			"zip_code": "316000"
		},
		{
			"id": "96",
			"province_id": "11",
			"name": "台州市",
			"zip_code": "318000"
		},
		{
			"id": "97",
			"province_id": "11",
			"name": "丽水市",
			"zip_code": "323000"
		},
		{
			"id": "98",
			"province_id": "12",
			"name": "合肥市",
			"zip_code": "230000"
		},
		{
			"id": "99",
			"province_id": "12",
			"name": "芜湖市",
			"zip_code": "241000"
		},
		{
			"id": "100",
			"province_id": "12",
			"name": "蚌埠市",
			"zip_code": "233000"
		},
		{
			"id": "101",
			"province_id": "12",
			"name": "淮南市",
			"zip_code": "232000"
		},
		{
			"id": "102",
			"province_id": "12",
			"name": "马鞍山市",
			"zip_code": "243000"
		},
		{
			"id": "103",
			"province_id": "12",
			"name": "淮北市",
			"zip_code": "235000"
		},
		{
			"id": "104",
			"province_id": "12",
			"name": "铜陵市",
			"zip_code": "244000"
		},
		{
			"id": "105",
			"province_id": "12",
			"name": "安庆市",
			"zip_code": "246000"
		},
		{
			"id": "106",
			"province_id": "12",
			"name": "黄山市",
			"zip_code": "242700"
		},
		{
			"id": "107",
			"province_id": "12",
			"name": "滁州市",
			"zip_code": "239000"
		},
		{
			"id": "108",
			"province_id": "12",
			"name": "阜阳市",
			"zip_code": "236100"
		},
		{
			"id": "109",
			"province_id": "12",
			"name": "宿州市",
			"zip_code": "234100"
		},
		{
			"id": "110",
			"province_id": "12",
			"name": "巢湖市",
			"zip_code": "238000"
		},
		{
			"id": "111",
			"province_id": "12",
			"name": "六安市",
			"zip_code": "237000"
		},
		{
			"id": "112",
			"province_id": "12",
			"name": "亳州市",
			"zip_code": "236800"
		},
		{
			"id": "113",
			"province_id": "12",
			"name": "池州市",
			"zip_code": "247100"
		},
		{
			"id": "114",
			"province_id": "12",
			"name": "宣城市",
			"zip_code": "366000"
		},
		{
			"id": "115",
			"province_id": "13",
			"name": "福州市",
			"zip_code": "350000"
		},
		{
			"id": "116",
			"province_id": "13",
			"name": "厦门市",
			"zip_code": "361000"
		},
		{
			"id": "117",
			"province_id": "13",
			"name": "莆田市",
			"zip_code": "351100"
		},
		{
			"id": "118",
			"province_id": "13",
			"name": "三明市",
			"zip_code": "365000"
		},
		{
			"id": "119",
			"province_id": "13",
			"name": "泉州市",
			"zip_code": "362000"
		},
		{
			"id": "120",
			"province_id": "13",
			"name": "漳州市",
			"zip_code": "363000"
		},
		{
			"id": "121",
			"province_id": "13",
			"name": "南平市",
			"zip_code": "353000"
		},
		{
			"id": "122",
			"province_id": "13",
			"name": "龙岩市",
			"zip_code": "364000"
		},
		{
			"id": "123",
			"province_id": "13",
			"name": "宁德市",
			"zip_code": "352100"
		},
		{
			"id": "124",
			"province_id": "14",
			"name": "南昌市",
			"zip_code": "330000"
		},
		{
			"id": "125",
			"province_id": "14",
			"name": "景德镇市",
			"zip_code": "333000"
		},
		{
			"id": "126",
			"province_id": "14",
			"name": "萍乡市",
			"zip_code": "337000"
		},
		{
			"id": "127",
			"province_id": "14",
			"name": "九江市",
			"zip_code": "332000"
		},
		{
			"id": "128",
			"province_id": "14",
			"name": "新余市",
			"zip_code": "338000"
		},
		{
			"id": "129",
			"province_id": "14",
			"name": "鹰潭市",
			"zip_code": "335000"
		},
		{
			"id": "130",
			"province_id": "14",
			"name": "赣州市",
			"zip_code": "341000"
		},
		{
			"id": "131",
			"province_id": "14",
			"name": "吉安市",
			"zip_code": "343000"
		},
		{
			"id": "132",
			"province_id": "14",
			"name": "宜春市",
			"zip_code": "336000"
		},
		{
			"id": "133",
			"province_id": "14",
			"name": "抚州市",
			"zip_code": "332900"
		},
		{
			"id": "134",
			"province_id": "14",
			"name": "上饶市",
			"zip_code": "334000"
		},
		{
			"id": "135",
			"province_id": "15",
			"name": "济南市",
			"zip_code": "250000"
		},
		{
			"id": "136",
			"province_id": "15",
			"name": "青岛市",
			"zip_code": "266000"
		},
		{
			"id": "137",
			"province_id": "15",
			"name": "淄博市",
			"zip_code": "255000"
		},
		{
			"id": "138",
			"province_id": "15",
			"name": "枣庄市",
			"zip_code": "277100"
		},
		{
			"id": "139",
			"province_id": "15",
			"name": "东营市",
			"zip_code": "257000"
		},
		{
			"id": "140",
			"province_id": "15",
			"name": "烟台市",
			"zip_code": "264000"
		},
		{
			"id": "141",
			"province_id": "15",
			"name": "潍坊市",
			"zip_code": "261000"
		},
		{
			"id": "142",
			"province_id": "15",
			"name": "济宁市",
			"zip_code": "272100"
		},
		{
			"id": "143",
			"province_id": "15",
			"name": "泰安市",
			"zip_code": "271000"
		},
		{
			"id": "144",
			"province_id": "15",
			"name": "威海市",
			"zip_code": "265700"
		},
		{
			"id": "145",
			"province_id": "15",
			"name": "日照市",
			"zip_code": "276800"
		},
		{
			"id": "146",
			"province_id": "15",
			"name": "莱芜市",
			"zip_code": "271100"
		},
		{
			"id": "147",
			"province_id": "15",
			"name": "临沂市",
			"zip_code": "276000"
		},
		{
			"id": "148",
			"province_id": "15",
			"name": "德州市",
			"zip_code": "253000"
		},
		{
			"id": "149",
			"province_id": "15",
			"name": "聊城市",
			"zip_code": "252000"
		},
		{
			"id": "150",
			"province_id": "15",
			"name": "滨州市",
			"zip_code": "256600"
		},
		{
			"id": "151",
			"province_id": "15",
			"name": "荷泽市",
			"zip_code": "255000"
		},
		{
			"id": "152",
			"province_id": "16",
			"name": "郑州市",
			"zip_code": "450000"
		},
		{
			"id": "153",
			"province_id": "16",
			"name": "开封市",
			"zip_code": "475000"
		},
		{
			"id": "154",
			"province_id": "16",
			"name": "洛阳市",
			"zip_code": "471000"
		},
		{
			"id": "155",
			"province_id": "16",
			"name": "平顶山市",
			"zip_code": "467000"
		},
		{
			"id": "156",
			"province_id": "16",
			"name": "安阳市",
			"zip_code": "454900"
		},
		{
			"id": "157",
			"province_id": "16",
			"name": "鹤壁市",
			"zip_code": "456600"
		},
		{
			"id": "158",
			"province_id": "16",
			"name": "新乡市",
			"zip_code": "453000"
		},
		{
			"id": "159",
			"province_id": "16",
			"name": "焦作市",
			"zip_code": "454100"
		},
		{
			"id": "160",
			"province_id": "16",
			"name": "濮阳市",
			"zip_code": "457000"
		},
		{
			"id": "161",
			"province_id": "16",
			"name": "许昌市",
			"zip_code": "461000"
		},
		{
			"id": "162",
			"province_id": "16",
			"name": "漯河市",
			"zip_code": "462000"
		},
		{
			"id": "163",
			"province_id": "16",
			"name": "三门峡市",
			"zip_code": "472000"
		},
		{
			"id": "164",
			"province_id": "16",
			"name": "南阳市",
			"zip_code": "473000"
		},
		{
			"id": "165",
			"province_id": "16",
			"name": "商丘市",
			"zip_code": "476000"
		},
		{
			"id": "166",
			"province_id": "16",
			"name": "信阳市",
			"zip_code": "464000"
		},
		{
			"id": "167",
			"province_id": "16",
			"name": "周口市",
			"zip_code": "466000"
		},
		{
			"id": "168",
			"province_id": "16",
			"name": "驻马店市",
			"zip_code": "463000"
		},
		{
			"id": "169",
			"province_id": "17",
			"name": "武汉市",
			"zip_code": "430000"
		},
		{
			"id": "170",
			"province_id": "17",
			"name": "黄石市",
			"zip_code": "435000"
		},
		{
			"id": "171",
			"province_id": "17",
			"name": "十堰市",
			"zip_code": "442000"
		},
		{
			"id": "172",
			"province_id": "17",
			"name": "宜昌市",
			"zip_code": "443000"
		},
		{
			"id": "173",
			"province_id": "17",
			"name": "襄樊市",
			"zip_code": "441000"
		},
		{
			"id": "174",
			"province_id": "17",
			"name": "鄂州市",
			"zip_code": "436000"
		},
		{
			"id": "175",
			"province_id": "17",
			"name": "荆门市",
			"zip_code": "448000"
		},
		{
			"id": "176",
			"province_id": "17",
			"name": "孝感市",
			"zip_code": "432100"
		},
		{
			"id": "177",
			"province_id": "17",
			"name": "荆州市",
			"zip_code": "434000"
		},
		{
			"id": "178",
			"province_id": "17",
			"name": "黄冈市",
			"zip_code": "438000"
		},
		{
			"id": "179",
			"province_id": "17",
			"name": "咸宁市",
			"zip_code": "437000"
		},
		{
			"id": "180",
			"province_id": "17",
			"name": "随州市",
			"zip_code": "441300"
		},
		{
			"id": "181",
			"province_id": "17",
			"name": "恩施土家族苗族自治州",
			"zip_code": "445000"
		},
		{
			"id": "182",
			"province_id": "17",
			"name": "神农架",
			"zip_code": "442400"
		},
		{
			"id": "183",
			"province_id": "18",
			"name": "长沙市",
			"zip_code": "410000"
		},
		{
			"id": "184",
			"province_id": "18",
			"name": "株洲市",
			"zip_code": "412000"
		},
		{
			"id": "185",
			"province_id": "18",
			"name": "湘潭市",
			"zip_code": "411100"
		},
		{
			"id": "186",
			"province_id": "18",
			"name": "衡阳市",
			"zip_code": "421000"
		},
		{
			"id": "187",
			"province_id": "18",
			"name": "邵阳市",
			"zip_code": "422000"
		},
		{
			"id": "188",
			"province_id": "18",
			"name": "岳阳市",
			"zip_code": "414000"
		},
		{
			"id": "189",
			"province_id": "18",
			"name": "常德市",
			"zip_code": "415000"
		},
		{
			"id": "190",
			"province_id": "18",
			"name": "张家界市",
			"zip_code": "427000"
		},
		{
			"id": "191",
			"province_id": "18",
			"name": "益阳市",
			"zip_code": "413000"
		},
		{
			"id": "192",
			"province_id": "18",
			"name": "郴州市",
			"zip_code": "423000"
		},
		{
			"id": "193",
			"province_id": "18",
			"name": "永州市",
			"zip_code": "425000"
		},
		{
			"id": "194",
			"province_id": "18",
			"name": "怀化市",
			"zip_code": "418000"
		},
		{
			"id": "195",
			"province_id": "18",
			"name": "娄底市",
			"zip_code": "417000"
		},
		{
			"id": "196",
			"province_id": "18",
			"name": "湘西土家族苗族自治州",
			"zip_code": "416000"
		},
		{
			"id": "197",
			"province_id": "19",
			"name": "广州市",
			"zip_code": "510000"
		},
		{
			"id": "198",
			"province_id": "19",
			"name": "韶关市",
			"zip_code": "521000"
		},
		{
			"id": "199",
			"province_id": "19",
			"name": "深圳市",
			"zip_code": "518000"
		},
		{
			"id": "200",
			"province_id": "19",
			"name": "珠海市",
			"zip_code": "519000"
		},
		{
			"id": "201",
			"province_id": "19",
			"name": "汕头市",
			"zip_code": "515000"
		},
		{
			"id": "202",
			"province_id": "19",
			"name": "佛山市",
			"zip_code": "528000"
		},
		{
			"id": "203",
			"province_id": "19",
			"name": "江门市",
			"zip_code": "529000"
		},
		{
			"id": "204",
			"province_id": "19",
			"name": "湛江市",
			"zip_code": "524000"
		},
		{
			"id": "205",
			"province_id": "19",
			"name": "茂名市",
			"zip_code": "525000"
		},
		{
			"id": "206",
			"province_id": "19",
			"name": "肇庆市",
			"zip_code": "526000"
		},
		{
			"id": "207",
			"province_id": "19",
			"name": "惠州市",
			"zip_code": "516000"
		},
		{
			"id": "208",
			"province_id": "19",
			"name": "梅州市",
			"zip_code": "514000"
		},
		{
			"id": "209",
			"province_id": "19",
			"name": "汕尾市",
			"zip_code": "516600"
		},
		{
			"id": "210",
			"province_id": "19",
			"name": "河源市",
			"zip_code": "517000"
		},
		{
			"id": "211",
			"province_id": "19",
			"name": "阳江市",
			"zip_code": "529500"
		},
		{
			"id": "212",
			"province_id": "19",
			"name": "清远市",
			"zip_code": "511500"
		},
		{
			"id": "213",
			"province_id": "19",
			"name": "东莞市",
			"zip_code": "511700"
		},
		{
			"id": "214",
			"province_id": "19",
			"name": "中山市",
			"zip_code": "528400"
		},
		{
			"id": "215",
			"province_id": "19",
			"name": "潮州市",
			"zip_code": "515600"
		},
		{
			"id": "216",
			"province_id": "19",
			"name": "揭阳市",
			"zip_code": "522000"
		},
		{
			"id": "217",
			"province_id": "19",
			"name": "云浮市",
			"zip_code": "527300"
		},
		{
			"id": "218",
			"province_id": "20",
			"name": "南宁市",
			"zip_code": "530000"
		},
		{
			"id": "219",
			"province_id": "20",
			"name": "柳州市",
			"zip_code": "545000"
		},
		{
			"id": "220",
			"province_id": "20",
			"name": "桂林市",
			"zip_code": "541000"
		},
		{
			"id": "221",
			"province_id": "20",
			"name": "梧州市",
			"zip_code": "543000"
		},
		{
			"id": "222",
			"province_id": "20",
			"name": "北海市",
			"zip_code": "536000"
		},
		{
			"id": "223",
			"province_id": "20",
			"name": "防城港市",
			"zip_code": "538000"
		},
		{
			"id": "224",
			"province_id": "20",
			"name": "钦州市",
			"zip_code": "535000"
		},
		{
			"id": "225",
			"province_id": "20",
			"name": "贵港市",
			"zip_code": "537100"
		},
		{
			"id": "226",
			"province_id": "20",
			"name": "玉林市",
			"zip_code": "537000"
		},
		{
			"id": "227",
			"province_id": "20",
			"name": "百色市",
			"zip_code": "533000"
		},
		{
			"id": "228",
			"province_id": "20",
			"name": "贺州市",
			"zip_code": "542800"
		},
		{
			"id": "229",
			"province_id": "20",
			"name": "河池市",
			"zip_code": "547000"
		},
		{
			"id": "230",
			"province_id": "20",
			"name": "来宾市",
			"zip_code": "546100"
		},
		{
			"id": "231",
			"province_id": "20",
			"name": "崇左市",
			"zip_code": "532200"
		},
		{
			"id": "232",
			"province_id": "21",
			"name": "海口市",
			"zip_code": "570000"
		},
		{
			"id": "233",
			"province_id": "21",
			"name": "三亚市",
			"zip_code": "572000"
		},
		{
			"id": "234",
			"province_id": "22",
			"name": "重庆市",
			"zip_code": "400000"
		},
		{
			"id": "235",
			"province_id": "23",
			"name": "成都市",
			"zip_code": "610000"
		},
		{
			"id": "236",
			"province_id": "23",
			"name": "自贡市",
			"zip_code": "643000"
		},
		{
			"id": "237",
			"province_id": "23",
			"name": "攀枝花市",
			"zip_code": "617000"
		},
		{
			"id": "238",
			"province_id": "23",
			"name": "泸州市",
			"zip_code": "646100"
		},
		{
			"id": "239",
			"province_id": "23",
			"name": "德阳市",
			"zip_code": "618000"
		},
		{
			"id": "240",
			"province_id": "23",
			"name": "绵阳市",
			"zip_code": "621000"
		},
		{
			"id": "241",
			"province_id": "23",
			"name": "广元市",
			"zip_code": "628000"
		},
		{
			"id": "242",
			"province_id": "23",
			"name": "遂宁市",
			"zip_code": "629000"
		},
		{
			"id": "243",
			"province_id": "23",
			"name": "内江市",
			"zip_code": "641000"
		},
		{
			"id": "244",
			"province_id": "23",
			"name": "乐山市",
			"zip_code": "614000"
		},
		{
			"id": "245",
			"province_id": "23",
			"name": "南充市",
			"zip_code": "637000"
		},
		{
			"id": "246",
			"province_id": "23",
			"name": "眉山市",
			"zip_code": "612100"
		},
		{
			"id": "247",
			"province_id": "23",
			"name": "宜宾市",
			"zip_code": "644000"
		},
		{
			"id": "248",
			"province_id": "23",
			"name": "广安市",
			"zip_code": "638000"
		},
		{
			"id": "249",
			"province_id": "23",
			"name": "达州市",
			"zip_code": "635000"
		},
		{
			"id": "250",
			"province_id": "23",
			"name": "雅安市",
			"zip_code": "625000"
		},
		{
			"id": "251",
			"province_id": "23",
			"name": "巴中市",
			"zip_code": "635500"
		},
		{
			"id": "252",
			"province_id": "23",
			"name": "资阳市",
			"zip_code": "641300"
		},
		{
			"id": "253",
			"province_id": "23",
			"name": "阿坝藏族羌族自治州",
			"zip_code": "624600"
		},
		{
			"id": "254",
			"province_id": "23",
			"name": "甘孜藏族自治州",
			"zip_code": "626000"
		},
		{
			"id": "255",
			"province_id": "23",
			"name": "凉山彝族自治州",
			"zip_code": "615000"
		},
		{
			"id": "256",
			"province_id": "24",
			"name": "贵阳市",
			"zip_code": "55000"
		},
		{
			"id": "257",
			"province_id": "24",
			"name": "六盘水市",
			"zip_code": "553000"
		},
		{
			"id": "258",
			"province_id": "24",
			"name": "遵义市",
			"zip_code": "563000"
		},
		{
			"id": "259",
			"province_id": "24",
			"name": "安顺市",
			"zip_code": "561000"
		},
		{
			"id": "260",
			"province_id": "24",
			"name": "铜仁地区",
			"zip_code": "554300"
		},
		{
			"id": "261",
			"province_id": "24",
			"name": "黔西南布依族苗族自治州",
			"zip_code": "551500"
		},
		{
			"id": "262",
			"province_id": "24",
			"name": "毕节地区",
			"zip_code": "551700"
		},
		{
			"id": "263",
			"province_id": "24",
			"name": "黔东南苗族侗族自治州",
			"zip_code": "551500"
		},
		{
			"id": "264",
			"province_id": "24",
			"name": "黔南布依族苗族自治州",
			"zip_code": "550100"
		},
		{
			"id": "265",
			"province_id": "25",
			"name": "昆明市",
			"zip_code": "650000"
		},
		{
			"id": "266",
			"province_id": "25",
			"name": "曲靖市",
			"zip_code": "655000"
		},
		{
			"id": "267",
			"province_id": "25",
			"name": "玉溪市",
			"zip_code": "653100"
		},
		{
			"id": "268",
			"province_id": "25",
			"name": "保山市",
			"zip_code": "678000"
		},
		{
			"id": "269",
			"province_id": "25",
			"name": "昭通市",
			"zip_code": "657000"
		},
		{
			"id": "270",
			"province_id": "25",
			"name": "丽江市",
			"zip_code": "674100"
		},
		{
			"id": "271",
			"province_id": "25",
			"name": "思茅市",
			"zip_code": "665000"
		},
		{
			"id": "272",
			"province_id": "25",
			"name": "临沧市",
			"zip_code": "677000"
		},
		{
			"id": "273",
			"province_id": "25",
			"name": "楚雄彝族自治州",
			"zip_code": "675000"
		},
		{
			"id": "274",
			"province_id": "25",
			"name": "红河哈尼族彝族自治州",
			"zip_code": "654400"
		},
		{
			"id": "275",
			"province_id": "25",
			"name": "文山壮族苗族自治州",
			"zip_code": "663000"
		},
		{
			"id": "276",
			"province_id": "25",
			"name": "西双版纳傣族自治州",
			"zip_code": "666200"
		},
		{
			"id": "277",
			"province_id": "25",
			"name": "大理白族自治州",
			"zip_code": "671000"
		},
		{
			"id": "278",
			"province_id": "25",
			"name": "德宏傣族景颇族自治州",
			"zip_code": "678400"
		},
		{
			"id": "279",
			"province_id": "25",
			"name": "怒江傈僳族自治州",
			"zip_code": "671400"
		},
		{
			"id": "280",
			"province_id": "25",
			"name": "迪庆藏族自治州",
			"zip_code": "674400"
		},
		{
			"id": "281",
			"province_id": "26",
			"name": "拉萨市",
			"zip_code": "850000"
		},
		{
			"id": "282",
			"province_id": "26",
			"name": "昌都地区",
			"zip_code": "854000"
		},
		{
			"id": "283",
			"province_id": "26",
			"name": "山南地区",
			"zip_code": "856000"
		},
		{
			"id": "284",
			"province_id": "26",
			"name": "日喀则地区",
			"zip_code": "857000"
		},
		{
			"id": "285",
			"province_id": "26",
			"name": "那曲地区",
			"zip_code": "852000"
		},
		{
			"id": "286",
			"province_id": "26",
			"name": "阿里地区",
			"zip_code": "859100"
		},
		{
			"id": "287",
			"province_id": "26",
			"name": "林芝地区",
			"zip_code": "860100"
		},
		{
			"id": "288",
			"province_id": "27",
			"name": "西安市",
			"zip_code": "710000"
		},
		{
			"id": "289",
			"province_id": "27",
			"name": "铜川市",
			"zip_code": "727000"
		},
		{
			"id": "290",
			"province_id": "27",
			"name": "宝鸡市",
			"zip_code": "721000"
		},
		{
			"id": "291",
			"province_id": "27",
			"name": "咸阳市",
			"zip_code": "712000"
		},
		{
			"id": "292",
			"province_id": "27",
			"name": "渭南市",
			"zip_code": "714000"
		},
		{
			"id": "293",
			"province_id": "27",
			"name": "延安市",
			"zip_code": "716000"
		},
		{
			"id": "294",
			"province_id": "27",
			"name": "汉中市",
			"zip_code": "723000"
		},
		{
			"id": "295",
			"province_id": "27",
			"name": "榆林市",
			"zip_code": "719000"
		},
		{
			"id": "296",
			"province_id": "27",
			"name": "安康市",
			"zip_code": "725000"
		},
		{
			"id": "297",
			"province_id": "27",
			"name": "商洛市",
			"zip_code": "711500"
		},
		{
			"id": "298",
			"province_id": "28",
			"name": "兰州市",
			"zip_code": "730000"
		},
		{
			"id": "299",
			"province_id": "28",
			"name": "嘉峪关市",
			"zip_code": "735100"
		},
		{
			"id": "300",
			"province_id": "28",
			"name": "金昌市",
			"zip_code": "737100"
		},
		{
			"id": "301",
			"province_id": "28",
			"name": "白银市",
			"zip_code": "730900"
		},
		{
			"id": "302",
			"province_id": "28",
			"name": "天水市",
			"zip_code": "741000"
		},
		{
			"id": "303",
			"province_id": "28",
			"name": "武威市",
			"zip_code": "733000"
		},
		{
			"id": "304",
			"province_id": "28",
			"name": "张掖市",
			"zip_code": "734000"
		},
		{
			"id": "305",
			"province_id": "28",
			"name": "平凉市",
			"zip_code": "744000"
		},
		{
			"id": "306",
			"province_id": "28",
			"name": "酒泉市",
			"zip_code": "735000"
		},
		{
			"id": "307",
			"province_id": "28",
			"name": "庆阳市",
			"zip_code": "744500"
		},
		{
			"id": "308",
			"province_id": "28",
			"name": "定西市",
			"zip_code": "743000"
		},
		{
			"id": "309",
			"province_id": "28",
			"name": "陇南市",
			"zip_code": "742100"
		},
		{
			"id": "310",
			"province_id": "28",
			"name": "临夏回族自治州",
			"zip_code": "731100"
		},
		{
			"id": "311",
			"province_id": "28",
			"name": "甘南藏族自治州",
			"zip_code": "747000"
		},
		{
			"id": "312",
			"province_id": "29",
			"name": "西宁市",
			"zip_code": "810000"
		},
		{
			"id": "313",
			"province_id": "29",
			"name": "海东地区",
			"zip_code": "810600"
		},
		{
			"id": "314",
			"province_id": "29",
			"name": "海北藏族自治州",
			"zip_code": "810300"
		},
		{
			"id": "315",
			"province_id": "29",
			"name": "黄南藏族自治州",
			"zip_code": "811300"
		},
		{
			"id": "316",
			"province_id": "29",
			"name": "海南藏族自治州",
			"zip_code": "813000"
		},
		{
			"id": "317",
			"province_id": "29",
			"name": "果洛藏族自治州",
			"zip_code": "814000"
		},
		{
			"id": "318",
			"province_id": "29",
			"name": "玉树藏族自治州",
			"zip_code": "815000"
		},
		{
			"id": "319",
			"province_id": "29",
			"name": "海西蒙古族藏族自治州",
			"zip_code": "817000"
		},
		{
			"id": "320",
			"province_id": "30",
			"name": "银川市",
			"zip_code": "750000"
		},
		{
			"id": "321",
			"province_id": "30",
			"name": "石嘴山市",
			"zip_code": "753000"
		},
		{
			"id": "322",
			"province_id": "30",
			"name": "吴忠市",
			"zip_code": "751100"
		},
		{
			"id": "323",
			"province_id": "30",
			"name": "固原市",
			"zip_code": "756000"
		},
		{
			"id": "324",
			"province_id": "30",
			"name": "中卫市",
			"zip_code": "751700"
		},
		{
			"id": "325",
			"province_id": "31",
			"name": "乌鲁木齐市",
			"zip_code": "830000"
		},
		{
			"id": "326",
			"province_id": "31",
			"name": "克拉玛依市",
			"zip_code": "834000"
		},
		{
			"id": "327",
			"province_id": "31",
			"name": "吐鲁番地区",
			"zip_code": "838000"
		},
		{
			"id": "328",
			"province_id": "31",
			"name": "哈密地区",
			"zip_code": "839000"
		},
		{
			"id": "329",
			"province_id": "31",
			"name": "昌吉回族自治州",
			"zip_code": "831100"
		},
		{
			"id": "330",
			"province_id": "31",
			"name": "博尔塔拉蒙古自治州",
			"zip_code": "833400"
		},
		{
			"id": "331",
			"province_id": "31",
			"name": "巴音郭楞蒙古自治州",
			"zip_code": "841000"
		},
		{
			"id": "332",
			"province_id": "31",
			"name": "阿克苏地区",
			"zip_code": "843000"
		},
		{
			"id": "333",
			"province_id": "31",
			"name": "克孜勒苏柯尔克孜自治州",
			"zip_code": "835600"
		},
		{
			"id": "334",
			"province_id": "31",
			"name": "喀什地区",
			"zip_code": "844000"
		},
		{
			"id": "335",
			"province_id": "31",
			"name": "和田地区",
			"zip_code": "848000"
		},
		{
			"id": "336",
			"province_id": "31",
			"name": "伊犁哈萨克自治州",
			"zip_code": "833200"
		},
		{
			"id": "337",
			"province_id": "31",
			"name": "塔城地区",
			"zip_code": "834700"
		},
		{
			"id": "338",
			"province_id": "31",
			"name": "阿勒泰地区",
			"zip_code": "836500"
		},
		{
			"id": "339",
			"province_id": "31",
			"name": "石河子市",
			"zip_code": "832000"
		},
		{
			"id": "340",
			"province_id": "31",
			"name": "阿拉尔市",
			"zip_code": "843300"
		},
		{
			"id": "341",
			"province_id": "31",
			"name": "图木舒克市",
			"zip_code": "843900"
		},
		{
			"id": "342",
			"province_id": "31",
			"name": "五家渠市",
			"zip_code": "831300"
		},
		{
			"id": "343",
			"province_id": "32",
			"name": "香港特别行政区",
			"zip_code": "000000"
		},
		{
			"id": "344",
			"province_id": "33",
			"name": "澳门特别行政区",
			"zip_code": "000000"
		},
		{
			"id": "345",
			"province_id": "34",
			"name": "台湾省",
			"zip_code": "000000"
		}
	],
	"DISTRICTS": [
		{
			"id": "1",
			"city_id": "1",
			"name": "东城区"
		},
		{
			"id": "2",
			"city_id": "1",
			"name": "西城区"
		},
		{
			"id": "3",
			"city_id": "1",
			"name": "崇文区"
		},
		{
			"id": "4",
			"city_id": "1",
			"name": "宣武区"
		},
		{
			"id": "5",
			"city_id": "1",
			"name": "朝阳区"
		},
		{
			"id": "6",
			"city_id": "1",
			"name": "丰台区"
		},
		{
			"id": "7",
			"city_id": "1",
			"name": "石景山区"
		},
		{
			"id": "8",
			"city_id": "1",
			"name": "海淀区"
		},
		{
			"id": "9",
			"city_id": "1",
			"name": "门头沟区"
		},
		{
			"id": "10",
			"city_id": "1",
			"name": "房山区"
		},
		{
			"id": "11",
			"city_id": "1",
			"name": "通州区"
		},
		{
			"id": "12",
			"city_id": "1",
			"name": "顺义区"
		},
		{
			"id": "13",
			"city_id": "1",
			"name": "昌平区"
		},
		{
			"id": "14",
			"city_id": "1",
			"name": "大兴区"
		},
		{
			"id": "15",
			"city_id": "1",
			"name": "怀柔区"
		},
		{
			"id": "16",
			"city_id": "1",
			"name": "平谷区"
		},
		{
			"id": "17",
			"city_id": "1",
			"name": "密云县"
		},
		{
			"id": "18",
			"city_id": "1",
			"name": "延庆县"
		},
		{
			"id": "19",
			"city_id": "2",
			"name": "和平区"
		},
		{
			"id": "20",
			"city_id": "2",
			"name": "河东区"
		},
		{
			"id": "21",
			"city_id": "2",
			"name": "河西区"
		},
		{
			"id": "22",
			"city_id": "2",
			"name": "南开区"
		},
		{
			"id": "23",
			"city_id": "2",
			"name": "河北区"
		},
		{
			"id": "24",
			"city_id": "2",
			"name": "红桥区"
		},
		{
			"id": "25",
			"city_id": "2",
			"name": "塘沽区"
		},
		{
			"id": "26",
			"city_id": "2",
			"name": "汉沽区"
		},
		{
			"id": "27",
			"city_id": "2",
			"name": "大港区"
		},
		{
			"id": "28",
			"city_id": "2",
			"name": "东丽区"
		},
		{
			"id": "29",
			"city_id": "2",
			"name": "西青区"
		},
		{
			"id": "30",
			"city_id": "2",
			"name": "津南区"
		},
		{
			"id": "31",
			"city_id": "2",
			"name": "北辰区"
		},
		{
			"id": "32",
			"city_id": "2",
			"name": "武清区"
		},
		{
			"id": "33",
			"city_id": "2",
			"name": "宝坻区"
		},
		{
			"id": "34",
			"city_id": "2",
			"name": "宁河县"
		},
		{
			"id": "35",
			"city_id": "2",
			"name": "静海县"
		},
		{
			"id": "36",
			"city_id": "2",
			"name": "蓟县"
		},
		{
			"id": "37",
			"city_id": "3",
			"name": "长安区"
		},
		{
			"id": "38",
			"city_id": "3",
			"name": "桥东区"
		},
		{
			"id": "39",
			"city_id": "3",
			"name": "桥西区"
		},
		{
			"id": "40",
			"city_id": "3",
			"name": "新华区"
		},
		{
			"id": "41",
			"city_id": "3",
			"name": "井陉矿区"
		},
		{
			"id": "42",
			"city_id": "3",
			"name": "裕华区"
		},
		{
			"id": "43",
			"city_id": "3",
			"name": "井陉县"
		},
		{
			"id": "44",
			"city_id": "3",
			"name": "正定县"
		},
		{
			"id": "45",
			"city_id": "3",
			"name": "栾城县"
		},
		{
			"id": "46",
			"city_id": "3",
			"name": "行唐县"
		},
		{
			"id": "47",
			"city_id": "3",
			"name": "灵寿县"
		},
		{
			"id": "48",
			"city_id": "3",
			"name": "高邑县"
		},
		{
			"id": "49",
			"city_id": "3",
			"name": "深泽县"
		},
		{
			"id": "50",
			"city_id": "3",
			"name": "赞皇县"
		},
		{
			"id": "51",
			"city_id": "3",
			"name": "无极县"
		},
		{
			"id": "52",
			"city_id": "3",
			"name": "平山县"
		},
		{
			"id": "53",
			"city_id": "3",
			"name": "元氏县"
		},
		{
			"id": "54",
			"city_id": "3",
			"name": "赵县"
		},
		{
			"id": "55",
			"city_id": "3",
			"name": "辛集市"
		},
		{
			"id": "56",
			"city_id": "3",
			"name": "藁城市"
		},
		{
			"id": "57",
			"city_id": "3",
			"name": "晋州市"
		},
		{
			"id": "58",
			"city_id": "3",
			"name": "新乐市"
		},
		{
			"id": "59",
			"city_id": "3",
			"name": "鹿泉市"
		},
		{
			"id": "60",
			"city_id": "4",
			"name": "路南区"
		},
		{
			"id": "61",
			"city_id": "4",
			"name": "路北区"
		},
		{
			"id": "62",
			"city_id": "4",
			"name": "古冶区"
		},
		{
			"id": "63",
			"city_id": "4",
			"name": "开平区"
		},
		{
			"id": "64",
			"city_id": "4",
			"name": "丰南区"
		},
		{
			"id": "65",
			"city_id": "4",
			"name": "丰润区"
		},
		{
			"id": "66",
			"city_id": "4",
			"name": "滦县"
		},
		{
			"id": "67",
			"city_id": "4",
			"name": "滦南县"
		},
		{
			"id": "68",
			"city_id": "4",
			"name": "乐亭县"
		},
		{
			"id": "69",
			"city_id": "4",
			"name": "迁西县"
		},
		{
			"id": "70",
			"city_id": "4",
			"name": "玉田县"
		},
		{
			"id": "71",
			"city_id": "4",
			"name": "唐海县"
		},
		{
			"id": "72",
			"city_id": "4",
			"name": "遵化市"
		},
		{
			"id": "73",
			"city_id": "4",
			"name": "迁安市"
		},
		{
			"id": "74",
			"city_id": "5",
			"name": "海港区"
		},
		{
			"id": "75",
			"city_id": "5",
			"name": "山海关区"
		},
		{
			"id": "76",
			"city_id": "5",
			"name": "北戴河区"
		},
		{
			"id": "77",
			"city_id": "5",
			"name": "青龙满族自治县"
		},
		{
			"id": "78",
			"city_id": "5",
			"name": "昌黎县"
		},
		{
			"id": "79",
			"city_id": "5",
			"name": "抚宁县"
		},
		{
			"id": "80",
			"city_id": "5",
			"name": "卢龙县"
		},
		{
			"id": "81",
			"city_id": "6",
			"name": "邯山区"
		},
		{
			"id": "82",
			"city_id": "6",
			"name": "丛台区"
		},
		{
			"id": "83",
			"city_id": "6",
			"name": "复兴区"
		},
		{
			"id": "84",
			"city_id": "6",
			"name": "峰峰矿区"
		},
		{
			"id": "85",
			"city_id": "6",
			"name": "邯郸县"
		},
		{
			"id": "86",
			"city_id": "6",
			"name": "临漳县"
		},
		{
			"id": "87",
			"city_id": "6",
			"name": "成安县"
		},
		{
			"id": "88",
			"city_id": "6",
			"name": "大名县"
		},
		{
			"id": "89",
			"city_id": "6",
			"name": "涉县"
		},
		{
			"id": "90",
			"city_id": "6",
			"name": "磁县"
		},
		{
			"id": "91",
			"city_id": "6",
			"name": "肥乡县"
		},
		{
			"id": "92",
			"city_id": "6",
			"name": "永年县"
		},
		{
			"id": "93",
			"city_id": "6",
			"name": "邱县"
		},
		{
			"id": "94",
			"city_id": "6",
			"name": "鸡泽县"
		},
		{
			"id": "95",
			"city_id": "6",
			"name": "广平县"
		},
		{
			"id": "96",
			"city_id": "6",
			"name": "馆陶县"
		},
		{
			"id": "97",
			"city_id": "6",
			"name": "魏县"
		},
		{
			"id": "98",
			"city_id": "6",
			"name": "曲周县"
		},
		{
			"id": "99",
			"city_id": "6",
			"name": "武安市"
		},
		{
			"id": "100",
			"city_id": "7",
			"name": "桥东区"
		},
		{
			"id": "101",
			"city_id": "7",
			"name": "桥西区"
		},
		{
			"id": "102",
			"city_id": "7",
			"name": "邢台县"
		},
		{
			"id": "103",
			"city_id": "7",
			"name": "临城县"
		},
		{
			"id": "104",
			"city_id": "7",
			"name": "内丘县"
		},
		{
			"id": "105",
			"city_id": "7",
			"name": "柏乡县"
		},
		{
			"id": "106",
			"city_id": "7",
			"name": "隆尧县"
		},
		{
			"id": "107",
			"city_id": "7",
			"name": "任县"
		},
		{
			"id": "108",
			"city_id": "7",
			"name": "南和县"
		},
		{
			"id": "109",
			"city_id": "7",
			"name": "宁晋县"
		},
		{
			"id": "110",
			"city_id": "7",
			"name": "巨鹿县"
		},
		{
			"id": "111",
			"city_id": "7",
			"name": "新河县"
		},
		{
			"id": "112",
			"city_id": "7",
			"name": "广宗县"
		},
		{
			"id": "113",
			"city_id": "7",
			"name": "平乡县"
		},
		{
			"id": "114",
			"city_id": "7",
			"name": "威县"
		},
		{
			"id": "115",
			"city_id": "7",
			"name": "清河县"
		},
		{
			"id": "116",
			"city_id": "7",
			"name": "临西县"
		},
		{
			"id": "117",
			"city_id": "7",
			"name": "南宫市"
		},
		{
			"id": "118",
			"city_id": "7",
			"name": "沙河市"
		},
		{
			"id": "119",
			"city_id": "8",
			"name": "新市区"
		},
		{
			"id": "120",
			"city_id": "8",
			"name": "北市区"
		},
		{
			"id": "121",
			"city_id": "8",
			"name": "南市区"
		},
		{
			"id": "122",
			"city_id": "8",
			"name": "满城县"
		},
		{
			"id": "123",
			"city_id": "8",
			"name": "清苑县"
		},
		{
			"id": "124",
			"city_id": "8",
			"name": "涞水县"
		},
		{
			"id": "125",
			"city_id": "8",
			"name": "阜平县"
		},
		{
			"id": "126",
			"city_id": "8",
			"name": "徐水县"
		},
		{
			"id": "127",
			"city_id": "8",
			"name": "定兴县"
		},
		{
			"id": "128",
			"city_id": "8",
			"name": "唐县"
		},
		{
			"id": "129",
			"city_id": "8",
			"name": "高阳县"
		},
		{
			"id": "130",
			"city_id": "8",
			"name": "容城县"
		},
		{
			"id": "131",
			"city_id": "8",
			"name": "涞源县"
		},
		{
			"id": "132",
			"city_id": "8",
			"name": "望都县"
		},
		{
			"id": "133",
			"city_id": "8",
			"name": "安新县"
		},
		{
			"id": "134",
			"city_id": "8",
			"name": "易县"
		},
		{
			"id": "135",
			"city_id": "8",
			"name": "曲阳县"
		},
		{
			"id": "136",
			"city_id": "8",
			"name": "蠡县"
		},
		{
			"id": "137",
			"city_id": "8",
			"name": "顺平县"
		},
		{
			"id": "138",
			"city_id": "8",
			"name": "博野县"
		},
		{
			"id": "139",
			"city_id": "8",
			"name": "雄县"
		},
		{
			"id": "140",
			"city_id": "8",
			"name": "涿州市"
		},
		{
			"id": "141",
			"city_id": "8",
			"name": "定州市"
		},
		{
			"id": "142",
			"city_id": "8",
			"name": "安国市"
		},
		{
			"id": "143",
			"city_id": "8",
			"name": "高碑店市"
		},
		{
			"id": "144",
			"city_id": "9",
			"name": "桥东区"
		},
		{
			"id": "145",
			"city_id": "9",
			"name": "桥西区"
		},
		{
			"id": "146",
			"city_id": "9",
			"name": "宣化区"
		},
		{
			"id": "147",
			"city_id": "9",
			"name": "下花园区"
		},
		{
			"id": "148",
			"city_id": "9",
			"name": "宣化县"
		},
		{
			"id": "149",
			"city_id": "9",
			"name": "张北县"
		},
		{
			"id": "150",
			"city_id": "9",
			"name": "康保县"
		},
		{
			"id": "151",
			"city_id": "9",
			"name": "沽源县"
		},
		{
			"id": "152",
			"city_id": "9",
			"name": "尚义县"
		},
		{
			"id": "153",
			"city_id": "9",
			"name": "蔚县"
		},
		{
			"id": "154",
			"city_id": "9",
			"name": "阳原县"
		},
		{
			"id": "155",
			"city_id": "9",
			"name": "怀安县"
		},
		{
			"id": "156",
			"city_id": "9",
			"name": "万全县"
		},
		{
			"id": "157",
			"city_id": "9",
			"name": "怀来县"
		},
		{
			"id": "158",
			"city_id": "9",
			"name": "涿鹿县"
		},
		{
			"id": "159",
			"city_id": "9",
			"name": "赤城县"
		},
		{
			"id": "160",
			"city_id": "9",
			"name": "崇礼县"
		},
		{
			"id": "161",
			"city_id": "10",
			"name": "双桥区"
		},
		{
			"id": "162",
			"city_id": "10",
			"name": "双滦区"
		},
		{
			"id": "163",
			"city_id": "10",
			"name": "鹰手营子矿区"
		},
		{
			"id": "164",
			"city_id": "10",
			"name": "承德县"
		},
		{
			"id": "165",
			"city_id": "10",
			"name": "兴隆县"
		},
		{
			"id": "166",
			"city_id": "10",
			"name": "平泉县"
		},
		{
			"id": "167",
			"city_id": "10",
			"name": "滦平县"
		},
		{
			"id": "168",
			"city_id": "10",
			"name": "隆化县"
		},
		{
			"id": "169",
			"city_id": "10",
			"name": "丰宁满族自治县"
		},
		{
			"id": "170",
			"city_id": "10",
			"name": "宽城满族自治县"
		},
		{
			"id": "171",
			"city_id": "10",
			"name": "围场满族蒙古族自治县"
		},
		{
			"id": "172",
			"city_id": "11",
			"name": "新华区"
		},
		{
			"id": "173",
			"city_id": "11",
			"name": "运河区"
		},
		{
			"id": "174",
			"city_id": "11",
			"name": "沧县"
		},
		{
			"id": "175",
			"city_id": "11",
			"name": "青县"
		},
		{
			"id": "176",
			"city_id": "11",
			"name": "东光县"
		},
		{
			"id": "177",
			"city_id": "11",
			"name": "海兴县"
		},
		{
			"id": "178",
			"city_id": "11",
			"name": "盐山县"
		},
		{
			"id": "179",
			"city_id": "11",
			"name": "肃宁县"
		},
		{
			"id": "180",
			"city_id": "11",
			"name": "南皮县"
		},
		{
			"id": "181",
			"city_id": "11",
			"name": "吴桥县"
		},
		{
			"id": "182",
			"city_id": "11",
			"name": "献县"
		},
		{
			"id": "183",
			"city_id": "11",
			"name": "孟村回族自治县"
		},
		{
			"id": "184",
			"city_id": "11",
			"name": "泊头市"
		},
		{
			"id": "185",
			"city_id": "11",
			"name": "任丘市"
		},
		{
			"id": "186",
			"city_id": "11",
			"name": "黄骅市"
		},
		{
			"id": "187",
			"city_id": "11",
			"name": "河间市"
		},
		{
			"id": "188",
			"city_id": "12",
			"name": "安次区"
		},
		{
			"id": "189",
			"city_id": "12",
			"name": "广阳区"
		},
		{
			"id": "190",
			"city_id": "12",
			"name": "固安县"
		},
		{
			"id": "191",
			"city_id": "12",
			"name": "永清县"
		},
		{
			"id": "192",
			"city_id": "12",
			"name": "香河县"
		},
		{
			"id": "193",
			"city_id": "12",
			"name": "大城县"
		},
		{
			"id": "194",
			"city_id": "12",
			"name": "文安县"
		},
		{
			"id": "195",
			"city_id": "12",
			"name": "大厂回族自治县"
		},
		{
			"id": "196",
			"city_id": "12",
			"name": "霸州市"
		},
		{
			"id": "197",
			"city_id": "12",
			"name": "三河市"
		},
		{
			"id": "198",
			"city_id": "13",
			"name": "桃城区"
		},
		{
			"id": "199",
			"city_id": "13",
			"name": "枣强县"
		},
		{
			"id": "200",
			"city_id": "13",
			"name": "武邑县"
		},
		{
			"id": "201",
			"city_id": "13",
			"name": "武强县"
		},
		{
			"id": "202",
			"city_id": "13",
			"name": "饶阳县"
		},
		{
			"id": "203",
			"city_id": "13",
			"name": "安平县"
		},
		{
			"id": "204",
			"city_id": "13",
			"name": "故城县"
		},
		{
			"id": "205",
			"city_id": "13",
			"name": "景县"
		},
		{
			"id": "206",
			"city_id": "13",
			"name": "阜城县"
		},
		{
			"id": "207",
			"city_id": "13",
			"name": "冀州市"
		},
		{
			"id": "208",
			"city_id": "13",
			"name": "深州市"
		},
		{
			"id": "209",
			"city_id": "14",
			"name": "小店区"
		},
		{
			"id": "210",
			"city_id": "14",
			"name": "迎泽区"
		},
		{
			"id": "211",
			"city_id": "14",
			"name": "杏花岭区"
		},
		{
			"id": "212",
			"city_id": "14",
			"name": "尖草坪区"
		},
		{
			"id": "213",
			"city_id": "14",
			"name": "万柏林区"
		},
		{
			"id": "214",
			"city_id": "14",
			"name": "晋源区"
		},
		{
			"id": "215",
			"city_id": "14",
			"name": "清徐县"
		},
		{
			"id": "216",
			"city_id": "14",
			"name": "阳曲县"
		},
		{
			"id": "217",
			"city_id": "14",
			"name": "娄烦县"
		},
		{
			"id": "218",
			"city_id": "14",
			"name": "古交市"
		},
		{
			"id": "219",
			"city_id": "15",
			"name": "城区"
		},
		{
			"id": "220",
			"city_id": "15",
			"name": "矿区"
		},
		{
			"id": "221",
			"city_id": "15",
			"name": "南郊区"
		},
		{
			"id": "222",
			"city_id": "15",
			"name": "新荣区"
		},
		{
			"id": "223",
			"city_id": "15",
			"name": "阳高县"
		},
		{
			"id": "224",
			"city_id": "15",
			"name": "天镇县"
		},
		{
			"id": "225",
			"city_id": "15",
			"name": "广灵县"
		},
		{
			"id": "226",
			"city_id": "15",
			"name": "灵丘县"
		},
		{
			"id": "227",
			"city_id": "15",
			"name": "浑源县"
		},
		{
			"id": "228",
			"city_id": "15",
			"name": "左云县"
		},
		{
			"id": "229",
			"city_id": "15",
			"name": "大同县"
		},
		{
			"id": "230",
			"city_id": "16",
			"name": "城区"
		},
		{
			"id": "231",
			"city_id": "16",
			"name": "矿区"
		},
		{
			"id": "232",
			"city_id": "16",
			"name": "郊区"
		},
		{
			"id": "233",
			"city_id": "16",
			"name": "平定县"
		},
		{
			"id": "234",
			"city_id": "16",
			"name": "盂县"
		},
		{
			"id": "235",
			"city_id": "17",
			"name": "城区"
		},
		{
			"id": "236",
			"city_id": "17",
			"name": "郊区"
		},
		{
			"id": "237",
			"city_id": "17",
			"name": "长治县"
		},
		{
			"id": "238",
			"city_id": "17",
			"name": "襄垣县"
		},
		{
			"id": "239",
			"city_id": "17",
			"name": "屯留县"
		},
		{
			"id": "240",
			"city_id": "17",
			"name": "平顺县"
		},
		{
			"id": "241",
			"city_id": "17",
			"name": "黎城县"
		},
		{
			"id": "242",
			"city_id": "17",
			"name": "壶关县"
		},
		{
			"id": "243",
			"city_id": "17",
			"name": "长子县"
		},
		{
			"id": "244",
			"city_id": "17",
			"name": "武乡县"
		},
		{
			"id": "245",
			"city_id": "17",
			"name": "沁县"
		},
		{
			"id": "246",
			"city_id": "17",
			"name": "沁源县"
		},
		{
			"id": "247",
			"city_id": "17",
			"name": "潞城市"
		},
		{
			"id": "248",
			"city_id": "18",
			"name": "城区"
		},
		{
			"id": "249",
			"city_id": "18",
			"name": "沁水县"
		},
		{
			"id": "250",
			"city_id": "18",
			"name": "阳城县"
		},
		{
			"id": "251",
			"city_id": "18",
			"name": "陵川县"
		},
		{
			"id": "252",
			"city_id": "18",
			"name": "泽州县"
		},
		{
			"id": "253",
			"city_id": "18",
			"name": "高平市"
		},
		{
			"id": "254",
			"city_id": "19",
			"name": "朔城区"
		},
		{
			"id": "255",
			"city_id": "19",
			"name": "平鲁区"
		},
		{
			"id": "256",
			"city_id": "19",
			"name": "山阴县"
		},
		{
			"id": "257",
			"city_id": "19",
			"name": "应县"
		},
		{
			"id": "258",
			"city_id": "19",
			"name": "右玉县"
		},
		{
			"id": "259",
			"city_id": "19",
			"name": "怀仁县"
		},
		{
			"id": "260",
			"city_id": "20",
			"name": "榆次区"
		},
		{
			"id": "261",
			"city_id": "20",
			"name": "榆社县"
		},
		{
			"id": "262",
			"city_id": "20",
			"name": "左权县"
		},
		{
			"id": "263",
			"city_id": "20",
			"name": "和顺县"
		},
		{
			"id": "264",
			"city_id": "20",
			"name": "昔阳县"
		},
		{
			"id": "265",
			"city_id": "20",
			"name": "寿阳县"
		},
		{
			"id": "266",
			"city_id": "20",
			"name": "太谷县"
		},
		{
			"id": "267",
			"city_id": "20",
			"name": "祁县"
		},
		{
			"id": "268",
			"city_id": "20",
			"name": "平遥县"
		},
		{
			"id": "269",
			"city_id": "20",
			"name": "灵石县"
		},
		{
			"id": "270",
			"city_id": "20",
			"name": "介休市"
		},
		{
			"id": "271",
			"city_id": "21",
			"name": "盐湖区"
		},
		{
			"id": "272",
			"city_id": "21",
			"name": "临猗县"
		},
		{
			"id": "273",
			"city_id": "21",
			"name": "万荣县"
		},
		{
			"id": "274",
			"city_id": "21",
			"name": "闻喜县"
		},
		{
			"id": "275",
			"city_id": "21",
			"name": "稷山县"
		},
		{
			"id": "276",
			"city_id": "21",
			"name": "新绛县"
		},
		{
			"id": "277",
			"city_id": "21",
			"name": "绛县"
		},
		{
			"id": "278",
			"city_id": "21",
			"name": "垣曲县"
		},
		{
			"id": "279",
			"city_id": "21",
			"name": "夏县"
		},
		{
			"id": "280",
			"city_id": "21",
			"name": "平陆县"
		},
		{
			"id": "281",
			"city_id": "21",
			"name": "芮城县"
		},
		{
			"id": "282",
			"city_id": "21",
			"name": "永济市"
		},
		{
			"id": "283",
			"city_id": "21",
			"name": "河津市"
		},
		{
			"id": "284",
			"city_id": "22",
			"name": "忻府区"
		},
		{
			"id": "285",
			"city_id": "22",
			"name": "定襄县"
		},
		{
			"id": "286",
			"city_id": "22",
			"name": "五台县"
		},
		{
			"id": "287",
			"city_id": "22",
			"name": "代县"
		},
		{
			"id": "288",
			"city_id": "22",
			"name": "繁峙县"
		},
		{
			"id": "289",
			"city_id": "22",
			"name": "宁武县"
		},
		{
			"id": "290",
			"city_id": "22",
			"name": "静乐县"
		},
		{
			"id": "291",
			"city_id": "22",
			"name": "神池县"
		},
		{
			"id": "292",
			"city_id": "22",
			"name": "五寨县"
		},
		{
			"id": "293",
			"city_id": "22",
			"name": "岢岚县"
		},
		{
			"id": "294",
			"city_id": "22",
			"name": "河曲县"
		},
		{
			"id": "295",
			"city_id": "22",
			"name": "保德县"
		},
		{
			"id": "296",
			"city_id": "22",
			"name": "偏关县"
		},
		{
			"id": "297",
			"city_id": "22",
			"name": "原平市"
		},
		{
			"id": "298",
			"city_id": "23",
			"name": "尧都区"
		},
		{
			"id": "299",
			"city_id": "23",
			"name": "曲沃县"
		},
		{
			"id": "300",
			"city_id": "23",
			"name": "翼城县"
		},
		{
			"id": "301",
			"city_id": "23",
			"name": "襄汾县"
		},
		{
			"id": "302",
			"city_id": "23",
			"name": "洪洞县"
		},
		{
			"id": "303",
			"city_id": "23",
			"name": "古县"
		},
		{
			"id": "304",
			"city_id": "23",
			"name": "安泽县"
		},
		{
			"id": "305",
			"city_id": "23",
			"name": "浮山县"
		},
		{
			"id": "306",
			"city_id": "23",
			"name": "吉县"
		},
		{
			"id": "307",
			"city_id": "23",
			"name": "乡宁县"
		},
		{
			"id": "308",
			"city_id": "23",
			"name": "大宁县"
		},
		{
			"id": "309",
			"city_id": "23",
			"name": "隰县"
		},
		{
			"id": "310",
			"city_id": "23",
			"name": "永和县"
		},
		{
			"id": "311",
			"city_id": "23",
			"name": "蒲县"
		},
		{
			"id": "312",
			"city_id": "23",
			"name": "汾西县"
		},
		{
			"id": "313",
			"city_id": "23",
			"name": "侯马市"
		},
		{
			"id": "314",
			"city_id": "23",
			"name": "霍州市"
		},
		{
			"id": "315",
			"city_id": "24",
			"name": "离石区"
		},
		{
			"id": "316",
			"city_id": "24",
			"name": "文水县"
		},
		{
			"id": "317",
			"city_id": "24",
			"name": "交城县"
		},
		{
			"id": "318",
			"city_id": "24",
			"name": "兴县"
		},
		{
			"id": "319",
			"city_id": "24",
			"name": "临县"
		},
		{
			"id": "320",
			"city_id": "24",
			"name": "柳林县"
		},
		{
			"id": "321",
			"city_id": "24",
			"name": "石楼县"
		},
		{
			"id": "322",
			"city_id": "24",
			"name": "岚县"
		},
		{
			"id": "323",
			"city_id": "24",
			"name": "方山县"
		},
		{
			"id": "324",
			"city_id": "24",
			"name": "中阳县"
		},
		{
			"id": "325",
			"city_id": "24",
			"name": "交口县"
		},
		{
			"id": "326",
			"city_id": "24",
			"name": "孝义市"
		},
		{
			"id": "327",
			"city_id": "24",
			"name": "汾阳市"
		},
		{
			"id": "328",
			"city_id": "25",
			"name": "新城区"
		},
		{
			"id": "329",
			"city_id": "25",
			"name": "回民区"
		},
		{
			"id": "330",
			"city_id": "25",
			"name": "玉泉区"
		},
		{
			"id": "331",
			"city_id": "25",
			"name": "赛罕区"
		},
		{
			"id": "332",
			"city_id": "25",
			"name": "土默特左旗"
		},
		{
			"id": "333",
			"city_id": "25",
			"name": "托克托县"
		},
		{
			"id": "334",
			"city_id": "25",
			"name": "和林格尔县"
		},
		{
			"id": "335",
			"city_id": "25",
			"name": "清水河县"
		},
		{
			"id": "336",
			"city_id": "25",
			"name": "武川县"
		},
		{
			"id": "337",
			"city_id": "26",
			"name": "东河区"
		},
		{
			"id": "338",
			"city_id": "26",
			"name": "昆都仑区"
		},
		{
			"id": "339",
			"city_id": "26",
			"name": "青山区"
		},
		{
			"id": "340",
			"city_id": "26",
			"name": "石拐区"
		},
		{
			"id": "341",
			"city_id": "26",
			"name": "白云矿区"
		},
		{
			"id": "342",
			"city_id": "26",
			"name": "九原区"
		},
		{
			"id": "343",
			"city_id": "26",
			"name": "土默特右旗"
		},
		{
			"id": "344",
			"city_id": "26",
			"name": "固阳县"
		},
		{
			"id": "345",
			"city_id": "26",
			"name": "达尔罕茂明安联合旗"
		},
		{
			"id": "346",
			"city_id": "27",
			"name": "海勃湾区"
		},
		{
			"id": "347",
			"city_id": "27",
			"name": "海南区"
		},
		{
			"id": "348",
			"city_id": "27",
			"name": "乌达区"
		},
		{
			"id": "349",
			"city_id": "28",
			"name": "红山区"
		},
		{
			"id": "350",
			"city_id": "28",
			"name": "元宝山区"
		},
		{
			"id": "351",
			"city_id": "28",
			"name": "松山区"
		},
		{
			"id": "352",
			"city_id": "28",
			"name": "阿鲁科尔沁旗"
		},
		{
			"id": "353",
			"city_id": "28",
			"name": "巴林左旗"
		},
		{
			"id": "354",
			"city_id": "28",
			"name": "巴林右旗"
		},
		{
			"id": "355",
			"city_id": "28",
			"name": "林西县"
		},
		{
			"id": "356",
			"city_id": "28",
			"name": "克什克腾旗"
		},
		{
			"id": "357",
			"city_id": "28",
			"name": "翁牛特旗"
		},
		{
			"id": "358",
			"city_id": "28",
			"name": "喀喇沁旗"
		},
		{
			"id": "359",
			"city_id": "28",
			"name": "宁城县"
		},
		{
			"id": "360",
			"city_id": "28",
			"name": "敖汉旗"
		},
		{
			"id": "361",
			"city_id": "29",
			"name": "科尔沁区"
		},
		{
			"id": "362",
			"city_id": "29",
			"name": "科尔沁左翼中旗"
		},
		{
			"id": "363",
			"city_id": "29",
			"name": "科尔沁左翼后旗"
		},
		{
			"id": "364",
			"city_id": "29",
			"name": "开鲁县"
		},
		{
			"id": "365",
			"city_id": "29",
			"name": "库伦旗"
		},
		{
			"id": "366",
			"city_id": "29",
			"name": "奈曼旗"
		},
		{
			"id": "367",
			"city_id": "29",
			"name": "扎鲁特旗"
		},
		{
			"id": "368",
			"city_id": "29",
			"name": "霍林郭勒市"
		},
		{
			"id": "369",
			"city_id": "30",
			"name": "东胜区"
		},
		{
			"id": "370",
			"city_id": "30",
			"name": "达拉特旗"
		},
		{
			"id": "371",
			"city_id": "30",
			"name": "准格尔旗"
		},
		{
			"id": "372",
			"city_id": "30",
			"name": "鄂托克前旗"
		},
		{
			"id": "373",
			"city_id": "30",
			"name": "鄂托克旗"
		},
		{
			"id": "374",
			"city_id": "30",
			"name": "杭锦旗"
		},
		{
			"id": "375",
			"city_id": "30",
			"name": "乌审旗"
		},
		{
			"id": "376",
			"city_id": "30",
			"name": "伊金霍洛旗"
		},
		{
			"id": "377",
			"city_id": "31",
			"name": "海拉尔区"
		},
		{
			"id": "378",
			"city_id": "31",
			"name": "阿荣旗"
		},
		{
			"id": "379",
			"city_id": "31",
			"name": "莫力达瓦达斡尔族自治旗"
		},
		{
			"id": "380",
			"city_id": "31",
			"name": "鄂伦春自治旗"
		},
		{
			"id": "381",
			"city_id": "31",
			"name": "鄂温克族自治旗"
		},
		{
			"id": "382",
			"city_id": "31",
			"name": "陈巴尔虎旗"
		},
		{
			"id": "383",
			"city_id": "31",
			"name": "新巴尔虎左旗"
		},
		{
			"id": "384",
			"city_id": "31",
			"name": "新巴尔虎右旗"
		},
		{
			"id": "385",
			"city_id": "31",
			"name": "满洲里市"
		},
		{
			"id": "386",
			"city_id": "31",
			"name": "牙克石市"
		},
		{
			"id": "387",
			"city_id": "31",
			"name": "扎兰屯市"
		},
		{
			"id": "388",
			"city_id": "31",
			"name": "额尔古纳市"
		},
		{
			"id": "389",
			"city_id": "31",
			"name": "根河市"
		},
		{
			"id": "390",
			"city_id": "32",
			"name": "临河区"
		},
		{
			"id": "391",
			"city_id": "32",
			"name": "五原县"
		},
		{
			"id": "392",
			"city_id": "32",
			"name": "磴口县"
		},
		{
			"id": "393",
			"city_id": "32",
			"name": "乌拉特前旗"
		},
		{
			"id": "394",
			"city_id": "32",
			"name": "乌拉特中旗"
		},
		{
			"id": "395",
			"city_id": "32",
			"name": "乌拉特后旗"
		},
		{
			"id": "396",
			"city_id": "32",
			"name": "杭锦后旗"
		},
		{
			"id": "397",
			"city_id": "33",
			"name": "集宁区"
		},
		{
			"id": "398",
			"city_id": "33",
			"name": "卓资县"
		},
		{
			"id": "399",
			"city_id": "33",
			"name": "化德县"
		},
		{
			"id": "400",
			"city_id": "33",
			"name": "商都县"
		},
		{
			"id": "401",
			"city_id": "33",
			"name": "兴和县"
		},
		{
			"id": "402",
			"city_id": "33",
			"name": "凉城县"
		},
		{
			"id": "403",
			"city_id": "33",
			"name": "察哈尔右翼前旗"
		},
		{
			"id": "404",
			"city_id": "33",
			"name": "察哈尔右翼中旗"
		},
		{
			"id": "405",
			"city_id": "33",
			"name": "察哈尔右翼后旗"
		},
		{
			"id": "406",
			"city_id": "33",
			"name": "四子王旗"
		},
		{
			"id": "407",
			"city_id": "33",
			"name": "丰镇市"
		},
		{
			"id": "408",
			"city_id": "34",
			"name": "乌兰浩特市"
		},
		{
			"id": "409",
			"city_id": "34",
			"name": "阿尔山市"
		},
		{
			"id": "410",
			"city_id": "34",
			"name": "科尔沁右翼前旗"
		},
		{
			"id": "411",
			"city_id": "34",
			"name": "科尔沁右翼中旗"
		},
		{
			"id": "412",
			"city_id": "34",
			"name": "扎赉特旗"
		},
		{
			"id": "413",
			"city_id": "34",
			"name": "突泉县"
		},
		{
			"id": "414",
			"city_id": "35",
			"name": "二连浩特市"
		},
		{
			"id": "415",
			"city_id": "35",
			"name": "锡林浩特市"
		},
		{
			"id": "416",
			"city_id": "35",
			"name": "阿巴嘎旗"
		},
		{
			"id": "417",
			"city_id": "35",
			"name": "苏尼特左旗"
		},
		{
			"id": "418",
			"city_id": "35",
			"name": "苏尼特右旗"
		},
		{
			"id": "419",
			"city_id": "35",
			"name": "东乌珠穆沁旗"
		},
		{
			"id": "420",
			"city_id": "35",
			"name": "西乌珠穆沁旗"
		},
		{
			"id": "421",
			"city_id": "35",
			"name": "太仆寺旗"
		},
		{
			"id": "422",
			"city_id": "35",
			"name": "镶黄旗"
		},
		{
			"id": "423",
			"city_id": "35",
			"name": "正镶白旗"
		},
		{
			"id": "424",
			"city_id": "35",
			"name": "正蓝旗"
		},
		{
			"id": "425",
			"city_id": "35",
			"name": "多伦县"
		},
		{
			"id": "426",
			"city_id": "36",
			"name": "阿拉善左旗"
		},
		{
			"id": "427",
			"city_id": "36",
			"name": "阿拉善右旗"
		},
		{
			"id": "428",
			"city_id": "36",
			"name": "额济纳旗"
		},
		{
			"id": "429",
			"city_id": "37",
			"name": "和平区"
		},
		{
			"id": "430",
			"city_id": "37",
			"name": "沈河区"
		},
		{
			"id": "431",
			"city_id": "37",
			"name": "大东区"
		},
		{
			"id": "432",
			"city_id": "37",
			"name": "皇姑区"
		},
		{
			"id": "433",
			"city_id": "37",
			"name": "铁西区"
		},
		{
			"id": "434",
			"city_id": "37",
			"name": "苏家屯区"
		},
		{
			"id": "435",
			"city_id": "37",
			"name": "东陵区"
		},
		{
			"id": "436",
			"city_id": "37",
			"name": "新城子区"
		},
		{
			"id": "437",
			"city_id": "37",
			"name": "于洪区"
		},
		{
			"id": "438",
			"city_id": "37",
			"name": "辽中县"
		},
		{
			"id": "439",
			"city_id": "37",
			"name": "康平县"
		},
		{
			"id": "440",
			"city_id": "37",
			"name": "法库县"
		},
		{
			"id": "441",
			"city_id": "37",
			"name": "新民市"
		},
		{
			"id": "442",
			"city_id": "38",
			"name": "中山区"
		},
		{
			"id": "443",
			"city_id": "38",
			"name": "西岗区"
		},
		{
			"id": "444",
			"city_id": "38",
			"name": "沙河口区"
		},
		{
			"id": "445",
			"city_id": "38",
			"name": "甘井子区"
		},
		{
			"id": "446",
			"city_id": "38",
			"name": "旅顺口区"
		},
		{
			"id": "447",
			"city_id": "38",
			"name": "金州区"
		},
		{
			"id": "448",
			"city_id": "38",
			"name": "长海县"
		},
		{
			"id": "449",
			"city_id": "38",
			"name": "瓦房店市"
		},
		{
			"id": "450",
			"city_id": "38",
			"name": "普兰店市"
		},
		{
			"id": "451",
			"city_id": "38",
			"name": "庄河市"
		},
		{
			"id": "452",
			"city_id": "39",
			"name": "铁东区"
		},
		{
			"id": "453",
			"city_id": "39",
			"name": "铁西区"
		},
		{
			"id": "454",
			"city_id": "39",
			"name": "立山区"
		},
		{
			"id": "455",
			"city_id": "39",
			"name": "千山区"
		},
		{
			"id": "456",
			"city_id": "39",
			"name": "台安县"
		},
		{
			"id": "457",
			"city_id": "39",
			"name": "岫岩满族自治县"
		},
		{
			"id": "458",
			"city_id": "39",
			"name": "海城市"
		},
		{
			"id": "459",
			"city_id": "40",
			"name": "新抚区"
		},
		{
			"id": "460",
			"city_id": "40",
			"name": "东洲区"
		},
		{
			"id": "461",
			"city_id": "40",
			"name": "望花区"
		},
		{
			"id": "462",
			"city_id": "40",
			"name": "顺城区"
		},
		{
			"id": "463",
			"city_id": "40",
			"name": "抚顺县"
		},
		{
			"id": "464",
			"city_id": "40",
			"name": "新宾满族自治县"
		},
		{
			"id": "465",
			"city_id": "40",
			"name": "清原满族自治县"
		},
		{
			"id": "466",
			"city_id": "41",
			"name": "平山区"
		},
		{
			"id": "467",
			"city_id": "41",
			"name": "溪湖区"
		},
		{
			"id": "468",
			"city_id": "41",
			"name": "明山区"
		},
		{
			"id": "469",
			"city_id": "41",
			"name": "南芬区"
		},
		{
			"id": "470",
			"city_id": "41",
			"name": "本溪满族自治县"
		},
		{
			"id": "471",
			"city_id": "41",
			"name": "桓仁满族自治县"
		},
		{
			"id": "472",
			"city_id": "42",
			"name": "元宝区"
		},
		{
			"id": "473",
			"city_id": "42",
			"name": "振兴区"
		},
		{
			"id": "474",
			"city_id": "42",
			"name": "振安区"
		},
		{
			"id": "475",
			"city_id": "42",
			"name": "宽甸满族自治县"
		},
		{
			"id": "476",
			"city_id": "42",
			"name": "东港市"
		},
		{
			"id": "477",
			"city_id": "42",
			"name": "凤城市"
		},
		{
			"id": "478",
			"city_id": "43",
			"name": "古塔区"
		},
		{
			"id": "479",
			"city_id": "43",
			"name": "凌河区"
		},
		{
			"id": "480",
			"city_id": "43",
			"name": "太和区"
		},
		{
			"id": "481",
			"city_id": "43",
			"name": "黑山县"
		},
		{
			"id": "482",
			"city_id": "43",
			"name": "义县"
		},
		{
			"id": "483",
			"city_id": "43",
			"name": "凌海市"
		},
		{
			"id": "484",
			"city_id": "43",
			"name": "北宁市"
		},
		{
			"id": "485",
			"city_id": "44",
			"name": "站前区"
		},
		{
			"id": "486",
			"city_id": "44",
			"name": "西市区"
		},
		{
			"id": "487",
			"city_id": "44",
			"name": "鲅鱼圈区"
		},
		{
			"id": "488",
			"city_id": "44",
			"name": "老边区"
		},
		{
			"id": "489",
			"city_id": "44",
			"name": "盖州市"
		},
		{
			"id": "490",
			"city_id": "44",
			"name": "大石桥市"
		},
		{
			"id": "491",
			"city_id": "45",
			"name": "海州区"
		},
		{
			"id": "492",
			"city_id": "45",
			"name": "新邱区"
		},
		{
			"id": "493",
			"city_id": "45",
			"name": "太平区"
		},
		{
			"id": "494",
			"city_id": "45",
			"name": "清河门区"
		},
		{
			"id": "495",
			"city_id": "45",
			"name": "细河区"
		},
		{
			"id": "496",
			"city_id": "45",
			"name": "阜新蒙古族自治县"
		},
		{
			"id": "497",
			"city_id": "45",
			"name": "彰武县"
		},
		{
			"id": "498",
			"city_id": "46",
			"name": "白塔区"
		},
		{
			"id": "499",
			"city_id": "46",
			"name": "文圣区"
		},
		{
			"id": "500",
			"city_id": "46",
			"name": "宏伟区"
		},
		{
			"id": "501",
			"city_id": "46",
			"name": "弓长岭区"
		},
		{
			"id": "502",
			"city_id": "46",
			"name": "太子河区"
		},
		{
			"id": "503",
			"city_id": "46",
			"name": "辽阳县"
		},
		{
			"id": "504",
			"city_id": "46",
			"name": "灯塔市"
		},
		{
			"id": "505",
			"city_id": "47",
			"name": "双台子区"
		},
		{
			"id": "506",
			"city_id": "47",
			"name": "兴隆台区"
		},
		{
			"id": "507",
			"city_id": "47",
			"name": "大洼县"
		},
		{
			"id": "508",
			"city_id": "47",
			"name": "盘山县"
		},
		{
			"id": "509",
			"city_id": "48",
			"name": "银州区"
		},
		{
			"id": "510",
			"city_id": "48",
			"name": "清河区"
		},
		{
			"id": "511",
			"city_id": "48",
			"name": "铁岭县"
		},
		{
			"id": "512",
			"city_id": "48",
			"name": "西丰县"
		},
		{
			"id": "513",
			"city_id": "48",
			"name": "昌图县"
		},
		{
			"id": "514",
			"city_id": "48",
			"name": "调兵山市"
		},
		{
			"id": "515",
			"city_id": "48",
			"name": "开原市"
		},
		{
			"id": "516",
			"city_id": "49",
			"name": "双塔区"
		},
		{
			"id": "517",
			"city_id": "49",
			"name": "龙城区"
		},
		{
			"id": "518",
			"city_id": "49",
			"name": "朝阳县"
		},
		{
			"id": "519",
			"city_id": "49",
			"name": "建平县"
		},
		{
			"id": "520",
			"city_id": "49",
			"name": "喀喇沁左翼蒙古族自治县"
		},
		{
			"id": "521",
			"city_id": "49",
			"name": "北票市"
		},
		{
			"id": "522",
			"city_id": "49",
			"name": "凌源市"
		},
		{
			"id": "523",
			"city_id": "50",
			"name": "连山区"
		},
		{
			"id": "524",
			"city_id": "50",
			"name": "龙港区"
		},
		{
			"id": "525",
			"city_id": "50",
			"name": "南票区"
		},
		{
			"id": "526",
			"city_id": "50",
			"name": "绥中县"
		},
		{
			"id": "527",
			"city_id": "50",
			"name": "建昌县"
		},
		{
			"id": "528",
			"city_id": "50",
			"name": "兴城市"
		},
		{
			"id": "529",
			"city_id": "51",
			"name": "南关区"
		},
		{
			"id": "530",
			"city_id": "51",
			"name": "宽城区"
		},
		{
			"id": "531",
			"city_id": "51",
			"name": "朝阳区"
		},
		{
			"id": "532",
			"city_id": "51",
			"name": "二道区"
		},
		{
			"id": "533",
			"city_id": "51",
			"name": "绿园区"
		},
		{
			"id": "534",
			"city_id": "51",
			"name": "双阳区"
		},
		{
			"id": "535",
			"city_id": "51",
			"name": "农安县"
		},
		{
			"id": "536",
			"city_id": "51",
			"name": "九台市"
		},
		{
			"id": "537",
			"city_id": "51",
			"name": "榆树市"
		},
		{
			"id": "538",
			"city_id": "51",
			"name": "德惠市"
		},
		{
			"id": "539",
			"city_id": "52",
			"name": "昌邑区"
		},
		{
			"id": "540",
			"city_id": "52",
			"name": "龙潭区"
		},
		{
			"id": "541",
			"city_id": "52",
			"name": "船营区"
		},
		{
			"id": "542",
			"city_id": "52",
			"name": "丰满区"
		},
		{
			"id": "543",
			"city_id": "52",
			"name": "永吉县"
		},
		{
			"id": "544",
			"city_id": "52",
			"name": "蛟河市"
		},
		{
			"id": "545",
			"city_id": "52",
			"name": "桦甸市"
		},
		{
			"id": "546",
			"city_id": "52",
			"name": "舒兰市"
		},
		{
			"id": "547",
			"city_id": "52",
			"name": "磐石市"
		},
		{
			"id": "548",
			"city_id": "53",
			"name": "铁西区"
		},
		{
			"id": "549",
			"city_id": "53",
			"name": "铁东区"
		},
		{
			"id": "550",
			"city_id": "53",
			"name": "梨树县"
		},
		{
			"id": "551",
			"city_id": "53",
			"name": "伊通满族自治县"
		},
		{
			"id": "552",
			"city_id": "53",
			"name": "公主岭市"
		},
		{
			"id": "553",
			"city_id": "53",
			"name": "双辽市"
		},
		{
			"id": "554",
			"city_id": "54",
			"name": "龙山区"
		},
		{
			"id": "555",
			"city_id": "54",
			"name": "西安区"
		},
		{
			"id": "556",
			"city_id": "54",
			"name": "东丰县"
		},
		{
			"id": "557",
			"city_id": "54",
			"name": "东辽县"
		},
		{
			"id": "558",
			"city_id": "55",
			"name": "东昌区"
		},
		{
			"id": "559",
			"city_id": "55",
			"name": "二道江区"
		},
		{
			"id": "560",
			"city_id": "55",
			"name": "通化县"
		},
		{
			"id": "561",
			"city_id": "55",
			"name": "辉南县"
		},
		{
			"id": "562",
			"city_id": "55",
			"name": "柳河县"
		},
		{
			"id": "563",
			"city_id": "55",
			"name": "梅河口市"
		},
		{
			"id": "564",
			"city_id": "55",
			"name": "集安市"
		},
		{
			"id": "565",
			"city_id": "56",
			"name": "八道江区"
		},
		{
			"id": "566",
			"city_id": "56",
			"name": "抚松县"
		},
		{
			"id": "567",
			"city_id": "56",
			"name": "靖宇县"
		},
		{
			"id": "568",
			"city_id": "56",
			"name": "长白朝鲜族自治县"
		},
		{
			"id": "569",
			"city_id": "56",
			"name": "江源县"
		},
		{
			"id": "570",
			"city_id": "56",
			"name": "临江市"
		},
		{
			"id": "571",
			"city_id": "57",
			"name": "宁江区"
		},
		{
			"id": "572",
			"city_id": "57",
			"name": "前郭尔罗斯蒙古族自治县"
		},
		{
			"id": "573",
			"city_id": "57",
			"name": "长岭县"
		},
		{
			"id": "574",
			"city_id": "57",
			"name": "乾安县"
		},
		{
			"id": "575",
			"city_id": "57",
			"name": "扶余县"
		},
		{
			"id": "576",
			"city_id": "58",
			"name": "洮北区"
		},
		{
			"id": "577",
			"city_id": "58",
			"name": "镇赉县"
		},
		{
			"id": "578",
			"city_id": "58",
			"name": "通榆县"
		},
		{
			"id": "579",
			"city_id": "58",
			"name": "洮南市"
		},
		{
			"id": "580",
			"city_id": "58",
			"name": "大安市"
		},
		{
			"id": "581",
			"city_id": "59",
			"name": "延吉市"
		},
		{
			"id": "582",
			"city_id": "59",
			"name": "图们市"
		},
		{
			"id": "583",
			"city_id": "59",
			"name": "敦化市"
		},
		{
			"id": "584",
			"city_id": "59",
			"name": "珲春市"
		},
		{
			"id": "585",
			"city_id": "59",
			"name": "龙井市"
		},
		{
			"id": "586",
			"city_id": "59",
			"name": "和龙市"
		},
		{
			"id": "587",
			"city_id": "59",
			"name": "汪清县"
		},
		{
			"id": "588",
			"city_id": "59",
			"name": "安图县"
		},
		{
			"id": "589",
			"city_id": "60",
			"name": "道里区"
		},
		{
			"id": "590",
			"city_id": "60",
			"name": "南岗区"
		},
		{
			"id": "591",
			"city_id": "60",
			"name": "道外区"
		},
		{
			"id": "592",
			"city_id": "60",
			"name": "香坊区"
		},
		{
			"id": "593",
			"city_id": "60",
			"name": "动力区"
		},
		{
			"id": "594",
			"city_id": "60",
			"name": "平房区"
		},
		{
			"id": "595",
			"city_id": "60",
			"name": "松北区"
		},
		{
			"id": "596",
			"city_id": "60",
			"name": "呼兰区"
		},
		{
			"id": "597",
			"city_id": "60",
			"name": "依兰县"
		},
		{
			"id": "598",
			"city_id": "60",
			"name": "方正县"
		},
		{
			"id": "599",
			"city_id": "60",
			"name": "宾县"
		},
		{
			"id": "600",
			"city_id": "60",
			"name": "巴彦县"
		},
		{
			"id": "601",
			"city_id": "60",
			"name": "木兰县"
		},
		{
			"id": "602",
			"city_id": "60",
			"name": "通河县"
		},
		{
			"id": "603",
			"city_id": "60",
			"name": "延寿县"
		},
		{
			"id": "604",
			"city_id": "60",
			"name": "阿城市"
		},
		{
			"id": "605",
			"city_id": "60",
			"name": "双城市"
		},
		{
			"id": "606",
			"city_id": "60",
			"name": "尚志市"
		},
		{
			"id": "607",
			"city_id": "60",
			"name": "五常市"
		},
		{
			"id": "608",
			"city_id": "61",
			"name": "龙沙区"
		},
		{
			"id": "609",
			"city_id": "61",
			"name": "建华区"
		},
		{
			"id": "610",
			"city_id": "61",
			"name": "铁锋区"
		},
		{
			"id": "611",
			"city_id": "61",
			"name": "昂昂溪区"
		},
		{
			"id": "612",
			"city_id": "61",
			"name": "富拉尔基区"
		},
		{
			"id": "613",
			"city_id": "61",
			"name": "碾子山区"
		},
		{
			"id": "614",
			"city_id": "61",
			"name": "梅里斯达斡尔族区"
		},
		{
			"id": "615",
			"city_id": "61",
			"name": "龙江县"
		},
		{
			"id": "616",
			"city_id": "61",
			"name": "依安县"
		},
		{
			"id": "617",
			"city_id": "61",
			"name": "泰来县"
		},
		{
			"id": "618",
			"city_id": "61",
			"name": "甘南县"
		},
		{
			"id": "619",
			"city_id": "61",
			"name": "富裕县"
		},
		{
			"id": "620",
			"city_id": "61",
			"name": "克山县"
		},
		{
			"id": "621",
			"city_id": "61",
			"name": "克东县"
		},
		{
			"id": "622",
			"city_id": "61",
			"name": "拜泉县"
		},
		{
			"id": "623",
			"city_id": "61",
			"name": "讷河市"
		},
		{
			"id": "624",
			"city_id": "62",
			"name": "鸡冠区"
		},
		{
			"id": "625",
			"city_id": "62",
			"name": "恒山区"
		},
		{
			"id": "626",
			"city_id": "62",
			"name": "滴道区"
		},
		{
			"id": "627",
			"city_id": "62",
			"name": "梨树区"
		},
		{
			"id": "628",
			"city_id": "62",
			"name": "城子河区"
		},
		{
			"id": "629",
			"city_id": "62",
			"name": "麻山区"
		},
		{
			"id": "630",
			"city_id": "62",
			"name": "鸡东县"
		},
		{
			"id": "631",
			"city_id": "62",
			"name": "虎林市"
		},
		{
			"id": "632",
			"city_id": "62",
			"name": "密山市"
		},
		{
			"id": "633",
			"city_id": "63",
			"name": "向阳区"
		},
		{
			"id": "634",
			"city_id": "63",
			"name": "工农区"
		},
		{
			"id": "635",
			"city_id": "63",
			"name": "南山区"
		},
		{
			"id": "636",
			"city_id": "63",
			"name": "兴安区"
		},
		{
			"id": "637",
			"city_id": "63",
			"name": "东山区"
		},
		{
			"id": "638",
			"city_id": "63",
			"name": "兴山区"
		},
		{
			"id": "639",
			"city_id": "63",
			"name": "萝北县"
		},
		{
			"id": "640",
			"city_id": "63",
			"name": "绥滨县"
		},
		{
			"id": "641",
			"city_id": "64",
			"name": "尖山区"
		},
		{
			"id": "642",
			"city_id": "64",
			"name": "岭东区"
		},
		{
			"id": "643",
			"city_id": "64",
			"name": "四方台区"
		},
		{
			"id": "644",
			"city_id": "64",
			"name": "宝山区"
		},
		{
			"id": "645",
			"city_id": "64",
			"name": "集贤县"
		},
		{
			"id": "646",
			"city_id": "64",
			"name": "友谊县"
		},
		{
			"id": "647",
			"city_id": "64",
			"name": "宝清县"
		},
		{
			"id": "648",
			"city_id": "64",
			"name": "饶河县"
		},
		{
			"id": "649",
			"city_id": "65",
			"name": "萨尔图区"
		},
		{
			"id": "650",
			"city_id": "65",
			"name": "龙凤区"
		},
		{
			"id": "651",
			"city_id": "65",
			"name": "让胡路区"
		},
		{
			"id": "652",
			"city_id": "65",
			"name": "红岗区"
		},
		{
			"id": "653",
			"city_id": "65",
			"name": "大同区"
		},
		{
			"id": "654",
			"city_id": "65",
			"name": "肇州县"
		},
		{
			"id": "655",
			"city_id": "65",
			"name": "肇源县"
		},
		{
			"id": "656",
			"city_id": "65",
			"name": "林甸县"
		},
		{
			"id": "657",
			"city_id": "65",
			"name": "杜尔伯特蒙古族自治县"
		},
		{
			"id": "658",
			"city_id": "66",
			"name": "伊春区"
		},
		{
			"id": "659",
			"city_id": "66",
			"name": "南岔区"
		},
		{
			"id": "660",
			"city_id": "66",
			"name": "友好区"
		},
		{
			"id": "661",
			"city_id": "66",
			"name": "西林区"
		},
		{
			"id": "662",
			"city_id": "66",
			"name": "翠峦区"
		},
		{
			"id": "663",
			"city_id": "66",
			"name": "新青区"
		},
		{
			"id": "664",
			"city_id": "66",
			"name": "美溪区"
		},
		{
			"id": "665",
			"city_id": "66",
			"name": "金山屯区"
		},
		{
			"id": "666",
			"city_id": "66",
			"name": "五营区"
		},
		{
			"id": "667",
			"city_id": "66",
			"name": "乌马河区"
		},
		{
			"id": "668",
			"city_id": "66",
			"name": "汤旺河区"
		},
		{
			"id": "669",
			"city_id": "66",
			"name": "带岭区"
		},
		{
			"id": "670",
			"city_id": "66",
			"name": "乌伊岭区"
		},
		{
			"id": "671",
			"city_id": "66",
			"name": "红星区"
		},
		{
			"id": "672",
			"city_id": "66",
			"name": "上甘岭区"
		},
		{
			"id": "673",
			"city_id": "66",
			"name": "嘉荫县"
		},
		{
			"id": "674",
			"city_id": "66",
			"name": "铁力市"
		},
		{
			"id": "675",
			"city_id": "67",
			"name": "永红区"
		},
		{
			"id": "676",
			"city_id": "67",
			"name": "向阳区"
		},
		{
			"id": "677",
			"city_id": "67",
			"name": "前进区"
		},
		{
			"id": "678",
			"city_id": "67",
			"name": "东风区"
		},
		{
			"id": "679",
			"city_id": "67",
			"name": "郊区"
		},
		{
			"id": "680",
			"city_id": "67",
			"name": "桦南县"
		},
		{
			"id": "681",
			"city_id": "67",
			"name": "桦川县"
		},
		{
			"id": "682",
			"city_id": "67",
			"name": "汤原县"
		},
		{
			"id": "683",
			"city_id": "67",
			"name": "抚远县"
		},
		{
			"id": "684",
			"city_id": "67",
			"name": "同江市"
		},
		{
			"id": "685",
			"city_id": "67",
			"name": "富锦市"
		},
		{
			"id": "686",
			"city_id": "68",
			"name": "新兴区"
		},
		{
			"id": "687",
			"city_id": "68",
			"name": "桃山区"
		},
		{
			"id": "688",
			"city_id": "68",
			"name": "茄子河区"
		},
		{
			"id": "689",
			"city_id": "68",
			"name": "勃利县"
		},
		{
			"id": "690",
			"city_id": "69",
			"name": "东安区"
		},
		{
			"id": "691",
			"city_id": "69",
			"name": "阳明区"
		},
		{
			"id": "692",
			"city_id": "69",
			"name": "爱民区"
		},
		{
			"id": "693",
			"city_id": "69",
			"name": "西安区"
		},
		{
			"id": "694",
			"city_id": "69",
			"name": "东宁县"
		},
		{
			"id": "695",
			"city_id": "69",
			"name": "林口县"
		},
		{
			"id": "696",
			"city_id": "69",
			"name": "绥芬河市"
		},
		{
			"id": "697",
			"city_id": "69",
			"name": "海林市"
		},
		{
			"id": "698",
			"city_id": "69",
			"name": "宁安市"
		},
		{
			"id": "699",
			"city_id": "69",
			"name": "穆棱市"
		},
		{
			"id": "700",
			"city_id": "70",
			"name": "爱辉区"
		},
		{
			"id": "701",
			"city_id": "70",
			"name": "嫩江县"
		},
		{
			"id": "702",
			"city_id": "70",
			"name": "逊克县"
		},
		{
			"id": "703",
			"city_id": "70",
			"name": "孙吴县"
		},
		{
			"id": "704",
			"city_id": "70",
			"name": "北安市"
		},
		{
			"id": "705",
			"city_id": "70",
			"name": "五大连池市"
		},
		{
			"id": "706",
			"city_id": "71",
			"name": "北林区"
		},
		{
			"id": "707",
			"city_id": "71",
			"name": "望奎县"
		},
		{
			"id": "708",
			"city_id": "71",
			"name": "兰西县"
		},
		{
			"id": "709",
			"city_id": "71",
			"name": "青冈县"
		},
		{
			"id": "710",
			"city_id": "71",
			"name": "庆安县"
		},
		{
			"id": "711",
			"city_id": "71",
			"name": "明水县"
		},
		{
			"id": "712",
			"city_id": "71",
			"name": "绥棱县"
		},
		{
			"id": "713",
			"city_id": "71",
			"name": "安达市"
		},
		{
			"id": "714",
			"city_id": "71",
			"name": "肇东市"
		},
		{
			"id": "715",
			"city_id": "71",
			"name": "海伦市"
		},
		{
			"id": "716",
			"city_id": "72",
			"name": "呼玛县"
		},
		{
			"id": "717",
			"city_id": "72",
			"name": "塔河县"
		},
		{
			"id": "718",
			"city_id": "72",
			"name": "漠河县"
		},
		{
			"id": "719",
			"city_id": "73",
			"name": "黄浦区"
		},
		{
			"id": "720",
			"city_id": "73",
			"name": "卢湾区"
		},
		{
			"id": "721",
			"city_id": "73",
			"name": "徐汇区"
		},
		{
			"id": "722",
			"city_id": "73",
			"name": "长宁区"
		},
		{
			"id": "723",
			"city_id": "73",
			"name": "静安区"
		},
		{
			"id": "724",
			"city_id": "73",
			"name": "普陀区"
		},
		{
			"id": "725",
			"city_id": "73",
			"name": "闸北区"
		},
		{
			"id": "726",
			"city_id": "73",
			"name": "虹口区"
		},
		{
			"id": "727",
			"city_id": "73",
			"name": "杨浦区"
		},
		{
			"id": "728",
			"city_id": "73",
			"name": "闵行区"
		},
		{
			"id": "729",
			"city_id": "73",
			"name": "宝山区"
		},
		{
			"id": "730",
			"city_id": "73",
			"name": "嘉定区"
		},
		{
			"id": "731",
			"city_id": "73",
			"name": "浦东新区"
		},
		{
			"id": "732",
			"city_id": "73",
			"name": "金山区"
		},
		{
			"id": "733",
			"city_id": "73",
			"name": "松江区"
		},
		{
			"id": "734",
			"city_id": "73",
			"name": "青浦区"
		},
		{
			"id": "735",
			"city_id": "73",
			"name": "南汇区"
		},
		{
			"id": "736",
			"city_id": "73",
			"name": "奉贤区"
		},
		{
			"id": "737",
			"city_id": "73",
			"name": "崇明县"
		},
		{
			"id": "738",
			"city_id": "74",
			"name": "玄武区"
		},
		{
			"id": "740",
			"city_id": "74",
			"name": "秦淮区"
		},
		{
			"id": "741",
			"city_id": "74",
			"name": "建邺区"
		},
		{
			"id": "742",
			"city_id": "74",
			"name": "鼓楼区"
		},
		{
			"id": "743",
			"city_id": "74",
			"name": "下关区"
		},
		{
			"id": "744",
			"city_id": "74",
			"name": "浦口区"
		},
		{
			"id": "745",
			"city_id": "74",
			"name": "栖霞区"
		},
		{
			"id": "746",
			"city_id": "74",
			"name": "雨花台区"
		},
		{
			"id": "747",
			"city_id": "74",
			"name": "江宁区"
		},
		{
			"id": "748",
			"city_id": "74",
			"name": "六合区"
		},
		{
			"id": "749",
			"city_id": "74",
			"name": "溧水县"
		},
		{
			"id": "750",
			"city_id": "74",
			"name": "高淳县"
		},
		{
			"id": "751",
			"city_id": "75",
			"name": "崇安区"
		},
		{
			"id": "752",
			"city_id": "75",
			"name": "南长区"
		},
		{
			"id": "753",
			"city_id": "75",
			"name": "北塘区"
		},
		{
			"id": "754",
			"city_id": "75",
			"name": "锡山区"
		},
		{
			"id": "755",
			"city_id": "75",
			"name": "惠山区"
		},
		{
			"id": "756",
			"city_id": "75",
			"name": "滨湖区"
		},
		{
			"id": "757",
			"city_id": "75",
			"name": "江阴市"
		},
		{
			"id": "758",
			"city_id": "75",
			"name": "宜兴市"
		},
		{
			"id": "759",
			"city_id": "76",
			"name": "鼓楼区"
		},
		{
			"id": "760",
			"city_id": "76",
			"name": "云龙区"
		},
		{
			"id": "761",
			"city_id": "76",
			"name": "九里区"
		},
		{
			"id": "762",
			"city_id": "76",
			"name": "贾汪区"
		},
		{
			"id": "763",
			"city_id": "76",
			"name": "泉山区"
		},
		{
			"id": "764",
			"city_id": "76",
			"name": "丰县"
		},
		{
			"id": "765",
			"city_id": "76",
			"name": "沛县"
		},
		{
			"id": "766",
			"city_id": "76",
			"name": "铜山县"
		},
		{
			"id": "767",
			"city_id": "76",
			"name": "睢宁县"
		},
		{
			"id": "768",
			"city_id": "76",
			"name": "新沂市"
		},
		{
			"id": "769",
			"city_id": "76",
			"name": "邳州市"
		},
		{
			"id": "770",
			"city_id": "77",
			"name": "天宁区"
		},
		{
			"id": "771",
			"city_id": "77",
			"name": "钟楼区"
		},
		{
			"id": "772",
			"city_id": "77",
			"name": "戚墅堰区"
		},
		{
			"id": "773",
			"city_id": "77",
			"name": "新北区"
		},
		{
			"id": "774",
			"city_id": "77",
			"name": "武进区"
		},
		{
			"id": "775",
			"city_id": "77",
			"name": "溧阳市"
		},
		{
			"id": "776",
			"city_id": "77",
			"name": "金坛市"
		},
		{
			"id": "777",
			"city_id": "78",
			"name": "沧浪区"
		},
		{
			"id": "778",
			"city_id": "78",
			"name": "平江区"
		},
		{
			"id": "779",
			"city_id": "78",
			"name": "金阊区"
		},
		{
			"id": "780",
			"city_id": "78",
			"name": "虎丘区"
		},
		{
			"id": "781",
			"city_id": "78",
			"name": "吴中区"
		},
		{
			"id": "782",
			"city_id": "78",
			"name": "相城区"
		},
		{
			"id": "783",
			"city_id": "78",
			"name": "常熟市"
		},
		{
			"id": "784",
			"city_id": "78",
			"name": "张家港市"
		},
		{
			"id": "785",
			"city_id": "78",
			"name": "昆山市"
		},
		{
			"id": "786",
			"city_id": "78",
			"name": "吴江市"
		},
		{
			"id": "787",
			"city_id": "78",
			"name": "太仓市"
		},
		{
			"id": "788",
			"city_id": "79",
			"name": "崇川区"
		},
		{
			"id": "789",
			"city_id": "79",
			"name": "港闸区"
		},
		{
			"id": "790",
			"city_id": "79",
			"name": "海安县"
		},
		{
			"id": "791",
			"city_id": "79",
			"name": "如东县"
		},
		{
			"id": "792",
			"city_id": "79",
			"name": "启东市"
		},
		{
			"id": "793",
			"city_id": "79",
			"name": "如皋市"
		},
		{
			"id": "794",
			"city_id": "79",
			"name": "通州市"
		},
		{
			"id": "795",
			"city_id": "79",
			"name": "海门市"
		},
		{
			"id": "796",
			"city_id": "80",
			"name": "连云区"
		},
		{
			"id": "797",
			"city_id": "80",
			"name": "新浦区"
		},
		{
			"id": "798",
			"city_id": "80",
			"name": "海州区"
		},
		{
			"id": "799",
			"city_id": "80",
			"name": "赣榆县"
		},
		{
			"id": "800",
			"city_id": "80",
			"name": "东海县"
		},
		{
			"id": "801",
			"city_id": "80",
			"name": "灌云县"
		},
		{
			"id": "802",
			"city_id": "80",
			"name": "灌南县"
		},
		{
			"id": "803",
			"city_id": "81",
			"name": "清河区"
		},
		{
			"id": "804",
			"city_id": "81",
			"name": "楚州区"
		},
		{
			"id": "805",
			"city_id": "81",
			"name": "淮阴区"
		},
		{
			"id": "806",
			"city_id": "81",
			"name": "清浦区"
		},
		{
			"id": "807",
			"city_id": "81",
			"name": "涟水县"
		},
		{
			"id": "808",
			"city_id": "81",
			"name": "洪泽县"
		},
		{
			"id": "809",
			"city_id": "81",
			"name": "盱眙县"
		},
		{
			"id": "810",
			"city_id": "81",
			"name": "金湖县"
		},
		{
			"id": "811",
			"city_id": "82",
			"name": "亭湖区"
		},
		{
			"id": "812",
			"city_id": "82",
			"name": "盐都区"
		},
		{
			"id": "813",
			"city_id": "82",
			"name": "响水县"
		},
		{
			"id": "814",
			"city_id": "82",
			"name": "滨海县"
		},
		{
			"id": "815",
			"city_id": "82",
			"name": "阜宁县"
		},
		{
			"id": "816",
			"city_id": "82",
			"name": "射阳县"
		},
		{
			"id": "817",
			"city_id": "82",
			"name": "建湖县"
		},
		{
			"id": "818",
			"city_id": "82",
			"name": "东台市"
		},
		{
			"id": "819",
			"city_id": "82",
			"name": "大丰市"
		},
		{
			"id": "820",
			"city_id": "83",
			"name": "广陵区"
		},
		{
			"id": "821",
			"city_id": "83",
			"name": "邗江区"
		},
		{
			"id": "822",
			"city_id": "83",
			"name": "维扬区"
		},
		{
			"id": "823",
			"city_id": "83",
			"name": "宝应县"
		},
		{
			"id": "824",
			"city_id": "83",
			"name": "仪征市"
		},
		{
			"id": "825",
			"city_id": "83",
			"name": "高邮市"
		},
		{
			"id": "826",
			"city_id": "83",
			"name": "江都市"
		},
		{
			"id": "827",
			"city_id": "84",
			"name": "京口区"
		},
		{
			"id": "828",
			"city_id": "84",
			"name": "润州区"
		},
		{
			"id": "829",
			"city_id": "84",
			"name": "丹徒区"
		},
		{
			"id": "830",
			"city_id": "84",
			"name": "丹阳市"
		},
		{
			"id": "831",
			"city_id": "84",
			"name": "扬中市"
		},
		{
			"id": "832",
			"city_id": "84",
			"name": "句容市"
		},
		{
			"id": "833",
			"city_id": "85",
			"name": "海陵区"
		},
		{
			"id": "834",
			"city_id": "85",
			"name": "高港区"
		},
		{
			"id": "835",
			"city_id": "85",
			"name": "兴化市"
		},
		{
			"id": "836",
			"city_id": "85",
			"name": "靖江市"
		},
		{
			"id": "837",
			"city_id": "85",
			"name": "泰兴市"
		},
		{
			"id": "838",
			"city_id": "85",
			"name": "姜堰市"
		},
		{
			"id": "839",
			"city_id": "86",
			"name": "宿城区"
		},
		{
			"id": "840",
			"city_id": "86",
			"name": "宿豫区"
		},
		{
			"id": "841",
			"city_id": "86",
			"name": "沭阳县"
		},
		{
			"id": "842",
			"city_id": "86",
			"name": "泗阳县"
		},
		{
			"id": "843",
			"city_id": "86",
			"name": "泗洪县"
		},
		{
			"id": "844",
			"city_id": "87",
			"name": "上城区"
		},
		{
			"id": "845",
			"city_id": "87",
			"name": "下城区"
		},
		{
			"id": "846",
			"city_id": "87",
			"name": "江干区"
		},
		{
			"id": "847",
			"city_id": "87",
			"name": "拱墅区"
		},
		{
			"id": "848",
			"city_id": "87",
			"name": "西湖区"
		},
		{
			"id": "849",
			"city_id": "87",
			"name": "滨江区"
		},
		{
			"id": "850",
			"city_id": "87",
			"name": "萧山区"
		},
		{
			"id": "851",
			"city_id": "87",
			"name": "余杭区"
		},
		{
			"id": "852",
			"city_id": "87",
			"name": "桐庐县"
		},
		{
			"id": "853",
			"city_id": "87",
			"name": "淳安县"
		},
		{
			"id": "854",
			"city_id": "87",
			"name": "建德市"
		},
		{
			"id": "855",
			"city_id": "87",
			"name": "富阳市"
		},
		{
			"id": "856",
			"city_id": "87",
			"name": "临安市"
		},
		{
			"id": "857",
			"city_id": "88",
			"name": "海曙区"
		},
		{
			"id": "858",
			"city_id": "88",
			"name": "江东区"
		},
		{
			"id": "859",
			"city_id": "88",
			"name": "江北区"
		},
		{
			"id": "860",
			"city_id": "88",
			"name": "北仑区"
		},
		{
			"id": "861",
			"city_id": "88",
			"name": "镇海区"
		},
		{
			"id": "862",
			"city_id": "88",
			"name": "鄞州区"
		},
		{
			"id": "863",
			"city_id": "88",
			"name": "象山县"
		},
		{
			"id": "864",
			"city_id": "88",
			"name": "宁海县"
		},
		{
			"id": "865",
			"city_id": "88",
			"name": "余姚市"
		},
		{
			"id": "866",
			"city_id": "88",
			"name": "慈溪市"
		},
		{
			"id": "867",
			"city_id": "88",
			"name": "奉化市"
		},
		{
			"id": "868",
			"city_id": "89",
			"name": "鹿城区"
		},
		{
			"id": "869",
			"city_id": "89",
			"name": "龙湾区"
		},
		{
			"id": "870",
			"city_id": "89",
			"name": "瓯海区"
		},
		{
			"id": "871",
			"city_id": "89",
			"name": "洞头县"
		},
		{
			"id": "872",
			"city_id": "89",
			"name": "永嘉县"
		},
		{
			"id": "873",
			"city_id": "89",
			"name": "平阳县"
		},
		{
			"id": "874",
			"city_id": "89",
			"name": "苍南县"
		},
		{
			"id": "875",
			"city_id": "89",
			"name": "文成县"
		},
		{
			"id": "876",
			"city_id": "89",
			"name": "泰顺县"
		},
		{
			"id": "877",
			"city_id": "89",
			"name": "瑞安市"
		},
		{
			"id": "878",
			"city_id": "89",
			"name": "乐清市"
		},
		{
			"id": "879",
			"city_id": "90",
			"name": "秀城区"
		},
		{
			"id": "880",
			"city_id": "90",
			"name": "秀洲区"
		},
		{
			"id": "881",
			"city_id": "90",
			"name": "嘉善县"
		},
		{
			"id": "882",
			"city_id": "90",
			"name": "海盐县"
		},
		{
			"id": "883",
			"city_id": "90",
			"name": "海宁市"
		},
		{
			"id": "884",
			"city_id": "90",
			"name": "平湖市"
		},
		{
			"id": "885",
			"city_id": "90",
			"name": "桐乡市"
		},
		{
			"id": "886",
			"city_id": "91",
			"name": "吴兴区"
		},
		{
			"id": "887",
			"city_id": "91",
			"name": "南浔区"
		},
		{
			"id": "888",
			"city_id": "91",
			"name": "德清县"
		},
		{
			"id": "889",
			"city_id": "91",
			"name": "长兴县"
		},
		{
			"id": "890",
			"city_id": "91",
			"name": "安吉县"
		},
		{
			"id": "891",
			"city_id": "92",
			"name": "越城区"
		},
		{
			"id": "892",
			"city_id": "92",
			"name": "绍兴县"
		},
		{
			"id": "893",
			"city_id": "92",
			"name": "新昌县"
		},
		{
			"id": "894",
			"city_id": "92",
			"name": "诸暨市"
		},
		{
			"id": "895",
			"city_id": "92",
			"name": "上虞市"
		},
		{
			"id": "896",
			"city_id": "92",
			"name": "嵊州市"
		},
		{
			"id": "897",
			"city_id": "93",
			"name": "婺城区"
		},
		{
			"id": "898",
			"city_id": "93",
			"name": "金东区"
		},
		{
			"id": "899",
			"city_id": "93",
			"name": "武义县"
		},
		{
			"id": "900",
			"city_id": "93",
			"name": "浦江县"
		},
		{
			"id": "901",
			"city_id": "93",
			"name": "磐安县"
		},
		{
			"id": "902",
			"city_id": "93",
			"name": "兰溪市"
		},
		{
			"id": "903",
			"city_id": "93",
			"name": "义乌市"
		},
		{
			"id": "904",
			"city_id": "93",
			"name": "东阳市"
		},
		{
			"id": "905",
			"city_id": "93",
			"name": "永康市"
		},
		{
			"id": "906",
			"city_id": "94",
			"name": "柯城区"
		},
		{
			"id": "907",
			"city_id": "94",
			"name": "衢江区"
		},
		{
			"id": "908",
			"city_id": "94",
			"name": "常山县"
		},
		{
			"id": "909",
			"city_id": "94",
			"name": "开化县"
		},
		{
			"id": "910",
			"city_id": "94",
			"name": "龙游县"
		},
		{
			"id": "911",
			"city_id": "94",
			"name": "江山市"
		},
		{
			"id": "912",
			"city_id": "95",
			"name": "定海区"
		},
		{
			"id": "913",
			"city_id": "95",
			"name": "普陀区"
		},
		{
			"id": "914",
			"city_id": "95",
			"name": "岱山县"
		},
		{
			"id": "915",
			"city_id": "95",
			"name": "嵊泗县"
		},
		{
			"id": "916",
			"city_id": "96",
			"name": "椒江区"
		},
		{
			"id": "917",
			"city_id": "96",
			"name": "黄岩区"
		},
		{
			"id": "918",
			"city_id": "96",
			"name": "路桥区"
		},
		{
			"id": "919",
			"city_id": "96",
			"name": "玉环县"
		},
		{
			"id": "920",
			"city_id": "96",
			"name": "三门县"
		},
		{
			"id": "921",
			"city_id": "96",
			"name": "天台县"
		},
		{
			"id": "922",
			"city_id": "96",
			"name": "仙居县"
		},
		{
			"id": "923",
			"city_id": "96",
			"name": "温岭市"
		},
		{
			"id": "924",
			"city_id": "96",
			"name": "临海市"
		},
		{
			"id": "925",
			"city_id": "97",
			"name": "莲都区"
		},
		{
			"id": "926",
			"city_id": "97",
			"name": "青田县"
		},
		{
			"id": "927",
			"city_id": "97",
			"name": "缙云县"
		},
		{
			"id": "928",
			"city_id": "97",
			"name": "遂昌县"
		},
		{
			"id": "929",
			"city_id": "97",
			"name": "松阳县"
		},
		{
			"id": "930",
			"city_id": "97",
			"name": "云和县"
		},
		{
			"id": "931",
			"city_id": "97",
			"name": "庆元县"
		},
		{
			"id": "932",
			"city_id": "97",
			"name": "景宁畲族自治县"
		},
		{
			"id": "933",
			"city_id": "97",
			"name": "龙泉市"
		},
		{
			"id": "934",
			"city_id": "98",
			"name": "瑶海区"
		},
		{
			"id": "935",
			"city_id": "98",
			"name": "庐阳区"
		},
		{
			"id": "936",
			"city_id": "98",
			"name": "蜀山区"
		},
		{
			"id": "937",
			"city_id": "98",
			"name": "包河区"
		},
		{
			"id": "938",
			"city_id": "98",
			"name": "长丰县"
		},
		{
			"id": "939",
			"city_id": "98",
			"name": "肥东县"
		},
		{
			"id": "940",
			"city_id": "98",
			"name": "肥西县"
		},
		{
			"id": "941",
			"city_id": "99",
			"name": "镜湖区"
		},
		{
			"id": "942",
			"city_id": "99",
			"name": "马塘区"
		},
		{
			"id": "943",
			"city_id": "99",
			"name": "新芜区"
		},
		{
			"id": "944",
			"city_id": "99",
			"name": "鸠江区"
		},
		{
			"id": "945",
			"city_id": "99",
			"name": "芜湖县"
		},
		{
			"id": "946",
			"city_id": "99",
			"name": "繁昌县"
		},
		{
			"id": "947",
			"city_id": "99",
			"name": "南陵县"
		},
		{
			"id": "948",
			"city_id": "100",
			"name": "龙子湖区"
		},
		{
			"id": "949",
			"city_id": "100",
			"name": "蚌山区"
		},
		{
			"id": "950",
			"city_id": "100",
			"name": "禹会区"
		},
		{
			"id": "951",
			"city_id": "100",
			"name": "淮上区"
		},
		{
			"id": "952",
			"city_id": "100",
			"name": "怀远县"
		},
		{
			"id": "953",
			"city_id": "100",
			"name": "五河县"
		},
		{
			"id": "954",
			"city_id": "100",
			"name": "固镇县"
		},
		{
			"id": "955",
			"city_id": "101",
			"name": "大通区"
		},
		{
			"id": "956",
			"city_id": "101",
			"name": "田家庵区"
		},
		{
			"id": "957",
			"city_id": "101",
			"name": "谢家集区"
		},
		{
			"id": "958",
			"city_id": "101",
			"name": "八公山区"
		},
		{
			"id": "959",
			"city_id": "101",
			"name": "潘集区"
		},
		{
			"id": "960",
			"city_id": "101",
			"name": "凤台县"
		},
		{
			"id": "961",
			"city_id": "102",
			"name": "金家庄区"
		},
		{
			"id": "962",
			"city_id": "102",
			"name": "花山区"
		},
		{
			"id": "963",
			"city_id": "102",
			"name": "雨山区"
		},
		{
			"id": "964",
			"city_id": "102",
			"name": "当涂县"
		},
		{
			"id": "965",
			"city_id": "103",
			"name": "杜集区"
		},
		{
			"id": "966",
			"city_id": "103",
			"name": "相山区"
		},
		{
			"id": "967",
			"city_id": "103",
			"name": "烈山区"
		},
		{
			"id": "968",
			"city_id": "103",
			"name": "濉溪县"
		},
		{
			"id": "969",
			"city_id": "104",
			"name": "铜官山区"
		},
		{
			"id": "970",
			"city_id": "104",
			"name": "狮子山区"
		},
		{
			"id": "971",
			"city_id": "104",
			"name": "郊区"
		},
		{
			"id": "972",
			"city_id": "104",
			"name": "铜陵县"
		},
		{
			"id": "973",
			"city_id": "105",
			"name": "迎江区"
		},
		{
			"id": "974",
			"city_id": "105",
			"name": "大观区"
		},
		{
			"id": "975",
			"city_id": "105",
			"name": "郊区"
		},
		{
			"id": "976",
			"city_id": "105",
			"name": "怀宁县"
		},
		{
			"id": "977",
			"city_id": "105",
			"name": "枞阳县"
		},
		{
			"id": "978",
			"city_id": "105",
			"name": "潜山县"
		},
		{
			"id": "979",
			"city_id": "105",
			"name": "太湖县"
		},
		{
			"id": "980",
			"city_id": "105",
			"name": "宿松县"
		},
		{
			"id": "981",
			"city_id": "105",
			"name": "望江县"
		},
		{
			"id": "982",
			"city_id": "105",
			"name": "岳西县"
		},
		{
			"id": "983",
			"city_id": "105",
			"name": "桐城市"
		},
		{
			"id": "984",
			"city_id": "106",
			"name": "屯溪区"
		},
		{
			"id": "985",
			"city_id": "106",
			"name": "黄山区"
		},
		{
			"id": "986",
			"city_id": "106",
			"name": "徽州区"
		},
		{
			"id": "987",
			"city_id": "106",
			"name": "歙县"
		},
		{
			"id": "988",
			"city_id": "106",
			"name": "休宁县"
		},
		{
			"id": "989",
			"city_id": "106",
			"name": "黟县"
		},
		{
			"id": "990",
			"city_id": "106",
			"name": "祁门县"
		},
		{
			"id": "991",
			"city_id": "107",
			"name": "琅琊区"
		},
		{
			"id": "992",
			"city_id": "107",
			"name": "南谯区"
		},
		{
			"id": "993",
			"city_id": "107",
			"name": "来安县"
		},
		{
			"id": "994",
			"city_id": "107",
			"name": "全椒县"
		},
		{
			"id": "995",
			"city_id": "107",
			"name": "定远县"
		},
		{
			"id": "996",
			"city_id": "107",
			"name": "凤阳县"
		},
		{
			"id": "997",
			"city_id": "107",
			"name": "天长市"
		},
		{
			"id": "998",
			"city_id": "107",
			"name": "明光市"
		},
		{
			"id": "999",
			"city_id": "108",
			"name": "颍州区"
		},
		{
			"id": "1000",
			"city_id": "108",
			"name": "颍东区"
		},
		{
			"id": "1001",
			"city_id": "108",
			"name": "颍泉区"
		},
		{
			"id": "1002",
			"city_id": "108",
			"name": "临泉县"
		},
		{
			"id": "1003",
			"city_id": "108",
			"name": "太和县"
		},
		{
			"id": "1004",
			"city_id": "108",
			"name": "阜南县"
		},
		{
			"id": "1005",
			"city_id": "108",
			"name": "颍上县"
		},
		{
			"id": "1006",
			"city_id": "108",
			"name": "界首市"
		},
		{
			"id": "1007",
			"city_id": "109",
			"name": "埇桥区"
		},
		{
			"id": "1008",
			"city_id": "109",
			"name": "砀山县"
		},
		{
			"id": "1009",
			"city_id": "109",
			"name": "萧县"
		},
		{
			"id": "1010",
			"city_id": "109",
			"name": "灵璧县"
		},
		{
			"id": "1011",
			"city_id": "109",
			"name": "泗县"
		},
		{
			"id": "1012",
			"city_id": "110",
			"name": "居巢区"
		},
		{
			"id": "1013",
			"city_id": "110",
			"name": "庐江县"
		},
		{
			"id": "1014",
			"city_id": "110",
			"name": "无为县"
		},
		{
			"id": "1015",
			"city_id": "110",
			"name": "含山县"
		},
		{
			"id": "1016",
			"city_id": "110",
			"name": "和县"
		},
		{
			"id": "1017",
			"city_id": "111",
			"name": "金安区"
		},
		{
			"id": "1018",
			"city_id": "111",
			"name": "裕安区"
		},
		{
			"id": "1019",
			"city_id": "111",
			"name": "寿县"
		},
		{
			"id": "1020",
			"city_id": "111",
			"name": "霍邱县"
		},
		{
			"id": "1021",
			"city_id": "111",
			"name": "舒城县"
		},
		{
			"id": "1022",
			"city_id": "111",
			"name": "金寨县"
		},
		{
			"id": "1023",
			"city_id": "111",
			"name": "霍山县"
		},
		{
			"id": "1024",
			"city_id": "112",
			"name": "谯城区"
		},
		{
			"id": "1025",
			"city_id": "112",
			"name": "涡阳县"
		},
		{
			"id": "1026",
			"city_id": "112",
			"name": "蒙城县"
		},
		{
			"id": "1027",
			"city_id": "112",
			"name": "利辛县"
		},
		{
			"id": "1028",
			"city_id": "113",
			"name": "贵池区"
		},
		{
			"id": "1029",
			"city_id": "113",
			"name": "东至县"
		},
		{
			"id": "1030",
			"city_id": "113",
			"name": "石台县"
		},
		{
			"id": "1031",
			"city_id": "113",
			"name": "青阳县"
		},
		{
			"id": "1032",
			"city_id": "114",
			"name": "宣州区"
		},
		{
			"id": "1033",
			"city_id": "114",
			"name": "郎溪县"
		},
		{
			"id": "1034",
			"city_id": "114",
			"name": "广德县"
		},
		{
			"id": "1035",
			"city_id": "114",
			"name": "泾县"
		},
		{
			"id": "1036",
			"city_id": "114",
			"name": "绩溪县"
		},
		{
			"id": "1037",
			"city_id": "114",
			"name": "旌德县"
		},
		{
			"id": "1038",
			"city_id": "114",
			"name": "宁国市"
		},
		{
			"id": "1039",
			"city_id": "115",
			"name": "鼓楼区"
		},
		{
			"id": "1040",
			"city_id": "115",
			"name": "台江区"
		},
		{
			"id": "1041",
			"city_id": "115",
			"name": "仓山区"
		},
		{
			"id": "1042",
			"city_id": "115",
			"name": "马尾区"
		},
		{
			"id": "1043",
			"city_id": "115",
			"name": "晋安区"
		},
		{
			"id": "1044",
			"city_id": "115",
			"name": "闽侯县"
		},
		{
			"id": "1045",
			"city_id": "115",
			"name": "连江县"
		},
		{
			"id": "1046",
			"city_id": "115",
			"name": "罗源县"
		},
		{
			"id": "1047",
			"city_id": "115",
			"name": "闽清县"
		},
		{
			"id": "1048",
			"city_id": "115",
			"name": "永泰县"
		},
		{
			"id": "1049",
			"city_id": "115",
			"name": "平潭县"
		},
		{
			"id": "1050",
			"city_id": "115",
			"name": "福清市"
		},
		{
			"id": "1051",
			"city_id": "115",
			"name": "长乐市"
		},
		{
			"id": "1052",
			"city_id": "116",
			"name": "思明区"
		},
		{
			"id": "1053",
			"city_id": "116",
			"name": "海沧区"
		},
		{
			"id": "1054",
			"city_id": "116",
			"name": "湖里区"
		},
		{
			"id": "1055",
			"city_id": "116",
			"name": "集美区"
		},
		{
			"id": "1056",
			"city_id": "116",
			"name": "同安区"
		},
		{
			"id": "1057",
			"city_id": "116",
			"name": "翔安区"
		},
		{
			"id": "1058",
			"city_id": "117",
			"name": "城厢区"
		},
		{
			"id": "1059",
			"city_id": "117",
			"name": "涵江区"
		},
		{
			"id": "1060",
			"city_id": "117",
			"name": "荔城区"
		},
		{
			"id": "1061",
			"city_id": "117",
			"name": "秀屿区"
		},
		{
			"id": "1062",
			"city_id": "117",
			"name": "仙游县"
		},
		{
			"id": "1063",
			"city_id": "118",
			"name": "梅列区"
		},
		{
			"id": "1064",
			"city_id": "118",
			"name": "三元区"
		},
		{
			"id": "1065",
			"city_id": "118",
			"name": "明溪县"
		},
		{
			"id": "1066",
			"city_id": "118",
			"name": "清流县"
		},
		{
			"id": "1067",
			"city_id": "118",
			"name": "宁化县"
		},
		{
			"id": "1068",
			"city_id": "118",
			"name": "大田县"
		},
		{
			"id": "1069",
			"city_id": "118",
			"name": "尤溪县"
		},
		{
			"id": "1070",
			"city_id": "118",
			"name": "沙县"
		},
		{
			"id": "1071",
			"city_id": "118",
			"name": "将乐县"
		},
		{
			"id": "1072",
			"city_id": "118",
			"name": "泰宁县"
		},
		{
			"id": "1073",
			"city_id": "118",
			"name": "建宁县"
		},
		{
			"id": "1074",
			"city_id": "118",
			"name": "永安市"
		},
		{
			"id": "1075",
			"city_id": "119",
			"name": "鲤城区"
		},
		{
			"id": "1076",
			"city_id": "119",
			"name": "丰泽区"
		},
		{
			"id": "1077",
			"city_id": "119",
			"name": "洛江区"
		},
		{
			"id": "1078",
			"city_id": "119",
			"name": "泉港区"
		},
		{
			"id": "1079",
			"city_id": "119",
			"name": "惠安县"
		},
		{
			"id": "1080",
			"city_id": "119",
			"name": "安溪县"
		},
		{
			"id": "1081",
			"city_id": "119",
			"name": "永春县"
		},
		{
			"id": "1082",
			"city_id": "119",
			"name": "德化县"
		},
		{
			"id": "1083",
			"city_id": "119",
			"name": "金门县"
		},
		{
			"id": "1084",
			"city_id": "119",
			"name": "石狮市"
		},
		{
			"id": "1085",
			"city_id": "119",
			"name": "晋江市"
		},
		{
			"id": "1086",
			"city_id": "119",
			"name": "南安市"
		},
		{
			"id": "1087",
			"city_id": "120",
			"name": "芗城区"
		},
		{
			"id": "1088",
			"city_id": "120",
			"name": "龙文区"
		},
		{
			"id": "1089",
			"city_id": "120",
			"name": "云霄县"
		},
		{
			"id": "1090",
			"city_id": "120",
			"name": "漳浦县"
		},
		{
			"id": "1091",
			"city_id": "120",
			"name": "诏安县"
		},
		{
			"id": "1092",
			"city_id": "120",
			"name": "长泰县"
		},
		{
			"id": "1093",
			"city_id": "120",
			"name": "东山县"
		},
		{
			"id": "1094",
			"city_id": "120",
			"name": "南靖县"
		},
		{
			"id": "1095",
			"city_id": "120",
			"name": "平和县"
		},
		{
			"id": "1096",
			"city_id": "120",
			"name": "华安县"
		},
		{
			"id": "1097",
			"city_id": "120",
			"name": "龙海市"
		},
		{
			"id": "1098",
			"city_id": "121",
			"name": "延平区"
		},
		{
			"id": "1099",
			"city_id": "121",
			"name": "顺昌县"
		},
		{
			"id": "1100",
			"city_id": "121",
			"name": "浦城县"
		},
		{
			"id": "1101",
			"city_id": "121",
			"name": "光泽县"
		},
		{
			"id": "1102",
			"city_id": "121",
			"name": "松溪县"
		},
		{
			"id": "1103",
			"city_id": "121",
			"name": "政和县"
		},
		{
			"id": "1104",
			"city_id": "121",
			"name": "邵武市"
		},
		{
			"id": "1105",
			"city_id": "121",
			"name": "武夷山市"
		},
		{
			"id": "1106",
			"city_id": "121",
			"name": "建瓯市"
		},
		{
			"id": "1107",
			"city_id": "121",
			"name": "建阳市"
		},
		{
			"id": "1108",
			"city_id": "122",
			"name": "新罗区"
		},
		{
			"id": "1109",
			"city_id": "122",
			"name": "长汀县"
		},
		{
			"id": "1110",
			"city_id": "122",
			"name": "永定县"
		},
		{
			"id": "1111",
			"city_id": "122",
			"name": "上杭县"
		},
		{
			"id": "1112",
			"city_id": "122",
			"name": "武平县"
		},
		{
			"id": "1113",
			"city_id": "122",
			"name": "连城县"
		},
		{
			"id": "1114",
			"city_id": "122",
			"name": "漳平市"
		},
		{
			"id": "1115",
			"city_id": "123",
			"name": "蕉城区"
		},
		{
			"id": "1116",
			"city_id": "123",
			"name": "霞浦县"
		},
		{
			"id": "1117",
			"city_id": "123",
			"name": "古田县"
		},
		{
			"id": "1118",
			"city_id": "123",
			"name": "屏南县"
		},
		{
			"id": "1119",
			"city_id": "123",
			"name": "寿宁县"
		},
		{
			"id": "1120",
			"city_id": "123",
			"name": "周宁县"
		},
		{
			"id": "1121",
			"city_id": "123",
			"name": "柘荣县"
		},
		{
			"id": "1122",
			"city_id": "123",
			"name": "福安市"
		},
		{
			"id": "1123",
			"city_id": "123",
			"name": "福鼎市"
		},
		{
			"id": "1124",
			"city_id": "124",
			"name": "东湖区"
		},
		{
			"id": "1125",
			"city_id": "124",
			"name": "西湖区"
		},
		{
			"id": "1126",
			"city_id": "124",
			"name": "青云谱区"
		},
		{
			"id": "1127",
			"city_id": "124",
			"name": "湾里区"
		},
		{
			"id": "1128",
			"city_id": "124",
			"name": "青山湖区"
		},
		{
			"id": "1129",
			"city_id": "124",
			"name": "南昌县"
		},
		{
			"id": "1130",
			"city_id": "124",
			"name": "新建县"
		},
		{
			"id": "1131",
			"city_id": "124",
			"name": "安义县"
		},
		{
			"id": "1132",
			"city_id": "124",
			"name": "进贤县"
		},
		{
			"id": "1133",
			"city_id": "125",
			"name": "昌江区"
		},
		{
			"id": "1134",
			"city_id": "125",
			"name": "珠山区"
		},
		{
			"id": "1135",
			"city_id": "125",
			"name": "浮梁县"
		},
		{
			"id": "1136",
			"city_id": "125",
			"name": "乐平市"
		},
		{
			"id": "1137",
			"city_id": "126",
			"name": "安源区"
		},
		{
			"id": "1138",
			"city_id": "126",
			"name": "湘东区"
		},
		{
			"id": "1139",
			"city_id": "126",
			"name": "莲花县"
		},
		{
			"id": "1140",
			"city_id": "126",
			"name": "上栗县"
		},
		{
			"id": "1141",
			"city_id": "126",
			"name": "芦溪县"
		},
		{
			"id": "1142",
			"city_id": "127",
			"name": "庐山区"
		},
		{
			"id": "1143",
			"city_id": "127",
			"name": "浔阳区"
		},
		{
			"id": "1144",
			"city_id": "127",
			"name": "九江县"
		},
		{
			"id": "1145",
			"city_id": "127",
			"name": "武宁县"
		},
		{
			"id": "1146",
			"city_id": "127",
			"name": "修水县"
		},
		{
			"id": "1147",
			"city_id": "127",
			"name": "永修县"
		},
		{
			"id": "1148",
			"city_id": "127",
			"name": "德安县"
		},
		{
			"id": "1149",
			"city_id": "127",
			"name": "星子县"
		},
		{
			"id": "1150",
			"city_id": "127",
			"name": "都昌县"
		},
		{
			"id": "1151",
			"city_id": "127",
			"name": "湖口县"
		},
		{
			"id": "1152",
			"city_id": "127",
			"name": "彭泽县"
		},
		{
			"id": "1153",
			"city_id": "127",
			"name": "瑞昌市"
		},
		{
			"id": "1154",
			"city_id": "128",
			"name": "渝水区"
		},
		{
			"id": "1155",
			"city_id": "128",
			"name": "分宜县"
		},
		{
			"id": "1156",
			"city_id": "129",
			"name": "月湖区"
		},
		{
			"id": "1157",
			"city_id": "129",
			"name": "余江县"
		},
		{
			"id": "1158",
			"city_id": "129",
			"name": "贵溪市"
		},
		{
			"id": "1159",
			"city_id": "130",
			"name": "章贡区"
		},
		{
			"id": "1160",
			"city_id": "130",
			"name": "赣县"
		},
		{
			"id": "1161",
			"city_id": "130",
			"name": "信丰县"
		},
		{
			"id": "1162",
			"city_id": "130",
			"name": "大余县"
		},
		{
			"id": "1163",
			"city_id": "130",
			"name": "上犹县"
		},
		{
			"id": "1164",
			"city_id": "130",
			"name": "崇义县"
		},
		{
			"id": "1165",
			"city_id": "130",
			"name": "安远县"
		},
		{
			"id": "1166",
			"city_id": "130",
			"name": "龙南县"
		},
		{
			"id": "1167",
			"city_id": "130",
			"name": "定南县"
		},
		{
			"id": "1168",
			"city_id": "130",
			"name": "全南县"
		},
		{
			"id": "1169",
			"city_id": "130",
			"name": "宁都县"
		},
		{
			"id": "1170",
			"city_id": "130",
			"name": "于都县"
		},
		{
			"id": "1171",
			"city_id": "130",
			"name": "兴国县"
		},
		{
			"id": "1172",
			"city_id": "130",
			"name": "会昌县"
		},
		{
			"id": "1173",
			"city_id": "130",
			"name": "寻乌县"
		},
		{
			"id": "1174",
			"city_id": "130",
			"name": "石城县"
		},
		{
			"id": "1175",
			"city_id": "130",
			"name": "瑞金市"
		},
		{
			"id": "1176",
			"city_id": "130",
			"name": "南康市"
		},
		{
			"id": "1177",
			"city_id": "131",
			"name": "吉州区"
		},
		{
			"id": "1178",
			"city_id": "131",
			"name": "青原区"
		},
		{
			"id": "1179",
			"city_id": "131",
			"name": "吉安县"
		},
		{
			"id": "1180",
			"city_id": "131",
			"name": "吉水县"
		},
		{
			"id": "1181",
			"city_id": "131",
			"name": "峡江县"
		},
		{
			"id": "1182",
			"city_id": "131",
			"name": "新干县"
		},
		{
			"id": "1183",
			"city_id": "131",
			"name": "永丰县"
		},
		{
			"id": "1184",
			"city_id": "131",
			"name": "泰和县"
		},
		{
			"id": "1185",
			"city_id": "131",
			"name": "遂川县"
		},
		{
			"id": "1186",
			"city_id": "131",
			"name": "万安县"
		},
		{
			"id": "1187",
			"city_id": "131",
			"name": "安福县"
		},
		{
			"id": "1188",
			"city_id": "131",
			"name": "永新县"
		},
		{
			"id": "1189",
			"city_id": "131",
			"name": "井冈山市"
		},
		{
			"id": "1190",
			"city_id": "132",
			"name": "袁州区"
		},
		{
			"id": "1191",
			"city_id": "132",
			"name": "奉新县"
		},
		{
			"id": "1192",
			"city_id": "132",
			"name": "万载县"
		},
		{
			"id": "1193",
			"city_id": "132",
			"name": "上高县"
		},
		{
			"id": "1194",
			"city_id": "132",
			"name": "宜丰县"
		},
		{
			"id": "1195",
			"city_id": "132",
			"name": "靖安县"
		},
		{
			"id": "1196",
			"city_id": "132",
			"name": "铜鼓县"
		},
		{
			"id": "1197",
			"city_id": "132",
			"name": "丰城市"
		},
		{
			"id": "1198",
			"city_id": "132",
			"name": "樟树市"
		},
		{
			"id": "1199",
			"city_id": "132",
			"name": "高安市"
		},
		{
			"id": "1200",
			"city_id": "133",
			"name": "临川区"
		},
		{
			"id": "1201",
			"city_id": "133",
			"name": "南城县"
		},
		{
			"id": "1202",
			"city_id": "133",
			"name": "黎川县"
		},
		{
			"id": "1203",
			"city_id": "133",
			"name": "南丰县"
		},
		{
			"id": "1204",
			"city_id": "133",
			"name": "崇仁县"
		},
		{
			"id": "1205",
			"city_id": "133",
			"name": "乐安县"
		},
		{
			"id": "1206",
			"city_id": "133",
			"name": "宜黄县"
		},
		{
			"id": "1207",
			"city_id": "133",
			"name": "金溪县"
		},
		{
			"id": "1208",
			"city_id": "133",
			"name": "资溪县"
		},
		{
			"id": "1209",
			"city_id": "133",
			"name": "东乡县"
		},
		{
			"id": "1210",
			"city_id": "133",
			"name": "广昌县"
		},
		{
			"id": "1211",
			"city_id": "134",
			"name": "信州区"
		},
		{
			"id": "1212",
			"city_id": "134",
			"name": "上饶县"
		},
		{
			"id": "1213",
			"city_id": "134",
			"name": "广丰县"
		},
		{
			"id": "1214",
			"city_id": "134",
			"name": "玉山县"
		},
		{
			"id": "1215",
			"city_id": "134",
			"name": "铅山县"
		},
		{
			"id": "1216",
			"city_id": "134",
			"name": "横峰县"
		},
		{
			"id": "1217",
			"city_id": "134",
			"name": "弋阳县"
		},
		{
			"id": "1218",
			"city_id": "134",
			"name": "余干县"
		},
		{
			"id": "1219",
			"city_id": "134",
			"name": "鄱阳县"
		},
		{
			"id": "1220",
			"city_id": "134",
			"name": "万年县"
		},
		{
			"id": "1221",
			"city_id": "134",
			"name": "婺源县"
		},
		{
			"id": "1222",
			"city_id": "134",
			"name": "德兴市"
		},
		{
			"id": "1223",
			"city_id": "135",
			"name": "历下区"
		},
		{
			"id": "1224",
			"city_id": "135",
			"name": "市中区"
		},
		{
			"id": "1225",
			"city_id": "135",
			"name": "槐荫区"
		},
		{
			"id": "1226",
			"city_id": "135",
			"name": "天桥区"
		},
		{
			"id": "1227",
			"city_id": "135",
			"name": "历城区"
		},
		{
			"id": "1228",
			"city_id": "135",
			"name": "长清区"
		},
		{
			"id": "1229",
			"city_id": "135",
			"name": "平阴县"
		},
		{
			"id": "1230",
			"city_id": "135",
			"name": "济阳县"
		},
		{
			"id": "1231",
			"city_id": "135",
			"name": "商河县"
		},
		{
			"id": "1232",
			"city_id": "135",
			"name": "章丘市"
		},
		{
			"id": "1233",
			"city_id": "136",
			"name": "市南区"
		},
		{
			"id": "1234",
			"city_id": "136",
			"name": "市北区"
		},
		{
			"id": "1235",
			"city_id": "136",
			"name": "四方区"
		},
		{
			"id": "1236",
			"city_id": "136",
			"name": "黄岛区"
		},
		{
			"id": "1237",
			"city_id": "136",
			"name": "崂山区"
		},
		{
			"id": "1238",
			"city_id": "136",
			"name": "李沧区"
		},
		{
			"id": "1239",
			"city_id": "136",
			"name": "城阳区"
		},
		{
			"id": "1240",
			"city_id": "136",
			"name": "胶州市"
		},
		{
			"id": "1241",
			"city_id": "136",
			"name": "即墨市"
		},
		{
			"id": "1242",
			"city_id": "136",
			"name": "平度市"
		},
		{
			"id": "1243",
			"city_id": "136",
			"name": "胶南市"
		},
		{
			"id": "1244",
			"city_id": "136",
			"name": "莱西市"
		},
		{
			"id": "1245",
			"city_id": "137",
			"name": "淄川区"
		},
		{
			"id": "1246",
			"city_id": "137",
			"name": "张店区"
		},
		{
			"id": "1247",
			"city_id": "137",
			"name": "博山区"
		},
		{
			"id": "1248",
			"city_id": "137",
			"name": "临淄区"
		},
		{
			"id": "1249",
			"city_id": "137",
			"name": "周村区"
		},
		{
			"id": "1250",
			"city_id": "137",
			"name": "桓台县"
		},
		{
			"id": "1251",
			"city_id": "137",
			"name": "高青县"
		},
		{
			"id": "1252",
			"city_id": "137",
			"name": "沂源县"
		},
		{
			"id": "1253",
			"city_id": "138",
			"name": "市中区"
		},
		{
			"id": "1254",
			"city_id": "138",
			"name": "薛城区"
		},
		{
			"id": "1255",
			"city_id": "138",
			"name": "峄城区"
		},
		{
			"id": "1256",
			"city_id": "138",
			"name": "台儿庄区"
		},
		{
			"id": "1257",
			"city_id": "138",
			"name": "山亭区"
		},
		{
			"id": "1258",
			"city_id": "138",
			"name": "滕州市"
		},
		{
			"id": "1259",
			"city_id": "139",
			"name": "东营区"
		},
		{
			"id": "1260",
			"city_id": "139",
			"name": "河口区"
		},
		{
			"id": "1261",
			"city_id": "139",
			"name": "垦利县"
		},
		{
			"id": "1262",
			"city_id": "139",
			"name": "利津县"
		},
		{
			"id": "1263",
			"city_id": "139",
			"name": "广饶县"
		},
		{
			"id": "1264",
			"city_id": "140",
			"name": "芝罘区"
		},
		{
			"id": "1265",
			"city_id": "140",
			"name": "福山区"
		},
		{
			"id": "1266",
			"city_id": "140",
			"name": "牟平区"
		},
		{
			"id": "1267",
			"city_id": "140",
			"name": "莱山区"
		},
		{
			"id": "1268",
			"city_id": "140",
			"name": "长岛县"
		},
		{
			"id": "1269",
			"city_id": "140",
			"name": "龙口市"
		},
		{
			"id": "1270",
			"city_id": "140",
			"name": "莱阳市"
		},
		{
			"id": "1271",
			"city_id": "140",
			"name": "莱州市"
		},
		{
			"id": "1272",
			"city_id": "140",
			"name": "蓬莱市"
		},
		{
			"id": "1273",
			"city_id": "140",
			"name": "招远市"
		},
		{
			"id": "1274",
			"city_id": "140",
			"name": "栖霞市"
		},
		{
			"id": "1275",
			"city_id": "140",
			"name": "海阳市"
		},
		{
			"id": "1276",
			"city_id": "141",
			"name": "潍城区"
		},
		{
			"id": "1277",
			"city_id": "141",
			"name": "寒亭区"
		},
		{
			"id": "1278",
			"city_id": "141",
			"name": "坊子区"
		},
		{
			"id": "1279",
			"city_id": "141",
			"name": "奎文区"
		},
		{
			"id": "1280",
			"city_id": "141",
			"name": "临朐县"
		},
		{
			"id": "1281",
			"city_id": "141",
			"name": "昌乐县"
		},
		{
			"id": "1282",
			"city_id": "141",
			"name": "青州市"
		},
		{
			"id": "1283",
			"city_id": "141",
			"name": "诸城市"
		},
		{
			"id": "1284",
			"city_id": "141",
			"name": "寿光市"
		},
		{
			"id": "1285",
			"city_id": "141",
			"name": "安丘市"
		},
		{
			"id": "1286",
			"city_id": "141",
			"name": "高密市"
		},
		{
			"id": "1287",
			"city_id": "141",
			"name": "昌邑市"
		},
		{
			"id": "1288",
			"city_id": "142",
			"name": "市中区"
		},
		{
			"id": "1289",
			"city_id": "142",
			"name": "任城区"
		},
		{
			"id": "1290",
			"city_id": "142",
			"name": "微山县"
		},
		{
			"id": "1291",
			"city_id": "142",
			"name": "鱼台县"
		},
		{
			"id": "1292",
			"city_id": "142",
			"name": "金乡县"
		},
		{
			"id": "1293",
			"city_id": "142",
			"name": "嘉祥县"
		},
		{
			"id": "1294",
			"city_id": "142",
			"name": "汶上县"
		},
		{
			"id": "1295",
			"city_id": "142",
			"name": "泗水县"
		},
		{
			"id": "1296",
			"city_id": "142",
			"name": "梁山县"
		},
		{
			"id": "1297",
			"city_id": "142",
			"name": "曲阜市"
		},
		{
			"id": "1298",
			"city_id": "142",
			"name": "兖州市"
		},
		{
			"id": "1299",
			"city_id": "142",
			"name": "邹城市"
		},
		{
			"id": "1300",
			"city_id": "143",
			"name": "泰山区"
		},
		{
			"id": "1301",
			"city_id": "143",
			"name": "岱岳区"
		},
		{
			"id": "1302",
			"city_id": "143",
			"name": "宁阳县"
		},
		{
			"id": "1303",
			"city_id": "143",
			"name": "东平县"
		},
		{
			"id": "1304",
			"city_id": "143",
			"name": "新泰市"
		},
		{
			"id": "1305",
			"city_id": "143",
			"name": "肥城市"
		},
		{
			"id": "1306",
			"city_id": "144",
			"name": "环翠区"
		},
		{
			"id": "1307",
			"city_id": "144",
			"name": "文登市"
		},
		{
			"id": "1308",
			"city_id": "144",
			"name": "荣成市"
		},
		{
			"id": "1309",
			"city_id": "144",
			"name": "乳山市"
		},
		{
			"id": "1310",
			"city_id": "145",
			"name": "东港区"
		},
		{
			"id": "1311",
			"city_id": "145",
			"name": "岚山区"
		},
		{
			"id": "1312",
			"city_id": "145",
			"name": "五莲县"
		},
		{
			"id": "1313",
			"city_id": "145",
			"name": "莒县"
		},
		{
			"id": "1314",
			"city_id": "146",
			"name": "莱城区"
		},
		{
			"id": "1315",
			"city_id": "146",
			"name": "钢城区"
		},
		{
			"id": "1316",
			"city_id": "147",
			"name": "兰山区"
		},
		{
			"id": "1317",
			"city_id": "147",
			"name": "罗庄区"
		},
		{
			"id": "1318",
			"city_id": "147",
			"name": "河东区"
		},
		{
			"id": "1319",
			"city_id": "147",
			"name": "沂南县"
		},
		{
			"id": "1320",
			"city_id": "147",
			"name": "郯城县"
		},
		{
			"id": "1321",
			"city_id": "147",
			"name": "沂水县"
		},
		{
			"id": "1322",
			"city_id": "147",
			"name": "苍山县"
		},
		{
			"id": "1323",
			"city_id": "147",
			"name": "费县"
		},
		{
			"id": "1324",
			"city_id": "147",
			"name": "平邑县"
		},
		{
			"id": "1325",
			"city_id": "147",
			"name": "莒南县"
		},
		{
			"id": "1326",
			"city_id": "147",
			"name": "蒙阴县"
		},
		{
			"id": "1327",
			"city_id": "147",
			"name": "临沭县"
		},
		{
			"id": "1328",
			"city_id": "148",
			"name": "德城区"
		},
		{
			"id": "1329",
			"city_id": "148",
			"name": "陵县"
		},
		{
			"id": "1330",
			"city_id": "148",
			"name": "宁津县"
		},
		{
			"id": "1331",
			"city_id": "148",
			"name": "庆云县"
		},
		{
			"id": "1332",
			"city_id": "148",
			"name": "临邑县"
		},
		{
			"id": "1333",
			"city_id": "148",
			"name": "齐河县"
		},
		{
			"id": "1334",
			"city_id": "148",
			"name": "平原县"
		},
		{
			"id": "1335",
			"city_id": "148",
			"name": "夏津县"
		},
		{
			"id": "1336",
			"city_id": "148",
			"name": "武城县"
		},
		{
			"id": "1337",
			"city_id": "148",
			"name": "乐陵市"
		},
		{
			"id": "1338",
			"city_id": "148",
			"name": "禹城市"
		},
		{
			"id": "1339",
			"city_id": "149",
			"name": "东昌府区"
		},
		{
			"id": "1340",
			"city_id": "149",
			"name": "阳谷县"
		},
		{
			"id": "1341",
			"city_id": "149",
			"name": "莘县"
		},
		{
			"id": "1342",
			"city_id": "149",
			"name": "茌平县"
		},
		{
			"id": "1343",
			"city_id": "149",
			"name": "东阿县"
		},
		{
			"id": "1344",
			"city_id": "149",
			"name": "冠县"
		},
		{
			"id": "1345",
			"city_id": "149",
			"name": "高唐县"
		},
		{
			"id": "1346",
			"city_id": "149",
			"name": "临清市"
		},
		{
			"id": "1347",
			"city_id": "150",
			"name": "滨城区"
		},
		{
			"id": "1348",
			"city_id": "150",
			"name": "惠民县"
		},
		{
			"id": "1349",
			"city_id": "150",
			"name": "阳信县"
		},
		{
			"id": "1350",
			"city_id": "150",
			"name": "无棣县"
		},
		{
			"id": "1351",
			"city_id": "150",
			"name": "沾化县"
		},
		{
			"id": "1352",
			"city_id": "150",
			"name": "博兴县"
		},
		{
			"id": "1353",
			"city_id": "150",
			"name": "邹平县"
		},
		{
			"id": "1354",
			"city_id": "151",
			"name": "牡丹区"
		},
		{
			"id": "1355",
			"city_id": "151",
			"name": "曹县"
		},
		{
			"id": "1356",
			"city_id": "151",
			"name": "单县"
		},
		{
			"id": "1357",
			"city_id": "151",
			"name": "成武县"
		},
		{
			"id": "1358",
			"city_id": "151",
			"name": "巨野县"
		},
		{
			"id": "1359",
			"city_id": "151",
			"name": "郓城县"
		},
		{
			"id": "1360",
			"city_id": "151",
			"name": "鄄城县"
		},
		{
			"id": "1361",
			"city_id": "151",
			"name": "定陶县"
		},
		{
			"id": "1362",
			"city_id": "151",
			"name": "东明县"
		},
		{
			"id": "1363",
			"city_id": "152",
			"name": "中原区"
		},
		{
			"id": "1364",
			"city_id": "152",
			"name": "二七区"
		},
		{
			"id": "1365",
			"city_id": "152",
			"name": "管城回族区"
		},
		{
			"id": "1366",
			"city_id": "152",
			"name": "金水区"
		},
		{
			"id": "1367",
			"city_id": "152",
			"name": "上街区"
		},
		{
			"id": "1368",
			"city_id": "152",
			"name": "惠济区"
		},
		{
			"id": "1369",
			"city_id": "152",
			"name": "中牟县"
		},
		{
			"id": "1370",
			"city_id": "152",
			"name": "巩义市"
		},
		{
			"id": "1371",
			"city_id": "152",
			"name": "荥阳市"
		},
		{
			"id": "1372",
			"city_id": "152",
			"name": "新密市"
		},
		{
			"id": "1373",
			"city_id": "152",
			"name": "新郑市"
		},
		{
			"id": "1374",
			"city_id": "152",
			"name": "登封市"
		},
		{
			"id": "1375",
			"city_id": "153",
			"name": "龙亭区"
		},
		{
			"id": "1376",
			"city_id": "153",
			"name": "顺河回族区"
		},
		{
			"id": "1377",
			"city_id": "153",
			"name": "鼓楼区"
		},
		{
			"id": "1378",
			"city_id": "153",
			"name": "南关区"
		},
		{
			"id": "1379",
			"city_id": "153",
			"name": "郊区"
		},
		{
			"id": "1380",
			"city_id": "153",
			"name": "杞县"
		},
		{
			"id": "1381",
			"city_id": "153",
			"name": "通许县"
		},
		{
			"id": "1382",
			"city_id": "153",
			"name": "尉氏县"
		},
		{
			"id": "1383",
			"city_id": "153",
			"name": "开封县"
		},
		{
			"id": "1384",
			"city_id": "153",
			"name": "兰考县"
		},
		{
			"id": "1385",
			"city_id": "154",
			"name": "老城区"
		},
		{
			"id": "1386",
			"city_id": "154",
			"name": "西工区"
		},
		{
			"id": "1387",
			"city_id": "154",
			"name": "廛河回族区"
		},
		{
			"id": "1388",
			"city_id": "154",
			"name": "涧西区"
		},
		{
			"id": "1389",
			"city_id": "154",
			"name": "吉利区"
		},
		{
			"id": "1390",
			"city_id": "154",
			"name": "洛龙区"
		},
		{
			"id": "1391",
			"city_id": "154",
			"name": "孟津县"
		},
		{
			"id": "1392",
			"city_id": "154",
			"name": "新安县"
		},
		{
			"id": "1393",
			"city_id": "154",
			"name": "栾川县"
		},
		{
			"id": "1394",
			"city_id": "154",
			"name": "嵩县"
		},
		{
			"id": "1395",
			"city_id": "154",
			"name": "汝阳县"
		},
		{
			"id": "1396",
			"city_id": "154",
			"name": "宜阳县"
		},
		{
			"id": "1397",
			"city_id": "154",
			"name": "洛宁县"
		},
		{
			"id": "1398",
			"city_id": "154",
			"name": "伊川县"
		},
		{
			"id": "1399",
			"city_id": "154",
			"name": "偃师市"
		},
		{
			"id": "1400",
			"city_id": "155",
			"name": "新华区"
		},
		{
			"id": "1401",
			"city_id": "155",
			"name": "卫东区"
		},
		{
			"id": "1402",
			"city_id": "155",
			"name": "石龙区"
		},
		{
			"id": "1403",
			"city_id": "155",
			"name": "湛河区"
		},
		{
			"id": "1404",
			"city_id": "155",
			"name": "宝丰县"
		},
		{
			"id": "1405",
			"city_id": "155",
			"name": "叶县"
		},
		{
			"id": "1406",
			"city_id": "155",
			"name": "鲁山县"
		},
		{
			"id": "1407",
			"city_id": "155",
			"name": "郏县"
		},
		{
			"id": "1408",
			"city_id": "155",
			"name": "舞钢市"
		},
		{
			"id": "1409",
			"city_id": "155",
			"name": "汝州市"
		},
		{
			"id": "1410",
			"city_id": "156",
			"name": "文峰区"
		},
		{
			"id": "1411",
			"city_id": "156",
			"name": "北关区"
		},
		{
			"id": "1412",
			"city_id": "156",
			"name": "殷都区"
		},
		{
			"id": "1413",
			"city_id": "156",
			"name": "龙安区"
		},
		{
			"id": "1414",
			"city_id": "156",
			"name": "安阳县"
		},
		{
			"id": "1415",
			"city_id": "156",
			"name": "汤阴县"
		},
		{
			"id": "1416",
			"city_id": "156",
			"name": "滑县"
		},
		{
			"id": "1417",
			"city_id": "156",
			"name": "内黄县"
		},
		{
			"id": "1418",
			"city_id": "156",
			"name": "林州市"
		},
		{
			"id": "1419",
			"city_id": "157",
			"name": "鹤山区"
		},
		{
			"id": "1420",
			"city_id": "157",
			"name": "山城区"
		},
		{
			"id": "1421",
			"city_id": "157",
			"name": "淇滨区"
		},
		{
			"id": "1422",
			"city_id": "157",
			"name": "浚县"
		},
		{
			"id": "1423",
			"city_id": "157",
			"name": "淇县"
		},
		{
			"id": "1424",
			"city_id": "158",
			"name": "红旗区"
		},
		{
			"id": "1425",
			"city_id": "158",
			"name": "卫滨区"
		},
		{
			"id": "1426",
			"city_id": "158",
			"name": "凤泉区"
		},
		{
			"id": "1427",
			"city_id": "158",
			"name": "牧野区"
		},
		{
			"id": "1428",
			"city_id": "158",
			"name": "新乡县"
		},
		{
			"id": "1429",
			"city_id": "158",
			"name": "获嘉县"
		},
		{
			"id": "1430",
			"city_id": "158",
			"name": "原阳县"
		},
		{
			"id": "1431",
			"city_id": "158",
			"name": "延津县"
		},
		{
			"id": "1432",
			"city_id": "158",
			"name": "封丘县"
		},
		{
			"id": "1433",
			"city_id": "158",
			"name": "长垣县"
		},
		{
			"id": "1434",
			"city_id": "158",
			"name": "卫辉市"
		},
		{
			"id": "1435",
			"city_id": "158",
			"name": "辉县市"
		},
		{
			"id": "1436",
			"city_id": "159",
			"name": "解放区"
		},
		{
			"id": "1437",
			"city_id": "159",
			"name": "中站区"
		},
		{
			"id": "1438",
			"city_id": "159",
			"name": "马村区"
		},
		{
			"id": "1439",
			"city_id": "159",
			"name": "山阳区"
		},
		{
			"id": "1440",
			"city_id": "159",
			"name": "修武县"
		},
		{
			"id": "1441",
			"city_id": "159",
			"name": "博爱县"
		},
		{
			"id": "1442",
			"city_id": "159",
			"name": "武陟县"
		},
		{
			"id": "1443",
			"city_id": "159",
			"name": "温县"
		},
		{
			"id": "1444",
			"city_id": "159",
			"name": "济源市"
		},
		{
			"id": "1445",
			"city_id": "159",
			"name": "沁阳市"
		},
		{
			"id": "1446",
			"city_id": "159",
			"name": "孟州市"
		},
		{
			"id": "1447",
			"city_id": "160",
			"name": "华龙区"
		},
		{
			"id": "1448",
			"city_id": "160",
			"name": "清丰县"
		},
		{
			"id": "1449",
			"city_id": "160",
			"name": "南乐县"
		},
		{
			"id": "1450",
			"city_id": "160",
			"name": "范县"
		},
		{
			"id": "1451",
			"city_id": "160",
			"name": "台前县"
		},
		{
			"id": "1452",
			"city_id": "160",
			"name": "濮阳县"
		},
		{
			"id": "1453",
			"city_id": "161",
			"name": "魏都区"
		},
		{
			"id": "1454",
			"city_id": "161",
			"name": "许昌县"
		},
		{
			"id": "1455",
			"city_id": "161",
			"name": "鄢陵县"
		},
		{
			"id": "1456",
			"city_id": "161",
			"name": "襄城县"
		},
		{
			"id": "1457",
			"city_id": "161",
			"name": "禹州市"
		},
		{
			"id": "1458",
			"city_id": "161",
			"name": "长葛市"
		},
		{
			"id": "1459",
			"city_id": "162",
			"name": "源汇区"
		},
		{
			"id": "1460",
			"city_id": "162",
			"name": "郾城区"
		},
		{
			"id": "1461",
			"city_id": "162",
			"name": "召陵区"
		},
		{
			"id": "1462",
			"city_id": "162",
			"name": "舞阳县"
		},
		{
			"id": "1463",
			"city_id": "162",
			"name": "临颍县"
		},
		{
			"id": "1464",
			"city_id": "163",
			"name": "市辖区"
		},
		{
			"id": "1465",
			"city_id": "163",
			"name": "湖滨区"
		},
		{
			"id": "1466",
			"city_id": "163",
			"name": "渑池县"
		},
		{
			"id": "1467",
			"city_id": "163",
			"name": "陕县"
		},
		{
			"id": "1468",
			"city_id": "163",
			"name": "卢氏县"
		},
		{
			"id": "1469",
			"city_id": "163",
			"name": "义马市"
		},
		{
			"id": "1470",
			"city_id": "163",
			"name": "灵宝市"
		},
		{
			"id": "1471",
			"city_id": "164",
			"name": "宛城区"
		},
		{
			"id": "1472",
			"city_id": "164",
			"name": "卧龙区"
		},
		{
			"id": "1473",
			"city_id": "164",
			"name": "南召县"
		},
		{
			"id": "1474",
			"city_id": "164",
			"name": "方城县"
		},
		{
			"id": "1475",
			"city_id": "164",
			"name": "西峡县"
		},
		{
			"id": "1476",
			"city_id": "164",
			"name": "镇平县"
		},
		{
			"id": "1477",
			"city_id": "164",
			"name": "内乡县"
		},
		{
			"id": "1478",
			"city_id": "164",
			"name": "淅川县"
		},
		{
			"id": "1479",
			"city_id": "164",
			"name": "社旗县"
		},
		{
			"id": "1480",
			"city_id": "164",
			"name": "唐河县"
		},
		{
			"id": "1481",
			"city_id": "164",
			"name": "新野县"
		},
		{
			"id": "1482",
			"city_id": "164",
			"name": "桐柏县"
		},
		{
			"id": "1483",
			"city_id": "164",
			"name": "邓州市"
		},
		{
			"id": "1484",
			"city_id": "165",
			"name": "梁园区"
		},
		{
			"id": "1485",
			"city_id": "165",
			"name": "睢阳区"
		},
		{
			"id": "1486",
			"city_id": "165",
			"name": "民权县"
		},
		{
			"id": "1487",
			"city_id": "165",
			"name": "睢县"
		},
		{
			"id": "1488",
			"city_id": "165",
			"name": "宁陵县"
		},
		{
			"id": "1489",
			"city_id": "165",
			"name": "柘城县"
		},
		{
			"id": "1490",
			"city_id": "165",
			"name": "虞城县"
		},
		{
			"id": "1491",
			"city_id": "165",
			"name": "夏邑县"
		},
		{
			"id": "1492",
			"city_id": "165",
			"name": "永城市"
		},
		{
			"id": "1493",
			"city_id": "166",
			"name": "浉河区"
		},
		{
			"id": "1494",
			"city_id": "166",
			"name": "平桥区"
		},
		{
			"id": "1495",
			"city_id": "166",
			"name": "罗山县"
		},
		{
			"id": "1496",
			"city_id": "166",
			"name": "光山县"
		},
		{
			"id": "1497",
			"city_id": "166",
			"name": "新县"
		},
		{
			"id": "1498",
			"city_id": "166",
			"name": "商城县"
		},
		{
			"id": "1499",
			"city_id": "166",
			"name": "固始县"
		},
		{
			"id": "1500",
			"city_id": "166",
			"name": "潢川县"
		},
		{
			"id": "1501",
			"city_id": "166",
			"name": "淮滨县"
		},
		{
			"id": "1502",
			"city_id": "166",
			"name": "息县"
		},
		{
			"id": "1503",
			"city_id": "167",
			"name": "川汇区"
		},
		{
			"id": "1504",
			"city_id": "167",
			"name": "扶沟县"
		},
		{
			"id": "1505",
			"city_id": "167",
			"name": "西华县"
		},
		{
			"id": "1506",
			"city_id": "167",
			"name": "商水县"
		},
		{
			"id": "1507",
			"city_id": "167",
			"name": "沈丘县"
		},
		{
			"id": "1508",
			"city_id": "167",
			"name": "郸城县"
		},
		{
			"id": "1509",
			"city_id": "167",
			"name": "淮阳县"
		},
		{
			"id": "1510",
			"city_id": "167",
			"name": "太康县"
		},
		{
			"id": "1511",
			"city_id": "167",
			"name": "鹿邑县"
		},
		{
			"id": "1512",
			"city_id": "167",
			"name": "项城市"
		},
		{
			"id": "1513",
			"city_id": "168",
			"name": "驿城区"
		},
		{
			"id": "1514",
			"city_id": "168",
			"name": "西平县"
		},
		{
			"id": "1515",
			"city_id": "168",
			"name": "上蔡县"
		},
		{
			"id": "1516",
			"city_id": "168",
			"name": "平舆县"
		},
		{
			"id": "1517",
			"city_id": "168",
			"name": "正阳县"
		},
		{
			"id": "1518",
			"city_id": "168",
			"name": "确山县"
		},
		{
			"id": "1519",
			"city_id": "168",
			"name": "泌阳县"
		},
		{
			"id": "1520",
			"city_id": "168",
			"name": "汝南县"
		},
		{
			"id": "1521",
			"city_id": "168",
			"name": "遂平县"
		},
		{
			"id": "1522",
			"city_id": "168",
			"name": "新蔡县"
		},
		{
			"id": "1523",
			"city_id": "169",
			"name": "江岸区"
		},
		{
			"id": "1524",
			"city_id": "169",
			"name": "江汉区"
		},
		{
			"id": "1525",
			"city_id": "169",
			"name": "硚口区"
		},
		{
			"id": "1526",
			"city_id": "169",
			"name": "汉阳区"
		},
		{
			"id": "1527",
			"city_id": "169",
			"name": "武昌区"
		},
		{
			"id": "1528",
			"city_id": "169",
			"name": "青山区"
		},
		{
			"id": "1529",
			"city_id": "169",
			"name": "洪山区"
		},
		{
			"id": "1530",
			"city_id": "169",
			"name": "东西湖区"
		},
		{
			"id": "1531",
			"city_id": "169",
			"name": "汉南区"
		},
		{
			"id": "1532",
			"city_id": "169",
			"name": "蔡甸区"
		},
		{
			"id": "1533",
			"city_id": "169",
			"name": "江夏区"
		},
		{
			"id": "1534",
			"city_id": "169",
			"name": "黄陂区"
		},
		{
			"id": "1535",
			"city_id": "169",
			"name": "新洲区"
		},
		{
			"id": "1536",
			"city_id": "170",
			"name": "黄石港区"
		},
		{
			"id": "1537",
			"city_id": "170",
			"name": "西塞山区"
		},
		{
			"id": "1538",
			"city_id": "170",
			"name": "下陆区"
		},
		{
			"id": "1539",
			"city_id": "170",
			"name": "铁山区"
		},
		{
			"id": "1540",
			"city_id": "170",
			"name": "阳新县"
		},
		{
			"id": "1541",
			"city_id": "170",
			"name": "大冶市"
		},
		{
			"id": "1542",
			"city_id": "171",
			"name": "茅箭区"
		},
		{
			"id": "1543",
			"city_id": "171",
			"name": "张湾区"
		},
		{
			"id": "1544",
			"city_id": "171",
			"name": "郧县"
		},
		{
			"id": "1545",
			"city_id": "171",
			"name": "郧西县"
		},
		{
			"id": "1546",
			"city_id": "171",
			"name": "竹山县"
		},
		{
			"id": "1547",
			"city_id": "171",
			"name": "竹溪县"
		},
		{
			"id": "1548",
			"city_id": "171",
			"name": "房县"
		},
		{
			"id": "1549",
			"city_id": "171",
			"name": "丹江口市"
		},
		{
			"id": "1550",
			"city_id": "172",
			"name": "西陵区"
		},
		{
			"id": "1551",
			"city_id": "172",
			"name": "伍家岗区"
		},
		{
			"id": "1552",
			"city_id": "172",
			"name": "点军区"
		},
		{
			"id": "1553",
			"city_id": "172",
			"name": "猇亭区"
		},
		{
			"id": "1554",
			"city_id": "172",
			"name": "夷陵区"
		},
		{
			"id": "1555",
			"city_id": "172",
			"name": "远安县"
		},
		{
			"id": "1556",
			"city_id": "172",
			"name": "兴山县"
		},
		{
			"id": "1557",
			"city_id": "172",
			"name": "秭归县"
		},
		{
			"id": "1558",
			"city_id": "172",
			"name": "长阳土家族自治县"
		},
		{
			"id": "1559",
			"city_id": "172",
			"name": "五峰土家族自治县"
		},
		{
			"id": "1560",
			"city_id": "172",
			"name": "宜都市"
		},
		{
			"id": "1561",
			"city_id": "172",
			"name": "当阳市"
		},
		{
			"id": "1562",
			"city_id": "172",
			"name": "枝江市"
		},
		{
			"id": "1563",
			"city_id": "173",
			"name": "襄城区"
		},
		{
			"id": "1564",
			"city_id": "173",
			"name": "樊城区"
		},
		{
			"id": "1565",
			"city_id": "173",
			"name": "襄阳区"
		},
		{
			"id": "1566",
			"city_id": "173",
			"name": "南漳县"
		},
		{
			"id": "1567",
			"city_id": "173",
			"name": "谷城县"
		},
		{
			"id": "1568",
			"city_id": "173",
			"name": "保康县"
		},
		{
			"id": "1569",
			"city_id": "173",
			"name": "老河口市"
		},
		{
			"id": "1570",
			"city_id": "173",
			"name": "枣阳市"
		},
		{
			"id": "1571",
			"city_id": "173",
			"name": "宜城市"
		},
		{
			"id": "1572",
			"city_id": "174",
			"name": "梁子湖区"
		},
		{
			"id": "1573",
			"city_id": "174",
			"name": "华容区"
		},
		{
			"id": "1574",
			"city_id": "174",
			"name": "鄂城区"
		},
		{
			"id": "1575",
			"city_id": "175",
			"name": "东宝区"
		},
		{
			"id": "1576",
			"city_id": "175",
			"name": "掇刀区"
		},
		{
			"id": "1577",
			"city_id": "175",
			"name": "京山县"
		},
		{
			"id": "1578",
			"city_id": "175",
			"name": "沙洋县"
		},
		{
			"id": "1579",
			"city_id": "175",
			"name": "钟祥市"
		},
		{
			"id": "1580",
			"city_id": "176",
			"name": "孝南区"
		},
		{
			"id": "1581",
			"city_id": "176",
			"name": "孝昌县"
		},
		{
			"id": "1582",
			"city_id": "176",
			"name": "大悟县"
		},
		{
			"id": "1583",
			"city_id": "176",
			"name": "云梦县"
		},
		{
			"id": "1584",
			"city_id": "176",
			"name": "应城市"
		},
		{
			"id": "1585",
			"city_id": "176",
			"name": "安陆市"
		},
		{
			"id": "1586",
			"city_id": "176",
			"name": "汉川市"
		},
		{
			"id": "1587",
			"city_id": "177",
			"name": "沙市区"
		},
		{
			"id": "1588",
			"city_id": "177",
			"name": "荆州区"
		},
		{
			"id": "1589",
			"city_id": "177",
			"name": "公安县"
		},
		{
			"id": "1590",
			"city_id": "177",
			"name": "监利县"
		},
		{
			"id": "1591",
			"city_id": "177",
			"name": "江陵县"
		},
		{
			"id": "1592",
			"city_id": "177",
			"name": "石首市"
		},
		{
			"id": "1593",
			"city_id": "177",
			"name": "洪湖市"
		},
		{
			"id": "1594",
			"city_id": "177",
			"name": "松滋市"
		},
		{
			"id": "1595",
			"city_id": "178",
			"name": "黄州区"
		},
		{
			"id": "1596",
			"city_id": "178",
			"name": "团风县"
		},
		{
			"id": "1597",
			"city_id": "178",
			"name": "红安县"
		},
		{
			"id": "1598",
			"city_id": "178",
			"name": "罗田县"
		},
		{
			"id": "1599",
			"city_id": "178",
			"name": "英山县"
		},
		{
			"id": "1600",
			"city_id": "178",
			"name": "浠水县"
		},
		{
			"id": "1601",
			"city_id": "178",
			"name": "蕲春县"
		},
		{
			"id": "1602",
			"city_id": "178",
			"name": "黄梅县"
		},
		{
			"id": "1603",
			"city_id": "178",
			"name": "麻城市"
		},
		{
			"id": "1604",
			"city_id": "178",
			"name": "武穴市"
		},
		{
			"id": "1605",
			"city_id": "179",
			"name": "咸安区"
		},
		{
			"id": "1606",
			"city_id": "179",
			"name": "嘉鱼县"
		},
		{
			"id": "1607",
			"city_id": "179",
			"name": "通城县"
		},
		{
			"id": "1608",
			"city_id": "179",
			"name": "崇阳县"
		},
		{
			"id": "1609",
			"city_id": "179",
			"name": "通山县"
		},
		{
			"id": "1610",
			"city_id": "179",
			"name": "赤壁市"
		},
		{
			"id": "1611",
			"city_id": "180",
			"name": "曾都区"
		},
		{
			"id": "1612",
			"city_id": "180",
			"name": "广水市"
		},
		{
			"id": "1613",
			"city_id": "181",
			"name": "恩施市"
		},
		{
			"id": "1614",
			"city_id": "181",
			"name": "利川市"
		},
		{
			"id": "1615",
			"city_id": "181",
			"name": "建始县"
		},
		{
			"id": "1616",
			"city_id": "181",
			"name": "巴东县"
		},
		{
			"id": "1617",
			"city_id": "181",
			"name": "宣恩县"
		},
		{
			"id": "1618",
			"city_id": "181",
			"name": "咸丰县"
		},
		{
			"id": "1619",
			"city_id": "181",
			"name": "来凤县"
		},
		{
			"id": "1620",
			"city_id": "181",
			"name": "鹤峰县"
		},
		{
			"id": "1621",
			"city_id": "182",
			"name": "仙桃市"
		},
		{
			"id": "1622",
			"city_id": "182",
			"name": "潜江市"
		},
		{
			"id": "1623",
			"city_id": "182",
			"name": "天门市"
		},
		{
			"id": "1624",
			"city_id": "182",
			"name": "神农架林区"
		},
		{
			"id": "1625",
			"city_id": "183",
			"name": "芙蓉区"
		},
		{
			"id": "1626",
			"city_id": "183",
			"name": "天心区"
		},
		{
			"id": "1627",
			"city_id": "183",
			"name": "岳麓区"
		},
		{
			"id": "1628",
			"city_id": "183",
			"name": "开福区"
		},
		{
			"id": "1629",
			"city_id": "183",
			"name": "雨花区"
		},
		{
			"id": "1630",
			"city_id": "183",
			"name": "长沙县"
		},
		{
			"id": "1631",
			"city_id": "183",
			"name": "望城县"
		},
		{
			"id": "1632",
			"city_id": "183",
			"name": "宁乡县"
		},
		{
			"id": "1633",
			"city_id": "183",
			"name": "浏阳市"
		},
		{
			"id": "1634",
			"city_id": "184",
			"name": "荷塘区"
		},
		{
			"id": "1635",
			"city_id": "184",
			"name": "芦淞区"
		},
		{
			"id": "1636",
			"city_id": "184",
			"name": "石峰区"
		},
		{
			"id": "1637",
			"city_id": "184",
			"name": "天元区"
		},
		{
			"id": "1638",
			"city_id": "184",
			"name": "株洲县"
		},
		{
			"id": "1639",
			"city_id": "184",
			"name": "攸县"
		},
		{
			"id": "1640",
			"city_id": "184",
			"name": "茶陵县"
		},
		{
			"id": "1641",
			"city_id": "184",
			"name": "炎陵县"
		},
		{
			"id": "1642",
			"city_id": "184",
			"name": "醴陵市"
		},
		{
			"id": "1643",
			"city_id": "185",
			"name": "雨湖区"
		},
		{
			"id": "1644",
			"city_id": "185",
			"name": "岳塘区"
		},
		{
			"id": "1645",
			"city_id": "185",
			"name": "湘潭县"
		},
		{
			"id": "1646",
			"city_id": "185",
			"name": "湘乡市"
		},
		{
			"id": "1647",
			"city_id": "185",
			"name": "韶山市"
		},
		{
			"id": "1648",
			"city_id": "186",
			"name": "珠晖区"
		},
		{
			"id": "1649",
			"city_id": "186",
			"name": "雁峰区"
		},
		{
			"id": "1650",
			"city_id": "186",
			"name": "石鼓区"
		},
		{
			"id": "1651",
			"city_id": "186",
			"name": "蒸湘区"
		},
		{
			"id": "1652",
			"city_id": "186",
			"name": "南岳区"
		},
		{
			"id": "1653",
			"city_id": "186",
			"name": "衡阳县"
		},
		{
			"id": "1654",
			"city_id": "186",
			"name": "衡南县"
		},
		{
			"id": "1655",
			"city_id": "186",
			"name": "衡山县"
		},
		{
			"id": "1656",
			"city_id": "186",
			"name": "衡东县"
		},
		{
			"id": "1657",
			"city_id": "186",
			"name": "祁东县"
		},
		{
			"id": "1658",
			"city_id": "186",
			"name": "耒阳市"
		},
		{
			"id": "1659",
			"city_id": "186",
			"name": "常宁市"
		},
		{
			"id": "1660",
			"city_id": "187",
			"name": "双清区"
		},
		{
			"id": "1661",
			"city_id": "187",
			"name": "大祥区"
		},
		{
			"id": "1662",
			"city_id": "187",
			"name": "北塔区"
		},
		{
			"id": "1663",
			"city_id": "187",
			"name": "邵东县"
		},
		{
			"id": "1664",
			"city_id": "187",
			"name": "新邵县"
		},
		{
			"id": "1665",
			"city_id": "187",
			"name": "邵阳县"
		},
		{
			"id": "1666",
			"city_id": "187",
			"name": "隆回县"
		},
		{
			"id": "1667",
			"city_id": "187",
			"name": "洞口县"
		},
		{
			"id": "1668",
			"city_id": "187",
			"name": "绥宁县"
		},
		{
			"id": "1669",
			"city_id": "187",
			"name": "新宁县"
		},
		{
			"id": "1670",
			"city_id": "187",
			"name": "城步苗族自治县"
		},
		{
			"id": "1671",
			"city_id": "187",
			"name": "武冈市"
		},
		{
			"id": "1672",
			"city_id": "188",
			"name": "岳阳楼区"
		},
		{
			"id": "1673",
			"city_id": "188",
			"name": "云溪区"
		},
		{
			"id": "1674",
			"city_id": "188",
			"name": "君山区"
		},
		{
			"id": "1675",
			"city_id": "188",
			"name": "岳阳县"
		},
		{
			"id": "1676",
			"city_id": "188",
			"name": "华容县"
		},
		{
			"id": "1677",
			"city_id": "188",
			"name": "湘阴县"
		},
		{
			"id": "1678",
			"city_id": "188",
			"name": "平江县"
		},
		{
			"id": "1679",
			"city_id": "188",
			"name": "汨罗市"
		},
		{
			"id": "1680",
			"city_id": "188",
			"name": "临湘市"
		},
		{
			"id": "1681",
			"city_id": "189",
			"name": "武陵区"
		},
		{
			"id": "1682",
			"city_id": "189",
			"name": "鼎城区"
		},
		{
			"id": "1683",
			"city_id": "189",
			"name": "安乡县"
		},
		{
			"id": "1684",
			"city_id": "189",
			"name": "汉寿县"
		},
		{
			"id": "1685",
			"city_id": "189",
			"name": "澧县"
		},
		{
			"id": "1686",
			"city_id": "189",
			"name": "临澧县"
		},
		{
			"id": "1687",
			"city_id": "189",
			"name": "桃源县"
		},
		{
			"id": "1688",
			"city_id": "189",
			"name": "石门县"
		},
		{
			"id": "1689",
			"city_id": "189",
			"name": "津市市"
		},
		{
			"id": "1690",
			"city_id": "190",
			"name": "永定区"
		},
		{
			"id": "1691",
			"city_id": "190",
			"name": "武陵源区"
		},
		{
			"id": "1692",
			"city_id": "190",
			"name": "慈利县"
		},
		{
			"id": "1693",
			"city_id": "190",
			"name": "桑植县"
		},
		{
			"id": "1694",
			"city_id": "191",
			"name": "资阳区"
		},
		{
			"id": "1695",
			"city_id": "191",
			"name": "赫山区"
		},
		{
			"id": "1696",
			"city_id": "191",
			"name": "南县"
		},
		{
			"id": "1697",
			"city_id": "191",
			"name": "桃江县"
		},
		{
			"id": "1698",
			"city_id": "191",
			"name": "安化县"
		},
		{
			"id": "1699",
			"city_id": "191",
			"name": "沅江市"
		},
		{
			"id": "1700",
			"city_id": "192",
			"name": "北湖区"
		},
		{
			"id": "1701",
			"city_id": "192",
			"name": "苏仙区"
		},
		{
			"id": "1702",
			"city_id": "192",
			"name": "桂阳县"
		},
		{
			"id": "1703",
			"city_id": "192",
			"name": "宜章县"
		},
		{
			"id": "1704",
			"city_id": "192",
			"name": "永兴县"
		},
		{
			"id": "1705",
			"city_id": "192",
			"name": "嘉禾县"
		},
		{
			"id": "1706",
			"city_id": "192",
			"name": "临武县"
		},
		{
			"id": "1707",
			"city_id": "192",
			"name": "汝城县"
		},
		{
			"id": "1708",
			"city_id": "192",
			"name": "桂东县"
		},
		{
			"id": "1709",
			"city_id": "192",
			"name": "安仁县"
		},
		{
			"id": "1710",
			"city_id": "192",
			"name": "资兴市"
		},
		{
			"id": "1711",
			"city_id": "193",
			"name": "芝山区"
		},
		{
			"id": "1712",
			"city_id": "193",
			"name": "冷水滩区"
		},
		{
			"id": "1713",
			"city_id": "193",
			"name": "祁阳县"
		},
		{
			"id": "1714",
			"city_id": "193",
			"name": "东安县"
		},
		{
			"id": "1715",
			"city_id": "193",
			"name": "双牌县"
		},
		{
			"id": "1716",
			"city_id": "193",
			"name": "道县"
		},
		{
			"id": "1717",
			"city_id": "193",
			"name": "江永县"
		},
		{
			"id": "1718",
			"city_id": "193",
			"name": "宁远县"
		},
		{
			"id": "1719",
			"city_id": "193",
			"name": "蓝山县"
		},
		{
			"id": "1720",
			"city_id": "193",
			"name": "新田县"
		},
		{
			"id": "1721",
			"city_id": "193",
			"name": "江华瑶族自治县"
		},
		{
			"id": "1722",
			"city_id": "194",
			"name": "鹤城区"
		},
		{
			"id": "1723",
			"city_id": "194",
			"name": "中方县"
		},
		{
			"id": "1724",
			"city_id": "194",
			"name": "沅陵县"
		},
		{
			"id": "1725",
			"city_id": "194",
			"name": "辰溪县"
		},
		{
			"id": "1726",
			"city_id": "194",
			"name": "溆浦县"
		},
		{
			"id": "1727",
			"city_id": "194",
			"name": "会同县"
		},
		{
			"id": "1728",
			"city_id": "194",
			"name": "麻阳苗族自治县"
		},
		{
			"id": "1729",
			"city_id": "194",
			"name": "新晃侗族自治县"
		},
		{
			"id": "1730",
			"city_id": "194",
			"name": "芷江侗族自治县"
		},
		{
			"id": "1731",
			"city_id": "194",
			"name": "靖州苗族侗族自治县"
		},
		{
			"id": "1732",
			"city_id": "194",
			"name": "通道侗族自治县"
		},
		{
			"id": "1733",
			"city_id": "194",
			"name": "洪江市"
		},
		{
			"id": "1734",
			"city_id": "195",
			"name": "娄星区"
		},
		{
			"id": "1735",
			"city_id": "195",
			"name": "双峰县"
		},
		{
			"id": "1736",
			"city_id": "195",
			"name": "新化县"
		},
		{
			"id": "1737",
			"city_id": "195",
			"name": "冷水江市"
		},
		{
			"id": "1738",
			"city_id": "195",
			"name": "涟源市"
		},
		{
			"id": "1739",
			"city_id": "196",
			"name": "吉首市"
		},
		{
			"id": "1740",
			"city_id": "196",
			"name": "泸溪县"
		},
		{
			"id": "1741",
			"city_id": "196",
			"name": "凤凰县"
		},
		{
			"id": "1742",
			"city_id": "196",
			"name": "花垣县"
		},
		{
			"id": "1743",
			"city_id": "196",
			"name": "保靖县"
		},
		{
			"id": "1744",
			"city_id": "196",
			"name": "古丈县"
		},
		{
			"id": "1745",
			"city_id": "196",
			"name": "永顺县"
		},
		{
			"id": "1746",
			"city_id": "196",
			"name": "龙山县"
		},
		{
			"id": "1747",
			"city_id": "197",
			"name": "东山区"
		},
		{
			"id": "1748",
			"city_id": "197",
			"name": "荔湾区"
		},
		{
			"id": "1749",
			"city_id": "197",
			"name": "越秀区"
		},
		{
			"id": "1750",
			"city_id": "197",
			"name": "海珠区"
		},
		{
			"id": "1751",
			"city_id": "197",
			"name": "天河区"
		},
		{
			"id": "1752",
			"city_id": "197",
			"name": "芳村区"
		},
		{
			"id": "1753",
			"city_id": "197",
			"name": "白云区"
		},
		{
			"id": "1754",
			"city_id": "197",
			"name": "黄埔区"
		},
		{
			"id": "1755",
			"city_id": "197",
			"name": "番禺区"
		},
		{
			"id": "1756",
			"city_id": "197",
			"name": "花都区"
		},
		{
			"id": "1757",
			"city_id": "197",
			"name": "增城市"
		},
		{
			"id": "1758",
			"city_id": "197",
			"name": "从化市"
		},
		{
			"id": "1759",
			"city_id": "198",
			"name": "武江区"
		},
		{
			"id": "1760",
			"city_id": "198",
			"name": "浈江区"
		},
		{
			"id": "1761",
			"city_id": "198",
			"name": "曲江区"
		},
		{
			"id": "1762",
			"city_id": "198",
			"name": "始兴县"
		},
		{
			"id": "1763",
			"city_id": "198",
			"name": "仁化县"
		},
		{
			"id": "1764",
			"city_id": "198",
			"name": "翁源县"
		},
		{
			"id": "1765",
			"city_id": "198",
			"name": "乳源瑶族自治县"
		},
		{
			"id": "1766",
			"city_id": "198",
			"name": "新丰县"
		},
		{
			"id": "1767",
			"city_id": "198",
			"name": "乐昌市"
		},
		{
			"id": "1768",
			"city_id": "198",
			"name": "南雄市"
		},
		{
			"id": "1769",
			"city_id": "199",
			"name": "罗湖区"
		},
		{
			"id": "1770",
			"city_id": "199",
			"name": "福田区"
		},
		{
			"id": "1771",
			"city_id": "199",
			"name": "南山区"
		},
		{
			"id": "1772",
			"city_id": "199",
			"name": "宝安区"
		},
		{
			"id": "1773",
			"city_id": "199",
			"name": "龙岗区"
		},
		{
			"id": "1774",
			"city_id": "199",
			"name": "盐田区"
		},
		{
			"id": "1775",
			"city_id": "200",
			"name": "香洲区"
		},
		{
			"id": "1776",
			"city_id": "200",
			"name": "斗门区"
		},
		{
			"id": "1777",
			"city_id": "200",
			"name": "金湾区"
		},
		{
			"id": "1778",
			"city_id": "201",
			"name": "龙湖区"
		},
		{
			"id": "1779",
			"city_id": "201",
			"name": "金平区"
		},
		{
			"id": "1780",
			"city_id": "201",
			"name": "濠江区"
		},
		{
			"id": "1781",
			"city_id": "201",
			"name": "潮阳区"
		},
		{
			"id": "1782",
			"city_id": "201",
			"name": "潮南区"
		},
		{
			"id": "1783",
			"city_id": "201",
			"name": "澄海区"
		},
		{
			"id": "1784",
			"city_id": "201",
			"name": "南澳县"
		},
		{
			"id": "1785",
			"city_id": "202",
			"name": "禅城区"
		},
		{
			"id": "1786",
			"city_id": "202",
			"name": "南海区"
		},
		{
			"id": "1787",
			"city_id": "202",
			"name": "顺德区"
		},
		{
			"id": "1788",
			"city_id": "202",
			"name": "三水区"
		},
		{
			"id": "1789",
			"city_id": "202",
			"name": "高明区"
		},
		{
			"id": "1790",
			"city_id": "203",
			"name": "蓬江区"
		},
		{
			"id": "1791",
			"city_id": "203",
			"name": "江海区"
		},
		{
			"id": "1792",
			"city_id": "203",
			"name": "新会区"
		},
		{
			"id": "1793",
			"city_id": "203",
			"name": "台山市"
		},
		{
			"id": "1794",
			"city_id": "203",
			"name": "开平市"
		},
		{
			"id": "1795",
			"city_id": "203",
			"name": "鹤山市"
		},
		{
			"id": "1796",
			"city_id": "203",
			"name": "恩平市"
		},
		{
			"id": "1797",
			"city_id": "204",
			"name": "赤坎区"
		},
		{
			"id": "1798",
			"city_id": "204",
			"name": "霞山区"
		},
		{
			"id": "1799",
			"city_id": "204",
			"name": "坡头区"
		},
		{
			"id": "1800",
			"city_id": "204",
			"name": "麻章区"
		},
		{
			"id": "1801",
			"city_id": "204",
			"name": "遂溪县"
		},
		{
			"id": "1802",
			"city_id": "204",
			"name": "徐闻县"
		},
		{
			"id": "1803",
			"city_id": "204",
			"name": "廉江市"
		},
		{
			"id": "1804",
			"city_id": "204",
			"name": "雷州市"
		},
		{
			"id": "1805",
			"city_id": "204",
			"name": "吴川市"
		},
		{
			"id": "1806",
			"city_id": "205",
			"name": "茂南区"
		},
		{
			"id": "1807",
			"city_id": "205",
			"name": "茂港区"
		},
		{
			"id": "1808",
			"city_id": "205",
			"name": "电白县"
		},
		{
			"id": "1809",
			"city_id": "205",
			"name": "高州市"
		},
		{
			"id": "1810",
			"city_id": "205",
			"name": "化州市"
		},
		{
			"id": "1811",
			"city_id": "205",
			"name": "信宜市"
		},
		{
			"id": "1812",
			"city_id": "206",
			"name": "端州区"
		},
		{
			"id": "1813",
			"city_id": "206",
			"name": "鼎湖区"
		},
		{
			"id": "1814",
			"city_id": "206",
			"name": "广宁县"
		},
		{
			"id": "1815",
			"city_id": "206",
			"name": "怀集县"
		},
		{
			"id": "1816",
			"city_id": "206",
			"name": "封开县"
		},
		{
			"id": "1817",
			"city_id": "206",
			"name": "德庆县"
		},
		{
			"id": "1818",
			"city_id": "206",
			"name": "高要市"
		},
		{
			"id": "1819",
			"city_id": "206",
			"name": "四会市"
		},
		{
			"id": "1820",
			"city_id": "207",
			"name": "惠城区"
		},
		{
			"id": "1821",
			"city_id": "207",
			"name": "惠阳区"
		},
		{
			"id": "1822",
			"city_id": "207",
			"name": "博罗县"
		},
		{
			"id": "1823",
			"city_id": "207",
			"name": "惠东县"
		},
		{
			"id": "1824",
			"city_id": "207",
			"name": "龙门县"
		},
		{
			"id": "1825",
			"city_id": "208",
			"name": "梅江区"
		},
		{
			"id": "1826",
			"city_id": "208",
			"name": "梅县"
		},
		{
			"id": "1827",
			"city_id": "208",
			"name": "大埔县"
		},
		{
			"id": "1828",
			"city_id": "208",
			"name": "丰顺县"
		},
		{
			"id": "1829",
			"city_id": "208",
			"name": "五华县"
		},
		{
			"id": "1830",
			"city_id": "208",
			"name": "平远县"
		},
		{
			"id": "1831",
			"city_id": "208",
			"name": "蕉岭县"
		},
		{
			"id": "1832",
			"city_id": "208",
			"name": "兴宁市"
		},
		{
			"id": "1833",
			"city_id": "209",
			"name": "城区"
		},
		{
			"id": "1834",
			"city_id": "209",
			"name": "海丰县"
		},
		{
			"id": "1835",
			"city_id": "209",
			"name": "陆河县"
		},
		{
			"id": "1836",
			"city_id": "209",
			"name": "陆丰市"
		},
		{
			"id": "1837",
			"city_id": "210",
			"name": "源城区"
		},
		{
			"id": "1838",
			"city_id": "210",
			"name": "紫金县"
		},
		{
			"id": "1839",
			"city_id": "210",
			"name": "龙川县"
		},
		{
			"id": "1840",
			"city_id": "210",
			"name": "连平县"
		},
		{
			"id": "1841",
			"city_id": "210",
			"name": "和平县"
		},
		{
			"id": "1842",
			"city_id": "210",
			"name": "东源县"
		},
		{
			"id": "1843",
			"city_id": "211",
			"name": "江城区"
		},
		{
			"id": "1844",
			"city_id": "211",
			"name": "阳西县"
		},
		{
			"id": "1845",
			"city_id": "211",
			"name": "阳东县"
		},
		{
			"id": "1846",
			"city_id": "211",
			"name": "阳春市"
		},
		{
			"id": "1847",
			"city_id": "212",
			"name": "清城区"
		},
		{
			"id": "1848",
			"city_id": "212",
			"name": "佛冈县"
		},
		{
			"id": "1849",
			"city_id": "212",
			"name": "阳山县"
		},
		{
			"id": "1850",
			"city_id": "212",
			"name": "连山壮族瑶族自治县"
		},
		{
			"id": "1851",
			"city_id": "212",
			"name": "连南瑶族自治县"
		},
		{
			"id": "1852",
			"city_id": "212",
			"name": "清新县"
		},
		{
			"id": "1853",
			"city_id": "212",
			"name": "英德市"
		},
		{
			"id": "1854",
			"city_id": "212",
			"name": "连州市"
		},
		{
			"id": "1855",
			"city_id": "215",
			"name": "湘桥区"
		},
		{
			"id": "1856",
			"city_id": "215",
			"name": "潮安县"
		},
		{
			"id": "1857",
			"city_id": "215",
			"name": "饶平县"
		},
		{
			"id": "1858",
			"city_id": "216",
			"name": "榕城区"
		},
		{
			"id": "1859",
			"city_id": "216",
			"name": "揭东县"
		},
		{
			"id": "1860",
			"city_id": "216",
			"name": "揭西县"
		},
		{
			"id": "1861",
			"city_id": "216",
			"name": "惠来县"
		},
		{
			"id": "1862",
			"city_id": "216",
			"name": "普宁市"
		},
		{
			"id": "1863",
			"city_id": "217",
			"name": "云城区"
		},
		{
			"id": "1864",
			"city_id": "217",
			"name": "新兴县"
		},
		{
			"id": "1865",
			"city_id": "217",
			"name": "郁南县"
		},
		{
			"id": "1866",
			"city_id": "217",
			"name": "云安县"
		},
		{
			"id": "1867",
			"city_id": "217",
			"name": "罗定市"
		},
		{
			"id": "1868",
			"city_id": "218",
			"name": "兴宁区"
		},
		{
			"id": "1869",
			"city_id": "218",
			"name": "青秀区"
		},
		{
			"id": "1870",
			"city_id": "218",
			"name": "江南区"
		},
		{
			"id": "1871",
			"city_id": "218",
			"name": "西乡塘区"
		},
		{
			"id": "1872",
			"city_id": "218",
			"name": "良庆区"
		},
		{
			"id": "1873",
			"city_id": "218",
			"name": "邕宁区"
		},
		{
			"id": "1874",
			"city_id": "218",
			"name": "武鸣县"
		},
		{
			"id": "1875",
			"city_id": "218",
			"name": "隆安县"
		},
		{
			"id": "1876",
			"city_id": "218",
			"name": "马山县"
		},
		{
			"id": "1877",
			"city_id": "218",
			"name": "上林县"
		},
		{
			"id": "1878",
			"city_id": "218",
			"name": "宾阳县"
		},
		{
			"id": "1879",
			"city_id": "218",
			"name": "横县"
		},
		{
			"id": "1880",
			"city_id": "219",
			"name": "城中区"
		},
		{
			"id": "1881",
			"city_id": "219",
			"name": "鱼峰区"
		},
		{
			"id": "1882",
			"city_id": "219",
			"name": "柳南区"
		},
		{
			"id": "1883",
			"city_id": "219",
			"name": "柳北区"
		},
		{
			"id": "1884",
			"city_id": "219",
			"name": "柳江县"
		},
		{
			"id": "1885",
			"city_id": "219",
			"name": "柳城县"
		},
		{
			"id": "1886",
			"city_id": "219",
			"name": "鹿寨县"
		},
		{
			"id": "1887",
			"city_id": "219",
			"name": "融安县"
		},
		{
			"id": "1888",
			"city_id": "219",
			"name": "融水苗族自治县"
		},
		{
			"id": "1889",
			"city_id": "219",
			"name": "三江侗族自治县"
		},
		{
			"id": "1890",
			"city_id": "220",
			"name": "秀峰区"
		},
		{
			"id": "1891",
			"city_id": "220",
			"name": "叠彩区"
		},
		{
			"id": "1892",
			"city_id": "220",
			"name": "象山区"
		},
		{
			"id": "1893",
			"city_id": "220",
			"name": "七星区"
		},
		{
			"id": "1894",
			"city_id": "220",
			"name": "雁山区"
		},
		{
			"id": "1895",
			"city_id": "220",
			"name": "阳朔县"
		},
		{
			"id": "1896",
			"city_id": "220",
			"name": "临桂县"
		},
		{
			"id": "1897",
			"city_id": "220",
			"name": "灵川县"
		},
		{
			"id": "1898",
			"city_id": "220",
			"name": "全州县"
		},
		{
			"id": "1899",
			"city_id": "220",
			"name": "兴安县"
		},
		{
			"id": "1900",
			"city_id": "220",
			"name": "永福县"
		},
		{
			"id": "1901",
			"city_id": "220",
			"name": "灌阳县"
		},
		{
			"id": "1902",
			"city_id": "220",
			"name": "龙胜各族自治县"
		},
		{
			"id": "1903",
			"city_id": "220",
			"name": "资源县"
		},
		{
			"id": "1904",
			"city_id": "220",
			"name": "平乐县"
		},
		{
			"id": "1905",
			"city_id": "220",
			"name": "荔蒲县"
		},
		{
			"id": "1906",
			"city_id": "220",
			"name": "恭城瑶族自治县"
		},
		{
			"id": "1907",
			"city_id": "221",
			"name": "万秀区"
		},
		{
			"id": "1908",
			"city_id": "221",
			"name": "蝶山区"
		},
		{
			"id": "1909",
			"city_id": "221",
			"name": "长洲区"
		},
		{
			"id": "1910",
			"city_id": "221",
			"name": "苍梧县"
		},
		{
			"id": "1911",
			"city_id": "221",
			"name": "藤县"
		},
		{
			"id": "1912",
			"city_id": "221",
			"name": "蒙山县"
		},
		{
			"id": "1913",
			"city_id": "221",
			"name": "岑溪市"
		},
		{
			"id": "1914",
			"city_id": "222",
			"name": "海城区"
		},
		{
			"id": "1915",
			"city_id": "222",
			"name": "银海区"
		},
		{
			"id": "1916",
			"city_id": "222",
			"name": "铁山港区"
		},
		{
			"id": "1917",
			"city_id": "222",
			"name": "合浦县"
		},
		{
			"id": "1918",
			"city_id": "223",
			"name": "港口区"
		},
		{
			"id": "1919",
			"city_id": "223",
			"name": "防城区"
		},
		{
			"id": "1920",
			"city_id": "223",
			"name": "上思县"
		},
		{
			"id": "1921",
			"city_id": "223",
			"name": "东兴市"
		},
		{
			"id": "1922",
			"city_id": "224",
			"name": "钦南区"
		},
		{
			"id": "1923",
			"city_id": "224",
			"name": "钦北区"
		},
		{
			"id": "1924",
			"city_id": "224",
			"name": "灵山县"
		},
		{
			"id": "1925",
			"city_id": "224",
			"name": "浦北县"
		},
		{
			"id": "1926",
			"city_id": "225",
			"name": "港北区"
		},
		{
			"id": "1927",
			"city_id": "225",
			"name": "港南区"
		},
		{
			"id": "1928",
			"city_id": "225",
			"name": "覃塘区"
		},
		{
			"id": "1929",
			"city_id": "225",
			"name": "平南县"
		},
		{
			"id": "1930",
			"city_id": "225",
			"name": "桂平市"
		},
		{
			"id": "1931",
			"city_id": "226",
			"name": "玉州区"
		},
		{
			"id": "1932",
			"city_id": "226",
			"name": "容县"
		},
		{
			"id": "1933",
			"city_id": "226",
			"name": "陆川县"
		},
		{
			"id": "1934",
			"city_id": "226",
			"name": "博白县"
		},
		{
			"id": "1935",
			"city_id": "226",
			"name": "兴业县"
		},
		{
			"id": "1936",
			"city_id": "226",
			"name": "北流市"
		},
		{
			"id": "1937",
			"city_id": "227",
			"name": "右江区"
		},
		{
			"id": "1938",
			"city_id": "227",
			"name": "田阳县"
		},
		{
			"id": "1939",
			"city_id": "227",
			"name": "田东县"
		},
		{
			"id": "1940",
			"city_id": "227",
			"name": "平果县"
		},
		{
			"id": "1941",
			"city_id": "227",
			"name": "德保县"
		},
		{
			"id": "1942",
			"city_id": "227",
			"name": "靖西县"
		},
		{
			"id": "1943",
			"city_id": "227",
			"name": "那坡县"
		},
		{
			"id": "1944",
			"city_id": "227",
			"name": "凌云县"
		},
		{
			"id": "1945",
			"city_id": "227",
			"name": "乐业县"
		},
		{
			"id": "1946",
			"city_id": "227",
			"name": "田林县"
		},
		{
			"id": "1947",
			"city_id": "227",
			"name": "西林县"
		},
		{
			"id": "1948",
			"city_id": "227",
			"name": "隆林各族自治县"
		},
		{
			"id": "1949",
			"city_id": "228",
			"name": "八步区"
		},
		{
			"id": "1950",
			"city_id": "228",
			"name": "昭平县"
		},
		{
			"id": "1951",
			"city_id": "228",
			"name": "钟山县"
		},
		{
			"id": "1952",
			"city_id": "228",
			"name": "富川瑶族自治县"
		},
		{
			"id": "1953",
			"city_id": "229",
			"name": "金城江区"
		},
		{
			"id": "1954",
			"city_id": "229",
			"name": "南丹县"
		},
		{
			"id": "1955",
			"city_id": "229",
			"name": "天峨县"
		},
		{
			"id": "1956",
			"city_id": "229",
			"name": "凤山县"
		},
		{
			"id": "1957",
			"city_id": "229",
			"name": "东兰县"
		},
		{
			"id": "1958",
			"city_id": "229",
			"name": "罗城仫佬族自治县"
		},
		{
			"id": "1959",
			"city_id": "229",
			"name": "环江毛南族自治县"
		},
		{
			"id": "1960",
			"city_id": "229",
			"name": "巴马瑶族自治县"
		},
		{
			"id": "1961",
			"city_id": "229",
			"name": "都安瑶族自治县"
		},
		{
			"id": "1962",
			"city_id": "229",
			"name": "大化瑶族自治县"
		},
		{
			"id": "1963",
			"city_id": "229",
			"name": "宜州市"
		},
		{
			"id": "1964",
			"city_id": "230",
			"name": "兴宾区"
		},
		{
			"id": "1965",
			"city_id": "230",
			"name": "忻城县"
		},
		{
			"id": "1966",
			"city_id": "230",
			"name": "象州县"
		},
		{
			"id": "1967",
			"city_id": "230",
			"name": "武宣县"
		},
		{
			"id": "1968",
			"city_id": "230",
			"name": "金秀瑶族自治县"
		},
		{
			"id": "1969",
			"city_id": "230",
			"name": "合山市"
		},
		{
			"id": "1970",
			"city_id": "231",
			"name": "江洲区"
		},
		{
			"id": "1971",
			"city_id": "231",
			"name": "扶绥县"
		},
		{
			"id": "1972",
			"city_id": "231",
			"name": "宁明县"
		},
		{
			"id": "1973",
			"city_id": "231",
			"name": "龙州县"
		},
		{
			"id": "1974",
			"city_id": "231",
			"name": "大新县"
		},
		{
			"id": "1975",
			"city_id": "231",
			"name": "天等县"
		},
		{
			"id": "1976",
			"city_id": "231",
			"name": "凭祥市"
		},
		{
			"id": "1977",
			"city_id": "232",
			"name": "秀英区"
		},
		{
			"id": "1978",
			"city_id": "232",
			"name": "龙华区"
		},
		{
			"id": "1979",
			"city_id": "232",
			"name": "琼山区"
		},
		{
			"id": "1980",
			"city_id": "232",
			"name": "美兰区"
		},
		{
			"id": "1981",
			"city_id": "233",
			"name": "五指山市"
		},
		{
			"id": "1982",
			"city_id": "233",
			"name": "琼海市"
		},
		{
			"id": "1983",
			"city_id": "233",
			"name": "儋州市"
		},
		{
			"id": "1984",
			"city_id": "233",
			"name": "文昌市"
		},
		{
			"id": "1985",
			"city_id": "233",
			"name": "万宁市"
		},
		{
			"id": "1986",
			"city_id": "233",
			"name": "东方市"
		},
		{
			"id": "1987",
			"city_id": "233",
			"name": "定安县"
		},
		{
			"id": "1988",
			"city_id": "233",
			"name": "屯昌县"
		},
		{
			"id": "1989",
			"city_id": "233",
			"name": "澄迈县"
		},
		{
			"id": "1990",
			"city_id": "233",
			"name": "临高县"
		},
		{
			"id": "1991",
			"city_id": "233",
			"name": "白沙黎族自治县"
		},
		{
			"id": "1992",
			"city_id": "233",
			"name": "昌江黎族自治县"
		},
		{
			"id": "1993",
			"city_id": "233",
			"name": "乐东黎族自治县"
		},
		{
			"id": "1994",
			"city_id": "233",
			"name": "陵水黎族自治县"
		},
		{
			"id": "1995",
			"city_id": "233",
			"name": "保亭黎族苗族自治县"
		},
		{
			"id": "1996",
			"city_id": "233",
			"name": "琼中黎族苗族自治县"
		},
		{
			"id": "1997",
			"city_id": "233",
			"name": "西沙群岛"
		},
		{
			"id": "1998",
			"city_id": "233",
			"name": "南沙群岛"
		},
		{
			"id": "1999",
			"city_id": "233",
			"name": "中沙群岛的岛礁及其海域"
		},
		{
			"id": "2000",
			"city_id": "234",
			"name": "万州区"
		},
		{
			"id": "2001",
			"city_id": "234",
			"name": "涪陵区"
		},
		{
			"id": "2002",
			"city_id": "234",
			"name": "渝中区"
		},
		{
			"id": "2003",
			"city_id": "234",
			"name": "大渡口区"
		},
		{
			"id": "2004",
			"city_id": "234",
			"name": "江北区"
		},
		{
			"id": "2005",
			"city_id": "234",
			"name": "沙坪坝区"
		},
		{
			"id": "2006",
			"city_id": "234",
			"name": "九龙坡区"
		},
		{
			"id": "2007",
			"city_id": "234",
			"name": "南岸区"
		},
		{
			"id": "2008",
			"city_id": "234",
			"name": "北碚区"
		},
		{
			"id": "2009",
			"city_id": "234",
			"name": "万盛区"
		},
		{
			"id": "2010",
			"city_id": "234",
			"name": "双桥区"
		},
		{
			"id": "2011",
			"city_id": "234",
			"name": "渝北区"
		},
		{
			"id": "2012",
			"city_id": "234",
			"name": "巴南区"
		},
		{
			"id": "2013",
			"city_id": "234",
			"name": "黔江区"
		},
		{
			"id": "2014",
			"city_id": "234",
			"name": "长寿区"
		},
		{
			"id": "2015",
			"city_id": "234",
			"name": "綦江县"
		},
		{
			"id": "2016",
			"city_id": "234",
			"name": "潼南县"
		},
		{
			"id": "2017",
			"city_id": "234",
			"name": "铜梁县"
		},
		{
			"id": "2018",
			"city_id": "234",
			"name": "大足县"
		},
		{
			"id": "2019",
			"city_id": "234",
			"name": "荣昌县"
		},
		{
			"id": "2020",
			"city_id": "234",
			"name": "璧山县"
		},
		{
			"id": "2021",
			"city_id": "234",
			"name": "梁平县"
		},
		{
			"id": "2022",
			"city_id": "234",
			"name": "城口县"
		},
		{
			"id": "2023",
			"city_id": "234",
			"name": "丰都县"
		},
		{
			"id": "2024",
			"city_id": "234",
			"name": "垫江县"
		},
		{
			"id": "2025",
			"city_id": "234",
			"name": "武隆县"
		},
		{
			"id": "2026",
			"city_id": "234",
			"name": "忠县"
		},
		{
			"id": "2027",
			"city_id": "234",
			"name": "开县"
		},
		{
			"id": "2028",
			"city_id": "234",
			"name": "云阳县"
		},
		{
			"id": "2029",
			"city_id": "234",
			"name": "奉节县"
		},
		{
			"id": "2030",
			"city_id": "234",
			"name": "巫山县"
		},
		{
			"id": "2031",
			"city_id": "234",
			"name": "巫溪县"
		},
		{
			"id": "2032",
			"city_id": "234",
			"name": "石柱土家族自治县"
		},
		{
			"id": "2033",
			"city_id": "234",
			"name": "秀山土家族苗族自治县"
		},
		{
			"id": "2034",
			"city_id": "234",
			"name": "酉阳土家族苗族自治县"
		},
		{
			"id": "2035",
			"city_id": "234",
			"name": "彭水苗族土家族自治县"
		},
		{
			"id": "2036",
			"city_id": "234",
			"name": "江津市"
		},
		{
			"id": "2037",
			"city_id": "234",
			"name": "合川市"
		},
		{
			"id": "2038",
			"city_id": "234",
			"name": "永川市"
		},
		{
			"id": "2039",
			"city_id": "234",
			"name": "南川市"
		},
		{
			"id": "2040",
			"city_id": "235",
			"name": "锦江区"
		},
		{
			"id": "2041",
			"city_id": "235",
			"name": "青羊区"
		},
		{
			"id": "2042",
			"city_id": "235",
			"name": "金牛区"
		},
		{
			"id": "2043",
			"city_id": "235",
			"name": "武侯区"
		},
		{
			"id": "2044",
			"city_id": "235",
			"name": "成华区"
		},
		{
			"id": "2045",
			"city_id": "235",
			"name": "龙泉驿区"
		},
		{
			"id": "2046",
			"city_id": "235",
			"name": "青白江区"
		},
		{
			"id": "2047",
			"city_id": "235",
			"name": "新都区"
		},
		{
			"id": "2048",
			"city_id": "235",
			"name": "温江区"
		},
		{
			"id": "2049",
			"city_id": "235",
			"name": "金堂县"
		},
		{
			"id": "2050",
			"city_id": "235",
			"name": "双流县"
		},
		{
			"id": "2051",
			"city_id": "235",
			"name": "郫县"
		},
		{
			"id": "2052",
			"city_id": "235",
			"name": "大邑县"
		},
		{
			"id": "2053",
			"city_id": "235",
			"name": "蒲江县"
		},
		{
			"id": "2054",
			"city_id": "235",
			"name": "新津县"
		},
		{
			"id": "2055",
			"city_id": "235",
			"name": "都江堰市"
		},
		{
			"id": "2056",
			"city_id": "235",
			"name": "彭州市"
		},
		{
			"id": "2057",
			"city_id": "235",
			"name": "邛崃市"
		},
		{
			"id": "2058",
			"city_id": "235",
			"name": "崇州市"
		},
		{
			"id": "2059",
			"city_id": "236",
			"name": "自流井区"
		},
		{
			"id": "2060",
			"city_id": "236",
			"name": "贡井区"
		},
		{
			"id": "2061",
			"city_id": "236",
			"name": "大安区"
		},
		{
			"id": "2062",
			"city_id": "236",
			"name": "沿滩区"
		},
		{
			"id": "2063",
			"city_id": "236",
			"name": "荣县"
		},
		{
			"id": "2064",
			"city_id": "236",
			"name": "富顺县"
		},
		{
			"id": "2065",
			"city_id": "237",
			"name": "东区"
		},
		{
			"id": "2066",
			"city_id": "237",
			"name": "西区"
		},
		{
			"id": "2067",
			"city_id": "237",
			"name": "仁和区"
		},
		{
			"id": "2068",
			"city_id": "237",
			"name": "米易县"
		},
		{
			"id": "2069",
			"city_id": "237",
			"name": "盐边县"
		},
		{
			"id": "2070",
			"city_id": "238",
			"name": "江阳区"
		},
		{
			"id": "2071",
			"city_id": "238",
			"name": "纳溪区"
		},
		{
			"id": "2072",
			"city_id": "238",
			"name": "龙马潭区"
		},
		{
			"id": "2073",
			"city_id": "238",
			"name": "泸县"
		},
		{
			"id": "2074",
			"city_id": "238",
			"name": "合江县"
		},
		{
			"id": "2075",
			"city_id": "238",
			"name": "叙永县"
		},
		{
			"id": "2076",
			"city_id": "238",
			"name": "古蔺县"
		},
		{
			"id": "2077",
			"city_id": "239",
			"name": "旌阳区"
		},
		{
			"id": "2078",
			"city_id": "239",
			"name": "中江县"
		},
		{
			"id": "2079",
			"city_id": "239",
			"name": "罗江县"
		},
		{
			"id": "2080",
			"city_id": "239",
			"name": "广汉市"
		},
		{
			"id": "2081",
			"city_id": "239",
			"name": "什邡市"
		},
		{
			"id": "2082",
			"city_id": "239",
			"name": "绵竹市"
		},
		{
			"id": "2083",
			"city_id": "240",
			"name": "涪城区"
		},
		{
			"id": "2084",
			"city_id": "240",
			"name": "游仙区"
		},
		{
			"id": "2085",
			"city_id": "240",
			"name": "三台县"
		},
		{
			"id": "2086",
			"city_id": "240",
			"name": "盐亭县"
		},
		{
			"id": "2087",
			"city_id": "240",
			"name": "安县"
		},
		{
			"id": "2088",
			"city_id": "240",
			"name": "梓潼县"
		},
		{
			"id": "2089",
			"city_id": "240",
			"name": "北川羌族自治县"
		},
		{
			"id": "2090",
			"city_id": "240",
			"name": "平武县"
		},
		{
			"id": "2091",
			"city_id": "240",
			"name": "江油市"
		},
		{
			"id": "2092",
			"city_id": "241",
			"name": "市中区"
		},
		{
			"id": "2093",
			"city_id": "241",
			"name": "元坝区"
		},
		{
			"id": "2094",
			"city_id": "241",
			"name": "朝天区"
		},
		{
			"id": "2095",
			"city_id": "241",
			"name": "旺苍县"
		},
		{
			"id": "2096",
			"city_id": "241",
			"name": "青川县"
		},
		{
			"id": "2097",
			"city_id": "241",
			"name": "剑阁县"
		},
		{
			"id": "2098",
			"city_id": "241",
			"name": "苍溪县"
		},
		{
			"id": "2099",
			"city_id": "242",
			"name": "船山区"
		},
		{
			"id": "2100",
			"city_id": "242",
			"name": "安居区"
		},
		{
			"id": "2101",
			"city_id": "242",
			"name": "蓬溪县"
		},
		{
			"id": "2102",
			"city_id": "242",
			"name": "射洪县"
		},
		{
			"id": "2103",
			"city_id": "242",
			"name": "大英县"
		},
		{
			"id": "2104",
			"city_id": "243",
			"name": "市中区"
		},
		{
			"id": "2105",
			"city_id": "243",
			"name": "东兴区"
		},
		{
			"id": "2106",
			"city_id": "243",
			"name": "威远县"
		},
		{
			"id": "2107",
			"city_id": "243",
			"name": "资中县"
		},
		{
			"id": "2108",
			"city_id": "243",
			"name": "隆昌县"
		},
		{
			"id": "2109",
			"city_id": "244",
			"name": "市中区"
		},
		{
			"id": "2110",
			"city_id": "244",
			"name": "沙湾区"
		},
		{
			"id": "2111",
			"city_id": "244",
			"name": "五通桥区"
		},
		{
			"id": "2112",
			"city_id": "244",
			"name": "金口河区"
		},
		{
			"id": "2113",
			"city_id": "244",
			"name": "犍为县"
		},
		{
			"id": "2114",
			"city_id": "244",
			"name": "井研县"
		},
		{
			"id": "2115",
			"city_id": "244",
			"name": "夹江县"
		},
		{
			"id": "2116",
			"city_id": "244",
			"name": "沐川县"
		},
		{
			"id": "2117",
			"city_id": "244",
			"name": "峨边彝族自治县"
		},
		{
			"id": "2118",
			"city_id": "244",
			"name": "马边彝族自治县"
		},
		{
			"id": "2119",
			"city_id": "244",
			"name": "峨眉山市"
		},
		{
			"id": "2120",
			"city_id": "245",
			"name": "顺庆区"
		},
		{
			"id": "2121",
			"city_id": "245",
			"name": "高坪区"
		},
		{
			"id": "2122",
			"city_id": "245",
			"name": "嘉陵区"
		},
		{
			"id": "2123",
			"city_id": "245",
			"name": "南部县"
		},
		{
			"id": "2124",
			"city_id": "245",
			"name": "营山县"
		},
		{
			"id": "2125",
			"city_id": "245",
			"name": "蓬安县"
		},
		{
			"id": "2126",
			"city_id": "245",
			"name": "仪陇县"
		},
		{
			"id": "2127",
			"city_id": "245",
			"name": "西充县"
		},
		{
			"id": "2128",
			"city_id": "245",
			"name": "阆中市"
		},
		{
			"id": "2129",
			"city_id": "246",
			"name": "东坡区"
		},
		{
			"id": "2130",
			"city_id": "246",
			"name": "仁寿县"
		},
		{
			"id": "2131",
			"city_id": "246",
			"name": "彭山县"
		},
		{
			"id": "2132",
			"city_id": "246",
			"name": "洪雅县"
		},
		{
			"id": "2133",
			"city_id": "246",
			"name": "丹棱县"
		},
		{
			"id": "2134",
			"city_id": "246",
			"name": "青神县"
		},
		{
			"id": "2135",
			"city_id": "247",
			"name": "翠屏区"
		},
		{
			"id": "2136",
			"city_id": "247",
			"name": "宜宾县"
		},
		{
			"id": "2137",
			"city_id": "247",
			"name": "南溪县"
		},
		{
			"id": "2138",
			"city_id": "247",
			"name": "江安县"
		},
		{
			"id": "2139",
			"city_id": "247",
			"name": "长宁县"
		},
		{
			"id": "2140",
			"city_id": "247",
			"name": "高县"
		},
		{
			"id": "2141",
			"city_id": "247",
			"name": "珙县"
		},
		{
			"id": "2142",
			"city_id": "247",
			"name": "筠连县"
		},
		{
			"id": "2143",
			"city_id": "247",
			"name": "兴文县"
		},
		{
			"id": "2144",
			"city_id": "247",
			"name": "屏山县"
		},
		{
			"id": "2145",
			"city_id": "248",
			"name": "广安区"
		},
		{
			"id": "2146",
			"city_id": "248",
			"name": "岳池县"
		},
		{
			"id": "2147",
			"city_id": "248",
			"name": "武胜县"
		},
		{
			"id": "2148",
			"city_id": "248",
			"name": "邻水县"
		},
		{
			"id": "2149",
			"city_id": "248",
			"name": "华蓥市"
		},
		{
			"id": "2150",
			"city_id": "249",
			"name": "通川区"
		},
		{
			"id": "2151",
			"city_id": "249",
			"name": "达县"
		},
		{
			"id": "2152",
			"city_id": "249",
			"name": "宣汉县"
		},
		{
			"id": "2153",
			"city_id": "249",
			"name": "开江县"
		},
		{
			"id": "2154",
			"city_id": "249",
			"name": "大竹县"
		},
		{
			"id": "2155",
			"city_id": "249",
			"name": "渠县"
		},
		{
			"id": "2156",
			"city_id": "249",
			"name": "万源市"
		},
		{
			"id": "2157",
			"city_id": "250",
			"name": "雨城区"
		},
		{
			"id": "2158",
			"city_id": "250",
			"name": "名山县"
		},
		{
			"id": "2159",
			"city_id": "250",
			"name": "荥经县"
		},
		{
			"id": "2160",
			"city_id": "250",
			"name": "汉源县"
		},
		{
			"id": "2161",
			"city_id": "250",
			"name": "石棉县"
		},
		{
			"id": "2162",
			"city_id": "250",
			"name": "天全县"
		},
		{
			"id": "2163",
			"city_id": "250",
			"name": "芦山县"
		},
		{
			"id": "2164",
			"city_id": "250",
			"name": "宝兴县"
		},
		{
			"id": "2165",
			"city_id": "251",
			"name": "巴州区"
		},
		{
			"id": "2166",
			"city_id": "251",
			"name": "通江县"
		},
		{
			"id": "2167",
			"city_id": "251",
			"name": "南江县"
		},
		{
			"id": "2168",
			"city_id": "251",
			"name": "平昌县"
		},
		{
			"id": "2169",
			"city_id": "252",
			"name": "雁江区"
		},
		{
			"id": "2170",
			"city_id": "252",
			"name": "安岳县"
		},
		{
			"id": "2171",
			"city_id": "252",
			"name": "乐至县"
		},
		{
			"id": "2172",
			"city_id": "252",
			"name": "简阳市"
		},
		{
			"id": "2173",
			"city_id": "253",
			"name": "汶川县"
		},
		{
			"id": "2174",
			"city_id": "253",
			"name": "理县"
		},
		{
			"id": "2175",
			"city_id": "253",
			"name": "茂县"
		},
		{
			"id": "2176",
			"city_id": "253",
			"name": "松潘县"
		},
		{
			"id": "2177",
			"city_id": "253",
			"name": "九寨沟县"
		},
		{
			"id": "2178",
			"city_id": "253",
			"name": "金川县"
		},
		{
			"id": "2179",
			"city_id": "253",
			"name": "小金县"
		},
		{
			"id": "2180",
			"city_id": "253",
			"name": "黑水县"
		},
		{
			"id": "2181",
			"city_id": "253",
			"name": "马尔康县"
		},
		{
			"id": "2182",
			"city_id": "253",
			"name": "壤塘县"
		},
		{
			"id": "2183",
			"city_id": "253",
			"name": "阿坝县"
		},
		{
			"id": "2184",
			"city_id": "253",
			"name": "若尔盖县"
		},
		{
			"id": "2185",
			"city_id": "253",
			"name": "红原县"
		},
		{
			"id": "2186",
			"city_id": "254",
			"name": "康定县"
		},
		{
			"id": "2187",
			"city_id": "254",
			"name": "泸定县"
		},
		{
			"id": "2188",
			"city_id": "254",
			"name": "丹巴县"
		},
		{
			"id": "2189",
			"city_id": "254",
			"name": "九龙县"
		},
		{
			"id": "2190",
			"city_id": "254",
			"name": "雅江县"
		},
		{
			"id": "2191",
			"city_id": "254",
			"name": "道孚县"
		},
		{
			"id": "2192",
			"city_id": "254",
			"name": "炉霍县"
		},
		{
			"id": "2193",
			"city_id": "254",
			"name": "甘孜县"
		},
		{
			"id": "2194",
			"city_id": "254",
			"name": "新龙县"
		},
		{
			"id": "2195",
			"city_id": "254",
			"name": "德格县"
		},
		{
			"id": "2196",
			"city_id": "254",
			"name": "白玉县"
		},
		{
			"id": "2197",
			"city_id": "254",
			"name": "石渠县"
		},
		{
			"id": "2198",
			"city_id": "254",
			"name": "色达县"
		},
		{
			"id": "2199",
			"city_id": "254",
			"name": "理塘县"
		},
		{
			"id": "2200",
			"city_id": "254",
			"name": "巴塘县"
		},
		{
			"id": "2201",
			"city_id": "254",
			"name": "乡城县"
		},
		{
			"id": "2202",
			"city_id": "254",
			"name": "稻城县"
		},
		{
			"id": "2203",
			"city_id": "254",
			"name": "得荣县"
		},
		{
			"id": "2204",
			"city_id": "255",
			"name": "西昌市"
		},
		{
			"id": "2205",
			"city_id": "255",
			"name": "木里藏族自治县"
		},
		{
			"id": "2206",
			"city_id": "255",
			"name": "盐源县"
		},
		{
			"id": "2207",
			"city_id": "255",
			"name": "德昌县"
		},
		{
			"id": "2208",
			"city_id": "255",
			"name": "会理县"
		},
		{
			"id": "2209",
			"city_id": "255",
			"name": "会东县"
		},
		{
			"id": "2210",
			"city_id": "255",
			"name": "宁南县"
		},
		{
			"id": "2211",
			"city_id": "255",
			"name": "普格县"
		},
		{
			"id": "2212",
			"city_id": "255",
			"name": "布拖县"
		},
		{
			"id": "2213",
			"city_id": "255",
			"name": "金阳县"
		},
		{
			"id": "2214",
			"city_id": "255",
			"name": "昭觉县"
		},
		{
			"id": "2215",
			"city_id": "255",
			"name": "喜德县"
		},
		{
			"id": "2216",
			"city_id": "255",
			"name": "冕宁县"
		},
		{
			"id": "2217",
			"city_id": "255",
			"name": "越西县"
		},
		{
			"id": "2218",
			"city_id": "255",
			"name": "甘洛县"
		},
		{
			"id": "2219",
			"city_id": "255",
			"name": "美姑县"
		},
		{
			"id": "2220",
			"city_id": "255",
			"name": "雷波县"
		},
		{
			"id": "2221",
			"city_id": "256",
			"name": "南明区"
		},
		{
			"id": "2222",
			"city_id": "256",
			"name": "云岩区"
		},
		{
			"id": "2223",
			"city_id": "256",
			"name": "花溪区"
		},
		{
			"id": "2224",
			"city_id": "256",
			"name": "乌当区"
		},
		{
			"id": "2225",
			"city_id": "256",
			"name": "白云区"
		},
		{
			"id": "2226",
			"city_id": "256",
			"name": "小河区"
		},
		{
			"id": "2227",
			"city_id": "256",
			"name": "开阳县"
		},
		{
			"id": "2228",
			"city_id": "256",
			"name": "息烽县"
		},
		{
			"id": "2229",
			"city_id": "256",
			"name": "修文县"
		},
		{
			"id": "2230",
			"city_id": "256",
			"name": "清镇市"
		},
		{
			"id": "2231",
			"city_id": "257",
			"name": "钟山区"
		},
		{
			"id": "2232",
			"city_id": "257",
			"name": "六枝特区"
		},
		{
			"id": "2233",
			"city_id": "257",
			"name": "水城县"
		},
		{
			"id": "2234",
			"city_id": "257",
			"name": "盘县"
		},
		{
			"id": "2235",
			"city_id": "258",
			"name": "红花岗区"
		},
		{
			"id": "2236",
			"city_id": "258",
			"name": "汇川区"
		},
		{
			"id": "2237",
			"city_id": "258",
			"name": "遵义县"
		},
		{
			"id": "2238",
			"city_id": "258",
			"name": "桐梓县"
		},
		{
			"id": "2239",
			"city_id": "258",
			"name": "绥阳县"
		},
		{
			"id": "2240",
			"city_id": "258",
			"name": "正安县"
		},
		{
			"id": "2241",
			"city_id": "258",
			"name": "道真仡佬族苗族自治县"
		},
		{
			"id": "2242",
			"city_id": "258",
			"name": "务川仡佬族苗族自治县"
		},
		{
			"id": "2243",
			"city_id": "258",
			"name": "凤冈县"
		},
		{
			"id": "2244",
			"city_id": "258",
			"name": "湄潭县"
		},
		{
			"id": "2245",
			"city_id": "258",
			"name": "余庆县"
		},
		{
			"id": "2246",
			"city_id": "258",
			"name": "习水县"
		},
		{
			"id": "2247",
			"city_id": "258",
			"name": "赤水市"
		},
		{
			"id": "2248",
			"city_id": "258",
			"name": "仁怀市"
		},
		{
			"id": "2249",
			"city_id": "259",
			"name": "西秀区"
		},
		{
			"id": "2250",
			"city_id": "259",
			"name": "平坝县"
		},
		{
			"id": "2251",
			"city_id": "259",
			"name": "普定县"
		},
		{
			"id": "2252",
			"city_id": "259",
			"name": "镇宁布依族苗族自治县"
		},
		{
			"id": "2253",
			"city_id": "259",
			"name": "关岭布依族苗族自治县"
		},
		{
			"id": "2254",
			"city_id": "259",
			"name": "紫云苗族布依族自治县"
		},
		{
			"id": "2255",
			"city_id": "260",
			"name": "铜仁市"
		},
		{
			"id": "2256",
			"city_id": "260",
			"name": "江口县"
		},
		{
			"id": "2257",
			"city_id": "260",
			"name": "玉屏侗族自治县"
		},
		{
			"id": "2258",
			"city_id": "260",
			"name": "石阡县"
		},
		{
			"id": "2259",
			"city_id": "260",
			"name": "思南县"
		},
		{
			"id": "2260",
			"city_id": "260",
			"name": "印江土家族苗族自治县"
		},
		{
			"id": "2261",
			"city_id": "260",
			"name": "德江县"
		},
		{
			"id": "2262",
			"city_id": "260",
			"name": "沿河土家族自治县"
		},
		{
			"id": "2263",
			"city_id": "260",
			"name": "松桃苗族自治县"
		},
		{
			"id": "2264",
			"city_id": "260",
			"name": "万山特区"
		},
		{
			"id": "2265",
			"city_id": "261",
			"name": "兴义市"
		},
		{
			"id": "2266",
			"city_id": "261",
			"name": "兴仁县"
		},
		{
			"id": "2267",
			"city_id": "261",
			"name": "普安县"
		},
		{
			"id": "2268",
			"city_id": "261",
			"name": "晴隆县"
		},
		{
			"id": "2269",
			"city_id": "261",
			"name": "贞丰县"
		},
		{
			"id": "2270",
			"city_id": "261",
			"name": "望谟县"
		},
		{
			"id": "2271",
			"city_id": "261",
			"name": "册亨县"
		},
		{
			"id": "2272",
			"city_id": "261",
			"name": "安龙县"
		},
		{
			"id": "2273",
			"city_id": "262",
			"name": "毕节市"
		},
		{
			"id": "2274",
			"city_id": "262",
			"name": "大方县"
		},
		{
			"id": "2275",
			"city_id": "262",
			"name": "黔西县"
		},
		{
			"id": "2276",
			"city_id": "262",
			"name": "金沙县"
		},
		{
			"id": "2277",
			"city_id": "262",
			"name": "织金县"
		},
		{
			"id": "2278",
			"city_id": "262",
			"name": "纳雍县"
		},
		{
			"id": "2279",
			"city_id": "262",
			"name": "威宁彝族回族苗族自治县"
		},
		{
			"id": "2280",
			"city_id": "262",
			"name": "赫章县"
		},
		{
			"id": "2281",
			"city_id": "263",
			"name": "凯里市"
		},
		{
			"id": "2282",
			"city_id": "263",
			"name": "黄平县"
		},
		{
			"id": "2283",
			"city_id": "263",
			"name": "施秉县"
		},
		{
			"id": "2284",
			"city_id": "263",
			"name": "三穗县"
		},
		{
			"id": "2285",
			"city_id": "263",
			"name": "镇远县"
		},
		{
			"id": "2286",
			"city_id": "263",
			"name": "岑巩县"
		},
		{
			"id": "2287",
			"city_id": "263",
			"name": "天柱县"
		},
		{
			"id": "2288",
			"city_id": "263",
			"name": "锦屏县"
		},
		{
			"id": "2289",
			"city_id": "263",
			"name": "剑河县"
		},
		{
			"id": "2290",
			"city_id": "263",
			"name": "台江县"
		},
		{
			"id": "2291",
			"city_id": "263",
			"name": "黎平县"
		},
		{
			"id": "2292",
			"city_id": "263",
			"name": "榕江县"
		},
		{
			"id": "2293",
			"city_id": "263",
			"name": "从江县"
		},
		{
			"id": "2294",
			"city_id": "263",
			"name": "雷山县"
		},
		{
			"id": "2295",
			"city_id": "263",
			"name": "麻江县"
		},
		{
			"id": "2296",
			"city_id": "263",
			"name": "丹寨县"
		},
		{
			"id": "2297",
			"city_id": "264",
			"name": "都匀市"
		},
		{
			"id": "2298",
			"city_id": "264",
			"name": "福泉市"
		},
		{
			"id": "2299",
			"city_id": "264",
			"name": "荔波县"
		},
		{
			"id": "2300",
			"city_id": "264",
			"name": "贵定县"
		},
		{
			"id": "2301",
			"city_id": "264",
			"name": "瓮安县"
		},
		{
			"id": "2302",
			"city_id": "264",
			"name": "独山县"
		},
		{
			"id": "2303",
			"city_id": "264",
			"name": "平塘县"
		},
		{
			"id": "2304",
			"city_id": "264",
			"name": "罗甸县"
		},
		{
			"id": "2305",
			"city_id": "264",
			"name": "长顺县"
		},
		{
			"id": "2306",
			"city_id": "264",
			"name": "龙里县"
		},
		{
			"id": "2307",
			"city_id": "264",
			"name": "惠水县"
		},
		{
			"id": "2308",
			"city_id": "264",
			"name": "三都水族自治县"
		},
		{
			"id": "2309",
			"city_id": "265",
			"name": "五华区"
		},
		{
			"id": "2310",
			"city_id": "265",
			"name": "盘龙区"
		},
		{
			"id": "2311",
			"city_id": "265",
			"name": "官渡区"
		},
		{
			"id": "2312",
			"city_id": "265",
			"name": "西山区"
		},
		{
			"id": "2313",
			"city_id": "265",
			"name": "东川区"
		},
		{
			"id": "2314",
			"city_id": "265",
			"name": "呈贡县"
		},
		{
			"id": "2315",
			"city_id": "265",
			"name": "晋宁县"
		},
		{
			"id": "2316",
			"city_id": "265",
			"name": "富民县"
		},
		{
			"id": "2317",
			"city_id": "265",
			"name": "宜良县"
		},
		{
			"id": "2318",
			"city_id": "265",
			"name": "石林彝族自治县"
		},
		{
			"id": "2319",
			"city_id": "265",
			"name": "嵩明县"
		},
		{
			"id": "2320",
			"city_id": "265",
			"name": "禄劝彝族苗族自治县"
		},
		{
			"id": "2321",
			"city_id": "265",
			"name": "寻甸回族彝族自治县"
		},
		{
			"id": "2322",
			"city_id": "265",
			"name": "安宁市"
		},
		{
			"id": "2323",
			"city_id": "266",
			"name": "麒麟区"
		},
		{
			"id": "2324",
			"city_id": "266",
			"name": "马龙县"
		},
		{
			"id": "2325",
			"city_id": "266",
			"name": "陆良县"
		},
		{
			"id": "2326",
			"city_id": "266",
			"name": "师宗县"
		},
		{
			"id": "2327",
			"city_id": "266",
			"name": "罗平县"
		},
		{
			"id": "2328",
			"city_id": "266",
			"name": "富源县"
		},
		{
			"id": "2329",
			"city_id": "266",
			"name": "会泽县"
		},
		{
			"id": "2330",
			"city_id": "266",
			"name": "沾益县"
		},
		{
			"id": "2331",
			"city_id": "266",
			"name": "宣威市"
		},
		{
			"id": "2332",
			"city_id": "267",
			"name": "红塔区"
		},
		{
			"id": "2333",
			"city_id": "267",
			"name": "江川县"
		},
		{
			"id": "2334",
			"city_id": "267",
			"name": "澄江县"
		},
		{
			"id": "2335",
			"city_id": "267",
			"name": "通海县"
		},
		{
			"id": "2336",
			"city_id": "267",
			"name": "华宁县"
		},
		{
			"id": "2337",
			"city_id": "267",
			"name": "易门县"
		},
		{
			"id": "2338",
			"city_id": "267",
			"name": "峨山彝族自治县"
		},
		{
			"id": "2339",
			"city_id": "267",
			"name": "新平彝族傣族自治县"
		},
		{
			"id": "2340",
			"city_id": "267",
			"name": "元江哈尼族彝族傣族自治县"
		},
		{
			"id": "2341",
			"city_id": "268",
			"name": "隆阳区"
		},
		{
			"id": "2342",
			"city_id": "268",
			"name": "施甸县"
		},
		{
			"id": "2343",
			"city_id": "268",
			"name": "腾冲县"
		},
		{
			"id": "2344",
			"city_id": "268",
			"name": "龙陵县"
		},
		{
			"id": "2345",
			"city_id": "268",
			"name": "昌宁县"
		},
		{
			"id": "2346",
			"city_id": "269",
			"name": "昭阳区"
		},
		{
			"id": "2347",
			"city_id": "269",
			"name": "鲁甸县"
		},
		{
			"id": "2348",
			"city_id": "269",
			"name": "巧家县"
		},
		{
			"id": "2349",
			"city_id": "269",
			"name": "盐津县"
		},
		{
			"id": "2350",
			"city_id": "269",
			"name": "大关县"
		},
		{
			"id": "2351",
			"city_id": "269",
			"name": "永善县"
		},
		{
			"id": "2352",
			"city_id": "269",
			"name": "绥江县"
		},
		{
			"id": "2353",
			"city_id": "269",
			"name": "镇雄县"
		},
		{
			"id": "2354",
			"city_id": "269",
			"name": "彝良县"
		},
		{
			"id": "2355",
			"city_id": "269",
			"name": "威信县"
		},
		{
			"id": "2356",
			"city_id": "269",
			"name": "水富县"
		},
		{
			"id": "2357",
			"city_id": "270",
			"name": "古城区"
		},
		{
			"id": "2358",
			"city_id": "270",
			"name": "玉龙纳西族自治县"
		},
		{
			"id": "2359",
			"city_id": "270",
			"name": "永胜县"
		},
		{
			"id": "2360",
			"city_id": "270",
			"name": "华坪县"
		},
		{
			"id": "2361",
			"city_id": "270",
			"name": "宁蒗彝族自治县"
		},
		{
			"id": "2362",
			"city_id": "271",
			"name": "翠云区"
		},
		{
			"id": "2363",
			"city_id": "271",
			"name": "普洱哈尼族彝族自治县"
		},
		{
			"id": "2364",
			"city_id": "271",
			"name": "墨江哈尼族自治县"
		},
		{
			"id": "2365",
			"city_id": "271",
			"name": "景东彝族自治县"
		},
		{
			"id": "2366",
			"city_id": "271",
			"name": "景谷傣族彝族自治县"
		},
		{
			"id": "2367",
			"city_id": "271",
			"name": "镇沅彝族哈尼族拉祜族自治县"
		},
		{
			"id": "2368",
			"city_id": "271",
			"name": "江城哈尼族彝族自治县"
		},
		{
			"id": "2369",
			"city_id": "271",
			"name": "孟连傣族拉祜族佤族自治县"
		},
		{
			"id": "2370",
			"city_id": "271",
			"name": "澜沧拉祜族自治县"
		},
		{
			"id": "2371",
			"city_id": "271",
			"name": "西盟佤族自治县"
		},
		{
			"id": "2372",
			"city_id": "272",
			"name": "临翔区"
		},
		{
			"id": "2373",
			"city_id": "272",
			"name": "凤庆县"
		},
		{
			"id": "2374",
			"city_id": "272",
			"name": "云县"
		},
		{
			"id": "2375",
			"city_id": "272",
			"name": "永德县"
		},
		{
			"id": "2376",
			"city_id": "272",
			"name": "镇康县"
		},
		{
			"id": "2377",
			"city_id": "272",
			"name": "双江拉祜族佤族布朗族傣族自治县"
		},
		{
			"id": "2378",
			"city_id": "272",
			"name": "耿马傣族佤族自治县"
		},
		{
			"id": "2379",
			"city_id": "272",
			"name": "沧源佤族自治县"
		},
		{
			"id": "2380",
			"city_id": "273",
			"name": "楚雄市"
		},
		{
			"id": "2381",
			"city_id": "273",
			"name": "双柏县"
		},
		{
			"id": "2382",
			"city_id": "273",
			"name": "牟定县"
		},
		{
			"id": "2383",
			"city_id": "273",
			"name": "南华县"
		},
		{
			"id": "2384",
			"city_id": "273",
			"name": "姚安县"
		},
		{
			"id": "2385",
			"city_id": "273",
			"name": "大姚县"
		},
		{
			"id": "2386",
			"city_id": "273",
			"name": "永仁县"
		},
		{
			"id": "2387",
			"city_id": "273",
			"name": "元谋县"
		},
		{
			"id": "2388",
			"city_id": "273",
			"name": "武定县"
		},
		{
			"id": "2389",
			"city_id": "273",
			"name": "禄丰县"
		},
		{
			"id": "2390",
			"city_id": "274",
			"name": "个旧市"
		},
		{
			"id": "2391",
			"city_id": "274",
			"name": "开远市"
		},
		{
			"id": "2392",
			"city_id": "274",
			"name": "蒙自县"
		},
		{
			"id": "2393",
			"city_id": "274",
			"name": "屏边苗族自治县"
		},
		{
			"id": "2394",
			"city_id": "274",
			"name": "建水县"
		},
		{
			"id": "2395",
			"city_id": "274",
			"name": "石屏县"
		},
		{
			"id": "2396",
			"city_id": "274",
			"name": "弥勒县"
		},
		{
			"id": "2397",
			"city_id": "274",
			"name": "泸西县"
		},
		{
			"id": "2398",
			"city_id": "274",
			"name": "元阳县"
		},
		{
			"id": "2399",
			"city_id": "274",
			"name": "红河县"
		},
		{
			"id": "2400",
			"city_id": "274",
			"name": "金平苗族瑶族傣族自治县"
		},
		{
			"id": "2401",
			"city_id": "274",
			"name": "绿春县"
		},
		{
			"id": "2402",
			"city_id": "274",
			"name": "河口瑶族自治县"
		},
		{
			"id": "2403",
			"city_id": "275",
			"name": "文山县"
		},
		{
			"id": "2404",
			"city_id": "275",
			"name": "砚山县"
		},
		{
			"id": "2405",
			"city_id": "275",
			"name": "西畴县"
		},
		{
			"id": "2406",
			"city_id": "275",
			"name": "麻栗坡县"
		},
		{
			"id": "2407",
			"city_id": "275",
			"name": "马关县"
		},
		{
			"id": "2408",
			"city_id": "275",
			"name": "丘北县"
		},
		{
			"id": "2409",
			"city_id": "275",
			"name": "广南县"
		},
		{
			"id": "2410",
			"city_id": "275",
			"name": "富宁县"
		},
		{
			"id": "2411",
			"city_id": "276",
			"name": "景洪市"
		},
		{
			"id": "2412",
			"city_id": "276",
			"name": "勐海县"
		},
		{
			"id": "2413",
			"city_id": "276",
			"name": "勐腊县"
		},
		{
			"id": "2414",
			"city_id": "277",
			"name": "大理市"
		},
		{
			"id": "2415",
			"city_id": "277",
			"name": "漾濞彝族自治县"
		},
		{
			"id": "2416",
			"city_id": "277",
			"name": "祥云县"
		},
		{
			"id": "2417",
			"city_id": "277",
			"name": "宾川县"
		},
		{
			"id": "2418",
			"city_id": "277",
			"name": "弥渡县"
		},
		{
			"id": "2419",
			"city_id": "277",
			"name": "南涧彝族自治县"
		},
		{
			"id": "2420",
			"city_id": "277",
			"name": "巍山彝族回族自治县"
		},
		{
			"id": "2421",
			"city_id": "277",
			"name": "永平县"
		},
		{
			"id": "2422",
			"city_id": "277",
			"name": "云龙县"
		},
		{
			"id": "2423",
			"city_id": "277",
			"name": "洱源县"
		},
		{
			"id": "2424",
			"city_id": "277",
			"name": "剑川县"
		},
		{
			"id": "2425",
			"city_id": "277",
			"name": "鹤庆县"
		},
		{
			"id": "2426",
			"city_id": "278",
			"name": "瑞丽市"
		},
		{
			"id": "2427",
			"city_id": "278",
			"name": "潞西市"
		},
		{
			"id": "2428",
			"city_id": "278",
			"name": "梁河县"
		},
		{
			"id": "2429",
			"city_id": "278",
			"name": "盈江县"
		},
		{
			"id": "2430",
			"city_id": "278",
			"name": "陇川县"
		},
		{
			"id": "2431",
			"city_id": "279",
			"name": "泸水县"
		},
		{
			"id": "2432",
			"city_id": "279",
			"name": "福贡县"
		},
		{
			"id": "2433",
			"city_id": "279",
			"name": "贡山独龙族怒族自治县"
		},
		{
			"id": "2434",
			"city_id": "279",
			"name": "兰坪白族普米族自治县"
		},
		{
			"id": "2435",
			"city_id": "280",
			"name": "香格里拉县"
		},
		{
			"id": "2436",
			"city_id": "280",
			"name": "德钦县"
		},
		{
			"id": "2437",
			"city_id": "280",
			"name": "维西傈僳族自治县"
		},
		{
			"id": "2438",
			"city_id": "281",
			"name": "城关区"
		},
		{
			"id": "2439",
			"city_id": "281",
			"name": "林周县"
		},
		{
			"id": "2440",
			"city_id": "281",
			"name": "当雄县"
		},
		{
			"id": "2441",
			"city_id": "281",
			"name": "尼木县"
		},
		{
			"id": "2442",
			"city_id": "281",
			"name": "曲水县"
		},
		{
			"id": "2443",
			"city_id": "281",
			"name": "堆龙德庆县"
		},
		{
			"id": "2444",
			"city_id": "281",
			"name": "达孜县"
		},
		{
			"id": "2445",
			"city_id": "281",
			"name": "墨竹工卡县"
		},
		{
			"id": "2446",
			"city_id": "282",
			"name": "昌都县"
		},
		{
			"id": "2447",
			"city_id": "282",
			"name": "江达县"
		},
		{
			"id": "2448",
			"city_id": "282",
			"name": "贡觉县"
		},
		{
			"id": "2449",
			"city_id": "282",
			"name": "类乌齐县"
		},
		{
			"id": "2450",
			"city_id": "282",
			"name": "丁青县"
		},
		{
			"id": "2451",
			"city_id": "282",
			"name": "察雅县"
		},
		{
			"id": "2452",
			"city_id": "282",
			"name": "八宿县"
		},
		{
			"id": "2453",
			"city_id": "282",
			"name": "左贡县"
		},
		{
			"id": "2454",
			"city_id": "282",
			"name": "芒康县"
		},
		{
			"id": "2455",
			"city_id": "282",
			"name": "洛隆县"
		},
		{
			"id": "2456",
			"city_id": "282",
			"name": "边坝县"
		},
		{
			"id": "2457",
			"city_id": "283",
			"name": "乃东县"
		},
		{
			"id": "2458",
			"city_id": "283",
			"name": "扎囊县"
		},
		{
			"id": "2459",
			"city_id": "283",
			"name": "贡嘎县"
		},
		{
			"id": "2460",
			"city_id": "283",
			"name": "桑日县"
		},
		{
			"id": "2461",
			"city_id": "283",
			"name": "琼结县"
		},
		{
			"id": "2462",
			"city_id": "283",
			"name": "曲松县"
		},
		{
			"id": "2463",
			"city_id": "283",
			"name": "措美县"
		},
		{
			"id": "2464",
			"city_id": "283",
			"name": "洛扎县"
		},
		{
			"id": "2465",
			"city_id": "283",
			"name": "加查县"
		},
		{
			"id": "2466",
			"city_id": "283",
			"name": "隆子县"
		},
		{
			"id": "2467",
			"city_id": "283",
			"name": "错那县"
		},
		{
			"id": "2468",
			"city_id": "283",
			"name": "浪卡子县"
		},
		{
			"id": "2469",
			"city_id": "284",
			"name": "日喀则市"
		},
		{
			"id": "2470",
			"city_id": "284",
			"name": "南木林县"
		},
		{
			"id": "2471",
			"city_id": "284",
			"name": "江孜县"
		},
		{
			"id": "2472",
			"city_id": "284",
			"name": "定日县"
		},
		{
			"id": "2473",
			"city_id": "284",
			"name": "萨迦县"
		},
		{
			"id": "2474",
			"city_id": "284",
			"name": "拉孜县"
		},
		{
			"id": "2475",
			"city_id": "284",
			"name": "昂仁县"
		},
		{
			"id": "2476",
			"city_id": "284",
			"name": "谢通门县"
		},
		{
			"id": "2477",
			"city_id": "284",
			"name": "白朗县"
		},
		{
			"id": "2478",
			"city_id": "284",
			"name": "仁布县"
		},
		{
			"id": "2479",
			"city_id": "284",
			"name": "康马县"
		},
		{
			"id": "2480",
			"city_id": "284",
			"name": "定结县"
		},
		{
			"id": "2481",
			"city_id": "284",
			"name": "仲巴县"
		},
		{
			"id": "2482",
			"city_id": "284",
			"name": "亚东县"
		},
		{
			"id": "2483",
			"city_id": "284",
			"name": "吉隆县"
		},
		{
			"id": "2484",
			"city_id": "284",
			"name": "聂拉木县"
		},
		{
			"id": "2485",
			"city_id": "284",
			"name": "萨嘎县"
		},
		{
			"id": "2486",
			"city_id": "284",
			"name": "岗巴县"
		},
		{
			"id": "2487",
			"city_id": "285",
			"name": "那曲县"
		},
		{
			"id": "2488",
			"city_id": "285",
			"name": "嘉黎县"
		},
		{
			"id": "2489",
			"city_id": "285",
			"name": "比如县"
		},
		{
			"id": "2490",
			"city_id": "285",
			"name": "聂荣县"
		},
		{
			"id": "2491",
			"city_id": "285",
			"name": "安多县"
		},
		{
			"id": "2492",
			"city_id": "285",
			"name": "申扎县"
		},
		{
			"id": "2493",
			"city_id": "285",
			"name": "索县"
		},
		{
			"id": "2494",
			"city_id": "285",
			"name": "班戈县"
		},
		{
			"id": "2495",
			"city_id": "285",
			"name": "巴青县"
		},
		{
			"id": "2496",
			"city_id": "285",
			"name": "尼玛县"
		},
		{
			"id": "2497",
			"city_id": "286",
			"name": "普兰县"
		},
		{
			"id": "2498",
			"city_id": "286",
			"name": "札达县"
		},
		{
			"id": "2499",
			"city_id": "286",
			"name": "噶尔县"
		},
		{
			"id": "2500",
			"city_id": "286",
			"name": "日土县"
		},
		{
			"id": "2501",
			"city_id": "286",
			"name": "革吉县"
		},
		{
			"id": "2502",
			"city_id": "286",
			"name": "改则县"
		},
		{
			"id": "2503",
			"city_id": "286",
			"name": "措勤县"
		},
		{
			"id": "2504",
			"city_id": "287",
			"name": "林芝县"
		},
		{
			"id": "2505",
			"city_id": "287",
			"name": "工布江达县"
		},
		{
			"id": "2506",
			"city_id": "287",
			"name": "米林县"
		},
		{
			"id": "2507",
			"city_id": "287",
			"name": "墨脱县"
		},
		{
			"id": "2508",
			"city_id": "287",
			"name": "波密县"
		},
		{
			"id": "2509",
			"city_id": "287",
			"name": "察隅县"
		},
		{
			"id": "2510",
			"city_id": "287",
			"name": "朗县"
		},
		{
			"id": "2511",
			"city_id": "288",
			"name": "新城区"
		},
		{
			"id": "2512",
			"city_id": "288",
			"name": "碑林区"
		},
		{
			"id": "2513",
			"city_id": "288",
			"name": "莲湖区"
		},
		{
			"id": "2514",
			"city_id": "288",
			"name": "灞桥区"
		},
		{
			"id": "2515",
			"city_id": "288",
			"name": "未央区"
		},
		{
			"id": "2516",
			"city_id": "288",
			"name": "雁塔区"
		},
		{
			"id": "2517",
			"city_id": "288",
			"name": "阎良区"
		},
		{
			"id": "2518",
			"city_id": "288",
			"name": "临潼区"
		},
		{
			"id": "2519",
			"city_id": "288",
			"name": "长安区"
		},
		{
			"id": "2520",
			"city_id": "288",
			"name": "蓝田县"
		},
		{
			"id": "2521",
			"city_id": "288",
			"name": "周至县"
		},
		{
			"id": "2522",
			"city_id": "288",
			"name": "户县"
		},
		{
			"id": "2523",
			"city_id": "288",
			"name": "高陵县"
		},
		{
			"id": "2524",
			"city_id": "289",
			"name": "王益区"
		},
		{
			"id": "2525",
			"city_id": "289",
			"name": "印台区"
		},
		{
			"id": "2526",
			"city_id": "289",
			"name": "耀州区"
		},
		{
			"id": "2527",
			"city_id": "289",
			"name": "宜君县"
		},
		{
			"id": "2528",
			"city_id": "290",
			"name": "渭滨区"
		},
		{
			"id": "2529",
			"city_id": "290",
			"name": "金台区"
		},
		{
			"id": "2530",
			"city_id": "290",
			"name": "陈仓区"
		},
		{
			"id": "2531",
			"city_id": "290",
			"name": "凤翔县"
		},
		{
			"id": "2532",
			"city_id": "290",
			"name": "岐山县"
		},
		{
			"id": "2533",
			"city_id": "290",
			"name": "扶风县"
		},
		{
			"id": "2534",
			"city_id": "290",
			"name": "眉县"
		},
		{
			"id": "2535",
			"city_id": "290",
			"name": "陇县"
		},
		{
			"id": "2536",
			"city_id": "290",
			"name": "千阳县"
		},
		{
			"id": "2537",
			"city_id": "290",
			"name": "麟游县"
		},
		{
			"id": "2538",
			"city_id": "290",
			"name": "凤县"
		},
		{
			"id": "2539",
			"city_id": "290",
			"name": "太白县"
		},
		{
			"id": "2540",
			"city_id": "291",
			"name": "秦都区"
		},
		{
			"id": "2541",
			"city_id": "291",
			"name": "杨凌区"
		},
		{
			"id": "2542",
			"city_id": "291",
			"name": "渭城区"
		},
		{
			"id": "2543",
			"city_id": "291",
			"name": "三原县"
		},
		{
			"id": "2544",
			"city_id": "291",
			"name": "泾阳县"
		},
		{
			"id": "2545",
			"city_id": "291",
			"name": "乾县"
		},
		{
			"id": "2546",
			"city_id": "291",
			"name": "礼泉县"
		},
		{
			"id": "2547",
			"city_id": "291",
			"name": "永寿县"
		},
		{
			"id": "2548",
			"city_id": "291",
			"name": "彬县"
		},
		{
			"id": "2549",
			"city_id": "291",
			"name": "长武县"
		},
		{
			"id": "2550",
			"city_id": "291",
			"name": "旬邑县"
		},
		{
			"id": "2551",
			"city_id": "291",
			"name": "淳化县"
		},
		{
			"id": "2552",
			"city_id": "291",
			"name": "武功县"
		},
		{
			"id": "2553",
			"city_id": "291",
			"name": "兴平市"
		},
		{
			"id": "2554",
			"city_id": "292",
			"name": "临渭区"
		},
		{
			"id": "2555",
			"city_id": "292",
			"name": "华县"
		},
		{
			"id": "2556",
			"city_id": "292",
			"name": "潼关县"
		},
		{
			"id": "2557",
			"city_id": "292",
			"name": "大荔县"
		},
		{
			"id": "2558",
			"city_id": "292",
			"name": "合阳县"
		},
		{
			"id": "2559",
			"city_id": "292",
			"name": "澄城县"
		},
		{
			"id": "2560",
			"city_id": "292",
			"name": "蒲城县"
		},
		{
			"id": "2561",
			"city_id": "292",
			"name": "白水县"
		},
		{
			"id": "2562",
			"city_id": "292",
			"name": "富平县"
		},
		{
			"id": "2563",
			"city_id": "292",
			"name": "韩城市"
		},
		{
			"id": "2564",
			"city_id": "292",
			"name": "华阴市"
		},
		{
			"id": "2565",
			"city_id": "293",
			"name": "宝塔区"
		},
		{
			"id": "2566",
			"city_id": "293",
			"name": "延长县"
		},
		{
			"id": "2567",
			"city_id": "293",
			"name": "延川县"
		},
		{
			"id": "2568",
			"city_id": "293",
			"name": "子长县"
		},
		{
			"id": "2569",
			"city_id": "293",
			"name": "安塞县"
		},
		{
			"id": "2570",
			"city_id": "293",
			"name": "志丹县"
		},
		{
			"id": "2571",
			"city_id": "293",
			"name": "吴旗县"
		},
		{
			"id": "2572",
			"city_id": "293",
			"name": "甘泉县"
		},
		{
			"id": "2573",
			"city_id": "293",
			"name": "富县"
		},
		{
			"id": "2574",
			"city_id": "293",
			"name": "洛川县"
		},
		{
			"id": "2575",
			"city_id": "293",
			"name": "宜川县"
		},
		{
			"id": "2576",
			"city_id": "293",
			"name": "黄龙县"
		},
		{
			"id": "2577",
			"city_id": "293",
			"name": "黄陵县"
		},
		{
			"id": "2578",
			"city_id": "294",
			"name": "汉台区"
		},
		{
			"id": "2579",
			"city_id": "294",
			"name": "南郑县"
		},
		{
			"id": "2580",
			"city_id": "294",
			"name": "城固县"
		},
		{
			"id": "2581",
			"city_id": "294",
			"name": "洋县"
		},
		{
			"id": "2582",
			"city_id": "294",
			"name": "西乡县"
		},
		{
			"id": "2583",
			"city_id": "294",
			"name": "勉县"
		},
		{
			"id": "2584",
			"city_id": "294",
			"name": "宁强县"
		},
		{
			"id": "2585",
			"city_id": "294",
			"name": "略阳县"
		},
		{
			"id": "2586",
			"city_id": "294",
			"name": "镇巴县"
		},
		{
			"id": "2587",
			"city_id": "294",
			"name": "留坝县"
		},
		{
			"id": "2588",
			"city_id": "294",
			"name": "佛坪县"
		},
		{
			"id": "2589",
			"city_id": "295",
			"name": "榆阳区"
		},
		{
			"id": "2590",
			"city_id": "295",
			"name": "神木县"
		},
		{
			"id": "2591",
			"city_id": "295",
			"name": "府谷县"
		},
		{
			"id": "2592",
			"city_id": "295",
			"name": "横山县"
		},
		{
			"id": "2593",
			"city_id": "295",
			"name": "靖边县"
		},
		{
			"id": "2594",
			"city_id": "295",
			"name": "定边县"
		},
		{
			"id": "2595",
			"city_id": "295",
			"name": "绥德县"
		},
		{
			"id": "2596",
			"city_id": "295",
			"name": "米脂县"
		},
		{
			"id": "2597",
			"city_id": "295",
			"name": "佳县"
		},
		{
			"id": "2598",
			"city_id": "295",
			"name": "吴堡县"
		},
		{
			"id": "2599",
			"city_id": "295",
			"name": "清涧县"
		},
		{
			"id": "2600",
			"city_id": "295",
			"name": "子洲县"
		},
		{
			"id": "2601",
			"city_id": "296",
			"name": "汉滨区"
		},
		{
			"id": "2602",
			"city_id": "296",
			"name": "汉阴县"
		},
		{
			"id": "2603",
			"city_id": "296",
			"name": "石泉县"
		},
		{
			"id": "2604",
			"city_id": "296",
			"name": "宁陕县"
		},
		{
			"id": "2605",
			"city_id": "296",
			"name": "紫阳县"
		},
		{
			"id": "2606",
			"city_id": "296",
			"name": "岚皋县"
		},
		{
			"id": "2607",
			"city_id": "296",
			"name": "平利县"
		},
		{
			"id": "2608",
			"city_id": "296",
			"name": "镇坪县"
		},
		{
			"id": "2609",
			"city_id": "296",
			"name": "旬阳县"
		},
		{
			"id": "2610",
			"city_id": "296",
			"name": "白河县"
		},
		{
			"id": "2611",
			"city_id": "297",
			"name": "商州区"
		},
		{
			"id": "2612",
			"city_id": "297",
			"name": "洛南县"
		},
		{
			"id": "2613",
			"city_id": "297",
			"name": "丹凤县"
		},
		{
			"id": "2614",
			"city_id": "297",
			"name": "商南县"
		},
		{
			"id": "2615",
			"city_id": "297",
			"name": "山阳县"
		},
		{
			"id": "2616",
			"city_id": "297",
			"name": "镇安县"
		},
		{
			"id": "2617",
			"city_id": "297",
			"name": "柞水县"
		},
		{
			"id": "2618",
			"city_id": "298",
			"name": "城关区"
		},
		{
			"id": "2619",
			"city_id": "298",
			"name": "七里河区"
		},
		{
			"id": "2620",
			"city_id": "298",
			"name": "西固区"
		},
		{
			"id": "2621",
			"city_id": "298",
			"name": "安宁区"
		},
		{
			"id": "2622",
			"city_id": "298",
			"name": "红古区"
		},
		{
			"id": "2623",
			"city_id": "298",
			"name": "永登县"
		},
		{
			"id": "2624",
			"city_id": "298",
			"name": "皋兰县"
		},
		{
			"id": "2625",
			"city_id": "298",
			"name": "榆中县"
		},
		{
			"id": "2626",
			"city_id": "300",
			"name": "金川区"
		},
		{
			"id": "2627",
			"city_id": "300",
			"name": "永昌县"
		},
		{
			"id": "2628",
			"city_id": "301",
			"name": "白银区"
		},
		{
			"id": "2629",
			"city_id": "301",
			"name": "平川区"
		},
		{
			"id": "2630",
			"city_id": "301",
			"name": "靖远县"
		},
		{
			"id": "2631",
			"city_id": "301",
			"name": "会宁县"
		},
		{
			"id": "2632",
			"city_id": "301",
			"name": "景泰县"
		},
		{
			"id": "2633",
			"city_id": "302",
			"name": "秦城区"
		},
		{
			"id": "2634",
			"city_id": "302",
			"name": "北道区"
		},
		{
			"id": "2635",
			"city_id": "302",
			"name": "清水县"
		},
		{
			"id": "2636",
			"city_id": "302",
			"name": "秦安县"
		},
		{
			"id": "2637",
			"city_id": "302",
			"name": "甘谷县"
		},
		{
			"id": "2638",
			"city_id": "302",
			"name": "武山县"
		},
		{
			"id": "2639",
			"city_id": "302",
			"name": "张家川回族自治县"
		},
		{
			"id": "2640",
			"city_id": "303",
			"name": "凉州区"
		},
		{
			"id": "2641",
			"city_id": "303",
			"name": "民勤县"
		},
		{
			"id": "2642",
			"city_id": "303",
			"name": "古浪县"
		},
		{
			"id": "2643",
			"city_id": "303",
			"name": "天祝藏族自治县"
		},
		{
			"id": "2644",
			"city_id": "304",
			"name": "甘州区"
		},
		{
			"id": "2645",
			"city_id": "304",
			"name": "肃南裕固族自治县"
		},
		{
			"id": "2646",
			"city_id": "304",
			"name": "民乐县"
		},
		{
			"id": "2647",
			"city_id": "304",
			"name": "临泽县"
		},
		{
			"id": "2648",
			"city_id": "304",
			"name": "高台县"
		},
		{
			"id": "2649",
			"city_id": "304",
			"name": "山丹县"
		},
		{
			"id": "2650",
			"city_id": "305",
			"name": "崆峒区"
		},
		{
			"id": "2651",
			"city_id": "305",
			"name": "泾川县"
		},
		{
			"id": "2652",
			"city_id": "305",
			"name": "灵台县"
		},
		{
			"id": "2653",
			"city_id": "305",
			"name": "崇信县"
		},
		{
			"id": "2654",
			"city_id": "305",
			"name": "华亭县"
		},
		{
			"id": "2655",
			"city_id": "305",
			"name": "庄浪县"
		},
		{
			"id": "2656",
			"city_id": "305",
			"name": "静宁县"
		},
		{
			"id": "2657",
			"city_id": "306",
			"name": "肃州区"
		},
		{
			"id": "2658",
			"city_id": "306",
			"name": "金塔县"
		},
		{
			"id": "2659",
			"city_id": "306",
			"name": "安西县"
		},
		{
			"id": "2660",
			"city_id": "306",
			"name": "肃北蒙古族自治县"
		},
		{
			"id": "2661",
			"city_id": "306",
			"name": "阿克塞哈萨克族自治县"
		},
		{
			"id": "2662",
			"city_id": "306",
			"name": "玉门市"
		},
		{
			"id": "2663",
			"city_id": "306",
			"name": "敦煌市"
		},
		{
			"id": "2664",
			"city_id": "307",
			"name": "西峰区"
		},
		{
			"id": "2665",
			"city_id": "307",
			"name": "庆城县"
		},
		{
			"id": "2666",
			"city_id": "307",
			"name": "环县"
		},
		{
			"id": "2667",
			"city_id": "307",
			"name": "华池县"
		},
		{
			"id": "2668",
			"city_id": "307",
			"name": "合水县"
		},
		{
			"id": "2669",
			"city_id": "307",
			"name": "正宁县"
		},
		{
			"id": "2670",
			"city_id": "307",
			"name": "宁县"
		},
		{
			"id": "2671",
			"city_id": "307",
			"name": "镇原县"
		},
		{
			"id": "2672",
			"city_id": "308",
			"name": "安定区"
		},
		{
			"id": "2673",
			"city_id": "308",
			"name": "通渭县"
		},
		{
			"id": "2674",
			"city_id": "308",
			"name": "陇西县"
		},
		{
			"id": "2675",
			"city_id": "308",
			"name": "渭源县"
		},
		{
			"id": "2676",
			"city_id": "308",
			"name": "临洮县"
		},
		{
			"id": "2677",
			"city_id": "308",
			"name": "漳县"
		},
		{
			"id": "2678",
			"city_id": "308",
			"name": "岷县"
		},
		{
			"id": "2679",
			"city_id": "309",
			"name": "武都区"
		},
		{
			"id": "2680",
			"city_id": "309",
			"name": "成县"
		},
		{
			"id": "2681",
			"city_id": "309",
			"name": "文县"
		},
		{
			"id": "2682",
			"city_id": "309",
			"name": "宕昌县"
		},
		{
			"id": "2683",
			"city_id": "309",
			"name": "康县"
		},
		{
			"id": "2684",
			"city_id": "309",
			"name": "西和县"
		},
		{
			"id": "2685",
			"city_id": "309",
			"name": "礼县"
		},
		{
			"id": "2686",
			"city_id": "309",
			"name": "徽县"
		},
		{
			"id": "2687",
			"city_id": "309",
			"name": "两当县"
		},
		{
			"id": "2688",
			"city_id": "310",
			"name": "临夏市"
		},
		{
			"id": "2689",
			"city_id": "310",
			"name": "临夏县"
		},
		{
			"id": "2690",
			"city_id": "310",
			"name": "康乐县"
		},
		{
			"id": "2691",
			"city_id": "310",
			"name": "永靖县"
		},
		{
			"id": "2692",
			"city_id": "310",
			"name": "广河县"
		},
		{
			"id": "2693",
			"city_id": "310",
			"name": "和政县"
		},
		{
			"id": "2694",
			"city_id": "310",
			"name": "东乡族自治县"
		},
		{
			"id": "2695",
			"city_id": "310",
			"name": "积石山保安族东乡族撒拉族自治县"
		},
		{
			"id": "2696",
			"city_id": "311",
			"name": "合作市"
		},
		{
			"id": "2697",
			"city_id": "311",
			"name": "临潭县"
		},
		{
			"id": "2698",
			"city_id": "311",
			"name": "卓尼县"
		},
		{
			"id": "2699",
			"city_id": "311",
			"name": "舟曲县"
		},
		{
			"id": "2700",
			"city_id": "311",
			"name": "迭部县"
		},
		{
			"id": "2701",
			"city_id": "311",
			"name": "玛曲县"
		},
		{
			"id": "2702",
			"city_id": "311",
			"name": "碌曲县"
		},
		{
			"id": "2703",
			"city_id": "311",
			"name": "夏河县"
		},
		{
			"id": "2704",
			"city_id": "312",
			"name": "城东区"
		},
		{
			"id": "2705",
			"city_id": "312",
			"name": "城中区"
		},
		{
			"id": "2706",
			"city_id": "312",
			"name": "城西区"
		},
		{
			"id": "2707",
			"city_id": "312",
			"name": "城北区"
		},
		{
			"id": "2708",
			"city_id": "312",
			"name": "大通回族土族自治县"
		},
		{
			"id": "2709",
			"city_id": "312",
			"name": "湟中县"
		},
		{
			"id": "2710",
			"city_id": "312",
			"name": "湟源县"
		},
		{
			"id": "2711",
			"city_id": "313",
			"name": "平安县"
		},
		{
			"id": "2712",
			"city_id": "313",
			"name": "民和回族土族自治县"
		},
		{
			"id": "2713",
			"city_id": "313",
			"name": "乐都县"
		},
		{
			"id": "2714",
			"city_id": "313",
			"name": "互助土族自治县"
		},
		{
			"id": "2715",
			"city_id": "313",
			"name": "化隆回族自治县"
		},
		{
			"id": "2716",
			"city_id": "313",
			"name": "循化撒拉族自治县"
		},
		{
			"id": "2717",
			"city_id": "314",
			"name": "门源回族自治县"
		},
		{
			"id": "2718",
			"city_id": "314",
			"name": "祁连县"
		},
		{
			"id": "2719",
			"city_id": "314",
			"name": "海晏县"
		},
		{
			"id": "2720",
			"city_id": "314",
			"name": "刚察县"
		},
		{
			"id": "2721",
			"city_id": "315",
			"name": "同仁县"
		},
		{
			"id": "2722",
			"city_id": "315",
			"name": "尖扎县"
		},
		{
			"id": "2723",
			"city_id": "315",
			"name": "泽库县"
		},
		{
			"id": "2724",
			"city_id": "315",
			"name": "河南蒙古族自治县"
		},
		{
			"id": "2725",
			"city_id": "316",
			"name": "共和县"
		},
		{
			"id": "2726",
			"city_id": "316",
			"name": "同德县"
		},
		{
			"id": "2727",
			"city_id": "316",
			"name": "贵德县"
		},
		{
			"id": "2728",
			"city_id": "316",
			"name": "兴海县"
		},
		{
			"id": "2729",
			"city_id": "316",
			"name": "贵南县"
		},
		{
			"id": "2730",
			"city_id": "317",
			"name": "玛沁县"
		},
		{
			"id": "2731",
			"city_id": "317",
			"name": "班玛县"
		},
		{
			"id": "2732",
			"city_id": "317",
			"name": "甘德县"
		},
		{
			"id": "2733",
			"city_id": "317",
			"name": "达日县"
		},
		{
			"id": "2734",
			"city_id": "317",
			"name": "久治县"
		},
		{
			"id": "2735",
			"city_id": "317",
			"name": "玛多县"
		},
		{
			"id": "2736",
			"city_id": "318",
			"name": "玉树县"
		},
		{
			"id": "2737",
			"city_id": "318",
			"name": "杂多县"
		},
		{
			"id": "2738",
			"city_id": "318",
			"name": "称多县"
		},
		{
			"id": "2739",
			"city_id": "318",
			"name": "治多县"
		},
		{
			"id": "2740",
			"city_id": "318",
			"name": "囊谦县"
		},
		{
			"id": "2741",
			"city_id": "318",
			"name": "曲麻莱县"
		},
		{
			"id": "2742",
			"city_id": "319",
			"name": "格尔木市"
		},
		{
			"id": "2743",
			"city_id": "319",
			"name": "德令哈市"
		},
		{
			"id": "2744",
			"city_id": "319",
			"name": "乌兰县"
		},
		{
			"id": "2745",
			"city_id": "319",
			"name": "都兰县"
		},
		{
			"id": "2746",
			"city_id": "319",
			"name": "天峻县"
		},
		{
			"id": "2747",
			"city_id": "320",
			"name": "兴庆区"
		},
		{
			"id": "2748",
			"city_id": "320",
			"name": "西夏区"
		},
		{
			"id": "2749",
			"city_id": "320",
			"name": "金凤区"
		},
		{
			"id": "2750",
			"city_id": "320",
			"name": "永宁县"
		},
		{
			"id": "2751",
			"city_id": "320",
			"name": "贺兰县"
		},
		{
			"id": "2752",
			"city_id": "320",
			"name": "灵武市"
		},
		{
			"id": "2753",
			"city_id": "321",
			"name": "大武口区"
		},
		{
			"id": "2754",
			"city_id": "321",
			"name": "惠农区"
		},
		{
			"id": "2755",
			"city_id": "321",
			"name": "平罗县"
		},
		{
			"id": "2756",
			"city_id": "322",
			"name": "利通区"
		},
		{
			"id": "2757",
			"city_id": "322",
			"name": "盐池县"
		},
		{
			"id": "2758",
			"city_id": "322",
			"name": "同心县"
		},
		{
			"id": "2759",
			"city_id": "322",
			"name": "青铜峡市"
		},
		{
			"id": "2760",
			"city_id": "323",
			"name": "原州区"
		},
		{
			"id": "2761",
			"city_id": "323",
			"name": "西吉县"
		},
		{
			"id": "2762",
			"city_id": "323",
			"name": "隆德县"
		},
		{
			"id": "2763",
			"city_id": "323",
			"name": "泾源县"
		},
		{
			"id": "2764",
			"city_id": "323",
			"name": "彭阳县"
		},
		{
			"id": "2765",
			"city_id": "324",
			"name": "沙坡头区"
		},
		{
			"id": "2766",
			"city_id": "324",
			"name": "中宁县"
		},
		{
			"id": "2767",
			"city_id": "324",
			"name": "海原县"
		},
		{
			"id": "2768",
			"city_id": "325",
			"name": "天山区"
		},
		{
			"id": "2769",
			"city_id": "325",
			"name": "沙依巴克区"
		},
		{
			"id": "2770",
			"city_id": "325",
			"name": "新市区"
		},
		{
			"id": "2771",
			"city_id": "325",
			"name": "水磨沟区"
		},
		{
			"id": "2772",
			"city_id": "325",
			"name": "头屯河区"
		},
		{
			"id": "2773",
			"city_id": "325",
			"name": "达坂城区"
		},
		{
			"id": "2774",
			"city_id": "325",
			"name": "东山区"
		},
		{
			"id": "2775",
			"city_id": "325",
			"name": "乌鲁木齐县"
		},
		{
			"id": "2776",
			"city_id": "326",
			"name": "独山子区"
		},
		{
			"id": "2777",
			"city_id": "326",
			"name": "克拉玛依区"
		},
		{
			"id": "2778",
			"city_id": "326",
			"name": "白碱滩区"
		},
		{
			"id": "2779",
			"city_id": "326",
			"name": "乌尔禾区"
		},
		{
			"id": "2780",
			"city_id": "327",
			"name": "吐鲁番市"
		},
		{
			"id": "2781",
			"city_id": "327",
			"name": "鄯善县"
		},
		{
			"id": "2782",
			"city_id": "327",
			"name": "托克逊县"
		},
		{
			"id": "2783",
			"city_id": "328",
			"name": "哈密市"
		},
		{
			"id": "2784",
			"city_id": "328",
			"name": "巴里坤哈萨克自治县"
		},
		{
			"id": "2785",
			"city_id": "328",
			"name": "伊吾县"
		},
		{
			"id": "2786",
			"city_id": "329",
			"name": "昌吉市"
		},
		{
			"id": "2787",
			"city_id": "329",
			"name": "阜康市"
		},
		{
			"id": "2788",
			"city_id": "329",
			"name": "米泉市"
		},
		{
			"id": "2789",
			"city_id": "329",
			"name": "呼图壁县"
		},
		{
			"id": "2790",
			"city_id": "329",
			"name": "玛纳斯县"
		},
		{
			"id": "2791",
			"city_id": "329",
			"name": "奇台县"
		},
		{
			"id": "2792",
			"city_id": "329",
			"name": "吉木萨尔县"
		},
		{
			"id": "2793",
			"city_id": "329",
			"name": "木垒哈萨克自治县"
		},
		{
			"id": "2794",
			"city_id": "330",
			"name": "博乐市"
		},
		{
			"id": "2795",
			"city_id": "330",
			"name": "精河县"
		},
		{
			"id": "2796",
			"city_id": "330",
			"name": "温泉县"
		},
		{
			"id": "2797",
			"city_id": "331",
			"name": "库尔勒市"
		},
		{
			"id": "2798",
			"city_id": "331",
			"name": "轮台县"
		},
		{
			"id": "2799",
			"city_id": "331",
			"name": "尉犁县"
		},
		{
			"id": "2800",
			"city_id": "331",
			"name": "若羌县"
		},
		{
			"id": "2801",
			"city_id": "331",
			"name": "且末县"
		},
		{
			"id": "2802",
			"city_id": "331",
			"name": "焉耆回族自治县"
		},
		{
			"id": "2803",
			"city_id": "331",
			"name": "和静县"
		},
		{
			"id": "2804",
			"city_id": "331",
			"name": "和硕县"
		},
		{
			"id": "2805",
			"city_id": "331",
			"name": "博湖县"
		},
		{
			"id": "2806",
			"city_id": "332",
			"name": "阿克苏市"
		},
		{
			"id": "2807",
			"city_id": "332",
			"name": "温宿县"
		},
		{
			"id": "2808",
			"city_id": "332",
			"name": "库车县"
		},
		{
			"id": "2809",
			"city_id": "332",
			"name": "沙雅县"
		},
		{
			"id": "2810",
			"city_id": "332",
			"name": "新和县"
		},
		{
			"id": "2811",
			"city_id": "332",
			"name": "拜城县"
		},
		{
			"id": "2812",
			"city_id": "332",
			"name": "乌什县"
		},
		{
			"id": "2813",
			"city_id": "332",
			"name": "阿瓦提县"
		},
		{
			"id": "2814",
			"city_id": "332",
			"name": "柯坪县"
		},
		{
			"id": "2815",
			"city_id": "333",
			"name": "阿图什市"
		},
		{
			"id": "2816",
			"city_id": "333",
			"name": "阿克陶县"
		},
		{
			"id": "2817",
			"city_id": "333",
			"name": "阿合奇县"
		},
		{
			"id": "2818",
			"city_id": "333",
			"name": "乌恰县"
		},
		{
			"id": "2819",
			"city_id": "334",
			"name": "喀什市"
		},
		{
			"id": "2820",
			"city_id": "334",
			"name": "疏附县"
		},
		{
			"id": "2821",
			"city_id": "334",
			"name": "疏勒县"
		},
		{
			"id": "2822",
			"city_id": "334",
			"name": "英吉沙县"
		},
		{
			"id": "2823",
			"city_id": "334",
			"name": "泽普县"
		},
		{
			"id": "2824",
			"city_id": "334",
			"name": "莎车县"
		},
		{
			"id": "2825",
			"city_id": "334",
			"name": "叶城县"
		},
		{
			"id": "2826",
			"city_id": "334",
			"name": "麦盖提县"
		},
		{
			"id": "2827",
			"city_id": "334",
			"name": "岳普湖县"
		},
		{
			"id": "2828",
			"city_id": "334",
			"name": "伽师县"
		},
		{
			"id": "2829",
			"city_id": "334",
			"name": "巴楚县"
		},
		{
			"id": "2830",
			"city_id": "334",
			"name": "塔什库尔干塔吉克自治县"
		},
		{
			"id": "2831",
			"city_id": "335",
			"name": "和田市"
		},
		{
			"id": "2832",
			"city_id": "335",
			"name": "和田县"
		},
		{
			"id": "2833",
			"city_id": "335",
			"name": "墨玉县"
		},
		{
			"id": "2834",
			"city_id": "335",
			"name": "皮山县"
		},
		{
			"id": "2835",
			"city_id": "335",
			"name": "洛浦县"
		},
		{
			"id": "2836",
			"city_id": "335",
			"name": "策勒县"
		},
		{
			"id": "2837",
			"city_id": "335",
			"name": "于田县"
		},
		{
			"id": "2838",
			"city_id": "335",
			"name": "民丰县"
		},
		{
			"id": "2839",
			"city_id": "336",
			"name": "伊宁市"
		},
		{
			"id": "2840",
			"city_id": "336",
			"name": "奎屯市"
		},
		{
			"id": "2841",
			"city_id": "336",
			"name": "伊宁县"
		},
		{
			"id": "2842",
			"city_id": "336",
			"name": "察布查尔锡伯自治县"
		},
		{
			"id": "2843",
			"city_id": "336",
			"name": "霍城县"
		},
		{
			"id": "2844",
			"city_id": "336",
			"name": "巩留县"
		},
		{
			"id": "2845",
			"city_id": "336",
			"name": "新源县"
		},
		{
			"id": "2846",
			"city_id": "336",
			"name": "昭苏县"
		},
		{
			"id": "2847",
			"city_id": "336",
			"name": "特克斯县"
		},
		{
			"id": "2848",
			"city_id": "336",
			"name": "尼勒克县"
		},
		{
			"id": "2849",
			"city_id": "337",
			"name": "塔城市"
		},
		{
			"id": "2850",
			"city_id": "337",
			"name": "乌苏市"
		},
		{
			"id": "2851",
			"city_id": "337",
			"name": "额敏县"
		},
		{
			"id": "2852",
			"city_id": "337",
			"name": "沙湾县"
		},
		{
			"id": "2853",
			"city_id": "337",
			"name": "托里县"
		},
		{
			"id": "2854",
			"city_id": "337",
			"name": "裕民县"
		},
		{
			"id": "2855",
			"city_id": "337",
			"name": "和布克赛尔蒙古自治县"
		},
		{
			"id": "2856",
			"city_id": "338",
			"name": "阿勒泰市"
		},
		{
			"id": "2857",
			"city_id": "338",
			"name": "布尔津县"
		},
		{
			"id": "2858",
			"city_id": "338",
			"name": "富蕴县"
		},
		{
			"id": "2859",
			"city_id": "338",
			"name": "福海县"
		},
		{
			"id": "2860",
			"city_id": "338",
			"name": "哈巴河县"
		},
		{
			"id": "2861",
			"city_id": "338",
			"name": "青河县"
		},
		{
			"id": "2862",
			"city_id": "338",
			"name": "吉木乃县"
		},
		{
			"id": "2863",
			"city_id": "214",
			"name": "中山市"
		},
		{
			"id": "2864",
			"city_id": "75",
			"name": "新吴区"
		},
		{
			"id": "2865",
			"city_id": "75",
			"name": "梁溪区"
		},
		{
			"id": "2866",
			"city_id": "213",
			"name": "东莞市"
		},
		{
			"id": "2867",
			"city_id": "343",
			"name": "香港特别行政区"
		},
		{
			"id": "2868",
			"city_id": "344",
			"name": "澳门特别行政区"
		},
		{
			"id": "2869",
			"city_id": "345",
			"name": "台湾省"
		},
		{
			"id": "2870",
			"city_id": "78",
			"name": "姑苏区"
		},
		{
			"id": "2871",
			"city_id": "233",
			"name": "天涯区"
		},
		{
			"id": "2872",
			"city_id": "241",
			"name": "利州区"
		},
		{
			"id": "2873",
			"city_id": "78",
			"name": "工业园区"
		}
	]
}`