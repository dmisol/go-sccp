package sccp

// SCCP release causes, ITU-T Q.713.
const (
	EndUserOriginated            byte = 0x00 // End user originated
	EndUserCongestion            byte = 0x01 // End user congestion
	EndUserFailure               byte = 0x02 // End user failure
	SCCPUserOriginated           byte = 0x03 // SCCP user originated
	RemoteProcedureError         byte = 0x04 // Remote procedure error
	InconsistentConnectionData   byte = 0x05 // Inconsistent connection data
	AccessFailure                byte = 0x06 // Access failure
	AccessCongestion             byte = 0x07 // Access congestion
	SubsystemFailure             byte = 0x08 // Subsystem failure
	SubsystemCongestion          byte = 0x09 // Subsystem congestion
	MTPFailure                   byte = 0x0A // MTP failure
	NetworkCongestion            byte = 0x0B // Network congestion
	ResetTimerExpiration         byte = 0x0C // Expiration of reset timer
	ReceiveInactivityTimerExpiry byte = 0x0D // Expiration of receive inactivity timer
	Reserved                     byte = 0x0E // Reserved
	Unqualified                  byte = 0x0F // Unqualified
	SCCPFailure                  byte = 0x10 // SCCP failure
)
