package loan

type documentService struct{}

func NewDocumentService() DocumentService {
    return &documentService{}
}

func (s *documentService) StoreEvidence(evidence *Evidence) error {
    return nil
}

func (s *documentService) GenerateInvoice(payment *PaymentPeriod) error {
    return nil
}

func (s *documentService) GenerateStatement(loanID string) error {
    return nil
} 