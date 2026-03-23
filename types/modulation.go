package types

type Modulation int

const (
	ModulationFM Modulation = iota
	ModulationNFM
	ModulationAM
	ModulationWFM
	ModulationLSB
	ModulationUSB
	ModulationCW
	ModulationC4FM
	ModulationDSTAR
	ModulationP25
	ModulationNXDN
	ModulationDMR
	ModulationYSF
	ModulationFUSION
	ModulationPOCSAG
	ModulationDPMR
	ModulationTETRA
)
