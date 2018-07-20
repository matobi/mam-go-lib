LDFLAGS = -extldflags -static -s -w

#build:
#	env CGO_ENABLED=0 vgo build -ldflags "${LDFLAGS}" -o ./buildtarget/mam-golib-example ./cmd/mam-golib-example

#clean:
#	rm -rf ./buildtarget/*

test:
	vgo test github.com/matobi/mam-go-lib/pkg/conf/test
