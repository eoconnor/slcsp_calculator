# SLCSP Calculator

This Golang program parses a series of CSV files containing data on health plans, and determines the second-lowest cost silver-level plan for a series of input zip codes. Details can be found here: https://homework.adhoc.team/slcsp.

# Installation and Use

On a machine with Go installed (this program uses version 1.13), cd to the directory where you want to install the program and perform the following steps:

- `git clone https://github.com/eoconnor/slcsp_calculator.git`
- `cd slcsp_calculator`
- `go install`
- `slcsp_calculator`

The program will run and the output will be written to the terminal window. There will be a series of logging statements, followed by the line `*** OUTPUT:`, followed by the CSV-formatted output.

# Implementation Notes

All application logic is contained in the file `slcsp_calculator.go`. The file `slcsp_calculator_test.go` contains some unit tests. The input CSV-formatted data files are located in the `data` directory.

The program uses the Golang CSV parsing library (https://golang.org/pkg/encoding/csv/), which supports the CSV format described in RFC 4180.

# To Do's

This section lists items that could use some additional work.

- Currently, the code is all in a single file. It could be broken up into separate packages that provide a better logical grouping; e.g., primary business logic, utility functions, etc.
- Moar tests! I have a few for the `getRateArea` and `contains` helper functions, but all the significant logic is untested. This would have required mocking the CSV parsing functionality so that we could supply our own test input to that logic.
- In a real-world scenario, data of this type would be loaded into a relational database for persistence and easier lookup. I used lists and maps to mimic that kind of access.
- Better input data validation; e.g., there's no logic that would catch a non-numeric value in the `rate_area` column in `plans.csv`.
- To improve performance, we could get rid of the `getInputZips` function and just iterate over the records from `slcsp.csv`.
- The index values used to grab column values from the various file type records are hard-coded throughout the code. They should be defined as constants in one location.
