package utils

import pb "github.com/daniarmas/api_go/pkg"

func ParsePlatformType(tp *string) *pb.PlatformType {
	if *tp == "IOS" {
		return pb.PlatformType_IOS.Enum()
	} else if *tp == "Android" {
		return pb.PlatformType_Android.Enum()
	} else {
		return pb.PlatformType_PlatformTypeUnspecified.Enum()
	}
}

func ParseOrderStatusType(tp *string) *pb.OrderStatusType {
	switch *tp {
	case "OrderStatusTypePending":
		return pb.OrderStatusType_OrderStatusTypePending.Enum()
	case "OrderStatusTypeRejected":
		return pb.OrderStatusType_OrderStatusTypeRejected.Enum()
	case "OrderStatusTypeApproved":
		return pb.OrderStatusType_OrderStatusTypeApproved.Enum()
	case "OrderStatusTypeReceived":
		return pb.OrderStatusType_OrderStatusTypeReceived.Enum()
	case "OrderStatusTypeCanceled":
		return pb.OrderStatusType_OrderStatusTypeCanceled.Enum()
	default:
		return pb.OrderStatusType_OrderStatusTypeUnspecified.Enum()
	}
}

func ParseAppType(tp *string) *pb.AppType {
	if *tp == "App" {
		return pb.AppType_App.Enum()
	} else if *tp == "BusinessApp" {
		return pb.AppType_BusinessApp.Enum()
	} else {
		return pb.AppType_AppTypeUnspecified.Enum()
	}
}

func ParseItemStatusType(tp *string) *pb.ItemStatusType {
	if *tp == "Available" {
		return pb.ItemStatusType_Available.Enum()
	} else if *tp == "Unavailable" {
		return pb.ItemStatusType_Unavailable.Enum()
	} else if *tp == "Deprecated" {
		return pb.ItemStatusType_Deprecated.Enum()
	} else {
		return pb.ItemStatusType_ItemStatusTypeUnspecified.Enum()
	}
}

func ParseSearchMunicipalityType(tp string) *pb.SearchMunicipalityType {
	if tp == "More" {
		return pb.SearchMunicipalityType_More.Enum()
	} else if tp == "NoMore" {
		return pb.SearchMunicipalityType_NoMore.Enum()
	} else {
		return pb.SearchMunicipalityType_SearchMunicipalityTypeUnspecified.Enum()
	}
}

func ParseOrderResidenceType(tp *string) *pb.ResidenceType {
	switch *tp {
	case "ResidenceTypeHouse":
		return pb.ResidenceType_ResidenceTypeHouse.Enum()
	case "ResidenceTypeApartament":
		return pb.ResidenceType_ResidenceTypeApartment.Enum()
	default:
		return pb.ResidenceType_ResidenceTypeUnspecified.Enum()
	}
}

func ParseDeliveryType(tp *string) *pb.OrderType {
	switch *tp {
	case "DeliveryTypePickUp":
		return pb.OrderType_OrderTypePickUp.Enum()
	case "DeliveryTypeHomeDelivery":
		return pb.OrderType_OrderTypeHomeDelivery.Enum()
	default:
		return pb.OrderType_OrderTypeUnspecified.Enum()
	}
}

func ParseResidenceType(tp string) *pb.UserAddress_UserAddressType {
	if tp == "House" {
		return pb.UserAddress_HOUSE.Enum()
	} else if tp == "Apartament" {
		return pb.UserAddress_APARTAMENT.Enum()
	} else {
		return pb.UserAddress_UNSPECIFIED.Enum()
	}
}
