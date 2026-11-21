package api

// maxPages caps how many pages allPages will fetch, so a bogus total_pages
// from the API can't make us loop forever (TMDB itself caps at 500).
const maxPages = 500

// allPages calls fetch for page 1, 2, ... until the reported total page count
// (or maxPages) is reached, and concatenates the results. fetch returns the
// page's results and the total number of pages.
func allPages[T any](fetch func(page int) ([]T, int, error)) ([]T, error) {
	var all []T
	for page := 1; page <= maxPages; page++ {
		results, totalPages, err := fetch(page)
		if err != nil {
			return nil, err
		}
		all = append(all, results...)
		if page >= totalPages {
			break
		}
	}
	return all, nil
}
