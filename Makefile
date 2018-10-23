# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGENERATE=$(GOCMD) generate
    
all: test install

test: test_protocol

install: generate install_spi install_i2c install_mihome install_ener314

generate:
	$(GOGENERATE)  ./rpc/protobuf/...

install_i2c:
	$(GOINSTALL) -tags "rpi i2c" ./cmd/bme280/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/bme680/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/tsl2561/...

install_spi:
	$(GOINSTALL) -tags "rpi spi" ./cmd/bme280/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/bme680/...

install_mihome:
	$(GOINSTALL) -tags "rpi spi" ./cmd/mihome/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/mihome-client/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/mihome-service/...

install_ener314:
	$(GOINSTALL) -tags "rpi" ./cmd/ener314/...

test_protocol: 
	$(GOTEST) -tags rpi ./protocol/...

clean: 
	$(GOCLEAN)
