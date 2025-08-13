### Meteora ☄️ DbcGOSDK

The work done so far aims to maintain correctness & closeness w/ the JS/TS SDK. In that fact, directories, files, and variable names shares closeness with the JS/TS SDK... making comparison very easy. Task forwards from here is to introduce the project to the Meteora dev team, maintain & optimises where possible.

## Running test...

Tests are located in ./helpers/ and ./maths

- to run all test in a directory:

 > go test ./helpers/

- to run a single test file:

> go test ./helpers/curve_test.go

- to run a specific test func:

> go test ./helpers/ -run "TestBuildCurve"

- run a specific sub-test:

> go test ./helpers/ -run "TestBuildCurve/xxxxx"


The idl was generated w/ [solana-anchor-go](https://github.com/fragmetric-labs/solana-anchor-go) from the guys are Fragmetric. The dependency is also inlcuded in the go.mod file w/ [`go tool`](https://www.bytesizego.com/blog/go-124-tool-directive).
