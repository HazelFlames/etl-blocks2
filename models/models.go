package models

import (
	"etl-blocks2/db"
	"fmt"
	"log"
	"strconv"
)

type BQData struct {
	ClientBQ_id  int
	Block_id     int
	Block_parent int
	Block_name   string
}

type Client struct {
	Client_id   int
	Client_name string
}

type Client_2 struct {
	Client_id    int
	Block_id     int
	Block_parent int
	Client_name  string
}

type Mapping struct {
	Client_id  int
	Mapping_id int
}

type Areas struct {
	Client_id       int
	Areas_name      string
	Areas_branch_id int
}

type Areas_2 struct {
	Client_id    int
	Block_id     int
	Block_parent int
	Areas_name   string
}

type Branch struct {
	Client_id    int
	Branch_id    int
	Branch_name  string
	Area_farm_id int
}

type Branch_2 struct {
	Client_id    int
	Block_id     int
	Block_parent int
	Branch_name  string
}

type Farm struct {
	Client_id int
	Farm_id   int
	Farm_name string
}

type Farm_2 struct {
	Client_id    int
	Block_id     int
	Block_parent int
	Farm_name    string
}

var FarmList []int

func VerifyFarmID(farmID int) bool {
	var status bool = false
	for _, x := range FarmList {
		if x == farmID {
			status = true
			break
		}
	}
	return status
}

var BranchList []int

func VerifyBranchID(branchID int) bool {
	var status bool = false
	for _, x := range BranchList {
		if x == branchID {
			status = true
			break
		}
	}
	return status
}

func ReadPg() {

	db := db.DbConect()
	defer db.Close()

	block_id := 1

	client := Client{}
	var clientPg []Client

	client_2 := Client_2{}
	var clientPg_2 []Client_2

	mapping := Mapping{}
	var mappingPg []Mapping

	areas := Areas{}
	var areasPg []Areas

	areas_2 := Areas_2{}
	var areas_2_Pg []Areas_2

	branch := Branch{}
	var branchPg []Branch

	branch_2 := Branch_2{}
	var branch_2_Pg []Branch_2

	farm := Farm{}
	var farmPg []Farm

	farm_2 := Farm_2{}
	var farmPg_2 []Farm_2

	bqData := BQData{}
	var bqDataArr []BQData

	var block_id_cli int
	var block_id_farm int
	var block_id_branch int

	//client
	queryClient := fmt.Sprintf(
		"select cli.client_id, cli.name as client_name " +
			"from client cli " +
			"where cli.client_id in (163, 545);")

	rows, err := db.Query(queryClient)
	if err != nil {
		log.Println("Error:", err.Error())
	}

	for rows.Next() {
		err := rows.Scan(&client.Client_id, &client.Client_name)
		if err != nil {
			log.Println("Error:", err.Error())
		}

		clientPg = append(clientPg, client)

		client_2 = Client_2{client.Client_id, block_id, 0, client.Client_name}
		clientPg_2 = append(clientPg_2, client_2)

		bqData = BQData{client.Client_id, block_id, 0, client.Client_name}
		bqDataArr = append(bqDataArr, bqData)

		block_id_cli = block_id

		block_id++

		c := strconv.Itoa(client.Client_id)

		//mapping
		queryMapping := fmt.Sprintf(
			"select ma.client_id, ma.id " +
				"from mapping_areas ma " +
				"where ma.client_id = " + c)

		rows, err := db.Query(queryMapping)
		if err != nil {
			log.Println("Error:", err.Error())
		}

		for rows.Next() {
			err := rows.Scan(&mapping.Client_id, &mapping.Mapping_id)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			mappingPg = append(mappingPg, mapping)

			//areas
			queryAreas := fmt.Sprintf(
				"select ar.name as areas_name,  " + strconv.Itoa(mapping.Client_id) + "as Client_id," + "ar.area_branch_id " +
					"from areas ar " +
					"where mapping_area_id = " + strconv.Itoa(mapping.Mapping_id))

			rows, err := db.Query(queryAreas)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			for rows.Next() {
				err := rows.Scan(&areas.Areas_name, &areas.Client_id, &areas.Areas_branch_id)
				if err != nil {
					log.Println("Error:", err.Error())
				}

				//branch
				queryBranch := fmt.Sprintf(
					"select " + strconv.Itoa(mapping.Client_id) + " as Client_id , " + "b.area_branch_id as branch_id, b.name as branch_name, b.area_farm_id " +
						"from area_branch b " +
						"where b.area_branch_id = " + strconv.Itoa(areas.Areas_branch_id))

				rows, err := db.Query(queryBranch)
				if err != nil {
					log.Println("Error:", err.Error())
				}
				for rows.Next() {

					err := rows.Scan(&branch.Client_id, &branch.Branch_id, &branch.Branch_name, &branch.Area_farm_id)
					if err != nil {
						log.Println("Error:", err.Error())
					}

					//farm
					queryFarm := fmt.Sprintf(
						"select distinct " + strconv.Itoa(mapping.Client_id) + " as Client_id , f.area_farm_id, " + "f.name as farm_name " +
							"from area_farm f " +
							"where f.area_farm_id = " + strconv.Itoa(branch.Area_farm_id))

					rows, err := db.Query(queryFarm)
					if err != nil {
						log.Println("Error:", err.Error())
					}

					for rows.Next() {

						err := rows.Scan(&farm.Client_id, &farm.Farm_id, &farm.Farm_name)
						if err != nil {
							log.Println("Error:", err.Error())
						}

						farmPg = append(farmPg, farm)

						farm_2 = Farm_2{farm.Client_id, block_id, block_id_cli, farm.Farm_name}

						bqData = BQData{farm.Client_id, block_id, block_id_cli, farm.Farm_name}

						//if para validar se farm.Farm_id ja existe no array
						if !VerifyFarmID(farm.Farm_id) { //Element is not present in the slice
							FarmList = append(FarmList, farm.Farm_id)
							farmPg_2 = append(farmPg_2, farm_2)
							bqDataArr = append(bqDataArr, bqData)
							block_id_farm = block_id

						}

						// bqData = BQData{farm.Client_id, block_id, block_id_cli, farm.Farm_name}

						// bqDataArr = append(bqDataArr, bqData)

						block_id++

					}
					branchPg = append(branchPg, branch)

					branch_2 = Branch_2{branch.Client_id, block_id, block_id_farm, branch.Branch_name}

					bqData = BQData{branch.Client_id, block_id, block_id_farm, branch.Branch_name}

					if !VerifyBranchID(branch.Branch_id) { //Element is not present in the slice
						BranchList = append(BranchList, branch.Branch_id)
						branch_2_Pg = append(branch_2_Pg, branch_2)
						bqDataArr = append(bqDataArr, bqData)
						block_id_branch = block_id
					}

					// bqData = BQData{branch.Client_id, block_id, branch.Branch_name}
					// bqDataArr = append(bqDataArr, bqData)

					block_id++

				}
				areasPg = append(areasPg, areas)

				areas_2 = Areas_2{areas.Client_id, block_id, block_id_branch, areas.Areas_name}

				areas_2_Pg = append(areas_2_Pg, areas_2)

				bqData = BQData{areas.Client_id, block_id, block_id_branch, areas.Areas_name}
				bqDataArr = append(bqDataArr, bqData)

			}
			block_id++

		}

	}

	fmt.Println("cli", clientPg_2)
	fmt.Println("=========")
	fmt.Println("map", mappingPg)
	fmt.Println("=========")
	fmt.Println("areas", areas_2_Pg)
	fmt.Println("=========")
	fmt.Println("branch", branch_2_Pg)
	fmt.Println("=========")
	fmt.Println("farm", farmPg_2)
	fmt.Println("=========")
	fmt.Println("bqdata", bqDataArr)

}
