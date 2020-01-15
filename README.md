# SLCSP Calculator

This program parses a series of CSV files containing data on health plans, and determines the second-lowest cost silver-level plan for a series of input zip codes. Details can be found here: https://homework.adhoc.team/slcsp.

# Installation and Use



# Implementation Notes

This program uses the Golang CSV parsing library (https://golang.org/pkg/encoding/csv/), which supports the CSV format described in RFC 4180.

# To Do's

This section lists items that could use some additional work.

- Currently, the code is all in a single file. It could be broken up into separate packages that provide a better logical grouping; e.g., primary business logic, utility functions, etc.
- Moar tests! I have a few for the `getRateArea` and `contains` helper functions, but all the significant logic is untested. This would have required mocking the CSV parsing functionality so that we could supply our own test input to that logic.
- In a real-world scenario, data of this type would be loaded into a relational database for persistence and easier lookup. I used lists and maps to mimic that kind of access.
- Better input data validation; e.g., there's no logic that would catch a non-numeric value in the `rate_area` column in `plans.csv`.
- To improve performance, we could get rid of the `getInputZips` function and just iterate over the records from `slcsp.csv`.
- The index values used to grab column values from the various file type records are hard-coded throughout the code. They should be defined as constants in one location.
