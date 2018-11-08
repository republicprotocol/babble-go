# Simulate the environment used by Travis CI so that we can run local tests to
# find and resolve issues that are consistent with the Travis CI environment.
# This is helpful because Travis CI often finds issues that our own local tests
# do not.

# go vet ./...
# golint -set_exit_status `go list ./... | grep -Ev "(vendor)"`

go build ./...

# Test and generate cover profiles
GOMAXPROCS=1 CI=true ginkgo --cover adapter/db  \
                                    adapter/rpc \
                                    core/addr   \
                                    core/gossip

# Merge cover profiles into one root cover profile
covermerge adapter/db/db.coverprofile      \
           adapter/rpc/rpc.coverprofile    \
           core/addr/addr.coverprofile     \
           core/gossip/gossip.coverprofile \
           > babble.coverprofile

# Remove auto-generated protobuf files
sed -i '/.pb.go/d' babble.coverprofile

# Remove marshaling files
sed -i '/marshal.go/d' babble.coverprofile