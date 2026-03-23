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

var modulationNames = map[Modulation]string{
	ModulationFM:     "FM",
	ModulationNFM:    "NFM",
	ModulationAM:     "AM",
	ModulationWFM:    "WFM",
	ModulationLSB:    "LSB",
	ModulationUSB:    "USB",
	ModulationCW:     "CW",
	ModulationC4FM:   "C4FM",
	ModulationDSTAR:  "D-STAR",
	ModulationP25:    "P25",
	ModulationNXDN:   "NXDN",
	ModulationDMR:    "DMR",
	ModulationYSF:    "YSF",
	ModulationFUSION: "FUSION",
	ModulationPOCSAG: "POCSAG",
	ModulationDPMR:   "DPMR",
	ModulationTETRA:  "TETRA",
}

func (m Modulation) String() string {
	if s, ok := modulationNames[m]; ok {
		return s
	}
	return "unknown"
}
