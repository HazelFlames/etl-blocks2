package models

import (
	"etl-blocks2/db"
	"fmt"
	"log"
	"math"
	"strconv"

	geojson "github.com/paulmach/go.geojson"
)

type BQData struct {
	Block_id     int
	ClientBQ_id  int
	Block_parent int
	Block_name   string
	Block_bounds geojson.Geometry
	Block_abrv   string
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
					min_lat = coord[0]
					max_lat = coord[0]
					min_long = coord[1]
					max_long = coord[1]
				} else {
					min_lat = math.Min(min_lat, coord[0])
					max_lat = math.Max(max_lat, coord[0])
					min_long = math.Min(min_long, coord[1])
					max_long = math.Max(max_long, coord[1])
				}
				n++
			}
		}
	}
	a_point := fmt.Sprintf("[%f,%f]", min_lat, min_long)
	b_point := fmt.Sprintf("[%f,%f]", min_lat, max_long)
	c_point := fmt.Sprintf("[%f,%f]", max_lat, min_long)
	d_point := fmt.Sprintf("[%f,%f]", max_lat, max_long)
	pol := fmt.Sprintf("[[ %s, %s, %s, %s ]]", a_point, b_point, c_point, d_point)

	rawAreaJSON := []byte(fmt.Sprintf(`{ "type": "Polygon", "coordinates": %s}`, pol))

	polygon_area, err := geojson.UnmarshalGeometry(rawAreaJSON)
	if err != nil {
		fmt.Println("Geo error: ", err)
		return nil
	}
	return polygon_area
}

func ReadPg() {

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
		"select cli.client_id, cli.name as client_name " +
			"from client cli " +
			"where cli.client_id in (163,545);")

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
				"select ar.name as areas_name,  " + strconv.Itoa(mapping.Client_id) + "as Client_id, ar.id, ar.area_branch_id, ST_AsGeoJSON(ar.bounds) as bounds " +
					"from areas ar " +
					"where mapping_area_id = " + strconv.Itoa(mapping.Mapping_id))

			areas_rows, err := db.Query(queryAreas)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			for areas_rows.Next() {

				err := areas_rows.Scan(&areas.Areas_name, &areas.Client_id, &areas.Areas_id, &areas.Areas_branch_id, &areas.Bounds)
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
							bqData = BQData{block_id, farm.Client_id, block_id_cli, farm.Farm_name,
								*Polygon_GetArea(poli_farms), abrv_farm}
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
							bqDataArr[farm_bqs[farm.Farm_id]].Block_bounds = *Polygon_GetArea(poli_farms)
						}

					}

					abrv_branch := branch.Branch_name[0:3]

					if branch_blks[branch.Branch_id] == 0 {
						poli_branches = nil
						poli_branches = append(poli_branches, areas.Bounds)
						bqData = BQData{block_id, branch.Client_id, block_id_farm, branch.Branch_name,
							*Polygon_GetArea(poli_branches), abrv_branch}
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
						bqDataArr[branch_bqs[branch.Branch_id]].Block_bounds = *Polygon_GetArea(poli_branches)
					}

				}

				poli_areas = append(poli_areas, areas.Bounds)

				abrv_areas := areas.Areas_name[0:3]

				bqData = BQData{block_id, areas.Client_id, block_id_branch, areas.Areas_name, areas.Bounds, abrv_areas}
				bqDataArr = append(bqDataArr, bqData)

				//tabela separada
				table_id = TableId{block_id, areas.Areas_id, "Areas"}
				table_id_arr = append(table_id_arr, table_id)
				//

				block_id++

			}

		}

		abrv_client := client.Client_name[0:3]
		bqData = BQData{block_id_cli, client.Client_id, 0, client.Client_name, *Polygon_GetArea(poli_areas), abrv_client}
		bqDataArr = append(bqDataArr, bqData)
		client_bqs[client.Client_id] = len(bqDataArr) - 1

		//tabela separada
		table_id = TableId{block_id_cli, client.Client_id, "Client"}
		table_id_arr = append(table_id_arr, table_id)
		//

	}

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

}
