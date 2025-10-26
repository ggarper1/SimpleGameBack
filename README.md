# Simple Game Go Backend 🎮

This is a backend for a [simple iPhone game]([https://pages.github.com/](https://github.com/ggarper1/SimpleGameFront) I've created.

## Directory Structure 🗃️
```
├── README.md
├── go.mod
├── src: directory containing backend code.
│   │
│   ├── routes: directory containing files with endpoints.
│   │
│   ├── services: directory containing game logic and other type of logic.
│   │   └── map_generator.go: code for generating game map.
│   │
│   └── storage: directory containing data structures.
│       ├── models: files containing code defining models stored in database.
│       │
│       └── objects: files containing data structures to use in logic and endpoints, that is, data structures used during endpoint handling that are never stored.
│           ├── line.go: code for struct representing a line.
│           ├── piece.go: code for struct representing a piece.
│           ├── point.go: code for struct representing a point.
│           └── segment.go:  code for struct representing a segment.
│
└── tests: directory containing tests for various parts of the backend.
    └── objects_tests: tests for structs in objects directory.
        ├── dataset_generator.py: python script for generating test dataset for objects tests.
        ├── line_test.go
        ├── point_test.go
        ├── segment_test.go
        └── utils.go: test utils.
```
