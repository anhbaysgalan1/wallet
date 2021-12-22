package transfer

const (
	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain.
	ErrNonceTooLow = "nonce too low"

	// ErrNonceTooHigh is returned if the nonce of a transaction is higher than the
	// next one expected based on the local chain.
	ErrNonceTooHigh = "nonce too high"

	// ErrGasLimitReached is returned by the gas pool if the amount of gas required
	// by a transaction is higher than what's left in the block.
	ErrGasLimitReached = "gas limit reached"

	// ErrInsufficientFundsForTransfer is returned if the transaction sender doesn't
	// have enough funds for transfer(topmost call only).
	ErrInsufficientFundsForTransfer = "insufficient funds for transfer"

	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = "insufficient funds for gas * price + value"

	// ErrGasUintOverflow is returned when calculating gas usage.
	ErrGasUintOverflow = "gas uint64 overflow"

	// ErrIntrinsicGas is returned if the transaction is specified to use less gas
	// than required to start the invocation.
	ErrIntrinsicGas = "intrinsic gas too low"
)
