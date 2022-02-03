package models

import (
	"etl-blocks2/db"
	"fmt"
	"log"
	"strconv"
)

type PgData struct {
	Client_id      int
	Client_name    string
	Mapping_id     int
	Mapping_name   string
	Area_farm_id   int
	Farm_name      string
	Area_branch_id int
	Branch_name    string
	Areas_id       int
	Areas_name     string
}

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
	Client_id  int
	Areas_name string
}

type Branch struct {
	Client_id    int
	Branch_name  string
	Area_farm_id int
}

type Branch_2 struct {
	Client_id   int
	Branch_name string
}

type Farm struct {
	Client_id int
	Farm_name string
}

func ReadPg() {

	db := db.DbConect()
	defer db.Close()

	queryClient := fmt.Sprintf(
		"select cli.client_id, cli.name as client_name " +
			"from client cli " +
			"where cli.client_id in (163);")

	rows, err := db.Query(queryClient)
	if err != nil {
		log.Println("Error:", err.Error())
	}

	client := Client{}
	var clientPg []Client

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

	//var blocklist []BQData

	for rows.Next() {
		err := rows.Scan(&client.Client_id, &client.Client_name)
		if err != nil {
			log.Println("Error:", err.Error())
		}

		clientPg = append(clientPg, client)

		c := strconv.Itoa(client.Client_id)

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

				areasPg = append(areasPg, areas)

				areas_2 = Areas_2{areas.Client_id, areas.Areas_name}

				areas_2_Pg = append(areas_2_Pg, areas_2)

				queryBranch := fmt.Sprintf(
					"select " + strconv.Itoa(mapping.Client_id) + " as Client_id , " + "b.name as branch_name, b.area_farm_id " +
						"from area_branch b " +
						"where b.area_branch_id = " + strconv.Itoa(areas.Areas_branch_id))

				rows, err := db.Query(queryBranch)
				if err != nil {
					log.Println("Error:", err.Error())
				}
				for rows.Next() {

					err := rows.Scan(&branch.Client_id, &branch.Branch_name, &branch.Area_farm_id)
					if err != nil {
						log.Println("Error:", err.Error())
					}

					branchPg = append(branchPg, branch)

					branch_2 = Branch_2{branch.Client_id, branch.Branch_name}

					branch_2_Pg = append(branch_2_Pg, branch_2)

					queryFarm := fmt.Sprintf(
						"select " + strconv.Itoa(mapping.Client_id) + " as Client_id , " + "f.name as farm_name " +
							"from area_farm f " +
							"where f.area_farm_id = " + strconv.Itoa(branch.Area_farm_id))

					rows, err := db.Query(queryFarm)
					if err != nil {
						log.Println("Error:", err.Error())
					}

					for rows.Next() {

						err := rows.Scan(&farm.Client_id, &farm.Farm_name)
						if err != nil {
							log.Println("Error:", err.Error())
						}

						farmPg = append(farmPg, farm)
					}

				}

			}

		}

	}
	fmt.Println("cli", clientPg)
	fmt.Println("=========")
	fmt.Println("map", mappingPg)
	fmt.Println("=========")
	fmt.Println("areas", areas_2_Pg)
	fmt.Println("=========")
	fmt.Println("branch", branch_2_Pg)
	fmt.Println("=========")
	fmt.Println("farm", farmPg)
}
