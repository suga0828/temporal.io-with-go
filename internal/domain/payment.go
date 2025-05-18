package domain

// TransferDetails contains all the information needed for a money transfer
type TransferDetails struct {
	SourceAccount string // Account number to withdraw from
	TargetAccount string // Account number to deposit to
	Amount        int    // Amount to transfer
	ReferenceID   string // Unique reference ID for the transaction
}
