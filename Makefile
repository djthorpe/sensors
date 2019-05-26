# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
PKG_CONFIG_PATH="/opt/vc/lib/pkgconfig"

# App parameters
GOPI=github.com/djthorpe/gopi
GOLDFLAGS += -X $(GOPI).GitTag=$(shell git describe --tags)
GOLDFLAGS += -X $(GOPI).GitBranch=$(shell git name-rev HEAD --name-only --always)
GOLDFLAGS += -X $(GOPI).GitHash=$(shell git rev-parse HEAD)
GOLDFLAGS += -X $(GOPI).GoBuildTime=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOFLAGS = -ldflags "-s -w $(GOLDFLAGS)" 

all: test install

test: test_protocol

install: mihome mihome-rpi spi i2c 

i2c:
	$(GOINSTALL) -tags "rpi i2c" ./cmd/bme280/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/bme680/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/tsl2561/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/ads1015/...

spi:
	$(GOINSTALL) -tags "rpi spi" ./cmd/bme280/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/bme680/...

protobuf:
	$(GOGENERATE)  ./rpc/protobuf/...

mihome-rpi: mihome
	PKG_CONFIG_PATH=$(PKG_CONFIG_PATH) $(GOINSTALL) -tags "rpi spi" $(GOFLAGS) ./cmd/mihome-service/...

mihome: protobuf
	$(GOINSTALL) -tags "rpi spi" $(GOFLAGS) ./cmd/mihome-client/...
	$(GOINSTALL) -tags "rpi spi" $(GOFLAGS) ./cmd/mihome-reset/...

test_protocol: 
	$(GOTEST) -tags rpi ./protocol/...

clean: 
	$(GOCLEAN)
