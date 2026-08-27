package service

import "github.com/Euhcslel/SagaWeb/internal/repository"

// GetDealersList возвращает список дилеров.
func GetDealersList() ([]UserInfo, error) {
	dealers, err := repository.GetDealersList()
	if err != nil {
		return nil, err
	}
	var dealersInfo []UserInfo
	for _, dealer := range dealers {
		address := ""
		if dealer.Address != nil {
			address = *dealer.Address
		}

		dealersInfo = append(dealersInfo, UserInfo{
			ID:       dealer.User.ID,
			Fullname: dealer.Company.Name,
			Phone:    dealer.User.PhoneNumber,
			Email:    dealer.User.Email,
			IsDealer: true,

			Dealer: &DealerInfo{
				CompanyName: dealer.Company.Name,
				Address:     address,
			},
		})

	}
	return dealersInfo, nil
}
