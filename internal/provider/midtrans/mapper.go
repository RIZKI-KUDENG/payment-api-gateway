package midtrans

func NormalizeStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus{
		case "capture":
			if fraudStatus == "challenge" {
				return "pending"
			}
			if fraudStatus == "accept" {
				return "success"
			}
			return "pending"
		case "settlement":
			return "success"
		case "pending":
			return "pending"
		case "deny","cancel":
			return "pending"
		case "expired":
			return "expired"
		case "refund", "partial_refund":
			return "refunded"
		default:
			return "unknown"
	}
}
