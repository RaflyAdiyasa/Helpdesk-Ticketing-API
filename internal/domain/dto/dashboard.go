package dto

type DashboardSummary struct {
	TotalTicket      int64
	OpenTicket       int64
	InProgressTicket int64
	ClosedTicket     int64
	TotalUser        int64
}
