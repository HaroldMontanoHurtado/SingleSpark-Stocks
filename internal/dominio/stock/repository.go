package stock

type Repository interface {
    Save(s *Stock) error
    SaveBatch(stocks []*Stock) error
    List(limit, offset int) ([]*Stock, error)
    FindByTicker(ticker string) ([]*Stock, error)
}
