package sensor

// Interface for the Bosch BME280 chip
type BME280 interface {

	// Get Version
	ChipIDVersion() (uint8, uint8)

	// Get Mode
	Mode() BME280Mode

	// Return IIR filter co-officient
	Filter() BME280Filter

	// Return standby time
	Standby() BME280Standby

	// Return oversampling values osrs_t, osrs_p, osrs_h
	Oversample() (BME280Oversample, BME280Oversample, BME280Oversample)

	// Return current measuring and updating value
	Status() (bool, bool, error)

	// Reset
	SoftReset() error

	// Set BME280 mode
	SetMode(mode BME280Mode) error

	// Set Oversampling
	SetOversample(osrs_t, osrs_p, osrs_h BME280Oversample) error

	// Set Filter
	SetFilter(filter BME280Filter) error

	// Set Standby mode
	SetStandby(t_sb BME280Standby) error

	// Return raw sample data for temperature, pressure and humidity
	ReadSample() (float64, float64, float64, error)

	// Return altitude in meters for given pressure
	AltitudeForPressure(atmospheric, sealevel float64) float64
}
