package main

const (
	// SatInBTC represents number of satoshis in 1 bitcoin
	SatInBTC = uint64(100000000)

	// WavesNodeURL is an URL for Waves Node
	WavesNodeURL = "https://nodes.wavesnodes.com"

	// WavesMonitorTick interval in seconds
	WavesMonitorTick = 10

	// PricesURL is URL for crypo prices
	PricesURL = "https://min-api.cryptocompare.com/data/price?fsym=WAVES&tsyms=BTC,ETH,HRK,USD,EUR,NGN,JPY"

	// AnoteId is Anote's WAVES asset Id
	AnoteId = "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf"

	// AINTId is AINT's WAVES asset Id
	AINTId = "66DUhUoJaoZcstkKpcoN3FUcqjB6v8VJd5ZQd6RsPxhv"

	// WavesFee represents fee amount in Waves
	WavesFee = 100000

	// TelPollerTimeout is Telegram poller timeout in seconds
	TelPollerTimeout = 30

	// TelAnonOps group for error logging
	TelAnonOps = -1001213539865

	// TelAnonTeam group for team messages
	TelAnonTeam = -1001280228955

	// SigTick represents a ticker in seconds for signal handler
	SigTick = 10
)
