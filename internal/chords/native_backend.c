/*
 * Builds the canonical C chord sources as part of the Go package so ordinary
 * `go build` and `go test` commands do not depend on a prebuilt native library.
 */
#include "../../src/theory/chord_library.c"
#include "../../src/backend/chord_api.c"
