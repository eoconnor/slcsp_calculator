package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	FILE_INPUT_ZIP_CODES = "data/slcsp.csv"
	FILE_ALL_ZIP_CODES   = "data/zips.csv"
	FILE_PLANS           = "data/plans.csv"

	PLAN_LEVEL_BRONZE       = "bronze"
	PLAN_LEVEL_SILVER       = "silver"
	PLAN_LEVEL_GOLD         = "gold"
	PLAN_LEVEL_PLATINUM     = "platinum"
	PLAN_LEVEL_CATASTROPHIC = "catastrophic"
)

type Plan struct {
	ID       string
	Level    string
	Rate     float64
	RateArea string
}

// ByRate implements sort.Interface for []Plan based on the Rate field.
type ByRate []Plan

func (a ByRate) Len() int           { return len(a) }
func (a ByRate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRate) Less(i, j int) bool { return a[i].Rate < a[j].Rate }

func main() {
	rateAreasByZip, err := getRateAreasByZip()
	if err != nil {
		log.Fatal(err)
	}

	plansByRateArea, err := getPlansByRateArea()
	if err != nil {
		log.Fatal(err)
	}

	inputZips, err := getInputZips()
	if err != nil {
		log.Fatal(err)
	}

	// This is our output
	zipsAndRates := make([]string, len(inputZips))

	for zipIndex, inputZip := range inputZips {
		// How many rate areas are in this zip code?
		rateAreas := rateAreasByZip[inputZip]
		if len(rateAreas) != 1 {
			// Either zero or > 1; no return value for this zip
			log.Printf("Zip code %s has %d rate areas; no result\n", inputZip, len(rateAreas))
			zipsAndRates[zipIndex] = fmt.Sprintf("%s,", inputZip)
			continue
		}

		// Find all of the silver plans for this zip's rate area
		plans := plansByRateArea[rateAreas[0]]
		silverPlans := make([]Plan, 0, len(plans))
		for _, plan := range plans {
			if plan.Level == PLAN_LEVEL_SILVER {
				silverPlans = append(silverPlans, plan)
			}
		}

		// Need at least two to have a second lowest value
		if len(silverPlans) < 2 {
			log.Printf("Zip code %s has %d silver plan(s); no result\n", inputZip, len(silverPlans))
			zipsAndRates[zipIndex] = fmt.Sprintf("%s,", inputZip)
			continue
		}

		// Sort the plans by rate and get the lowest rate
		sort.Sort(ByRate(silverPlans))
		lowestRate := silverPlans[0].Rate

		// Multiple plans might have the lowest rate, so don't assume that the second rate is the second lowest
		var secondLowestRate float64
		for i := 1; i < len(silverPlans); i++ {
			if silverPlans[i].Rate > lowestRate {
				secondLowestRate = silverPlans[i].Rate
				break
			}
		}

		if secondLowestRate == 0 {
			log.Printf("Zip code %s has no second lowest rate\n", inputZip)
			zipsAndRates[zipIndex] = fmt.Sprintf("%s,", inputZip)
		} else {
			log.Printf("Zip code %s has second lowest silver plan rate of %.2f\n", inputZip, secondLowestRate)
			zipsAndRates[zipIndex] = fmt.Sprintf("%s,%.2f", inputZip, secondLowestRate)
		}
	}

	fmt.Printf("\n*** OUTPUT:\n\nzipcode,rate\n%s\n", strings.Join(zipsAndRates, "\n"))
}

// getInputZips reads the file of input zip codes and returns them in a list
func getInputZips() ([]string, error) {
	var inputZips []string

	records, err := getFileRecords(FILE_INPUT_ZIP_CODES)
	if err != nil {
		return inputZips, err
	}

	inputZips = make([]string, len(records))
	for i, record := range records {
		inputZips[i] = record[0]
	}

	return inputZips, nil
}

// getRateAreasByZip reads the file containing all zip codes and their associated counties & rate areas
// and returns them as a map where they key is a zip code and the value is a list of all the rate areas
// for that zip code. This allows us to easily look up all rate areas for each of the input zip codes.
func getRateAreasByZip() (map[string][]string, error) {
	rateAreasByZip := make(map[string][]string)

	records, err := getFileRecords(FILE_ALL_ZIP_CODES)
	if err != nil {
		return rateAreasByZip, err
	}

	for _, record := range records {
		zip := record[0]
		rateArea := getRateArea(record[1], record[4])

		// Have we seen this zip before?
		if rateAreas, ok := rateAreasByZip[zip]; ok {
			// Yes; have we seen this rate area?
			if !contains(rateAreas, rateArea) {
				// No, add it to the list
				rateAreasByZip[zip] = append(rateAreas, rateArea)
			}
		} else {
			// No, add it to the map with this rate area
			rateAreasByZip[zip] = []string{rateArea}
		}
	}

	return rateAreasByZip, nil
}

// getPlansByRateArea reads the file containing all the plans and returns them as a map where the
// key is the rate area string (e.g., "al 1") and the value is a list of plans for that rate area.
func getPlansByRateArea() (map[string][]Plan, error) {
	plansByRateArea := make(map[string][]Plan)

	records, err := getFileRecords(FILE_PLANS)
	if err != nil {
		return plansByRateArea, err
	}

	for _, record := range records {
		rateArea := getRateArea(record[1], record[4])
		rate, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return plansByRateArea, err
		}
		plan := Plan{
			ID:       record[0],
			Level:    strings.ToLower(record[2]),
			Rate:     rate,
			RateArea: rateArea,
		}

		// Have we seen this rate area before?
		if plans, ok := plansByRateArea[rateArea]; ok {
			// Yes, add this plan to the list
			plansByRateArea[rateArea] = append(plans, plan)
		} else {
			// No, add it to the map with this plan
			plansByRateArea[rateArea] = []Plan{plan}
		}
	}

	return plansByRateArea, nil
}

// getFileRecords is a helper function for parsing a CSV file and returning the records as a list.
// It omits the header row from the results.
func getFileRecords(filePathname string) ([][]string, error) {
	var records [][]string

	csvFile, err := os.Open(filePathname)
	if err != nil {
		return records, err
	}

	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	records, err = csvReader.ReadAll()
	if err != nil {
		return records, err
	}

	// Leave off the header row
	return records[1:], nil
}

// getRateArea is a helper function for taking a state and rate area number and returning a standard
// format string (e.g., "al 1"). Useful for ensuring that we're comparing rate areas without having
// to worry about things like case.
func getRateArea(state, rateAreaNumber string) string {
	return fmt.Sprintf("%s %s", strings.ToLower(state), rateAreaNumber)
}

// contains is a helpful function for determining if the input string slice the input string.
// Comparison is case-sensitive.
func contains(strSlice []string, strToFind string) bool {
	for _, aString := range strSlice {
		if aString == strToFind {
			return true
		}
	}

	return false
}
