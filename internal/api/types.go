package api

// Search results

type MovieResult struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Overview    string  `json:"overview"`
	VoteAverage float64 `json:"vote_average"`
	Rating      float64 `json:"rating,omitempty"` // user's rating (from rated lists)
}

type TVResult struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	FirstAirDate string  `json:"first_air_date"`
	Overview     string  `json:"overview"`
	VoteAverage  float64 `json:"vote_average"`
	Rating       float64 `json:"rating,omitempty"`
}

type PersonResult struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	KnownForDepartment string `json:"known_for_department"`
}

type SearchMoviesResponse struct {
	Page         int           `json:"page"`
	Results      []MovieResult `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

type SearchTVResponse struct {
	Page         int        `json:"page"`
	Results      []TVResult `json:"results"`
	TotalPages   int        `json:"total_pages"`
	TotalResults int        `json:"total_results"`
}

type SearchPersonResponse struct {
	Page         int            `json:"page"`
	Results      []PersonResult `json:"results"`
	TotalPages   int            `json:"total_pages"`
	TotalResults int            `json:"total_results"`
}

// Credits / Filmography

type CastCredit struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`        // movie
	Name        string `json:"name"`         // tv
	MediaType   string `json:"media_type"`
	ReleaseDate string `json:"release_date"` // movie
	FirstAirDate string `json:"first_air_date"` // tv
	Character   string `json:"character"`
}

type CombinedCreditsResponse struct {
	Cast []CastCredit `json:"cast"`
}

// TV details

type TVSeason struct {
	ID           int    `json:"id"`
	SeasonNumber int    `json:"season_number"`
	Name         string `json:"name"`
	EpisodeCount int    `json:"episode_count"`
	AirDate      string `json:"air_date"`
}

type TVDetails struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	FirstAirDate string     `json:"first_air_date"`
	Seasons      []TVSeason `json:"seasons"`
	Overview     string     `json:"overview"`
}

type TVEpisode struct {
	ID            int     `json:"id"`
	EpisodeNumber int     `json:"episode_number"`
	SeasonNumber  int     `json:"season_number"`
	Name          string  `json:"name"`
	AirDate       string  `json:"air_date"`
	Overview      string  `json:"overview"`
	VoteAverage   float64 `json:"vote_average"`
}

type SeasonDetails struct {
	ID           int         `json:"id"`
	SeasonNumber int         `json:"season_number"`
	Name         string      `json:"name"`
	Episodes     []TVEpisode `json:"episodes"`
}

// Auth

type RequestTokenResponse struct {
	Success      bool   `json:"success"`
	RequestToken string `json:"request_token"`
}

type SessionResponse struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id"`
}

type AccountResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// API status response

type StatusResponse struct {
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}
