package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	// set up path and handlers
	http.HandleFunc("/echo", handleEcho)
	http.HandleFunc("/invert", handleInvert)
	http.HandleFunc("/flatten", handleFlatten)
	http.HandleFunc("/sum", handleSum)
	http.HandleFunc("/multiply", handleMultiply)

	fmt.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("error starting server: %s\n", err.Error())
	}
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	handleCSVParsing(w, r, func(matrix [][]string) [][]string {
		return matrix
	})
}

func handleInvert(w http.ResponseWriter, r *http.Request) {
	handleCSVParsing(w, r, func(matrix [][]string) [][]string {
		inverted := make([][]string, len(matrix[0]))
		for i := range inverted {
			inverted[i] = make([]string, len(matrix))
			for j := range inverted[i] {
				inverted[i][j] = matrix[j][i]
			}
		}
		return inverted
	})
}

func handleFlatten(w http.ResponseWriter, r *http.Request) {
	handleCSVParsing(w, r, func(matrix [][]string) [][]string {
		flat := [][]string{{}}
		for _, row := range matrix {
			for _, value := range row {
				flat[0] = append(flat[0], value)
			}
		}
		return flat
	})
}

func handleSum(w http.ResponseWriter, r *http.Request) {
	handleCSVParsing(w, r, func(matrix [][]string) [][]string {
		total := 0
		for _, row := range matrix {
			for _, value := range row {
				num, _ := strconv.Atoi(value)
				total += num
			}
		}
		return [][]string{{strconv.Itoa(total)}}
	})
}

func handleMultiply(w http.ResponseWriter, r *http.Request) {
	handleCSVParsing(w, r, func(matrix [][]string) [][]string {
		product := 1
		for _, row := range matrix {
			for _, value := range row {
				num, _ := strconv.Atoi(value)
				product *= num
			}
		}
		return [][]string{{strconv.Itoa(product)}}
	})
}

// generic handler to handle the csv opearations
func handleCSVParsing(w http.ResponseWriter, r *http.Request, operation func(matrix [][]string) [][]string) {
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error %s", err.Error())))
		return
	}

	if len(records) == 0 {
		w.Write([]byte("error: Empty CSV file"))
		return
	}

	numColumns := len(records[0])
	for i, row := range records {
		if len(row) != numColumns {
			w.Write([]byte(fmt.Sprintf("error: row %d has %d columns, expected %d columns", i+1, len(row), numColumns)))
			return
		}
	}

	result := operation(records)

	var response string
	for _, row := range result {
		response = fmt.Sprintf("%s%s\n", response, strings.Join(row, ","))
	}

	fmt.Fprint(w, response)
}
