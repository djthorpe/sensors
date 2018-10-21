# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
    
all: test install

test: test_protocol

install: install_spi install_i2c install_mihome install_ener314

install_i2c:
	$(GOINSTALL) -tags "rpi i2c" ./cmd/bme280/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/bme680/...
	$(GOINSTALL) -tags "rpi i2c" ./cmd/tsl2561/...

install_spi:
	$(GOINSTALL) -tags "rpi spi" ./cmd/bme280/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/bme680/...

install_mihome:
	$(GOINSTALL) -tags "rpi spi" ./cmd/mihome/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/mihome_client/...
	$(GOINSTALL) -tags "rpi spi" ./cmd/mihome_service/...

install_ener314:
	$(GOINSTALL) -tags "rpi" ./cmd/ener314/...

test_protocol: 
	$(GOTEST) -tags rpi ./protocol/...

clean: 
	$(GOCLEAN)
