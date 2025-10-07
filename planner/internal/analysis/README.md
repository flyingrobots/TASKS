# Package `analysis`

The `analysis` package provides functionalities for performing a census and analysis of a codebase.

## `RunCensus` Function

The `RunCensus` function recursively walks a specified directory path, identifies all files, and categorizes them based on their properties (e.g., file extension). It returns a `CodebaseAnalysis` struct containing the discovered information.

## `CodebaseAnalysis` Struct

The `CodebaseAnalysis` struct holds the results of a codebase scan, including:
- `Files`: A list of all file paths found.
- `GoFiles`: A list of file paths specifically for Go source files.

This package is the first step in the T.A.S.K.S. planner's process, providing an initial understanding of the codebase's structure.
