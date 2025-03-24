package service

import (
	"go-rest-api/domain"
	"go-rest-api/dto"

	"github.com/google/uuid"
)

type CustomerService struct {
	customerRepo domain.CustomerRepository
}

func NewCustomerService(customerRepo domain.CustomerRepository) domain.CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) GetAllCustomers() ([]dto.CustomerResponse, error) {
	customers, err := s.customerRepo.FindAll()
	if err != nil {
		return nil, err
	}

	customerResponses := make([]dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = dto.CustomerResponse{
			ID:   customer.ID,
			Code: customer.Code,
			Name: customer.Name,
		}
	}

	return customerResponses, nil
}

func (s *CustomerService) GetCustomerByID(id uuid.UUID) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.CustomerResponse{
		ID:   customer.ID,
		Code: customer.Code,
		Name: customer.Name,
	}, nil
}

func (s *CustomerService) GetCustomerByCode(code string) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	return &dto.CustomerResponse{
		ID:   customer.ID,
		Code: customer.Code,
		Name: customer.Name,
	}, nil
}

func (s *CustomerService) CreateCustomer(req dto.CustomerCreateRequest) (*dto.CustomerResponse, error) {
	customer := domain.Customer{
		ID:   uuid.New(),
		Code: req.Code,
		Name: req.Name,
	}

	err := s.customerRepo.Create(&customer)
	if err != nil {
		return nil, err
	}

	return &dto.CustomerResponse{
		ID:   customer.ID,
		Code: customer.Code,
		Name: customer.Name,
	}, nil
}

func (s *CustomerService) UpdateCustomer(id uuid.UUID, req dto.CustomerUpdateRequest) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	customer.Code = req.Code
	customer.Name = req.Name

	err = s.customerRepo.Update(customer)
	if err != nil {
		return nil, err
	}

	return &dto.CustomerResponse{
		ID:   customer.ID,
		Code: customer.Code,
		Name: customer.Name,
	}, nil
}

func (s *CustomerService) DeleteCustomer(id uuid.UUID) error {
	return s.customerRepo.Delete(id)
}
