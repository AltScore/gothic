package model

type AddMovieRequest struct {
	Title string
	Year  int
}

type AddMovieResponse struct {
	ID    int
	Title string
	Year  int
}
