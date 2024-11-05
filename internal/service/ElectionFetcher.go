package service

type parser interface {
	Parse() (int, int, error)
}

type ElectionFetcher struct {
	parser                 parser
	prevFirstCandidateRes  int
	prevSecondCandidateRes int
}

func NewElectionFetcher(parser parser) *ElectionFetcher {
	return &ElectionFetcher{parser: parser}
}

func (ef *ElectionFetcher) Fetch() (int, int, bool, error) {
	curFirstCandidateRes, curSecondCandidateRes, err := ef.parser.Parse()
	if err != nil {
		return 0, 0, false, err
	}
	if curFirstCandidateRes == ef.prevFirstCandidateRes && curSecondCandidateRes == ef.prevSecondCandidateRes {
		return curFirstCandidateRes, curSecondCandidateRes, false, nil
	}

	ef.prevFirstCandidateRes = curFirstCandidateRes
	ef.prevSecondCandidateRes = curSecondCandidateRes

	return curFirstCandidateRes, curSecondCandidateRes, true, nil
}
