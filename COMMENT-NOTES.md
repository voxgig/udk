# Implementation rationale

The Go representation cannot distinguish a missing value from JSON null wherever both become nil. Test adapters must make that representational boundary explicit rather than infer an undefined value that the host representation cannot hold.

List property updates must avoid sparse holes: out-of-range positions append or prepend, and removal shifts remaining elements. These operations preserve a dense list representation across ports.

Sources: [Go test adapter](go/test/utility/struct_runner_test.go), [Go structural utilities](go/utility/struct/voxgigstruct.go).
