package months

type Month int

const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// Mapear los valores del enum a sus nombres en inglés
func (m Month) String() string {
	return [...]string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}[m-1]
}

var OrderMonths = []Month{
	January, February, March, April, May, June,
	July, August, September, October, November, December,
}
