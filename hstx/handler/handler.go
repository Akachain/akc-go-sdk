package handler

// Handler ...
type Handler struct {
	SuperAdminHanler *SuperAdminHanler
	AdminHanler      *AdminHanler
	ProposalHanler   *ProposalHanler
	ApprovalHanler   *ApprovalHanler
}

// InitHandler ...
func (h *Handler) InitHandler() {
	h.SuperAdminHanler = new(SuperAdminHanler)
	h.AdminHanler = new(AdminHanler)
	h.ProposalHanler = new(ProposalHanler)
	h.ApprovalHanler = new(ApprovalHanler)
}
