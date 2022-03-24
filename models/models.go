package models

import (
	"encoding/json"
	"etl-blocks2/db"
	"etl-blocks2/dbRedis"
	"fmt"
	"log"
	"math"
	"strconv"

	geojson "github.com/paulmach/go.geojson"
)

type BQData struct {
	Block_id      int              `json:"block_id"`
	ClientBQ_id   int              `json:"client_id"`
	Block_parent  int              `json:"block_parent"`
	Block_name    string           `json:"block_name"`
	Block_bounds  geojson.Geometry `json:"bounds"`
	Block_abrv    string           `json:"abvr"`
	Centroid      geojson.Geometry `json:"centroid"`
	Centroid_text string           `json:"centroid_text"`
}

type Client struct {
	Client_id   int
	Client_name string
}

type Mapping struct {
	Client_id  int
	Mapping_id int
}

type Areas struct {
	Client_id       int
	Areas_id        int
	Areas_name      string
	Areas_branch_id int
	Bounds          geojson.Geometry
	Centroid        geojson.Geometry
	Centroid_text   string
	Block_abrv      string
}

type Branch struct {
	Client_id    int
	Branch_id    int
	Branch_name  string
	Area_farm_id int
	Bounds       float64
	Block_abrv   string
}

type Farm struct {
	Client_id  int
	Farm_id    int
	Farm_name  string
	Block_abrv string
}

type TableId struct {
	Block_id int
	Type_id  int
	Type     string
}

func Polygon_GetArea(geoareas []geojson.Geometry) *geojson.Geometry {

	min_lat := 0.0
	max_lat := 0.0
	min_long := 0.0
	max_long := 0.0

	n := 0
	for _, geo := range geoareas {
		polygon := geo.Polygon
		for _, elements := range polygon {
			for _, coord := range elements {
				if n == 0 {
					min_long = coord[0]
					max_long = coord[0]
					min_lat = coord[1]
					max_lat = coord[1]
				} else {
					min_long = math.Min(min_long, coord[0])
					max_long = math.Max(max_long, coord[0])
					min_lat = math.Min(min_lat, coord[1])
					max_lat = math.Max(max_lat, coord[1])
				}
				n++
			}
		}
	}
	a_point := fmt.Sprintf("[%f,%f]", min_long, min_lat)
	b_point := fmt.Sprintf("[%f,%f]", max_long, min_lat)
	c_point := fmt.Sprintf("[%f,%f]", min_long, max_lat)
	d_point := fmt.Sprintf("[%f,%f]", max_long, max_lat)
	pol := fmt.Sprintf("[[ %s, %s, %s, %s ]]", a_point, b_point, c_point, d_point)

	rawAreaJSON := []byte(fmt.Sprintf(`{ "type": "Polygon", "coordinates": %s}`, pol))

	polygon_area, err := geojson.UnmarshalGeometry(rawAreaJSON)
	if err != nil {
		fmt.Println("Geo error: ", err)
		return nil
	}
	return polygon_area
}

func ReadPg() []BQData {

	// Criando hashtables para armazenar os blocos
	client_blks := make(map[int]int)
	client_bqs := make(map[int]int)
	branch_blks := make(map[int]int)
	branch_bqs := make(map[int]int)
	farm_blks := make(map[int]int)
	farm_bqs := make(map[int]int)

	block_id := 1

	db := db.DbConect()
	defer db.Close()

	client := Client{}
	mapping := Mapping{}
	areas := Areas{}
	branch := Branch{}
	farm := Farm{}

	bqData := BQData{}
	var bqDataArr []BQData

	var block_id_cli int
	var block_id_farm int
	var block_id_branch int

	var poli_areas []geojson.Geometry
	var poli_branches []geojson.Geometry
	var poli_farms []geojson.Geometry

	table_id := TableId{}
	var table_id_arr []TableId

	//client
	queryClient := fmt.Sprintf(
		"select distinct cli.client_id, cli.name as client_name " +
			"from client cli " +
			"left join mapping_areas ma on ma.client_id = cli.client_id " +
			"where cli.situation_id in (1, 2, 7) and ma.id is not null and cli.client_id in (163, 16, 19, 812);")

	client_rows, err := db.Query(queryClient)
	if err != nil {
		log.Println("Error:", err.Error())
	}

	for client_rows.Next() {

		err := client_rows.Scan(&client.Client_id, &client.Client_name)
		if err != nil {
			log.Println("Error:", err.Error())
		}

		client_blks[client.Client_id] = block_id
		block_id_cli = block_id
		block_id++

		//mapping
		queryMapping := fmt.Sprintf(
			"select ma.client_id, ma.id " +
				"from mapping_areas ma " +
				"where ma.client_id = " + strconv.Itoa(client.Client_id))

		mapping_rows, err := db.Query(queryMapping)
		if err != nil {
			log.Println("Error:", err.Error())
		}

		poli_areas = nil
		for mapping_rows.Next() {

			err := mapping_rows.Scan(&mapping.Client_id, &mapping.Mapping_id)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			//areas
			queryAreas := fmt.Sprintf(
				"select ar.name as areas_name,  " + strconv.Itoa(mapping.Client_id) + "as Client_id, ar.id, ar.area_branch_id, ST_AsGeoJSON(ar.bounds) as bounds, " +
					"ST_AsGeoJSON(ST_Centroid(ar.bounds)) as centroid, st_astext(ST_Centroid(bounds)) as centroid_text " +
					"from areas ar " +
					"where mapping_area_id = " + strconv.Itoa(mapping.Mapping_id))

			areas_rows, err := db.Query(queryAreas)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			for areas_rows.Next() {

				err := areas_rows.Scan(&areas.Areas_name, &areas.Client_id, &areas.Areas_id, &areas.Areas_branch_id, &areas.Bounds, &areas.Centroid, &areas.Centroid_text)
				if err != nil {
					log.Println("Error:", err.Error())
				}

				//branch
				queryBranch := fmt.Sprintf(
					"select " + strconv.Itoa(mapping.Client_id) + " as Client_id , " + "b.area_branch_id as branch_id, b.name as branch_name, b.area_farm_id " +
						"from area_branch b " +
						"where b.area_branch_id = " + strconv.Itoa(areas.Areas_branch_id))

				branch_rows, err := db.Query(queryBranch)
				if err != nil {
					log.Println("Error:", err.Error())
				}

				for branch_rows.Next() {

					err := branch_rows.Scan(&branch.Client_id, &branch.Branch_id, &branch.Branch_name, &branch.Area_farm_id)
					if err != nil {
						log.Println("Error:", err.Error())
					}

					//farm
					queryFarm := fmt.Sprintf(
						"select distinct " + strconv.Itoa(mapping.Client_id) + " as Client_id , f.area_farm_id, " + "f.name as farm_name " +
							"from area_farm f " +
							"where f.area_farm_id = " + strconv.Itoa(branch.Area_farm_id))

					farm_rows, err := db.Query(queryFarm)
					if err != nil {
						log.Println("Error:", err.Error())
					}

					for farm_rows.Next() {

						err := farm_rows.Scan(&farm.Client_id, &farm.Farm_id, &farm.Farm_name)
						if err != nil {
							log.Println("Error:", err.Error())
						}

						abrv_farm := farm.Farm_name[0:3]

						if farm_blks[farm.Farm_id] == 0 {
							poli_farms = nil
							poli_farms = append(poli_farms, areas.Bounds)
							bounds_farms := *Polygon_GetArea(poli_farms)
							centroid := Centroid(bounds_farms.Polygon)
							centroidText := CentroidText(centroid)
							bqData = BQData{
								Block_id:      block_id,
								ClientBQ_id:   farm.Client_id,
								Block_parent:  block_id_cli,
								Block_name:    farm.Farm_name,
								Block_bounds:  *Polygon_GetArea(poli_farms),
								Block_abrv:    abrv_farm,
								Centroid:      *centroid,
								Centroid_text: centroidText,
							}
							bqDataArr = append(bqDataArr, bqData)

							//tabela separada
							table_id = TableId{block_id, farm.Farm_id, "Farm"}
							table_id_arr = append(table_id_arr, table_id)
							//
							farm_blks[farm.Farm_id] = block_id
							farm_bqs[farm.Farm_id] = len(bqDataArr) - 1
							block_id_farm = block_id
							block_id++
						} else {
							block_id_farm = farm_blks[farm.Farm_id]
							poli_farms = nil
							poli_farms = append(poli_farms, bqDataArr[farm_bqs[farm.Farm_id]].Block_bounds)
							poli_farms = append(poli_farms, areas.Bounds)
							bounds_farms := *Polygon_GetArea(poli_farms)
							centroid := Centroid(bounds_farms.Polygon)
							centroidText := CentroidText(centroid)
							bqDataArr[farm_bqs[farm.Farm_id]].Block_bounds = bounds_farms
							bqDataArr[farm_bqs[farm.Farm_id]].Centroid = *centroid
							bqDataArr[farm_bqs[farm.Farm_id]].Centroid_text = centroidText
						}

					}

					var abrv_branch string
					if len(branch.Branch_name) > 1 && len(branch.Branch_name) < 3 {
						abrv_branch = branch.Branch_name[0:2]
					} else if len(branch.Branch_name) < 2 {
						abrv_branch = branch.Branch_name[0:1]
					} else {
						abrv_branch = branch.Branch_name[0:3]
					}

					if branch_blks[branch.Branch_id] == 0 {
						poli_branches = nil
						poli_branches = append(poli_branches, areas.Bounds)
						bounds_branches := *Polygon_GetArea(poli_branches)
						centroid := Centroid(bounds_branches.Polygon)
						centroidText := CentroidText(centroid)
						bqData = BQData{
							Block_id:      block_id,
							ClientBQ_id:   branch.Client_id,
							Block_parent:  block_id_farm,
							Block_name:    branch.Branch_name,
							Block_bounds:  *Polygon_GetArea(poli_branches),
							Block_abrv:    abrv_branch,
							Centroid:      *centroid,
							Centroid_text: centroidText}
						bqDataArr = append(bqDataArr, bqData)

						//tabela separada
						table_id = TableId{block_id, branch.Branch_id, "Branch"}
						table_id_arr = append(table_id_arr, table_id)
						//
						branch_blks[branch.Branch_id] = block_id
						branch_bqs[branch.Branch_id] = len(bqDataArr) - 1
						block_id_branch = block_id
						block_id++
					} else {
						block_id_branch = branch_blks[branch.Branch_id]
						poli_branches = nil
						poli_branches = append(poli_branches, bqDataArr[branch_bqs[branch.Branch_id]].Block_bounds)
						poli_branches = append(poli_branches, areas.Bounds)
						bounds_branches := *Polygon_GetArea(poli_branches)
						centroid := Centroid(bounds_branches.Polygon)
						centroidText := CentroidText(centroid)
						bqDataArr[branch_bqs[branch.Branch_id]].Block_bounds = bounds_branches
						bqDataArr[branch_bqs[branch.Branch_id]].Centroid = *centroid
						bqDataArr[branch_bqs[branch.Branch_id]].Centroid_text = centroidText
					}

				}

				poli_areas = append(poli_areas, areas.Bounds)

				var abrv_areas string
				if len(areas.Areas_name) > 1 && len(areas.Areas_name) < 3 {
					abrv_areas = areas.Areas_name[0:2]
				} else if len(areas.Areas_name) < 2 {
					abrv_areas = areas.Areas_name[0:1]
				} else {
					abrv_areas = areas.Areas_name[0:3]
				}

				bqData = BQData{block_id, areas.Client_id, block_id_branch, areas.Areas_name, areas.Bounds, abrv_areas, areas.Centroid, areas.Centroid_text}
				bqDataArr = append(bqDataArr, bqData)

				//tabela separada
				table_id = TableId{block_id, areas.Areas_id, "Areas"}
				table_id_arr = append(table_id_arr, table_id)
				//

				block_id++

			}

		}

		bounds_client := *Polygon_GetArea(poli_areas)
		centroid := Centroid(bounds_client.Polygon)
		centroidText := CentroidText(centroid)

		abrv_client := client.Client_name[0:3]
		bqData = BQData{
			Block_id:      block_id_cli,
			ClientBQ_id:   client.Client_id,
			Block_parent:  0,
			Block_name:    client.Client_name,
			Block_bounds:  bounds_client,
			Block_abrv:    abrv_client,
			Centroid:      *centroid,
			Centroid_text: centroidText}
		bqDataArr = append(bqDataArr, bqData)
		client_bqs[client.Client_id] = len(bqDataArr) - 1

		//tabela separada
		table_id = TableId{block_id_cli, client.Client_id, "Client"}
		table_id_arr = append(table_id_arr, table_id)
		//

	}

	redis := dbRedis.ConnectRedis()
	defer redis.Close()

	fmt.Println("bqdata", bqDataArr)

	// Printing Clients
	fmt.Println("\nClients:")
	for _, i := range client_bqs {
		fmt.Println(bqDataArr[i])
	}

	// Printing Farms
	fmt.Println("\nFarms:")
	for _, i := range farm_bqs {
		fmt.Println(bqDataArr[i])
	}

	// Printing Branches
	fmt.Println("\nBranches:")
	for _, i := range branch_bqs {
		fmt.Println(bqDataArr[i])
	}

	//Tabela separada com o id de cada coisa, o block_id de cada coisa e a identificação
	//fmt.Println("\nTableId: ", table_id_arr)

	for _, v := range bqDataArr {
		json, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err)
		}

		err = redis.Set(strconv.Itoa(v.Block_id)+":"+strconv.Itoa(v.Block_parent), json, 0).Err()
		if err != nil {
			fmt.Println(err)
		}

	}

	return bqDataArr

}
