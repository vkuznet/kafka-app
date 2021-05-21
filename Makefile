VERSION=`git rev-parse --short HEAD`
flags=-ldflags="-s -w -X web.version=${VERSION}"
odir=`cat ${PKG_CONFIG_PATH}/oci8.pc | grep "libdir=" | sed -e "s,libdir=,,"`

all: build

vet:
	go vet .

build:
	go clean; rm -rf pkg kafka-app*; go build ${flags}

build_all: build build_osx build_linux build_power8 build_arm64

build_osx:
	go clean; rm -rf pkg kafka-app_osx; GOOS=darwin go build ${flags}
	mv kafka-app kafka-app_osx

build_linux:
	go clean; rm -rf pkg kafka-app_linux; GOOS=linux go build ${flags}
	mv kafka-app kafka-app_linux

build_power8:
	go clean; rm -rf pkg kafka-app_power8; GOARCH=ppc64le GOOS=linux go build ${flags}
	mv kafka-app kafka-app_power8

build_arm64:
	go clean; rm -rf pkg kafka-app_arm64; GOARCH=arm64 GOOS=linux go build ${flags}
	mv kafka-app kafka-app_arm64

install:
	go install

clean:
	go clean; rm -rf pkg

test: test-dbs test-sql test-validator test-bulk test-http

test-dbs:
	cd test && rm -f /tmp/dbs-test.db && sqlite3 /tmp/dbs-test.db < ../static/schema/sqlite-schema.sql && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -v -run TestDBS
test-bulk:
	cd test && rm -f /tmp/dbs-test.db && sqlite3 /tmp/dbs-test.db < ../static/schema/sqlite-schema.sql && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -v -run Bulk
test-sql:
	cd test && rm -f /tmp/dbs-test.db && sqlite3 /tmp/dbs-test.db < ../static/schema/sqlite-schema.sql && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -v -run SQL
test-validator:
	cd test && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -v -run Validator
test-http:
	cd test && rm -f /tmp/dbs-test.db && sqlite3 /tmp/dbs-test.db < ../static/schema/sqlite-schema.sql && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -v -run HTTP
bench:
	cd test && rm -f /tmp/dbs-test.db && sqlite3 /tmp/dbs-test.db < ../static/schema/sqlite-schema.sql && LD_LIBRARY_PATH=${odir} DYLD_LIBRARY_PATH=${odir} go test -run Benchmark -bench=.
