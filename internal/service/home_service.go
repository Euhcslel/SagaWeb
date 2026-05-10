package service

import "github.com/Euhcslel/SagaWeb/internal/repository"

func GetDealersList() ([]UserInfo, error) {
	dealers, err := repository.GetDealersList()
	if err != nil {
		return nil, err
	}
	var dealersInfo []UserInfo
	for _, dealer := range dealers {
		dealersInfo = append(dealersInfo, UserInfo{
			ID:       dealer.User.ID,
			Fullname: dealer.Company.Name,
			Phone:    dealer.User.PhoneNumber,
			Email:    dealer.User.Email,
			IsDealer: true,

			Dealer: &DealerInfo{
				CompanyName: dealer.Company.Name,
				Address:     *dealer.Address,
			},
		})

	}
	return dealersInfo, nil
}
