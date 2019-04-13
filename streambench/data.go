package streambench

// Data is the stucture we'll be shoveling data into and out of using a
// marshaller or encoder.
type Data struct {
	ID       string
	Name     string
	Email    string
	Whatever int
}
